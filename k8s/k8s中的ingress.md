<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [k8s中的ingress](#k8s%E4%B8%AD%E7%9A%84ingress)
  - [什么是ingress](#%E4%BB%80%E4%B9%88%E6%98%AFingress)
  - [理解Ingress 实现](#%E7%90%86%E8%A7%A3ingress-%E5%AE%9E%E7%8E%B0)
    - [Ingress Controller](#ingress-controller)
  - [配置ingress规则](#%E9%85%8D%E7%BD%AEingress%E8%A7%84%E5%88%99)
    - [转发到单个后端服务上](#%E8%BD%AC%E5%8F%91%E5%88%B0%E5%8D%95%E4%B8%AA%E5%90%8E%E7%AB%AF%E6%9C%8D%E5%8A%A1%E4%B8%8A)
    - [同一个域名，不同的URL路径被转发到不同的服务上](#%E5%90%8C%E4%B8%80%E4%B8%AA%E5%9F%9F%E5%90%8D%E4%B8%8D%E5%90%8C%E7%9A%84url%E8%B7%AF%E5%BE%84%E8%A2%AB%E8%BD%AC%E5%8F%91%E5%88%B0%E4%B8%8D%E5%90%8C%E7%9A%84%E6%9C%8D%E5%8A%A1%E4%B8%8A)
    - [不同的域名（虚拟主机名）被转发到不同的服务上](#%E4%B8%8D%E5%90%8C%E7%9A%84%E5%9F%9F%E5%90%8D%E8%99%9A%E6%8B%9F%E4%B8%BB%E6%9C%BA%E5%90%8D%E8%A2%AB%E8%BD%AC%E5%8F%91%E5%88%B0%E4%B8%8D%E5%90%8C%E7%9A%84%E6%9C%8D%E5%8A%A1%E4%B8%8A)
    - [不使用域名的转发规则](#%E4%B8%8D%E4%BD%BF%E7%94%A8%E5%9F%9F%E5%90%8D%E7%9A%84%E8%BD%AC%E5%8F%91%E8%A7%84%E5%88%99)
  - [四层、七层负载均衡的区别](#%E5%9B%9B%E5%B1%82%E4%B8%83%E5%B1%82%E8%B4%9F%E8%BD%BD%E5%9D%87%E8%A1%A1%E7%9A%84%E5%8C%BA%E5%88%AB)
    - [四层负载](#%E5%9B%9B%E5%B1%82%E8%B4%9F%E8%BD%BD)
    - [七层负载](#%E4%B8%83%E5%B1%82%E8%B4%9F%E8%BD%BD)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## k8s中的ingress

### 什么是ingress

`Ingress`对象，其实就是对“反向代理”的一种抽象，简单的说就是一个全局的负载均衡器，可以通过访问URL定位到后端的`Service`。  

有了`Ingress`这个抽象，K8S就不需要关心`Ingress`的细节了，实际使用时，只需要选择一个具体的`Ingress Controller`部署就行了，业界常用的反向代理项目有：`Nginx、HAProxy、Envoy、Traefik`，都已经成为了K8S专门维护的`Ingress Controller`。  

`Service`是基于四层TCP和UDP协议转发的，而`Ingress`可以基于七层的HTTP和HTTPS协议转发，可以通过域名和路径做到更细粒度的划分。  

![ingress](/img/k8s/ingress_7.jpg?raw=true)

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

关于ingress的部署可以参考另一篇文章[k8s发布go应用](https://www.cnblogs.com/ricklz/p/14071965.html)  

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

这样定义就能将`liz-test.com`转发到，集群的`go-app-svc`service的8000端口上。   

#### 同一个域名，不同的URL路径被转发到不同的服务上

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
      - path: /test1
        backend:
          serviceName: go-app-svc-1
          servicePort: 8001
      - path: /test2
        backend:
          serviceName: go-app-svc-2
          servicePort: 8002
```

通过设置不同的path，将同一个host中不同的路径转发到不同的service中。   

#### 不同的域名（虚拟主机名）被转发到不同的服务上

```
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: web
spec:
  rules:
  - host: liz-test-1.com
    http:
      paths:
      - backend:
          serviceName: go-app-svc-1
          servicePort: 8001
  - host: liz-test-2.com
    http:
      paths:
      - backend:
          serviceName: go-app-svc-2
          servicePort: 8002
```

设置不同的host，然后转发不同的host到不同的service中。   

#### 不使用域名的转发规则

这种配置用于一个网站不使用域名直接提供服务的场景，此时通过任意一台运行`ingress-controller`的Node都能访问到后端的服务。  

```
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: web-1
spec:
  rules:
   - http:
      paths:
      - path: /test1
        backend:
          serviceName: go-app-svc
          servicePort: 8000
```

访问`<ingress-controller-ip>/demo`即可  

注意，使用无域名的Ingress转发规则时，将默认禁用非安全HTTP，强制启用HTTPS。所以必须使用https方式访问。   

```
$ curl -k https://192.168.56.202/test1
<h1>hello world</h1><div>你好</div>
```

可以在Ingress的定义中设置一个annotation“ingress.kubernetes.io/ssl-redirect=false”来关闭强制启用HTTPS的设置

```
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: web-2
  annotations:
    ingress.kubernetes.io/ssl-redirect: "false"
spec:
  rules:
   - http:
      paths:
      - path: /test1
        backend:
          serviceName: go-app-svc
          servicePort: 8000
```

### 四层、七层负载均衡的区别

上面我们经常提到四层和七层的网络负载，不过到底是什么呢，已经他们之间的区别？  

更详细的可查看文章[四层、七层负载均衡的区别](https://www.jianshu.com/p/fa937b8e6712)  

#### 四层负载  

四层就是基于IP+端口的负载均衡  

所谓四层负载均衡，也就是主要通过报文中的目标地址和端口，再加上负载均衡设备设置的服务器选择方式，决定最终选择的内部服务器。  

以常见的TCP为例，负载均衡设备在接收到第一个来自客户端的SYN 请求时，即通过上述方式选择一个最佳的服务器，并对报文中目标IP地址进行修改(改为后端服务器IP），直接转发给该服务器。TCP的连接建立，即三次握手是客户端和服务器直接建立的，负载均衡设备只是起到一个类似路由器的转发动作。在某些部署情况下，为保证服务器回包可以正确返回给负载均衡设备，在转发报文的同时可能还会对报文原来的源地址进行修改。   

#### 七层负载  

七层就是基于URL等应用层信息的负载均衡  

所谓七层负载均衡，也称为“内容交换”，也就是主要通过报文中的真正有意义的应用层内容，再加上负载均衡设备设置的服务器选择方式，决定最终选择的内部服务器。  

以常见的TCP为例，负载均衡设备如果要根据真正的应用层内容再选择服务器，只能先代理最终的服务器和客户端建立连接(三次握手)后，才可能接受到客户端发送的真正应用层内容的报文，然后再根据该报文中的特定字段，再加上负载均衡设备设置的服务器选择方式，决定最终选择的内部服务器。负载均衡设备在这种情况下，更类似于一个代理服务器。负载均衡和前端的客户端以及后端的服务器会分别建立TCP连接。所以从这个技术原理上来看，七层负载均衡明显的对负载均衡设备的要求更高，处理七层的能力也必然会低于四层模式的部署方式。  

![ingress](/img/k8s/ingress_1.webp?raw=true)

应用场景  

七层应用负载的好处，是使得整个网络更智能化。例如访问一个网站的用户流量，可以通过七层的方式，将对图片类的请求转发到特定的图片服务器并可以使用缓存技术；将对文字类的请求可以转发到特定的文字服务器并可以使用压缩技术。当然这只是七层应用的一个小案例，从技术原理上，这种方式可以对客户端的请求和服务器的响应进行任意意义上的修改，极大的提升了应用系统在网络层的灵活性。很多在后台，例如Nginx或者Apache上部署的功能可以前移到负载均衡设备上，例如客户请求中的Header重写，服务器响应中的关键字过滤或者内容插入等功能。  

另外一个常常被提到功能就是安全性。网络中最常见的SYN Flood攻击，即黑客控制众多源客户端，使用虚假IP地址对同一目标发送SYN攻击，通常这种攻击会大量发送SYN报文，耗尽服务器上的相关资源，以达到Denial of Service(DoS)的目的。从技术原理上也可以看出，四层模式下这些SYN攻击都会被转发到后端的服务器上；而七层模式下这些SYN攻击自然在负载均衡设备上就截止，不会影响后台服务器的正常运营。另外负载均衡设备可以在七层层面设定多种策略，过滤特定报文，例如SQL Injection等应用层面的特定攻击手段，从应用层面进一步提高系统整体安全。  

现在的7层负载均衡，主要还是着重于应用HTTP协议，所以其应用范围主要是众多的网站或者内部信息平台等基于B/S开发的系统。 4层负载均衡则对应其他TCP应用，例如基于C/S开发的ERP等系统。  

### 参考
【Kubernetes的Ingress是啥】https://www.cnblogs.com/chenqionghe/p/11726231.html  
【理解k8s 的 Ingress】https://www.jianshu.com/p/189fab1845c5  
【Ingress】https://www.huaweicloud.com/zhishi/Ingress.html   
【四层、七层负载均衡的区别】https://cloud.tencent.com/developer/article/1082047  
【四层、七层负载均衡的区别】https://www.jianshu.com/p/fa937b8e6712  

