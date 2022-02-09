<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Redis 中数据的持久化](#redis-%E4%B8%AD%E6%95%B0%E6%8D%AE%E7%9A%84%E6%8C%81%E4%B9%85%E5%8C%96)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [AOF 持久化](#aof-%E6%8C%81%E4%B9%85%E5%8C%96)
    - [什么是 AOF 持久化](#%E4%BB%80%E4%B9%88%E6%98%AF-aof-%E6%8C%81%E4%B9%85%E5%8C%96)
    - [为什么要后记录日志呢](#%E4%B8%BA%E4%BB%80%E4%B9%88%E8%A6%81%E5%90%8E%E8%AE%B0%E5%BD%95%E6%97%A5%E5%BF%97%E5%91%A2)
    - [AOF 的潜在风险](#aof-%E7%9A%84%E6%BD%9C%E5%9C%A8%E9%A3%8E%E9%99%A9)
    - [AOF 文件的写入和同步](#aof-%E6%96%87%E4%BB%B6%E7%9A%84%E5%86%99%E5%85%A5%E5%92%8C%E5%90%8C%E6%AD%A5)
    - [AOF 文件重写机制](#aof-%E6%96%87%E4%BB%B6%E9%87%8D%E5%86%99%E6%9C%BA%E5%88%B6)
    - [AOF 的数据还原](#aof-%E7%9A%84%E6%95%B0%E6%8D%AE%E8%BF%98%E5%8E%9F)
  - [RDB 持久化](#rdb-%E6%8C%81%E4%B9%85%E5%8C%96)
    - [什么是 RDB 持久化](#%E4%BB%80%E4%B9%88%E6%98%AF-rdb-%E6%8C%81%E4%B9%85%E5%8C%96)
    - [RDB 如何做内存快照](#rdb-%E5%A6%82%E4%BD%95%E5%81%9A%E5%86%85%E5%AD%98%E5%BF%AB%E7%85%A7)
    - [快照时发生数据修改](#%E5%BF%AB%E7%85%A7%E6%97%B6%E5%8F%91%E7%94%9F%E6%95%B0%E6%8D%AE%E4%BF%AE%E6%94%B9)
    - [多久做一次快照](#%E5%A4%9A%E4%B9%85%E5%81%9A%E4%B8%80%E6%AC%A1%E5%BF%AB%E7%85%A7)
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

这两种情况可通过调整 AOF 文件的写入磁盘的时机来避免  

#### AOF 文件的写入和同步

AOF 文件持久化的功能分成三个步骤，文件追加(append),文件写入，文件同步(sync)。  

AOF 文件在写入磁盘之前是先写入到 aof_buf 缓冲区中，然后通过调用 flushAppendOnlyFile 将缓冲区中的内容保存到 AOF 文件中。  

写入的策略通过 appendfsync 来进行配置  

- Always：同步写回 每次写操作命令执行完后，同步将 AOF 日志数据写回硬盘；  

- Everysec：每秒写回 每次写操作命令执行完后，先将命令写入到 AOF 文件的内核缓冲区，然后每隔一秒将缓冲区里的内容写回到硬盘；

- No：操作系统控制的写回 Redis 不在控制命令的写会时机，交由系统控制。每次写操作命令执行完成之后，命令会被放入到 AOF 文件的内核缓冲区，之后什么时候写入到磁盘，交由系统控制。  

#### AOF 文件重写机制

因为每次执行的命令都会被写入到 AOF 文件中，随着系统的运行，越来越多的文件会被写入到 AOF 文件中，这样 AOF 文件势必会变得很大，这种情况该如何去处理呢？   

为了解决这种情况，Redis 中引入了重写的机制  

什么是重写呢？  

因为 AOF 文件中记录的是每个命令的操作记录，举个🌰，比如当一个键值对被多条写命令反复修改时，AOF文件会记录相应的多条命令，那么重写机制，就是根据这个键值对当前的最新状态，为它生成对应的写入命令，保存成一行操作命令。这样就精简了 AOF 文件的大小。       

```
192.168.56.118:6379> set name "xiaoming"
OK
192.168.56.118:6379> get name
"xiaoming"
192.168.56.118:6379> set name "xiaozhang"
OK
192.168.56.118:6379> set name "xiaoli"
OK

# 重写后就是
192.168.56.118:6379> set name "xiaoli"
```

简单来讲就是多变一，就是把 AOF 中日志根据当前键值的状态，合并成一条操作命令。  

重写之后的文件会保存到新的 AOF 文件中，这时候旧的 AOF 文件和新的 AOF 文件中键值对的状态是一样的。然后新的 AOF 文件会替换掉旧的 AOF 文件，这样 重写操作一直在进行，AOF 文件就不至于变的过大。  

重写是后台进行的， AOF 的重写会放到子进程中进行的，使用子进程的优点：   

1、子进程处理 AOF 期间，不会影响 Redis 主线程对数据的处理；  

2、子进程拥有所在线程的数据副本，使用进程能够避免锁的使用，保证数据的安全。   

这里来看下，AOF 的处理流程  

AOF 重写也有一个缓冲区，当服务节接收到新的命令的是，如果在正在进行 AOF 重写，命令同样也会被发送到 AOF 缓冲区   

<img src="/img/redis/redis-aof-rewrite.png"  alt="redis" align="center" />

子进程执行 AOF 重写的过程,服务端进程主要处理以下内容   

1、接收并处理客户端发送的命令；  

2、将执行后的命令写入到 AOF 缓冲区；  

3、将执行后的命令也写入到 AOF 重写缓冲区；   

AOF 缓冲区和  AOF 重写缓冲区中的内容会被定期的同步到 AOF 文件和 AOF 重写文件中  

当子进程完成重写的时候，会给父进程发送一个信号，这时候父进程主要主要进行下面的两步操作：  

1、将 AOF 重写缓冲区中的内容全部写入到 AOF 重写文件中，这时候重写 AOF 文件保存的数据状态是和服务端数据库的状态一致的；  

2、将 AOF 重写文件替换旧的 AOF 文件；  

通过 AOF 的重写操作，新的 AOF 文件不断的替换旧的 AOF 文件，这样就能控制 AOF 文件的大小  

#### AOF 的数据还原

AOF 文件包了重建数据库索引锁需要的全部命令，所以只需要读入并重新执行一遍 AOF 文件中保存的命令，即可还原服务关闭之前数据库的状态。  
 
### RDB 持久化

#### 什么是 RDB 持久化  

RDB(Redis database)：实现方式是将存在 Redis 内存中的数据写入到 RDB 文件中保存到磁盘上从而实现持久化的。  

和 AOF 不同的是 RDB 保存的是数据而不是操作，在进行数据恢复的时候，直接把 RDB 的文件读入到内存，即可完成数据恢复。  

<img src="/img/redis/redis-rdb.png"  alt="redis" align="center" />

#### RDB 如何做内存快照

Redis 中对于如何备份数据到 RDB 文件中，提供了两种方式  

- 1、save: 在主线程中执行，不过这种会阻塞 Redis 服务进程；  

- 2、bgsave: 主线程会 fork 出一个子进程来负责处理 RDB 文件的创建，不会阻塞主线程的命令操作，这也是 Redis 中 RDB 文件生成的默认配置；  

对于 save 和 bgsave 这两种快照方式，服务端是禁止这两种方式同时执行的，防止产生竞争条件。  

Redis 中可以使用 save 选项，来配置服务端执行 BGSAVE 命令的间隔时间   

```
#
# Save the DB on disk:
#
#   save <seconds> <changes>
#
#   Will save the DB if both the given number of seconds and the given
#   number of write operations against the DB occurred.
#
#   In the example below the behaviour will be to save:
#   after 900 sec (15 min) if at least 1 key changed
#   after 300 sec (5 min) if at least 10 keys changed
#   after 60 sec if at least 10000 keys changed
#
#   Note: you can disable saving completely by commenting out all "save" lines.
#
#   It is also possible to remove all the previously configured save
#   points by adding a save directive with a single empty string argument
#   like in the following example:
#
#   save ""

save 900 1
save 300 10
save 60 10000
```

`save 900 1` 就是服务端在900秒，读数据进行了至少1次修改，就会触发一次 BGSAVE 命令  

`save 300 10` 就是服务端在300秒，读数据进行了至少10次修改，就会触发一次 BGSAVE 命令    

#### 快照时发生数据修改

举个栗子🌰：我们在t时刻开始对内存数据内进行快照，假定目前有 2GB 的数据需要同步，磁盘写入的速度是 `0.1GB/s` 那么，快照的时间就是 20s，那就是在 `t+20s` 完成快照。     

如果在 t+6s 的时候修改一个还没有写入磁盘的内存数据 test 为 test-hello。那么就会破坏快照的完整性了，因为 t 时刻备份的数据已经被修改了。当然是希望在备份期间数据不能被修改。  

如果不能被修改，就意味这在快照期间不能对数据进行修改操作，就如上面的栗子，快照需要进行20s,期间不允许处理数据更新操作，这显然也是不合理的。  

这里需要聊一下 bgsave 是可以避免阻塞，不过需要注意的是避免阻塞和正常读写操作是有区别的。避免阻塞主线程确实没有阻塞可以处理读操作，但是为了保护快照的完整性，是不能修改快照期间的数据的。  

这里就需要引入一种新的处理方案，写时复制技术（Copy-On-Write, COW），在执行快照的同时，正常处理写操作。  

bgsave 子进程是由主线程 fork 生成的，所以是可以共享主线程的内存的，bgsave子进程运行后会读取主线程中的内存数据，并且写入到 RDB 文件中。  

写复制技术就是，如果主线程在内存快照期间修改了一块内存，那么这块内存会被复制一份，生成该数据的副本，然后 bgsave 子进程在把这段内存写入到 RDB 文件中。这样就可以在快照期间进行数据的修改了。   

<img src="/img/redis/redis-rdb-cow.png"  alt="redis" align="center" />

#### 多久做一次快照

对于快照，如果做的太频繁，可能会出现前一次快照还没有处理完成，后面的快照数据马上就进来了，同时过于频繁的快照也会增加磁盘的压力。   

如果间隔时间过久，服务器在两次快照期间宕机，丢失的数据大小会随着快照间隔时间的增长而增加。   

是否可以选择增量式快照呢？选择增量式快照，我们就需要记住每个键值对的状态，如果键值对很多，同样也会引入很多内存空间，这对于内存资源宝贵的Redis来说，有些得不偿失。  

相较于 AOF 来对比，RDB 是会在数据恢复时，速度更快。但是 RDB 的内存快照同步频率不太好控制，过多过少都有问题。  

Redis 4.0中提出了一个混合使用 AOF 日志和内存快照的方法。简单来说，内存快照以一定的频率执行，在两次快照之间，使用AOF日志记录这期间的所有命令操作。  

通过混合使用AOF日志和内存快照的方法，RDB 快照的频率不需要过于频繁，在两次 RDB 快照期间，使用 AOF 日志来记录，这样也不用考虑 AOF 的文件过大问题，在下一次 RDB 快照开始的时候就可以删除 AOF 文件了。  

<img src="/img/redis/redis-aof-and-rdb.png"  alt="redis" align="center" />

### 参考

【Redis核心技术与实战】https://time.geekbang.org/column/intro/100056701    
【Redis设计与实现】https://book.douban.com/subject/25900156/  
【过期键与持久化】https://segmentfault.com/a/1190000017526315    

                                 
