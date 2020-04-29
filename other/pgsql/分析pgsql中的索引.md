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
Lehman和Yao的论文中，修改了B树的结构，不管是内部节点还是叶子节点，都有一个指针指向兄弟节点。同时还引入了“High Key”（下述HK）用于描述当前子节点的最大值。  

![](https://img2020.cnblogs.com/blog/1237626/202004/1237626-20200429090003097-1404023179.png)

其中的k1就代表一个HK,其值是p0以及p0子节点的最大值。HK并不作为索引结构中的一个元组，只是标记了一个最大的范围。同理，对于上述的2n个节点，每个节点都存在一个指针指向右兄弟节点，Pi的子节点取值范围为（Ki-1，Ki]。  

然后来了解下：表和元组的组织方式  
PostgreSQL的索引结构，也是按照这种方式进行存储的。  

![](https://img2020.cnblogs.com/blog/1237626/202004/1237626-20200429093159100-190402728.png)

PageHeaderDate:是长度为20字节的页头数据，包括该文件块的一般信息，如：  

- 空闲空间的起始和结束位置  
- Special space的开始位置  
- 项指针的开始位置  
- 标志信息，如是否存在空闲项指针、是否所有的元组都可见  


Freespase是指未分配的空间（空闲空间）  

新插入的元组及其对应的Linp元素都会从Freespase空间来分配，Linp从Freespase头部开始分配，新元组是从尾部开始分配。  

Special space 是特殊空间：  

用于存放与索引相关的特定函数。由于索引文件的文件块结构和普通的文件的相同，因此Special space在不同表文件块中并没有使用，
其内容被置为空。







 

### 参考

【深入浅出PostgreSQL B-Tree索引结构】https://yq.aliyun.com/articles/53701   
【PostgreSQL内核分析——BTree索引】https://www.cnblogs.com/scu-cjx/p/9960483.html    