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

### helm使用











### 参考
【YAML 模版老去？Helm Chart 或将应用分发事实标准】https://www.infoq.cn/article/dwc0ipnguogq4kbap*9g  
【Helm V3使用指北】http://www.wangyapu.com/2020/04/10/helm_user_guide/  