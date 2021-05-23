<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [helm使用](#helm%E4%BD%BF%E7%94%A8)
  - [什么是helm](#%E4%BB%80%E4%B9%88%E6%98%AFhelm)
  - [安装helm](#%E5%AE%89%E8%A3%85helm)
    - [Helm V2 & V3 架构设计](#helm-v2--v3-%E6%9E%B6%E6%9E%84%E8%AE%BE%E8%AE%A1)
    - [配置kube config](#%E9%85%8D%E7%BD%AEkube-config)
  - [helm使用](#helm%E4%BD%BF%E7%94%A8-1)
    - [添加仓库](#%E6%B7%BB%E5%8A%A0%E4%BB%93%E5%BA%93)
    - [helm安装nginx](#helm%E5%AE%89%E8%A3%85nginx)
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

需要借助于`kubectl`中的`Kubeconfig`，来实现对远端`k8s`集群的访问操作  

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

上面`helm`安装好之后，以及配置好本地的`kubectl`之后，我们就可以对远端`k8s`集群，使用`helm`进行操作了  

查看版本

```go
$ helm version
version.BuildInfo{Version:"v3.4.2", GitCommit:"23dd3af5e19a02d4f4baa5b2f242645a1a3af629", GitTreeState:"clean", GoVersion:"go1.14.13"}
```

#### 添加仓库

在`Helm`中默认是不会添加`Chart`仓库，所以这里我们需要手动添加，下面是添加一些常用的`Charts`库，命令如下：   

```go
helm repo add  elastic    https://helm.elastic.co       
helm repo add  gitlab     https://charts.gitlab.io       
helm repo add  harbor     https://helm.goharbor.io       
helm repo add  bitnami    https://charts.bitnami.com/bitnami       
helm repo add  incubator  https://kubernetes-charts-incubator.storage.googleapis.com       
helm repo add  stable     https://kubernetes-charts.storage.googleapis.com       

// 添加国内仓库       
helm repo add stable http://mirror.azure.cn/kubernetes/charts       
helm repo add aliyun https://kubernetes.oss-cn-hangzhou.aliyuncs.com/charts     
  
// 执行更新命令，将仓库中的信息进行同步：
helm repo update       

// 查看仓库信息
helm repo list   
```

#### helm安装nginx

通过`Helm`在`Repo`中查询可安装的`Nginx`包  

```go
$ helm search repo nginx
NAME                            	CHART VERSION	APP VERSION	DESCRIPTION                                       
bitnami/nginx                   	8.9.0        	1.19.10    	Chart for the nginx server                        
bitnami/nginx-ingress-controller	7.6.6        	0.46.0     	Chart for the nginx Ingress controller            
stable/nginx-ingress            	1.41.2       	v0.34.1    	An nginx Ingress controller that uses ConfigMap...
stable/nginx-ldapauth-proxy     	0.1.6        	1.13.5     	DEPRECATED - nginx proxy with ldapauth            
stable/nginx-lego               	0.3.1        	           	Chart for nginx-ingress-controller and kube-lego  
bitnami/kong                    	3.7.3        	2.4.1      	Kong is a scalable, open source API layer (aka ...
```

创建`Namespace`并且部署应用  

```go
//  创建命名空间test
$ kubectl create namespace test

// 查看创建的命名空间
$ kubectl get ns

// 选择一个chart在k8s上部署我们的应用
$ helm install nginx bitnami/nginx -n test

// 查看应用状态
$ helm status nginx -n test
$ helm list -n test

// 查看pod的状态
$ kubectl get pods -n test
NAME                     READY   STATUS    RESTARTS   AGE
nginx-7b9d7c59ff-69mgz   1/1     Running   0          2m39s
```

查看部署的结果

```go
$ kubectl get svc -n test
NAME    TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
nginx   LoadBalancer   xx.xxx.xxx.xx   <pending>     80:31998/TCP   6m3s
```

访问结果

<img src="/img/helm_2.jpg" alt="helm" align=center />










### 参考
【YAML 模版老去？Helm Chart 或将应用分发事实标准】https://www.infoq.cn/article/dwc0ipnguogq4kbap*9g  
【Helm V3使用指北】http://www.wangyapu.com/2020/04/10/helm_user_guide/  
【Helm v3安装与应用】https://blog.51cto.com/wutengfei/2569465  