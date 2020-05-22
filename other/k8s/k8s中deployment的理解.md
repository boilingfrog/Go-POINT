## 理解Deployment

### pod与controllers的关系

controllers：在集群上管理和运行容器的  
通过label-selector相关联  
Pod通过控制器实现应用的运维，如伸缩，滚动升级等  

### Deployment功能与应用

部署无状态应用  
管理Pod和ReplicaSet  
具有上线部署、副本设定、滚动升级、回滚等功能  
提供声明式更新，例如只更新一个新的Image  
应用场景：Web服务，微服务

### 使用Deployment部署一个应用

创建
````
kubectl create deployment web --image=nginx:1.14 
kubectl get deploy,pods
````
发布
````
kubectl expose deployment web --port=80 --type=NodePort --target-port=80 --name=web
kubectl get service
````

### 升级与回滚

升级
````
kubectl set image deployment/web nginx=nginx:1.15
````
# 查看升级状态
````
kubectl rollout status deployment/web 
````
回滚
````
kubectl rollout history deployment/web
kubectl rollout undo deployment/web
kubectl rollout undo deployment/web --revision=2
````

### 应用弹性

````
$ kubectl scale deployment nginx-deployment --replicas=1
````