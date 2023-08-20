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