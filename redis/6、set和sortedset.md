<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [set 和 sorted set](#set-%E5%92%8C-sorted-set)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [set](#set)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## set 和 sorted set

### 前言

前面在几个文章聊到了 `list，string，hash` 等结构的实现，这次来聊一下 set 和 `sorted set` 的细节。  

### set

Redis 的 Set 是 String 类型的无序集合，集合成员是唯一的。  

底层实现主要用到了两种数据结构 hashtable 和 inset(整数集合)。    

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




### 参考

【Redis核心技术与实战】https://time.geekbang.org/column/intro/100056701    
【Redis设计与实现】https://book.douban.com/subject/25900156/    

