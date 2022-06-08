<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Redis 中的分布式锁如何使用](#redis-%E4%B8%AD%E7%9A%84%E5%88%86%E5%B8%83%E5%BC%8F%E9%94%81%E5%A6%82%E4%BD%95%E4%BD%BF%E7%94%A8)
  - [分布式锁的使用场景](#%E5%88%86%E5%B8%83%E5%BC%8F%E9%94%81%E7%9A%84%E4%BD%BF%E7%94%A8%E5%9C%BA%E6%99%AF)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Redis 中的分布式锁如何使用

### 分布式锁的使用场景

为了保证我们线上服务的并发性和安全性，目前我们的服务一般抛弃了单体应用，采用的都是扩展性很强的分布式架构。    

对于可变共享资源的访问，同一时刻，只能由一个线程或者进程去访问操作。这时候我们就需要做个标识，如果当前有线程或者进程在操作共享变量，我们就做个标记，标识当前资源正在被操作中， 其它的线程或者进程，就不能进行操作了。当前操作完成之后，删除标记，这样其他的线程或者进程，就能来申请共享变量的操作。通过上面的标记来保证同一时刻共享变量只能由一个线程或者进行持有。  

- 对于单体应用：多个线程之间访问可变共享变量，比较容易处理，可简单使用内存来存储标示即可；  

- 分布式应用：这种场景下比较麻烦，因为多个应用，部署的地址可能在不同的机房，一个在北京一个在上海。不能简单的存储标示在内存中了，这时候需要使用公共内存来记录该标示，栗如 Redis，MySQL 。。。   



### 参考

【Redis核心技术与实战】https://time.geekbang.org/column/intro/100056701    
【Redis设计与实现】https://book.douban.com/subject/25900156/   
【EVAL简介】http://www.redis.cn/commands/eval.html   
【Redis学习笔记】https://github.com/boilingfrog/Go-POINT/tree/master/redis
【Redis Lua脚本调试器】http://www.redis.cn/topics/ldb.html    
【redis中Lua脚本的使用】https://boilingfrog.github.io/2022/06/06/Redis%E4%B8%AD%E7%9A%84%E5%8E%9F%E5%AD%90%E6%93%8D%E4%BD%9C(2)-redis%E4%B8%AD%E4%BD%BF%E7%94%A8Lua%E8%84%9A%E6%9C%AC%E4%BF%9D%E8%AF%81%E5%91%BD%E4%BB%A4%E5%8E%9F%E5%AD%90%E6%80%A7/  


