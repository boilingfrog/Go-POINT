## 理解ConfigMap

### 什么是ConfigMap

首先来弄明白为什么需要`ConfigMap`。对于应用开发来讲，特别是后端开发。我们需要连接数据库，mysql,redis..。这些链接，我们在测试环境，生成环境用的是多套。这就意味着，我们的代码中有一些配置需要经常改动。所以`ConfigMap`就是解决这些问题的。  

`ConfigMap`是一种API对象，用来将非机密性的数据保存到键值对中。使用时，`Pods`可以将其用作环境变量、命令行参数或者存储卷中的配置文件。  

`ConfigMap`将环境配置信息和容器镜像解耦，便于应用配置的修改。    

不过需要注意的是`ConfigMap`本身不提供加密功能。如果要存储的数据是机密的，使用`Secret`，或者使用其他第三方工具来保证你的数据的私密性，而不是用`ConfigMap`。  

`ConfigMap`在设计上不是用来保存大量数据的。在`ConfigMap`中保存的数据不可超过`1 MiB`。如果你需要保存超出此尺寸限制的数据，你可能希望考虑挂载存储卷或者使用独立的数据库或者文件服务。  

### 参考

【ConfigMap】https://kubernetes.io/zh/docs/concepts/configuration/configmap/  