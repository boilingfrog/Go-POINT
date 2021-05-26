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
  - [helm的核心概念](#helm%E7%9A%84%E6%A0%B8%E5%BF%83%E6%A6%82%E5%BF%B5)
    - [Chart](#chart)
    - [Config](#config)
    - [Repository](#repository)
    - [Release](#release)
  - [基本使用](#%E5%9F%BA%E6%9C%AC%E4%BD%BF%E7%94%A8)
    - [chart的目录](#chart%E7%9A%84%E7%9B%AE%E5%BD%95)
    - [模板管理](#%E6%A8%A1%E6%9D%BF%E7%AE%A1%E7%90%86)
    - [模板部署](#%E6%A8%A1%E6%9D%BF%E9%83%A8%E7%BD%B2)
    - [卸载应用](#%E5%8D%B8%E8%BD%BD%E5%BA%94%E7%94%A8)
    - [自定义参数安装应用](#%E8%87%AA%E5%AE%9A%E4%B9%89%E5%8F%82%E6%95%B0%E5%AE%89%E8%A3%85%E5%BA%94%E7%94%A8)
  - [应用发布顺序依赖](#%E5%BA%94%E7%94%A8%E5%8F%91%E5%B8%83%E9%A1%BA%E5%BA%8F%E4%BE%9D%E8%B5%96)
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

### helm的核心概念

#### Chart

`Helm`采用`Chart`的格式来标准化描述一个应用（K8S 资源文件集合），`Chart`有自身标准的目录结构，可以将目录打包成版本化的压缩包进行部署。就像我们下载一个软件包之后，就可以在电脑上直接安装一样，同理`Chart`包可以通过`Helm`部署到任意的`K8S`集群中。  

#### Config

`Config`指应用配置参数，在`Chart`中由`values.yaml`和命令行参数组成。`Chart`采用`Go Template`的特性 + `values.yaml`对部署的模板文件进行参数渲染，也可以通过`Helm Client`的命令`–set key=value`的方式进行参数赋值。    

#### Repository

类似于`Docker Hub`，`Helm`官方、阿里云等社区都提供了`Helm Repository`，我们可以通过`helm repo add`导入仓库地址，便可以检索仓库并选择别人已经制作好的`Chart`包，开箱即用。  

#### Release

`Release`代表`Chart`在集群中的运行实例，同一个集群的同一个`Namespace`下`Release`名称是唯一的。`Helm`围绕`Release`对应用提供了强大的生命周期管理能力，包括`Release`的查询、安装、更新、删除、回滚等。    

### 基本使用

#### chart的目录

```go
chart-demo/
├── Chart.yaml
├── charts
├── templates
│   ├── NOTES.txt
│   ├── _helpers.tpl
│   ├── deployment.yaml
│   ├── hpa.yaml
│   ├── ingress.yaml
│   ├── service.yaml
│   ├── serviceaccount.yaml
│   └── tests
│       └── test-connection.yaml
└── values.yaml
```

#### 模板管理

创建 Chart 骨架

```go
helm create ./chart-demo
```

Chart 打包

````go
helm package ./chart-demo
````

获取 Chart 包元数据信息

````go
helm inspect chart ./chart-demo
````

本地渲染模板文件

```go
helm template ${chart-demo-release-name} --namespace ${namespace} ./chart-demo
```

查询 Chart 依赖信息

```go
helm dependency list ./chart-demo
```

检查依赖和模板配置是否正确

```go
$ helm lint chart-demo
==> Linting charts/chart-demo/
[INFO] Chart.yaml: icon is recommended

1 chart(s) linted, 0 chart(s) failed
```

#### 模板部署  

查询 Release 列表

```go
helm list --namespace xxxx
```

Chart 安装

````go
helm install ${chart-demo-release-name} ./chart-demo --namespace ${namespace}
````

Chart 版本升级

```go
helm upgrade ${chart-demo-release-name} ./chart-demo-new-version --namespace ${namespace}
```

Chart 版本回滚

```go
helm rollback ${chart-demo-release-name} ${revision} --namespace ${namespace}
```

查看 Release 历史版本

```go
helm history ${chart-demo-release-name} --namespace ${namespace}
```

#### 卸载应用

卸载应用，并保留安装记录

```go
helm uninstall ${chart-demo-release-name} -n ${namespace} --keep-history
```

查看全部应用（包含安装和卸载的应用）  

```go
helm list -n ${namespace} --all
```

卸载应用，不保留安装记录

```go
helm delete ${chart-demo-release-name} -n ${namespace}
```

**${}中的替换成自己的名字**

#### 自定义参数安装应用

`Helm`中支持使用自定义`yaml`文件和`--set`命令参数对要安装的应用进行参数配置，使用如下：  

方式一：使用自定义`values.yaml`文件安装应用

我们知道chart的目录结构中有一个`values.yaml`，里面就是用来放参数的配置文件，修改对应的`values.yaml`就可以了  

```go
// 展示对应配置参数信息
$ helm show values bitnami/nginx
image:
  registry: docker.io
  repository: bitnami/nginx
  tag: 1.19.10-debian-10-r14
...
```

方式二：使用`--set`配置参数进行安装  

`--set`参数是在使用`helm`命令时候添加的参数，可以在执行`helm`安装与更新应用时使用，多个参数间用,隔开，使用如下：  

注意：如果配置文件和`--set`同时使用，则`--set`设置的参数会覆盖配置文件中的参数配置。  

```go
// 使用set创建一个release
helm install --set 'registry.registry=docker.io,registry.repository=bitnami/nginx' nginx bitnami/nginx -n blog

// 更新一个release
helm upgrade --set 'servers[0].port=8080' nginx bitnami/nginx -n blog
```

### 应用发布顺序依赖  

虽然`Chart`可以通过`requirements.yaml`来管理依赖关系，并按照顺序下发模板资源，但是并无法控制子`Chart`之间的发布顺序。例如服务 B 部署必须依赖服务 A 的资源全部`Ready`。可以通过自定义子`Chart`之间的依赖顺序，在产品层控制每个子`Chart`的发布过程。  

### 参考
【YAML 模版老去？Helm Chart 或将应用分发事实标准】https://www.infoq.cn/article/dwc0ipnguogq4kbap*9g  
【Helm V3使用指北】http://www.wangyapu.com/2020/04/10/helm_user_guide/  
【Helm v3安装与应用】https://blog.51cto.com/wutengfei/2569465  