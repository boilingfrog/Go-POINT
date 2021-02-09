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

`Deadline`返回Context被取消的时间






### 参考

【Go Context的踩坑经历】https://zhuanlan.zhihu.com/p/34417106  
【深度解密Go语言之context】https://www.cnblogs.com/qcrao-2018/p/11007503.html   
【深入理解Golang之context】https://juejin.cn/post/6844904070667321357  
【上下文 Context】https://draveness.me/golang/docs/part3-runtime/ch06-concurrency/golang-context/  
