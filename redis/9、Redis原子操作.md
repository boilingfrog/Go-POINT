<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Redis 如何应对并发访问](#redis-%E5%A6%82%E4%BD%95%E5%BA%94%E5%AF%B9%E5%B9%B6%E5%8F%91%E8%AE%BF%E9%97%AE)
  - [Redis 中处理并发的方案](#redis-%E4%B8%AD%E5%A4%84%E7%90%86%E5%B9%B6%E5%8F%91%E7%9A%84%E6%96%B9%E6%A1%88)
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

<img src="/img/redis/redis-rmw.png"  alt="redis" align="center" />


处理高并发的时候，一般会有两种操作  

1、使用 Redis 实现一把分布式锁，通过锁来保护热点数据；  

- 栗如，我们现在有个秒杀的场景，并发量可能是3000，但是我们商品的库存数量是一定的，为了防止超卖，我们就需要在减库存的时候加上锁，当第一个请求过来的时候，先判断锁时候存在，不存在就加锁，然后去处理秒杀的业务，并且在处理完成的时候，释放锁，如果判断锁存在，就轮训等待锁被释放。  

2、使用原子性。  

- 栗如，对于上面秒杀减库存的操作，查询商品库存，染回  

### 参考

【Redis核心技术与实战】https://time.geekbang.org/column/intro/100056701    
【Redis设计与实现】https://book.douban.com/subject/25900156/   
