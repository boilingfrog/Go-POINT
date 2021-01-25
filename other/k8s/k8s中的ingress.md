## k8s中的ingress

### 什么是ingress

`Ingress`对象，其实就是对“反向代理”的一种抽象，简单的说就是一个全局的负载均衡器，可以通过访问URL定位到后端的`Service`。  

有了`Ingress`这个抽象，K8S就不需要关心`Ingress`的细节了，实际使用时，只需要选择一个具体的`Ingress Controller`部署就行了，业界常用的反向代理项目有：`Nginx、HAProxy、Envoy、Traefik`，都已经成为了K8S专门维护的`Ingress Controller`。  

`Service`是基于四层TCP和UDP协议转发的，而`Ingress`可以基于七层的HTTP和HTTPS协议转发，可以通过域名和路径做到更细粒度的划分。  

![channel](/img/ingress_7.jpg?raw=true)

### 理解Ingress 实现

k8s 有一个贯穿始终的设计理念，即需求和供给的分离。`Ingress Controller`和 `Ingress` 的实现也很好的实践了这一点。 要理解k8s ，时刻记住 需求供给分离的设计理念。

### Pod与Ingress的关系

通过Service关联Pod  
基于域名访问  
通过Ingress Controller实现Pod的负载均衡  
支持TCP/UDP 4层和HTTP 7层  

Ingress是一个规







### 参考
【Kubernetes的Ingress是啥】https://www.cnblogs.com/chenqionghe/p/11726231.html  
【理解k8s 的 Ingress】https://www.jianshu.com/p/189fab1845c5  
【Ingress】https://www.huaweicloud.com/zhishi/Ingress.html 

