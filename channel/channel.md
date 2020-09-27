## channel

### 前言

`channel`作为go中并发的一神器，深入研究下吧。  

### 设计的原理

go中的并发形式

#### 共享内存

多线程共享内存。其实就是Java或者C++等语言中的多线程开发。单个的goutine代码是顺序执行，而并发编程时，创建多个goroutine，但我们并
不能确定不同的goroutine之间的执行顺序，多个goroutine之间大部分情况是代码交叉执行，在执行过程中，可能会修改或读取共享内存变量，这样
就会产生数据竞争,但是我们可以用锁去消除数据的竞争。  

当然这种在go中是不推荐的  

#### csp

CSP 是`Communicating Sequential Process`的简称，中文可以叫做通信顺序进程，是一种并发编程模型，由`Tony Hoare`于 1977 年提出。
简单来说，CSP 模型由并发执行的实体（线程或者进程）所组成，实体之间通过发送消息进行通信，这里发送消息时使用的就是通道，或者叫 channel。
CSP 模型的关键是关注 channel，而不关注发送消息的实体。Go 语言实现了 CSP 部分理论。  

> Do not communicate by sharing memory; instead, share memory by communicating.

go中还是推荐使用csp的方式来处理并发的  

channel是go实现csp的关键  

### channel

Golang中使用 CSP中 channel 这个概念。channel 是被单独创建并且可以在进程之间传递，它的通信模式类似于 boss-worker 模式的，一个实体通
过将消息发送到channel 中，然后又监听这个 channel 的实体处理，两个实体之间是匿名的，这个就实现实体中间的解耦，其中 channel 是同步的一
个消息被发送到 channel 中，最终是一定要被另外的实体消费掉的。  










### 参考

【Go的CSP并发模型】https://www.jianshu.com/p/a3c9a05466e1  
【goroutine, channel 和 CSP】http://www.moye.me/2017/05/05/go-concurrency-patterns/  
【通过同步和加锁解决多线程的线程安全问题】https://blog.ailemon.me/2019/05/15/solving-multithreaded-thread-safety-problems-by-synchronization-and-locking/  