<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [为什么 K8s 会抛弃 docker](#%E4%B8%BA%E4%BB%80%E4%B9%88-k8s-%E4%BC%9A%E6%8A%9B%E5%BC%83-docker)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [CRI](#cri)
  - [containerd](#containerd)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 为什么 K8s 会抛弃 docker   

### 前言

在这之前先来了解下，k8s 是如何和 docker 进行交互的。   

### CRI  

kubelet 调用下层容器运行时的执行过程，并不会直接调用 Docker 的 API，而是通过 CRI（Container Runtime Interface，容器运行时接口）的 gRPC 接口来间接执行的。    

<img src="/img/k8s/k8s-cri.png"  alt="k8s" />    

为什么要引入 CRI？   

把 kubelet 对容器的操作，统一地抽象成一个接口，这样 kubelet 只需要和这个接口打交道，而不用关心底层容器，底层容器它们就只需要自己提供一个该接口的实现，然后对 kubelet 暴露出 gRPC 服务即可，这样底层容器就能很容器的进行切换了，而不是仅限于 Docker 这种容器了。  

同时，引入 CRI 接口，这样就不会受限于 docker 了，可以随时切换到其它的运行时。意味着容器运行时和镜像的实现与 Docker 项目完全剥离，让其他玩家不依赖 Docker 实现自己的运行时成为可能。   

### containerd

面对挑战，Docker 采取了“断臂求生”的策略，推动自身重构，将原有单一架构的 `Docker Engine` 拆分成多个模块，其中 `Docker daemon` 部分捐赠给 CNCF，containerd 形成。   

作为 CNCF 的托管项目，containerd 必须符合 CRI 标准。   

但是 docker 本身还是没有实现 CRI 标准。所以 k8s 中引入了一个叫作 dockershim 的组件。它会把 CRI 请求里的内容拿出来，然后组装成 `Docker API` 请求发给 `Docker Daemon`。    

```
kubelet --> dockershim （在 kubelet 进程中） --> dockerd --> containerd
```

containerd 作为 k8s 容器运行时   

```
kubelet --> cri plugin（在 containerd 进程中） --> containerd
```

显然都是通过 containerd 来管理容器的，所以这两种的调用效果最终是一样的。第二种，去掉了 dockershim 调用链更短了，性能更好了，同时因为不用维护 dockershim 了，维护性难度也大大减少了。    

所以 Docker 被抛弃的原因就很显而易见了。不过弃用 Docker 对 k8s 和 Docker 的影响不大，因为它们都已经将底层改为开源 containerd，原有的 Docker 镜像和容器仍然可以正常运行。唯一的变化是 K8s 绕过了 Docker，直接调用 Docker 内部的 containerd。    

<img src="/img/k8s/k8s-containerd.png"  alt="k8s" />

### 参考

【深入剖析 Kubernetes】https://time.geekbang.org/column/intro/100015201?code=UhApqgxa4VLIA591OKMTemuH1%2FWyLNNiHZ2CRYYdZzY%3D  
【K8s 为什么要弃用 Docker？】https://mp.weixin.qq.com/s/qEKyEseD370xWI-2yIyUzg     
【Docker与k8s的恩怨情仇】https://www.cnblogs.com/powertoolsteam/p/14980851.html     
【k8s为什么会抛弃docker】https://boilingfrog.github.io/2023/01/07/k8s%E4%B8%BA%E4%BB%80%E4%B9%88%E4%BC%9A%E6%8A%9B%E5%BC%83docker/   



