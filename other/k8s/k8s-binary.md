<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


  - [二进制部署k8s](#%E4%BA%8C%E8%BF%9B%E5%88%B6%E9%83%A8%E7%BD%B2k8s)
    - [准备工作](#%E5%87%86%E5%A4%87%E5%B7%A5%E4%BD%9C)
      - [关闭防火墙](#%E5%85%B3%E9%97%AD%E9%98%B2%E7%81%AB%E5%A2%99)
      - [关闭 swap 分区](#%E5%85%B3%E9%97%AD-swap-%E5%88%86%E5%8C%BA)
      - [关闭 SELinux](#%E5%85%B3%E9%97%AD-selinux)
      - [更新系统时间](#%E6%9B%B4%E6%96%B0%E7%B3%BB%E7%BB%9F%E6%97%B6%E9%97%B4)
    - [秘钥免密码](#%E7%A7%98%E9%92%A5%E5%85%8D%E5%AF%86%E7%A0%81)
    - [docker安装](#docker%E5%AE%89%E8%A3%85)
    - [部署的命令](#%E9%83%A8%E7%BD%B2%E7%9A%84%E5%91%BD%E4%BB%A4)
    - [安装etcd](#%E5%AE%89%E8%A3%85etcd)
      - [创建证书](#%E5%88%9B%E5%BB%BA%E8%AF%81%E4%B9%A6)
      - [生成证书](#%E7%94%9F%E6%88%90%E8%AF%81%E4%B9%A6)
      - [部署Etcd](#%E9%83%A8%E7%BD%B2etcd)
    - [在Node安装Docker](#%E5%9C%A8node%E5%AE%89%E8%A3%85docker)
    - [Flannel网络](#flannel%E7%BD%91%E7%BB%9C)
- [cat /usr/lib/systemd/system/flanneld.service](#cat-usrlibsystemdsystemflanneldservice)
    - [配置](#%E9%85%8D%E7%BD%AE)
    - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 二进制部署k8s

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
$ ssh-copy-id root@192.168.56.101
```

### docker安装

docker 安装
```
// 安装
$ yum -y install docker

// 设置开机启动
$ sudo systemctl enable docker

// 启动docker 
$ sudo systemctl start docker
```

docker-compose安装  

```
$ sudo curl -L "https://github.com/docker/compose/releases/download/1.27.4/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose

// 将可执行权限应用于二进制文件
$ sudo chmod +x /usr/local/bin/docker-compose

// 创建软连接
$ sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose
```


### 部署的命令


192.168.56.101 kube-master

192.168.56.102 kube-node1

192.168.56.103 kube-node2

### 安装etcd

集群中每台机器都需要安装，先在一台节点安装配置，之后把配置文件scp到其他机器就行了  

#### 创建证书

使用cfssl来生成自签证书  

安装

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

创建需要的文件

```
# cat ca-config.json
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

# cat ca-csr.json
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

# cat server-csr.json
{
    "CN": "etcd",
    "hosts": [
    "192.168.56.101",
    "192.168.56.102",
    "192.168.56.103"
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

移动证书文件

```
$ cp ca*pem server*pem /opt/etcd/ssl
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
# cat /opt/etcd/cfg/etcd   
#[Member]
ETCD_NAME="etcd01"
ETCD_DATA_DIR="/var/lib/etcd/default.etcd"
ETCD_LISTEN_PEER_URLS="https://192.168.56.101:2380"
ETCD_LISTEN_CLIENT_URLS="https://192.168.56.101:2379"

#[Clustering]
ETCD_INITIAL_ADVERTISE_PEER_URLS="https://192.168.56.101:2380"
ETCD_ADVERTISE_CLIENT_URLS="https://192.168.56.101:2379"
ETCD_INITIAL_CLUSTER="etcd01=https://192.168.56.101:2380,etcd02=https://192.168.56.102:2380,etcd03=https://192.168.56.103:2380"
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


```
#[Member]
ETCD_NAME="etcd01"
ETCD_DATA_DIR="/var/lib/etcd/default.etcd"
ETCD_LISTEN_PEER_URLS="https://192.168.56.101:2380"
ETCD_LISTEN_CLIENT_URLS="https://192.168.56.101:2379"

#[Clustering]
ETCD_INITIAL_ADVERTISE_PEER_URLS="https://192.168.56.101:2380"
ETCD_ADVERTISE_CLIENT_URLS="https://192.168.56.101:2379"
ETCD_INITIAL_CLUSTER="etcd01=https://192.168.56.101:2380,etcd02=https://192.168.56.102:2380,etcd03=https://192.168.56.103:2380"
ETCD_INITIAL_CLUSTER_TOKEN="etcd-cluster"
ETCD_INITIAL_CLUSTER_STATE="new"
```


```
--ca-file=ca.pem --cert-file=server.pem --key-file=server-key.pem \
--endpoints="https://192.168.56.101:2379,https://192.168.56.102:2379,https://192.168.56.103:2379" \
cluster-health



--ca-file=/opt/etcd/ssl/ca.pem --cert-file=/opt/etcd/ssl/server.pem --key-file=/opt/etcd/ssl/server-key.pem \
--endpoints="https://192.168.56.101:2379,https://192.168.56.102:2379,https://192.168.56.103:2379" \
cluster-health
```

启动

```
$ systemctl start etcd
$ systemctl enable etcd
```

查看状态 

```
$ journalctl -u etcd
```

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

配置子网
```go
/opt/etcd/bin/etcdctl \
--ca-file=/opt/etcd/ssl/ca.pem --cert-file=/opt/etcd/ssl/server.pem --key-file=/opt/etcd/ssl/server-key.pem \
--endpoints="https://192.168.56.201:2379,https://192.168.56.202:2379,https://192.168.56.203:2379" \
set /coreos.com/network/config  '{ "Network": "172.17.0.0/16", "Backend": {"Type": "vxlan"}}'
```

每个节点都需要操作  

```
$ wget https://github.com/coreos/flannel/releases/download/v0.9.1/flannel-v0.9.1-linux-amd64.tar.gz
$ tar zxvf flannel-v0.9.1-linux-amd64.tar.gz
$ mv flanneld mk-docker-opts.sh /opt/kubernetes/bin
```

配置Flannel

```
# cat /opt/kubernetes/cfg/flanneld
FLANNEL_OPTIONS="--etcd-endpoints=https://192.168.56.101:2379,https://192.168.56.102:2379,https://192.168.56.103:2379 -etcd-cafile=/opt/etcd/ssl/ca.pem -etcd-certfile=/opt/etcd/ssl/server.pem -etcd-keyfile=/opt/etcd/ssl/server-key.pem"
```

systemd管理Flannel 

````
# cat /usr/lib/systemd/system/flanneld.service
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
# cat /usr/lib/systemd/system/docker.service 

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
# ps -ef |grep docker
root      2859     1  0 18:13 ?        00:00:05 /usr/bin/dockerd-current --add-runtime docker-runc=/usr/libexec/docker/docker-runc-current --default-runtime=docker-runc --exec-opt native.cgroupdriver=systemd --userland-proxy-path=/usr/libexec/docker/docker-proxy-current --init-path=/usr/libexec/docker/docker-init-current --seccomp-profile=/etc/docker/seccomp.json --selinux-enabled --log-driver=journald --signature-verification=false --storage-driver overlay2 --bip=172.17.96.1/24 --ip-masq=false --mtu=1450
root      2865  2859  0 18:13 ?        00:00:02 /usr/bin/docker-containerd-current -l unix:///var/run/docker/libcontainerd/docker-containerd.sock --metrics-interval=0 --start-timeout 2m --state-dir /var/run/docker/libcontainerd/containerd --shim docker-containerd-shim --runtime docker-runc --runtime-args --systemd-cgroup=true
root      5799  1753  0 18:49 pts/0    00:00:00 grep --color=auto docker

# ip addr
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
    inet 192.168.56.103/24 brd 192.168.56.255 scope global noprefixroute enp0s8
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
# ping 172.17.72.0
PING 172.17.72.0 (172.17.72.0) 56(84) bytes of data.
64 bytes from 172.17.72.0: icmp_seq=1 ttl=64 time=1.25 ms
64 bytes from 172.17.72.0: icmp_seq=2 ttl=64 time=0.507 ms
64 bytes from 172.17.72.0: icmp_seq=3 ttl=64 time=1.17 ms
64 bytes from 172.17.72.0: icmp_seq=4 ttl=64 time=1.84 ms
64 bytes from 172.17.72.0: icmp_seq=5 ttl=64 time=0.426 ms
```

### 配置apiserver组件

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
$ cat /opt/kubernetes/cfg/token.csv
674c457d4dcf2eefe4920d7dbb6b0ddc,kubelet-bootstrap,10001,"system:kubelet-bootstrap"
```

第一列：随机字符串，自己可生成  
第二列：用户名  
第三列：UID  
第四列：用户组  

创建apiserver配置文件：  

```
KUBE_APISERVER_OPTS="--logtostderr=true \
--v=4 \
--etcd-servers=https://192.168.56.101:2379,https://192.168.56.102:2379,https://192.168.56.103:2379 \
--bind-address=192.168.56.101 \
--secure-port=6443 \
--advertise-address=192.168.56.101 \
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
$ cat /usr/lib/systemd/system/kube-apiserver.service 
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


```
{
    "CN": "kubernetes",
    "hosts": [
      "10.0.0.1",
      "127.0.0.1",
      "192.168.56.101",
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
```

创建kubelet配置文件：

```
KUBELET_OPTS="--logtostderr=true \
--v=4 \
--hostname-override=192.168.56.102 \
--kubeconfig=/opt/kubernetes/cfg/kubelet.kubeconfig \
--bootstrap-kubeconfig=/opt/kubernetes/cfg/bootstrap.kubeconfig \
--config=/opt/kubernetes/cfg/kubelet.config \
--cert-dir=/opt/kubernetes/ssl \
--pod-infra-container-image=registry.cn-hangzhou.aliyuncs.com/google-containers/pause-amd64:3.0"
```

参数说明：

- --hostname-override 在集群中显示的主机名
- --kubeconfig 指定kubeconfig文件位置，会自动生成
- --bootstrap-kubeconfig 指定刚才生成的bootstrap.kubeconfig文件
- --cert-dir 颁发证书存放位置
- --pod-infra-container-image 管理Pod网络的镜像

`kubelet.config`配置  
```
$ cat /opt/kubernetes/cfg/kubelet.config

kind: KubeletConfiguration
apiVersion: kubelet.config.k8s.io/v1beta1
address: 192.168.56.102
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
/opt/kubernetes/cfg/kubelet.config



```
KUBE_PROXY_OPTS="--logtostderr=true \
--v=4 \
--hostname-override=192.168.56.102 \
--cluster-cidr=10.0.0.0/24 \
--kubeconfig=/opt/kubernetes/cfg/kube-proxy.kubeconfig"
```



### 参考
【二进制安装部署kubernetes集群---超详细教程】https://www.cnblogs.com/along21/p/10044931.html  
【etcd时间同步】https://bingohuang.com/etcd-operation-2/  
【Kubernetes v1.12/v1.13 二进制部署集群（HTTPS+RBAC）】https://blog.51cto.com/lizhenliang/2325770  