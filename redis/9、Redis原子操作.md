<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Redis 如何应对并发访问](#redis-%E5%A6%82%E4%BD%95%E5%BA%94%E5%AF%B9%E5%B9%B6%E5%8F%91%E8%AE%BF%E9%97%AE)
  - [Redis 中处理并发的方案](#redis-%E4%B8%AD%E5%A4%84%E7%90%86%E5%B9%B6%E5%8F%91%E7%9A%84%E6%96%B9%E6%A1%88)
    - [原子性](#%E5%8E%9F%E5%AD%90%E6%80%A7)
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

#### 分布式锁

### 参考

【Redis核心技术与实战】https://time.geekbang.org/column/intro/100056701    
【Redis设计与实现】https://book.douban.com/subject/25900156/   
