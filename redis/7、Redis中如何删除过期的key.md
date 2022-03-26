<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Redis 中 key 的过期删除策略](#redis-%E4%B8%AD-key-%E7%9A%84%E8%BF%87%E6%9C%9F%E5%88%A0%E9%99%A4%E7%AD%96%E7%95%A5)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [Redis 中 key 的过期删除策略](#redis-%E4%B8%AD-key-%E7%9A%84%E8%BF%87%E6%9C%9F%E5%88%A0%E9%99%A4%E7%AD%96%E7%95%A5-1)
    - [1、定时删除](#1%E5%AE%9A%E6%97%B6%E5%88%A0%E9%99%A4)
    - [2、惰性删除](#2%E6%83%B0%E6%80%A7%E5%88%A0%E9%99%A4)
    - [3、定期删除](#3%E5%AE%9A%E6%9C%9F%E5%88%A0%E9%99%A4)
    - [Redis 中过期删除策略](#redis-%E4%B8%AD%E8%BF%87%E6%9C%9F%E5%88%A0%E9%99%A4%E7%AD%96%E7%95%A5)
    - [从库是否会脏读主库创建的过期键](#%E4%BB%8E%E5%BA%93%E6%98%AF%E5%90%A6%E4%BC%9A%E8%84%8F%E8%AF%BB%E4%B8%BB%E5%BA%93%E5%88%9B%E5%BB%BA%E7%9A%84%E8%BF%87%E6%9C%9F%E9%94%AE)
  - [内存淘汰机制](#%E5%86%85%E5%AD%98%E6%B7%98%E6%B1%B0%E6%9C%BA%E5%88%B6)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Redis 中 key 的过期删除策略

### 前言

Redis 中的 key 设置一个过期时间，在过期时间到的时候，Redis 是如何清除这个 key 的呢？  

这来分析下 Redis 中的过期删除策略和内存淘汰机制   

### Redis 中 key 的过期删除策略

Redis 中提供了三种过期删除的策略  

#### 1、定时删除

在设置某个 key 的过期时间同时，我们创建一个定时器，让定时器在该过期时间到来时，立即执行对其进行删除的操作。  

优点：  

通过使用定时器，可以保证过期 key 可以被尽快的删除，并且释放过期 key 所占用的内存  

缺点：  

对 CPU 是不友好的，当过期键比较多的时候，删除过期 key 会占用相当一部分的 CPU 资源，对服务器的响应时间和吞吐量造成影响。  

#### 2、惰性删除

惰性删除，当一个键值对过期的时候，只有再次用到这个键值对的时候才去检查删除这个键值对，也就是如果用不着，这个键值对就会一直存在。  

优点：  

对 CPU 是友好的，只有在取出键值对的时候才会进行过期检查，这样就不会把 CPU 资源花费在其他无关简要的键值对的过期删除上。  

缺点：  

如果一些键值对永远不会被再次用到，那么将不会被删除，最终会造成内存泄漏，无用的垃圾数据占用了大量的资源，但是服务器却不能去删除。  

看下源码  

```
// https://github.com/redis/redis/blob/6.2/src/db.c#L1541
// 当访问到 key 的时候，会调用这个函数，因为有的 key 虽然已经被删除了，但是还可能存在于内存中

// key 仍然有效，函数返回值为0，否则，如果 key 过期，函数返回1。
int expireIfNeeded(redisDb *db, robj *key) {
    // 没有过期
    if (!keyIsExpired(db,key)) return 0;

    // 从库的过期是主库控制的，是不会进行删除操作的
    // 不过仍然尝试返回正确的信息给调用者，也就是说，如果我们认为 key 仍然有效，则返回0，如果我们认为此时 key 已经过期，则返回1。
    if (server.masterhost != NULL) return 1;
    ...
    /* Delete the key */
    // 删除 key 
    deleteExpiredKeyAndPropagate(db,key);
    return 1;
}
```

可以看到每次操作对应的 key 是会检查 key 是否过期，如果过期则会删除对应的 key 。   

如果过期键是主库创建的，那么从库进行检查是不会进行删除操作的,只是会根据 key 的过期时间返回过期或者未过期的状态。那么从库是不是会读取到过期的键值对信息呢，具体见下文分析。   

#### 3、定期删除

定期删除是对上面两种删除策略的一种整合和折中  

每个一段时间就对一些 key 进行采样检查，检查是否过期，如果过期就进行删除  

1、采样一定个数个数的key，可以进行配置，并将其中过期的key全部删除；  

2、如果过期 key 的占比超过`可接受的过期 key 的百分比`，则重复删除的过程，直到过期key的比例降至`可接受的过期 key 的百分比`以下。   

优点：  

定期删除，通过控制定期删除执行的时长和频率，可以减少删除操作对 CPU 的影响，同时也能较少因过期键带来的内存的浪费。  

缺点：  

执行的频率不太好控制  

频率过快对 CPU 不友好，如果过慢了就会对内存不太友好，过期的键值对不能及时的被删除掉  

同时如果一个键值对过期了，但是没有被删除，这时候业务再次获取到这个键值对，那么就会获取到被删除的数据了，这肯定是不合理的。   

看下源码实现  

```
// https://github.com/redis/redis/blob/6.2/src/server.c#L1853
// 这个函数处理我们需要在Redis数据库中增量执行的“后台”操作，例如活动键过期，调整大小，重哈希。
void databasesCron(void) {
    // 通过随机抽样来过期
    // 这里区分了主节点和从节点的处理
    if (server.active_expire_enabled) {
        if (iAmMaster()) {
            activeExpireCycle(ACTIVE_EXPIRE_CYCLE_SLOW);
        } else {
            expireSlaveKeys();
        }
    }
    ...
}

// https://github.com/redis/redis/blob/6.2/src/expire.c#L113
void activeExpireCycle(int type) {
    // 根据配置的超时工作调整运行参数。默认工作量为1，最大可配置工作量为10
    unsigned long
    effort = server.active_expire_effort-1, /* Rescale from 0 to 9. */
    // 采样的 key 的数量
    config_keys_per_loop = ACTIVE_EXPIRE_CYCLE_KEYS_PER_LOOP +
                           ACTIVE_EXPIRE_CYCLE_KEYS_PER_LOOP/4*effort,
    // 占比CPU时间，默认是25%，最大43%，如果是100%，那除了定时删除其他的工作都做不了了，所以要做限制
    config_cycle_slow_time_perc = ACTIVE_EXPIRE_CYCLE_SLOW_TIME_PERC +
                                  2*effort,
    // 可接受的过期 key 的百分比
    config_cycle_acceptable_stale = ACTIVE_EXPIRE_CYCLE_ACCEPTABLE_STALE-
                                    effort;
    ...
    //慢速定期删除的执行时长
    timelimit = config_cycle_slow_time_perc*1000000/server.hz/100;
    timelimit_exit = 0;
    ...
    // 在 key 过期时积累一些全局统计信息，以便了解逻辑上已经过期但仍存在于数据库中的 key 的数量
    long total_sampled = 0;
    long total_expired = 0;

    for (j = 0; j < dbs_per_call && timelimit_exit == 0; j++) {
        ...
        // 如果超过 config_cycle_acceptable_stale 的key过期了，则重复删除的过程，直到过期key的比例降至 config_cycle_acceptable_stale 以下。  
        // 存储在 config_cycle_acceptable_stale 中的百分比不是固定的，而是取决于Redis配置的“expire efforce”  
        do {
            /* If there is nothing to expire try next DB ASAP. */
            if ((num = dictSize(db->expires)) == 0) {
                db->avg_ttl = 0;
                break;
            }
            ...
            // 采样的 key 的数量 
            if (num > config_keys_per_loop)
                num = config_keys_per_loop;
            ...
            while (sampled < num && checked_buckets < max_buckets) {
                for (int table = 0; table < 2; table++) {
                    ...
                    while(de) {
                        /* Get the next entry now since this entry may get
                         * deleted. */
                        dictEntry *e = de;
                        de = de->next;

                        ttl = dictGetSignedIntegerVal(e)-now;
                        // 过期检查，并对过期键进行删除
                        if (activeExpireCycleTryExpire(db,e,now)) expired++;
                        ...
                    }
                }
                db->expires_cursor++;
            }
         ...
        // 判断过期 key 的占比是否大于 config_cycle_acceptable_stale，如果大于持续进行过期 key 的删除
        } while (sampled == 0 ||
                 (expired*100/sampled) > config_cycle_acceptable_stale);
    }
    ...
}

// 检查删除由从节点创建的有过期的时间的 key 
void expireSlaveKeys(void) {
    // 从主库同步的 key，过期时间由主库维护，主库同步 DEL 操作到从库。
    // 从库如果是 READ-WRITE 模式，就可以继续写入数据。从库自己写入的数据就需要自己来维护其过期操作。
    if (slaveKeysWithExpire == NULL ||
        dictSize(slaveKeysWithExpire) == 0) return;
     ...
}
```

惰性删除过程    

1、固定的时间执行一次定期删除；  

2、采样一定个数个数的key，可以进行配置，并将其中过期的key全部删除；  

3、如果过期 key 的占比超过`可接受的过期 key 的百分比`，则重复删除的过程，直到过期key的比例降至`可接受的过期 key 的百分比`以下；     

4、对于从库创建的过期 key 同样从库是不能进行删除的。   

#### Redis 中过期删除策略

上面讨论的三种策略，都有或多或少的问题。Redis 中实际采用的策略是惰性删除加定期删除的组合方式。  

组合方式的使用  

定期删除，获取 CPU 和 内存的使用平衡，针对过期的 KEY 可能得不到及时的删除，当 KEY 被再次获取的时候，通过惰性删除再做一次过期检查，来避免业务获取到过期内容。   

#### 从库是否会脏读主库创建的过期键  

从上面惰性删除和定期删除的源码阅读中，我们可以发现，从库对于主库的过期键是不能主动进行删除的。如果一个主库创建的过期键值对，已经过期了，主库在进行定期删除的时候，没有及时的删除掉，这时候从库请求了这个键值对，当执行惰性删除的时候，因为是主库创建的键值对，这时候是不能在从库中删除的，那么是不是就意味着从库会读取到已经过期的数据呢？  

答案肯定不是的  

[how-redis-replication-deals-with-expires-on-keys](https://redis.io/docs/manual/replication/#how-redis-replication-deals-with-expires-on-keys)  

> How Redis replication deals with expires on keys
> Redis expires allow keys to have a limited time to live. Such a feature depends on the ability of an instance to count the time, however Redis slaves correctly replicate keys with expires, even when such keys are altered using Lua scripts.
> To implement such a feature Redis cannot rely on the ability of the master and slave to have synchronized clocks, since this is a problem that cannot be solved and would result into race conditions and diverging data sets, so Redis uses three main techniques in order to make the replication of expired keys able to work:
> 1.Slaves don’t expire keys, instead they wait for masters to expire the keys. When a master expires a key (or evict it because of LRU), it synthesizes a DEL command which is transmitted to all the slaves.
> 2.However because of master-driven expire, sometimes slaves may still have in memory keys that are already logically expired, since the master was not able to provide the DEL command in time. In order to deal with that the slave uses its logical clock in order to report that a key does not exist only for read operations that don’t violate the consistency of the data set (as new commands from the master will arrive). In this way slaves avoid to report logically expired keys are still existing. In practical terms, an HTML fragments cache that uses slaves to scale will avoid returning items that are already older than the desired time to live.
> 3.During Lua scripts executions no keys expires are performed. As a Lua script runs, conceptually the time in the master is frozen, so that a given key will either exist or not for all the time the script runs. This prevents keys to expire in the middle of a script, and is needed in order to send the same script to the slave in a way that is guaranteed to have the same effects in the data set.
> Once a slave is promoted to a master it will start to expire keys independently, and will not require any help from its old master.

上面是官方文档中针对这一问题的描述  

大概意思就是从节点不会主动删除过期键，从节点会等待主节点触发键过期。当主节点触发键过期时，主节点会同步一个del命令给所有的从节点。  

因为是主节点驱动删除的，所以从节点会获取到已经过期的键值对。从节点需要根据自己本地的逻辑时钟来判断减值是否过期，从而实现数据集合的一致性读操作。     


### 内存淘汰机制

### 参考

【Redis核心技术与实战】https://time.geekbang.org/column/intro/100056701    
【Redis设计与实现】https://book.douban.com/subject/25900156/   