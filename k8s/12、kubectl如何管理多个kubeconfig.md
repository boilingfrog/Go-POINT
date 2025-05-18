<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [kubectl 如何管理多个 kubeconfig](#kubectl-%E5%A6%82%E4%BD%95%E7%AE%A1%E7%90%86%E5%A4%9A%E4%B8%AA-kubeconfig)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [如何使用](#%E5%A6%82%E4%BD%95%E4%BD%BF%E7%94%A8)
  - [编写优化脚本](#%E7%BC%96%E5%86%99%E4%BC%98%E5%8C%96%E8%84%9A%E6%9C%AC)
  - [第三方工具](#%E7%AC%AC%E4%B8%89%E6%96%B9%E5%B7%A5%E5%85%B7)
  - [参考](#%E5%8F%82%E8%80%83)

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

切换上下文

```
$ kubectl config current-context
test
```

切换上下文

```
$ kubectl config use-context production
Switched to context "production".
```

可以看到， clusters 的 name(标识集群的名字) ，和 contexts 的 name(标识上下文的名字)，。在我们切换配置文件的时候，帮助最大.clusters 的 name,帮我们区分是那个集群，contexts 的 name，用来帮助在切换集群的时候，进行上下文的区分。   

弄明白了，kubeconfig 的文件构成和如何设置之后，我们尝试使用脚本来帮助我们处理。  

script.sh

```
#!/bin/bash

# 设置 KUBECONFIG 环境变量
export KUBECONFIG=./kube/config-test:./kube/config-production

# 列出所有上下文
kubectl config get-contexts

# 提示用户选择上下文
echo "请输入要使用的 k8s 集群的 NAME："
read context_name

# 切换到选择的上下文
kubectl config use-context $context_name

# 显示当前上下文
kubectl config current-context
```

添加执行的权限

```
 chmod +x ./kube/script.sh
```

执行  

```
$ ./kube/script.sh         
CURRENT   NAME         CLUSTER                 AUTHINFO             NAMESPACE
          production   kubernetes-production   kubernetes-admin-1   
*         test         kubernetes-test         kubernetes-admin     
请输入要使用的 k8s 集群的 NAME：
test
Switched to context "test".
test
```

### 编写优化脚本

上面展示了 kubeconfig 中文件的构成，已经如何进行多个 config 的切换，这里添加一个更加方便的脚本。  

指定文件的名字，文件的名字就是我们集群的上下文。    

创建两个集群的 yaml, config-test.yaml 和  config-production.yaml 

```
#!/bin/bash

# 设置 KUBECONFIG 环境变量
export KUBECONFIG=$(ls ./kube/config-*.yaml | tr '\n' ':' | sed 's/:$//')


# 获取所有上下文
contexts=$(kubectl config get-contexts -o name)

# 重命名上下文
for context in $contexts; do
  for file in ./kube/config-*.yaml; do
    if grep -q $context $file; then
      suffix=$(basename $file | sed 's/config-//;s/.yaml//')
      kubectl config rename-context $context $suffix
    fi
  done
done

# 列出所有上下文
kubectl config get-contexts

# 提示用户选择上下文
echo "请输入要使用的 k8s 集群的 NAME："
read context_name

# 切换到选择的上下文
kubectl config use-context $context_name

# 显示当前上下文
kubectl config current-context
```

### 第三方工具

如果有了上下文信息，这里也可以使用第三方工具进行快速的切换。 这里推荐使用 kubectx。具体安装地址参考[https://github.com/ahmetb/kubectx]

```
$ kubectx
production
test

$ kubectx test
Switched to context "test".
```

这里再放一个切换合并上下文的脚本  

```
#!/bin/bash

# 设置 KUBECONFIG 环境变量
export KUBECONFIG=$(ls ./kube/config-*.yaml | tr '\n' ':' | sed 's/:$//')


# 获取所有上下文
contexts=$(kubectl config get-contexts -o name)

# 重命名上下文
for context in $contexts; do
  for file in ./kube/config-*.yaml; do
    if grep -q $context $file; then
      suffix=$(basename $file | sed 's/config-//;s/.yaml//')
      kubectl config rename-context $context $suffix
    fi
  done
done

# 列出所有上下文
kubectl config get-contexts

# 合并 kubeconfig 文件
kubectl config view --merge --flatten > ~./config
```

我们可以将合并之后的 config, 放到 ~/.kube/config 文件中。借助于 kubectx 就能很好的进行 k8s 上下文的切换了。   

### 参考

【配置对多集群的访问】https://kubernetes.io/zh-cn/docs/tasks/access-application-cluster/configure-access-multiple-clusters/    
【如何管理和切换多个 Kubernetes kubeconfig 文件】https://blog.csdn.net/Mint6/article/details/142572046  
【kubectx】https://github.com/ahmetb/kubectx  
【kubectl 如何管理多个 kubeconfig】https://github.com/boilingfrog/Go-POINT/blob/master/k8s/12%E3%80%81kubectl%E5%A6%82%E4%BD%95%E7%AE%A1%E7%90%86%E5%A4%9A%E4%B8%AAkubeconfig.md  