<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Redis 保持如何保持高性能](#redis-%E4%BF%9D%E6%8C%81%E5%A6%82%E4%BD%95%E4%BF%9D%E6%8C%81%E9%AB%98%E6%80%A7%E8%83%BD)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [如何应对 Redis 变慢](#%E5%A6%82%E4%BD%95%E5%BA%94%E5%AF%B9-redis-%E5%8F%98%E6%85%A2)
    - [Redis自身操作特性的影响](#redis%E8%87%AA%E8%BA%AB%E6%93%8D%E4%BD%9C%E7%89%B9%E6%80%A7%E7%9A%84%E5%BD%B1%E5%93%8D)
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





### 参考

【Redis核心技术与实战】https://time.geekbang.org/column/intro/100056701    
【Redis设计与实现】https://book.douban.com/subject/25900156/  
