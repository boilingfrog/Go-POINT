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

