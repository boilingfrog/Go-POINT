<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [MySQL 中的锁理解](#mysql-%E4%B8%AD%E7%9A%84%E9%94%81%E7%90%86%E8%A7%A3)
  - [锁的类型](#%E9%94%81%E7%9A%84%E7%B1%BB%E5%9E%8B)
  - [全局锁](#%E5%85%A8%E5%B1%80%E9%94%81)
  - [行锁](#%E8%A1%8C%E9%94%81)
  - [表锁](#%E8%A1%A8%E9%94%81)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## MySQL 中的锁理解

### 锁的类型

MySQL 找那个根据加锁的范围，大致可以分成全局锁，表级锁和行级锁。   

### 全局锁

全局锁，就是对整个数据库加锁。      

加锁

```
flush tables with read lock
```

解锁

```
unlock tables
```

全局锁会让整个库处于只读状态，之后所有的更新操作都会被阻塞：   

- 数据更新语句（数据的增删改）；    

- 数据定义语句（包括建表、修改表结构等）和更新类事务的提交语句。

如果对主库加锁，那么执行期间就不能执行更新，业务基本上就停摆了；   

如果对从库加锁，那么执行期间，从库就不能执行主库同步过来的 binlog，会导致主从延迟。       

全局锁的典型使用场景是，做全库逻辑备份。也就是把整库每个表都select出来存成文本。

不过为什么要在备份的时候加锁，不加锁的话，备份系统备份的得到的库不是一个逻辑时间点，这个视图是逻辑不一致的。  

官方自带的逻辑备份工具是 mysqldump。当 mysqldump 使用参数 `–single-transaction` 的时候，导数据之前就会启动一个事务，来确保拿到一致性视图。而由于 MVCC 的支持，这个过程中数据是可以正常更新的。   

对于 MyISAM 这种不支持事务的引擎，mysqldump 工具就不能用了，所以 全局锁 虽然缺点很多，但是还是有存在的必要。   

### 表锁


### 行锁



### 参考

【高性能MySQL(第3版)】https://book.douban.com/subject/23008813/    
【MySQL 实战 45 讲】https://time.geekbang.org/column/100020801  
【MySQL技术内幕】https://book.douban.com/subject/24708143/    
【MySQL学习笔记】https://github.com/boilingfrog/Go-POINT/tree/master/mysql    



