
- [EXPLAIN分析pgsql的性能](#explain%E5%88%86%E6%9E%90pgsql%E7%9A%84%E6%80%A7%E8%83%BD)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [EXPLAIN命令](#explain%E5%91%BD%E4%BB%A4)
    - [EXPLAIN -- 显示一个语句的执行计划](#explain----%E6%98%BE%E7%A4%BA%E4%B8%80%E4%B8%AA%E8%AF%AD%E5%8F%A5%E7%9A%84%E6%89%A7%E8%A1%8C%E8%AE%A1%E5%88%92)
    - [命令详解](#%E5%91%BD%E4%BB%A4%E8%AF%A6%E8%A7%A3)
    - [EXPLAIN输出结果展示](#explain%E8%BE%93%E5%87%BA%E7%BB%93%E6%9E%9C%E5%B1%95%E7%A4%BA)
    - [analyze](#analyze)
    - [buffers](#buffers)
    - [全表扫描](#%E5%85%A8%E8%A1%A8%E6%89%AB%E6%8F%8F)
    - [索引扫描](#%E7%B4%A2%E5%BC%95%E6%89%AB%E6%8F%8F)
    - [位图扫描](#%E4%BD%8D%E5%9B%BE%E6%89%AB%E6%8F%8F)
    - [条件过滤](#%E6%9D%A1%E4%BB%B6%E8%BF%87%E6%BB%A4)
    - [Nestloop join](#nestloop-join)
    - [Hash join](#hash-join)
    - [Merge Join](#merge-join)
    - [Nested Loop，Hash JOin，Merge Join对比](#nested-loophash-joinmerge-join%E5%AF%B9%E6%AF%94)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

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

#### 命令详解

VERBOSE选项用于显示计划的附加信息。这些附加的信息有：计划中每个节点输出的各个列，如果触发器被触发，还会输出触发器的名称。该选项默认为FALSE。

COSTS选项显示每个计划节点的启动成本和总成本，以及估计行数和每行宽度。该选项默认是TRUE。

BUFFERS选项显示关于缓存区的信息。该选项只能与ANALYZE参数一起使用。显示的缓存区信息包括共享快块，本地块和临时读和写的块数。共享块、本地块和临时块分别
包含表和索引、临时快和临时索引、以及在排序和物化计划中使用的磁盘块。上层节点显示出来的数据块包含其所有子节点使用的块数。该选项默认为FALSE。

#### EXPLAIN输出结果展示

````sql
explain select  * from test1

QUERY PLAN
-----------------------------------------------------------------------------------------------------------------------
Seq Scan on test1  (cost=0.00..146666.56 rows=7999956 width=33)
````

`Seq Scan` 表示的是全表扫描，就是从头扫尾扫描一遍表里面的数据。  
`(cost=0.00..146666.56 rows=7999956 width=33)`的内容可以分成三部分  

- cost=0.00..146666.56 cost后面有两个数字,中间使用`..`分割，第一个数字`0.00`表示启动的成本，也就是返回第一行需要多少cost值；第二行表示返回所有数据的成本。  
- rows=7999956:表示会返回7999956行  
- width=33:表示每行的数据宽度为33字节  

其中`cost`描述的是一个`sql`执行的代价。

#### analyze

通过`analyze`可以看到更加精确的执行计划。  

````sql
explain  analyze select * from test1

QUERY PLAN
-----------------------------------------------------------------------------------------------------------------------
Seq Scan on test1  (cost=0.00..146666.56 rows=7999956 width=33) (actual time=0.012..1153.317 rows=8000001 loops=1)
Planning time: 0.049 ms
Execution time: 1637.480 ms
````

加了`analyze`可以看到实际的启动时间，`(actual time=0.012..1153.317 rows=8000001 loops=1)`其中:  

- actual time=0.012..1153.317:`0.012`表示的是启动的时间,`..`后面的时间表示返回所有行需要的时间  
- rows=8000001:表示返回的行数  

#### buffers

通过使用buffers来查看缓存区命中的情况  

````sql
explain  (analyze,buffers) select * from test1

QUERY PLAN
-----------------------------------------------------------------------------------------------------------------------
Seq Scan on test1  (cost=0.00..146666.56 rows=7999956 width=33) (actual time=0.013..1166.464 rows=8000001 loops=1)
  Buffers: shared hit=777 read=65890
Planning time: 0.049 ms
Execution time: 1747.163 ms
````

其中会多出一行`Buffers: shared hit=777 read=65890`。  

- shared hit=777：表示的在共享内存里面直接读到`777`个块；
- read=65890：表示

#### 全表扫描

全表扫描在`pgsql`中叫做顺序扫描(`seq scan`)，全表扫描就是把表的的所有的数据从头到尾读取一遍，然后从数据块中找到符合条件的数据块。


#### 索引扫描

索引是为了加快数据查询的速度索引而增加的(`Index Scan`)。索引扫描也就是我们的查询条件使用到了我们创建的索引，当然什么是索引，自行查阅资料吧。

#### 位图扫描

位图扫描也是走索引的一种方式。方法是扫描索引，那满足条件的行或者块在内存中建立一个位图，扫描完索引后，再根据位图到表的数据文件中把相应的数据读出来。
如果走了两个索引，可以把两个索引进行’and‘或‘or’计算，合并到一个位图，再到表的数据文件中把数据读出来。  

#### 条件过滤

就是在where后面加上过滤条件，扫描数据行，会找出满足条件过滤的行。条件过滤执行计划中显示为‘Filter’。

#### Nestloop join

对于被连接的数据子集较小的情况，`Nested Loop`是个较好的选择。`Nested Loop`是连表查询最朴素的一种连接方式。在嵌套循环的时候，内表被外表驱动，外表中返回的每一
行，都要在内表中检索找寻和它匹配的行。整个查询返回的结果集不能太大（>10000不适合），要把返回子集比较小的表作为外表，同时内表中连接查询的字段，最好能命中索
引，不然会有性能问题。  

执行的过程：确定一个驱动表（outer table），另一个表为inner table，驱动表中的每一行数据会去inner表中，查找检索数据。注意，驱动表的每一行都去inner表中
检索，索引驱动表的数据不能太大。对于inner表中的数据就没有限制了，只要创建的索引合适，inner表中数据的大小对查询的性能影响不大。

测试下：  

创建数据表

````sql
create table test1
(
    id          bigserial not null
        constraint test1_pk
            primary key,
    name        text      not null,
    category_id bigint    not null
);
create index test1_category_id_index
    on test1 (category_id);
  
create table test2
(
    id          bigserial not null
        constraint test2_pk
            primary key,
    name        text      not null,
    category_id bigint    not null
);

create index test2_category_id_index
    on test2 (category_id);
````

插入数据,test1插入8000000条，test2插入6000000条  

````sql
do $$
declare
v_idx integer := 1;
begin
  while v_idx < 8000000 loop
       v_idx = v_idx+1;
    insert into test1 (name, category_id) values ( random()*(20000000000000000-10)+10,random()*(8000000-10)+10);
  end loop;
end $$;


do $$
declare
v_idx integer := 1;
begin
  while v_idx < 6000000 loop
       v_idx = v_idx+1;
    insert into test2 (name, category_id) values ( random()*(20000000000000000-10)+10,random()*(6000000-10)+10);
  end loop;
end $$;
````

验证下查询  

````sql
explain select a.id,b.id from test1 a,test2 b where a.category_id=b.category_id and a.id<10000

QUERY PLAN
-----------------------------------------------------------------------------------------------------------------------
Gather  (cost=1170.58..67037.87 rows=14666 width=16)
  Workers Planned: 2
  ->  Nested Loop  (cost=170.58..64571.27 rows=6111 width=16)
        ->  Parallel Bitmap Heap Scan on test1 a  (cost=170.15..24939.15 rows=3748 width=16)
              Recheck Cond: (id < 10000)
              ->  Bitmap Index Scan on test1_pk  (cost=0.00..167.90 rows=8996 width=0)
                    Index Cond: (id < 10000)
        ->  Index Scan using test2_category_id_index on test2 b  (cost=0.43..10.55 rows=2 width=16)
              Index Cond: (category_id = a.category_id)
````

可以看到当有高选择性索引或进行限制性，查询优化器会自动选择`Nested Loop`


#### Hash join

优化器使用两个表较小的表，并利用连接键在内存中建立散列表，然后扫描较大的表并探测散列表，找出与散列表匹配的列。  

这种方式适用于较小的表可以完全放入到内存中，这样总成本就是访问两个表的成本之和。但是如果表都很大，不能放入到内存中，优化器会将它分割成若干个不同的分区，把不能放入到内存的
部分写入到临时段。此时要求有较大的临时段从而尽量提高I/O 的性能。它能够很好的工作于没有索引的大表和并行查询的环境中，并提供最好的性能。

优化器会自动选择较小的表，建立散列表，然后扫描另个较大的表。`Hash Join`只能应用于等值连接(如WHERE A.COL3 = B.COL4)，这是由Hash的特点决定的。  

测试下  

````sql
drop index test1_category_id_index;
drop index test2_category_id_index;
````
删除关联查询的字段的索引  

````sql
explain select a.id,b.id from test1 a,test2 b where a.category_id=b.category_id and a.id<10000

QUERY PLAN
-----------------------------------------------------------------------------------------------------------------------
Hash Join  (cost=214297.39..248684.97 rows=14666 width=16)
  Hash Cond: (a.category_id = b.category_id)
  ->  Index Scan using test1_pk on test1 a  (cost=0.43..335.86 rows=8996 width=16)
        Index Cond: (id < 10000)
  ->  Hash  (cost=109999.98..109999.98 rows=5999998 width=16)
        ->  Seq Scan on test2 b  (cost=0.00..109999.98 rows=5999998 width=16)
````

因为没有索引了，并且查询的条数可以完全放入到内存里面，所以查询优化器就选择使用`Hash join`了，对于选择那个表建立散列表，要看查询的条件。如上面的限制条件`a.id<10000`,
限制了a表查询的数据条数，那么a表条数较少，然后就在a表建立散列表，然后扫描b表。

#### Merge Join

通常情况下散列连接的效果比合并连接的效果好，如果源数据上有索引，或者结果已经排过序，在执行顺序合并连接时就不需要排序了，这时合并连接的性能会优于散列连接。  

`Merge join`的操作步骤：  
1' 对连接的每个表做table access full;  
2' 对table access full的结果进行排序;  
3' 进行merge join对排序结果进行合并。  

`Merge Join`可适于于非等值Join（>，<，>=，<=，但是不包含!=，也即<>）  

````sql
explain   select a.id,b.id from test1 a,test2 b where a.category_id=b.category_id 

QUERY PLAN
-----------------------------------------------------------------------------------------------------------------------
Merge Join  (cost=2342944.84..2575285.64 rows=13041968 width=16)
  Merge Cond: (b.category_id = a.category_id)
  ->  Sort  (cost=1005574.67..1020574.66 rows=5999998 width=16)
        Sort Key: b.category_id
        ->  Seq Scan on test2 b  (cost=0.00..124999.98 rows=5999998 width=16)
              Filter: (id < 100000000)
  ->  Materialize  (cost=1337364.94..1377364.72 rows=7999956 width=16)
        ->  Sort  (cost=1337364.94..1357364.83 rows=7999956 width=16)
              Sort Key: a.category_id
              ->  Seq Scan on test1 a  (cost=0.00..146666.56 rows=7999956 width=16)
````

`category_id`上面是没有索引的，这时候查询选择了`Merge Join`，上面的`Sort Key: a.category_id`，就是对a表的`category_id`字段排序。


#### Nested Loop，Hash JOin，Merge Join对比


|     类别                |     Nested Loop                | Hash Join                      |             Merge Join         |
| -----------------------| -------------------------------| -------------------------------| -------------------------------|
| 使用条件                | 任何条件                        | 等值连接（=）                    | 等值或非等值连接(>，<，=，>=，<=)，‘<>’除外              |
| 相关资源                | CPU、磁盘I/O                     | 内存、临时空间                    | 内存、临时空间               |
| 特点                   | 当有高选择性索引或进行限制性搜索时效率比较高，能够快速返回第一次的搜索结果。|当缺乏索引或者索引条件模糊时，Hash Join比Nested Loop有效。通常比Merge Join快,如果有索引，或者结果已经被排序了，这时候Merge Join的查询更快。在数据仓库环境下，如果表的纪录数多，效率高。|当缺乏索引或者索引条件模糊时，Merge Join比Nested Loop有效。非等值连接时，Merge Join比Hash Join更有效|
| 缺点                   | 当索引丢失或者查询条件限制不够时，效率很低；当表的纪录数多时，效率低。|为建立哈希表，需要大量内存。第一次的结果返回较慢。|所有的表都需要排序。它为最优化的吞吐量而设计，并且在结果没有全部找到前不返回数据。|




### 参考
【EXPLAIN】http://www.postgres.cn/docs/9.5/sql-explain.html  


