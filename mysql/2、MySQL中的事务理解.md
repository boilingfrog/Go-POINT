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

- 读未提交：一个事务还没提交时，它的变更就能被别的事务看到；    

- 读提交：一个事务提交之后，它的变更才能被其他的事务看到；   

- 可重读读：一个事务执行的过程中看到的数据，总是跟这个事务在启动时看到的数据是一致的，在此隔离级别下，未提交的变更对其它事务也是不可见的；   

- 串行化：这是事务的最高级别，在每条读据上，加锁，写就加写锁，读就加读锁，当出现读写冲突的时候，后面的事务必须等前面的事务执行成功，才能继续执行。    

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



### 参考

【高性能MySQL(第3版)】https://book.douban.com/subject/23008813/    
【MySQL 实战 45 讲】https://time.geekbang.org/column/100020801  
【MySQL技术内幕】https://book.douban.com/subject/24708143/    



