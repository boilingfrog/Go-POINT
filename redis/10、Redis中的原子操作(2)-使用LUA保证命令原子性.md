<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Redis 如何应对并发访问](#redis-%E5%A6%82%E4%BD%95%E5%BA%94%E5%AF%B9%E5%B9%B6%E5%8F%91%E8%AE%BF%E9%97%AE)
  - [使用 LUA 脚本](#%E4%BD%BF%E7%94%A8-lua-%E8%84%9A%E6%9C%AC)
    - [Redis 中如何使用 LUA 脚本](#redis-%E4%B8%AD%E5%A6%82%E4%BD%95%E4%BD%BF%E7%94%A8-lua-%E8%84%9A%E6%9C%AC)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Redis 如何应对并发访问

上个文章中，我们分析了Redis 中命令的执行是单线程的，虽然 Redis6.0 版本之后，引入了 I/O 多线程，但是对于 Redis 命令的还是单线程去执行的。所以如果业务中，我们只用 Redis 中的单命令去处理业务的话，命令的原子性是可以得到保障的。  

但是很多业务场景中，需要多个命令组合的使用，例如前面介绍的 `读取-修改-写回` 场景，这时候就不能保证组合命令的原子性了。所以这时候 LUA 就登场了。

### 使用 LUA 脚本

Redis 在 2.6 版本推出了 lua 脚本功能。  

引入 lua 脚本的优点：  

1、减少网络开销。可以将多个请求通过脚本的形式一次发送，减少网络时延。  

2、原子操作。Redis会将整个脚本作为一个整体执行，中间不会被其他请求插入。因此在脚本运行过程中无需担心会出现竞态条件，无需使用事务。  

3、复用。客户端发送的脚本会永久存在redis中，这样其他客户端可以复用这一脚本，而不需要使用代码完成相同的逻辑。  

关于 lua 的语法和 lua 是一门什么样的语言，可以自行 google。  

#### Redis 中如何使用 LUA 脚本

redis 中支持 LUA 脚本的几个命令  

````
EVAL
EVALSHA
SCRIPT LOAD - SCRIPT EXISTS
SCRIPT FLUSH
SCRIPT KILL
````


### 参考

【Redis核心技术与实战】https://time.geekbang.org/column/intro/100056701    
【Redis设计与实现】https://book.douban.com/subject/25900156/   
【字符串命令的实现】https://mcgrady-forever.github.io/2018/02/10/redis-analysis-t-string/     
【Redis 多线程网络模型全面揭秘】https://segmentfault.com/a/1190000039223696     
【高性能IO模型分析-Reactor模式和Proactor模式】https://zhuanlan.zhihu.com/p/95662364  
【什么是事件驱动架构？】https://www.redhat.com/zh/topics/integration/what-is-event-driven-architecture  
【事件驱动架构】https://help.aliyun.com/document_detail/207135.html   
【Comparing Two High-Performance I/O Design Patterns】https://www.artima.com/articles/comparing-two-high-performance-io-design-patterns  
【如何深刻理解Reactor和Proactor？】https://www.zhihu.com/question/26943938  
【Go netpoller 原生网络模型之源码全面揭秘】https://strikefreedom.top/go-netpoll-io-multiplexing-reactor  
【Redis中使用Lua脚本】https://zhuanlan.zhihu.com/p/77484377    
【Lua 是怎样一门语言？】https://www.zhihu.com/question/19841006  


