<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Redis 基础知识点](#redis-%E5%9F%BA%E7%A1%80%E7%9F%A5%E8%AF%86%E7%82%B9)
  - [为什么 Redis 比较快](#%E4%B8%BA%E4%BB%80%E4%B9%88-redis-%E6%AF%94%E8%BE%83%E5%BF%AB)
  - [为什么单线程还能很快](#%E4%B8%BA%E4%BB%80%E4%B9%88%E5%8D%95%E7%BA%BF%E7%A8%8B%E8%BF%98%E8%83%BD%E5%BE%88%E5%BF%AB)
  - [基于多路复用的高性能I/O模型](#%E5%9F%BA%E4%BA%8E%E5%A4%9A%E8%B7%AF%E5%A4%8D%E7%94%A8%E7%9A%84%E9%AB%98%E6%80%A7%E8%83%BDio%E6%A8%A1%E5%9E%8B)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Redis 基础知识点

### 为什么 Redis 比较快 

Redis 中的查询速度为什么那么快呢？  

1、因为它是内存数据库；  

2、归功于它的数据结构；  

3、Redis 中是单线程；     

4、Redis 中使用了多路复用。  

### 为什么单线程还能很快

Redis 是单线程，主要是指 Redis 的网络IO和键值对读写是由一个线程来完成的，这也是 Redis 对外提供键值存储服务的主要流程。   

多线程必然会面临对于共享资源的访问，这时候通常的做法就是加锁，虽然是多线程，这时候就会变成串行的访问。也就是多线程编程模式会面临的共享资源的并发访问控制问题。  

同时多线程也会引入同步原语来保护共享资源的并发访问，代码的可维护性和易读性将会下降。   

### 基于多路复用的高性能I/O模型

Linux 中的IO多路复用机制是指一个线程处理多个IO流，就是我们经常听到的`select/epoll`机制。简单来说，在Redis只运行单线程的情况下，该机制允许内核中，同时存在多个监听套接字和已连接套接字。内核会一直监听这些套接字上的连接请求或数据请求。一旦有请求到达，就会交给 Redis 线程处理，这就实现了一个 Redis 线程处理多个IO流的效果。  

当然请求到达 Redis 线程，`select/epoll`提供了基于事件的回调机制，即针对不同事件，调用响应的回调函数。这些回调事件会被放入到一个事件队列，Redis 单线程对该事件队列不断进行处理。这样一来，Redis无需一直轮询是否有请求实际发生，这就可以避免造成CPU资源浪费。同时，Redis 在对事件队列中的事件进行处理时，会调用相应的处理函数，这就实现了基于事件的回调。因为 Redis 一直在对事件队列进行处理，所以能及时响应客户端请求，提升 Redis 的响应性能。  


### 参考

【Redis核心技术与实战】https://time.geekbang.org/column/intro/100056701    