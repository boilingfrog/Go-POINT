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
  - [发布应用](#%E5%8F%91%E5%B8%83%E5%BA%94%E7%94%A8)
  - [遇到的问题](#%E9%81%87%E5%88%B0%E7%9A%84%E9%97%AE%E9%A2%98)
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

```shell script
$ KUBECONFIG=~/.kube/config
```

2、完成`Kubeconfig`配置后，依次执行以下命令查看并切换`context`以访问本集群

````shell script
$ kubectl config get-contexts
$ kubectl config use-context cls-3jju4zdc-context-default
````

3、执行以下命令，测试是否可正常访问集群  

```shell script
$ kubectl get node
```

### helm使用

上面`helm`安装好之后，以及配置好本地的`kubectl`之后，我们就可以对远端`k8s`集群，使用`helm`进行操作了  

查看版本

```shell script
$ helm version
version.BuildInfo{Version:"v3.4.2", GitCommit:"23dd3af5e19a02d4f4baa5b2f242645a1a3af629", GitTreeState:"clean", GoVersion:"go1.14.13"}
```

#### 添加仓库

在`Helm`中默认是不会添加`Chart`仓库，所以这里我们需要手动添加，下面是添加一些常用的`Charts`库，命令如下：   

```shell script
helm repo add  elastic    https://helm.elastic.co       
helm repo add  gitlab     https://charts.gitlab.io       
helm repo add  harbor     https://helm.goharbor.io       
helm repo add  bitnami    https://charts.bitnami.com/bitnami       
helm repo add  incubator  https://kubernetes-charts-incubator.storage.googleapis.com       
helm repo add  stable     https://kubernetes-charts.storage.googleapis.com       

# 添加国内仓库       
helm repo add stable http://mirror.azure.cn/kubernetes/charts       
helm repo add aliyun https://kubernetes.oss-cn-hangzhou.aliyuncs.com/charts     
  
# 执行更新命令，将仓库中的信息进行同步：
helm repo update       

# 查看仓库信息
helm repo list   
```

#### helm安装nginx

通过`Helm`在`Repo`中查询可安装的`Nginx`包  

```shell script
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

```shell script
#  创建命名空间test
$ kubectl create namespace test

# 查看创建的命名空间
$ kubectl get ns

# 选择一个chart在k8s上部署我们的应用
$ helm install nginx bitnami/nginx -n test

# 查看应用状态
$ helm status nginx -n test
$ helm list -n test

# 查看pod的状态
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

```shell script
chart-demo/
├── Chart.yaml # chart原数据信息
├── charts # 应用依赖集合
├── templates # k8s资源模板集合
│   ├── NOTES.txt
│   ├── _helpers.tpl
│   ├── deployment.yaml
│   ├── hpa.yaml
│   ├── ingress.yaml
│   ├── service.yaml
│   ├── serviceaccount.yaml
│   └── tests
│       └── test-connection.yaml
└── values.yaml # 资源配置文件
```

#### 模板管理

创建 Chart 骨架

```shell script
helm create ./chart-demo
```

Chart 打包

````shell script
helm package ./chart-demo
````

获取 Chart 包元数据信息

````shell script
helm inspect chart ./chart-demo
````

本地渲染模板文件

```shell script
helm template ${chart-demo-release-name} --namespace ${namespace} ./chart-demo
```

查询 Chart 依赖信息

```shell script
helm dependency list ./chart-demo
```

检查依赖和模板配置是否正确

```shell script
$ helm lint chart-demo
==> Linting charts/chart-demo/
[INFO] Chart.yaml: icon is recommended

1 chart(s) linted, 0 chart(s) failed
```

#### 模板部署  

查询 Release 列表

```shell script
helm list --namespace xxxx
```

Chart 安装

````shell script
helm install ${chart-demo-release-name} ./chart-demo --namespace ${namespace}
````

Chart 版本升级

```shell script
helm upgrade ${chart-demo-release-name} ./chart-demo-new-version --namespace ${namespace}
```

Chart 版本回滚

```shell script
helm rollback ${chart-demo-release-name} ${revision} --namespace ${namespace}
```

查看 Release 历史版本

```shell script
helm history ${chart-demo-release-name} --namespace ${namespace}
```

#### 卸载应用

卸载应用，并保留安装记录

```shell script
helm uninstall ${chart-demo-release-name} -n ${namespace} --keep-history
```

查看全部应用（包含安装和卸载的应用）  

```shell script
helm list -n ${namespace} --all
```

卸载应用，不保留安装记录

```shell script
helm delete ${chart-demo-release-name} -n ${namespace}
```

**${}中的替换成自己的名字**

#### 自定义参数安装应用

`Helm`中支持使用自定义`yaml`文件和`--set`命令参数对要安装的应用进行参数配置，使用如下：  

方式一：使用自定义`values.yaml`文件安装应用

我们知道chart的目录结构中有一个`values.yaml`，里面就是用来放参数的配置文件，修改对应的`values.yaml`就可以了  

```shell script
# 展示对应配置参数信息
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

```shell script
# 使用set创建一个release
helm install --set 'registry.registry=docker.io,registry.repository=bitnami/nginx' nginx bitnami/nginx -n blog

# 更新一个release
helm upgrade --set 'servers[0].port=8080' nginx bitnami/nginx -n blog
```

### 应用发布顺序依赖  

虽然`Chart`可以通过`requirements.yaml`来管理依赖关系，并按照顺序下发模板资源，但是并无法控制子`Chart`之间的发布顺序。例如服务 B 部署必须依赖服务 A 的资源全部`Ready`。可以通过自定义子`Chart`之间的依赖顺序，在产品层控制每个子`Chart`的发布过程。  

### 发布应用

上面我们使用helm初始化了一个`charts`结构，我们使用上面初始化的好的结构，把我之前打包的一个镜像，发布到k8s环境中  

镜像`liz2019/main-test:1.1.72`

开始部署，为了方便查看结果，我们使用`NodePort`的类型部署，修改`values.yaml`  

```
service:
#  type: ClusterIP
#  port: 80
  type: NodePort
  port: 80
```

部署

```
$ helm upgrade --install --force --wait --namespace test  --set image.repository=liz2019/main-test --set image.tag=1.1.72 chart-demo  ./chart-demo

Release "chart-demo" does not exist. Installing it now.
NAME: chart-demo
LAST DEPLOYED: Tue Jun 15 16:27:22 2021
NAMESPACE: test
STATUS: deployed
REVISION: 1
NOTES:
1. Get the application URL by running these commands:
  export NODE_PORT=$(kubectl get --namespace test -o jsonpath="{.spec.ports[0].nodePort}" services chart-demo)
  export NODE_IP=$(kubectl get nodes --namespace test -o jsonpath="{.items[0].status.addresses[0].address}")
  echo http://$NODE_IP:$NODE_PORT
```

<img src="/img/helm_6.jpg" alt="helm" align=center />

查看下结果

<img src="/img/helm_7.jpg" alt="helm" align=center />

可以看到已经成功部署了

### 遇到的问题

```
UPGRADE FAILED: failed to replace object: Service "chart-demo" is invalid: spec.clusterIP: Invalid value: "": field is immutable
```

这种尝试设置不可变字段为空值的`Chart`，在`Helm 2`中可能正常工作，因为`Helm 2`仅仅比对新旧配置。只要你没有后续设置`clusterIP`为非空就不会出问题。但是到了`Helm 3`中则肯定会出现问题。  

**强制更新的行为**

指定`--force`标记时，可以在升级时，进行必要的强制更新：  

- 对于Helm 2，当PATCH操作失败时，会删除、再重建目标资源

- 对于Helm 3，会用PUT操作来替换（replace/overwrite）目标资源

如果三方合并出现问题，有可能通过强制更新解决。对于`Helm 3`来说，更多情况下无济于事，主要是K8S限制某些字段一旦创建即不可变更。

在Helm 3中，即使强制更新，你也可能遇到类似下面的错误：

1、`ersistentVolumeClaim "ng" is invalid: spec: Forbidden: is immutable after creation except resources.requests for bound claims`

2、`failed to replace object: Service "ng" is invalid: spec.clusterIP: Invalid value: "": field is immutable`

`PUT`操作解决不了不可变字段的问题，然而`Helm 2`删除后再创建，则规避了不可变字段问题，但会引发其它问题：

- PVC删除，PV级联删除么？数据怎么办

- Service删除，会导致暂时的服务不可用么？

### 参考
【YAML 模版老去？Helm Chart 或将应用分发事实标准】https://www.infoq.cn/article/dwc0ipnguogq4kbap*9g  
【Helm V3使用指北】http://www.wangyapu.com/2020/04/10/helm_user_guide/  
【Helm v3安装与应用】https://blog.51cto.com/wutengfei/2569465  
【基于Helm的Kubernetes资源管理】https://blog.gmem.cc/helm    