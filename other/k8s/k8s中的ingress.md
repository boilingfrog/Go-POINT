## k8s中的ingress

### 什么是ingress

`Ingress`对象，其实就是对“反向代理”的一种抽象，简单的说就是一个全局的负载均衡器，可以通过访问URL定位到后端的`Service`。  

有了`Ingress`这个抽象，K8S就不需要关心`Ingress`的细节了，实际使用时，只需要选择一个具体的`Ingress Controller`部署就行了，业界常用的反向代理项目有：`Nginx、HAProxy、Envoy、Traefik`，都已经成为了K8S专门维护的`Ingress Controller`。  

`Service`是基于四层TCP和UDP协议转发的，而`Ingress`可以基于七层的HTTP和HTTPS协议转发，可以通过域名和路径做到更细粒度的划分。  

![channel](/img/ingress_7.jpg?raw=true)

### 理解Ingress 实现

k8s 有一个贯穿始终的设计理念，即需求和供给的分离。`Ingress Controller`和 `Ingress` 的实现也很好的实践了这一点。 要理解k8s ，时刻记住 需求供给分离的设计理念。

为使用Ingress，需要创建`Ingress Controller`（带一个默认backend服务）和Ingres s策略设置来共同完成。 

#### Ingress Controller

在定义Ingress策略之前，需要先部署`Ingress Controller`，以实现为所有后端Servi ce都提供一个统一的入口。`Ingress Controller`需要实现基于不同HTTP URL向后转发的负 载分发规则，并可以灵活设置7层负载分发策略。如果公有云服务商能够提供该类型的HTT P路由LoadBalancer，则也可设置其为`Ingress Controller`。    

在Kubernetes中，`Ingress Controller`将以Pod的形式运行，监控`API Server`的`/ingress`接口后端的`backend services`，如果Service发生变化，则`Ingress Controller`应自动 更新其转发规则。    

为了让`Ingress Controller`正常启动，还需要为它配置一个默认的backend，用于在 客户端访问的URL地址不存在时，返回一个正确的404应答。这个backend服务用任何应用 实现都可以，只要满足对根路径“/”的访问返回404应答，并且提供`/healthz`路径以使`kubelet`完成对它的健康检查。  

`注意事项`

1、一个集群中可以有多个 `Ingress Controller`， 在Ingress 中可以指定使用哪一个`Ingress Controller`；  
2、多个Ingress 规则可能出现竞争；   
3、`Ingress Controller` 本身需要以hostport 或者 service形式暴露出来。 云端可以使用云供应商lb 服务；    
4、Ingress 可以为多个命名空间服务。  

### 配置ingress规则

#### 转发到单个后端服务上

无需定义rule,直接指定到需要转发的service上就好了。  

```
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: web
spec:
  rules:
  - host: liz-test.com
    http:
      paths:
      - backend:
          serviceName: go-app-svc
          servicePort: 8000
```




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

