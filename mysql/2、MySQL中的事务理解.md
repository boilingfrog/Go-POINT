<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [MySQL 中的事务](#mysql-%E4%B8%AD%E7%9A%84%E4%BA%8B%E5%8A%A1)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [隔离性](#%E9%9A%94%E7%A6%BB%E6%80%A7)
    - [事务的隔离级别](#%E4%BA%8B%E5%8A%A1%E7%9A%84%E9%9A%94%E7%A6%BB%E7%BA%A7%E5%88%AB)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## MySQL 中的事务

### 前言

MySQL 中的事务操作，要么修改都成功，要么就什么也不做，这就是事务的目的。事务有四大特性 ACID，原子性，一致性，隔离性，持久性。   

- A(Atomic),原子性：指的是整个数据库事务操作是不可分割的工作单元，要么全部执行，要么都不执行；   

- C(Consistent),一致性：指的是事务将数据库从一种状态转换成下一种一致性状态，在事务开始之前和事务结束之后，数据库的完整性约束没有被破坏；   

数据的完整性： 实体完整性、列完整性（如字段的类型、大小、长度要符合要求）、外键约束等;  

业务的一致性：例如在银行转账时，不管事务成功还是失败，双方钱的总额不变。   

如果事务执行过程中，每个操作失败了，系统可以撤销事务，系统可以撤销事务，返回系统初始化的状态。   

- I(isolation): 隔离性还有其它称呼，如并发控制，可串行化，锁等，事务的隔离性要求每个读写事务对象对其他事务操作对象能够相互隔离，即事务提交之前对其它事务都不可见；   

- D(durability), 持久性: 指的是一旦数据提交，对数据库中数据的改变就是永久的。及时发生宕机，数据库也能回复恢复。    

下面来一一分析下这几个特性    

### 隔离性

#### 事务的隔离级别

MySQL 中标准的事务隔离级别包括：读未提交（read uncommitted）、读提交（read committed）、可重复读（repeatable read）和串行化（serializable ）。   

- 读未提交：一个事务还没提交时，它的变更就能被别的事务看到，读取未提交的数据也叫做脏读；      

- 读提交：一个事务提交之后，它的变更才能被其他的事务看到；   

- 可重复读：MySQL 中默认的实物隔离级别，一个事务执行的过程中看到的数据，总是跟这个事务在启动时看到的数据是一致的，在此隔离级别下，未提交的变更对其它事务也是不可见的，此隔离级别基本上避免了幻读；   

什么是幻读    

> The so-called phantom problem occurs within a transaction when the same query produces different sets of rows at different times. For example, if a SELECT is executed twice, but returns a row the second time that was not returned the first time, the row is a “phantom” row.

简单的讲就是，幻读指的是一个事务在前后两次查询同一个范围的时候，后一次查询看到了前一次查询没有看到的行。   

- 串行化：这是事务的最高级别，顾名思义就是对于同一行记录，“写”会加“写锁”，“读”会加“读锁”。当出现读写锁冲突的时候，后访问的事务必须等前一个事务执行完成，才能继续执行。     

串行化，不是所有的事务都串行执行，没冲突的事务是可以并发执行的。   

下面来详细的介绍下 `读提交` 和 `可重复读`    

准备下数据表  

```
create table user
(
id int auto_increment primary key,
username varchar(64) not null,
age int not null
);
insert into user values(2, "小张", 1);
```

来分析下下面的栗子  

<img src="/img/mysql/mysql-acid-durability.png"  alt="mysql" />  

**读未提交**   

V1、V2，V3 的值都是2，虽然事务1的修改还没有提交，但是读未提交的隔离能够看到事务未提交的数据，所以 V1 看到的数据就是 2 了。   

**读提交**  

V1 的值是1，V2 是2，V3 是2。因为事务1提交了，读提交可以看到提交的数据，所以 V2 的值就是2，V3 查询的结果肯定也是2了。   

**可重复读**

V1、V2 的值是1，V3 的值是 2。   

虽然事务1提交了，但是 V2 还是在事务2 中没有提交，根据可重复读的要求，一个事务执行的过程中看到的数据，总是跟这个事务在启动时看到的数据是一致的，所以 V2 也是 1。   

**串行化**

V1、V2 的值是1，V3 的值是 2。因为事务2，先启动查询，所以事务1必须等到事务2提交之后才能提交事务的修改，所以 V1、V2 的值是1，因为 V3 的查询时在事务1提交之后，所以 V3 查询的值就是2。   

#### 事务隔离是如何实现

在了解了四种隔离级别，下面来聊聊这几种隔离级别是如何实现的。   

首先来介绍一个非常重要的概念 `Read View`。   

`Read View` 是一个数据库的内部快照，用于 InnoDB 中 MVCC 机制。

##### 可重复度

可重复读主要是通过 `Read View` 来实现的，



### 参考

【高性能MySQL(第3版)】https://book.douban.com/subject/23008813/    
【MySQL 实战 45 讲】https://time.geekbang.org/column/100020801  
【MySQL技术内幕】https://book.douban.com/subject/24708143/    
【MySQL总结--MVCC（read view和undo log）】https://blog.csdn.net/huangzhilin2015/article/details/115195777     
【深入理解 MySQL 事务：隔离级别、ACID 特性及其实现原理】https://blog.csdn.net/qq_35246620/article/details/61200815     



