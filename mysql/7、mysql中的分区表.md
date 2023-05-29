<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [MySQL 中的分区表](#mysql-%E4%B8%AD%E7%9A%84%E5%88%86%E5%8C%BA%E8%A1%A8)
  - [InnoDB 逻辑存储结构](#innodb-%E9%80%BB%E8%BE%91%E5%AD%98%E5%82%A8%E7%BB%93%E6%9E%84)
    - [表空间](#%E8%A1%A8%E7%A9%BA%E9%97%B4)
    - [段](#%E6%AE%B5)
    - [区](#%E5%8C%BA)
  - [分区别表的概念](#%E5%88%86%E5%8C%BA%E5%88%AB%E8%A1%A8%E7%9A%84%E6%A6%82%E5%BF%B5)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## MySQL 中的分区表

### InnoDB 逻辑存储结构

InnoDB 存储引擎的存储逻辑，所有数据都会被逻辑的保存在一个空间中，这个空间称为 InnoDB（Tablespace） 表空间，本质上是一个或多个磁盘文件组成的虚拟文件系统。InnoDB 表空间不仅仅存储了表和索引，它还保存了回滚日志`（redo log）`、插入缓冲`（insert buffer）`、双写缓冲`（doublewrite buffer）`以及其他内部数据结构。   

表空间，由段（segment），区（extend）,页（page）组层。页在一些文档中也被称为块（block）。   

<img src="/img/mysql/mysql-innodb-space.jpg"  alt="mysql" />   

#### 表空间

表空间可以看做是 InnoDB 存储引擎逻辑结构的最高层，所有的数据都存放在表空间中。   

在默认情况下，InnoDB 存储引擎有一个共享的表空间 ibdatal 即所有的数据都存放在这个表空间中。如果用户开启了 innodb_file_per_table，则每张表内的数据都可以单独的存放到一个表空间中。     

不过即使开启了 innodb_file_per_table，每张表内的存放的数据只是数据，索引和插入缓存 Bitmap 页，其它类的数据回滚信息，插入缓存页，系统事务信息，二次写缓冲，还是会放在原来的共享表空间中。这也能说明一个问题，及时开始了 innodb_file_per_table，贡献表空间的数据还是会不断的增大。   

#### 段

表空间是由各个段组层的，常见的段由数据段，索引段，回滚段等。  

对于 InnoDB 存储引擎的存储数据存储结构 B+ 树，只有叶子节点存储数据，其它节点存储索引信息。也就是叶子节点存储的是数据段，非叶子节点存储的是索引段。   

在 InnoDB 存储引擎中，对段的管理都是由引擎自身所完成的，我们一般不能也没有必要对其进行操作。    

#### 区


### 分区别表的概念

分区是将一个表的数据按照某个特定的键值，把不同的数据划分到不同的分区中。类似于数据库的分片处理

### 参考

【高性能MySQL(第3版)】https://book.douban.com/subject/23008813/    
【MySQL 实战 45 讲】https://time.geekbang.org/column/100020801  
【MySQL技术内幕】https://book.douban.com/subject/24708143/    
【MySQL学习笔记】https://github.com/boilingfrog/Go-POINT/tree/master/mysql    
【InnoDB 表空间】https://blog.csdn.net/u010647035/article/details/105009979  
