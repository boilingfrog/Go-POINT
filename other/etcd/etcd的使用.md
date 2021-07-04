<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [etcd的使用](#etcd%E7%9A%84%E4%BD%BF%E7%94%A8)
  - [什么是etcd](#%E4%BB%80%E4%B9%88%E6%98%AFetcd)
  - [etcd的特点](#etcd%E7%9A%84%E7%89%B9%E7%82%B9)
  - [etcd的应用场景](#etcd%E7%9A%84%E5%BA%94%E7%94%A8%E5%9C%BA%E6%99%AF)
    - [服务注册与发现](#%E6%9C%8D%E5%8A%A1%E6%B3%A8%E5%86%8C%E4%B8%8E%E5%8F%91%E7%8E%B0)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## etcd的使用

### 什么是etcd

ETCD是一个分布式、可靠的`key-value`存储的分布式系统，用于存储分布式系统中的关键数据；当然，它不仅仅用于存储，还提供配置共享及服务发现；基于Go语言实现  。

### etcd的特点

- 完全复制：集群中的每个节点都可以使用完整的存档

- 高可用性：Etcd可用于避免硬件的单点故障或网络问题

- 一致性：每次读取都会返回跨多主机的最新写入

- 简单：包括一个定义良好、面向用户的API（gRPC）

- 安全：实现了带有可选的客户端证书身份验证的自动化TLS

- 可靠：使用Raft算法实现了强一致、高可用的服务存储目录

### etcd的应用场景

#### 服务注册与发现

服务发现还能注册

服务注册发现解决的是分布式系统中最常见的问题之一，即在同一个分布式系统中，找到我们需要的目标服务，建立连接，然后完成整个链路的调度。   

本质上来说，服务发现就是想要了解集群中是否有进程在监听 udp 或 tcp 端口，并且通过名字就可以查找和连接。要解决服务发现的问题，需要有下面三大支柱，缺一不可。  

1、**一个强一致性、高可用的服务存储目录**。基于Raft算法的etcd天生就是这样一个强一致性高可用的服务存储目录。  

2、**一种注册服务和监控服务健康状态的机制**。用户可以在etcd中注册服务，并且对注册的服务设置`key TTL`，定时保持服务的心跳以达到监控健康状态的效果。   

3、**一种查找和连接服务的机制**。通过在 etcd 指定的主题下注册的服务也能在对应的主题下查找到。为了确保连接，我们可以在每个服务机器上都部署一个Proxy模式的etcd，这样就可以确保能访问etcd集群的服务都能互相连接。    

### 参考

【一文入门ETCD】https://juejin.cn/post/6844904031186321416   
【etcd：从应用场景到实现原理的全方位解读】https://www.infoq.cn/article/etcd-interpretation-application-scenario-implement-principle   
【Etcd 架构与实现解析】http://jolestar.com/etcd-architecture/   