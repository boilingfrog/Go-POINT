<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [MongoDB 中的分布式集群架构](#mongodb-%E4%B8%AD%E7%9A%84%E5%88%86%E5%B8%83%E5%BC%8F%E9%9B%86%E7%BE%A4%E6%9E%B6%E6%9E%84)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [Replica Set 副本集模式](#replica-set-%E5%89%AF%E6%9C%AC%E9%9B%86%E6%A8%A1%E5%BC%8F)
  - [Sharding 分片模式](#sharding-%E5%88%86%E7%89%87%E6%A8%A1%E5%BC%8F)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## MongoDB 中的分布式集群架构

### 前言

前面我们了解了 MongoDB 中的索引，事务，锁等知识点。线上使用的 MongoDB 大部分的场景我们都会考虑使用分布式结构，这里我们来了解一下 MongoDB 中的分布式架构。   

MongoDB 中常用的分布式架构有下面几种：  

1、Replica Set 副本集模式：一个 Primary 节点用于写入数据，其它的 Secondary 节点用于查询数据，适合读写少的场景，是目前较为主流的架构方式，Primary 节点挂了，会自动从 Secondary 节点选出新的 Primary 节点，提供数据写入操作；  

2、Master-Slaver 主从副本的模式：也是主节点写入，数据同步到 Slave 节点，Slave 节点提供数据查询，最大的问题就是可用性差，`MongoDB 3.6` 起已不推荐使用主从模式，自 `MongoDB 3.2` 起，分片群集组件已弃用主从复制。因为 `Master-Slave` 其中 Master 宕机后不能自动恢复，只能靠人为操作，可靠性也差，操作不当就存在丢数据的风险，这种模式被 `Replica Set` 所替代 ；     

3、Sharding 分片模式：将不同的数据分配在不同的机器中，也就是数据的横向扩展，单个机器只存储整个数据中的一部分，这样通过横向增加机器的数量来提高集群的存储和计算能力。    

因为 `Master-Slaver` 模式已经在新版本中弃用了，下面主要来介绍下 `Replica Set` 模式和 `Sharding` 模式。   

### Replica Set 副本集模式

MongoDB 中的 `Replica Set` 副本集模式，可以简单理解为一主多从的集群，包括一个主节点（primary）和多个副本节点（Secondaries）。   

主节点只有一个，所有的写操作都在主节点中进行，副本节点可以有多个，通过同步主节点的操作日志（oplog）来备份主节点数据。   

在主节点挂掉之后，有选举节点功能



### Sharding 分片模式


### 参考

【replication】https://www.mongodb.com/docs/manual/replication/     
【MongoDB 副本集之入门篇】https://jelly.jd.com/article/5f990ebbbfbee00150eb620a     