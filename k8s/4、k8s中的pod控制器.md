<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [k8s 中 Pod 的控制器](#k8s-%E4%B8%AD-pod-%E7%9A%84%E6%8E%A7%E5%88%B6%E5%99%A8)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [Replication Controller](#replication-controller)
  - [ReplicaSet](#replicaset)
  - [Deployment](#deployment)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## k8s 中 Pod 的控制器

### 前言

Pod 是 Kubernetes 集群中能够被创建和管理的最小部署单元。所以需要有工具去操作和管理它们的生命周期,这里就需要用到控制器了。  

Pod 控制器由 master 的 `kube-controller-manager` 组件提供，常见的此类控制器有 `Replication Controller、ReplicaSet、Deployment、DaemonSet、StatefulSet、Job` 和 CronJob 等，它们分别以不同的方式管理 Pod 资源对象。    

### Replication Controller

RC 是 k8s 集群中最早的保证 Pod 高可用的 API 对象。它的作用就是保证集群中有指定数目的 pod 运行。

当前运行的 pod 数目少于指定的数目，RC 就会启动新的 pod 副本，保证运行 pod 数量等于指定数目。

当前运行的 pod 数目大于指定的数目，RC 就会杀死多余的 pod 副本。     

直接上栗子  

```
cat <<EOF >./pod-rc.yaml
apiVersion: v1
kind: ReplicationController
metadata:
  name: nginx
spec:
  replicas: 3
  selector:
    app: nginx
  template:
    metadata:
      name: nginx
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx
        ports:
        - containerPort: 80
EOF
```

在新版的 Kubernetes 中建议使用 ReplicaSet (RS)来取代 ReplicationController。  

ReplicaSet 跟 ReplicationController 没有本质的不同，只是名字不一样，但 ReplicaSet 支持集合式 selector。  

关于 ReplicationController 这里也不展开讨论了，主要看下 ReplicaSet。  

### ReplicaSet  

RS 是新一代 RC，提供同样的高可用能力，区别主要在于 RS 后来居上，能支持支持集合式 selector。  

副本集对象一般不单独使用，而是作为 Deployment 的理想状态参数使用。    

下面看下 Deployment 中是如何使用 ReplicaSet 的。  

### Deployment  

一个 Deployment 为 Pod 和 ReplicaSet 提供声明式的更新能力,每一个 Deployment 都对应集群中的一次部署。  

一般使用 Deployment 来管理 RS，可以用来创建一个新的服务，更新一个新的服务，也可以用来滚动升级一个服务。  

滚动升级一个服务，滚动升级一个服务，实际是创建一个新的 RS，然后逐渐将新 RS 中副本数增加到理想状态，将旧 RS 中的副本数减小到 0 的复合操作；这样一个复合操作用一个 RS 是不太好描述的，所以用一个更通用的 Deployment 来描述。    

举个栗子  

```
cat <<EOF >./nginx-deployment.yaml

apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels: # 这里定义需要管理的 pod，通过 Pod 的标签进行匹配
      app: nginx
  template:
    metadata:
      labels: # 运行的 pod 的标签
        app: nginx
    spec:
      containers: # pod 中运行的容器
      - name: nginx
        image: nginx:1.14.2 
        ports:
        - containerPort: 80
EOF
```


### 参考

【Deployments】https://kubernetes.io/zh-cn/docs/concepts/workloads/controllers/deployment/    