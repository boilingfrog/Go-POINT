<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [helm使用](#helm%E4%BD%BF%E7%94%A8)
  - [什么是helm](#%E4%BB%80%E4%B9%88%E6%98%AFhelm)
  - [安装helm](#%E5%AE%89%E8%A3%85helm)
    - [Helm V2 & V3 架构设计](#helm-v2--v3-%E6%9E%B6%E6%9E%84%E8%AE%BE%E8%AE%A1)
    - [配置kube config](#%E9%85%8D%E7%BD%AEkube-config)
  - [helm使用](#helm%E4%BD%BF%E7%94%A8-1)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## helm使用

### 什么是helm 

`Helm`是`Deis`开发的一个用于`Kubernetes`应用的包管理工具，主要用来管理`Charts`。有点类似于`Ubuntu`中的`APT`或`CentOS`中的`YUM`。    

`Helm chart`是用来封装`Kubernetes`原生应用程序的`YAML`文件，可以在你部署应用的时候自定义应用程序的一些`metadata`，便与应用程序的分发。      

### 安装helm

#### Helm V2 & V3 架构设计

`Helm V2`到`V3`经历了较大的变革，其中最大的改动就是移除了`Tiller`组件，所有功能都通过`Helm CLI`与`ApiServer`直接交互。`Tiller`在`V2`的架构中扮演着重要的角色，但是它与`K8S`的设计理念是冲突的。  

- 1、围绕 Tiller 管理应用的生命周期不利于扩展。  

- 2、增加用户的使用壁垒，服务端需要部署 Tiller 组件才可以使用，侵入性强。  

- 3、K8S 的 RBAC 变得毫无用处，Tiller 拥有过大的 RBAC 权限存在安全风险。  

- 4、造成多租户场景下架构设计与实现复杂。    

<img src="/img/helm_1.jpeg" alt="helm" align=center />

所以对应v3版本的helm直接安装客户端就好了。   

关于helm版本的安装选择需要根据自己的k8s版本做选择，具体的版本策略[Helm版本支持策略](https://helm.sh/zh/docs/topics/version_skew/)  

具体的安装参考[安装Helm](https://helm.sh/zh/docs/intro/install/)就好了  

我的k8s版本是`v1.19.9`，所以本次我选择的helm版本是`v3.4.2`  

#### 配置kube config

1、配置`KUBECONFIG`变量  

```go
$ KUBECONFIG=~/.kube/config
```

2、完成`Kubeconfig`配置后，依次执行以下命令查看并切换`context`以访问本集群

````go
$ kubectl config get-contexts
$ kubectl config use-context cls-3jju4zdc-context-default
````

3、执行以下命令，测试是否可正常访问集群  

```go
$ kubectl get node
```

### helm使用











### 参考
【YAML 模版老去？Helm Chart 或将应用分发事实标准】https://www.infoq.cn/article/dwc0ipnguogq4kbap*9g  
【Helm V3使用指北】http://www.wangyapu.com/2020/04/10/helm_user_guide/  