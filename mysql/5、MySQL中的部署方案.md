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

在主从复制中，从库利用主库上的 binlog 进行重播，实现主从同步，复制的过程中蛀主要使用到了 `dump thread，I/O thread，sql thread` 这三个线程。    

`IO thread`: 在从库执行 `start slave` 语句时创建，负责连接主库，请求 binlog，接收 binlog 并写入 relay-log；   

`dump thread`：用户主库同步 binlog 给从库，负责响应从 `IO thread` 的请求。主库会给每个从库的连接创建一个 `dump thread`，然后同步 binlog 给从库，如果该线程追上了主库，就会进行休眠状态，当主库有新的更新操作，这个线程就会被唤醒；    

`sql thread`：读取 `relay log` 执行命令实现从库数据的更新。   

来看下复制的流程：   

1、主库收到更新命令，执行更新操作，生成 binlog;  

2、从库在主从之间建立长连接；   

3、主库 dump_thread 从本地读取 binlog 传送刚给从库；  

4、从库从主库获取到 binlog 后存储到本地，成为 `relay log`（中继日志）；  

5、sql_thread 线程读取 `relay log` 解析、执行命令更新数据。       













### 参考

【高性能MySQL(第3版)】https://book.douban.com/subject/23008813/    
【MySQL 实战 45 讲】https://time.geekbang.org/column/100020801  
【MySQL技术内幕】https://book.douban.com/subject/24708143/    
【MySQL学习笔记】https://github.com/boilingfrog/Go-POINT/tree/master/mysql    
【MySQL文档】https://dev.mysql.com/doc/refman/8.0/en/replication.html  
【浅谈 MySQL binlog 主从同步】http://www.linkedkeeper.com/1503.html     
