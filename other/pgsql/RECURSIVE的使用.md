## RECURSIVE

### 前言

简单的探究下CTE的执行步骤及使用方法，以及RECURSIVE的使用。

### CTE or WITH

WITH语句通常被称为通用表表达式（Common Table Expressions）或者CTEs。  

WITH语句作为一个辅助语句依附于主语句，WITH语句和主语句都可以是SELECT，INSERT，UPDATE，DELETE中的任何一种语句。  

举个栗子  
````sql
WITH result AS (
    SELECT d.user_id
    FROM documents d
    GROUP BY d.user_id
),info as(
    SELECT t.*,json_build_object('id', ur.id, 'name', ur.name) AS user_info
    FROM result t
    LEFT JOIN users ur on ur.id = t.user_id
    WHERE ur.id IS NOT NULL
)select * from info
````
定义了两个WITH辅助语句，result和info。result查询出符合要求的user信息，然后info对这个信息进行组装，组装出我们需要的
数据信息。

当然不用这个也是可以的，不过CTE主要的还是做数据的过滤。什么意思呢，我们可以定义多层级的CTE，然后一层层的查询过滤组装。最终筛选出我们
需要的数据，当然你可能会问为什么不一次性拿出所有的数据呢，当然如果数据很大，我们通过多层次的数据过滤组装，在效率上也更好。

### 在WITH中使用数据修改语句

WITH中可以不仅可以使用SELECT语句，同时还能使用DELETE，UPDATE，INSERT语句。因此，可以使用WITH，在一条SQL语句中进行不同的操作，如下例所示。

````sql
WITH moved_rows AS (
  DELETE FROM products
  WHERE
    "date" >= '2010-10-01'
  AND "date" < '2010-11-01'
  RETURNING *
)
INSERT INTO products_log
SELECT * FROM moved_rows;
````
本例通过WITH中的DELETE语句从products表中删除了一个月的数据，并通过RETURNING子句将删除的数据集赋给moved_rows这一CTE，最后在主语句中通过INSERT将删除的商品插入products_log中。

如果WITH里面使用的不是SELECT语句，并且没有通过RETURNING子句返回结果集，则主查询中不可以引用该CTE，但主查询和WITH语句仍然可以继续执行。这种情况可以实现将多个不相关的语句放在一个SQL语句里，实现了在不显式使用事务的情况下保证WITH语句和主语句的事务性，如下例所示。

````sql
WITH d AS (
  DELETE FROM foo
),
u as (
  UPDATE foo SET a = 1
  WHERE b = 2
)
DELETE FROM bar;
````




### 参考
【SQL优化（五） PostgreSQL （递归）CTE 通用表表达式】http://www.jasongj.com/sql/cte/  
