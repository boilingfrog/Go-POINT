<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [MySQL 中的分区表](#mysql-%E4%B8%AD%E7%9A%84%E5%88%86%E5%8C%BA%E8%A1%A8)
  - [InnoDB 逻辑存储结构](#innodb-%E9%80%BB%E8%BE%91%E5%AD%98%E5%82%A8%E7%BB%93%E6%9E%84)
    - [表空间 (Tablespace)](#%E8%A1%A8%E7%A9%BA%E9%97%B4-tablespace)
    - [段 (segment)](#%E6%AE%B5-segment)
    - [区 (extent)](#%E5%8C%BA-extent)
    - [页 (page)](#%E9%A1%B5-page)
    - [行 (row)](#%E8%A1%8C-row)
  - [InnoDB 数据页结构](#innodb-%E6%95%B0%E6%8D%AE%E9%A1%B5%E7%BB%93%E6%9E%84)
  - [分区别表的概念](#%E5%88%86%E5%8C%BA%E5%88%AB%E8%A1%A8%E7%9A%84%E6%A6%82%E5%BF%B5)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## MySQL 中的分区表

### InnoDB 逻辑存储结构

InnoDB 存储引擎的存储逻辑，所有数据都会被逻辑的保存在一个空间中，这个空间称为 InnoDB（Tablespace） 表空间，本质上是一个或多个磁盘文件组成的虚拟文件系统。InnoDB 表空间不仅仅存储了表和索引，它还保存了回滚日志`（redo log）`、插入缓冲`（insert buffer）`、双写缓冲`（doublewrite buffer）`以及其他内部数据结构。   

表空间，由段（segment），区（extend）,页（page）组层。页在一些文档中也被称为块（block）。   

<img src="/img/mysql/mysql-innodb-space.jpg"  alt="mysql" />   

#### 表空间 (Tablespace)

表空间可以看做是 InnoDB 存储引擎逻辑结构的最高层，所有的数据都存放在表空间中。   

在默认情况下，InnoDB 存储引擎有一个共享的表空间 ibdatal 即所有的数据都存放在这个表空间中。如果用户开启了 innodb_file_per_table，则每张表内的数据都可以单独的存放到一个表空间中。     

不过即使开启了 innodb_file_per_table，每张表内的存放的数据只是数据，索引和插入缓存 Bitmap 页，其它类的数据回滚信息，插入缓存页，系统事务信息，二次写缓冲，还是会放在原来的共享表空间中。这也能说明一个问题，及时开始了 innodb_file_per_table，贡献表空间的数据还是会不断的增大。   

#### 段 (segment)

表空间是由各个段组层的，常见的段由数据段，索引段，回滚段等。  

对于 InnoDB 存储引擎的存储数据存储结构 B+ 树，只有叶子节点存储数据，其它节点存储索引信息。也就是叶子节点存储的是数据段，非叶子节点存储的是索引段。   

在 InnoDB 存储引擎中，对段的管理都是由引擎自身所完成的，我们一般不能也没有必要对其进行操作。    

#### 区 (extent)

区是由连续页组成的空间，在任何情况下每个区的大小都是 1MB,为了保证区中页的连续性，InnoDB 存储引擎会一次性的从磁盘申请 4~5 个区。在默认情况下，InnoDB 存储引擎的大小为 16KB,即一个区中一共有 64 个连续的页。   

`InnoDB 1.0.x` 版本开始引入了压缩页，每个页的大小可以通过参数 KEY_BLOCK_SIZE 设置为 2K,4K,8K,因此每个区对应的页的数量就应该为 512,256,128。    

`InnoDB 1.0.x` 版本新增加了参数 innodb_page_size ,通过该参数可以将默认页的大小设置为 4k,8k,但是页中的数据库不会压缩，这时候区中页的数量同样是 256,128。不管页的大小怎么变化，区的大小总是1M.   

不过需要注意的是，用户启动了参数 innodb_file_per_table 设置了单独的表空间，创建的表默认大小是 96KB。不过一个区中有64个连续的页，创建的表至少是 1MB 才对？   

原因是每个段开始时，会先用32个页大小的碎片页来存放数据，在使用完，这些页之后才是64个连续页的申请，这样的目的，对于一些小表，或者 undo 这类的段，可以在开始申请较少的空间，节省磁盘容量的开销。  

#### 页 (page)

和大多数的数据库一样，InnoDB 有页 (page) 的概念(也称为块)，页是 InnoDB 磁盘管理的最小单位。在 InnoDB 存储引擎中，默认每个页的大小为 16KB。  

从 InnoDB 1.2.x 版本开始，可以通过 innodb_page_size 修改页的大小，可以设置为 4k,8k,16k。如果修改完成，所有表中页的大小都是 innodb_page_size，不能对其进行再次修改，除非使用 mysqldump 导入和导出操作产生的新的库。   

innoDB存储引擎中，常见的页类型有：

1、数据页（B-tree Node）；  

2、undo页（undo Log Page）；  

3、系统页 （System Page）；  

4、事物数据页 （Transaction System Page）；  

5、插入缓冲位图页（Insert Buffer Bitmap）；  

6、插入缓冲空闲列表页（Insert Buffer Free List）；  

7、未压缩的二进制大对象页（Uncompressed BLOB Page）；  

6、压缩的二进制大对象页 （compressed BLOB Page）。     

#### 行 (row)

InnoDB 存储引擎是按行进行存放的，每个页存放的数据是有硬性要求的，最多允许存放 `16KB/2-200`，即7992行记录。   

### InnoDB 数据页结构

页是 InnoDB 存储引擎管理数据库的最小单位，这里我们再来看下 InnoDB 数据页的内部结构。   

InnoDB 数据页由下面7部分组层  

1、File Header (文件头)；  

2、Page Header (页头)；  

3、Infimum 和 Supremum Records;  

4、User Records(用户记录，机行记录)；  

5、Free Space (空闲空间)；  

6、Page Directory（页目录）；  

7、File Trailer（文件结尾信息）。   

其中 `File Header`，`Page Header`，`File Trailer` 用来记录该页一些空间，大小是固定的，分别为38，56，8 字节。  

`User Records,Free Space,Page Directory` 这部分为实际的行记录存储空间，大小是动态的。   

<img src="/img/mysql/innodb-table-space.jpg"  alt="mysql" />   

1、`File Header`：用来记录页的一些头部信息；  

2、`Page Header`：用来记录数据页的状态信息；  

3、Infimum 和 Supremum Records：在 InnoDB 存引擎中，每个数据页由连个虚拟的行记录，用来限定记录的边界。Infimum 记录的是比页中任何主键都要小的值，Supremum 中记录的是页中最大值的边界。这两个值随着页的创建而创建并且永远都不会被删除。   

<img src="/img/mysql/innodb-infimum-supremun.jpg"  alt="mysql" />   



### 分区别表的概念

分区是将一个表的数据按照某个特定的键值，把不同的数据划分到不同的分区中。类似于数据库的分片处理

### 参考

【高性能MySQL(第3版)】https://book.douban.com/subject/23008813/    
【MySQL 实战 45 讲】https://time.geekbang.org/column/100020801  
【MySQL技术内幕】https://book.douban.com/subject/24708143/    
【MySQL学习笔记】https://github.com/boilingfrog/Go-POINT/tree/master/mysql    
【InnoDB 表空间】https://blog.csdn.net/u010647035/article/details/105009979  
