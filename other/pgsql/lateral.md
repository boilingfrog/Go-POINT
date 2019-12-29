## pgsql中的lateral

我们先来看官方对lateral的定义

> 可以在出现于FROM中的子查询前放置关键词LATERAL。这允许它们引用前面的FROM项提供的列（如果没有
> LATERAL，每一个子查询将被独立计算，并且因此不能被其他FROM项交叉引用）。  
> 出现在FROM中的表函数的前面也可以被放上关键词LATERAL,但对于关键词的不可选的，在任何情况下函
> 数的参数都可以包含前面FROM项提供额列的引用。  
> 一个LATERAL项可以出现在FROM列表项层，或者出现在一个JOIN树中。在后者如果出现在JOIN的右部，
> 那么可以引用在JOIN左部分的任何项。  
> 如果一个FROM项包含LATERAL交叉引用，计算过程中：对于提供交叉引用列的FROM项的每一行，或者
> 多个提供这些列的多个FROM项进行集合，LATERAL项将被使用该行或者行集中的列值进行计算。得到结
> 果行将和它们被计算出来的行进行正常的连接。对于来自这些列的源表的每一行或行集，该过程将重复。

手册上提到：
````
SELECT * FROM foo, LATERAL (SELECT * FROM bar WHERE bar.id = foo.bar_id) ss;    
  
在 LATERAL (这里可以关联(引用)lateral左边的表或子句)  
  
所以允许:   
  
LATERAL (SELECT * FROM bar WHERE bar.id = foo.bar_id)  
````

### 带有LATERAL的SQL的计算步骤

1、逐行提取被lateral子句关联（引用）的FROM或JOIN的ITEM（也叫source table）的记录（s）
中的column(s)  
for each row of the FROM item providing the cross-referenced column(s),  
or set of rows of multiple FROM items providing the columns,  
2、使用以上提取的columns(s),关联计算lateral子句中的ITEM  
the LATERAL item is evaluated using that row or row set'us values of the columns.  
3、lateral的计算结果row(s),与所有from,join ITEM(S)正常的进行join计算
The resulting row(s) are joined as usual with the rows they were computed from.  
4、从1到3开始循环，直到所有的source table的行取尽。  
This is repeated for each row or set of rows from the column source table(s).


















