<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [声明式API](#%E5%A3%B0%E6%98%8E%E5%BC%8Fapi)
  - [声明式和命令式的对比](#%E5%A3%B0%E6%98%8E%E5%BC%8F%E5%92%8C%E5%91%BD%E4%BB%A4%E5%BC%8F%E7%9A%84%E5%AF%B9%E6%AF%94)
  - [Kubernetes 声明式 API 的工作原理](#kubernetes-%E5%A3%B0%E6%98%8E%E5%BC%8F-api-%E7%9A%84%E5%B7%A5%E4%BD%9C%E5%8E%9F%E7%90%86)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 声明式API

### 声明式和命令式的对比

命令式  

命令式有时也称为指令式，命令式的场景下，计算机只会机械的完成指定的命令操作，执行的结果就取决于执行的命令是否正确。    

声明式  

声明式也称为描述式或者申明式，这种方式告诉计算机想要的，由计算机自己去设定执行的路径，需要计算机有一定的`智能`。    

最常见的声明式栗子就是数据库，查询的 sql 就表示我们想要的结果集，数据库运行查询 sql 的时候，会帮我们处理查询，并且返回查询的结果。数据库在查询的时候，会进行索引匹配，做查询优化等处理，再返回数据结果的时候，同时使用最优的查询路径。如果我们自己去处理这些操作，就需要写很多代码了，而不是仅仅通过一行代码就能解决。    

在运维中，命令式的思路弊端就很多了  

1、编写脚本的难度很大，需要想明白每一步的操作，处理每一步可能遇到的各种情况，高度依赖编写者的经验；  

2、脚本依赖部署的环境，不同的环境可能执行的结果就是不一样的,同时依赖稳定的环境，在分布式环境下，不太容易实现；     

3、运维的流程很难具备事务性，如果脚本的执行被以外打断，会产生一些意想不到的中间状态；  

4、需要另外编写运维的文档，如果多人维护，协作就很困难。   

声明式的思路就能解决上面的这些情况了  

声明式的运维表现，就是编写一个配置文件，描述想要的部署结果，然后平台去解析这个配置文件生成部署的结果。部署的配置文件比过程化的脚本更容易理解，也更方便开发人员编写。   

配置文件更容易审错和多人维护，并且如果当前的部署失败，我们只需要更改部署文件为之前的版本，重新部署就能回退到之前的版本。   

Kubernetes 的设计用的就是声明式的设计思想。  

Kubernetes 不仅仅是一个编排系统，实际上它消除了编排的需要。 编排的技术定义是执行已定义的工作流程：首先执行 A，然后执行 B，再执行 C。 而 Kubernetes 包含了一组独立可组合的控制过程，可以连续地将当前状态驱动到所提供的预期状态。 你不需要在乎如何从 A 移动到 C，也不需要集中控制，这使得系统更易于使用且功能更强大、 系统更健壮，更为弹性和可扩展。这正是声明式的设计思想的体现。  

###  Kubernetes 声明式 API 的工作原理

当一个 YAML 文件提交给 Kubernetes 之后，它究竟是如何创建出一个 API 对象的呢？  

在 Kubernetes 项目中，一个 API 对象在 Etcd 里的完整资源路径，是由：Group（API 组）、Version（API 版本）和 Resource（API 资源类型）三个部分组成的。   

<img src="/img/k8s/k8s-api.webp"  alt="k8s" />    

可以看到 Kubernetes 中 API 对象组织方式，其实是层层递进的。   

例如下面的栗子  

```
apiVersion: batch/v2alpha1
kind: CronJob
...
```

在这个 YAML 文件中，“CronJob”就是这个 API 对象的资源类型（Resource），“batch”就是它的组（Group），v2alpha1 就是它的版本（Version）。   

这个 YAML 会被 Kubernetes 转化成一个 CronJob 对象，这里来看下是如何转化的。  

1、首先匹配 API 对象的组；  

Kubernetes 中的核心 API 对象，例如 Pod、Node。不需要 Group,(Group 为 "")，这些 API 对象会在 `/api` 这个层级进行匹配；  

非核心的 Api 对象类似 CronJob 会在 `/apis` 这个层级进行匹配，上面 CronJob 的 Group 是 batch，所以就会找到 `/apis/batch`；   

2、匹配 API 的版本号；  

上面栗子 CronJob 的版本号是 v2alpha1，所以当前匹配的路径就是 `/apis/batch/v2alpha1`;   

和我们 API 中的设计原则一样，Kubernetes 中的版本号也是用来进行 API 版本管理的。   

3、匹配 API 对象的资源类型；   

在匹配好版本之后，Kubernetes 根据路径和资源类型，就能确定创建的是 `/apis/batch/v2alpha1` 下的 CronJob 对象。    

当发起创建 CronJob 的 POST 请求之后，编写的 YAML 的信息就被提交给 APIServer 处理。   

APIServer 如何处理 API 请求    

1、请求过滤请求，进行一些前置性的工作，比如授权、超时处理、审计等；  

2、API 路由匹配；  

进入 MUX 和 Routes 流程，MUX 和 Routes 的主要作用是完成 APIServer 的 URL 和 Handler 绑定，Handler 找到对应的 CronJob 类型定义。  

3、根据提交的 YAML 文件创建资源；   

上面的栗子，根据这个 CronJob 类型定义，使用提交的 YAML 文件里的字段，创建一个 CronJob 对象。   

4、使用准入控制器，进行变更操作或验证操作；    

准入控制器 是一段代码，它会在请求通过认证和鉴权之后、对象被持久化之前拦截到达 API 服务器的请求。

准入控制器可以执行验证（Validating） 和/或变更（Mutating） 操作。  

5、序列化，保存到 Etcd 中。  

APIServer 会把验证过的 API 对象转换成用户最初提交的版本，进行序列化操作，并调用 Etcd 的 API 把它保存起来。  

### 参考

【深入剖析 Kubernetes】https://time.geekbang.org/column/intro/100015201?code=UhApqgxa4VLIA591OKMTemuH1%2FWyLNNiHZ2CRYYdZzY%3D  
【k8s 声明式 API】https://www.51cto.com/article/712066.html     
【声明式对比命令式】https://cloud.tencent.com/developer/article/1080886  
【Kubernetes 对象管理】https://kubernetes.io/zh-cn/docs/concepts/overview/working-with-objects/object-management/  





