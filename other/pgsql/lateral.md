## pgsql中的lateral


- [什么是LATERAL](#%e4%bb%80%e4%b9%88%e6%98%afLATERAL)
- [带有LATERAL的SQL的计算步骤](#%e5%b8%a6%e6%9c%89LATERAL%e7%9a%84SQL%e7%9a%84%e8%ae%a1%e7%ae%97%e6%ad%a5%e9%aa%a4)
- [LATERAL在OUTER JOIN中的使用限制（或定义限制）](#LATERAL%e5%9c%a8OUTER+JOIN%e4%b8%ad%e7%9a%84%e4%bd%bf%e7%94%a8%e9%99%90%e5%88%b6%ef%bc%88%e6%88%96%e5%ae%9a%e4%b9%89%e9%99%90%e5%88%b6%ef%bc%89)
- [LATERAL的几个简单的例子](#LATERAL%e7%9a%84%e5%87%a0%e4%b8%aa%e7%ae%80%e5%8d%95%e7%9a%84%e4%be%8b%e5%ad%90)
- [总结](#%e6%80%bb%e7%bb%93)

### 举几个我经常使用的栗子
首先说下场景：  
有个一个商品表goods，还有一个表一个评价表evaluations。商品表和评价表是一对多的。  
1、在一个后台，我想查询商品的信息，同时查询这个商品的评价的数量。  
我们可以通过这样来实现
````
SELECT 
    g.*,
    COUNT(e.*) as num
FROM goods as g
LEFT JOIN evaluation as e on e.goods_id=g.id
WHERE 1=1  GROUP BY g.id
````
通过左连接，加上分组就能实现了  
那么也可以使用lateral来实现  
````
SELECT 
    g.*,
    e.num
FROM goods as g
LEFT JOIN LATERAL(
    SELECT COUNT(ev.id) as num FROM  evaluation AS ev 
    WHERE   ev.goods_id=g.id
) AS e ON TRUE
WHERE 1=1 
````
就这样好像lateral的优势不是那么明显。  
2、我们查询评论数目大于3的商品的信息
````
SELECT 
    g.*,
    COUNT(e.*) as num
FROM goods as g
LEFT JOIN evaluation as e on e.goods_id=g.id
HAVING COUNT(e.*)>3   GROUP BY g.id 
````
这样就不行了，查询不到了。  
这时候就需要使用LATERAL  
````
SELECT 
    g.*,
    e.num
FROM goods as g
LEFT JOIN LATERAL(
    SELECT COUNT(ev.id) as num FROM  evaluation AS ev 
    WHERE   ev.goods_id=g.id
) AS e ON TRUE
WHERE 1=1 AND num>3
````
3、然后我们再次查询这些商品的信息，希望找到**黄金会员**评论的商品信息  
这时候LATERAL的优势就更加明显了  
````
SELECT 
    g.*,
    e.num
FROM goods as g
LEFT JOIN LATERAL(
     SELECT COUNT(ev.id) as num
     FROM  evaluation AS ev 
     LEFT JOIN  users u  on u.id=ev.user_id
     WHERE   ev.goods_id=g.id AND u.grade=9
) AS e ON TRUE
WHERE 1=1 AND num>0
````

### 什么是LATERAL

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

### LATERAL在OUTER JOIN中的使用限制（或定义限制）
由于lateral的计算步骤是从source table逐条展开的，所以OUTER JOIN时只能使用source table
作为whole端，LATERAL内的ITEM不能作为WHOLE端。  
因此lateral只能在left join的右边。或者right join的左边。因此不能是WHOLE端。
````
The column source table(s) must be INNER or LEFT joined to the LATERAL item,   
  
else there would not be a well-defined set of rows from which to compute each set of rows for the LATERAL item.   
  
Thus, although a construct such as X RIGHT JOIN LATERAL Y is syntactically valid,   
  
it is not actually allowed for Y to reference X.  
````

### LATERAL的几个简单的例子

1、  
A trivial example of LATERAL is
````
SELECT * FROM foo, LATERAL (SELECT * FROM bar WHERE bar.id = foo.bar_id) ss;
````
This is not especially useful since it has exactly the same result as the more conventional
````
SELECT * FROM foo, bar WHERE bar.id = foo.bar_id;
````
一个LATERAL项可以出现在FROM列表项层  
 2、  
**LATERAL** is primarily useful when the cross-referenced column is necessary 
for computing the row(s) to be joined. A common application is providing an 
argument value for a set-returning function. For example, supposing that 
vertices(polygon) returns the set of vertices of a polygon, we could identify 
close-together vertices of polygons stored in a table with:

````
SELECT p1.id, p2.id, v1, v2
FROM polygons p1, polygons p2,
     LATERAL vertices(p1.poly) v1,
     LATERAL vertices(p2.poly) v2
WHERE (v1 <-> v2) < 10 AND p1.id != p2.id;
````
This query could also be written
````
SELECT p1.id, p2.id, v1, v2
FROM polygons p1 CROSS JOIN LATERAL vertices(p1.poly) v1,
     polygons p2 CROSS JOIN LATERAL vertices(p2.poly) v2
WHERE (v1 <-> v2) < 10 AND p1.id != p2.id;
````
函数调用，支持应用函数左边的ITEM（S）。所以可以看消除LATERAL，语义是一样的。  
(As already mentioned, the 
LATERAL key word is unnecessary in this example, but we use it for clarity.)

3、  
It is often particularly handy to LEFT JOIN to a LATERAL subquery, so that 
source rows will appear in the result even if the LATERAL subquery produces 
no rows for them. For example, if get_product_names() returns the names of 
products made by a manufacturer, but some manufacturers in our table currently
produce no products, we could find out which ones those are like this:
````
SELECT m.name
FROM manufacturers m LEFT JOIN LATERAL get_product_names(m.id) pname ON true
WHERE pname IS NULL;
````
lateral的查询结果也是可以作为整个语句的查询条件的

### 总结
1、lateral 可以出现在FROM的列表项层，也可以出现在JOIN数树中，如果出现在JOIN的右部分，那么
可以引用在JOIN左部分的任何项。  
2、由于lateral的计算步骤是从source table逐条展开的，所以OUTER JOIN时只能使用source table 
作为whole端，LATERAL内的ITEM不能作为WHOLE端。  
3、LATERAL 关键词可以在前缀一个 SELECT FROM 子项. 这能让 SELECT 子项在FROM项出现之前就引
用到FROM项中的列. (没有 LATERAL 的话, 每一个 SELECT 子项彼此都是独立的，因此不能够对其
它的 FROM 项进行交叉引用.)  
4、当一个 FROM 项包含 LATERAL 交叉引用的时候，查询的计算过程如下: 对于FROM项提供给交叉引
用列的每一行，或者多个FROM像提供给引用列的行的集合, LATERAL 项都会使用行或者行的集合的列
值来进行计算. 计算出来的结果集像往常一样被加入到联合查询之中. 这一过程会在列的来源表的行或
者行的集合上重复进行.

### 参考
【PostgreSQL 9.3 add LATERAL support - LATERAL的语法和用法介绍】https://github.com/digoal/blog/blob/master/201210/20121008_01.md?spm=a2c4e.10696291.0.0.408619a4cXorB6&file=20121008_01.md  
【LATERAL】https://www.postgresql.org/docs/devel/queries-table-expressions.html#QUERIES-LATERAL  





