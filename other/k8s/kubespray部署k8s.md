## kubespray部署k8s

### 依赖的镜像


```go
k8s.gcr.io/kube-proxy:v1.19.9
k8s.gcr.io/kube-controller-manager:v1.19.9
k8s.gcr.io/kube-scheduler:v1.19.9
k8s.gcr.io/kube-apiserver:v1.19.9
k8s.gcr.io/coredns:1.7.0
k8s.gcr.io/pause:3.3
k8s.gcr.io/dns/k8s-dns-node-cache:1.17.1
k8s.gcr.io/cpa/cluster-proportional-autoscaler-amd64:1.8.3
k8s.gcr.io/kube-scheduler:v1.19.9
```

镜像打包

打包导入到tar包中 

```go
$ docker save -o ./images.tar webapp:1.0 nginx:1.12 mysql:5.7

```

导出镜像
```
$  docker load -i ./images.tar
```

