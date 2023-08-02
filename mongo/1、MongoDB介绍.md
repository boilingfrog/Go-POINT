<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [什么是 MongoDB](#%E4%BB%80%E4%B9%88%E6%98%AF-mongodb)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [MongoDB 对比关系型数据库 MySQL](#mongodb-%E5%AF%B9%E6%AF%94%E5%85%B3%E7%B3%BB%E5%9E%8B%E6%95%B0%E6%8D%AE%E5%BA%93-mysql)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 什么是 MongoDB

### 前言

MongoDB 是一个开源、高性能、无模式的以 JSON 为数据模型的文档数据库，是 NoSQL 数据库产品的一种。   

MongoDB 中记录的是一个文档，由字段和键值对组成的数据结构。MongoDB 文档类似于 JSON 对象，字段的值可以包括其它文档，数组和文档数组。   

<img src="/img/mongo/mongo-json.jpg"  alt="mongo" />     

使用文档的优点：      

1、文档（即对象）对应于许多编程语言中的内置数据类型；  

2、嵌入式文档和数组减少了对昂贵连接的需求；  

3、动态模式支持流畅的多态性。    

MongoDB 的主要特性  

**1、高性能**   

MongoDB 提供高性能的数据持久化。   

- 对嵌入式数据模型的支持能够减少数据库系统上的 I/0 操作；   

- 支持索引查询，嵌入式文档和数组的键也支持索引的创建。   

**2、丰富的查询语句**  

MongoDB 支持丰富的查询语句,除了简单的 CRUD 操作，还支持复杂的查询类似 数据聚合，文本搜索和地理空间搜索。  

**3、高可用**

MongoDB 中的副本集给 MongoDB 提供了冗余和数据高可用的特性。   

**4、水平扩展**

MongoDB 中支持分片保证了 MongoDB 

当数据量很大的时候有两种方式解决：  

1、垂直扩展，增加单机 CPU、内存等；  

2、水平扩展，将数据分散到不同机器上，分摊数据压力。  

MongoDB 除了支持垂直扩展，MongoDB 中的 `Sharded Cluster` 集群实现了水平扩展的能力。   

### MongoDB 对比关系型数据库 MySQL

|              |    MongoDB          |      MySQL                  |
| ------------ | -----------------   | --------------------------  |
| 数据模型      |    文档模型           |  关系模型                    |
| 存储方式      | 以类json的文档格式存储 |  不同的存储引擎有自己的存储方式  |
| 高可用       |    复制集             |  集群模式                    |
| 横向扩展能力  | 通过原生分片完善支持     |  数据分区或者应用侵入式        |
| 数据容量     | 没有理论上限            |  千万，亿                    |
| join操作     | MongoDB没有Join       |  MySQL支持join               |


### 参考

【MongoDB简介】https://docs.mongoing.com/mongo-introduction      