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








### 参考
【二进制安装部署kubernetes集群---超详细教程】https://www.cnblogs.com/along21/p/10044931.html