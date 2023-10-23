<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [MongoDB 中的锁](#mongodb-%E4%B8%AD%E7%9A%84%E9%94%81)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [MongoDB 中锁的类型](#mongodb-%E4%B8%AD%E9%94%81%E7%9A%84%E7%B1%BB%E5%9E%8B)
  - [锁的让渡释放](#%E9%94%81%E7%9A%84%E8%AE%A9%E6%B8%A1%E9%87%8A%E6%94%BE)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## MongoDB 中的锁

### 前言

MongoDB 是一种常见的文档型数据库，因为其高性能、高可用、高扩展性等特点，被广泛应用于各种场景。   

在多线程的访问下，可能会出现多线程同时操作一个集合的情况，进而出现数据冲突的情况，为了保证数据的一致性，MongoDB 采用了锁机制来保证数据的一致性。   

下面来看看 MongoDB 中的锁机制。   

### MongoDB 中锁的类型

MongoDB 中使用多粒度锁定，它允许操作锁定在全局，数据库或集合级别，同时允许各个存储引擎在集合级别一下实现自己的并发控制(例如，WiredTiger 中的文档级别)。

MongoDB 中使用一个 readers-writer 锁，它允许并发多个读操作访问数据库，但是只提供唯一写操作访问。    

当一个读锁存在时，其它的读操作可以继续，不会被阻塞，如果一个写锁占有这个资源的时候，其它所有的读操作和写操作都会被阻塞。也就是读读不阻塞，读写阻塞，写写阻塞。    

MongoDB 中的锁首先提供了读写锁，即共享锁（Shared, S）（读锁）以及排他锁（Exclusive, X）（写锁），同时，为了解决多层级资源之间的互斥关系，提高多层级资源请求的效率，还在此基础上提供了意向锁（Intent Lock）。即锁可以划分为4中类型：   

1、共享锁，读锁（S），允许多个线程同时读取一个集合，读读不互斥；   

2、排他锁，写锁（X），允许一个线程写入数据，写写互斥，读写互斥；   

3、意向共享锁，IS，表示意向读取；  

4、意向排他锁，IX，表示意向写入；  

什么是意向锁呢？

如果另一个任务企图在某表级别上应用共享或排他锁，则受由第一个任务控制的表级别意向锁的阻塞，第二个任务在锁定该表前不需要检查各个页或行锁，而只需检查表上的意向锁。   

简单的讲就是意向锁是为了快速判断，表里面是否有记录被加锁。     

例如，当以写入模式（X模式）锁定集合时，相应的数据库锁和全局锁都必须以意向独占（IX）模式锁定。一个数据库可以同时以IS和IX模式进行锁定，但独占（X）锁无法与其他模式并存，共享（S）锁只能与意向共享（IS）锁并存。     

MongoDB 中的锁是公平的，所有的请求都会排队获取相应的锁。但是 MongoDB 为了优化吞吐量，在执行某个请求的时候，会同时执行和它兼容的其它请求。比如一个队列一个请求队列需要的锁如下，执行IS请求的同时，会同时执行和它相容的其他S和IS请求。等这一批请求的S锁释放后，再执行X锁的请求。   

```
IS → IS → X → X → S → IS
```

这种处理机制保证了在相对公平的前提下，提高了吞吐量，不会让某一类的请求长时间的等待。    

对于长时间的读或者写操作，某些条件下，mongodb 会临时的 让渡 锁，以防止长时间的阻塞。

### 锁的让渡释放   

对于大多数读取和写入操作，WiredTiger 使用乐观并发控制。WiredTiger 仅在全局、数据库和集合级别使用意向锁。当存储引擎检测到两个操作之间的冲突时，其中一个将导致写冲突，从而使 MongoDB 可以在不可见的情况下重新尝试该操作。  

在某些情况下，读写操作可以释放它们持有的锁。以防止长时间的阻塞。   

长时间运行的读取和写入操作，比如查询、更新和删除，在许多情况下都会释放锁。MongoDB操作也可以在影响多个文档的写入操作中，在单个文档修改之间释放锁。  

对于支持文档级并发控制的存储引擎，比如 WiredTiger，在访问存储时通常不需要释放锁，因为在全局、数据库和集合级别保持的意向锁不会阻塞其他读取者和写入者。然而，操作会定期释放锁，以便：

1、避免长时间存储事务，因为这些事务可能需要在内存中保存大量数据；  

2、充当中断点，以便可以终止长时间运行的操作；  

3、允许需要对集合进行独占访问的操作，比如索引/集合的删除和创建。  

### 常见操作使用的锁类型

下面列列举一下 MongoDB 中的常见查询对应的锁类型   

### 如果定位 MongoDB 中锁操作  

当查询有慢查询出现的时候，有时候会出现锁的阻塞等待，紧急情况，需要我们快速定位并且结束当前操作。   

使用 `db.currentOp()` 就能查看当前数据库正在执行的操作。   

```
db.currentOp()

{
    "inprog" : [
        {
            "opid" : 6222,   #进程号
            "active" : true, #是否活动状态
            "secs_running" : 3,#操作运行了多少秒
            "microsecs_running" : NumberLong(3662328),#操作持续时间(以微秒为单位)。MongoDB通过从操作开始时减去当前时间来计算这个值。
            "op" : "getmore",#操作类型，包括(insert/query/update/remove/getmore/command)
            "ns" : "local.oplog.rs",#命名空间
            "query" : {#如果op是查询操作，这里将显示查询内容；也有说这里显示具体的操作语句的
                 
            },
            "client" : "192.168.91.132:45745",#连接的客户端信息
            "desc" : "conn5",#数据库的连接信息
            "threadId" : "0x7f1370cb4700",#线程ID
            "connectionId" : 5,#数据库的连接ID
            "waitingForLock" : false,#是否等待获取锁
            "numYields" : 0,
            "lockStats" : {
                "timeLockedMicros" : {#持有的锁时间微秒
                    "r" : NumberLong(141),#整个MongoDB实例的全局读锁
                    "w" : NumberLong(0)#整个MongoDB实例的全局写锁
                },
                "timeAcquiringMicros" : {#为了获得锁，等待的微秒时间
                    "r" : NumberLong(16),#整个MongoDB实例的全局读锁
                    "w" : NumberLong(0)#整个MongoDB实例的全局写锁
                }
            }
        }
    ]
}
```

来看下几个主要的字段含义    

- client：发起请求的客户端；  

- opid: 操作的唯一标识；  

- secs_running：该操作已经执行的时间，单位：微妙。如果该字段的返回值很大，就需要查询请求是否合理；   
 
- op：操作类型。通常是query、insert、update、delete、command中的一种；  

- query/ns：这个字段能看出是对哪个集合正在执行什么操作。    

当当发现




### 参考

【mongodb锁表命令-相关文档】https://www.volcengine.com/theme/900385-M-7-1   
【mongo 中的锁】https://www.jinmuinfo.com/community/MongoDB/docs/15-faq/03-concurrency.html  
【FAQ: Concurrency】https://www.mongodb.com/docs/manual/faq/concurrency/      