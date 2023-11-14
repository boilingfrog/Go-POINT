<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [MongoDB 事务](#mongodb-%E4%BA%8B%E5%8A%A1)
  - [前言](#%E5%89%8D%E8%A8%80)
    - [如何使用](#%E5%A6%82%E4%BD%95%E4%BD%BF%E7%94%A8)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## MongoDB 事务

### 前言

在 MongoDB 中，对单个文档的操作都是原子的。因为可以在单个文档结构中使用内嵌文档和数据获得数据之间的关系，所以不必跨多个文档和集合进行范式化，这种
结构特性，避免了很多场景中的对多文档事务的需求。   

对于需要多个文档进行原子读写的场景，MongoDB 中引入了多文档事务和分布式事务。   

- 在4.0版本中，MongoDB支持副本集上的多文档事务；  

- 在4.2版本中，MongoDB 引入了分布式事务，增加了对分片集群上多文档事务的支持，并合并了对副本集上多文档事务的现有支持，事务可以跨多个操作、集合、数据库、文档和分片使用。   

#### 如何使用

MongoDB 中从 4.0 开始支持了事务，这里来看看 MongoDB 中的事务是如何使用的呢？   

首先登陆 MongoDB 命令行   

```
mongo -u <name> --port 27017 --host 127.0.0.1 admin  -p <pass>
```

1、打开 session

```
session = db.getMongo().startSession( { readPreference: { mode: "primary" } } );
```

2、将需要操作的 collection 进行变量绑定  

```
testCollection = session.getDatabase("gleeman").test_explain;
test1Collection = session.getDatabase("gleeman").test_explain_1;
```

3、开始事务标注，指定MVCC的模式，写模式  

```
session.startTransaction( { readConcern: { level: "snapshot" }, writeConcern: { w: "majority" } } );
```

4、拼接执行语句，将需要执行的语句进行事务封装  

```
try (ClientSession clientSession = client.startSession()) {
    clientSession.startTransaction();
    collection.insertOne(clientSession, docOne);
    collection.insertOne(clientSession, docTwo);
    clientSession.commitTransaction();
}
```

5、提交事务  

```
session.commitTransaction();
```

6、关闭session

```
session.endSession();
```

### 事务的原理  

MongoDB 中的 WiredTiger 存储引擎是目前使用最广泛的，这里主要介绍下 WiredTiger 中事务的实现原理。   

WiredTiger 存储引擎支持 `read-uncommitted 、read-committed` 和 `snapshot` 3 种事务隔离级别，MongoDB 启动时默认选择 `snapshot` 隔离。      

#### 事务和复复制集以及存储引擎之间的关系  

1、事务和复制集  

复制集配置下，MongoDB 整个事务在提交时，会记录一条 oplog，包含了事务所有的操作，备节点拉取 oplog，并在本地重放事务操作。事务 oplog 包含了事务操作的 `lsid，txnNumber`，以及事务内所有的操作日志（ `applyOps` 字段）。   

WiredTiger 是如何实现事务和 ACID 呢。WiredTiger 事务主要使用了三个技术 snapshot(事务快照)、MVCC (多版本并发控制)和 `redo log`(重做日志)。同时为了实现这三个技术，还定义了一个基于这三个技术的事务对象和全局事务管理器。        

```
wt_transaction{
	transaction_id:    // 本次事务的全局唯一的ID，用于标示事务修改数据的版本号
	snapshot_object:   // 当前事务开始或者操作时刻其他正在执行且并未提交的事务集合,用于事务隔离
	operation_array:   // 本次事务中已执行的操作列表,用于事务回滚。
	redo_log_buf:      // 操作日志缓冲区。用于事务提交后的持久化
	State:             // 事务当前状态
}

```

WiredTiger 中的 MVCC 是基于 `key/value` 中 value 值的链表，每个链表单元中存储有当先版本操作的事务 ID 和操作修改后的值。   

```
wt_mvcc{
	transaction_id:    // 本次修改事务的ID
	value:             // 本次修改后的值
}
```

WiredTiger 中数据修改都是在这个链表中进行 append 操作，每次对值的修改都是 append 到链表头，每次读取值的时候读是从链表头根据对应的修改事务 transaction_id 和本次事务的 snapshot 来判断是否可读，如果不可读，向链表尾方向移动，直到找到都事务可以读到的数据版本。    

什么是 snapshot 呢？   

事务开始或者结束操作之前都会对整个 WiredTiger 引擎内部正在执行的或者将要执行的事务进行一次快照，保存当时整个引擎的事务状态，确定那些事务是对自己可见的，哪些事务是自己不可见的。   

### WiredTiger 事务过程

一般事务有三个阶段：开启事务，执行事务，提交事务。如果事务执行失败，会进行事务的回滚操作，事务正常执行，最近进行事务的提交 (commit) 即可。     

#### 事务开启

事务开启的过程中，首先会为事务创建一个事务对象并把这个对象加入到全局的事务管理器当中，然后根据配置确定事务的隔离级别和 redo_log 的刷盘方式，并将事务状态设置成执行状态，最后判断事务的隔离级别，如果是 snapshot 级的事务隔离，在本次事务执行之前会创建一个系统并发事务的 snapshot 截屏，，保存当时整个引擎的事务状态，确定那些事务是对自己可见的，哪些事务是自己不可见的。

#### 事务执行  

事务在执行阶段，如果是读操作，不做任何处理，因为读操作不需要回滚和提交。如果是写操作，WiredTiger 会对每个操作做详细的记录。   

这里就会用使用到上面介绍的事务对象(wt_transaction)中的 operation_array 和 redo_log_buf。   

operation_array：主要记录本次事务中已经提交的操作列表，数组单元中，会包含一个指向 MVCC list 对应修改版本值的指针，用于事务的回滚。   

redo_log_buf: 操作日志缓冲区。用于事务提交后的持久化。   

来描述下具体的更新操作过程：   

1、创建一个 mvcclist 的值对象 update；  

2、根据事务对象的 transaction_id 和事务状态判断是为本次事务创建写的事务id,如果没有，为本次事务分配一个事务id,并将事务的状态设置成  HAS_TXN_ID 状态；  

3、将本次事务的 ID 设置到 update 单元中作为 mvcc 版本号；  

4、同时会创建一个 operation 对象，这个对象的指针会指向 update,这个对象会加入到 operation_array 中，用来进行操作事务的回滚；  

5、update 会被加入到 mvcclist 的头部；  

6、最后会写一条 redo_log 到本次事务的 redo_log_buffer 当中。   

#### 事务提交  

事务提交  

提交事务对象中的 redo_log_buf 中的数据到 redo_log_file(重做日志中)，并将 redo_log_file 持久化到磁盘上，清除提交事务对象的 snapshot，再将事务对象的transaction_id 设置成 WT_TNX_NODE，保证其他事务在创建 snapshot 时本次事务的状态是已提交的状态。     

#### 事务回滚

WiredTiger 引擎对事务的回滚过程比较简单，首先遍历 operation_array ，对每个数组单元对应的 update 事务 id 设置一个 WT_TXN_ABORTED ,标识 mvcc 对应的事务单元被回滚，在其它事务进行 mvcc 读操作的时候，跳过这个放弃的值即可，整个过程是一个无锁的操作，高效，简洁。   






### 参考

【MongoDB事务】https://docs.mongoing.com/transactions     
【WiredTiger的事务实现详解 】https://blog.csdn.net/daaikuaichuan/article/details/97893552  


