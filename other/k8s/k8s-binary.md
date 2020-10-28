## 二进制部署k8s


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


gpasswd -a k8s wheel Adding user k8s to group wheel




``````
cfssl gencert -ca=/opt/k8s/cert/ca.pem \
-ca-key=/opt/k8s/cert/ca-key.pem \
-config=/opt/k8s/cert/ca-config.json \
-profile=kubernetes admin-csr.json | cfssljson -bare admin
``````





NODE_IPS=("192.168.56.101" "192.168.56.102" "192.168.56.103")
for node_ip in ${NODE_IPS[@]};do
    echo ">>> ${node_ip}"
    scp /root/kubernetes/client/bin/kubectl k8s@${node_ip}:/opt/k8s/bin/
    ssh k8s@${node_ip} "chmod +x /opt/k8s/bin/*"
    ssh k8s@${node_ip} "mkdir -p ~/.kube"
    scp ~/.kube/config k8s@${node_ip}:~/.kube/config
    ssh root@${node_ip} "mkdir -p ~/.kube"
    scp ~/.kube/config root@${node_ip}:~/.kube/config
done



kubectl config set-cluster kubernetes \
--certificate-authority=/opt/k8s/cert/ca.pem \
--embed-certs=true \
--server=https://192.168.10.10:8443 \
--kubeconfig=/root/.kube/kubectl.kubeconfig




 kubectl config set-credentials kube-admin \
--client-certificate=/opt/k8s/cert/admin.pem \
--client-key=/opt/k8s/cert/admin-key.pem \
--embed-certs=true \
--kubeconfig=/root/.kube/kubectl.kubeconfig



kubectl config set-context kube-admin@kubernetes \
--cluster=kubernetes \
--user=kube-admin \
--kubeconfig=/root/.kube/kubectl.kubeconfig



bash /opt/k8s/script/kubectl_environment.sh 192.168.56.101 192.168.56.102 192.168.56.103





#### 部署etcd

etcd 是基于 Raft 的分布式 key-value 存储系统，由 CoreOS 开发，常用于服务发现、共享配置以及并发控制（如 leader 选举、分布式锁等）。
kubernetes 使用 etcd 存储所有运行数据。  

本文档介绍部署一个三节点高可用 etcd 集群的步骤：  

① 下载和分发 etcd 二进制文件  

② 创建 etcd 集群各节点的 x509 证书，用于加密客户端(如 etcdctl) 与 etcd 集群、etcd 集群之间的数据流；  

③ 创建 etcd 的 systemd unit 文件，配置服务参数；  

④ 检查集群工作状态；  


##### 下载etcd二进制文件

```
[root@kube-master ~]# wget https://github.com/coreos/etcd/releases/download/v3.3.7/etcd-v3.3.7-linux-amd64.tar.gz
[root@kube-master ~]# tar -xvf etcd-v3.3.7-linux-amd64.tar.gz
```




cat > etcd-csr.json <<EOF
{
    "CN": "etcd",
    "hosts": [
        "127.0.0.1",
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
            "ST": "BeiJing",
            "L": "BeiJing",
            "O": "k8s",
            "OU": "4Paradigm"
        }
    ]
}
EOF



```go
cfssl gencert -ca=/opt/k8s/cert/ca.pem \
-ca-key=/opt/k8s/cert/ca-key.pem \
-config=/opt/k8s/cert/ca-config.json \
-profile=kubernetes etcd-csr.json | cfssljson -bare etcd
```

修改权限

```go
chown -R k8s /opt/etcd/cert/
chmod +x -R /opt/etcd/cert/
```

#### 创建etcd证书和密匙

创建证书签名请求

````

````
NODE_IPS=("192.168.56.101" "192.168.56.102" "192.168.56.103")
for node_ip in ${NODE_IPS[@]};do
        echo ">>> ${node_ip}"
        scp /root/etcd-v3.3.7-linux-amd64/etcd* k8s@${node_ip}:/opt/k8s/bin
        ssh k8s@${node_ip} "chmod +x /opt/k8s/bin/*"
        ssh root@${node_ip} "mkdir -p /opt/etcd/cert && chown -R k8s /opt/etcd/cert"
        scp /opt/etcd/cert/etcd*.pem k8s@${node_ip}:/opt/etcd/cert/
done



NODE_NAMES=("etcd0" "etcd1" "etcd2")
NODE_IPS=("192.168.56.101" "192.168.56.102" "192.168.56.103")
#替换模板文件中的变量，为各节点创建 systemd unit 文件
for (( i=0; i < 3; i++ ));do
        sed -e "s/##NODE_NAME##/${NODE_NAMES[i]}/g" -e "s/##NODE_IP##/${NODE_IPS[i]}/g" /opt/etcd/etcd.service.template > /opt/etcd/etcd-${NODE_IPS[i]}.service
done
#分发生成的 systemd unit 和etcd的配置文件：
for node_ip in ${NODE_IPS[@]};do
        echo ">>> ${node_ip}"
        ssh root@${node_ip} "mkdir -p /opt/lib/etcd && chown -R k8s /opt/lib/etcd"
        scp /opt/etcd/etcd-${node_ip}.service root@${node_ip}:/etc/systemd/system/etcd.service
done



NODE_IPS=("192.168.56.101" "192.168.56.102" "192.168.56.103")
#启动 etcd 服务
for node_ip in ${NODE_IPS[@]};do
        echo ">>> ${node_ip}"
        ssh root@${node_ip} "systemctl daemon-reload && systemctl enable etcd && systemctl start etcd"
