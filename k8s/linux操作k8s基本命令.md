## linux中操作k8s的基本命令 


最近工作中使用到了k8s，那么就来总结下平时使用到的基本的命令 
- [获取某个namespace下的pod](#%e8%8e%b7%e5%8f%96%e6%9f%90%e4%b8%aanamespace%e4%b8%8b%e7%9a%84pod)
- [获取某个namespace下的pod,展示出ip和pod信息](#%e8%8e%b7%e5%8f%96%e6%9f%90%e4%b8%aanamespace%e4%b8%8b%e7%9a%84pod%2c%e5%b1%95%e7%a4%ba%e5%87%baip%e5%92%8cpod%e4%bf%a1%e6%81%af)
- [查看节点控制台的日志](#%e6%9f%a5%e7%9c%8b%e8%8a%82%e7%82%b9%e6%8e%a7%e5%88%b6%e5%8f%b0%e7%9a%84%e6%97%a5%e5%bf%97)
 
### 获取某个namespace下的pod

kubectl get pods -n namespace

````
# kubectl get pods -n handle
NAME                              READY   STATUS    RESTARTS   AGE
access-7754f795dd-f267n           1/1     Running   0          2d12h
account-78fc5f5bf4-xb96c          1/1     Running   0          2d12h
admin-bd8d5f6bb-fkc4l             1/1     Running   0          2d12h
cores-77c5f6f696-k26hf            1/1     Running   0          2d12h
file-7b94fb9fb7-m6x4v             1/1     Running   0          2d12h
handle-55989bc69b-b2rp7           1/1     Running   0          2d12h
handleapp-fddcf85b8-dn7t2         1/1     Running   0          25d
index-5b87c9fd5b-q6htq            1/1     Running   0          2d12h
log-statistics-8697f4987b-ptqn6   1/1     Running   1          60d
notification-66b9ddd5c4-f2ktq     1/1     Running   0          2d12h
open-74554f48-rclh6               1/1     Running   0          2d12h
search-7d469f95fb-r29rw           1/1     Running   0          2d12h
sequence-7d5bf65f9d-zt7xh         1/1     Running   0          2d12h
````

### 获取某个namespace下的pod,展示出ip和pod信息
kubectl get pods --all-namespaces -o wide
````
# kubectl get pods -n handle -o wide
NAME                              READY   STATUS    RESTARTS   AGE     IP             NODE           NOMINATED NODE   READINESS GATES
access-7754f795dd-f267n           1/1     Running   0          2d12h   172.20.1.174   192.168.1.13   <none>           <none>
account-78fc5f5bf4-xb96c          1/1     Running   0          2d12h   172.20.1.179   192.168.1.13   <none>           <none>
admin-bd8d5f6bb-fkc4l             1/1     Running   0          2d12h   172.20.1.172   192.168.1.13   <none>           <none>
cores-77c5f6f696-k26hf            1/1     Running   0          2d12h   172.20.1.173   192.168.1.13   <none>           <none>
file-7b94fb9fb7-m6x4v             1/1     Running   0          2d12h   172.20.1.178   192.168.1.13   <none>           <none>
handle-55989bc69b-b2rp7           1/1     Running   0          2d12h   172.20.1.176   192.168.1.13   <none>           <none>
handleapp-fddcf85b8-dn7t2         1/1     Running   0          25d     172.20.1.113   192.168.1.13   <none>           <none>
index-5b87c9fd5b-q6htq            1/1     Running   0          2d12h   172.20.2.27    192.168.1.12   <none>           <none>
log-statistics-8697f4987b-ptqn6   1/1     Running   1          60d     172.20.0.72    192.168.1.11   <none>           <none>
notification-66b9ddd5c4-f2ktq     1/1     Running   0          2d12h   172.20.2.28    192.168.1.12   <none>           <none>
open-74554f48-rclh6               1/1     Running   0          2d12h   172.20.1.177   192.168.1.13   <none>           <none>
search-7d469f95fb-r29rw           1/1     Running   0          2d12h   172.20.1.180   192.168.1.13   <none>           <none>
sequence-7d5bf65f9d-zt7xh         1/1     Running   0          2d12h   172.20.1.175   192.168.1.13   <none>           <none>

```` 

### 查看节点控制台的日志
kubectl logs -f POD-NAME　-n namespace
````
# kubectl logs -f  handle-55989bc69b-b2rp7  -n handle
2020/01/10 21:06:45.667471 [INFO][dbcache] table: users, prefix: 9bc6
2020/01/10 21:06:45.667911 [INFO][dbcache] table: settings, prefix: 2e5d
2020/01/10 21:06:45.742152 [INFO][handle-55989bc69b-b2rp7]["pkg.jimu.io/libs/util/version.go:24"] version: 680a3e2
2020/01/10 21:06:45.746331 [INFO][handle-55989bc69b-b2rp7]["pkg.jimu.io/libs/prometheus/prometheus.go:29"] Metrics listening on :3001
2020/01/10 21:06:45.788389 [INFO][handle-55989bc69b-b2rp7]["pkg.jimu.io/vendor/github.com/teapots/teapot/app.go:130"] Teapot listening on 0.0.0.0:80 in [prod] mode
2020/01/10 21:06:55.341161 [INFO][644c924b853ac088][680a3e2] [REQ_BEG] PUT handle.dev.jimu.io/enterprises/2650675293668250675 192.168.1.122
```` 