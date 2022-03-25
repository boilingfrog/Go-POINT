<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Redis 中 key 的过期删除策略](#redis-%E4%B8%AD-key-%E7%9A%84%E8%BF%87%E6%9C%9F%E5%88%A0%E9%99%A4%E7%AD%96%E7%95%A5)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [Redis 中 key 的过期删除策略](#redis-%E4%B8%AD-key-%E7%9A%84%E8%BF%87%E6%9C%9F%E5%88%A0%E9%99%A4%E7%AD%96%E7%95%A5-1)
    - [1、定时删除](#1%E5%AE%9A%E6%97%B6%E5%88%A0%E9%99%A4)
    - [2、惰性删除](#2%E6%83%B0%E6%80%A7%E5%88%A0%E9%99%A4)
    - [3、定期删除](#3%E5%AE%9A%E6%9C%9F%E5%88%A0%E9%99%A4)
    - [Redis 中过期删除策略](#redis-%E4%B8%AD%E8%BF%87%E6%9C%9F%E5%88%A0%E9%99%A4%E7%AD%96%E7%95%A5)
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

    // 如果是从库返回1
    // 从库的过期是主库控制的
    if (server.masterhost != NULL) return 1;
    ...
    /* Delete the key */
    // 删除 key 
    deleteExpiredKeyAndPropagate(db,key);
    return 1;
}
```

可以看到每次操作对应的 key 是会检查 key 是否过期，如果过期则会删除对应的 key 。   

#### 3、定期删除

定期删除是对上面两种删除策略的一种整合和折中  

每个一段时间就对一些 key 进行采样检查，检查是否过期，如果过期就进行删除  

1、采样ACTIVE_EXPIRE_CYCLE_LOOKUPS_PER_LOOP个数的key，并将其中过期的key全部删除；  

2、如果超过25%的key过期了，则重复删除的过程，直到过期key的比例降至25%以下。   

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
    config_keys_per_loop = ACTIVE_EXPIRE_CYCLE_KEYS_PER_LOOP +
                           ACTIVE_EXPIRE_CYCLE_KEYS_PER_LOOP/4*effort,
    config_cycle_fast_duration = ACTIVE_EXPIRE_CYCLE_FAST_DURATION +
                                 ACTIVE_EXPIRE_CYCLE_FAST_DURATION/4*effort,
    config_cycle_slow_time_perc = ACTIVE_EXPIRE_CYCLE_SLOW_TIME_PERC +
                                  2*effort,
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
        do {
            /* If there is nothing to expire try next DB ASAP. */
            if ((num = dictSize(db->expires)) == 0) {
                db->avg_ttl = 0;
                break;
            }
            ...
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
        // 判断过期 key 的占比是否满足 config_cycle_acceptable_stale
        } while (sampled == 0 ||
                 (expired*100/sampled) > config_cycle_acceptable_stale);
    }
    ...
}
```

#### Redis 中过期删除策略

上面讨论的三种策略，都有或多或少的问题。Redis 中实际采用的策略是惰性删除加定期删除的组合方式。  

组合方式的使用  

定期删除，获取 CPU 和 内存的使用平衡，针对过期的 KEY 可能得不到及时的删除，当 KEY 被再次获取的时候，通过惰性删除再做一次过期检查，来避免业务获取到过期内容。   



### 内存淘汰机制



 

### 参考

【Redis核心技术与实战】https://time.geekbang.org/column/intro/100056701    
【Redis设计与实现】https://book.douban.com/subject/25900156/   