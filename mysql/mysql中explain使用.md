<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [MySQL 中的 Explain 使用](#mysql-%E4%B8%AD%E7%9A%84-explain-%E4%BD%BF%E7%94%A8)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [分析](#%E5%88%86%E6%9E%90)
    - [1、id](#1id)
    - [2、select_type](#2select_type)
    - [3、table](#3table)
    - [4、type](#4type)
    - [5、possible_keys](#5possible_keys)
    - [6、key](#6key)
    - [7、key_len](#7key_len)
    - [8、ref](#8ref)
    - [9、rows](#9rows)
    - [10、Extra](#10extra)
  - [总结](#%E6%80%BB%E7%BB%93)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## MySQL 中的 Explain 使用  

### 前言

在开发的过程中,对于我们写的sql语句，我们有时候会考虑sql语句的性能，那么explain就是首选。Explain命令在解决数据库性能上是第一推荐使用命令，大部分的性能问题可以通过此命令来简单的解决，Explain可以用来查看 SQL 语句的执行效 果，可以帮助选择更好的索引和优化查询语句，写出更好的优化语句。

Explain语法：`explain select … from … [where ...]`   

### 分析

例如：`explain select * from news;`

输出：

```
+----+-------------+-------+-------+-------------------+---------+---------+-------+------+-------+
| id | select_type | table | type | possible_keys | key | key_len | ref | rows | Extra |
+----+-------------+-------+-------+-------------------+---------+---------+-------+------+-------+

```

下面对各个属性进行了解：  

#### 1、id

实际上每次 select 查询都会对应一个 id，它代表着 SQL 执行的顺序，如果 id 值越大，说明对应的 SQL 语句执行的优先级越高。在一些复杂的查询 SQL 语句中常常包含一些子查询，那么id序号就会递增，如果出现嵌套查询，我们可以发现最里层的查询对应的 id 最大，因此也优先被执行。  

#### 2、select_type  

select_type 表示执行的计划对应的查询时什么类型。  

会有下面几种类型：   

1、SIMPLE: 简单的 select 查询，不包含子查询或者UNION查询；  

2、PRIMARY: 查询中包含子查询，最为层查询被标记为 PRIMARY；   

3、SUBQUERY：在 select 或 where 列表中包含了子查询，该子查询被标记为 SUBQUERY；   

4、UNION: UNION 之后出现的 select 语句对应的查询类型会标记此类型;   

5、UNION RESULT: UNION 的结果；      

6、SUBQUERY: 子查询中的第一个 select；  

7、DEPENDENT SUBQUERY: 子查询中的第一个 select，取决于外面的查询。   

#### 3、table

table 表示表名称，表示要查询哪张表，不一样是真实的表，可能是表的别名或者临时表。   

#### 4、type  

这列最重要，显示了连接使用了哪种类别,有无使用索引，是使用 Explain 命令分析性能瓶颈的关键项之一。   

type 字段可能的值   

- system: 表示只有一行记录(等于系统表)，这是 const 类型的特列，基本上不太会遇到，可忽略；

- const: 表中只有一行数据数据被匹配，代表该查询使用主键或者唯一键进行查询，可以直接返回需要查询的记录，效率最高；

- eq_ref: 类似 ref，区别就在使用的索引是唯一索引，对于每个索引键值，表中只有一条记录匹配，简单来说，就是多表连接中使用 `primary key` 或者 `unique key` 作为关联条件;

- ref: 使用非唯一索引扫描或者唯一索引的前缀扫描，返回匹配某个单独值的记录行;

- index: 对索引列进行全索引扫描，这种就需要优化了；

- all: 对整张表进行全表扫描，这是最糟糕的情况。    

type 值的评判标准，结果值从好到坏依次是  

`system > const > eq_ref > ref > fulltext > ref_or_null > index_merge > unique_subquery > index_subquery > range > index > ALL`

一般来说，得保证查询至少达到 range 级别，最好能达到 ref，否则就可能会出现性能问题。    

#### 5、possible_keys

指出MySQL能使用哪个索引在表中找到记录，查询涉及到的字段上若存在索引，则该索引将被列出，但不一定被查询使用。   

#### 6、key

显示 MySQL 实际决定使用的键（索引）。如果没有选择索引，键是 NULL。  

#### 7、key_len 

显示 MySQL 决定使用的键长度。如果键是 NULL，则长度为 NULL。使用的索引的长度。在不损失精确性的情况下，长度越短越好。   

#### 8、ref 

显示使用哪个列或常数与 key 一起从表中选择行。  

#### 9、rows  

显示 MySQL 认为它执行查询时必须检查的行数。  

#### 10、Extra  

包含 MySQL 解决查询的详细信息，也是关键参考项之一。   

1、Using index  

该值表示在进行数据查询的时候使用了覆盖索引，而不用进行回表的查询了。

索引确实能够提高查询的效率，但二级索引会有某些情况会存在二次查询也就是回表的问题，这种情况合理的使用覆盖索引，能够提高索引的效率，减少回表的查询。

覆盖索引将需要查询的值建立联合索引，这样索引中就能包含查询的值，这样查询如果只查询 索引中的值和主键值 就不用进行二次查询了，因为当前索引中的信息已经能够满足当前查询值的请求。
   
2、Using where  

表示查询的时候未找到可用的索引，需要通过 where 条件过滤所需要的的数据，但要注意的是并不是所有带 where 语句的查询都会显示 `Using where`。  

Using where 的作用只是提醒我们 MySQL 将用 where 子句来过滤结果集。这个一般发生在 MySQL 服务器，而不是存储引擎层。一般发生在不能走索引扫描的情况下或者走索引扫描，但是有些查询条件不在索引当中的情况下。   

3、Using temporary

表示 MySQL 需要使用临时表来存储结果集，常见于排序和分组查询。  

这个值表示使用了内部的临时表，一个查询可能会用到多个临时表，有很多原因都会导致MySQL在执行查询期间创建临时表。两个常见的原因是在来自不同表的上使用了DISTINCT,或者使用了不同的 ORDER BY和 GROUP BY 列。  

4、Using filesort  

表示查询后结果需要使用临时表来存储，一般在排序或者分组查询时用到。   

5、Using join buffer  

改值强调了在获取连接条件时没有使用索引，并且需要连接缓冲区来存储中间结果。如果出现了这个值，那应该注意，根据查询的具体情况可能需要添加索引来改进能。  

6、Using Index Condition

在 `MySQL 5.6` 版本后加入的新特性`（Index Condition Pushdown）`，索引下推，主要是用来通过减少回表的次数，提高查询的性能。简单点讲就是在索引遍历过程中，对索引中包含的字段先做判断，直接过滤掉不满足条件的记录，减少回表次数。   

`Using index condition` 的出现，意味着 MySQL 在执行这条语句时使用了`Index Condition Pushdown`(ICP)。    

### 总结

面对一些慢查询，通过 Explain 能够帮助我们分析出一些性能的瓶颈。   
