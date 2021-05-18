<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [kubespray部署k8s](#kubespray%E9%83%A8%E7%BD%B2k8s)
  - [准备](#%E5%87%86%E5%A4%87)
    - [需要关闭防火墙](#%E9%9C%80%E8%A6%81%E5%85%B3%E9%97%AD%E9%98%B2%E7%81%AB%E5%A2%99)
    - [配置hosts](#%E9%85%8D%E7%BD%AEhosts)
    - [处理镜像](#%E5%A4%84%E7%90%86%E9%95%9C%E5%83%8F)
    - [配置文件](#%E9%85%8D%E7%BD%AE%E6%96%87%E4%BB%B6)
  - [运行](#%E8%BF%90%E8%A1%8C)
    - [通过对应的镜像](#%E9%80%9A%E8%BF%87%E5%AF%B9%E5%BA%94%E7%9A%84%E9%95%9C%E5%83%8F)
    - [运行代码](#%E8%BF%90%E8%A1%8C%E4%BB%A3%E7%A0%81)
  - [查看结果](#%E6%9F%A5%E7%9C%8B%E7%BB%93%E6%9E%9C)
  - [出现的问题](#%E5%87%BA%E7%8E%B0%E7%9A%84%E9%97%AE%E9%A2%98)
    - [墙](#%E5%A2%99)
    - [错误的配置](#%E9%94%99%E8%AF%AF%E7%9A%84%E9%85%8D%E7%BD%AE)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## kubespray部署k8s

### 准备

[kubespray项目地址](https://github.com/kubernetes-sigs/kubespray)  

releases版本：v2.15.1  

#### 需要关闭防火墙  

具体命令，自行google  

#### 配置hosts

```go
cat >>/etc/hosts<<EOF
10.10.110.xxx k8s-master
10.10.110.xxx k8s-node1
10.10.110.xxx k8s-node2
EOF
```

#### 处理镜像

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
```go
$  docker load -i ./images.tar
```

#### 配置文件

配置文件  

```go
[all]
k8s-master ansible_host=10.10.110.xxx ip=10.10.110.xxx ansible_ssh_port=666
k8s-node1 ansible_host=10.10.110.xxx ip=10.10.110.xxx ansible_ssh_port=666
k8s-node2 ansible_host=10.10.11o.xxx ip=10.10.110.xxx ansible_ssh_port=666

[kube-master]
k8s-master

[etcd]
k8s-master
k8s-node1
k8s-node2

[kube-node]
k8s-node1
k8s-node2

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

### 查看结果

运行的结果  

```go
PLAY RECAP *************************************************************************************************************************************************************************************************
xxxxxxxxxxxx.zs            : ok=391  changed=49   unreachable=0    failed=0    skipped=639  rescued=0    ignored=1   
xxxxxxxxxxxx.zs            : ok=391  changed=49   unreachable=0    failed=0    skipped=638  rescued=0    ignored=1   
xxxxxxxxxxxx.zs            : ok=525  changed=75   unreachable=0    failed=0    skipped=1079 rescued=0    ignored=2   
localhost                  : ok=1    changed=0    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0   
```

查看

```go
➜  ~ kubectl get nodes
NAME              STATUS   ROLES    AGE     VERSION
xxxxxxxxxxxx.zs   Ready    <none>   6m24s   v1.19.9
xxxxxxxxxxxx.zs   Ready    <none>   6m25s   v1.19.9
xxxxxxxxxxxx.zs   Ready    master   8m10s   v1.19.9
```

### 出现的问题

#### 墙

一些镜像在跑的时候国内的网络拉不上来，需要手动处理，可参考上面的处理办法，对于一些版本的镜像，在`docker-hub`中搜索，就能解决。  

#### 错误的配置

`inventory.ini`文件没有配置`hostname`就开始搞了，结果`kubeadmin init`报错  

```go
ReadString: expects \" or n, but found 1, error found in #10 byte of ...|s.local\",10,\"base-se|..., bigger context
```
