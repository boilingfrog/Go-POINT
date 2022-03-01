<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [使用 Redis 实现消息队列](#%E4%BD%BF%E7%94%A8-redis-%E5%AE%9E%E7%8E%B0%E6%B6%88%E6%81%AF%E9%98%9F%E5%88%97)
  - [基于List的消息队列](#%E5%9F%BA%E4%BA%8Elist%E7%9A%84%E6%B6%88%E6%81%AF%E9%98%9F%E5%88%97)
  - [基于 Streams 的消息队列](#%E5%9F%BA%E4%BA%8E-streams-%E7%9A%84%E6%B6%88%E6%81%AF%E9%98%9F%E5%88%97)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 使用 Redis 实现消息队列

Redis 中也是可以实现消息队列  

不过谈到消息队列，我们会经常遇到下面的几个问题  

1、消息如何防止丢失；    

2、消息的重复发送如何处理；  

3、消息的顺序性问题；  

关于 mq 中如何处理这几个问题，可参看[RabbitMQ，RocketMQ，Kafka 事务性，消息丢失，消息顺序性和消息重复发送的处理策略](https://www.cnblogs.com/ricklz/p/15747565.html)  

### 基于List的消息队列

对于 List  

使用 LPUSH 写入数据，使用 RPOP 读出数据  

```
127.0.0.1:6379> LPUSH test "ceshi-1"
(integer) 1

127.0.0.1:6379> RPOP test
"ceshi-1"
```

使用 RPOP 客户端就需要一直轮询，来监测是否有值可以读出，可以使用 BRPOP 可以进行阻塞式读取，客户端在没有读到队列数据时，自动阻塞，直到有新的数据写入队列，再开始读取新数据。  

```
127.0.0.1:6379> BRPOP test 10
```

后面的 10 是监听的时间，单位是秒，10秒没数据，就退出。  

如果客户端从队列中拿到一条消息时，但是还没消费，客户端宕机了，这条消息就对应丢失了， Redis 中为了避免这种情况的出现，提供了 BRPOPLPUSH 命令，BRPOPLPUSH 会在消费一条消息的时候，同时把消息插入到另一个 List，这样如果消费者程序读了消息但没能正常处理，等它重启后，就可以从备份 List 中重新读取消息并进行处理了。  

不过 List 类型并不支持消费组的实现,Redis 从 5.0 版本开始提供的 Streams 数据类型，来支持消息队列的场景。  

### 基于 Streams 的消息队列

Streams 是 Redis 专门为消息队列设计的数据类型。  










### 参考

【Redis核心技术与实战】https://time.geekbang.org/column/intro/100056701    
【Redis设计与实现】https://book.douban.com/subject/25900156/  
【Redis Streams 介绍】http://www.redis.cn/topics/streams-intro.html    
