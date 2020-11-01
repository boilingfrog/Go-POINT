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

### 创建证书



```
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