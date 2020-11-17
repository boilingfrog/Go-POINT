<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [二进制部署k8s](#%E4%BA%8C%E8%BF%9B%E5%88%B6%E9%83%A8%E7%BD%B2k8s)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [准备工作](#%E5%87%86%E5%A4%87%E5%B7%A5%E4%BD%9C)
    - [关闭防火墙](#%E5%85%B3%E9%97%AD%E9%98%B2%E7%81%AB%E5%A2%99)
    - [关闭 swap 分区](#%E5%85%B3%E9%97%AD-swap-%E5%88%86%E5%8C%BA)
    - [关闭 SELinux](#%E5%85%B3%E9%97%AD-selinux)
    - [更新系统时间](#%E6%9B%B4%E6%96%B0%E7%B3%BB%E7%BB%9F%E6%97%B6%E9%97%B4)
  - [秘钥免密码](#%E7%A7%98%E9%92%A5%E5%85%8D%E5%AF%86%E7%A0%81)
  - [设置主机名称](#%E8%AE%BE%E7%BD%AE%E4%B8%BB%E6%9C%BA%E5%90%8D%E7%A7%B0)
  - [服务器角色](#%E6%9C%8D%E5%8A%A1%E5%99%A8%E8%A7%92%E8%89%B2)
  - [安装etcd](#%E5%AE%89%E8%A3%85etcd)
    - [创建证书](#%E5%88%9B%E5%BB%BA%E8%AF%81%E4%B9%A6)
    - [生成证书](#%E7%94%9F%E6%88%90%E8%AF%81%E4%B9%A6)
    - [部署Etcd](#%E9%83%A8%E7%BD%B2etcd)
  - [在Node安装Docker](#%E5%9C%A8node%E5%AE%89%E8%A3%85docker)
  - [Flannel网络](#flannel%E7%BD%91%E7%BB%9C)
  - [master节点部署组件](#master%E8%8A%82%E7%82%B9%E9%83%A8%E7%BD%B2%E7%BB%84%E4%BB%B6)
    - [生成证书](#%E7%94%9F%E6%88%90%E8%AF%81%E4%B9%A6-1)
    - [配置apiserver组件](#%E9%85%8D%E7%BD%AEapiserver%E7%BB%84%E4%BB%B6)
    - [部署scheduler组件](#%E9%83%A8%E7%BD%B2scheduler%E7%BB%84%E4%BB%B6)
    - [部署controller-manager组件](#%E9%83%A8%E7%BD%B2controller-manager%E7%BB%84%E4%BB%B6)
  - [在Node节点部署组件](#%E5%9C%A8node%E8%8A%82%E7%82%B9%E9%83%A8%E7%BD%B2%E7%BB%84%E4%BB%B6)
    - [创建kubeconfig文件](#%E5%88%9B%E5%BB%BAkubeconfig%E6%96%87%E4%BB%B6)
    - [部署kubelet组件](#%E9%83%A8%E7%BD%B2kubelet%E7%BB%84%E4%BB%B6)
    - [部署kube-proxy组件](#%E9%83%A8%E7%BD%B2kube-proxy%E7%BB%84%E4%BB%B6)
  - [查看集群状态](#%E6%9F%A5%E7%9C%8B%E9%9B%86%E7%BE%A4%E7%8A%B6%E6%80%81)
  - [测试](#%E6%B5%8B%E8%AF%95)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 二进制部署k8s

### 前言

开始学习k8s吧，作为小白，先从手动搭建开始吧，然后在慢慢了解每个组件，本文经尝试能够成功部署成功。，本次尝试是看了阿良老师的视频，并且尝试了很多次，才成功的。阿良老师的博客地址，在下文有链接。  

### 准备工作

#### 关闭防火墙

关闭服务，并设为开机不自启

```
$ sudo systemctl stop firewalld
$ sudo systemctl disable firewalld
```

清空防火墙规则

```
$ sudo iptables -F && sudo iptables -X && sudo iptables -F -t nat && sudo iptables -X -t nat
$ sudo iptables -P FORWARD ACCEPT
```

#### 关闭 swap 分区

如果开启了 swap 分区，kubelet 会启动失败(可以通过将参数 --fail-swap-on 设置为false 来忽略 swap on)，故需要在每台机器上关闭 swap 分区：

```
$ sudo swapoff -a
```

为了防止开机自动挂载 swap 分区，可以注释 /etc/fstab 中相应的条目：

```
$ sudo sed -i '/ swap / s/^\(.*\)$/#\1/g' /etc/fstab
```

#### 关闭 SELinux

关闭 SELinux，否则后续 K8S 挂载目录时可能报错 Permission denied ：

```
$ sudo setenforce 0
```

修改配置文件，永久生效；

```
$ grep SELINUX /etc/selinux/config

SELINUX=disabled
```

#### 更新系统时间

1、调整系统 TimeZone

```
$ sudo timedatectl set-timezone Asia/Shanghai
```

2、将当前的 UTC 时间写入硬件时钟

```
$ sudo timedatectl set-local-rtc 0
```

3、重启依赖于系统时间的服务

```
$ sudo systemctl restart rsyslog
$ sudo systemctl restart crond
```

4、更新时间

```
$ yum -y install ntpdate
$ sudo ntpdate cn.pool.ntp.org
```

### 秘钥免密码

```
$ ssh-copy-id root@192.168.56.201
```

### 设置主机名称


192.168.56.201 kube-master
192.168.56.202 kube-node1
192.168.56.203 kube-node2


设置永久主机名称，然后重新登录  

```
$ sudo hostnamectl set-hostname kube-master
$ sudo hostnamectl set-hostname kube-node1
$ sudo hostnamectl set-hostname kube-node2
```

修改 /etc/hostname 文件，添加主机名和 IP 的对应关系：

```
$ vim /etc/hosts

192.168.56.201 kube-master
192.168.56.202 kube-node1
192.168.56.203 kube-node2
```

### 服务器角色

|  角色   | ip  | 组件  |
|  ----  | ----  | ----  |
| kube-master  | 192.168.56.201 | kube-apiserver，kube-controller-manager，kube-scheduler，etcd |
| kube-node1  | 192.168.56.202 | kubelet，kube-proxy，docker，flannel，etcd |
| kube-node2  | 192.168.56.203 | kubelet，kube-proxy，docker，flannel，etcd |

### 安装etcd

集群中每台机器都需要安装，先在一台节点安装配置，之后把配置文件scp到其他机器就行了  

#### 创建证书

使用cfssl来生成自签证书  

安装

wget工具安装`yum -y install wget`

```
$ wget https://pkg.cfssl.org/R1.2/cfssl_linux-amd64
$ wget https://pkg.cfssl.org/R1.2/cfssljson_linux-amd64
$ wget https://pkg.cfssl.org/R1.2/cfssl-certinfo_linux-amd64
$ chmod +x cfssl_linux-amd64 cfssljson_linux-amd64 cfssl-certinfo_linux-amd64
$ mv cfssl_linux-amd64 /usr/local/bin/cfssl
$ mv cfssljson_linux-amd64 /usr/local/bin/cfssljson
$ mv cfssl-certinfo_linux-amd64 /usr/bin/cfssl-certinfo
```

#### 生成证书

证书一样，保证集群中使用的证书一致，在其中一台机器中生成需要的证书，然后scp到其他的机器中  

创建需要的文件

创建文件夹`/opt/etcd/ssl`，etcd的证书文件都放到这个文件夹中  

```
$ vi ca-config.json
{
  "signing": {
    "default": {
      "expiry": "87600h"
    },
    "profiles": {
      "www": {
         "expiry": "87600h",
         "usages": [
            "signing",
            "key encipherment",
            "server auth",
            "client auth"
        ]
      }
    }
  }
}

$ vi ca-csr.json
{
    "CN": "etcd CA",
    "key": {
        "algo": "rsa",
        "size": 2048
    },
    "names": [
        {
            "C": "CN",
            "L": "Beijing",
            "ST": "Beijing"
        }
    ]
}

$ vi server-csr.json
{
    "CN": "etcd",
    "hosts": [
    "192.168.56.201",
    "192.168.56.202",
    "192.168.56.203"
    ],
    "key": {
        "algo": "rsa",
        "size": 2048
    },
    "names": [
        {
            "C": "CN",
            "L": "BeiJing",
            "ST": "BeiJing"
        }
    ]
}
```
生成证书

```
$ cfssl gencert -initca ca-csr.json | cfssljson -bare ca -
$ cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=www server-csr.json | cfssljson -bare server
$  ls *pem
ca-key.pem  ca.pem  server-key.pem  server.pem
```

#### 部署Etcd

```
$ wget https://github.com/etcd-io/etcd/releases/download/v3.2.12/etcd-v3.2.12-linux-amd64.tar.gz

$ mkdir /opt/etcd/{bin,cfg,ssl} -p
$ tar zxvf etcd-v3.2.12-linux-amd64.tar.gz
$ mv etcd-v3.2.12-linux-amd64/{etcd,etcdctl} /opt/etcd/bin/
```

创建配置文件

```
$ vi /opt/etcd/cfg/etcd   
#[Member]
ETCD_NAME="etcd01"
ETCD_DATA_DIR="/var/lib/etcd/default.etcd"
ETCD_LISTEN_PEER_URLS="https://192.168.56.201:2380"
ETCD_LISTEN_CLIENT_URLS="https://192.168.56.201:2379"

#[Clustering]
ETCD_INITIAL_ADVERTISE_PEER_URLS="https://192.168.56.201:2380"
ETCD_ADVERTISE_CLIENT_URLS="https://192.168.56.201:2379"
ETCD_INITIAL_CLUSTER="etcd01=https://192.168.56.201:2380,etcd02=https://192.168.56.202:2380,etcd03=https://192.168.56.203:2380"
ETCD_INITIAL_CLUSTER_TOKEN="etcd-cluster"
ETCD_INITIAL_CLUSTER_STATE="new"
```

- ETCD_NAME 节点名称
- ETCD_DATA_DIR 数据目录
- ETCD_LISTEN_PEER_URLS 集群通信监听地址
- ETCD_LISTEN_CLIENT_URLS 客户端访问监听地址
- ETCD_INITIAL_ADVERTISE_PEER_URLS 集群通告地址
- ETCD_ADVERTISE_CLIENT_URLS 客户端通告地址
- ETCD_INITIAL_CLUSTER 集群节点地址
- ETCD_INITIAL_CLUSTER_TOKEN 集群Token
- ETCD_INITIAL_CLUSTER_STATE 加入集群的当前状态，new是新集群，existing表示加入已有集群

systemd管理etcd
```
$ vi /usr/lib/systemd/system/etcd.service 

[Unit]
Description=Etcd Server
After=network.target
After=network-online.target
Wants=network-online.target

[Service]
Type=notify
EnvironmentFile=/opt/etcd/cfg/etcd
ExecStart=/opt/etcd/bin/etcd \
--name=${ETCD_NAME} \
--data-dir=${ETCD_DATA_DIR} \
--listen-peer-urls=${ETCD_LISTEN_PEER_URLS} \
--listen-client-urls=${ETCD_LISTEN_CLIENT_URLS},http://127.0.0.1:2379 \
--advertise-client-urls=${ETCD_ADVERTISE_CLIENT_URLS} \
--initial-advertise-peer-urls=${ETCD_INITIAL_ADVERTISE_PEER_URLS} \
--initial-cluster=${ETCD_INITIAL_CLUSTER} \
--initial-cluster-token=${ETCD_INITIAL_CLUSTER_TOKEN} \
--initial-cluster-state=new \
--cert-file=/opt/etcd/ssl/server.pem \
--key-file=/opt/etcd/ssl/server-key.pem \
--peer-cert-file=/opt/etcd/ssl/server.pem \
--peer-key-file=/opt/etcd/ssl/server-key.pem \
--trusted-ca-file=/opt/etcd/ssl/ca.pem \
--peer-trusted-ca-file=/opt/etcd/ssl/ca.pem
Restart=on-failure
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
```

启动

```
$ systemctl start etcd
$ systemctl enable etcd
```

只安装了一台机器，etcd是启动不起来的，我们需要把其他机器的都安装下  

只需要把上面添加的文件，scp到目标机器，然后修改下ip即可  

查看状态 

```
$ journalctl -u etcd
```

三台机器都安装完成之后查看状态  

检查etcd集群状态  

```go
$ /opt/etcd/bin/etcdctl \
   --ca-file=/opt/etcd/ssl/ca.pem --cert-file=/opt/etcd/ssl/server.pem --key-file=/opt/etcd/ssl/server-key.pem \
   --endpoints="https://192.168.56.201:2379,https://192.168.56.202:2379,https://192.168.56.203:2379" \
   cluster-health

  member 8c78d744f172cba9 is healthy: got healthy result from https://192.168.56.202:2379
  member bdc976e03235ad9b is healthy: got healthy result from https://192.168.56.201:2379
  member c6274de5e02a53ad is healthy: got healthy result from https://192.168.56.203:2379
  cluster is healthy
```

### 在Node安装Docker

在node的每台机器中安装  

```
$ yum install -y yum-utils device-mapper-persistent-data lvm2
$ yum-config-manager \
    --add-repo \
    https://download.docker.com/linux/centos/docker-ce.repo
$ yum install docker-ce -y
$ curl -sSL https://get.daocloud.io/daotools/set_mirror.sh | sh -s http://bc437cce.m.daocloud.io
$ systemctl start docker
$ systemctl enable docker
```

### Flannel网络

![channel](/img/k8s_flannel_1.png?raw=true)

Flannel是部署在node中的，每个node都需要安装，master节点无需安装   

配置子网
```go
$ /opt/etcd/bin/etcdctl \
--ca-file=/opt/etcd/ssl/ca.pem --cert-file=/opt/etcd/ssl/server.pem --key-file=/opt/etcd/ssl/server-key.pem \
--endpoints="https://192.168.56.201:2379,https://192.168.56.202:2379,https://192.168.56.203:2379" \
set /coreos.com/network/config  '{ "Network": "172.17.0.0/16", "Backend": {"Type": "vxlan"}}'
```

在一个node中安装，然后把配置的文件scp到其他node中  

```
$ wget https://github.com/coreos/flannel/releases/download/v0.9.1/flannel-v0.9.1-linux-amd64.tar.gz
$ tar zxvf flannel-v0.9.1-linux-amd64.tar.gz
$ mv flanneld mk-docker-opts.sh /opt/kubernetes/bin
```

配置Flannel

```
$ vi /opt/kubernetes/cfg/flanneld
FLANNEL_OPTIONS="--etcd-endpoints=https://192.168.56.201:2379,https://192.168.56.202:2379,https://192.168.56.203:2379 -etcd-cafile=/opt/etcd/ssl/ca.pem -etcd-certfile=/opt/etcd/ssl/server.pem -etcd-keyfile=/opt/etcd/ssl/server-key.pem"
```

systemd管理Flannel 

````
$ vi /usr/lib/systemd/system/flanneld.service
[Unit]
Description=Flanneld overlay address etcd agent
After=network-online.target network.target
Before=docker.service

[Service]
Type=notify
EnvironmentFile=/opt/kubernetes/cfg/flanneld
ExecStart=/opt/kubernetes/bin/flanneld --ip-masq $FLANNEL_OPTIONS
ExecStartPost=/opt/kubernetes/bin/mk-docker-opts.sh -k DOCKER_NETWORK_OPTIONS -d /run/flannel/subnet.env
Restart=on-failure

[Install]
WantedBy=multi-user.target
````

配置Docker启动指定子网段

```
$ vi /usr/lib/systemd/system/docker.service 

[Unit]
Description=Docker Application Container Engine
Documentation=https://docs.docker.com
After=network-online.target firewalld.service
Wants=network-online.target

[Service]
Type=notify
EnvironmentFile=/run/flannel/subnet.env
ExecStart=/usr/bin/dockerd $DOCKER_NETWORK_OPTIONS
ExecReload=/bin/kill -s HUP $MAINPID
LimitNOFILE=infinity
LimitNPROC=infinity
LimitCORE=infinity
TimeoutStartSec=0
Delegate=yes
KillMode=process
Restart=on-failure
StartLimitBurst=3
StartLimitInterval=60s

[Install]
WantedBy=multi-user.target
```

这块主要加入flannel的配置

```
EnvironmentFile=/run/flannel/subnet.env
ExecStart=/usr/bin/dockerd $DOCKER_NETWORK_OPTIONS
```

重启flannel和docker：

```
# systemctl daemon-reload
# systemctl start flanneld
# systemctl enable flanneld
# systemctl restart docker
```

检查是否生效

```
$ ps -ef |grep docker
root      2859     1  0 18:13 ?        00:00:05 /usr/bin/dockerd-current --add-runtime docker-runc=/usr/libexec/docker/docker-runc-current --default-runtime=docker-runc --exec-opt native.cgroupdriver=systemd --userland-proxy-path=/usr/libexec/docker/docker-proxy-current --init-path=/usr/libexec/docker/docker-init-current --seccomp-profile=/etc/docker/seccomp.json --selinux-enabled --log-driver=journald --signature-verification=false --storage-driver overlay2 --bip=172.17.96.1/24 --ip-masq=false --mtu=1450
root      2865  2859  0 18:13 ?        00:00:02 /usr/bin/docker-containerd-current -l unix:///var/run/docker/libcontainerd/docker-containerd.sock --metrics-interval=0 --start-timeout 2m --state-dir /var/run/docker/libcontainerd/containerd --shim docker-containerd-shim --runtime docker-runc --runtime-args --systemd-cgroup=true
root      5799  1753  0 18:49 pts/0    00:00:00 grep --color=auto docker

$ ip addr
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
2: enp0s3: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP group default qlen 1000
    link/ether 08:00:27:73:3f:bf brd ff:ff:ff:ff:ff:ff
    inet 10.0.2.7/24 brd 10.0.2.255 scope global noprefixroute dynamic enp0s3
       valid_lft 416sec preferred_lft 416sec
3: enp0s8: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP group default qlen 1000
    link/ether 08:00:27:7b:1e:59 brd ff:ff:ff:ff:ff:ff
    inet 192.168.56.203/24 brd 192.168.56.255 scope global noprefixroute enp0s8
       valid_lft forever preferred_lft forever
4: flannel.1: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1450 qdisc noqueue state UNKNOWN group default 
    link/ether c6:65:64:77:fd:04 brd ff:ff:ff:ff:ff:ff
    inet 172.17.96.0/32 scope global flannel.1
       valid_lft forever preferred_lft forever
5: docker0: <NO-CARRIER,BROADCAST,MULTICAST,UP> mtu 1500 qdisc noqueue state DOWN group default 
    link/ether 02:42:9a:42:ae:5b brd ff:ff:ff:ff:ff:ff
    inet 172.17.96.1/24 scope global docker0
       valid_lft forever preferred_lft forever
```

测试不同节点互通，在当前节点访问另一个Node节点flannel网络  

```
$ ping 172.17.72.0
PING 172.17.72.0 (172.17.72.0) 56(84) bytes of data.
64 bytes from 172.17.72.0: icmp_seq=1 ttl=64 time=1.25 ms
64 bytes from 172.17.72.0: icmp_seq=2 ttl=64 time=0.507 ms
64 bytes from 172.17.72.0: icmp_seq=3 ttl=64 time=1.17 ms
64 bytes from 172.17.72.0: icmp_seq=4 ttl=64 time=1.84 ms
64 bytes from 172.17.72.0: icmp_seq=5 ttl=64 time=0.426 ms
```

本人在部署的时候，选了最高版本的flannel，但是发现部署一致不成功，google了发现，需要找和etcd版本符合的flannel版本进行部署。  

### master节点部署组件

master节点需要部署apiserver，   不过前面部署的etcd,docker,flannel需要保证正常运行  

#### 生成证书

将证书统一放在`/opt/kubernetes/ssl/`进行管理

创建ca证书

```go
$ vi ca-config.json
{
  "signing": {
    "default": {
      "expiry": "87600h"
    },
    "profiles": {
      "kubernetes": {
         "expiry": "87600h",
         "usages": [
            "signing",
            "key encipherment",
            "server auth",
            "client auth"
        ]
      }
    }
  }
}

$ vi ca-csr.json
{
    "CN": "kubernetes",
    "key": {
        "algo": "rsa",
        "size": 2048
    },
    "names": [
        {
            "C": "CN",
            "L": "Beijing",
            "ST": "Beijing",
            "O": "k8s",
            "OU": "System"
        }
    ]
}

$ cfssl gencert -initca ca-csr.json | cfssljson -bare ca -
```

生成apiserver证书

```go
$ vi server-csr.json
{
    "CN": "kubernetes",
    "hosts": [
      "10.0.0.1",
      "127.0.0.1",
      "192.168.56.201",
      "kubernetes",
      "kubernetes.default",
      "kubernetes.default.svc",
      "kubernetes.default.svc.cluster",
      "kubernetes.default.svc.cluster.local"
    ],
    "key": {
        "algo": "rsa",
        "size": 2048
    },
    "names": [
        {
            "C": "CN",
            "L": "BeiJing",
            "ST": "BeiJing",
            "O": "k8s",
            "OU": "System"
        }
    ]
}
$ cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=kubernetes server-csr.json | cfssljson -bare server
```

生成kube-proxy证书

```go
$ vi kube-proxy-csr.json
{
  "CN": "system:kube-proxy",
  "hosts": [],
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "CN",
      "L": "BeiJing",
      "ST": "BeiJing",
      "O": "k8s",
      "OU": "System"
    }
  ]
}

$ cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=kubernetes kube-proxy-csr.json | cfssljson -bare kube-proxy
```

最终生成以下证书文件 

```
$ ls *pem
ca-key.pem  ca.pem  kube-proxy-key.pem  kube-proxy.pem  server-key.pem  server.pem
```

#### 配置apiserver组件

下载二进制安装包 `https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG-1.12.md`

kubernetes-server-linux-amd64.tar.gz  

```
$ mkdir /opt/kubernetes/{bin,cfg,ssl} -p
$ tar zxvf kubernetes-server-linux-amd64.tar.gz
$ cd kubernetes/server/bin
$ cp kube-apiserver kube-scheduler kube-controller-manager kubectl /opt/kubernetes/bin
```

创建token   

```
$ vi /opt/kubernetes/cfg/token.csv
674c457d4dcf2eefe4920d7dbb6b0ddc,kubelet-bootstrap,10001,"system:kubelet-bootstrap"
```

第一列：随机字符串，自己可生成  
第二列：用户名  
第三列：UID  
第四列：用户组  

创建apiserver配置文件：  

```
$ vi /opt/kubernetes/cfg/kube-apiserver 

KUBE_APISERVER_OPTS="--logtostderr=true \
--v=4 \
--etcd-servers=https://192.168.56.201:2379,https://192.168.56.202:2379,https://192.168.56.203:2379 \
--bind-address=192.168.56.201 \
--secure-port=6443 \
--advertise-address=192.168.56.201 \
--allow-privileged=true \
--service-cluster-ip-range=10.0.0.0/24 \
--enable-admission-plugins=NamespaceLifecycle,LimitRanger,SecurityContextDeny,ServiceAccount,ResourceQuota,NodeRestriction \
--authorization-mode=RBAC,Node \
--enable-bootstrap-token-auth \
--token-auth-file=/opt/kubernetes/cfg/token.csv \
--service-node-port-range=30000-50000 \
--tls-cert-file=/opt/kubernetes/ssl/server.pem  \
--tls-private-key-file=/opt/kubernetes/ssl/server-key.pem \
--client-ca-file=/opt/kubernetes/ssl/ca.pem \
--service-account-key-file=/opt/kubernetes/ssl/ca-key.pem \
--etcd-cafile=/opt/etcd/ssl/ca.pem \
--etcd-certfile=/opt/etcd/ssl/server.pem \
--etcd-keyfile=/opt/etcd/ssl/server-key.pem"
```

配置好前面生成的证书，确保能连接etcd。  

参数说明：

- --logtostderr 启用日志
- --v 日志等级
- --etcd-servers etcd集群地址
- --bind-address 监听地址
- --secure-port https安全端口
- --advertise-address 集群通告地址
- --allow-privileged 启用授权
- --service-cluster-ip-range Service虚拟IP地址段
- --enable-admission-plugins 准入控制模块
- --authorization-mode 认证授权，启用RBAC授权和节点自管理
- --enable-bootstrap-token-auth 启用TLS bootstrap功能，后面会讲到
- --token-auth-file token文件
- --service-node-port-range Service Node类型默认分配端口范围

systemd管理apiserver：  

```
$ vi /usr/lib/systemd/system/kube-apiserver.service 
[Unit]
Description=Kubernetes API Server
Documentation=https://github.com/kubernetes/kubernetes

[Service]
EnvironmentFile=-/opt/kubernetes/cfg/kube-apiserver
ExecStart=/opt/kubernetes/bin/kube-apiserver $KUBE_APISERVER_OPTS
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

启动  

```
$ systemctl daemon-reload
$ systemctl enable kube-apiserver
$ systemctl restart kube-apiserver
```


#### 部署scheduler组件

创建schduler配置文件

```
$ vi /opt/kubernetes/cfg/kube-scheduler 

KUBE_SCHEDULER_OPTS="--logtostderr=true \
--v=4 \
--master=127.0.0.1:8080 \
--leader-elect"
```

参数说明：
- --master 连接本地apiserver
- --leader-elect 当该组件启动多个时，自动选举（HA）

systemd管理schduler组件  

```go
$ vi /usr/lib/systemd/system/kube-scheduler.service 
[Unit]
Description=Kubernetes Scheduler
Documentation=https://github.com/kubernetes/kubernetes

[Service]
EnvironmentFile=-/opt/kubernetes/cfg/kube-scheduler
ExecStart=/opt/kubernetes/bin/kube-scheduler $KUBE_SCHEDULER_OPTS
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

启动  

```
$ systemctl daemon-reload
$ systemctl enable kube-scheduler
$ systemctl restart kube-scheduler
```

#### 部署controller-manager组件

创建controller-manager配置文件  

```go
$ vi /opt/kubernetes/cfg/kube-controller-manager 
KUBE_CONTROLLER_MANAGER_OPTS="--logtostderr=true \
--v=4 \
--master=127.0.0.1:8080 \
--leader-elect=true \
--address=127.0.0.1 \
--service-cluster-ip-range=10.0.0.0/24 \
--cluster-name=kubernetes \
--cluster-signing-cert-file=/opt/kubernetes/ssl/ca.pem \
--cluster-signing-key-file=/opt/kubernetes/ssl/ca-key.pem  \
--root-ca-file=/opt/kubernetes/ssl/ca.pem \
--service-account-private-key-file=/opt/kubernetes/ssl/ca-key.pem"
```

systemd管理controller-manager组件  

```go
$ vi /usr/lib/systemd/system/kube-controller-manager.service 
[Unit]
Description=Kubernetes Controller Manager
Documentation=https://github.com/kubernetes/kubernetes

[Service]
EnvironmentFile=-/opt/kubernetes/cfg/kube-controller-manager
ExecStart=/opt/kubernetes/bin/kube-controller-manager $KUBE_CONTROLLER_MANAGER_OPTS
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

启动 

```go
$ systemctl daemon-reload
$ systemctl enable kube-controller-manager
$ systemctl restart kube-controller-manager
```

安装完成查看当前的组件的状态  

```
$ /opt/kubernetes/bin/kubectl get cs
NAME                 STATUS    MESSAGE              ERROR
scheduler            Healthy   ok                   
controller-manager   Healthy   ok                   
etcd-0               Healthy   {"health": "true"}   
etcd-2               Healthy   {"health": "true"}   
etcd-1               Healthy   {"health": "true"} 
```

出现Healthy表示组件都正常  

### 在Node节点部署组件

将kubelet-bootstrap用户绑定到系统集群角色

```
$ kubectl create clusterrolebinding kubelet-bootstrap \
  --clusterrole=system:node-bootstrapper \
  --user=kubelet-bootstrap
```

#### 创建kubeconfig文件

在生成kubernetes证书的目录下执行以下命令生成kubeconfig文件,在master节点生成`bootstrap.kubeconfig  kube-proxy.kubeconfig`然后scp到node节点  

可以借助一个sh来实现,目录`/opt/kubernetes/cfg`
```
$ vi kubeconfig.sh 

# 创建 TLS Bootstrapping Token
#BOOTSTRAP_TOKEN=$(head -c 16 /dev/urandom | od -An -t x | tr -d ' ')
BOOTSTRAP_TOKEN=674c457d4dcf2eefe4920d7dbb6b0ddc

cat > token.csv <<EOF
${BOOTSTRAP_TOKEN},kubelet-bootstrap,10001,"system:kubelet-bootstrap"
EOF

#----------------------

APISERVER=192.168.56.201
SSL_DIR=/opt/kubernetes/ssl

# 创建kubelet bootstrapping kubeconfig 
export KUBE_APISERVER="https://$APISERVER:6443"

# 设置集群参数
kubectl config set-cluster kubernetes \
  --certificate-authority=$SSL_DIR/ca.pem \
  --embed-certs=true \
  --server=${KUBE_APISERVER} \
  --kubeconfig=bootstrap.kubeconfig

# 设置客户端认证参数
kubectl config set-credentials kubelet-bootstrap \
  --token=${BOOTSTRAP_TOKEN} \
  --kubeconfig=bootstrap.kubeconfig

# 设置上下文参数
kubectl config set-context default \
  --cluster=kubernetes \
  --user=kubelet-bootstrap \
  --kubeconfig=bootstrap.kubeconfig

# 设置默认上下文
kubectl config use-context default --kubeconfig=bootstrap.kubeconfig

#----------------------

# 创建kube-proxy kubeconfig文件

kubectl config set-cluster kubernetes \
  --certificate-authority=$SSL_DIR/ca.pem \
  --embed-certs=true \
  --server=${KUBE_APISERVER} \
  --kubeconfig=kube-proxy.kubeconfig

kubectl config set-credentials kube-proxy \
  --client-certificate=$SSL_DIR/kube-proxy.pem \
  --client-key=$SSL_DIR/kube-proxy-key.pem \
  --embed-certs=true \
  --kubeconfig=kube-proxy.kubeconfig

kubectl config set-context default \
  --cluster=kubernetes \
  --user=kube-proxy \
  --kubeconfig=kube-proxy.kubeconfig

kubectl config use-context default --kubeconfig=kube-proxy.kubeconfig
```

其中BOOTSTRAP_TOKEN就是上文我们设置的`/opt/kubernetes/cfg/token.csv`中的token  

生成 

``
$ bash kubeconfig.sh
``

```go
$ ls
bootstrap.kubeconfig  kube-proxy.kubeconfig
```

#### 部署kubelet组件

需要在每个node节点部署  

把前面在master节点下载的`kubernetes-server-linux-amd64.tar.gz`中的二进制包中的kubelet和kube-proxy拷贝到/opt/kubernetes/bin目录下。  

创建kubelet配置文件  

```
$ vi /opt/kubernetes/cfg/kubelet
KUBELET_OPTS="--logtostderr=true \
--v=4 \
--hostname-override=192.168.56.202 \
--kubeconfig=/opt/kubernetes/cfg/kubelet.kubeconfig \
--bootstrap-kubeconfig=/opt/kubernetes/cfg/bootstrap.kubeconfig \
--config=/opt/kubernetes/cfg/kubelet.config \
--cert-dir=/opt/kubernetes/ssl \
--pod-infra-container-image=registry.cn-hangzhou.aliyuncs.com/google-containers/pause-amd64:3.0"
```

参数说明

- --hostname-override 在集群中显示的主机名
- --kubeconfig 指定kubeconfig文件位置，会自动生成
- --bootstrap-kubeconfig 指定刚才生成的bootstrap.kubeconfig文件
- --cert-dir 颁发证书存放位置
- --pod-infra-container-image 管理Pod网络的镜像

kubelet.config的配制

```
$ vi /opt/kubernetes/cfg/kubelet.config

kind: KubeletConfiguration
apiVersion: kubelet.config.k8s.io/v1beta1
address: 192.168.56.202
port: 10250
readOnlyPort: 10255
cgroupDriver: cgroupfs
clusterDNS: ["10.0.0.2"]
clusterDomain: cluster.local.
failSwapOn: false
authentication:
  anonymous:
    enabled: true 
```

systemd管理kubelet组件  

```
$ vi /usr/lib/systemd/system/kubelet.service 
[Unit]
Description=Kubernetes Kubelet
After=docker.service
Requires=docker.service

[Service]
EnvironmentFile=/opt/kubernetes/cfg/kubelet
ExecStart=/opt/kubernetes/bin/kubelet $KUBELET_OPTS
Restart=on-failure
KillMode=process

[Install]
WantedBy=multi-user.target
```

启动

```go
$ systemctl daemon-reload
$ systemctl enable kubelet
$ systemctl restart kubelet
```

在Master审批Node加入集群  

master节点会收到node的验证请求，我们需要在master节点approve下  

```go
$ kubectl get csr
$ kubectl certificate approve XXXXID
$ kubectl get node
```

查看

```go
$ kubectl get csr
NAME                                                   AGE       REQUESTOR           CONDITION
node-csr-gE0iy6gY71RqRlC1ZhGGnvLwKBjLGmnTNmEoFj51yU4   26s       kubelet-bootstrap   Pending
node-csr-jFvVR_RBv-swHScxKEf_yDt_J72twIBTeslF8Bv18LQ   24s       kubelet-bootstrap   Pending
```

Pending状态的需要approve

```go
$ kubectl certificate approve node-csr-gE0iy6gY71RqRlC1ZhGGnvLwKBjLGmnTNmEoFj51yU4
```

#### 部署kube-proxy组件

每个node节点都需要执行  

创建kube-proxy配置文件  

```go
$ vi /opt/kubernetes/cfg/kube-proxy
KUBE_PROXY_OPTS="--logtostderr=true \
--v=4 \
--hostname-override=192.168.56.102 \
--cluster-cidr=10.0.0.0/24 \
--kubeconfig=/opt/kubernetes/cfg/kube-proxy.kubeconfig"
```

systemd管理kube-proxy组件  

```
$ vi /usr/lib/systemd/system/kube-proxy.service 
[Unit]
Description=Kubernetes Proxy
After=network.target

[Service]
EnvironmentFile=-/opt/kubernetes/cfg/kube-proxy
ExecStart=/opt/kubernetes/bin/kube-proxy $KUBE_PROXY_OPTS
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

启动

```
$ systemctl daemon-reload
$ systemctl enable kube-proxy
$ systemctl restart kube-proxy
```

### 查看集群状态

```
$ kubectl get node
NAME             STATUS   ROLES    AGE    VERSION
192.168.56.202   Ready    <none>   18h    v1.18.12
192.168.56.203   Ready    <none>   171m   v1.18.12

$ kubectl get cs
NAME                 STATUS    MESSAGE              ERROR
controller-manager   Healthy   ok                   
scheduler            Healthy   ok                   
etcd-0               Healthy   {"health": "true"}   
etcd-2               Healthy   {"health": "true"}   
etcd-1               Healthy   {"health": "true"}   
```

### 测试

创建测试文件  

```
$ vi /opt/kubernetes/demo/nginx-ds.yml 

apiVersion: v1
kind: Service
metadata:
  name: nginx-ds
  labels:
    app: nginx-ds
spec:
  type: NodePort
  selector:
    app: nginx-ds
  ports:
  - name: http
    port: 80
    targetPort: 80
---

apiVersion: apps/v1 # for versions before 1.9.0 use apps/v1beta2
kind: Deployment
metadata:
  name: nginx-ds
spec:
  selector:
    matchLabels:
      app: nginx-ds
  replicas: 3 # tells deployment to run 2 pods matching the template
  template:
    metadata:
      labels:
        app: nginx-ds
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80

```

运行

```
$ kubectl apply -f nginx-ds.yml 
```

结果

```
$  kubectl get pod
NAME                        READY   STATUS              RESTARTS   AGE
nginx-ds-76d6f5ffdd-2hds5   0/1     ContainerCreating   0          89s
nginx-ds-76d6f5ffdd-gv87x   0/1     ContainerCreating   0          16s
nginx-ds-76d6f5ffdd-m5dp5   0/1     ContainerCreating   0          89s
```

过一会就会自动创建成功

```
$  kubectl get pods -o wide|grep nginx-ds
nginx-ds-76d6f5ffdd-2hds5           1/1     Running   0          34m   172.17.21.4   192.168.56.203   <none>           <none>
nginx-ds-76d6f5ffdd-gv87x           1/1     Running   0          34m   172.17.46.4   192.168.56.202   <none>           <none>
nginx-ds-76d6f5ffdd-m5dp5           1/1     Running   0          34m   172.17.46.3   192.168.56.202   <none>           <none>

$  kubectl get svc |grep nginx-ds
nginx-ds     NodePort    10.0.0.112   <none>        80:41675/TCP   35m

$ curl 192.168.56.203:41675/

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

已经可以访问了   

### 参考
【二进制安装部署kubernetes集群---超详细教程】https://www.cnblogs.com/along21/p/10044931.html  
【etcd时间同步】https://bingohuang.com/etcd-operation-2/  
【Kubernetes v1.12/v1.13 二进制部署集群（HTTPS+RBAC）】https://blog.51cto.com/lizhenliang/2325770  