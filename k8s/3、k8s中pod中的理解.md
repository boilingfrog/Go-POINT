<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [k8s中Pod的理解](#k8s%E4%B8%ADpod%E7%9A%84%E7%90%86%E8%A7%A3)
  - [基本概念](#%E5%9F%BA%E6%9C%AC%E6%A6%82%E5%BF%B5)
  - [k8s 为什么使用 Pod 作为最小的管理单元](#k8s-%E4%B8%BA%E4%BB%80%E4%B9%88%E4%BD%BF%E7%94%A8-pod-%E4%BD%9C%E4%B8%BA%E6%9C%80%E5%B0%8F%E7%9A%84%E7%AE%A1%E7%90%86%E5%8D%95%E5%85%83)
  - [如何使用 Pod](#%E5%A6%82%E4%BD%95%E4%BD%BF%E7%94%A8-pod)
    - [1、自主式 Pod](#1%E8%87%AA%E4%B8%BB%E5%BC%8F-pod)
    - [2、控制器管理的 Pod](#2%E6%8E%A7%E5%88%B6%E5%99%A8%E7%AE%A1%E7%90%86%E7%9A%84-pod)
  - [静态 Pod](#%E9%9D%99%E6%80%81-pod)
  - [Pod的生命周期](#pod%E7%9A%84%E7%94%9F%E5%91%BD%E5%91%A8%E6%9C%9F)
  - [Pod 如何直接暴露服务](#pod-%E5%A6%82%E4%BD%95%E7%9B%B4%E6%8E%A5%E6%9A%B4%E9%9C%B2%E6%9C%8D%E5%8A%A1)
    - [hostNetwork](#hostnetwork)
    - [hostPort](#hostport)
    - [hostNetwork 和 hostPort 的对比](#hostnetwork-%E5%92%8C-hostport-%E7%9A%84%E5%AF%B9%E6%AF%94)
  - [资源限制](#%E8%B5%84%E6%BA%90%E9%99%90%E5%88%B6)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## k8s中Pod的理解  

### 基本概念

Pod 是 Kubernetes 集群中能够被创建和管理的最小部署单元,它是虚拟存在的。Pod 是一组容器的集合，并且部署在同一个 Pod 里面的容器是亲密性很强的一组容器，Pod 里面的容器，共享网络和存储空间，Pod 是短暂的。  

k8s 中的 Pod 有下面两种使用方式   

1、一个 Pod 中运行一个容器，这是最常见的用法。一个 Pod 封装一个容器，k8s 直接对 Pod 管理即可；  

2、一个 Pod 中同时运行多个容器，通常是紧耦合的我们才会放到一起。同一个 Pod 中的多个容器可以使用 localhost 通信，他们共享网络和存储卷。不过这种用法不常见，只有在特定的场景中才会使用。   

### k8s 为什么使用 Pod 作为最小的管理单元  

k8s 中为什么不直接操作容器而是使用 Pod 作为最小的部署单元呢？   

为了管理容器，k8s 需要更多的信息，比如重启策略，它定义了容器终止后要采取的策略;或者是一个可用性探针，从应用程序的角度去探测是否一个进程还存活着。基于这些原因，k8s 架构师决定使用一个新的实体，也就是 Pod，而不是重载容器的信息添加更多属性，用来在逻辑上包装一个或者多个容器的管理所需要的信息。  

### 如何使用 Pod 

#### 1、自主式 Pod  

我们可以简单的快速的创建一个 Pod 类似下面：  

```
$ cat pod.yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  containers:
  - name: nginx
    image: nginx:1.14.2
    ports:
    - containerPort: 80
```

创建  

```
$ kubectl apply -f pod.yaml -n study-k8s

$ kubectl get pod -n study-k8s
NAME    READY   STATUS    RESTARTS   AGE
nginx   1/1     Running   0          2m
```

自主创建的 Pod ,因为没有加入控制器来管理，这样创建的 Pod，被删除，或者因为意外退出了，不会重启自愈，直接就会被删除了。所以，业务中，我们在创建 Pod 的时候都会加入控制器。   

#### 2、控制器管理的 Pod

因为我们的业务长场景的需要，我们需要 Pod 有滚动升级，副本管理，集群级别的自愈能力，这时候我们就不能单独的创建 Pod, 我们需要通过相应的控制器来创建 Pod,来实现 Pod 滚动升级，自愈等的能力。  

对于 Pod 使用，我们最常使用的就是通过 Deployment 来管理。  

Deployment 提供了一种对 Pod 和 ReplicaSet 的管理方式。Deployment 可以用来创建一个新的服务，更新一个新的服务，也可以用来滚动升级一个服务。借助于 ReplicaSet 也可以实现 Pod 的副本管理功能。   

滚动升级一个服务，实际是创建一个新的 RS，然后逐渐将新 RS 中副本数增加到理想状态，将旧 RS 中的副本数减小到 0 的复合操作；这样一个复合操作用一个 RS 是不太好描述的，所以用一个更通用的 Deployment 来描述。  

```
$ vi deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
```

运行  

```
$ kubectl apply -f deployment.yaml -n study-k8s
$ kubectl get pods -n study-k8s
NAME                                READY   STATUS    RESTARTS   AGE
nginx                               1/1     Running   0          67m
nginx-deployment-66b6c48dd5-24sbd   1/1     Running   0          59s
nginx-deployment-66b6c48dd5-wxkln   1/1     Running   0          59s
nginx-deployment-66b6c48dd5-xgzgh   1/1     Running   0          59s
```
 
因为上面定义了 `replicas: 3` 也就是副本数为3，当我们删除一个 Pod 的时候，马上就会有一个新的 Pod 被创建。    

同样经常使用的到的控制器还有 DaemonSet 和 StatefulSet  

DaemonSet：DaemonSet 确保全部（或者某些）节点上运行一个 Pod 的副本。 当有节点加入集群时， 也会为他们新增一个 Pod 。 当有节点从集群移除时，这些 Pod 也会被回收。删除 DaemonSet 将会删除它创建的所有 Pod。  

StatefulSet：用来管理有状态应用的工作负载，和 Deployment 类似， StatefulSet 管理基于相同容器规约的一组 Pod。但和 Deployment 不同的是， StatefulSet 为它们的每个 Pod 维护了一个有粘性的 ID。这些 Pod 是基于相同的规约来创建的， 但是不能相互替换：无论怎么调度，每个 Pod 都有一个永久不变的 ID。  

### 静态 Pod 

静态 Pod 是由 kubelet 进行管理的仅存在与特定 node 上的 Pod，他们不能通过 `api server` 进行管理，无法与 `rc,deployment,ds` 进行关联，并且 kubelet 无法对他们进行健康检查。   

静态 Pod 始终绑定在某一个kubelet，并且始终运行在同一个节点上。 kubelet会自动为每一个静态 Pod 在 Kubernetes 的 apiserver 上创建一个镜像 Pod（Mirror Pod），因此我们可以在 apiserver 中查询到该 Pod，但是不能通过 apiserver 进行控制（例如不能删除）。   

为什么需要静态 Pod ?

主要是用来对集群中的组件进行容器化操作，例如 `etcd  kube-apiserver  kube-controller-manager  kube-scheduler` 这些都是静态 Pod 资源。   

因为这些 Pod 不受 apiserver 的控制，就不会不小新被删掉的情况，同时 kube-apiserver 也不能自己去控制自己。静态 Pod 的存在将集群中的容器化操作提供了可能。    

静态 Pod 的创建有两种方式，配置文件和 HTTP 两种方式，具体参见。[静态 Pod 的创建](https://kubernetes.io/zh-cn/docs/tasks/configure-Pod-container/static-Pod/#static-Pod-creation)   

### Pod的生命周期

Pod 在运行的过程中会被定义为各种状态，了解一些状态能帮助我们了解 Pod 的调度策略。当 Pod 被创建之后，就会进入健康检查状态，当 Kubernetes 确定当前 Pod 已经能够接受外部的请求时，才会将流量打到新的 Pod 上并继续对外提供服务，在这期间如果发生了错误就可能会触发重启机制。    

不过 Pod 本身不具有自愈能力，如果 Pod 因为 Node 故障，或者是调度器本身故障，这个 Pod 就会被删除。所以 Pod 中一般使用控制器来管理 Pod ，来实现 Pod 的自愈能力和滚动更新的能力。    

Pod的重启策略包括  

- Always 只要失败，就会重启；  

- OnFile 当容器终止运行，且退出吗不是0，就会重启；  

- Never  从来不会重启。   

重启的时间，是以2n来算。比如（10s、20s、40s、...），其最长延迟为 5 分钟。 一旦某容器执行了 10 分钟并且没有出现问题，kubelet 对该容器的重启回退计时器执行 重置操作。  

管理Pod的重启策略是靠控制器来完成的。   

**Pod 的几种状态**  

使用的过程中，会经常遇到下面几种 Pod 的状态。    

`Pending`：Pod 创建已经提交给 k8s，但有一个或者多个容器尚未创建亦未运行，此阶段包括等待 Pod 被调度的时间和通过网络下载镜像的时间。这个状态可能就是在下载镜像；  

`Running`：Pod 已经绑定到一个节点上了，并且已经创建了所有容器。至少有一个容器仍在运行，或者正处于启动或重启状态；    

`Secceeded`：Pod 中的所有容器都已经成功终止，并且不会再重启；     

`Failed`：Pod 中所有的容器均已经终止，并且至少有一个容器已经在故障中终止；  

`Unkown`：因为某些原因无法取得 Pod 的状态。这种情况通常是因为与 Pod 所在主机通信失败。

当 pod 一直处于 Pending 状态，可通过 `kubectl describe pods <node_name> -n namespace` 来获取出错的信息  

### Pod 如何直接暴露服务

Pod 一般不直接对外暴露服务，一个 Pod 只是一个运行服务的实例，随时可能在节点上停止，然后再新的节点上用一个新的 IP 启动一个新的 Pod,因此不能使用确定的 IP 和端口号提供服务。这对于业务来说，就不能根据 Pod 的 IP 作为业务调度。kubernetes 就引入了 Service 的概 念，它为 Pod 提供一个入口，主要通过 Labels 标签来选择后端Pod，这时候不论后端 Pod 的 IP 地址如何变更，只要 Pod 的 Labels 标签没变，那么 业务通过 service 调度就不会存在问题。  

不过使用 hostNetwork 和 hostPort，可以直接暴露 node 的 ip 地址。  

#### hostNetwork  

这是一种直接定义 Pod 网络的方式，使用 hostNetwork 配置网络，Pod 中的所有容器就直接暴露在宿主机的网络环境中，这时候，Pod 的 PodIP 就是其所在 Node 的 IP。从原理上来说，当设定 Pod 的网络为 Host 时，是设定了 Pod 中 `pod-infrastructure`（或pause）容器的网络为 Host，Pod 内部其他容器的网络指向该容器。   

```
cat <<EOF >./pod-host.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
      hostNetwork: true
EOF
```

运行  

```shell
$ kubectl apply -f pod-host.yaml -n study-k8s
$ kubectl get pods -n study-k8s -o wide
NAME                                READY   STATUS    RESTARTS   AGE     IP             NODE              NOMINATED NODE   READINESS GATES
nginx-deployment-6d47cff9fd-5bzjq   1/1     Running   0          6m25s   192.168.56.111   kube-server8.zs   <none>           <none>

$ curl http://192.168.56.111/
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
    body {
        width: 35em;
        margin: 0 auto;
        font-family: Tahoma, Verdana, Arial, sans-serif;
    }
</style>
</head>
<body>
<h1>Welcome to nginx!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a>.<br/>
Commercial support is available at
<a href="http://nginx.com/">nginx.com</a>.</p>

<p><em>Thank you for using nginx.</em></p>
</body>
</html>
```

一般情况下除非知道需要某个特定应用占用特定宿主机上的特定端口时才使用 `hostNetwork: true` 的方式。   

#### hostPort

这是一种直接定义Pod网络的方式。  

hostPort 是直接将容器的端口与所调度的节点上的端口路由，这样用户就可以通过宿主机的 IP 加上端口来访问 Pod。    

```
cat <<EOF >./pod-hostPort.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
          hostPort: 8000
EOF
```

运行

```shell
$ kubectl apply -f pod-hostPort.yaml -n study-k8s
$ kubectl get pods -n study-k8s -o wide
NAME                                READY   STATUS    RESTARTS   AGE     IP             NODE              NOMINATED NODE   READINESS GATES
nginx-deployment-6d47cff9fd-5bzjq   1/1     Running   0          3m25s   192.168.56.111   kube-server8.zs   <none>           <none>

$ curl http://192.168.56.111:8000/
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
    body {
        width: 35em;
        margin: 0 auto;
        font-family: Tahoma, Verdana, Arial, sans-serif;
    }
</style>
</head>
<body>
<h1>Welcome to nginx!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a>.<br/>
Commercial support is available at
<a href="http://nginx.com/">nginx.com</a>.</p>

<p><em>Thank you for using nginx.</em></p>
</body>
</html>
```

这种网络方式可以用来做 `nginx [Ingress controller]`。外部流量都需要通过 `kubenretes node` 节点的 80 和 443 端口。  

#### hostNetwork 和 hostPort 的对比

相同点  

hostNetwork 和 hostPort 的本质都是暴露 Pod 所在的节点 IP 给终端用户。因为 Pod 的生命周期并不固定，随时都有被重构的可能。可以使用 DaemonSet 或者非亲缘性策略，来保证每个 node 节点只有一个Pod 被部署。  

不同点  

使用 hostNetwork，pod 实际上用的是 pod 宿主机的网络地址空间；  

使用 hostPort，pod IP 并非宿主机 IP，而是 cni 分配的 pod IP，跟其他普通的 pod 使用一样的 ip 分配方式，端口并非宿主机网络监听端口，只是使用了 DNAT 机制将 hostPort 指定的端口映射到了容器的端口之上（可以通过 iptables 命令进行查看）。  

### 资源限制

每个 Pod 都可以对其能使用的服务器上的计算资源设置限额，当前可以设置限额的计算资源有 CPU 和 Memory 两种，其中 CPU 的资源单位为 CPU（Core）的数量，是一个绝对值。    

对于容器来说一个 CPU 的配额已经是相当大的资源配额了，所以在 Kubernetes 里，通常以千分之一的 CPU 配额为最小单位，用 m 来表示。通常一个容器的CPU配额被 定义为 100-300m，即占用0.1-0.3个CPU。与 CPU 配额类似，Memory 配额也是一个绝对值，它的单位是内存字节数。    

对计算资源进行配额限定需要设定以下两个参数：  

- Requests：该资源的最小申请量，系统必须满足要求。  

- Limits：该资源最大允许使用的量，不能超过这个使用限制，当容器试图使用超过这个量的资源时，可能会被Kubernetes Kill并重启。

通常我们应该把 Requests 设置为一个比较小的数值，满足容器平时的工作负载情况下的资源需求，而把 Limits 设置为峰值负载情况下资源占用的最大量。下面是一个资源配额的简单定义：  

```
spec:
  containers:
  - name: db
    image: mysql
    resources:
      requests:
        memory: "64Mi"
        cpu: "250m"
      limits:
        memory: "128Mi"
        cpu: "500m"
```

最小0.25个CPU及64MB内存，最大0.5个CPU及128MB内存。

### 参考
【初识Kubernetes（K8s）：各种资源对象的理解和定义】https://blog.51cto.com/andyxu/2329257  
【Kubernetes系列学习文章 - Pod的深入理解（四）】https://cloud.tencent.com/developer/article/1443520  
【详解 Kubernetes Pod 的实现原理】https://draveness.me/kubernetes-pod/  
