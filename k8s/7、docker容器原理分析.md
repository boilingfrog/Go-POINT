<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [docker 容器原理分析](#docker-%E5%AE%B9%E5%99%A8%E5%8E%9F%E7%90%86%E5%88%86%E6%9E%90)
  - [docker 的工作方式](#docker-%E7%9A%84%E5%B7%A5%E4%BD%9C%E6%96%B9%E5%BC%8F)
    - [Namespace](#namespace)
    - [Cgroups](#cgroups)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## docker 容器原理分析

### docker 的工作方式

当我们的程序运行起来的时候，在计算机中的表现就是一个个的进程，一个 docker 容器中可以跑各种程序，容器技术的核心功能，就是通过约束和修改进程的动态表现，从而为其创造出一个“边界”。  

对于 Docker 等大多数 Linux 容器来说，Cgroups 技术是用来进行资源限制，而 Namespace 技术用来做隔离作用。  

容器是一种特殊的进程：  

docker 容器在创建进程时，指定了这个进程所需要启用的一组 Namespace 参数。这样，容器就只能“看”到当前 Namespace 所限定的资源、文件、设备、状态，或者配置。而对于宿主机以及其他不相关的程序，它就完全看不到了。  

容器中启动的进程，还是在宿主机中运行的，只是 docker 容器会给这些进程，添加各种各样的 Namespace 参数，使这些进程和宿主机中的其它进程隔离开来，感知不到有其它进程的存在。     

下面来看下 docker 容器中，Namespace 和 Cgroups 的具体作用  

#### Namespace

Namespace 是 Linux 中隔离内核资源的方式，通过 Namespace 可以这些进程只能看到自己 Namespace 的相关资源，这样和其它 Namespace 的进程起到了隔离的作用。  

`Linux namespaces` 是对全局系统资源的一种封装隔离，使得处于不同 namespace 的进程拥有独立的全局系统资源，改变一个 namespace 中的系统资源只会影响当前 namespace 里的进程，对其他 namespace 中的进程没有影响。   

docker 容器的实现正是用到了 Namespace 的隔离，docker 容器通过启动的时候给进程添加  Namespace 参数，这样，容器就只能“看”到当前 Namespace 所限定的资源、文件、设备、状态，或者配置。而对于宿主机以及其他不相关的程序，它就完全看不到了。  







#### Cgroups


### 参考

【深入剖析 Kubernetes】https://time.geekbang.org/column/intro/100015201?code=UhApqgxa4VLIA591OKMTemuH1%2FWyLNNiHZ2CRYYdZzY%3D   
【Linux Namespace】https://www.cnblogs.com/sparkdev/p/9365405.html  
【浅谈 Linux Namespace】https://xigang.github.io/2018/10/14/namespace-md/   





