## pgsql中的锁

### 前言

日常的工作中，对于同一个资源的操作，有时候我们难免要加上锁，以防止在操作中被别的进程给删除或者更改掉了。

### 用户可见的锁

> 从系统视图pg_locks中可见

用户可见的锁，用户自己能够主动调用的，可以在pg_locks中看到是否grant的锁。包括regular lock和咨询锁。

#### regular Lock
 
regular lock分为表级别和行级别两种。

##### 表级别

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
到前一个事务结束。



### 参考

【Postgresql锁机制（表锁和行锁）】https://blog.csdn.net/turbo_zone/article/details/84036511  
【postgresql行级锁for update测试】https://blog.csdn.net/shuoyu816/article/details/80086810  

