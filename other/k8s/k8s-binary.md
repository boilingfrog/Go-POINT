## 二进制部署k8s


### docker安装

```go
// 安装
yum -y install docker

// 设置开机启动
sudo systemctl enable docker
```

docker-compose

```go
curl -L https://github.com/docker/compose/releases/download/1.23.2/docker-compose-`uname -s`-`uname -m` -o /usr/bin/docker-compose

```


### 部署的命令




192.168.56.101 kube-master

192.168.56.102 kube-node1

192.168.56.103 kube-node2


gpasswd -a k8s wheel Adding user k8s to group wheel





### 参考
【二进制安装部署kubernetes集群---超详细教程】https://www.cnblogs.com/along21/p/10044931.html