<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [RabbitMQ 如何做分布式](#rabbitmq-%E5%A6%82%E4%BD%95%E5%81%9A%E5%88%86%E5%B8%83%E5%BC%8F)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [集群配置方案](#%E9%9B%86%E7%BE%A4%E9%85%8D%E7%BD%AE%E6%96%B9%E6%A1%88)
    - [cluster](#cluster)
      - [普通模式](#%E6%99%AE%E9%80%9A%E6%A8%A1%E5%BC%8F)
      - [镜像模式](#%E9%95%9C%E5%83%8F%E6%A8%A1%E5%BC%8F)
    - [federation](#federation)
    - [shovel](#shovel)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## RabbitMQ 如何做分布式

### 前言

前面几篇文章介绍了消息队列中遇到的问题，这篇来聊聊 RabbitMQ 的集群搭建。    

### 集群配置方案

RabbitMQ 中集群的部署方案有三种 cluster,federation,shovel。    

#### cluster

cluster 有两种模式，分别是普通模式和镜像模式   

##### 普通模式

cluster 普通模式(默认的集群模式)，所有节点中的元数据是一致的，RabbitMQ 中的元数据会被复制到每一个节点上。  

队列里面的数据只会存在创建它的节点上，其他节点除了存储元数据，还存储了指向 Queue 的主节点(owner node)的指针。  

集群中节点之间没有主从节点之分。      

<img src="/img/mq-rabbitmq-cluster.png"  alt="mq" align="center" />

举个栗子来说明下普通模式的消息传输：  

假设我们 RabbitMQ 中有是三个节点，分别是 `node1,node2,node3`。如果队列 queue1 的连接创建发生在 node1 中，那么该队列的元数据会被同步到所有的节点中，但是 queue1 中的消息，只会在 node1 中。    

- 如果一个消费者通过 node2 连接，然后来消费 queue1 中的消息?  

RabbitMQ 会临时在 node1、node2 间进行消息传输，因为非 owner 节点除了存储元数据，还会存储指向 Queue 的主节点(owner node)的指针。RabbitMQ 会根据这个指向，把 node1 中的消息实体取出并经过 node2 发送给 consumer 。 

- 如果一个生产者通过 node2 连接，然后来向 queue1 中生产数据?  
 
同理，RabbitMQ 会根据 node2 中的主节点(owner node)的指针，把消息转发送给 owner 节点 node1,最后插入的数据还是在 node1 中。  
 
<img src="/img/mq-rabbitmq-cluster-data.png"  alt="mq" align="center" />  

同时对于队列的创建，要平均的落在每个节点上，如果只在一个节点上创建队列，所有的消费，最终都会落到这个节点上，会产生瓶颈。     

存在的问题：  

如果 node1 节点故障了，那么 node2 节点无法取出 node1 中还未消费的消息实体。  

1、如果做了队列的持久化，消息不会被丢失，等到 node1 恢复了，就能接着进行消费，但是在恢复之前其他节点不能创建 node1 中已将创建的队列。   

2、如果没有做持久化，消息会丢失，但是 node1 中的队列，可以在其他节点重新创建，不用等待 node1 的恢复。   

普通模式不支持消息在每个节点上的复制，当然 RabbitMQ 中也提供了支持复制的模式，就是镜像模式(参见下文)。  

##### 镜像模式

镜像队列会在节点中同步队列的数据，最终的队列数据会存在于每个节点中，而不像普通模式中只会存在于创建它的节点中。  

优点很明显，当有主机宕机的时候，因为队列数据会同步到所有节点上，避免了普通模式中的单点故障。  

缺点就是性能不好，集群内部的同步通讯会占用大量的网络带宽，适合一些可靠性要求比较高的场景。   

针对镜像模式 RabbitMQ 也提供了几种模式，有效值为 `all，exactly，nodes` 默认为 all。  

- all 表示集群中所有的节点进行镜像；  

- exactly 表示指定个数的节点上进行镜像，节点个数由`ha-params`指定;  

- nodes 表示在指定的节点上进行镜像，节点名称由`ha-params`指定;    

所以针对普通队列和镜像队列，我们可以选择其中几个队列作为镜像队列，在性能和可靠性之间找到一个平衡。   

关于镜像模式中消息的复制，这里也用的很巧妙，值得借鉴  

1、master 节点向 slave 节点同步消息是通过组播 GM(Guaranteed Multicast) 来同步的。  

2、所有的消息经过 master 节点，master 对消息进行处理，同时也会通过 GM 广播给所有的 slave，slave收到消息之后在进行数据的同步操作。   

3、GM 实现的是一种可靠的组播通信协议，该协议能保证组播消息的原子性。具体如何实现呢？  

它的实现大致为:将所有的节点形成一个循环链表，每个节点都会监控位于自己左右两边的节点，当有节点新增时，相邻的节点保证当前广播的消息会复制到新的节点上 当有节点失效时，相邻的节点会接管以保证本次广播的消息会复制到所有的节点。  

因为是一个循环链表，所以 master 发出去的消息最后也会返回到 master 中，master 如果收到了自己发出的操作命令，这时候就可以确定命令已经同步到了所有的节点。  

<img src="/img/mq-rabbitmq-cluster-mirror.png"  alt="mq" align="center" />

#### federation

federation 插件的设计目标是使 RabbitMQ 在不同的 Broker 节点之间进行消息传递而无需建立集群。  

看了定义还是很迷糊，来举举栗子吧   

假设我们有一个 RabbitMQ 的集群，分别部署在不同的城市，那么我们假定分别是在北京，上海，广州。  

<img src="/img/mq-rabbitmq-federation.png"  alt="mq" align="center" />

如果一个现在有一个业务 clientA，部署的机器在北京，然后连接到北京节点的 broker1 。然后网络连通性也很好，发送消息到 broker1 中的 exchangeA 中，消息能够很快的发送到，就算在开启了 `publisher confirm` 机制或者事务机制的情况下，也能快速确认信息，这种情况下是没有问题的。  

如果一个现在有一个业务 clientB，部署的机器在上海，然后连接到北京节点的 broker1 。然后网络连通性不好，发送消息到 broker1 中的 exchangeA 中，因为网络不好，所以消息的确认有一定的延迟，这对于我们无疑使灾难，消息量大情况下，必然造成数据的阻塞，在开启了 `publisher confirm` 机制或者事务机制的情况下，这种情况将会更严重。     

当然如果把 clientB ，部署在北京的机房中，这个问题就解决了，但是多地容灾就不能实现了。   

针对这种情况如何解决呢，这时候 federation 就登场了。  

比如位于上海的业务 clientB，连接北京节点的 broker1。然后发送消息到 broker1 中的 exchangeA 中。这时候是存在网络连通性的问题的。   

- 1、让上海的业务 clientB，连接上海的节点 broker2；  

- 2、通过 Federation ，在北京节点的 broker1 和上海节点的 broker2 之间建立一条单向的 `Federation link`；  

- 3、Federation 插件会在上海节点的 broker2 中创建一个同名的交换器 exchangeA (具体名字可配置，默认同名), 同时也会创建一个内部交换器，通过路由键 rkA ,将这两个交换器进行绑定，同时也会在 broker2 中创建一个

    1、Federation 插件会在上海节点的 broker2 中创建一个同名的交换器 exchangeA (具体名字可配置，默认同名)；  
    
    2、Federation 插件会在上海节点的 broker2 中创建一个内部交换器，通过路由键 rkA ,将 exchangeA 和内部交换器进行绑定；   
    
    3、Federation 插件会在上海节点的 broker2 中创建队列，和内部交换器进行绑定，同时这个队列会和北京节点的 broker1 中的 exchangeA，建立一条 AMQP 链接，来实时的消费队列中的消息了；  

- 4、经过上面的流程，就相当于在上海节点 broker2 中的 exchangeA 和北京节点 broker1 中的 exchangeA 建立了`Federation link`；  

这样位于上海的业务 clientB 链接到上海的节点 broker2，然后发送消息到该节点中的 exchangeA，这个消息会通过`Federation link`，发送到北京节点 broker1 中的 exchangeA，所以可以减少网络连通性的问题。   


#### shovel

### 参考

【RabbitMQ分布式集群架构和高可用性（HA）】http://chyufly.github.io/blog/2016/04/10/rabbitmq-cluster/   
【RabbitMQ分布式部署方案简介】https://www.jianshu.com/p/c7a1a63b745d   
【RabbitMQ实战指南】https://book.douban.com/subject/27591386/      
【RabbitMQ两种集群模式配置管理】https://blog.csdn.net/fgf00/article/details/79558498    


