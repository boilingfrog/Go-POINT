<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [MySQL 中的索引](#mysql-%E4%B8%AD%E7%9A%84%E7%B4%A2%E5%BC%95)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [索引的实现](#%E7%B4%A2%E5%BC%95%E7%9A%84%E5%AE%9E%E7%8E%B0)
    - [哈希索引](#%E5%93%88%E5%B8%8C%E7%B4%A2%E5%BC%95)
    - [全文索引](#%E5%85%A8%E6%96%87%E7%B4%A2%E5%BC%95)
    - [B+ 树索引](#b-%E6%A0%91%E7%B4%A2%E5%BC%95)
  - [索引的分类](#%E7%B4%A2%E5%BC%95%E7%9A%84%E5%88%86%E7%B1%BB)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## MySQL 中的索引

### 前言  

上篇文章聊完了 MySQL 中的锁，这里接着来看下 MySQL 中的索引。  

一般当我们数据库中的某些查询比较慢的时候，正常情况下，一顿分析下来，大多数我们会考虑对这个查询加个索引，那么索引是如何工作的呢？为什么索引能加快查询的速度，下面来具体的分析下。   

在关系数据库中，索引是一种单独的、物理的对数据库表中一列或多列的值进行排序的一种存储结构，它是某个表中一列或若干列值的集合和相应的指向表中物理标识这些值的数据页的逻辑指针清单。  

索引的作用相当于图书的目录，可以根据目录中的页码快速找到所需的内容。  

### 索引的实现   

InnoDB 支持三种索引模型：  

1、哈希索引；  

2、全文索引；  

3、B+ 树索引。   

#### 哈希索引

哈希表也称为散列，是一种以键-值 (key-value) 存储数据的结构。输入查询的 key,就能找到对应的 value。哈希的思路很简单，把值放在数组里，用一个哈希函数把 key 换算成一个确定的位置，然后把 value 放在数组的这个位置。   

当然会存在哈希冲突，对个 key 在经过哈希算法处理后可能出现在哈希表中的同一个槽位，当出现哈希冲突的时候，可以使用链表来解决。这样发生冲突的数据都放入到链表中。在进行数据查询的时候，首先找到哈希表中的槽位，然后在链表中依次遍历找到对应的值。    

<img src="/img/mysql/mysql-hash.png"  alt="mysql" />  

哈希表的这种结构适合于等值查询的场景，在最优场景的下的时间复杂度能达到 `O(1)`。    

缺点也很明显，对于范围查询，这种结构就不能很好的支持了。    

#### 全文索引

全文索引就是将存储于数据库中的整本书或整篇文章中的任意内容信息找出来的技术，可以根据需求获取全文中的有关文章，节，段，句，词等信息，也能进行统计和分析。  

InnoDB 最早是不支持存储全文索引的，想要使用全文索引就要选用 MySIAM 存储引擎,从版本 `1.2.x` 开始才增加了全文索引支持。   

全文索引一般使用倒排索引(inverted index)实现，倒排索引和 B+ 树索引一样，也是一种索引结构。  

倒排索引在辅助表 (auxiliary table) 中存储了单词与单词自身在一个或多个文档中所在位置之间的映射。这样当全文索引，匹配到对应的单词就能知道对应的文档位置了。    

倒排索引是区别于正排索引的概念：  

正排索引：是以文档对象的唯一 ID 作为索引，以文档内容作为记录的结构。  

倒排索引：`Inverted index`，指的是将文档内容中的单词作为索引，将包含该词的文档 ID 作为记录的结构。    

倒排索引中文档和单词的关联关系，通常使用关联数据实现，主要有两种实现方式：  

1、`inverted file index`: 会记录单词和单词所在的文档 ID 之间的映射；  

2、`full inverted index`: 这种方式记录的更详细，除了会记录单词和单词所在的文档 ID 之间的映射，还会记录单词所在文档中的具体位置。   

下面来举个栗子  

<img src="/img/mysql/mysql-full-text-search-demo-1.jpg"  alt="mysql" />     

DocumentId 表示全文索引中文件的 id, Text 表示存储的内容。这些存储的文档就是全文索引需要查询的对象。     

`inverted file index`  

<img src="/img/mysql/mysql-full-text-search-demo-2.jpg"  alt="mysql" />    

可以看到关联中，记录了单词和 DocumentId 的映射，这样通过对应的单词就能找到，单词所在的文档，不用一个个遍历文档了。   

`full inverted index`

<img src="/img/mysql/mysql-full-text-search-demo-3.jpg"  alt="mysql" />  

这种方式和上面的 `inverted file index` 一样也会记录单词和文档的映射，只不过记录的更详细了，还记录了单词在文档中的具体位置。   

相比于 `inverted file index` 优点就是定位更精确，缺点也是很明显，需要用更多的空间。    

InnoDB 存储引擎支持全文索引采用 `full inverted index` 的方式，将`(DocumentId，Position)` 视为一个 `ilist`。  

因此在全文检索的表中，一共有两列，一列是 word 字段，另一个是 `ilist`，并且在 word 字段上设有索引。    

记录单词和 DocumentId 的映射关系的表称为 `Auxiliary Table`（辅助表）。   

辅助表是存在与磁盘上的持久化的表，由于磁盘 `I/O` 比较慢，因此提供 `FTS Index Cache`（全文检索索引缓存）来提高性能。`FTS Index Cache` 是一个红黑树结构，根据`（word, list）`排序，在有数据插入时，索引先更新到缓存中，而后 InnoDB 存储引擎会批量进行更新到辅助表中。    

#### B+ 树索引

B+ 树就是传统意义上的索引，这是目前关系型数据库中查找最为常用和最为有效的索引。B+ 树构造的索引类似于二叉树，根据键值快速 (Key Value) 快速找到数据。   

有一点需要注意的是，B+ 树索引并不能找到给定键值具体的行。B+ 树索引能找到的只是被查找数据行所在的页。然后把页读入到内存中，再在内存中查找，找到查询的目标数据。  

B+ 树是 B 树的变种，这里需要了解下 B 树。   

**为什么要引入 B 树或者 B+ 树呢？**   

红黑树等其它的数据结构也可以用来实现索引，为什么要使用 B 树或者 B+ 树，简单点讲就是为了减少磁盘的 I/O。     

一般来说，索引本身的数据量很大，全部放入到内存中是不太现实的，因此索引往往以索引文件的形式存储在磁盘中，磁盘 I/O 的消耗相比于内存中的读取还是大了很多的，在机械硬盘时代，从磁盘随机读一个数据块需要`10 ms`左右的寻址时间。     

为了让一个查询尽量少地读磁盘，就需要减少树的高度，就不能使用二叉树，而是使用 N 叉树了，这样就能在遍历较少节点的情况下也就是较少 I/O 查询的情况下找到目标值。    

所以 B 树和 B+ 树就被慢慢演变而来了。   

为什么用的是 B+ 树 而不是 B 树呢，这里来看下区别？   

`B-tree` 和 `B+` 树最重要的区别是 `B+` 树只有叶子节点存储数据，其他节点用于索引，而 `B-tree` 对于每个索引节点都有 `Data` 字段。   

<img src="/img/mysql/mysql-btree"  alt="mysql" />     

B 树简单的讲就是一种平衡查找树，使用自平衡二叉树就会存在，数据索引过多的情况下


<img src="/img/mysql/mysql-b+tree"  alt="mysql" />     




### 索引的分类

### 参考

【高性能MySQL(第3版)】https://book.douban.com/subject/23008813/    
【MySQL 实战 45 讲】https://time.geekbang.org/column/100020801  
【MySQL技术内幕】https://book.douban.com/subject/24708143/    
【MySQL学习笔记】https://github.com/boilingfrog/Go-POINT/tree/master/mysql    
【what-is-the-difference-between-mysql-innodb-b-tree-index-and-hash-index】https://medium.com/@mena.meseha/what-is-the-difference-between-mysql-innodb-b-tree-index-and-hash-index-ed8f2ce66d69  



