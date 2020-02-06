## k8s用到的命令  

- [查看pod的状态](#%e6%9f%a5%e7%9c%8bpod%e7%9a%84%e7%8a%b6%e6%80%81)
- [查看pod的详情](#%e6%9f%a5%e7%9c%8bpod%e7%9a%84%e8%af%a6%e6%83%85)
- [修改host名字](#%e4%bf%ae%e6%94%b9host%e5%90%8d%e5%ad%97)
- [重启pod](#%e9%87%8d%e5%90%afpod)
  - [使用yaml文件](#%e4%bd%bf%e7%94%a8yaml%e6%96%87%e4%bb%b6)
- [本地文件上传到linux服务器上](#%e6%9c%ac%e5%9c%b0%e6%96%87%e4%bb%b6%e4%b8%8a%e4%bc%a0%e5%88%b0linux%e6%9c%8d%e5%8a%a1%e5%99%a8%e4%b8%8a)
   - [从服务器上下载文件](#%e4%bb%8e%e6%9c%8d%e5%8a%a1%e5%99%a8%e4%b8%8a%e4%b8%8b%e8%bd%bd%e6%96%87%e4%bb%b6)
   - [上传本地文件到服务器](#%e4%b8%8a%e4%bc%a0%e6%9c%ac%e5%9c%b0%e6%96%87%e4%bb%b6%e5%88%b0%e6%9c%8d%e5%8a%a1%e5%99%a8)
   - [从服务器下载整个目录](#%e4%bb%8e%e6%9c%8d%e5%8a%a1%e5%99%a8%e4%b8%8b%e8%bd%bd%e6%95%b4%e4%b8%aa%e7%9b%ae%e5%bd%95)
   - [上传目录到服务器](#%e4%b8%8a%e4%bc%a0%e7%9b%ae%e5%bd%95%e5%88%b0%e6%9c%8d%e5%8a%a1%e5%99%a8)

### 查看pod的状态  
````
[root@k8s-master ~]# kubectl get pod
NAME                     READY   STATUS              RESTARTS   AGE
nginx-554b9c67f9-dc68t   0/1     ContainerCreating   0          65s
````
### 查看pod的详情
````
[root@k8s-master ~]# kubectl describe pod nginx
Name:           nginx-554b9c67f9-5jrsw
Namespace:      default
Priority:       0
Node:           k8s-node1/192.168.31.191
Start Time:     Tue, 04 Feb 2020 15:25:04 +0800
Labels:         app=nginx
                pod-template-hash=554b9c67f9
Annotations:    <none>
Status:         Pending
IP:             
Controlled By:  ReplicaSet/nginx-554b9c67f9
Containers:
  nginx:
    Container ID:   
    Image:          nginx
    Image ID:       
    Port:           <none>
    Host Port:      <none>
    State:          Waiting
      Reason:       ContainerCreating
    Ready:          False
    Restart Count:  0
    Environment:    <none>
    Mounts:
      /var/run/secrets/kubernetes.io/serviceaccount from default-token-2r4rs (ro)
Conditions:
  Type              Status
  Initialized       True 
  Ready             False 
  ContainersReady   False 
  PodScheduled      True 
Volumes:
  default-token-2r4rs:
    Type:        Secret (a volume populated by a Secret)
    SecretName:  default-token-2r4rs
    Optional:    false
QoS Class:       BestEffort
Node-Selectors:  <none>
Tolerations:     node.kubernetes.io/not-ready:NoExecute for 300s
                 node.kubernetes.io/unreachable:NoExecute for 300s
Events:
  Type    Reason     Age    From                Message
  ----    ------     ----   ----                -------
  Normal  Scheduled  2m53s  default-scheduler   Successfully assigned default/nginx-554b9c67f9-5jrsw to k8s-node1
  Normal  Pulling    2m52s  kubelet, k8s-node1  Pulling image "nginx"
````
### 修改host名字
````
hostname XXXXX
````
修改之后退出重新进入就可以了

### 重启pod
#### 使用yaml文件
在有 yaml 文件的情况下可以直接使用kubectl replace --force -f xxxx.yaml 来强制替
换Pod 的 API 对象，从而达到重启的目的。  
````
# kubectl replace --force -f kubernetes-dashboard.yaml 
secret "kubernetes-dashboard-certs" deleted
serviceaccount "kubernetes-dashboard" deleted
role.rbac.authorization.k8s.io "kubernetes-dashboard-minimal" deleted
rolebinding.rbac.authorization.k8s.io "kubernetes-dashboard-minimal" deleted
deployment.apps "kubernetes-dashboard" deleted
service "kubernetes-dashboard" deleted
secret/kubernetes-dashboard-certs replaced
serviceaccount/kubernetes-dashboard replaced
role.rbac.authorization.k8s.io/kubernetes-dashboard-minimal replaced
rolebinding.rbac.authorization.k8s.io/kubernetes-dashboard-minimal replaced
deployment.apps/kubernetes-dashboard replaced
^[[Aservice/kubernetes-dashboard replaced
````
### 本地文件上传到linux服务器上
scp -P 端口 c://xxxx.txt user@ip:/home/root
注意：  
-P 大写  
-i 公钥  
#### 从服务器上下载文件
````
scp username@servername:/path/filename /var/www/local_dir（本地目录）  
例如scp root@192.168.0.101:/var/www/test.txt  把192.168.0.101上的/var/www/test.txt 的文件下载到/var/www/local_dir（本地目录）  
````
#### 上传本地文件到服务器
````
scp /path/filename username@servername:/path   

例如scp /var/www/test.php  root@192.168.0.101:/var/www/  把本机/var/www/目录下的test.php文件上传到192.168.0.101这台服务器上的/var/www/目录中
````
#### 从服务器下载整个目录
````
scp -r username@servername:/var/www/remote_dir/（远程目录） /var/www/local_dir（本地目录）
````
#### 上传目录到服务器
````
scp  -r local_dir username@servername:remote_dir
例如：scp -r test  root@192.168.0.101:/var/www/   把当前目录下的test目录上传到服务器的/var/www/ 目录
````