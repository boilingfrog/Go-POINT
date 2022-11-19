<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [docker 容器原理分析](#docker-%E5%AE%B9%E5%99%A8%E5%8E%9F%E7%90%86%E5%88%86%E6%9E%90)
  - [docker 的工作方式](#docker-%E7%9A%84%E5%B7%A5%E4%BD%9C%E6%96%B9%E5%BC%8F)
    - [Namespace](#namespace)
      - [容器对比虚拟机](#%E5%AE%B9%E5%99%A8%E5%AF%B9%E6%AF%94%E8%99%9A%E6%8B%9F%E6%9C%BA)
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

##### 容器对比虚拟机

虚拟机

使用虚拟机，就需要使用 Hypervisor 来负责创建虚拟机，这个虚拟机是真实存在的，并且里面需要运行 `Guest OS` 才能执行用户的应用进程，这就不可避免会带来额外的资源消耗和占用。

虚拟机本身的运行就会占用一定的资源，同时虚拟机对宿主机文件的调用，就不可避免的需要经过虚拟化软件的连接和处理，这本身就是一层性能消耗，尤其对计算资源、网络和磁盘 I/O 的损耗非常大。

容器

容器化后的应用，还是宿主机上的一个普通的进程，不会存在虚拟化带来的性能损耗，同时容器使用的是 Namespace 进行隔离的，所以不需要单独的 `Guest OS`，这就使得容器额外的资源占用几乎可以忽略不计。

容器的缺点

基于 Linux Namespace 的隔离机制相比于虚拟化技术也有很多不足之处，其中最主要的问题就是：隔离得不彻底。

1、容器只是运行在宿主机中的一中特殊的进程，容器时间使用的还是同一个宿主机的操作系统内核；

尽管你可以在容器里通过 Mount Namespace 单独挂载其他不同版本的操作系统文件，比如 CentOS 或者 Ubuntu，但这并不能改变共享宿主机内核的事实。这意味着，如果你要在 Windows 宿主机上运行 Linux 容器，或者在低版本的 Linux 宿主机上运行高版本的 Linux 容器，都是行不通的。

2、在 Linux 内核中，有很多资源和对象是不能被 Namespace 化的，最典型的例子就是：时间；

如果容器中的程序使用了 `settimeofday(2)` 系统调用修改了时间，整个宿主机的时间都会被随之修改，所以容器中我们应该尽量避免这种操作。

3、容器共享宿主机内核，会给应用暴露出更大的攻击面。

在是生产环境中，不会把物理机中 Linux 的容器直接暴露在公网上。

#### Cgroups  

docker 容器中的进程使用 Namespace 来进行隔离，使得这些在容器中运行的进程像是运行在一个独立的环境中一样。但是，被隔离的进程还是运行在宿主机中的，如果这些进程没有对资源进行限制，这些进程可能会占用很多的系统资源，影响到其他的进程。Docker 使用 `Linux cgroups` 来限制容器中的进程允许使用的系统资源。  

`Linux Cgroups` 的全称是 `Linux Control Group`。它最主要的作用，就是限制一个进程组能够使用的资源上限，包括 CPU、内存、磁盘、网络带宽等等。   

在 Linux 中，Cgroups 给用户暴露出来的操作接口是文件系统，即它以文件和目录的方式组织在操作系统的 `/sys/fs/cgroup` 路径下。    

`centos 7.2` 下面的文件  

```
# mount -t cgroup
cgroup on /sys/fs/cgroup/systemd type cgroup (rw,nosuid,nodev,noexec,relatime,xattr,release_agent=/usr/lib/systemd/systemd-cgroups-agent,name=systemd)
cgroup on /sys/fs/cgroup/cpu,cpuacct type cgroup (rw,nosuid,nodev,noexec,relatime,cpuacct,cpu)
cgroup on /sys/fs/cgroup/net_cls type cgroup (rw,nosuid,nodev,noexec,relatime,net_cls)
cgroup on /sys/fs/cgroup/freezer type cgroup (rw,nosuid,nodev,noexec,relatime,freezer)
cgroup on /sys/fs/cgroup/memory type cgroup (rw,nosuid,nodev,noexec,relatime,memory)
...
```

`Linux Cgroups` 的设计还是比较易用的，简单粗暴地理解呢，它就是一个子系统目录加上一组资源限制文件的组合。而对于 Docker 等 Linux 容器项目来说，它们只需要在每个子系统下面，为每个容器创建一个控制组（即创建一个新目录），然后在启动容器进程之后，把这个进程的 PID 填写到对应控制组的 tasks 文件中就可以了。  

总结下来就是，一个正在运行的 Docker 容器，其实就是一个启用了多个 `Linux Namespace` 的应用进程，而这个进程能够使用的资源量，则受 Cgroups 配置的限制。  

### 容器中的文件







### 参考

【深入剖析 Kubernetes】https://time.geekbang.org/column/intro/100015201?code=UhApqgxa4VLIA591OKMTemuH1%2FWyLNNiHZ2CRYYdZzY%3D   
【Linux Namespace】https://www.cnblogs.com/sparkdev/p/9365405.html  
【浅谈 Linux Namespace】https://xigang.github.io/2018/10/14/namespace-md/   
【理解Docker（4）：Docker 容器使用 cgroups 限制资源使用 】https://www.cnblogs.com/sammyliu/p/5886833.html     