done
#检查启动结果,确保状态为 active (running)
for node_ip in ${NODE_IPS[@]};do
        echo ">>> ${node_ip}"
        ssh k8s@${node_ip} "systemctl status etcd|grep Active"
done
#验证服务状态,输出均为healthy 时表示集群服务正常
for node_ip in ${NODE_IPS[@]};do
        echo ">>> ${node_ip}"
        ETCDCTL_API=3 /opt/k8s/bin/etcdctl \
--endpoints=https://${node_ip}:2379 \
--cacert=/opt/k8s/cert/ca.pem \
--cert=/opt/etcd/cert/etcd.pem \
--key=/opt/etcd/cert/etcd-key.pem endpoint health
done 



cat > /opt/etcd/etcd.service.template <<EOF

[Unit]
Description=Etcd Server
After=network.target
After=network-online.target
Wants=network-online.target
Documentation=https://github.com/coreos
[Service]
User=k8s
Type=notify
WorkingDirectory=/opt/lib/etcd/
ExecStart=/opt/k8s/bin/etcd \
    --data-dir=/opt/lib/etcd \
    --name ##NODE_NAME## \
    --cert-file=/opt/etcd/cert/etcd.pem \
    --key-file=/opt/etcd/cert/etcd-key.pem \
    --trusted-ca-file=/opt/k8s/cert/ca.pem \
    --peer-cert-file=/opt/etcd/cert/etcd.pem \
    --peer-key-file=/opt/etcd/cert/etcd-key.pem \
    --peer-trusted-ca-file=/opt/k8s/cert/ca.pem \
    --peer-client-cert-auth \
    --client-cert-auth \
    --listen-peer-urls=https://##NODE_IP##:2380 \
    --initial-advertise-peer-urls=https://##NODE_IP##:2380 \
    --listen-client-urls=https://##NODE_IP##:2379,http://127.0.0.1:2379\
    --advertise-client-urls=https://##NODE_IP##:2379 \
    --initial-cluster-token=etcd-cluster-0 \
    --initial-cluster=etcd0=https://192.168.56.101:2380,etcd1=https://192.168.56.102:2380,etcd2=https://192.168.56.103:2380 \
    --initial-cluster-state=new
Restart=on-failure
RestartSec=5
LimitNOFILE=65536
[Install]
WantedBy=multi-user.target
EOF

flannel
```go
cfssl gencert -ca=/opt/k8s/cert/ca.pem \
-ca-key=/opt/k8s/cert/ca-key.pem \
-config=/opt/k8s/cert/ca-config.json \
-profile=kubernetes flanneld-csr.json | cfssljson -bare flanneld
```

```go
cat > /opt/k8s/script/scp_flannel.sh <<EOF

NODE_IPS=("192.168.56.101" "192.168.56.102" "192.168.56.103")
for node_ip in ${NODE_IPS[@]};do
    echo ">>> ${node_ip}"
    scp /root/flannel/{flanneld,mk-docker-opts.sh} k8s@${node_ip}:/opt/k8s/bin/
    ssh k8s@${node_ip} "chmod +x /opt/k8s/bin/*"
    ssh root@${node_ip} "mkdir -p /opt/flannel/cert && chown -R k8s /opt/flannel"
    scp /opt/flannel/cert/flanneld*.pem k8s@${node_ip}:/opt/flannel/cert
done
EOF
```


````go
etcdctl \
--endpoints="https://192.168.56.101:2379,https://192.168.56.102:2379,https://192.168.56.103:2379" \
--ca-file=/opt/k8s/cert/ca.pem \
--cert-file=/opt/flannel/cert/flanneld.pem \
--key-file=/opt/flannel/cert/flanneld-key.pem \
set /atomic.io/network/config '{"Network":"10.30.0.0/16","SubnetLen": 24, "Backend": {"Type": "vxlan"}}'
{"Network":"10.30.0.0/16","SubnetLen": 24, "Backend": {"Type": "vxlan"}}


````


cat > /opt/flannel/flanneld.service << EOF

[Unit]
Description=Flanneld overlay address etcd agent
After=network.target
After=network-online.target
Wants=network-online.target
After=etcd.service
Before=docker.service

[Service]
Type=notify
ExecStart=/opt/k8s/bin/flanneld \
-etcd-cafile=/opt/k8s/cert/ca.pem \
-etcd-certfile=/opt/flannel/cert/flanneld.pem \
-etcd-keyfile=/opt/flannel/cert/flanneld-key.pem \
-etcd-endpoints=https://192.168.156.101:2379,https://192.168.56.102:2379,https://192.168.56.103:2379 \
-etcd-prefix=/atomic.io/network \
-iface=eth1
ExecStartPost=/opt/k8s/bin/mk-docker-opts.sh -k DOCKER_NETWORK_OPTIONS -d /run/flannel/docker
Restart=on-failure

[Install]
WantedBy=multi-user.target
RequiredBy=docker.service
EOF



```go
NODE_IPS=("192.168.56.101" "192.168.56.102" "192.168.56.103")
for node_ip in ${NODE_IPS[@]};do
    echo ">>> ${node_ip}"
 #分发 flanneld systemd unit 文件到所有节点
    scp /opt/flannel/flanneld.service root@${node_ip}:/etc/systemd/system/
 #启动 flanneld 服务
    ssh root@${node_ip} "systemctl daemon-reload && systemctl enable flanneld && systemctl restart flanneld"
 #检查启动结果
    ssh k8s@${node_ip} "systemctl status flanneld|grep Active"
done
```



### 参考
【二进制安装部署kubernetes集群---超详细教程】https://www.cnblogs.com/along21/p/10044931.html  
【etcd时间同步】https://bingohuang.com/etcd-operation-2/