## 分析了解pgsql中的索引  

### 前言

pgsql中索引的支持类型好像还是蛮多的，一一来分析下  

### 索引

PostgreSQL提供了多种索引类型： B-tree、Hash、GiST、SP-GiST 、GIN 和 BRIN。每一种索引类型使用了 一种不同的算法来适应不同类型的查询。

#### B-tree

B-tree可以在可排序数据上的处理等值和范围查询。  
例如下面的集中场景:  
````
<
<=
=
>=
>

````

##### 实现

pgsql中B-tree的实现是根据《Effiicient Locking for Concurrent Operations on B-Trees》论文设计实现的。  
Lehman和Yao的论文中，修改了B树的结构，不管是内部节点还是叶子节点，都有一个指针指向兄弟节点。  
![](https://img2020.cnblogs.com/blog/1237626/202004/1237626-20200429090003097-1404023179.png)





### 参考

【深入浅出PostgreSQL B-Tree索引结构】https://yq.aliyun.com/articles/53701 
【PostgreSQL内核分析——BTree索引】https://www.cnblogs.com/scu-cjx/p/9960483.html   