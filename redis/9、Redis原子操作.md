<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Redis 如何应对并发访问](#redis-%E5%A6%82%E4%BD%95%E5%BA%94%E5%AF%B9%E5%B9%B6%E5%8F%91%E8%AE%BF%E9%97%AE)
  - [Redis 中处理并发的方案](#redis-%E4%B8%AD%E5%A4%84%E7%90%86%E5%B9%B6%E5%8F%91%E7%9A%84%E6%96%B9%E6%A1%88)
  - [原子性](#%E5%8E%9F%E5%AD%90%E6%80%A7)
    - [原子性的单命令](#%E5%8E%9F%E5%AD%90%E6%80%A7%E7%9A%84%E5%8D%95%E5%91%BD%E4%BB%A4)
    - [Redis 的IO模型分析](#redis-%E7%9A%84io%E6%A8%A1%E5%9E%8B%E5%88%86%E6%9E%90)
      - [thread-based architecture（基于线程的架构）](#thread-based-architecture%E5%9F%BA%E4%BA%8E%E7%BA%BF%E7%A8%8B%E7%9A%84%E6%9E%B6%E6%9E%84)
      - [event-driven architecture（事件驱动模型）](#event-driven-architecture%E4%BA%8B%E4%BB%B6%E9%A9%B1%E5%8A%A8%E6%A8%A1%E5%9E%8B)
        - [Reactor 模式](#reactor-%E6%A8%A1%E5%BC%8F)
        - [Proactor 模式](#proactor-%E6%A8%A1%E5%BC%8F)
    - [为什么 Redis 选择单线程](#%E4%B8%BA%E4%BB%80%E4%B9%88-redis-%E9%80%89%E6%8B%A9%E5%8D%95%E7%BA%BF%E7%A8%8B)
    - [使用 LUA 脚本](#%E4%BD%BF%E7%94%A8-lua-%E8%84%9A%E6%9C%AC)
  - [分布式锁](#%E5%88%86%E5%B8%83%E5%BC%8F%E9%94%81)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Redis 如何应对并发访问

### Redis 中处理并发的方案

业务中有时候我们会用 Redis 处理一些高并发的业务场景，例如，秒杀业务，对于库存的操作。。。   

先来分析下，并发场景下会发生什么问题   

并发问题主要发生在数据的修改上，对于客户端修改数据，一般分成下面两个步骤：  

1、客户端先把数据读取到本地，在本地进行修改；  

2、客户端修改完数据后，再写回Redis。  

我们把这个流程叫做`读取-修改-写回`操作（`Read-Modify-Write`，简称为 RMW 操作）。如果客户端并发进行 RMW 操作的时候，就需要保证 `读取-修改-写回`是一个原子操作，进行命令操作的时候，其他客户端不能对当前的数据进行操作。  

错误的栗子：  

统计一个页面的访问次数，每次刷新页面访问次数+1，这里使用 Redis 来记录访问次数。  

如果每次的`读取-修改-写回`操作不是一个原子操作，那么就可能存在下图的问题，客户端2在客户端1操作的中途，也获取 Redis 的值，也对值进行+1，操作，这样就导致最终数据的错误。  

<img src="/img/redis/redis-rmw.png"  alt="redis" align="center" />

对于上面的这种情况，一般会有两种方式解决：  

1、使用 Redis 实现一把分布式锁，通过锁来保护每次只有一个线程来操作临界资源；  

2、实现操作命令的原子性。  

- 栗如，对于上面的错误栗子，如果`读取-修改-写回`是一个原子性的命令，那么这个命令在操作过程中就不有别的线程同时读取操作数据，这样就能避免上面栗子出现的问题。  

下面从原子性和锁两个方面，具体分析下，对并发访问问题的处理   

### 原子性

为了实现并发控制要求的临界区代码互斥执行，Redis的原子操作采用了两种方法：  

1、借助于 Redis 中的原子性的单命令；    

2、把多个操作写到一个Lua脚本中，以原子性方式执行单个Lua脚本。  

#### 原子性的单命令

Redis 中的单个命令的执行都是原子性的，这里来具体的探讨下   

#### Redis 的IO模型分析

在探讨 Redis 原子性的时候，先来探讨下 Redis 中的 IO 模型  

在 web 服务中，处理 web 请求通常有两种体系结构，分别为：`thread-based architecture`（基于线程的架构）、`event-driven architecture`（事件驱动模型）  

##### thread-based architecture（基于线程的架构）

thread-based architecture（基于线程的架构）：这种比较容易理解，就是多线程并发模式，服务端在处理请求的时候，一个请求分配一个独立的线程来处理。   

因为每个请求分配一个独立的线程，所以单个线程的阻塞不会影响到其他的线程，能够提高程序的响应速度。  

不足的是，连接和线程之间始终保持一对一的关系，如果是一直处于 Keep-Alive 状态的长连接将会导致大量工作线程在空闲状态下等待，例如，文件系统访问，网络等。此外，成百上千的连接还可能会导致并发线程浪费大量内存的堆栈空间。  

##### event-driven architecture（事件驱动模型）  

事件驱动的体系结构由事件生产者和事件消费者组，是一种松耦合、分布式的驱动架构，生产者收集到某应用产生的事件后实时对事件采取必要的处理后路由至下游系统，无需等待系统响应，下游的事件消费者组收到是事件消息，异步的处理。  

事件驱动架构具有以下优势：  

- 降低耦合；  

降低事件生产者和订阅者的耦合性。事件生产者只需关注事件的发生，无需关注事件如何处理以及被分发给哪些订阅者。任何一个环节出现故障，不会影响其他业务正常运行。

- 异步执行；  

事件驱动架构适用于异步场景，即便是需求高峰期，收集各种来源的事件后保留在事件总线中，然后逐步分发传递事件，不会造成系统拥塞或资源过剩的情况。

- 可扩展性；
- 
事件驱动架构中路由和过滤能力支持划分服务，便于扩展和路由分发。  

Reactor 模式和 Proactor 模式都是 `event-driven architecture`（事件驱动模型）的实现方式，这里具体分析下  

###### Reactor 模式

Reactor 模式，是指通过一个或多个输入同时传递给服务处理器的服务请求的事件驱动处理模式。  

在处理⽹络 IO 的连接事件、读事件、写事件。Reactor 中引入了三类角色  

- reactor：监听和分配事件，连接事件交给 acceptor 处理，读写事件交给 handler 处理；

- acceptor：接收连接请求，接收连接后，会创建 handler ，处理网络连接上对后续读写事件的处理；  

- handler：处理读写事件。  

<img src="/img/redis/redis-reactor.png"  alt="redis" />


###### Proactor 模式

#### 为什么 Redis 选择单线程


比如对于上面的`读取-修改-写回`操作可以使用 Redis 中的原子计数器, INCRBY（自增）、DECRBR（自减）、INCR（加1） 和 DECR（减1） 等命令。  

这些命令可以直接帮助我们处理并发控制   

```
127.0.0.1:6379> incr test-1
(integer) 1
127.0.0.1:6379> incr test-1
(integer) 2
127.0.0.1:6379> incr test-1
(integer) 3
```

分析下源码，看看这个命令是如何实现的    

```
// https://github.com/redis/redis/blob/6.2/src/t_string.c#L617

void incrCommand(client *c) {
    incrDecrCommand(c,1);
}

void decrCommand(client *c) {
    incrDecrCommand(c,-1);
}

void incrbyCommand(client *c) {
    long long incr;

    if (getLongLongFromObjectOrReply(c, c->argv[2], &incr, NULL) != C_OK) return;
    incrDecrCommand(c,incr);
}

void decrbyCommand(client *c) {
    long long incr;

    if (getLongLongFromObjectOrReply(c, c->argv[2], &incr, NULL) != C_OK) return;
    incrDecrCommand(c,-incr);
}
```

可以看到 INCRBY（自增）、DECRBR（自减）、INCR（加1） 和 DECR（减1）这几个命令最终都是调用的 incrDecrCommand  

```
// https://github.com/redis/redis/blob/6.2/src/t_string.c#L579  
void incrDecrCommand(client *c, long long incr) {
    long long value, oldvalue;
    robj *o, *new;

    // 查找有没有对应的键值
    o = lookupKeyWrite(c->db,c->argv[1]);
    // 判断类型，如果value对象不是字符串类型，直接返回
    if (checkType(c,o,OBJ_STRING)) return;

    // 将字符串类型的value转换为longlong类型保存在value中
    if (getLongLongFromObjectOrReply(c,o,&value,NULL) != C_OK) return;

    // 备份旧的value
    oldvalue = value;

    // 判断 incr 的值是否超过longlong类型所能表示的范围
    // 长度的范围，十进制 64 位有符号整数
    if ((incr < 0 && oldvalue < 0 && incr < (LLONG_MIN-oldvalue)) ||
        (incr > 0 && oldvalue > 0 && incr > (LLONG_MAX-oldvalue))) {
        addReplyError(c,"increment or decrement would overflow");
        return;
    }
    // 计算新的 value值
    value += incr;

    if (o && o->refcount == 1 && o->encoding == OBJ_ENCODING_INT &&
        (value < 0 || value >= OBJ_SHARED_INTEGERS) &&
        value >= LONG_MIN && value <= LONG_MAX)
    {
        new = o;
        o->ptr = (void*)((long)value);
    } else {
        new = createStringObjectFromLongLongForValue(value);
        // 如果之前的 value 对象存在
        if (o) {
            // 重写为 new 的值  
            dbOverwrite(c->db,c->argv[1],new);
        } else {
            // 如果之前没有对应的 value,新设置 value 的值
            dbAdd(c->db,c->argv[1],new);
        }
    }
    // 进行通知
    signalModifiedKey(c,c->db,c->argv[1]);
    notifyKeyspaceEvent(NOTIFY_STRING,"incrby",c->argv[1],c->db->id);
    server.dirty++;
    addReply(c,shared.colon);
    addReply(c,new);
    addReply(c,shared.crlf);
}
```

#### 使用 LUA 脚本


### 分布式锁

### 参考

【Redis核心技术与实战】https://time.geekbang.org/column/intro/100056701    
【Redis设计与实现】https://book.douban.com/subject/25900156/   
【字符串命令的实现】https://mcgrady-forever.github.io/2018/02/10/redis-analysis-t-string/     
【Redis 多线程网络模型全面揭秘】https://segmentfault.com/a/1190000039223696     
【高性能IO模型分析-Reactor模式和Proactor模式】https://zhuanlan.zhihu.com/p/95662364  
【什么是事件驱动架构？】https://www.redhat.com/zh/topics/integration/what-is-event-driven-architecture  
【事件驱动架构】https://help.aliyun.com/document_detail/207135.html   


