<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Redis 的高可用集群](#redis-%E7%9A%84%E9%AB%98%E5%8F%AF%E7%94%A8%E9%9B%86%E7%BE%A4)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [几种常用的集群方案](#%E5%87%A0%E7%A7%8D%E5%B8%B8%E7%94%A8%E7%9A%84%E9%9B%86%E7%BE%A4%E6%96%B9%E6%A1%88)
  - [主从集群模式](#%E4%B8%BB%E4%BB%8E%E9%9B%86%E7%BE%A4%E6%A8%A1%E5%BC%8F)
    - [全量同步](#%E5%85%A8%E9%87%8F%E5%90%8C%E6%AD%A5)
    - [增量同步](#%E5%A2%9E%E9%87%8F%E5%90%8C%E6%AD%A5)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Redis 的高可用集群

### 前言

这里来了解一下，Redis 中常见的集群方案  

### 几种常用的集群方案

- 主从集群模式  

- 哨兵机制  

- 切片集群(分片集群)  

### 主从集群模式 

主从集群，主从库之间采用的是读写分离  

- 主库：所有的写操作都在读库发生，然后主库同步数据到从库，同时也可以进行读操作；    

- 从库：只负责读操作；

<img src="/img/redis/redis-read-write.png"  alt="redis" align="center" />

主库需要复制数据到从库，主从双方的数据库需要保存相同的数据，将这种情况称为"数据库状态一致"  

来看下如何同步之前先来了解下几个概念  

- 1、服务器的运行ID(run ID)：每个 Redis 服务器在运行期间都有自己的`run ID`，`run ID`在服务器启动的时候自动生成。  

从服务器会记录主服务器的`run ID`，这样如果发生断网重连，就能判断新连接上的主服务器是不是上次的那一个，这样来决定是否进行数据部分重传还是完整重新同步。  

- 2、复制偏移量 offset：主服务器和从服务器都会维护一个复制偏移量  

主服务器每次向从服务器中传递 N 个字节的时候，会将自己的复制偏移量加上 N。  

从服务器中收到主服务器的 N 个字节的数据，就会将自己额复制偏移量加上 N。  

通过主从服务器的偏移量对比可以很清楚的知道主从服务器的数据是否处于一致。  

如果不一致就需要进行增量同步了，具体参加下文的增量同步   

#### 全量同步

从服务器首次加入主服务器中发生的是全量同步 

如何进行第一次同步？  

<img src="/img/redis/redis-read-write-copy.png"  alt="redis" align="center" />

1、从服务器连接到主服务器，然后发送 psync 到主服务器，因为第一次复制，不知道主库`run ID`,所以`run ID`为？；    

2、主服务器接收到同步的响应，回复从服务器自己的`run ID`和复制进行进度 offset；  

3、主服务器开始同步所有数据到从库中，同步依赖 RDB 文件，主库会通过 bgsave 命令，生成 RDB 文件，然后将 RDB 文件传送到从库中；  

4、从库收到  RDB 文件,清除自己的数据，然后载入 RDB 文件；  

5、主库在同步的过程中不会被阻塞，仍然能接收到命令，但是新的命令是不能同步到从库的，所以主库会在内存中用专门的 `replication buffer`，记录 RDB 文件生成后收到的所有写操作，然后在 RDB 文件，同步完成之后，再将`replication buffer`中的命令发送到从库中，这样就保证了从库的数据同步。  

#### 增量同步

如果主从服务器之间发生了网络闪断，从从服务将会丢失一部分同步的命令。  

在旧版本，`Redis 2.8`之前，如果发生了网络闪断，就会进行一次全量复制。  

在 2.8 版本之后，引入了增量同步的技术，这里主要是用到了 `repl_backlog_buffer`  

当主从库断连后，主库会把断连期间收到的写操作命令，写入`replication buffer`，同时也会把这些操作命令也写入`repl_backlog_buffer`这个缓冲区。

### 参考

【Redis核心技术与实战】https://time.geekbang.org/column/intro/100056701    
【Redis设计与实现】https://book.douban.com/subject/25900156/  