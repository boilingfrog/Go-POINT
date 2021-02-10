<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [context](#context)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [什么是context](#%E4%BB%80%E4%B9%88%E6%98%AFcontext)
    - [为什么需要context呢](#%E4%B8%BA%E4%BB%80%E4%B9%88%E9%9C%80%E8%A6%81context%E5%91%A2)
  - [context底层设计](#context%E5%BA%95%E5%B1%82%E8%AE%BE%E8%AE%A1)
    - [Context](#context)
  - [几种context](#%E5%87%A0%E7%A7%8Dcontext)
    - [emptyCtx](#emptyctx)
    - [cancelCtx](#cancelctx)
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

### 几种context

#### emptyCtx

context之源头  

```go
// An emptyCtx is never canceled, has no values, and has no deadline. It is not
// struct{}, since vars of this type must have distinct addresses.
type emptyCtx int

func (*emptyCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (*emptyCtx) Done() <-chan struct{} {
	return nil
}

func (*emptyCtx) Err() error {
	return nil
}

func (*emptyCtx) Value(key interface{}) interface{} {
	return nil
}

func (e *emptyCtx) String() string {
	switch e {
	case background:
		return "context.Background"
	case todo:
		return "context.TODO"
	}
	return "unknown empty Context"
}
```

`emptyCtx`以`do nothing`的方式实现了`Context`接口。  

同时有两个`emptyCtx`的全局变量  

```go
var (
	background = new(emptyCtx)
	todo       = new(emptyCtx)
)
```

通过下面两个导出的函数（首字母大写）对外公开：  

```go
func Background() Context {
	return background
}

func TODO() Context {
	return todo
}
```

这两个我们在使用的时候如何区分呢？  

先来看下官方的解释  

```go
// Background returns a non-nil, empty Context. It is never canceled, has no
// values, and has no deadline. It is typically used by the main function,
// initialization, and tests, and as the top-level Context for incoming
// requests.

// TODO returns a non-nil, empty Context. Code should use context.TODO when
// it's unclear which Context to use or it is not yet available (because the
// surrounding function has not yet been extended to accept a Context
// parameter).
```

`Background` 适用于主函数、初始化以及测试中，作为一个顶层的`context`。  

`TODO`适用于不知道传递什么`context`的情形。  

也就是在未考虑清楚是否传递、如何传递context时用`TODO`，作为发起点的时候用`Background`。  

#### cancelCtx

cancel机制的灵魂  

cancelCtx的cancel机制是手工取消、超时取消的内部实现  

```go
// A cancelCtx can be canceled. When canceled, it also cancels any children
// that implement canceler.
type cancelCtx struct {
	Context

	mu       sync.Mutex            // protects following fields
	done     chan struct{}         // created lazily, closed by first cancel call
	children map[canceler]struct{} // set to nil by the first cancel call
	err      error                 // set to non-nil by the first cancel call
}
```

看下`Done`  

```go
func (c *cancelCtx) Done() <-chan struct{} {
	// 加锁
	c.mu.Lock()
	// 如果done为空，创建make(chan struct{})
	if c.done == nil {
		c.done = make(chan struct{})
	}
	d := c.done
	c.mu.Unlock()
	return d
}
```

这是个懒汉模式的函数，第一次调用的时候`c.done`才会被创建。  


重点看下`cancel`  

```go
// 关闭 channel,c.done；递归地取消它的所有子节点；从父节点从删除自己。
func (c *cancelCtx) cancel(removeFromParent bool, err error) {
	if err == nil {
		panic("context: internal error: missing cancel error")
	}
	// 加锁
	c.mu.Lock()
	// 已经取消了
	if c.err != nil {
		c.mu.Unlock()
		return // already canceled
	}
	c.err = err
	// 关闭channel
	if c.done == nil {
		c.done = closedchan
	} else {
		close(c.done)
	}
	// 递归子节点，一层层取消
	for child := range c.children {
		// NOTE: acquiring the child's lock while holding parent's lock.
		child.cancel(false, err)
	}
	// 将子节点置空
	c.children = nil
	c.mu.Unlock()

	if removeFromParent {
		// 从父节点中移除自己 
		removeChild(c.Context, c)
	}
}

// 从父节点删除context
func removeChild(parent Context, child canceler) {
	p, ok := parentCancelCtx(parent)
	if !ok {
		return
	}
	p.mu.Lock()
	if p.children != nil {
		// 删除child
		delete(p.children, child)
	}
	p.mu.Unlock()
}
```

这个函数的作用就是关闭channel，递归地取消它的所有子节点；从父节点从删除自己。达到的效果是通过关闭channel，将取消信号传递给了它的所有子节点。    

对外暴露的`WithCancel`就是对`cancelCtx`的应用  

```go
// broadcastCancel安排在父级被取消时取消子级。
func propagateCancel(parent Context, child canceler) {
	// done为nil说明只读的
	if parent.Done() == nil {
		return // parent is never canceled
	}
	// 找到可以取消的父 context
	if p, ok := parentCancelCtx(parent); ok {
		p.mu.Lock()
		if p.err != nil {
			// 父节点已经取消了,子节点也取消
			child.cancel(false, p.err)
		} else {
			// 父节点未取消
			if p.children == nil {
				p.children = make(map[canceler]struct{})
			}
			p.children[child] = struct{}{}
		}
		p.mu.Unlock()
	} else {
		go func() {
			select {
			// 如果父节点取消了
			case <-parent.Done():
				child.cancel(false, parent.Err())
			case <-child.Done():
			}
		}()
	}
}
```


### 参考

【Go Context的踩坑经历】https://zhuanlan.zhihu.com/p/34417106  
【深度解密Go语言之context】https://www.cnblogs.com/qcrao-2018/p/11007503.html   
【深入理解Golang之context】https://juejin.cn/post/6844904070667321357  
【上下文 Context】https://draveness.me/golang/docs/part3-runtime/ch06-concurrency/golang-context/  
【Golang Context深入理解】https://juejin.cn/post/6844903555145400334  
