<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [k8s 中如何进行 Pod 的调度](#k8s-%E4%B8%AD%E5%A6%82%E4%BD%95%E8%BF%9B%E8%A1%8C-pod-%E7%9A%84%E8%B0%83%E5%BA%A6)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [endpoint](#endpoint)
  - [服务发现](#%E6%9C%8D%E5%8A%A1%E5%8F%91%E7%8E%B0)
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


### 服务发现



### 负载均衡


### 参考

【kubernetes service 原理解析】https://zhuanlan.zhihu.com/p/111244353     
【service selector】https://blog.csdn.net/luanpeng825485697/article/details/84296765    


