<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Redis 中的事务](#redis-%E4%B8%AD%E7%9A%84%E4%BA%8B%E5%8A%A1)
  - [什么是事务](#%E4%BB%80%E4%B9%88%E6%98%AF%E4%BA%8B%E5%8A%A1)
    - [1、原子性(Atomicity)](#1%E5%8E%9F%E5%AD%90%E6%80%A7atomicity)
    - [2、一致性(Consistency)](#2%E4%B8%80%E8%87%B4%E6%80%A7consistency)
    - [3、隔离性(Isolation)](#3%E9%9A%94%E7%A6%BB%E6%80%A7isolation)
    - [4、持久性(Durability)](#4%E6%8C%81%E4%B9%85%E6%80%A7durability)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Redis 中的事务

### 什么是事务

数据库事务( transaction )是访问并可能操作各种数据项的一个数据库操作序列，这些操作**要么全部执行,要么全部不执行，是一个不可分割的工作单位**。事务由事务开始与事务结束之间执行的全部数据库操作组成。    

事务必须满足所谓的ACID属性  

#### 1、原子性(Atomicity)

事务中的全部操作在数据库中是不可分割的，要么全部完成，要么全部不执行；  

- 整个数据库事务是不可分割的工作单位；  

- 只有使数据库中所有的数据库操作都执行成功，才算整个事务成功；  

- 事务中任何一个 SQL 执行失败，已经执行成功的 SQL 也必须撤回，数据库应该退回到执行事务之前的状态；  

#### 2、一致性(Consistency)

事务的执行使数据从一个状态转换为另一个状态,在事务开始之前和事务结束之后，数据库的完整性约束没有被破坏。    

有点绕，这里举个栗子  

如果一个名字字段，在数据库中是唯一属性，执行了事务之后，涉及到了对该字段的修改，事务执行过程中发生了回滚，之后该字段变的不唯一了，这种情况下就是破坏了事务的一致性要求。  

因为上面事务执行的过程中导致，导致里面名字字段属性的前后不一致，即数据库的状态从一中状态变成了一中不一致的状态。  

上面的这个栗子就是数据库没有遵循一致性的表现。  

#### 3、隔离性(Isolation)

事务的隔离性要求每个读写事务的对象对其他事务的操作对象相互分离，即该事务提交前对其他事务都不可见。  

通常使用锁来实现，数据库系统中会提供一种粒度锁的策略，允许事务仅锁住一个实体对象的子集，以此来提高事务之间的并发度。  

#### 4、持久性(Durability)

对于任意已提交事务，系统必须保证该事务对数据库的改变不被丢失，即使数据库出现故障。  

当时如果一些人为的或者自然灾害导致数据库机房被破坏，比如火灾，机房爆炸等。这种情况下锁提交的事务可能会丢失。  

因此可以理解，持久性保证的事务系统的高可靠性，而不是高可用性。   


### 参考

【Redis核心技术与实战】https://time.geekbang.org/column/intro/100056701    
【Redis设计与实现】https://book.douban.com/subject/25900156/   
【Redis 的学习笔记】https://github.com/boilingfrog/Go-POINT/tree/master/redis    
【数据库事务】https://baike.baidu.com/item/%E6%95%B0%E6%8D%AE%E5%BA%93%E4%BA%8B%E5%8A%A1/9744607  

