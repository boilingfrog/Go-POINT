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

- 是可持久化的，可以保证数据不丢失。  

- 支持消息的多播、分组消费。  

- 支持消息的有序性。

来看下几个主要的命令  

```
XADD：插入消息，保证有序，可以自动生成全局唯一ID；

XREAD：用于读取消息，可以按ID读取数据； 

XREADGROUP：按消费组形式读取消息；

XPENDING和XACK：XPENDING命令可以用来查询每个消费组内所有消费者已读取但尚未确认的消息，而XACK命令用于向消息队列确认消息处理已完成。  
```

下面看几个常用的命令  

**XADD**

使用 XADD 向队列添加消息，如果指定的队列不存在，则创建一个队列，XADD 语法格式：  

```
$ XADD key ID field value [field value ...]
```

- key：队列名称，如果不存在就创建  

- ID：消息 id，我们使用 * 表示由 redis 生成，可以自定义，但是要自己保证递增性  

- field value：记录   

```
$ XADD teststream * name xiaohong surname xiaobai
"1646650328883-0"
```

可以看到 `1646650328883-0`就是自动生成的全局唯一消息ID   

**XREAD**

使用 XREAD 以阻塞或非阻塞方式获取消息列表  

```
$ XREAD [COUNT count] [BLOCK milliseconds] STREAMS key [key ...] id [id ...]
```

- count：数量  

- milliseconds：可选，阻塞毫秒数，没有设置就是非阻塞模式  

- key：队列名  

- id：消息 ID  

```
$ XREAD BLOCK 100 STREAMS  teststream 0
1) 1) "teststream"
   2) 1) 1) "1646650328883-0"
         2) 1) "name"
            2) "xiaohong"
            3) "surname"
            4) "xiaobai"
```

BLOCK 就是阻塞的毫秒数  

**XGROUP**

使用 XGROUP CREATE 创建消费者组  

```
$ XGROUP [CREATE key groupname id-or-$] [SETID key groupname id-or-$] [DESTROY key groupname] [DELCONSUMER key groupname consumername]
```

- key：队列名称，如果不存在就创建  

- groupname：组名  

- $：表示从尾部开始消费，只接受新消息，当前 Stream 消息会全部忽略  

从头开始消费  

```
$ XGROUP CREATE teststream test-consumer-group-name 0-0  
```

从尾部开始消费  

```
$ XGROUP CREATE teststream test-consumer-group-name $
```

**XREADGROUP GROUP**

使用 `XREADGROUP GROUP` 读取消费组中的消息  

```
$ XREADGROUP GROUP group consumer [COUNT count] [BLOCK milliseconds] [NOACK] STREAMS key [key ...] ID [ID ...]
```

- group：消费组名  

- consumer：消费者名  

- count：读取数量

- milliseconds：阻塞毫秒数  

- key：队列名  

- ID：消息 ID

```
$ XADD teststream * name xiaohong surname xiaobai
"1646653392799-0"

$ XREADGROUP GROUP test-consumer-group-name test-consumer-name COUNT 1 STREAMS teststream >
1) 1) "teststream"
   2) 1) 1) "1646653392799-0"
         2) 1) "name"
            2) "xiaohong"
            3) "surname"
            4) "xiaobai"
```

消息队列中的消息一旦被消费组里的一个消费者读取了，就不能再被该消费组内的其他消费者读取了。  

如果没有通过 XACK 命令告知消息已经成功消费了，该消息会一直存在，可以通过 XPENDING 命令查看已读取、但尚未确认处理完成的消息。   

```
$ XPENDING teststream test-consumer-group-name
1) (integer) 3
2) "1646653325535-0"
3) "1646653392799-0"
4) 1) 1) "test-consumer-name"
      2) "3"
```



### 参考

【Redis核心技术与实战】https://time.geekbang.org/column/intro/100056701    
【Redis设计与实现】https://book.douban.com/subject/25900156/  
【Redis Streams 介绍】http://www.redis.cn/topics/streams-intro.html      
【Centos7.6安装redis-6.0.8版本】https://blog.csdn.net/roc_wl/article/details/108662719    
