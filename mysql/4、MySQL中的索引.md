<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [MySQL 中的索引](#mysql-%E4%B8%AD%E7%9A%84%E7%B4%A2%E5%BC%95)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [索引的实现](#%E7%B4%A2%E5%BC%95%E7%9A%84%E5%AE%9E%E7%8E%B0)
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

文档和单词的关联关系，通常使用关联数据实现，主要有两种实现方式：  

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





### 索引的分类

### 参考

【高性能MySQL(第3版)】https://book.douban.com/subject/23008813/    
【MySQL 实战 45 讲】https://time.geekbang.org/column/100020801  
【MySQL技术内幕】https://book.douban.com/subject/24708143/    
【MySQL学习笔记】https://github.com/boilingfrog/Go-POINT/tree/master/mysql    



