<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [理解Secret](#%E7%90%86%E8%A7%A3secret)
  - [什么是Secret](#%E4%BB%80%E4%B9%88%E6%98%AFsecret)
  - [Secret的类型](#secret%E7%9A%84%E7%B1%BB%E5%9E%8B)
    - [Opaque Secret](#opaque-secret)
  - [Opaque Secret的使用](#opaque-secret%E7%9A%84%E4%BD%BF%E7%94%A8)
    - [将Secret挂载到Volume中](#%E5%B0%86secret%E6%8C%82%E8%BD%BD%E5%88%B0volume%E4%B8%AD)
      - [挂载的Secret会被自动更新](#%E6%8C%82%E8%BD%BD%E7%9A%84secret%E4%BC%9A%E8%A2%AB%E8%87%AA%E5%8A%A8%E6%9B%B4%E6%96%B0)
    - [将Secret导出到环境变量中](#%E5%B0%86secret%E5%AF%BC%E5%87%BA%E5%88%B0%E7%8E%AF%E5%A2%83%E5%8F%98%E9%87%8F%E4%B8%AD)
      - [Secret更新之后对应的环境变量不会被更新](#secret%E6%9B%B4%E6%96%B0%E4%B9%8B%E5%90%8E%E5%AF%B9%E5%BA%94%E7%9A%84%E7%8E%AF%E5%A2%83%E5%8F%98%E9%87%8F%E4%B8%8D%E4%BC%9A%E8%A2%AB%E6%9B%B4%E6%96%B0)
  - [不可更改的Secret](#%E4%B8%8D%E5%8F%AF%E6%9B%B4%E6%94%B9%E7%9A%84secret)
  - [Secret与Pod生命周期的关系](#secret%E4%B8%8Epod%E7%94%9F%E5%91%BD%E5%91%A8%E6%9C%9F%E7%9A%84%E5%85%B3%E7%B3%BB)
  - [Secret与ConfigMap对比](#secret%E4%B8%8Econfigmap%E5%AF%B9%E6%AF%94)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 理解Secret

### 什么是Secret

`Secret`对象类型用来保存敏感信息，例如密码、`OAuth`令牌和 SSH 密钥。 将这些信息放在`secret`中比放在`Pod`的定义或者 容器镜像中来说更加安全和灵活。   

`Secret`是一种包含少量敏感信息例如密码、令牌或密钥的对象。 这样的信息可能会被放在`Pod` 规约中或者镜像中。 用户可以创建`Secret`，同时系统也创建了一些`Secret`。

要使用`Secret`，Pod 需要引用`Secret`。 `Pod`可以用三种方式之一来使用 Secret：

- 作为挂载到一个或多个容器上的卷中的文件。

- 作为容器的环境变量

- 由kubelet在为Pod拉取镜像时使用

### Secret的类型

`Kubernetes`提供若干种内置的类型，用于一些常见的使用场景。 针对这些类型，`Kubernetes`所执行的合法性检查操作以及对其所实施的限制各不相同。   

|     内置类型                           | 用法                                         |
| ------------------------------        |  -------------------------------            |
| Opaque                                |    用户定义的任意数据                          |
| kubernetes.io/service-account-token   |    服务账号令牌                               |
| kubernetes.io/dockercfg               |    ~/.dockercfg 文件的序列化形式               |
| kubernetes.io/dockerconfigjson        |    ~/.docker/config.json 文件的序列化形式      |
| kubernetes.io/basic-auth              |    用于基本身份认证的凭据                       |
| kubernetes.io/ssh-auth                |    用于 SSH 身份认证的凭据                      |
| kubernetes.io/tls	                    |    用于 TLS 客户端或者服务器端的数据              |
| bootstrap.kubernetes.io/token	        |    启动引导令牌数据                             |

#### Opaque Secret

`Opaque`类型的数据是一个`map`类型，要求`value`是`base64`编码格式：  

创建我们需要的两个`base64`账号信息  

```
$ echo -n "root" | base64
  cm9vdA==
$ echo -n "123456" | base64
  MTIzNDU2
```

secrets.yml

```
apiVersion: v1
kind: Secret
metadata:
  name: secret-test
type: Opaque
data:
  password: MTIzNDU2
  username: cm9vdA==
```

创建 secret：`kubectl create -f secrets.yml`  

```
$ kubectl get secret
  NAME                         TYPE                                  DATA   AGE
  secret-test                  Opaque                                2      8s
```

可以看到我们刚刚创建的`secret`  

### Opaque Secret的使用

创建好`secret`之后，有两种方式来使用它：  

- 以 Volume 方式

- 以环境变量方式

#### 将Secret挂载到Volume中

````yaml
apiVersion: v1
kind: Pod
metadata:
  name: mypod
spec:
  containers:
  - name: mypod
    image: redis
    volumeMounts:
    - name: foo
      mountPath: "/etc/foo"
      readOnly: true
  volumes:
  - name: foo
    secret:
      secretName: secret-test
      items:
      - key: username
        path: my-group/my-username
````

成功运行之后，进入到pod中  

```
root@mypod:/etc/foo/my-group# ls
my-username

root@mypod:/etc/foo/my-group# cat my-username 
root
```

可以看到我们之前定义的secret的username已经被当前的pod引用，并成功挂载  

##### 挂载的Secret会被自动更新 

当已经存储于卷中被使用的`Secret`被更新时，被映射的键也将终将被更新。 组件`kubelet`在周期性同步时检查被挂载的`Secret`是不是最新的。 但是，它会使用其本地缓存的数值作为`Secret`的当前值。  

不过，使用`Secret`作为子路径卷挂载的容器 不会收到`Secret`更新。  

#### 将Secret导出到环境变量中

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: secret-env-pod
spec:
  containers:
  - name: mycontainer
    image: redis
    env:
      - name: SECRET_USERNAME
        valueFrom:
          secretKeyRef:
            name: secret-test
            key: username
      - name: SECRET_PASSWORD
        valueFrom:
          secretKeyRef:
            name: secret-test
            key: password
  restartPolicy: Never
```

启动pod，之后进入pod

```
root@secret-env-pod:/data# echo $SECRET_USERNAME
root
root@secret-env-pod:/data# echo $SECRET_PASSWORD
123456
```

可以看到secret的值已经成功写入到环境变量中了    

##### Secret更新之后对应的环境变量不会被更新 

如果某个容器已经在通过环境变量使用某`Secret`，对该`Secret`的更新不会被 容器马上看见，除非容器被重启。有一些第三方的解决方案能够在`Secret`发生变化时触发容器重启。  

### 不可更改的Secret

`Kubernetes`对于`Secret`和`ConfigMap`提供了一种不可变更的配置选择，对于大量使用`Secret`的集群（至少有成千上万各不相同的`Secret`供`Pod`挂载）， 禁止变更它们的数据有下列好处：  

- 防止意外（或非预期的）更新导致应用程序中断

- 通过将 Secret 标记为不可变来关闭 kube-apiserver 对其的监视，从而显著降低 kube-apiserver 的负载，提升集群性能。

将`Secret`的`immutable`字段设置为`true`创建不可更改的`Secret`  

```
apiVersion: v1
kind: Secret
metadata:
  ...
data:
  ...
immutable: true
```

1、当一个`secret`设置成不可更改，如果想要更改`secret`中的内容，就需要删除并且重新创建这个`secret`  

2、对于引用老的`secret`的pod，需要删除并且重新创建  



### Secret与Pod生命周期的关系 

通过API创建Pod时，不会检查引用的Secret是否存在。一旦Pod被调度，kubelet就会尝试获取该Secret的值。如果获取不到该Secret，或者暂时无法与API服务器建立连接，kubelet将会定期重试。kubelet将会报告关于 Pod 的事件，并解释它无法启动的原因。 一旦获取到Secret，kubelet将创建并挂载一个包含它的卷。在Pod的所有卷被挂载之前，Pod中的容器不会启动。  

### Secret与ConfigMap对比

**相同点：**

- key/value 的形式

- 属于某个特定的 namespace

- 可以导出到环境变量

- 可以通过目录 / 文件形式挂载 (支持挂载所有 key 和部分 key)

**不同点：**

- Secret 可以被 ServerAccount 关联 (使用)

- Secret 可以存储 register 的鉴权信息，用在 ImagePullSecret 参数中，用于拉取私有仓库的镜像

- Secret 支持 Base64 加密

- Secret 文件存储在 tmpfs 文件系统中，Pod 删除后 Secret 文件也会对应的删除。


### 参考

【Secret】https://feisky.gitbooks.io/kubernetes/content/concepts/secret.html    
【k8s官方对Secret的描述】https://kubernetes.io/zh/docs/concepts/configuration/secret/  