<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [k8s发布go应用](#k8s%E5%8F%91%E5%B8%83go%E5%BA%94%E7%94%A8)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [部署](#%E9%83%A8%E7%BD%B2)
    - [镜像打包](#%E9%95%9C%E5%83%8F%E6%89%93%E5%8C%85)
    - [编写yaml文件](#%E7%BC%96%E5%86%99yaml%E6%96%87%E4%BB%B6)
  - [配置ingress](#%E9%85%8D%E7%BD%AEingress)
  - [配置ingress转发策略](#%E9%85%8D%E7%BD%AEingress%E8%BD%AC%E5%8F%91%E7%AD%96%E7%95%A5)
  - [添加本机的host](#%E6%B7%BB%E5%8A%A0%E6%9C%AC%E6%9C%BA%E7%9A%84host)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## k8s发布go应用

### 前言




搭建了一套K8s,尝试发布一个go应用

### 部署

#### 镜像打包

之前已经打包过一个go的镜像了，这次就直接跳过了，打包记录`https://www.cnblogs.com/ricklz/p/12860434.html`  

#### 编写yaml文件

```go
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-app
spec:
  replicas: 2
  selector:
    matchLabels:
      app: go-app
  template:
    metadata:
      labels:
        app: go-app
    spec:
      containers:
        - name: go-app-container
          image: liz2019/test-docker-go-hub
          resources:
            limits:
              memory: "128Mi"
              cpu: "500m"
          ports:
            - containerPort: 8000
```

启动

```
$ kubectl apply -f kube-go.yaml

NAME                                READY   STATUS              RESTARTS   AGE
go-app-6498bff568-8n72g             0/1     ContainerCreating   0          62s
go-app-6498bff568-ncrmf             0/1     ContainerCreating   0          62s
```

暴露应用

```
$ kubectl expose deployment go-app --type=NodePort --name=go-app-svc --target-port=8000
```
查看

```go
$ kubectl get pods -o wide|grep go-app
go-app-58c4ff7448-j9j6m             1/1     Running   0          106s   172.17.32.5   192.168.56.202   <none>           <none>
go-app-58c4ff7448-p7xnj             1/1     Running   0          106s   172.17.65.4   192.168.56.203   <none>           <none>

$ kubectl get svc|grep go-app
go-app-svc   NodePort    10.0.0.247   <none>        8000:35100/TCP   93s
```

### 配置ingress

mandatory.yaml地址`https://github.com/boilingfrog/daily-test/blob/master/k8s/ingress/mandatory.yaml`

创建

```
$ kubectl apply -f ./mandatory.yaml
$ kubectl get pods -n ingress-nginx
$ kubectl get service -n ingress-nginx
```

### 配置ingress转发策略
首先查看service

```
]# kubectl get svc
NAME         TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)          AGE
go-app-svc   NodePort    10.0.0.247   <none>        8000:35100/TCP   2d13h
```

配置ingress

```
$ cat ingress-go.yaml 

apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: web-ingress
spec:
  rules:
  - host: www.liz-test.com
    http:
      paths:
      - backend:
          serviceName: go-app-svc
          servicePort: 8000
```


### 添加本机的host

查看

```
$ kubectl get pod -n ingress-nginx -o wide
NAME                                        READY   STATUS    RESTARTS   AGE   IP            NODE             NOMINATED NODE   READINESS GATES
default-http-backend-649b76fd96-v8s9v       0/1     Error     0          14m   <none>        192.168.56.202   <none>           <none>
nginx-ingress-controller-5b988cc9c6-kh2dg   0/1     Running   1          14m   172.17.32.5   192.168.56.202   <none>           <none>
```

nginx-ingress-controller是在192.168.56.202上面的，所以我们下面的host就配置到这个机器中。  


```
$ sudo vi /etc/hosts

// 根据ingress部署的iP
192.168.56.202 www.liz-test.com
```

访问结果

![channel](/img/ingress_6.png?raw=true)


### 参考
【K8S 安装 Ingress】https://www.jianshu.com/p/4370c00c040a  
【k8s ingress原理及ingress-nginx部署测试】https://segmentfault.com/a/1190000019908991  
【Kubernetes Deployment 故障排查常见方法[译]】https://www.qikqiak.com/post/troubleshooting-deployments/  
【kubernetes部署-ingress】https://blog.csdn.net/u013726175/article/details/88177110  