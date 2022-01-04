<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [理解ConfigMap](#%E7%90%86%E8%A7%A3configmap)
  - [什么是ConfigMap](#%E4%BB%80%E4%B9%88%E6%98%AFconfigmap)
  - [ConfigMap的创建](#configmap%E7%9A%84%E5%88%9B%E5%BB%BA)
    - [使用key-value 字符串创建](#%E4%BD%BF%E7%94%A8key-value-%E5%AD%97%E7%AC%A6%E4%B8%B2%E5%88%9B%E5%BB%BA)
    - [从env文件创建](#%E4%BB%8Eenv%E6%96%87%E4%BB%B6%E5%88%9B%E5%BB%BA)
    - [从目录创建](#%E4%BB%8E%E7%9B%AE%E5%BD%95%E5%88%9B%E5%BB%BA)
    - [通过Yaml/Json创建](#%E9%80%9A%E8%BF%87yamljson%E5%88%9B%E5%BB%BA)
  - [ConfigMap使用](#configmap%E4%BD%BF%E7%94%A8)
    - [用作环境变量](#%E7%94%A8%E4%BD%9C%E7%8E%AF%E5%A2%83%E5%8F%98%E9%87%8F)
    - [用作命令参数](#%E7%94%A8%E4%BD%9C%E5%91%BD%E4%BB%A4%E5%8F%82%E6%95%B0)
    - [使用volume将ConfigMap作为文件或目录直接挂载](#%E4%BD%BF%E7%94%A8volume%E5%B0%86configmap%E4%BD%9C%E4%B8%BA%E6%96%87%E4%BB%B6%E6%88%96%E7%9B%AE%E5%BD%95%E7%9B%B4%E6%8E%A5%E6%8C%82%E8%BD%BD)
    - [使用subpath将ConfigMap作为单独的文件挂载到目录](#%E4%BD%BF%E7%94%A8subpath%E5%B0%86configmap%E4%BD%9C%E4%B8%BA%E5%8D%95%E7%8B%AC%E7%9A%84%E6%96%87%E4%BB%B6%E6%8C%82%E8%BD%BD%E5%88%B0%E7%9B%AE%E5%BD%95)
  - [被挂载的ConfigMap内容会被自动更新](#%E8%A2%AB%E6%8C%82%E8%BD%BD%E7%9A%84configmap%E5%86%85%E5%AE%B9%E4%BC%9A%E8%A2%AB%E8%87%AA%E5%8A%A8%E6%9B%B4%E6%96%B0)
  - [不可变更的ConfigMap](#%E4%B8%8D%E5%8F%AF%E5%8F%98%E6%9B%B4%E7%9A%84configmap)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 理解ConfigMap

### 什么是ConfigMap

首先来弄明白为什么需要`ConfigMap`。对于应用开发来讲，特别是后端开发。我们需要连接数据库，mysql,redis..。这些连接，我们在测试环境，生成环境用的是多套。这就意味着，我们的代码中有一些配置需要经常改动。所以`ConfigMap`就是解决这些问题的。  

`ConfigMap`是一种API对象，用来将非机密性的数据保存到键值对中。使用时，`Pods`可以将其用作环境变量、命令行参数或者存储卷中的配置文件。  

`ConfigMap`将环境配置信息和容器镜像解耦，便于应用配置的修改。    

不过需要注意的是`ConfigMap`本身不提供加密功能。如果要存储的数据是机密的，使用`Secret`，或者使用其他第三方工具来保证你的数据的私密性，而不是用`ConfigMap`。  

`ConfigMap`在设计上不是用来保存大量数据的。在`ConfigMap`中保存的数据不可超过`1 MiB`。如果你需要保存超出此尺寸限制的数据，你可能希望考虑挂载存储卷或者使用独立的数据库或者文件服务。  

### ConfigMap的创建

可以使用`kubectl create configmap`从文件、目录或者`key-value`字符串创建等创建`ConfigMap`。也可以通过`kubectl create -f file`创建。  

#### 使用key-value 字符串创建

创建name和age

```
$ kubectl create configmap config-test-1 --from-literal=name=xiaoming  --from-literal=age=22 
```

查看结果

```
$ kubectl get configmap config-test-1 -o go-template='{{.data}}'
map[age:22 name:xiaoming]#  
```

#### 从env文件创建

创建env文件  

```
$ echo -e "name=xiaobai\nage=25" | tee config.env
name=xiaobai
age=25
```

写入`config-test-1`内容  

```
$ kkubectl create configmap config-test-2 --from-env-file=config.env
configmap/config-test-2 created
```

查看内容

```
$ kubectl get configmap config-test-2 -o go-template='{{.data}}'
map[age:25 name:xiaobai]#     
```

#### 从目录创建

创建对应的目录

```
$ mkdir config
$ echo 18>config/age  
$ echo xiaohua>config/name
```

写入`config-test-3`内容  

```
$ kubectl create configmap config-test-3 --from-file=config/
configmap/config-test-3 created
```

查看写入的内容

```
$ kubectl get configmap config-test-3 -o go-template='{{.data}}'
map[age:18
 name:xiaohua
]#       
```

#### 通过Yaml/Json创建

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: config-test-4
  namespace: default
data:
  name: xiaolong
  age: 16year # 需要是string
```

创建

```
$ kubectl create  -f config-test-4.yaml
configmap/config-test-4 created
```

### ConfigMap使用

`ConfigMap`可以通过三种方式在Pod中使用，三种分别方式为：设置环境变量、设置容器命令行参数以及在`Volume`中直接挂载文件或目录。  

需要注意的点  

> ConfigMap 必须在 Pod 引用它之前创建

> 使用 envFrom 时，将会自动忽略无效的键

> Pod 只能使用同一个命名空间内的 ConfigMap 

首先创建 ConfigMap：

```
$ kubectl create configmap special-config --from-literal=name=long --from-literal=realname=xiaolong
$ kubectl create configmap env-config --from-literal=log_level=INFO
```

#### 用作环境变量

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
spec:
  containers:
    - name: test-container
      image: busybox
      command: ["/bin/sh", "-c", "env"]
      env:
        - name: SPECIAL_NAME_KEY
          valueFrom:
            configMapKeyRef:
              name: special-config
              key: name
        - name: SPECIAL_REALNAME_KEY
          valueFrom:
            configMapKeyRef:
              name: special-config
              key: realname
      envFrom:
        - configMapRef:
            name: env-config
  restartPolicy: Never
```

运行之后查看日志

```
$ kubectl logs -f test-pod
HOSTNAME=test-pod
SPECIAL_NAME_KEY=long
log_level=INFO
SPECIAL_REALNAME_KEY=xiaolong
```

发现上面的值已成功写入了  

#### 用作命令参数

将`ConfigMap`用作命令行参数时，需要先把`ConfigMap`的数据保存在环境变量中，然后通过`$(VAR_NAME)`的方式引用环境变量。  

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: dapi-test-pod
spec:
  containers:
    - name: test-container
      image: gcr.io/google_containers/busybox
      command: ["/bin/sh", "-c", "echo $(SPECIAL_NAME_KEY) $(SPECIAL_REALNAME_KEY)" ]
      env:
        - name: SPECIAL_NAME_KEY
          valueFrom:
            configMapKeyRef:
              name: special-config
              key: name
        - name: SPECIAL_REALNAME_KEY
          valueFrom:
            configMapKeyRef:
              name: special-config
              key: realname
  restartPolicy: Never
```

运行之后打印日志

```
$ kubectl logs -f dapi-test-pod
long xiaolong
```

输出我们之前的写入的配置信息  

#### 使用volume将ConfigMap作为文件或目录直接挂载

将创建的`ConfigMap`直接挂载至 Pod 的`/etc/config`目录下，其中每一个`key-value`键值对都会生成一个文件，`key`为文件名，`value`为内容  

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: vol-test-pod
spec:
  containers:
    - name: test-container
      image: busybox
      command: ["/bin/sh", "-c", "cat /etc/config/name"]
      volumeMounts:
      - name: config-volume
        mountPath: /etc/config
  volumes:
    - name: config-volume
      configMap:
        name: special-config
  restartPolicy: Never
```

启动之后打印输出

```
$ kubectl logs -f vol-test-pod
long#       
```

#### 使用subpath将ConfigMap作为单独的文件挂载到目录

在一般情况下`configmap`挂载文件时，会先覆盖掉挂载目录，然后再将`congfigmap`中的内容作为文件挂载进行。如果想不对原来的文件夹下的文件造成覆盖，只是将`configmap`中的每个 key，按照文件的方式挂载到目录下，可以使用`subpath`参数。

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: dapi-test-pod
spec:
  containers:
    - name: test-container
      image: nginx
      command: ["/bin/sh","-c","sleep 36000"]
      volumeMounts:
      - name: config-volume
        mountPath: /etc/nginx/name
        subPath: name
  volumes:
    - name: config-volume
      configMap:
        name: special-config
        items:
        - key: name
          path: name
  restartPolicy: Never
```

进入到pod中查看  

```
root@dapi-test-pod:/etc/nginx# ls
conf.d		mime.types  name	scgi_params
fastcgi_params	modules     nginx.conf	uwsgi_params
```

### 被挂载的ConfigMap内容会被自动更新

当卷中使用的`ConfigMap`被更新时，所投射的键最终也会被更新。`kubelet`组件会在每次周期性同步时检查所挂载的`ConfigMap`是否为最新。 不过,`kubelet`使用的是其本地的高速缓存来获得`ConfigMap`的当前值。  

`ConfigMap`既可以通过`watch`操作实现内容传播（默认形式），也可实现基于`TTL`的缓存，还可以直接经过所有请求重定向到`API`服务器。 因此，从`ConfigMap`被更新的那一刻算起，到新的主键被投射到`Pod`中去，这一 时间跨度可能与 kubelet 的同步周期加上高速缓存的传播延迟相等。 这里的传播延迟取决于所选的高速缓存类型 （分别对应`watch`操作的传播延迟、高速缓存的`TTL`时长或者 0）。  

以环境变量方式使用的`ConfigMap`数据不会被自动更新。 更新这些数据需要重新启动`Pod`。  

### 不可变更的ConfigMap

`Kubernetes`特性 不可变更的`Secret`和`ConfigMap`提供了一种将各个`Secret`和`ConfigMap`设置为不可变更的选项。对于大量使用 ConfigMap 的 集群（至少有数万个各不相同的`ConfigMap`给`Pod` 挂载）而言，禁止更改`ConfigMap`的数据有以下好处：  

- 保护应用，使之免受意外（不想要的）更新所带来的负面影响。

- 通过大幅降低对 `kube-apiserver` 的压力提升集群性能，这是因为系统会关闭 对已标记为不可变更的`ConfigMap`的监视操作。

可以通过将`immutable`字段设置为`true`创建不可变更的`ConfigMap`。 例如：

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  ...
data:
  ...
immutable: true
```

一旦某`ConfigMap`被标记为不可变更，则 无法 逆转这一变化，，也无法更改`data`或`binaryData`字段的内容。你只能删除并重建`ConfigMap`。 因为现有的`Pod`会维护一个对已删除的`ConfigMap`的挂载点，建议重新创建这些`Pods`。  

### 参考

【ConfigMap】https://kubernetes.io/zh/docs/concepts/configuration/configmap/  
【ConfigMap】https://feisky.gitbooks.io/kubernetes/content/concepts/configmap.html  