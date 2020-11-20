<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [理解flannel网络](#%E7%90%86%E8%A7%A3flannel%E7%BD%91%E7%BB%9C)
  - [简介](#%E7%AE%80%E4%BB%8B)
  - [Kubernetes中的网络](#kubernetes%E4%B8%AD%E7%9A%84%E7%BD%91%E7%BB%9C)
  - [flannel](#flannel)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 理解flannel网络

### 简介 

Flannel是CoreOS团队针对Kubernetes设计的一个网络规划服务，简单来说，它的功能是让集群中的不同节点主机创建的Docker容器都具有全集群唯一的虚拟IP地址。  

### Kubernetes中的网络

Kubernetes 对 Pod 之间如何进行组网通信提出了要求，Kubernetes 对集群网络有以下要求  

- 所有的 Pod 之间可以在不使用 NAT 网络地址转换的情况下相互通信；
- 所有的 Node 之间可以在不使用 NAT 网络地址转换的情况下相互通信；
- 每个 Pod 看到的自己的 IP 和其他 Pod 看到的一致。

Kubernetes 网络模型设计基础原则：  

每个pod都有一个独立的ip地址，而且假定所有 Pod 都在一个可以直接连通的、扁平的网络空间中。所以不管它们是否运行在同一个 Node (宿主机) 中，都要
求它们可以直接通过对方的 IP 进行访问。设计这个原则的原因是，用户不需要额外考虑如何建立 Pod 之间的连接，也不需要考虑将容器端口映射到主机端口等问题。  

由于 Kubernetes 的网络模型是假设 Pod 之间访问时使用的是对方 Pod 的实际地址，所以一个 Pod 内部的应用程序看到自己的 IP 地址和端口与集群内其他
 Pod 看到的一样。它们都是 Pod 实际分配的 IP 地址 。这个 IP 地址和端口在 Pod 内部和外部都保持一致，我们可以不使用 NAT 来进行转换。  

我们知道Kubernetes集群内部存在三类IP，分别是：  

- Node IP：宿主机的IP地址
- Pod IP：使用网络插件创建的IP（如flannel），使跨主机的Pod可以互通
- Cluster IP：虚拟IP，通过iptables规则访问服务

在安装node节点的时候，节点上的进程是按照flannel -> docker -> kubelet -> kube-proxy的顺序启动的  

### flannel

Flannel是作为一个二进制文件的方式部署在每个node上，主要实现两个功能：  

- 为每个node分配subnet，容器将自动从该子网中获取IP地址
- 当有node加入到网络中时，为每个node增加路由配置

他的特点主要以下几点  

- 使集群中的不同 Node 主机创建的 Docker 容器都具有全集群唯一的虚拟 IP 地址；
- 创建一个新的虚拟网卡 flannel0 接收 docker 网桥的数据，通过维护路由表，对接收到的数据进行封包和转发（VXLAN； 
- 路由信息一般存放到 etcd 中：多个 Node 上的 Flanneld 依赖一个 etcd cluster 来做集中配置服务，etcd 保证了所有 Node 上 Flannel 所看到的配置是一致的。同时每个 Node 上的 Flannel 都可以监听 etcd 上的数据变化，实时感知集群中 Node 的变化；
- Flannel 首先会在 Node 上创建一个名为 flannel0 的网桥（VXLAN 类型的设备），并且在每个 Node 上运行一个名为 Flanneld 的代理。每个 Node 上的 Flannel 代理会从 etcd 上为当前 Node 申请一个 CIDR 地址块用来给该 Node 上的 Pod 分配地址；  
- Flannel 致力于给 Kubernetes 集群中的 Node 提供一个三层网络，它并不控制 Node 中的容器是如何进行组网的，仅仅关心流量如何在 Node 之间流转。  
- 建立一个覆盖网络（overlay network），这个覆盖网络会将数据包原封不动的传递到目标容器中。覆盖网络是建立在另一个网络之上并由其基础设施支持的虚拟网络。覆盖网络通过将一个分组封装在另一个分组内来将网络服务与底层基础设施分离。在将封装的数据包转发到端点后，将其解封装；  

![channel](/img/k8s_flannel_1.png?raw=true)







### 参考
【Kubernetes中的网络解析——以flannel为例】https://jimmysong.io/kubernetes-handbook/concepts/flannel.html  
【kubernetes网络模型之“小而美”flannel】https://zhuanlan.zhihu.com/p/79270447  
【Flannel网络原理】https://www.jianshu.com/p/165a256fb1da  