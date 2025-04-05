<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [kubectl 如何管理多个 kubeconfig](#kubectl-%E5%A6%82%E4%BD%95%E7%AE%A1%E7%90%86%E5%A4%9A%E4%B8%AA-kubeconfig)
  - [前言](#%E5%89%8D%E8%A8%80)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## kubectl 如何管理多个 kubeconfig

### 前言

连接k8s 的时候，一般我们会使用 kubeconfig，正常情况下，我们会有至少两套环境，一套 测试环境，一套线上环境。当然有时候我们还有会一套预发布环境，如何快速切换环境呢？这里我们来探讨下。  

### 如何使用

先来看下 kubeconfig 中的文件结构

```yaml
apiVersion: v1
clusters:
  - cluster:
      # 集群的api 访问地址
      server: https://111.111.777.4:6443
      # Base64 编码的 CA 证书数据，用于验证 API Server 的身份。
      certificate-authority-data: c2RzZHNkc2RzZHNkc2RzZHMxMjEy
    # 集群的名字
    name: kubernetes-production
contexts:
  - context:
      # 集群的名字
      cluster: kubernetes-production
      # 指定该上下文使用的用户名，这里是 "kubernetes-admin-1"。
      user: "kubernetes-admin-1"
     # 上下文的名称，这里是 "production"。
    name: production
# 当前激活的上下文名称，这里是 "production"。
current-context: production
kind: Config
preferences: {}
users:
  - name: "kubernetes-admin-1"
    user:
      client-certificate-data: c2RzZHNkc2RzZHNkc2RzZHMxMjEy
      client-key-data: c2RzZHNkc2RzZHNkc2RzZHMxMjEy
```

创建两个 kubeconfig 文件，一个是测试环境，一个是线上环境 

```yaml
apiVersion: v1
clusters:
  - cluster:
      server: https://111.111.777.4:6443
      certificate-authority-data: c2RzZHNkc2RzZHNkc2RzZHMxMjEy
    name: kubernetes-test
contexts:
  - context:
      cluster: kubernetes-test
      user: "kubernetes-admin"
    name: test
current-context: test
kind: Config
preferences: {}
users:
  - name: "kubernetes-admin"
    user:
      client-certificate-data: c2RzZHNkc2RzZHNkc2RzZHMxMjEy
      client-key-data: c2RzZHNkc2RzZHNkc2RzZHMxMjEy
```

执行两个 kubeconfig 文件的合并  


合并

```
$ export KUBECONFIG=./kube/config-test:./kube/config-production
```

获取当前的上下文信息

```
$ kubectl config get-contexts

CURRENT   NAME         CLUSTER                 AUTHINFO             NAMESPACE
          production   kubernetes-production   kubernetes-admin-1   
*         test         kubernetes-test         kubernetes-admin 
```

获取当前的上下文信息

```
$ kubectl config current-context
test
```

切换上下文

```
$ kubectl config use-context production
Switched to context "production".
```



### 参考

【配置对多集群的访问】https://kubernetes.io/zh-cn/docs/tasks/access-application-cluster/configure-access-multiple-clusters/    
【如何管理和切换多个 Kubernetes kubeconfig 文件】https://blog.csdn.net/Mint6/article/details/142572046  