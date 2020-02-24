## k8s中的service


### service存在的意义  
防止Pod失联（服务发现）  
定义一组Pod的访问策略（负载均衡）

### Pod与Service的关系

通过label-selector相关联  
通过Service实现Pod的负载均衡（ TCP/UDP 4层）  

### Service几种类型

ClusterIP：分配一个内部集群IP地址，只能在集群内部访问（同Namespace内的Pod），默认ServiceType。  
ClusterIP 模式的 Service 为你提供的，就是一个 Pod 的稳定的 IP 地址，即 VIP。  
NodePort：分配一个内部集群IP地址，并在每个节点上启用一个端口来暴露服务，可以在集群外部访问。  
访问地址：<NodeIP>:<NodePort>  
LoadBalancer：分配一个内部集群IP地址，并在每个节点上启用一个端口来暴露服务。  
除此之外，Kubernetes会请求底层云平台上的负载均衡器，将每个Node（[NodeIP]:[NodePort]）作为后端添加进去。  
