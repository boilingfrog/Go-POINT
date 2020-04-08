- [pgsql中的行锁](#pgsql%E4%B8%AD%E7%9A%84%E8%A1%8C%E9%94%81)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [用户可见的锁](#%E7%94%A8%E6%88%B7%E5%8F%AF%E8%A7%81%E7%9A%84%E9%94%81)
    - [regular Lock](#regular-lock)
    - [行级别](#%E8%A1%8C%E7%BA%A7%E5%88%AB)
        - [FOR UPDATE](#for-update)
        - [FOR NO KEY UPDATE](#for-no-key-update)
        - [FOR SHARE](#for-share)
        - [FOR KEY SHARE](#for-key-share)
  - [参考](#%E5%8F%82%E8%80%83)

## pgsql中的行锁

### 前言

日常的工作中，对于同一个资源的操作，有时候我们难免要加上锁，以防止在操作中被别的进程给删除或者更改掉了，那么更多还是行级锁，那么我们就来探究下。

### 用户可见的锁

> 从系统视图pg_locks中可见

用户可见的锁，用户自己能够主动调用的，可以在pg_locks中看到是否grant的锁。包括regular lock和咨询锁。

#### regular Lock
 
regular lock分为表级别和行级别两种。
 
#### 行级别

通过一些数据库操作自动获得一些行锁，行锁并不阻塞数据查询，只阻塞writes和locker，比如如下操作。

````sql
FOR UPDATE
FOR NO KEY UPFATE
FOR SHARE
FOR KEY SHARE
````
###### FOR UPDATE

FOR UPDATE锁可以使得SELECT语句获取行级锁，用于更新数据。锁定该行可以防止该行在本次的操作过程中，被其他的事务获取锁或者进行更改删除操作。就是说其他事务的操作会被阻塞直到当前事务结束；同样的，SELECT FOR UPDATE命令会等待直
到前一个事务结束。即尝试 UPDATE、DELETE、SELECT FOR UPDATE、SELECT FOR NO KEY UPDATE、SELECT FOR SHARE 或 SELECT FOR KEY SHARE 的其他事务将被阻塞。反过来，SELECT FOR UPDATE将等待已经在相同行上运行以上这些命令的
并发事务，并且接着锁定并且返回被更新的行（或者没有行，因为行可能已被删除）。不过，在一个REPEATABLE READ或SERIALIZABLE事务中，如果一个要被锁定的行在事务开始后被更改，将会抛出一个错误。  

任何在一行上的DELETE命令也会获得FOR UPDATE锁模式，在某些列上修改值的UPDATE也会获得该锁模式。当前UPDATE情况中被考虑的列集合是那些具有能用于外键的唯一索引的列（所以部分索引和表达式索引不被考虑），但是这种要求未来有可能会改变。

###### FOR NO KEY UPDATE

和FOR UPDATE命令类似，但是对于获取锁的要求更加宽松一些，在同一行中不会阻塞SELECT FOR KEY SHARE命令。同样在UPDATE命令的时候如果没有获取到FOR UPDATE锁的情况下会获取到该锁。

###### FOR SHARE

行为与FOR NO KEY UPDATE类似，不过它在每个检索到的行上获得一个共享锁而不是排他锁。一个共享锁会阻塞其他事务在这些行上执行UPDATE、DELETE、SELECT FOR UPDATE或者SELECT FOR NO KEY UPDATE，但是它不会阻止它们执行SELECT FOR SHARE或者SELECT FOR KEY SHARE。

###### FOR KEY SHARE

行为与FOR SHARE类似，不过锁较弱：SELECT FOR UPDATE会被阻塞，但是SELECT FOR NO KEY UPDATE不会被阻塞。一个键共享锁会阻塞其他事务执行修改键值的DELETE或者UPDATE，但不会阻塞其他UPDATE，也不会阻止SELECT FOR NO KEY UPDATE、SELECT FOR SHARE或者SELECT FOR KEY SHARE。
 
![](https://img2020.cnblogs.com/blog/1237626/202004/1237626-20200408142154977-1951304776.png)

举个栗子


### 参考

【Postgresql锁机制（表锁和行锁）】https://blog.csdn.net/turbo_zone/article/details/84036511  
【postgresql行级锁for update测试】https://blog.csdn.net/shuoyu816/article/details/80086810  
【PostgreSQL 锁解密 】https://www.oschina.net/translate/postgresql-locking-revealed  
【显式锁定】http://postgres.cn/docs/11/explicit-locking.html

