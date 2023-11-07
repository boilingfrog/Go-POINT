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





### 参考

【MongoDB事务】https://docs.mongoing.com/transactions     


