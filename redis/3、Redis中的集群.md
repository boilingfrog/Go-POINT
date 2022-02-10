<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Redis 的高可用集群](#redis-%E7%9A%84%E9%AB%98%E5%8F%AF%E7%94%A8%E9%9B%86%E7%BE%A4)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [几种常用的集群方案](#%E5%87%A0%E7%A7%8D%E5%B8%B8%E7%94%A8%E7%9A%84%E9%9B%86%E7%BE%A4%E6%96%B9%E6%A1%88)
  - [主从集群模式](#%E4%B8%BB%E4%BB%8E%E9%9B%86%E7%BE%A4%E6%A8%A1%E5%BC%8F)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Redis 的高可用集群

### 前言

这里来了解一下，Redis 中常见的集群方案  

### 几种常用的集群方案

- 主从集群模式  

- 哨兵机制  

- 切片集群(分片集群)  

### 主从集群模式 

主从集群，主从库之间才用的是读写分离  

- 主库：所有的写操作都在读库发生，然后主库同步数据到从库，同时也可以进行读操作；    

- 从库：只负责读操作；

<img src="/img/redis/redis-read-write.png"  alt="redis" align="center" />

### 参考

【Redis核心技术与实战】https://time.geekbang.org/column/intro/100056701    
【Redis设计与实现】https://book.douban.com/subject/25900156/  