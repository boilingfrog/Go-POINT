## EXPLAIN分析pgsql的性能

### 前言

对于pgsql中查询性能的分析，好像不想mysql中那么简单。当然pgsql中也是通过EXPLAIN进行
分析，那么就来认真中结下pgsql中explain的使用。

### EXPLAIN命令

#### EXPLAIN -- 显示一个语句的执行计划

````
EXPLAIN [ ( option [, ...] ) ] statement
EXPLAIN [ ANALYZE ] [ VERBOSE ] statement

这里 option可以是：

    ANALYZE [ boolean ]
    VERBOSE [ boolean ]
    COSTS [ boolean ]
    BUFFERS [ boolean ]
    TIMING [ boolean ]
    FORMAT { TEXT | XML | JSON | YAML }
````
ANALYZE选项通过实际执行的sql来获取相应的计划。这个是真正执行的，多以可以真是的看到执行计划花费了多少的时间，还有它返回的行数。  

当然对于分析插入更新的语句，我们我们是可以把ANALYZE放到事物里面的，分析后之后回滚。  

````
BEGIN;
EXPLAIN ANALYZE ...;
ROLLBACK;
````

#### VERBOSE

VERBOSE选项用于显示计划的附加信息。这些附加的信息有：计划中每个节点输出的各个列，如果触发器被触发，还会输出触发器的名称。该选项默认为FALSE。

COSTS选项显示每个计划节点的启动成本和总成本，以及估计行数和每行宽度。该选项默认是TRUE。

BUFFERS选项显示关于缓存区的信息。该选项只能与ANALYZE参数一起使用。显示的缓存区信息包括共享快块，本地块和临时读和写的块数。共享块、本地块和临时块分别
包含表和索引、临时快和临时索引、以及在排序和物化计划中使用的磁盘块。上层节点显示出来的数据块包含其所有子节点使用的块数。该选项默认为FALSE。

#### 全表扫描

全表扫描在pgsql中叫做顺序扫描(seq scan)，全表扫描就是把表的的所有的数据从头到尾读取一遍，然后从数据块中找到符合条件的数据块。


#### 索引扫描

索引是为了加快数据查询的速度索引而增加的。





### 参考
【EXPLAIN】http://www.postgres.cn/docs/9.5/sql-explain.html  


