<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [理解ConfigMap](#%E7%90%86%E8%A7%A3configmap)
  - [什么是ConfigMap](#%E4%BB%80%E4%B9%88%E6%98%AFconfigmap)
  - [ConfigMap的创建](#configmap%E7%9A%84%E5%88%9B%E5%BB%BA)
    - [使用key-value 字符串创建](#%E4%BD%BF%E7%94%A8key-value-%E5%AD%97%E7%AC%A6%E4%B8%B2%E5%88%9B%E5%BB%BA)
    - [从 env 文件创建](#%E4%BB%8E-env-%E6%96%87%E4%BB%B6%E5%88%9B%E5%BB%BA)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 理解ConfigMap

### 什么是ConfigMap

首先来弄明白为什么需要`ConfigMap`。对于应用开发来讲，特别是后端开发。我们需要连接数据库，mysql,redis..。这些链接，我们在测试环境，生成环境用的是多套。这就意味着，我们的代码中有一些配置需要经常改动。所以`ConfigMap`就是解决这些问题的。  

`ConfigMap`是一种API对象，用来将非机密性的数据保存到键值对中。使用时，`Pods`可以将其用作环境变量、命令行参数或者存储卷中的配置文件。  

`ConfigMap`将环境配置信息和容器镜像解耦，便于应用配置的修改。    

不过需要注意的是`ConfigMap`本身不提供加密功能。如果要存储的数据是机密的，使用`Secret`，或者使用其他第三方工具来保证你的数据的私密性，而不是用`ConfigMap`。  

`ConfigMap`在设计上不是用来保存大量数据的。在`ConfigMap`中保存的数据不可超过`1 MiB`。如果你需要保存超出此尺寸限制的数据，你可能希望考虑挂载存储卷或者使用独立的数据库或者文件服务。  

### ConfigMap的创建

可以使用`kubectl create configmap`从文件、目录或者`key-value`字符串创建等创建`ConfigMap`。也可以通过`kubectl create -f file`创建。  

#### 使用key-value 字符串创建

创建name和age

```
$ kubectl create configmap config-test --from-literal=name=xiaoming  --from-literal=age=22 
```

查看结果

```
$ kubectl get configmap config-test -o go-template='{{.data}}'
map[age:22 name:xiaoming]#  
```

#### 从 env 文件创建

创建env文件  

```
$ echo -e "name=xiaobai\nage=25" | tee config.env
name=xiaobai
age=25
```

写入`config-test-1`内容  

```
$ kkubectl create configmap config-test-1 --from-env-file=config.env
configmap/config-test-1 created
```

查看内容

```
$ kubectl get configmap config-test-1 -o go-template='{{.data}}'
map[age:25 name:xiaobai]#     
```




### 参考

【ConfigMap】https://kubernetes.io/zh/docs/concepts/configuration/configmap/  
【ConfigMap】https://feisky.gitbooks.io/kubernetes/content/concepts/configmap.html  