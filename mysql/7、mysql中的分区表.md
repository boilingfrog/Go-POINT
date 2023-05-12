## MySQL 中的分区表

### InnoDB 表空间

InnoDB 表空间（Tablespace）可以看做一个逻辑概念，InnoDB 把数据保存在表空间，本质上是一个或多个磁盘文件组成的虚拟文件系统。InnoDB 表空间不仅仅存储了表和索引，它还保存了回滚日志`（redo log）`、插入缓冲`（insert buffer）`、双写缓冲`（doublewrite buffer）`以及其他内部数据结构。   

从 InnoDB 存储引擎的逻辑存储结构来看，所有的数据都被逻辑的存放在一个空间中，也就是表空间。表空间，又有段（segment），区（extend）,页（page）组层。页在一些文档中也被称为块（block）。   

<img src="/img/mysql/mysql-innodb-space.jpg"  alt="mysql" />   


### 分区别表的概念

分区是将一个表的数据按照某个特定的键值，把不同的数据划分到不同的分区中。类似于数据库的分片处理

### 参考

【高性能MySQL(第3版)】https://book.douban.com/subject/23008813/    
【MySQL 实战 45 讲】https://time.geekbang.org/column/100020801  
【MySQL技术内幕】https://book.douban.com/subject/24708143/    
【MySQL学习笔记】https://github.com/boilingfrog/Go-POINT/tree/master/mysql    
【InnoDB 表空间】https://blog.csdn.net/u010647035/article/details/105009979  
