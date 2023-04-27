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

### 主备同步延迟

主备同步延迟，就是同一个事务，在备库执行完成的时间和主库执行完成的时间之间的差值。   

1、主库 A 执行完成一个事务，并且写入到 binlog ，记录这个时间为 T1；  

2、传递数据给从库，从库接收这个 binlog，接收完成的时间记为 T2;  

3、从库 B 执行完成这个接收的事务，这个时间记为 T3。   

主备延迟的时间就是 T3-T1 之间的时间差。    

通过 `show slave status` 命令能到 seconds_behind_master 这个值就表示当前备库延迟了多少秒。    

seconds_behind_master 的计算方式：  

1、每个事务的 binlog 都有一个时间字段，用于记录主库写入的时间；  

2、从库取出当前正在执行的事务的时间字段的值，计算他与当前系统时间的差值，就能得到 seconds_behind_master。   

简单的讲 seconds_behind_master 就是上面 `T3 -T1` 的时间差值。      

如果主从机器的时间设置的不一致，会不会导致主备延迟的不准确？  

答案是不会的，备库连接到主库，会通过 `SELECT UNIX_TIMESTAMP()`函数来获取当前主库的时间，如果这时候发现主库系统时间与自己的不一致，备库在执行 seconds_behind_master 计算的时候会主动扣减掉这差值。    

### 主从读写延迟

主从延迟可能存在的原因：  

1、从库的性能比主库所在的机器性能较差；   

从库的性能比较查，如果从库的复制能力，低于主库，那么在主库写入压力很大的情况下，就会造成从库长时间数据延迟的情况出现。   

2、从库的压力大；     

大量查询放在从库上，可能会导致从库上耗费了大量的 CPU 资源，进而影响了同步速度，造成主从延迟。     

3、大事务的执行；    

有事务产生的时候，主库必须要等待事务完成之后才能写入到 binlog，假定执行的事务是一个非常大的数据插入，这些数据传输到从库，从库同步这些数据也需要一定的时间，就会导致从节点出现数据延迟。  

4、主库异常发生主备切换。   

发生主备切换的时候，可能会出现延迟，主备切换会有下面两种策略：      

可靠性优先策略：   

1、首先判断下备库的 seconds_behind_master ，如果小于某个可以容忍的值，就进行下一步，否则持续重试这一步；  

2、把主库 A 改成只读状态，设置 readonly 为 true;  

3、判断备库 B 的 seconds_behind_master，直到这个值变成 0 为止；  

4、把备库 B 改成可读写状态，设置 readonly 为 false;   

5、业务请求切换到备库 B。    

这个切换的过程中是存在不可用时间的，在步骤 2 之后，主库 A 和备库 B 都处于 readonly 状态，这时候系统处于不可写状态，知道备库 B readonly 状态变成 false，这时候才能正常的接收写请求。    

步骤 3 判断 seconds_behind_master 为 0，这个操作是最耗时的，通过步骤 1 中的提前判断，可以确保 seconds_behind_master 的值足够小，能够减少步骤 3 的等待时间。    

可用性优先策略：   

如果把步骤4、5调整到最开始执行，不等主库的数据同步，直接把连接切到备库 B，让备库 B 可以直接读写，这样系统就几乎没有不可用时间了。   

这种策略能最大可能保障服务的可用性，但是会出现数据不一致的情况。  

下面来分析下一种数据不一致的情况：     

```go
CREATE TABLE `t` (
`id` int(11) NOT NULL AUTO_INCREMENT,
`a` int(11) DEFAULT NULL,
PRIMARY KEY (`id`)
) ENGINE=InnoDB;

insert into t values(1,1);
insert into t values(2,2);
```

可以看到上面的表定义了一个自增主键，同时插入了两条数据。   

假定这时候数据库压力很大，主从库有延迟。主库在插入一条数据 `INSERT INTO `t` (`a`) VALUES (3);` 发生了主备切换







### 参考

【高性能MySQL(第3版)】https://book.douban.com/subject/23008813/    
【MySQL 实战 45 讲】https://time.geekbang.org/column/100020801  
【MySQL技术内幕】https://book.douban.com/subject/24708143/    
【MySQL学习笔记】https://github.com/boilingfrog/Go-POINT/tree/master/mysql    
【MySQL文档】https://dev.mysql.com/doc/refman/8.0/en/replication.html  
【浅谈 MySQL binlog 主从同步】http://www.linkedkeeper.com/1503.html     
