## k8s中的service

### service存在的意义  
防止Pod失联（服务发现）  
定义一组Pod的访问策略（负载均衡）  

对于kubernetes整个集群来说，Pod的地址也可变的，也就是说如果一个Pod因为某些原因退出了，而由于其设置了副本数replicas大于1，那么该Pod就会在集
群的任意节点重新启动，这个重新启动的Pod的IP地址与原IP地址不同，这对于业务来说，就不能根据Pod的IP作为业务调度。kubernetes就引入了Service的概
念，它为Pod提供一个入口，主要通过Labels标签来选择后端Pod，这时候不论后端Pod的IP地址如何变更，只要Pod的Labels标签没变，那么 业务通过service
调度就不会存在问题。  

当声明Service的时候，会自动生成一个cluster IP，这个IP是虚拟IP。我们就可以通过这个IP来访问后端的Pod，当然，如果集群配置了DNS服务，比如现在
的CoreDNS，那么也可以通过Service的名字来访问，它会通过DNS自动解析Service的IP地址。  

### Pod与Service的关系

通过label-selector相关联  
通过Service实现Pod的负载均衡（ TCP/UDP 4层）  

### Service几种类型

- ClusterIP：分配一个内部集群IP地址，只能在集群内部访问（同Namespace内的Pod），默认ServiceType。ClusterIP 模式的 Service 为你提供的，就是一个 Pod 的稳定的 IP 地址，即 VIP。  
- NodePort：就是Node的基本port。选择该值，这个servce就可以通过`NodeIP:NodePort`访问这个Service服务，NodePort会路由到Cluster IP服务，这个Cluster IP会通过请求自动创建；  
- LoadBalance：使用云提供商的负载均衡器，可以向外部暴露服务，选择该值，外部的负载均衡器可以路由到NodePort服务和Cluster IP服务；  
- ExternalName:通过返回 CNAME 和它的值，可以将服务映射到 externalName 字段的内容，没有任何类型代理被创建，可以用于访问集群内其他没有Labels的Pod，也可以访问其他NameSpace里的Service。  

### 参考

【kubernetes中常用对象service的详细介绍】https://zhuanlan.zhihu.com/p/103413341   
【Istio 运维实战系列（2）：让人头大的『无头服务』-上】https://cloud.tencent.com/developer/article/1700748  
【如何为服务网格选择入口网关？- Kubernetes Ingress, Istio Gateway还是API Gateway？】https://mp.weixin.qq.com/s?__biz=MzU3MjI5ODgxMA==&mid=2247483759&idx=1&sn=d44c3194810c02eba81d427292fab2d9&chksm=fcd2423acba5cb2c55e3d9952a74d06e6e5803e8a9755885e66630092e1687b57ad09154dc5f&scene=21#wechat_redirect  
