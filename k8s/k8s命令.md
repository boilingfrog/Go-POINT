## k8s用到的命令  

查看pod的状态  
````
[root@k8s-master ~]# kubectl get pod
NAME                     READY   STATUS              RESTARTS   AGE
nginx-554b9c67f9-dc68t   0/1     ContainerCreating   0          65s
````
查看pod的详情
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