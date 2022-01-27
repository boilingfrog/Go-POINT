<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Redis 中数据的持久化](#redis-%E4%B8%AD%E6%95%B0%E6%8D%AE%E7%9A%84%E6%8C%81%E4%B9%85%E5%8C%96)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [AOF 持久化](#aof-%E6%8C%81%E4%B9%85%E5%8C%96)
    - [什么是 AOF 持久化](#%E4%BB%80%E4%B9%88%E6%98%AF-aof-%E6%8C%81%E4%B9%85%E5%8C%96)
    - [为什么要后记录日志呢](#%E4%B8%BA%E4%BB%80%E4%B9%88%E8%A6%81%E5%90%8E%E8%AE%B0%E5%BD%95%E6%97%A5%E5%BF%97%E5%91%A2)
    - [AOF 的潜在风险](#aof-%E7%9A%84%E6%BD%9C%E5%9C%A8%E9%A3%8E%E9%99%A9)
    - [AOF 中日志回写策略](#aof-%E4%B8%AD%E6%97%A5%E5%BF%97%E5%9B%9E%E5%86%99%E7%AD%96%E7%95%A5)
  - [RDB 持久化](#rdb-%E6%8C%81%E4%B9%85%E5%8C%96)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Redis 中数据的持久化

### 前言

我们知道 Redis 是内存数据库，所有操作都在内存上完成。内存的话，服务器断电，内存上面的数据就会丢失了。这个问题显然是需要解决的。    

Redis 中引入了持久化来避免数据的丢失，主要有两种持久化的方式 RDB 持久化和 AOF 持久化。  

### AOF 持久化

#### 什么是 AOF 持久化

AOF(Append Only File):通过保存数据库执行的命令来记录数据库的状态。  

<img src="/img/redis/redis-aof.png"  alt="redis" align="center" />

AOF日志对数据库命令的保存顺序是，Redis 先执行命令，把数据写入内存，然后才记录日志。   

#### 为什么要后记录日志呢

1、后写，能够避免记录到错误的命令。因为是先执行命令，后写入日志，只有命令执行成功了，命令才能被写入到日志中。  

2、避免阻塞当前的写操作，是在命令执行后才记录日志，所以不会阻塞当前的写操作。  

#### AOF 的潜在风险

- 1、如果命令执行成功，写入日志的时候宕机了，命令没有写入到日志中，这时候就有丢失数据的风险了，因为这时候没有写入日志，服务断电之后，这部分数据就丢失了。  

这种场景在别的地方也很常见，比如基于 MQ 实现分布式事务，也会出现`业务处理成功 + 事务消息发送失败`这种场景，[RabbitMQ，RocketMQ，Kafka 事务性，消息丢失和消息重复发送的处理策略](https://www.cnblogs.com/ricklz/p/15747565.html#%E5%9F%BA%E4%BA%8E-mq-%E5%AE%9E%E7%8E%B0%E7%9A%84%E5%88%86%E5%B8%83%E5%BC%8F%E4%BA%8B%E5%8A%A1)  

- 2、AOF 的日志写入也是在主线程进行的，如果磁盘的压力很大，写入速度变慢了，会影响后续的操作。   

这两种情况可以通过磁盘的回写时机来解决  

#### AOF 中日志回写策略



  



### RDB 持久化


### 参考

【Redis核心技术与实战】https://time.geekbang.org/column/intro/100056701    
【Redis设计与实现】https://book.douban.com/subject/25900156/  

                                 
