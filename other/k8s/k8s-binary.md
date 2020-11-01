<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [二进制部署k8s](#%E4%BA%8C%E8%BF%9B%E5%88%B6%E9%83%A8%E7%BD%B2k8s)
  - [秘钥免密码](#%E7%A7%98%E9%92%A5%E5%85%8D%E5%AF%86%E7%A0%81)
  - [docker安装](#docker%E5%AE%89%E8%A3%85)
  - [部署的命令](#%E9%83%A8%E7%BD%B2%E7%9A%84%E5%91%BD%E4%BB%A4)
  - [安装etcd](#%E5%AE%89%E8%A3%85etcd)
    - [创建证书](#%E5%88%9B%E5%BB%BA%E8%AF%81%E4%B9%A6)
    - [生成证书](#%E7%94%9F%E6%88%90%E8%AF%81%E4%B9%A6)
    - [部署Etcd](#%E9%83%A8%E7%BD%B2etcd)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 二进制部署k8s


### 秘钥免密码

```go
ssh-copy-id root@192.168.56.101
```

### docker安装

docker 安装
```go
// 安装
yum -y install docker

// 设置开机启动
sudo systemctl enable docker

// 启动docker 
 sudo systemctl start docker
```

docker-compose安装  

```go
sudo curl -L "https://github.com/docker/compose/releases/download/1.27.4/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose

// 将可执行权限应用于二进制文件
sudo chmod +x /usr/local/bin/docker-compose

// 创建软连接
sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose
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

```go
wget https://pkg.cfssl.org/R1.2/cfssl_linux-amd64
wget https://pkg.cfssl.org/R1.2/cfssljson_linux-amd64
wget https://pkg.cfssl.org/R1.2/cfssl-certinfo_linux-amd64
chmod +x cfssl_linux-amd64 cfssljson_linux-amd64 cfssl-certinfo_linux-amd64
mv cfssl_linux-amd64 /usr/local/bin/cfssl
mv cfssljson_linux-amd64 /usr/local/bin/cfssljson
mv cfssl-certinfo_linux-amd64 /usr/bin/cfssl-certinfo
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

```go
cfssl gencert -initca ca-csr.json | cfssljson -bare ca -
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=www server-csr.json | cfssljson -bare server
# ls *pem
ca-key.pem  ca.pem  server-key.pem  server.pem
```

移动证书文件

```go
cp ca*pem server*pem /opt/etcd/ssl
```

#### 部署Etcd

```go
wget https://github.com/etcd-io/etcd/releases/download/v3.2.12/etcd-v3.2.12-linux-amd64.tar.gz

mkdir /opt/etcd/{bin,cfg,ssl} -p
tar zxvf etcd-v3.2.12-linux-amd64.tar.gz
mv etcd-v3.2.12-linux-amd64/{etcd,etcdctl} /opt/etcd/bin/
```

创建配置文件

```go
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











### 参考
【二进制安装部署kubernetes集群---超详细教程】https://www.cnblogs.com/along21/p/10044931.html  
【etcd时间同步】https://bingohuang.com/etcd-operation-2/  
【Kubernetes v1.12/v1.13 二进制部署集群（HTTPS+RBAC）】https://blog.51cto.com/lizhenliang/2325770  