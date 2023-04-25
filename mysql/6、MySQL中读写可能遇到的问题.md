<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [MySQL 中读写分离可能遇到的问题](#mysql-%E4%B8%AD%E8%AF%BB%E5%86%99%E5%88%86%E7%A6%BB%E5%8F%AF%E8%83%BD%E9%81%87%E5%88%B0%E7%9A%84%E9%97%AE%E9%A2%98)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [读写分离的架构](#%E8%AF%BB%E5%86%99%E5%88%86%E7%A6%BB%E7%9A%84%E6%9E%B6%E6%9E%84)
    - [基于客户端实现读写分离](#%E5%9F%BA%E4%BA%8E%E5%AE%A2%E6%88%B7%E7%AB%AF%E5%AE%9E%E7%8E%B0%E8%AF%BB%E5%86%99%E5%88%86%E7%A6%BB)
    - [基于中间代理实现读写分离](#%E5%9F%BA%E4%BA%8E%E4%B8%AD%E9%97%B4%E4%BB%A3%E7%90%86%E5%AE%9E%E7%8E%B0%E8%AF%BB%E5%86%99%E5%88%86%E7%A6%BB)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## MySQL 中读写分离可能遇到的问题

### 前言

MySQL 中读写分离是经常用到了的架构了，通过读写分离实现横向扩展的能力，写入和更新操作在源服务器上进行，从服务器中进行数据的读取操作，通过增大从服务器的个数，能够极大的增强数据库的读取能力。   

### 读写分离的架构

常用的读写分离有下面两种实现：  

1、客户端实现读写分离；

2、基于中间代理层实现读写分离。   

#### 基于客户端实现读写分离

客户端主动做负载均衡，根据 `select、insert` 进行路由分类，读请求发送到读库中，写请求转发到写库中。  

这种方式的特点是性能较好，代码中直接实现，不需要额外的硬件支持，架构简单，排查问题更方便。     

缺点需要嵌入到代码中，需要开发人员去实现，运维无从干预，大型代码，实现读写分离需要改动的代码比较多。   

<img src="/img/mysql/mysql-client-readwrite.png"  alt="mysql" />    

#### 基于中间代理实现读写分离

中间代理层实现读写分离，在 MySQL 和客户端之间有一个中间代理层 proxy，客户端只连接 proxy， 由 proxy 根据请求类型和上下文决定请求的分发路由。   

<img src="/img/mysql/mysql-proxy-readwrite.png"  alt="mysql" />   

带 proxy 的架构，对客户端比较友好。客户端不需要关注后端细节，连接维护、后端信息维护等工作，都是由 proxy 完成的。但这样的话，对后端维护团队的要求会更高。而且，proxy 也需要有高可用架构。因此，带 proxy 架构的整体就相对比较复杂。  

不过那种部署方式都会遇到读写分离主从延迟的问题，因为主从延迟的存在，客户端刚执行完成一个更新事务，然后马上发起查询，如果选择查询的是从库，可能读取到的状态是更新之前的状态。    

### 主从读写延迟

主从延迟可能存在的原因：  

1、从库的性能比主库所在的机器性能较差；   

2、从库的压力大；   

3、大事务的执行

### 参考

【高性能MySQL(第3版)】https://book.douban.com/subject/23008813/    
【MySQL 实战 45 讲】https://time.geekbang.org/column/100020801  
【MySQL技术内幕】https://book.douban.com/subject/24708143/    
【MySQL学习笔记】https://github.com/boilingfrog/Go-POINT/tree/master/mysql    
【MySQL文档】https://dev.mysql.com/doc/refman/8.0/en/replication.html  
【浅谈 MySQL binlog 主从同步】http://www.linkedkeeper.com/1503.html     
