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

下载，打tag
```go
docker pull liz2019/kube-proxy:v1.19.9
docker pull liz2019/kube-controller-manager:v1.19.9
docker pull liz2019/kube-scheduler:v1.19.9
docker pull liz2019/kube-apiserver:v1.19.9
docker pull liz2019/coredns:1.7.0
docker pull liz2019/pause:3.3
docker pull liz2019/k8s-dns-node-cache:1.17.1
docker pull liz2019/cluster-proportional-autoscaler-amd64:1.8.3
docker pull liz2019/kube-scheduler:v1.19.9


docker tag liz2019/kube-proxy:v1.19.9 k8s.gcr.io/kube-proxy:v1.19.9
docker tag liz2019/kube-controller-manager:v1.19.9 k8s.gcr.io/kube-controller-manager:v1.19.9
docker tag liz2019/kube-scheduler:v1.19.9 k8s.gcr.io/kube-scheduler:v1.19.9
docker tag liz2019/kube-apiserver:v1.19.9 k8s.gcr.io/kube-apiserver:v1.19.9
docker tag liz2019/coredns:1.7.0 k8s.gcr.io/coredns:1.7.0
docker tag liz2019/pause:3.3 k8s.gcr.io/pause:3.3
docker tag liz2019/k8s-dns-node-cache:1.17.1 k8s.gcr.io/dns/k8s-dns-node-cache:1.17.1
docker tag liz2019/cluster-proportional-autoscaler-amd64:1.8.3 k8s.gcr.io/cpa/cluster-proportional-autoscaler-amd64:1.8.3
docker tag liz2019/kube-scheduler:v1.19.9 k8s.gcr.io/kube-scheduler:v1.19.9

```

````go
docker pull k8smx/kube-apiserver:v1.19.11
docker pull k8smx/kube-controller-manager:v1.19.11
docker pull k8smx/kube-scheduler:v1.19.11
docker pull k8smx/kube-proxy:v1.19.11
docker pull k8smx/pause:3.2
docker pull k8smx/etcd:3.4.13-0
docker pull k8smx/coredns:1.7.0

docker tag k8smx/kube-apiserver:v1.19.11 k8s.gcr.io/kube-apiserver:v1.19.11
docker tag k8smx/kube-controller-manager:v1.19.11 k8s.gcr.io/kube-controller-manager:v1.19.11
docker tag k8smx/kube-scheduler:v1.19.11 k8s.gcr.io/kube-scheduler:v1.19.11
docker tag k8smx/kube-proxy:v1.19.11 k8s.gcr.io/kube-proxy:v1.19.11
docker tag k8smx/pause:3.2 k8s.gcr.io/pause:3.2
docker tag k8smx/etcd:3.4.13-0 k8s.gcr.io/etcd:3.4.13-0
docker tag k8smx/coredns:1.7.0 k8s.gcr.io/coredns:1.7.0
````


镜像打包

打包导入到tar包中 

```go
$ docker save -o ./images.tar webapp:1.0 nginx:1.12 mysql:5.7

```

导出镜像
```
$  docker load -i ./images.tar
```

运行的命令 

```go
docker run --rm -it --mount type=bind,source=C:/goWork/src/kubespray/inventory/sample,dst=/inventory --mount type=bind,source=C:/Users/rickl/.ssh/id_rsa,dst=/root/.ssh/id_rsa quay.io/kubespray/kubespray:v2.15.1 bash
ansible-playbook -i /inventory/inventory.ini  --become --become-user=root cluster.yml
```
