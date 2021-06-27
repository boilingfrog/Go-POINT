<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [理解ConfigMap](#%E7%90%86%E8%A7%A3configmap)
  - [什么是Secret](#%E4%BB%80%E4%B9%88%E6%98%AFsecret)
  - [Secret类型](#secret%E7%B1%BB%E5%9E%8B)
    - [Opaque Secret](#opaque-secret)
  - [Opaque Secret的使用](#opaque-secret%E7%9A%84%E4%BD%BF%E7%94%A8)
    - [将Secret挂载到Volume中](#%E5%B0%86secret%E6%8C%82%E8%BD%BD%E5%88%B0volume%E4%B8%AD)
    - [将Secret导出到环境变量中](#%E5%B0%86secret%E5%AF%BC%E5%87%BA%E5%88%B0%E7%8E%AF%E5%A2%83%E5%8F%98%E9%87%8F%E4%B8%AD)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 理解ConfigMap

### 什么是Secret

`Secret`对象类型用来保存敏感信息，例如密码、`OAuth`令牌和 SSH 密钥。 将这些信息放在`secret`中比放在`Pod`的定义或者 容器镜像中来说更加安全和灵活。   

`Secret`是一种包含少量敏感信息例如密码、令牌或密钥的对象。 这样的信息可能会被放在`Pod` 规约中或者镜像中。 用户可以创建`Secret`，同时系统也创建了一些`Secret`。

### Secret类型

要使用`Secret`，Pod 需要引用 Secret。 Pod 可以用三种方式之一来使用 Secret：

- 作为挂载到一个或多个容器上的 卷 中的文件。

- 作为容器的环境变量

- 由 kubelet 在为 Pod 拉取镜像时使用

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

#### 将Secret导出到环境变量中




### 参考

【Secret】https://feisky.gitbooks.io/kubernetes/content/concepts/secret.html    
【k8s官方对Secret的描述】https://kubernetes.io/zh/docs/concepts/configuration/secret/  