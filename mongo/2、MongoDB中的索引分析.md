<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [MongoDB 的索引](#mongodb-%E7%9A%84%E7%B4%A2%E5%BC%95)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## MongoDB 的索引

### 前言

MongoDB 在使用的过程中，一些频繁的查询条件我们会考虑添加索引。   

MongoDB 中支持多种索引  

1、单键索引：在单个字段上面建立索引；   

2、联合索引：联合索引支持在多个字段建立索引匹配查询；  

3、多键索引：对数组或者嵌套文档的字段进行索引；  

4、地理空间索引：对包含地理坐标的字段进行索引；   

5、文本索引：对文本字段进行全文索引；   

6、哈希索引：将字段值进行哈希处理后进行索引；   

7、通配符索引：使用通配符对任意字段进行索引；    

下面来看下 MongoDB 中的索引实现的底层逻辑。    

### MongoDB 使用 B 树还是 B+ 树索引

先来看下 B 树和 B+ 树的区别。   

B 树  和 B+ 树最重要的区别是 B+ 树只有叶子节点存储数据，其他节点用于索引，而 B 树 对于每个索引节点都有 `Data` 字段。

<img src="/img/mysql/mysql-btree"  alt="mysql" />     

B 树简单的讲就是一种多叉平衡查找树，它类似于普通的平衡二叉树。不同的是 B 树 允许每个节点有更多的子节点，这样就能大大减少树的高度。

<img src="/img/mysql/mysql-b+tree"  alt="mysql" />         

B 树 结构图中可以看到每个节点中不仅包含数据的 key 值，还有 data 值。而每一个页的存储空间是有限的，如果 data 数据较大时将会导致每个节点（即一个页）能存储的 key 的数量很小，当存储的数据量很大时同样会导致 B 树的深度较大，增大查询时的磁盘 I/O 次数，进而影响查询效率。  

在 B+Tree 中，所有数据记录节点都是按照键值大小顺序存放在同一层的叶子节点上，而非叶子节点上只存储 key 值信息，这样可以大大加大每个节点存储的 key 值数量，降低 B+Tree 的高度。

B+ 树相比与 B 树：

1、非叶子节点只存储索引信息；

2、所有叶子节点都有一个链指针，所以B+ 树可以进行范围查询；

3、数据都放在叶子节点中。   

那么 MongoDB 使用的是什么索引呢？在网上搜索会发现很多文章 MongoDB 用的是 B 树，这个答案是不准确的。   

MongoDB 官网中有一段描述写的是MongoDB索引使用 B-tree 数据结构。  

> Indexes are special data structures that store a small portion of the collection's data set in an easy-to-traverse form. MongoDB indexes use a B-tree data structure.

> The index stores the value of a specific field or set of fields, ordered by the value of the field. The ordering of the index entries supports efficient equality matches and range-based query operations. In addition, MongoDB can return sorted results using the ordering in the index.

大致意思就是 MongoDB 使用的是 B-tree 数据结构，支持等值匹配和范围查询。可以使用索引的排序返回排序的结果。   

> 在很地方我们会看到  B-tree， B-tree 树即 B 树。B 即 Balanced 平衡，因为 B 树的原英文名称为 B-tree，而国内很多人喜欢把 B-tree 译作 B-树，这是个非常不好的直译，很容易让人产生误解，人们可能会以为 B-树 和 B树 是两种树。

上面的 B 树 和 B+ 树的对比我们知道，B 树因为没有 B+ 中叶子节点的链指针，所以  B 树是不支持的范围查询的。   

MongoDB 官网中的介绍中明确的表示 MongoDB 支持范围查询，所以我们可以得出结论用的就是 B+ 树。官网中讲的 B 树，指广义上的 B 树，因为 B+ 树也是 B 树的变种也能称为 B 树。   

MongoDB 从 3.2 开始就默认使用 WiredTiger 作为存储引擎。   

> WiredTiger maintains a table's data in memory using a data structure called a B-Tree ( B+ Tree to be specific), referring to the nodes of a B-Tree as pages. Internal pages carry only keys. The leaf pages store both keys and values.

根据 WiredTiger 官方文档的描述 WiredTiger 就是 B+ 树，非叶子节点上只存储 key 值信息，叶子节点会存储 key 和 data 的数据。  

文档地址[WiredTiger Tuning page size and compression](https://source.wiredtiger.com/3.0.0/tune_page_size_and_comp.html)

所以可以得出结论 MongoDB 默认的存储引擎 WiredTiger 目前使用的是 B+ 树索引结构。   

### 单键索引

单键索引：只针对一个键添加索引，是最简单的索引类型。   

``
db.test_table.createIndex({ name: 1 },{background: true})
``

其中 1 指定升序创建索引，-1 表示指定降序创建索引。   

background

### 联合索引

MongoDB 中支持联合索引，联合索引就是将多个键值对组合到一起创建索引，也叫做复合索引。   

1、最左匹配原则   

MongoDB 中联合索引的使用和 MySQL 中联合索引的使用类似，也有最左匹配原则。即最左优先，在检索数据时从联合索引的最左边开始匹配。具体的最左匹配原则可参见[MySQL 联合索引](https://www.cnblogs.com/ricklz/p/17262747.html#%E8%81%94%E5%90%88%E7%B4%A2%E5%BC%95)   

联合索引创建的时候有一个一个基本的原则就是将选择性最强的列放到最前面。   

选择性最高值得是数据的重复值最少，因为区分度高的列能够很容易过滤掉很多的数据。组合索引中第一次能够过滤掉很多的数据，后面的索引查询的数据范围就小了很多了。    

2、遵循 ESR 规则

MongoDB 中联合索引的使用，对于索引的创建的顺序有一个原则就是  ESR 规则。   

就是联合索引的排序顺序从左到右为  

1、精确（Equal）匹配的字段放最前面；  

2、排序（Sort）条件放中间；  

3、范围（Range）匹配的字段放最后面。   

同样适用：ES, ER。   












### 参考

【MongoDB简介】https://docs.mongoing.com/mongo-introduction      
【MySQL 中的索引】https://www.cnblogs.com/ricklz/p/17262747.html   
【performance-best-practices-indexing】https://www.mongodb.com/blog/post/performance-best-practices-indexing  
【tune_page_size_and_comp】https://source.wiredtiger.com/3.0.0/tune_page_size_and_comp.html       