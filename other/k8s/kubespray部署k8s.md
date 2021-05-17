<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [kubespray部署k8s](#kubespray%E9%83%A8%E7%BD%B2k8s)
  - [依赖的镜像](#%E4%BE%9D%E8%B5%96%E7%9A%84%E9%95%9C%E5%83%8F)
  - [配置文件](#%E9%85%8D%E7%BD%AE%E6%96%87%E4%BB%B6)
  - [运行](#%E8%BF%90%E8%A1%8C)
    - [通过对应的镜像](#%E9%80%9A%E8%BF%87%E5%AF%B9%E5%BA%94%E7%9A%84%E9%95%9C%E5%83%8F)
    - [运行代码](#%E8%BF%90%E8%A1%8C%E4%BB%A3%E7%A0%81)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

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

镜像打包

打包导入到tar包中 

```go
$ docker save -o ./images.tar webapp:1.0 nginx:1.12 mysql:5.7

```

导出镜像
```
$  docker load -i ./images.tar
```

### 配置文件

```toml
[all]
10.10.110.xx ansible_ssh_port=333
10.10.110.xx ansible_ssh_port=333
10.10.110.xx ansible_ssh_port=333

[kube-master]
10.10.149.198

[etcd]
10.10.149.198
10.10.130.113
10.10.135.26

[kube-node]
10.10.130.113
10.10.135.26

[calico-rr]

[k8s-cluster:children]
kube-master
kube-node
calico-rr
```

### 运行 

#### 通过对应的镜像

```go
$ docker pull quay.io/kubespray/kubespray:v2.15.1
$ docker run --rm -it --mount type=bind,source="$(pwd)"/inventory/sample,dst=/inventory \
  --mount type=bind,source="${HOME}"/.ssh/id_rsa,dst=/root/.ssh/id_rsa \
  quay.io/kubespray/kubespray:v2.15.1 bash
$ ansible-playbook -i /inventory/inventory.ini --private-key /root/.ssh/id_rsa cluster.yml
```

#### 运行代码

```go
$ ansible-playbook -i /inventory/inventory.ini --private-key /root/.ssh/id_rsa cluster.yml
```

运行出错，删除

```go
$ ansible-playbook -i /inventory/inventory.ini --private-key /root/.ssh/id_rsa reset.yml
```