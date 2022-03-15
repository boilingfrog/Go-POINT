<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [set 和 sorted set](#set-%E5%92%8C-sorted-set)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [set](#set)
    - [常见命令](#%E5%B8%B8%E8%A7%81%E5%91%BD%E4%BB%A4)
    - [set 的使用场景](#set-%E7%9A%84%E4%BD%BF%E7%94%A8%E5%9C%BA%E6%99%AF)
    - [看下源码实现](#%E7%9C%8B%E4%B8%8B%E6%BA%90%E7%A0%81%E5%AE%9E%E7%8E%B0)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## set 和 sorted set

### 前言

前面在几个文章聊到了 `list，string，hash` 等结构的实现，这次来聊一下 set 和 `sorted set` 的细节。  

### set

Redis 的 Set 是 String 类型的无序集合，集合成员是唯一的。  

底层实现主要用到了两种数据结构 hashtable 和 inset(整数集合)。    

集合中最大的成员数为2的32次方-1 (4294967295, 每个集合可存储40多亿个成员)。   

#### 常见命令

来看下几个常用的命令  

```
# 向集合添加一个或多个成员
SADD key member1 [member2]

# 获取集合的成员数
SCARD key

# 返回第一个集合与其他集合之间的差异。
SDIFF key1 [key2]

# 返回给定所有集合的差集并存储在 destination 中
SDIFFSTORE destination key1 [key2]

# 返回给定所有集合的交集
SINTER key1 [key2]

# 返回给定所有集合的交集并存储在 destination 中
SINTERSTORE destination key1 [key2]

# 判断 member 元素是否是集合 key 的成员
SISMEMBER key member

# 返回集合中的所有成员
SMEMBERS key

# 将 member 元素从 source 集合移动到 destination 集合
SMOVE source destination member

# 移除并返回集合中的一个随机元素
SPOP key

# 返回集合中一个或多个随机数
SRANDMEMBER key [count]

# 移除集合中一个或多个成员
SREM key member1 [member2]

# 返回所有给定集合的并集
SUNION key1 [key2]

# 所有给定集合的并集存储在 destination 集合中
SUNIONSTORE destination key1 [key2]

# 迭代集合中的元素
SSCAN key cursor [MATCH pattern] [COUNT count]
```

来个栗子  

```
127.0.0.1:6379>  SADD set-test xiaoming
(integer) 1
127.0.0.1:6379>  SADD set-test xiaoming
(integer) 0
127.0.0.1:6379>  SADD set-test xiaoming1
(integer) 1
127.0.0.1:6379>  SADD set-test xiaoming2


127.0.0.1:6379> SMEMBERS set-test
1) "xiaoming2"
2) "xiaoming"
3) "xiaoming1"
```

上面重复值的插入，只有第一次可以插入成功  

#### set 的使用场景

比较适用于聚合分类  

1、标签：比如我们博客网站常常使用到的兴趣标签，把一个个有着相同爱好，关注类似内容的用户利用一个标签把他们进行归并。  

2、共同好友功能，共同喜好，或者可以引申到二度好友之类的扩展应用。  

3、统计网站的独立IP。利用set集合当中元素不唯一性，可以快速实时统计访问网站的独立IP。  

不过对于 set 中的命令要合理的应用，不然很容易造成慢查询  

1、使用高效的命令，比如说，如果你需要返回一个 SET 中的所有成员时，不要使用 SMEMBERS 命令，而是要使用 SSCAN 多次迭代返回，避免一次返回大量数据，造成线程阻塞。  

2、当你需要执行排序、交集、并集操作时，可以在客户端完成，而不要用`SORT、SUNION、SINTER`这些命令，以免拖慢 Redis 实例。  

#### 看下源码实现  



### 参考

【Redis核心技术与实战】https://time.geekbang.org/column/intro/100056701    
【Redis设计与实现】https://book.douban.com/subject/25900156/    
【redis 集合（set）类型的使用和应用场景】https://www.oraclejsq.com/redisjc/040101720.html    

