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




