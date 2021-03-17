## 运行时信号量机制 semaphore

### 前言

最近在看源码，发现好多地方用到了这个`semaphore`。  

本文是在`go version go1.13.15 darwin/amd64`上进行的  

### 作用是什么  


  
### 如何使用

这个是go内部的包只能在内部使用  

`go/src/runtime/sema.go`  

`go/src/sync/runtime.go`   








### 参考

【同步原语】https://golang.design/under-the-hood/zh-cn/part2runtime/ch06sched/sync/  
【Go并发编程实战--信号量的使用方法和其实现原理】https://juejin.cn/post/6906677772479889422  
【Semaphore】https://github.com/cch123/golang-notes/blob/master/semaphore.md  