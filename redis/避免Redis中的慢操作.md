<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Redis 保持如何保持高性能](#redis-%E4%BF%9D%E6%8C%81%E5%A6%82%E4%BD%95%E4%BF%9D%E6%8C%81%E9%AB%98%E6%80%A7%E8%83%BD)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [如何应对 Redis 变慢](#%E5%A6%82%E4%BD%95%E5%BA%94%E5%AF%B9-redis-%E5%8F%98%E6%85%A2)
    - [Redis自身操作特性的影响](#redis%E8%87%AA%E8%BA%AB%E6%93%8D%E4%BD%9C%E7%89%B9%E6%80%A7%E7%9A%84%E5%BD%B1%E5%93%8D)
  - [在Redis中，还有哪些命令可以代替KEYS命令，实现对键值对的key的模糊查询呢？这些命令的复杂度会导致Redis变慢吗？](#%E5%9C%A8redis%E4%B8%AD%E8%BF%98%E6%9C%89%E5%93%AA%E4%BA%9B%E5%91%BD%E4%BB%A4%E5%8F%AF%E4%BB%A5%E4%BB%A3%E6%9B%BFkeys%E5%91%BD%E4%BB%A4%E5%AE%9E%E7%8E%B0%E5%AF%B9%E9%94%AE%E5%80%BC%E5%AF%B9%E7%9A%84key%E7%9A%84%E6%A8%A1%E7%B3%8A%E6%9F%A5%E8%AF%A2%E5%91%A2%E8%BF%99%E4%BA%9B%E5%91%BD%E4%BB%A4%E7%9A%84%E5%A4%8D%E6%9D%82%E5%BA%A6%E4%BC%9A%E5%AF%BC%E8%87%B4redis%E5%8F%98%E6%85%A2%E5%90%97)
  - [导致 Redis 变慢的场景](#%E5%AF%BC%E8%87%B4-redis-%E5%8F%98%E6%85%A2%E7%9A%84%E5%9C%BA%E6%99%AF)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Redis 保持如何保持高性能

### 前言

总结下 Redis 中可能导致满请求的操作  

### 如何应对 Redis 变慢

#### Redis自身操作特性的影响

- 1、慢查询命令   

用其他高效命令代替。比如说，如果你需要返回一个SET中的所有成员时，不要使用SMEMBERS命令，而是要使用SSCAN多次迭代返回，避免一次返回大量数据，造成线程阻塞。 
 
当你需要执行排序、交集、并集操作时，可以在客户端完成，而不要用SORT、SUNION、SINTER这些命令，以免拖慢Redis实例。  

因为KEYS命令需要遍历存储的键值对，所以操作延时高，一般线上环境不建议用   

- 2.过期key操作

Redis 中清除过期 key 也是需要浪费一定的性能的，如果同一时刻有很多相同过期时间的 key 过期了，会触发过期删除策略，然后同一时刻需要处理的删除的 key 就很多。  

所以我们使用过期 key 的时候，可以考虑给过期时间设置一个随机数，而不是相同的过期时间。  

### 在Redis中，还有哪些命令可以代替KEYS命令，实现对键值对的key的模糊查询呢？这些命令的复杂度会导致Redis变慢吗？  

Redis提供的SCAN命令，以及针对集合类型数据提供的SSCAN、HSCAN等，可以根据执行时设定的数量参数，返回指定数量的数据，这就可以避免像KEYS命令一样同时返回所有匹配的数据，不会导致Redis变慢。以HSCAN为例，我们可以执行下面的命令，从user这个Hash集合中返回key前缀以103开头的100个键值对。  

``
HSCAN user 0  match "103*" 100
``

### 导致 Redis 变慢的场景

1、使用复杂度过高的命令或一次查询全量数据；  

2、操作bigkey；  

3、大量key集中过期；  

4、内存达到maxmemory；  

5、客户端使用短连接和Redis相连；  

6、当Redis实例的数据量大时，无论是生成RDB，还是AOF重写，都会导致fork耗时严重；  

7、AOF的写回策略为always，导致每个操作都要同步刷回磁盘；  
 
8、Redis实例运行机器的内存不足，导致swap发生，Redis需要到swap分区读取数据；  

9、进程绑定CPU不合理；  

10、Redis实例运行机器上开启了透明内存大页机制；  

11、网卡压力过大。


### 参考

【Redis核心技术与实战】https://time.geekbang.org/column/intro/100056701    
【Redis设计与实现】https://book.douban.com/subject/25900156/  
