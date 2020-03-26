## RECURSIVE

### 前言

简单的探究下CTE的执行步骤及使用方法，以及RECURSIVE的使用。

### CTE or WITH

WITH语句通常被称为通用表表达式（Common Table Expressions）或者CTEs。  

WITH语句作为一个辅助语句依附于主语句，WITH语句和主语句都可以是SELECT，INSERT，UPDATE，DELETE中的任何一种语句。  

举个栗子  
````go
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



### 参考
【SQL优化（五） PostgreSQL （递归）CTE 通用表表达式】http://www.jasongj.com/sql/cte/  
