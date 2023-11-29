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

在主节点挂掉之后，有选举节点功能会自动从从节点中选出一个新的主节点，如果一个从节点，从节点也会自动从集群中剔除，保证集群的数据读操作不受影响。   

搭建一个副本集集群最少需要三个节点：一个主节点，两个备份节点，如果三个节点分布合理，基本可以保证线上数据 `99.9%` 安全。       

<img src="/img/mongo/mongo-rc.jpg"  alt="mongo" />     

在集群只有是三个节点的情况下，当主节点超过配置的 `electionTimeoutMillis` 时间段（默认情况下为10 秒）内未与集合中的其他成员进行通信时，主节点就会被认为是异常了，两个副本节点也会进行选举，重新选出一个新的主节点。      

<img src="/img/mongo/mongo-rc-2.jpg"  alt="mongo" />     

在默认复制设置的情况下，从一个集群开始选举新主节点到选举完成的中间时间通常不超过 12 秒。中间包括将主节点标记不可用，发起并完成选举所需要的时间。在选举的过程中，就意味着集群暂时不能提供写入操作，时间越久集群不可写入的时间也就是越久。     

关于副本节点的属性，这里来主要的介绍下：priority、hidden、slaveDelay、tags、votes。     

- priority

对于副本节点，可以通过该属性增大或者减小该节点被选举为主节点的可能性，取值范围是 0-1000（如果是arbiters，则取值只有0或者1），数据越大，成为主节点的可能性越大，如果被配置为0，那么他就不能被选举成为主节点，而且也不能主动发起选举。    

比如说集群中的某几台机器配置较高，希望主节点主要在这几台机器中产生，那么我们就可以通过设置 priority 的大小来实现。   

- hidden  

隐藏节点可以从主节点同步数据，但对客户端不可见，在mongo shell 执行 db.isMaster() 方法也不会展示该节点，隐藏节点必须Priority为0，即不可以被选举成为主节点。但是如果有配置选举权限的话，可以参与选举。  

因为隐藏节点对客户端不可见，所以对于备份数据或者一些定时脚本可以直接连到隐藏节点，有大的慢查询也不会影响到集群本身对外提供的服务。   

<img src="/img/mongo/mongo-rc-3.jpg"  alt="mongo" />       

- slaveDelay

延迟从主节点同步数据，比如延迟节点时间配置为 1 小时，现在的时间是 10 点钟，那么从节点同步到的数据就是 9 点之前的数据。  

隐藏节点有什么作用呢？其中有一个和重要的作用就是防止数据库误操作，比如当我们对数据的进行大批量的删除或者更新操作，为了防止出现意外，我们可能会考虑事先备份一下数据，当操作出现异常的时候，我们还能根据备份进行复原回滚操作。有了延迟节点，因为延迟节点还没及时同步到最新的数据，我们就可以基于延迟节点进行数据库的复原操作。   

<img src="/img/mongo/mongo-rc-4.jpg"  alt="mongo" />       


### Sharding 分片模式


### 参考

【replication】https://www.mongodb.com/docs/manual/replication/     
【MongoDB 副本集之入门篇】https://jelly.jd.com/article/5f990ebbbfbee00150eb620a     