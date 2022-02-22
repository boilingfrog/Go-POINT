<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Redis 在日常使用中的优化](#redis-%E5%9C%A8%E6%97%A5%E5%B8%B8%E4%BD%BF%E7%94%A8%E4%B8%AD%E7%9A%84%E4%BC%98%E5%8C%96)
  - [使用 String 类型内存开销大](#%E4%BD%BF%E7%94%A8-string-%E7%B1%BB%E5%9E%8B%E5%86%85%E5%AD%98%E5%BC%80%E9%94%80%E5%A4%A7)
    - [1、简单动态字符串](#1%E7%AE%80%E5%8D%95%E5%8A%A8%E6%80%81%E5%AD%97%E7%AC%A6%E4%B8%B2)
    - [2、RedisObject](#2redisobject)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Redis 在日常使用中的优化

### 使用 String 类型内存开销大

如果我们有大量的数据需要来保存，在选型数据类型我们就需要知道 String 的内存开销是很大的  

这里我们来分析下使用一个 String 类型需要用到的内存    

#### 1、简单动态字符串

Redis 中的 String,使用的是简单动态字符串（Simple Dynamic Strings，SDS）。  

来看下数据结构  

```
struct sdshdr {
    // 记录 buf 数组中已使用字节的数量
    // 等于 SDS 保存字符串的长度，不包含'\0'
    long len;
    
    // 记录buf数组中未使用字节的数量
    long free;
    
    // 字节数组，用于保存字符串
    char buf[];
};
```

如果，使用 SDS 存储了一个字符串 hello,对应的 len 就是5，同时也申请了5个为未使用的空间，所以 free 就是5。对于 buf 来说，len 和 free 的内存占用都是额外开销。         

<img src="/img/redis/redis-sds.png"  alt="redis" align="center" />

#### 2、RedisObject

因为 Redis 中有很多数据类型，对于这些不同的数据结构，Redis 为了能够统一处理，所以引入了 RedisObject。  

```
typedef struct redisObject {
    unsigned type:4; // 类型
    unsigned encoding:4; // 编码
    unsigned lru:LRU_BITS; // 最近被访问的时间
    int refcount; // 引用次数
    void *ptr; // 指向具体底层数据的指针
} robj;
```

<img src="/img/redis/redis-object.svg"  alt="redis" align="center" />

一个 RedisObject 包含了8字节的元数据和一个8字节指针，指针指向实际的数据内存地址。   

不过需要注意的是这里 Redis 做了优化  

1、当保存的数据是 Long 类型整数时，RedisObjec t中的指针就直接赋值为整数数据了，就不用使用额外的指针了。  

2、如果保存的是字符串数据，并且字符串大小小于等于44字节时，RedisObject中的元数据、指针和SDS是一块连续的内存区域，这样就可以避免内存碎片。这种布局方式也被称为 embstr 编码方式。  

3、如果保存的是字符串数据，并且字符串大小大于44字节时，Redis 就不再把 SDS 和 RedisObject 布局在一起了，而是会给 SDS 分配独立的空间，并用指针指向 SDS 结构。这种布局方式被称为 raw 编码模式。    

<img src="/img/redis/redis-object-string.jpeg"  alt="redis" align="center" />


### 参考

【Redis核心技术与实战】https://time.geekbang.org/column/intro/100056701    
【Redis设计与实现】https://book.douban.com/subject/25900156/  
【redis 一组kv实际内存占用计算】https://kernelmaker.github.io/Redis-StringMem    