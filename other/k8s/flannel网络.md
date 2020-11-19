## 理解flannel网络

### 简介 

Flannel是CoreOS团队针对Kubernetes设计的一个网络规划服务，简单来说，它的功能是让集群中的不同节点主机创建的Docker容器都具有全集群唯一的虚拟IP地址。  

### Flannel

Kubernetes 对 Pod 之间如何进行组网通信提出了要求，Kubernetes 对集群网络有以下要求  

- 所有的 Pod 之间可以在不使用 NAT 网络地址转换的情况下相互通信；
- 所有的 Node 之间可以在不使用 NAT 网络地址转换的情况下相互通信；
- 每个 Pod 看到的自己的 IP 和其他 Pod 看到的一致。

Kubernetes 网络模型设计基础原则：  

每个pod都有一个独立的ip地址，而且假定





### 参考
【Kubernetes中的网络解析——以flannel为例】https://jimmysong.io/kubernetes-handbook/concepts/flannel.html  
【kubernetes网络模型之“小而美”flannel】https://zhuanlan.zhihu.com/p/79270447  
【Flannel网络原理】https://www.jianshu.com/p/165a256fb1da  