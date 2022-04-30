<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Redis 如何应对并发访问](#redis-%E5%A6%82%E4%BD%95%E5%BA%94%E5%AF%B9%E5%B9%B6%E5%8F%91%E8%AE%BF%E9%97%AE)
  - [Redis 中处理并发的方案](#redis-%E4%B8%AD%E5%A4%84%E7%90%86%E5%B9%B6%E5%8F%91%E7%9A%84%E6%96%B9%E6%A1%88)
    - [原子性](#%E5%8E%9F%E5%AD%90%E6%80%A7)
      - [原子性的单命令](#%E5%8E%9F%E5%AD%90%E6%80%A7%E7%9A%84%E5%8D%95%E5%91%BD%E4%BB%A4)
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

#### 原子性

为了实现并发控制要求的临界区代码互斥执行，Redis的原子操作采用了两种方法：  

1、借助于 Redis 中的原子性的单命令；    

2、把多个操作写到一个Lua脚本中，以原子性方式执行单个Lua脚本。  

##### 原子性的单命令

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

##### 使用 LUA 脚本


#### 分布式锁

### 参考

【Redis核心技术与实战】https://time.geekbang.org/column/intro/100056701    
【Redis设计与实现】https://book.douban.com/subject/25900156/   
【字符串命令的实现】https://mcgrady-forever.github.io/2018/02/10/redis-analysis-t-string/     
