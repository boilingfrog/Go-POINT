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

RabbitMQ 会临时在 node1、node2 间进行消息传输，因为非 owner 节点除了存储元数据，还会存储指向 Queue 的主节点(owner node)的指针。RabbitMQ 把 node1 中的消息实体取出并经过 node2 发送给 consumer 。  


##### 镜像模式





#### federation

#### shovel




### 参考

【RabbitMQ分布式集群架构和高可用性（HA）】http://chyufly.github.io/blog/2016/04/10/rabbitmq-cluster/   
【RabbitMQ分布式部署方案简介】https://www.jianshu.com/p/c7a1a63b745d   
【RabbitMQ实战指南】https://book.douban.com/subject/27591386/      
【RabbitMQ两种集群模式配置管理】https://blog.csdn.net/fgf00/article/details/79558498    


