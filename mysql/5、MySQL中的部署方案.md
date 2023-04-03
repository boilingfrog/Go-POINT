<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [MySQL 中的集群部署方案](#mysql-%E4%B8%AD%E7%9A%84%E9%9B%86%E7%BE%A4%E9%83%A8%E7%BD%B2%E6%96%B9%E6%A1%88)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [Replication](#replication)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## MySQL 中的集群部署方案

### 前言

这里来聊聊，MySQL 中常用的部署方案。  

常用的 MySQL 的部署方案：   

1、MySQL Replication；  

2、InnoDB Cluster；     

3、InnoDB ReplicaSet;  

### Replication

`MySQL Replication`,主从复制集群，使一台 MySQL 数据库服务器的的数据能够复制到一台或者多台 MySQL 的数据库服务器中。  

优点：  

1、通过读写分离实现横向扩展的能力，写入和更新操作在源服务器上进行，从服务器中进行数据的读取操作，通过增大从服务器的个数，能够极大的增强数据库的读取能力；    

2、数据安全，因为副本可以暂停复制过程，所以可以在副本上运行备份服务而不会破坏相应的源数据；   

3、方便进行数据分析，可以在写库中创建实时数据，数据的分析操作在从库中进行，不会影响到源数据库的性能；   

实现原理  









### 参考

【高性能MySQL(第3版)】https://book.douban.com/subject/23008813/    
【MySQL 实战 45 讲】https://time.geekbang.org/column/100020801  
【MySQL技术内幕】https://book.douban.com/subject/24708143/    
【MySQL学习笔记】https://github.com/boilingfrog/Go-POINT/tree/master/mysql    
【MySQL文档】https://dev.mysql.com/doc/refman/8.0/en/replication.html  
