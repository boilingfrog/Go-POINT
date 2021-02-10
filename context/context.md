<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [context](#context)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [什么是context](#%E4%BB%80%E4%B9%88%E6%98%AFcontext)
    - [为什么需要context呢](#%E4%B8%BA%E4%BB%80%E4%B9%88%E9%9C%80%E8%A6%81context%E5%91%A2)
  - [context底层设计](#context%E5%BA%95%E5%B1%82%E8%AE%BE%E8%AE%A1)
    - [Context](#context)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## context

### 前言

之前浅读过，不过很快就忘记了，再次深入学习下。  

### 什么是context

go在Go 1.7 标准库引入context，主要用来在goroutine之间传递上下文信息，包括：取消信号、超时时间、截止时间、k-v 等。  

#### 为什么需要context呢

在并发程序中，由于超时、取消或者一些异常的情况，我们可能要抢占操作或者中断之后的操作。  

在Go里，我们不能直接杀死协程，协程的关闭一般会用 `channel+select` 方式来控制，但是如果某一个操作衍生出了很多的协程，并且相互关联。
或者某一个协程层级很深，有很深的子子协程，这时候使用`channel+select`就比较头疼了。  

所以context的适用机制  

- 上层任务取消后，所有的下层任务都会被取消；  
- 中间某一层的任务取消后，只会将当前任务的下层任务取消，而不会影响上层的任务以及同级任务。  

同时context也可以传值，不过这个很少用到，使用`context.Context`进行传递参数请求的所有参数一种非常差的设计，比较常见的使用场景是传递请求对应用户的认证令牌以及用于进行分布式追踪的请求ID。  

### context底层设计

#### Context

```go
type Context interface {
	// 返回context被取消的时间
	// 当没有设置Deadline时间，返回false
	Deadline() (deadline time.Time, ok bool)

	// 当context被关闭，返回一个被关闭的channel
	Done() <-chan struct{}

	// 在 channel Done 关闭后，返回 context 取消原因
	Err() error

	// 获取key对应的value
	Value(key interface{}) interface{}
}
```

`Deadline`返回Context被取消的时间,第一个返回式是截止时间，到了这个时间点，Context会自动发起取消请求；第二个返回值ok==false时表示没有设置截止时间，如果需要取消的话，需要调用取消函数进行取消。  

`Done`返回一个只读的channel,类型为struct{}，当该context被取消的时候，该channel会被关闭,同时对应的使用该context的routine也应该结束并返回。  

`Err`返回context结束的原因，只有在context被关闭的时候才会返回非空的值。  

1、如果context被取消，会返回Canceled错误；  
2、如果context超时，会返回DeadlineExceeded错误；  

`Value`获取之前存入key对应的value值。里面的值可以多次拿取。  






### 参考

【Go Context的踩坑经历】https://zhuanlan.zhihu.com/p/34417106  
【深度解密Go语言之context】https://www.cnblogs.com/qcrao-2018/p/11007503.html   
【深入理解Golang之context】https://juejin.cn/post/6844904070667321357  
【上下文 Context】https://draveness.me/golang/docs/part3-runtime/ch06-concurrency/golang-context/  
【Golang Context深入理解】https://juejin.cn/post/6844903555145400334  
