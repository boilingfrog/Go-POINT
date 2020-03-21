## EXPLAIN分析pgsql的性能

### 前言

对于pgsql中查询性能的分析，好像不想mysql中那么简单。当然pgsql中也是通过EXPLAIN进行
分析，那么就来认真中结下pgsql中explain的使用。

### EXPLAIN命令

EXPLAIN -- 显示一个语句的执行计划

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
ANALYZE选项通过实际执行的sql来获取相应的计划。这个是真正执行的，多以可以真是的看到执行计划花费了多少的
时间，还有它返回的行数。  

当然对于分析插入更新的语句，我们我们是可以把ANALYZE放到事物里面的，分析后之后回滚。  

````
BEGIN;
EXPLAIN ANALYZE ...;
ROLLBACK;
````





### 参考
【EXPLAIN】http://www.postgres.cn/docs/9.5/sql-explain.html  


