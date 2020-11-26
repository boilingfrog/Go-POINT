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
$ kubectl expose deployment go-app --type=NodePort --name=go-app-svc --target-port=3000
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