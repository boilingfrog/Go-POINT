## k8s用到的命令  

- [查看pod的状态](#%e6%9f%a5%e7%9c%8bpod%e7%9a%84%e7%8a%b6%e6%80%81)
- [查看pod的详情](#%e6%9f%a5%e7%9c%8bpod%e7%9a%84%e8%af%a6%e6%83%85)
- [修改host名字](#%e4%bf%ae%e6%94%b9host%e5%90%8d%e5%ad%97)
- [重启pod](#%e9%87%8d%e5%90%afpod)
  - [使用yaml文件](#%e4%bd%bf%e7%94%a8yaml%e6%96%87%e4%bb%b6)
- [重启pod](#%e5%a4%8d%e5%88%b6Slice%e5%92%8cMap%e6%b3%a8%e6%84%8f%e4%ba%8b%e9%a1%b9)

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
### 上传文件到