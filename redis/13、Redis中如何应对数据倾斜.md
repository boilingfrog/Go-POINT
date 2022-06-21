<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Redis 中如何应对数据倾斜](#redis-%E4%B8%AD%E5%A6%82%E4%BD%95%E5%BA%94%E5%AF%B9%E6%95%B0%E6%8D%AE%E5%80%BE%E6%96%9C)
  - [什么是数据倾斜](#%E4%BB%80%E4%B9%88%E6%98%AF%E6%95%B0%E6%8D%AE%E5%80%BE%E6%96%9C)
  - [数据倾斜产生的原因](#%E6%95%B0%E6%8D%AE%E5%80%BE%E6%96%9C%E4%BA%A7%E7%94%9F%E7%9A%84%E5%8E%9F%E5%9B%A0)
    - [bigkey导致倾斜](#bigkey%E5%AF%BC%E8%87%B4%E5%80%BE%E6%96%9C)
  - [如何应对数据倾斜](#%E5%A6%82%E4%BD%95%E5%BA%94%E5%AF%B9%E6%95%B0%E6%8D%AE%E5%80%BE%E6%96%9C)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Redis 中如何应对数据倾斜

### 什么是数据倾斜

如果 Redis 中的部署，采用的是切片集群，数据是会按照一定的规则分散到不同的实例中保存，比如，使用 `Redis Cluster` 或 `Codis`。  

数据倾斜会有下面两种情况：  

1、**数据量倾斜**：在某些情况下，实例上的数据分布不均衡，某个实例上的数据特别多。  

2、**数据访问倾斜**：虽然每个集群实例上的数据量相差不大，但是某个实例上的数据是热点数据，被访问得非常频繁。  

发生了数据倾斜会造成，那些数据量大的和访问高的实例节点，系统的负载升高，响应速度变慢。严重的情况造成内存资源耗尽，引起系统崩溃。   

### 数据倾斜产生的原因

数据倾斜会存在，数据量倾斜和数据访问倾斜这两种情况，这里来分析下具体产生的原因  

#### bigkey导致倾斜

什么是 `Big Key`：我们将含有较大数据或含有大量成员、列表数的 Key 称之为大Key。

- 一个STRING类型的Key，它的值为5MB（数据过大）

- 一个LIST类型的Key，它的列表数量为20000个（列表数量过多）

- 一个ZSET类型的Key，它的成员数量为10000个（成员数量过多）

- 一个HASH格式的Key，它的成员数量虽然只有1000个但这些成员的value总大小为100MB（成员体积过大）

`Big Key` 存在问题  

- 内存空间不均匀：如果采用切片集群的部署方案，容易造成某些实例节点的内存分配不均匀；

- 造成网络拥塞：读取 bigkey 意味着需要消耗更多的网络流量，可能会对 Redis 服务器造成影响；

- 过期删除：big key 不单读写慢，删除也慢，删除过期 big key 也比较耗时；

- 迁移困难：由于数据庞大，备份和还原也容易造成阻塞，操作失败；

如何避免 

对于`Big Key`可以从以下两个方面进行处理

合理优化数据结构  

1、对较大的数据进行压缩处理；

2、拆分集合：将大的集合拆分成小集合（如以时间进行分片）或者单个的数据。

- 选择其他的技术来存储 `big key`；  

- 使用其他的存储形式，考虑使用 cdn 或者文档性数据库 MongoDB。  

### 如何应对数据倾斜

### 参考

【Redis核心技术与实战】https://time.geekbang.org/column/intro/100056701    
【Redis设计与实现】https://book.douban.com/subject/25900156/   
【Redis 的学习笔记】https://github.com/boilingfrog/Go-POINT/tree/master/redis    
【数据库事务】https://baike.baidu.com/item/%E6%95%B0%E6%8D%AE%E5%BA%93%E4%BA%8B%E5%8A%A1/9744607  
【transactions】https://redis.io/docs/manual/transactions/  
【Redis中的事务分析】https://boilingfrog.github.io/2022/06/19/Redis%E4%B8%AD%E7%9A%84%E4%BA%8B%E5%8A%A1%E5%88%86%E6%9E%90/  

