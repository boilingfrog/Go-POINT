## k8s 中如何进行 Pod 的调度

### 前言

Service 资源主要用于为 Pod 对象提供一个固定、统一的访问接口及负载均衡的能力。     

service 是一组具有相同 label pod 集合的抽象，集群内外的各个服务可以通过 service 进行互相通信。  

当创建一个 service 对象时也会对应创建一个 endpoint 对象，endpoint 是用来做容器发现的，service 只是将多个 pod 进行关联，实际的路由转发都是由 kubernetes 中的 kube-proxy 组件来实现，因此，service 必须结合 kube-proxy 使用，kube-proxy 组件可以运行在 kubernetes 集群中的每一个节点上也可以只运行在单独的几个节点上，其会根据 service 和 endpoints 的变动来改变节点上 iptables 或者 ipvs 中保存的路由规则。

### 服务发现



### 负载均衡


### 参考

【kubernetes service 原理解析】https://zhuanlan.zhihu.com/p/111244353    


