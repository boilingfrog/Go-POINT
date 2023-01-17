<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [一条 SQL 的执行过程](#%E4%B8%80%E6%9D%A1-sql-%E7%9A%84%E6%89%A7%E8%A1%8C%E8%BF%87%E7%A8%8B)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [查询](#%E6%9F%A5%E8%AF%A2)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 一条 SQL 的执行过程  

### 前言

在开始学习 MySQL 中知识点的时候，首先来看下 SQL 在 MySQL 中的执行过程。    

### 查询

查询语句是我们经常用到的，那么一个简单的查询 sql，在 MySQL 中的执行过程是怎么样的呢？    

```
SELECT * FROM user WHERE id =1
```

栗如上面的这个简单的查询语句，来看下具体的查询逻辑。    

<img src="/img/mysql/mysql-query.png"  alt="mysql" />

**连接器**   

大多数基于网络的客户端/服务器的工具或者服务都有类似的架构。  

这里主要的工作就是管理和客户端的连接，同时进行连接的权限认证。   

为了避免线程被频繁的创建和销毁，影响性能，`MySQL5.5` 版本引入了线程池，会缓存创建的线程，不需要为每一个新建的连接，创建或销毁线程。可以使用线程池中少量的线程服务大量的连接。    

**查询缓存**

MySQL 查询缓存，为了提高相同 Query 语句的响应速度，会缓存特定 Query 的整个结果集信息，当后面有相同的查询，直接查询缓存，响应给客户端。      

不过当一个表有更新的时候，和这个表有关的查询缓存都会被删除，造成查询缓存的失效。所以一般不建议使用查询缓存。   

在 `MySQL 5.6` 开始，就已经默认禁用查询缓存了。在 `MySQL 8.0`，就已经删除查询缓存功能了。  

### 参考

【高性能MySQL(第3版)】https://book.douban.com/subject/23008813/    
【MySQL 实战 45 讲】https://time.geekbang.org/column/100020801  