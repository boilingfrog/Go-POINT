<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Redis 基础知识点](#redis-%E5%9F%BA%E7%A1%80%E7%9F%A5%E8%AF%86%E7%82%B9)
  - [为什么 Redis 比较快](#%E4%B8%BA%E4%BB%80%E4%B9%88-redis-%E6%AF%94%E8%BE%83%E5%BF%AB)
  - [为什么单线程还能很快](#%E4%B8%BA%E4%BB%80%E4%B9%88%E5%8D%95%E7%BA%BF%E7%A8%8B%E8%BF%98%E8%83%BD%E5%BE%88%E5%BF%AB)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Redis 基础知识点

### 为什么 Redis 比较快 

Redis 中的查询速度为什么那么快呢？  

1、因为它是内存数据库；  

2、归功于它的数据结构；  

3、Redis 中是单线程。   

### 为什么单线程还能很快

Redis 是单线程，主要是指 Redis 的网络IO和键值对读写是由一个线程来完成的，这也是 Redis 对外提供键值存储服务的主要流程。   

 