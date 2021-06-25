## 理解ConfigMap

### 什么是Secret

`Secret`对象类型用来保存敏感信息，例如密码、`OAuth`令牌和 SSH 密钥。 将这些信息放在`secret`中比放在`Pod`的定义或者 容器镜像中来说更加安全和灵活。   

`Secret`是一种包含少量敏感信息例如密码、令牌或密钥的对象。 这样的信息可能会被放在`Pod` 规约中或者镜像中。 用户可以创建`Secret`，同时系统也创建了一些`Secret`。

### Secret类型

要使用`Secret`，Pod 需要引用 Secret。 Pod 可以用三种方式之一来使用 Secret：

- 作为挂载到一个或多个容器上的 卷 中的文件。

- 作为容器的环境变量

- 由 kubelet 在为 Pod 拉取镜像时使用