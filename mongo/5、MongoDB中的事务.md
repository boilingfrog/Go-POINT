<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [MongoDB 事务](#mongodb-%E4%BA%8B%E5%8A%A1)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## MongoDB 事务

### 前言

在 MongoDB 中，对单个文档的操作都是原子的。因为可以在单个文档结构中使用内嵌文档和数据获得数据之间的关系，所以不必跨多个文档和集合进行范式化，这种
结构特性，避免了很多场景中的对多文档事务的需求。   

对于需要多个文档进行原子读写的场景，MongoDB 中引入了多文档事务和分布式事务。   

- 在4.0版本中，MongoDB支持副本集上的多文档事务；  

- 在4.2版本中，MongoDB 引入了分布式事务，增加了对分片集群上多文档事务的支持，并合并了对副本集上多文档事务的现有支持，事务可以跨多个操作、集合、数据库、文档和分片使用。   


### 参考

【MongoDB事务】https://docs.mongoing.com/transactions     


