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
          image: liz2019/test-static
          resources:
            limits:
              memory: "128Mi"
              cpu: "500m"
          ports:
            - containerPort: 3000
```




### 参考
【K8S 安装 Ingress】https://www.jianshu.com/p/4370c00c040a  