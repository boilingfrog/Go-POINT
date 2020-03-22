## pgsql中的索引的数据结构探究


### 前言

工欲善其器，必先利其器。对于pgsql的使用,我们肯定要用到索引，但是索引在pgsql中是怎么存在的呢、还是来探究下比较好，当我们使用的时候能做到知其然，知其所以然。

### pgsql中的几种索引

PostgreSQL提供了好几种索引类型：B-tree, Hash, GiST, SP-GiST和GIN 。

#### B-tree






### 参考
【索引类型】http://www.postgres.cn/docs/9.4/indexes-types.html   