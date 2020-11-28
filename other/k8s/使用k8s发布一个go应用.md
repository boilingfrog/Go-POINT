<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [k8s发布go应用](#k8s%E5%8F%91%E5%B8%83go%E5%BA%94%E7%94%A8)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [部署](#%E9%83%A8%E7%BD%B2)
    - [镜像打包](#%E9%95%9C%E5%83%8F%E6%89%93%E5%8C%85)
    - [编写yaml文件](#%E7%BC%96%E5%86%99yaml%E6%96%87%E4%BB%B6)
  - [配置service](#%E9%85%8D%E7%BD%AEservice)
  - [配置ingress](#%E9%85%8D%E7%BD%AEingress)
  - [配置ingress转发策略](#%E9%85%8D%E7%BD%AEingress%E8%BD%AC%E5%8F%91%E7%AD%96%E7%95%A5)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## k8s发布go应用

### 前言

搭建了一套K8s,尝试发布一个go应用

### 部署

#### 镜像打包

之前已经打包过一个go的镜像了，这次就直接跳过了  

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

### 配置service



### 配置ingress

配置ingress`https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.18.0/deploy/mandatory.yaml`
把 `extensions/v1beta1`修改成`apps/v1`

创建

```
$ kubectl apply -f ./mandatory.yaml
$ kubectl get pods -n ingress-nginx
$ kubectl get service -n ingress-nginx
```

测试  

```
curl http://10.98.51.155  //ip来自上一步
如果返回404则表示安装成功
```

### 配置ingress转发策略




### 参考
【K8S 安装 Ingress】https://www.jianshu.com/p/4370c00c040a  