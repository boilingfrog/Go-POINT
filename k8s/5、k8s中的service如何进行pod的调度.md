<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [k8s 中如何进行 Pod 的调度](#k8s-%E4%B8%AD%E5%A6%82%E4%BD%95%E8%BF%9B%E8%A1%8C-pod-%E7%9A%84%E8%B0%83%E5%BA%A6)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [endpoint](#endpoint)
  - [kube-proxy](#kube-proxy)
    - [userspace 模式](#userspace-%E6%A8%A1%E5%BC%8F)
    - [iptables](#iptables)
  - [负载均衡](#%E8%B4%9F%E8%BD%BD%E5%9D%87%E8%A1%A1)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## k8s 中如何进行 Pod 的调度

### 前言

Service 资源主要用于为 Pod 对象提供一个固定、统一的访问接口及负载均衡的能力。     

service 是一组具有相同 label pod 集合的抽象，集群内外的各个服务可以通过 service 进行互相通信。  

当创建一个 service 对象时也会对应创建一个 endpoint 对象，endpoint 是用来做容器发现的，service 只是将多个 pod 进行关联，实际的路由转发都是由 kubernetes 中的 kube-proxy 组件来实现，因此，service 必须结合 kube-proxy 使用，kube-proxy 组件可以运行在 kubernetes 集群中的每一个节点上也可以只运行在单独的几个节点上，其会根据 service 和 endpoints 的变动来改变节点上 iptables 或者 ipvs 中保存的路由规则。  

### endpoint

endpoint 是 k8s 集群中的一个资源对象，存储在 etcd 中，用来记录一个 service 对应的所有 pod 的访问地址。  

service 通过 selector 和 pod 建立关联。k8s 会根据 service 关联到 pod 的 podIP 信息组合成一个 endpoint。   

如果 service 没有 selector 字段，当一个 service 被创建的时候，`endpoint controller` 不会自动创建 endpoint。   

```
$ kubectl get svc -n study-k8s
NAME         TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
go-web-svc   ClusterIP   10.233.55.112   <none>        8000/TCP   9d

$ kubectl get endpoints -n study-k8s
NAME         ENDPOINTS                                                                AGE
go-web-svc   10.233.111.171:8000,10.233.111.172:8000,10.233.72.153:8000 + 2 more...   9d
```

栗如  

上面的 service `go-web-svc`，就有一个对应的 endpoint，ENDPOINTS 里面展示的就是 service 关联的 pod 的 ip 地址和端口。    

其中 `endpoint controller` 负载维护 endpoint 对象，主要的功能有下面几种  

1、负责生成和维护所有endpoint对象的控制器；  

2、负责监听 service 和对应 pod 的变化；  

3、监听到 service 被删除，则删除和该 service 同名的 endpoint 对象；  

4、监听到新的 service 被创建，则根据新建 service 信息获取相关 pod 列表，然后创建对应 endpoint 对象；  

5、监听到 service 被更新，则根据更新后的 service 信息获取相关 pod 列表，然后更新对应 endpoint 对象；  

6、监听到 pod 事件，则更新对应的 service 的 endpoint 对象，将 podIp 记录到 endpoint中；  

### kube-proxy  

kube-proxy 是一个简单的网络代理和负载均衡器，它的作用主要是负责 Service 的实现，具体来说，就是实现了内部从 Pod 到 Service 和外部的从 NodePort 向 Service 的访问，每台机器上都运行一个 kube-proxy 服务，它监听 API server 中 service 和 endpoint 的变化情况，并通过 iptables 等来为服务配置负载均衡（仅支持 TCP 和 UDP）。    

在 k8s 中提供相同服务的一组 pod 可以抽象成一个 service，通过 service 提供统一的服务对外提供服务，kube-proxy 存在于各个 node 节点上，负责为 service 提供 cluster 内部的服务发现和负载均衡，负责 Pod 的网络代理，它会定时从 etcd 中获取 service 信息来做相应的策略，维护网络规则和四层负载均衡工作。k8s 中集群内部的负载均衡就是由 kube-proxy 实现的，它是 k8s 中内部的负载均衡器，也是一个分布式代理服务器，可以在每个节点中部署一个，部署的节点越多，提供负载均衡能力的 Kube-proxy 就越多，高可用节点就越多。   

简单点讲就是 k8s 内部的 pod 要访问 service ，kube-proxy 会将请求转发到 service 所代表的一个具体 pod，也就是 service 关联的 Pod。  

同理对于外部访问 service 的请求，不论是 `Cluster IP+TargetPort` 的方式；还是用 Node 节点 `IP+NodePort` 的方式，都被 Node 节点的 Iptables 规则重定向到 Kube-proxy 监听 Service 服务代理端口。kube-proxy 接收到 Service 的访问请求后，根据负载策略，转发到后端的 Pod。   

kube-proxy 的路由转发规则是通过其后端的代理模块实现的，其中 kube-proxy 的代理模块目前有四种实现方案，userspace、iptables、ipvs、kernelspace 。     

#### userspace 模式  

userspace 模式在 `k8s v1.2` 后就已经被淘汰了，userspace 的作用就是在 proxy 的用户空间监听一个端口，所有的 svc 都转到这个端口，然后 proxy 内部应用层对其进行转发。proxy 会为每一个 svc 随机监听一个端口，并增加一个 iptables 规则。  

从客户端到 `ClusterIP:Port` 的报文都会通过 iptables 规则被重定向到 `Proxy Port`，Kube-Proxy 收到报文后，然后分发给对应的 Pod。  

<img src="/img/k8s/kube-proxy-userspace-mode.png"  alt="k8s" />   

userspace 模式下，流量的转发主要是在用户空间下完成的，上面提到了客户端的请求需要借助于 iptables 规则找到对应的 `Proxy Port`，因为 iptables 是在内核空间，这里就会请求就会有一次从用户态到内核态再返回到用户态的传递过程, 一定程度降低了服务性能。所以就会认为这种方式会有一定的性能损耗。  

#### iptables

首先来简单了解下 iptables：  

iptables 是 Linux 中最常用的一种防火墙工具，除了防火墙它还可以用作 IP 转发和简单的负载均衡功能。基于 Linux 中的 netfilter 内核模块实现。 Netfilter 在协议中添加了一些钩子，它允许内核模块通过这些钩子注册回调函数，这样经过钩子的所有数据都会被注册在响应钩子上的函数处理，包括修改数据包内容、给数据包打标记或者丢掉数据包等。iptables 是运行在用户态的一个程序，通过 netlink 和内核的 netfilter 框架打交道，具有足够的灵活性来处理各种常见的数据包操作和过滤需求。它允许将灵活的规则序列附加到内核的数据包处理管道中的各种钩子上。  

Netfilter 是 `Linux 2.4.x` 引入的一个子系统，它作为一个通用的、抽象的框架，提供一整套的 hook 函数的管理机制，使得诸如数据包过滤、网络地址转换(NAT)和基于协议类型的连接跟踪成为了可能。  

在 `kubernetes v1.2` 之后 iptables 成为默认代理模式，这种模式下，kube-proxy 会监视 `Kubernetes master` 对 Service 对象和 Endpoints 对象的添加和移除。 对每个 Service，它会安装 iptables 规则，从而捕获到达该 Service 的 clusterIP（虚拟 IP）和端口的请求，进而将请求重定向到 Service 的一组 backend 中的某个上面。因为流量转发都是在内核进行的，所以性能更高更加可靠。  

<img src="/img/k8s/services-iptables-overview.jpeg"  alt="k8s" />   

可以看到该模式下 iptables 来做用户态的入口，kube-proxy 只是持续监听 Service 以及 Endpoints 对象的变化， iptables 通过设置的转发策略，直接将对 VIP 的请求转发给后端 Pod，iptables 使用 DNAT 来完成转发，其采用了随机数实现负载均衡。  

该模式相比 userspace 模式，克服了请求在用户态-内核态反复传递的问题，性能上有所提升，但使用 iptables NAT 来完成转发，存在不可忽视的性能损耗，iptables 模式最主要的问题是在 service 数量大的时候会产生太多的 iptables 规则，使用非增量式更新会引入一定的时延，大规模情况下有明显的性能问题。  

### 负载均衡


### 参考

【kubernetes service 原理解析】https://zhuanlan.zhihu.com/p/111244353     
【service selector】https://blog.csdn.net/luanpeng825485697/article/details/84296765   
【一文看懂 Kube-proxy】https://zhuanlan.zhihu.com/p/337806843  
【Kubernetes 【网络组件】kube-proxy使用详解】https://blog.csdn.net/xixihahalelehehe/article/details/115370095     
【Service】https://jimmysong.io/kubernetes-handbook/concepts/service.html    


