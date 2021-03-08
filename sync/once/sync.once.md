<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [sync.once](#synconce)
  - [sync.once的作用](#synconce%E7%9A%84%E4%BD%9C%E7%94%A8)
  - [实现原理](#%E5%AE%9E%E7%8E%B0%E5%8E%9F%E7%90%86)
  - [总结](#%E6%80%BB%E7%BB%93)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## sync.once

### sync.once的作用

根据名字就大致能猜到这个函数的作用，就是使用`sync.once`的对象只能执行一次。  

我们在`errgroup`就能看到它的身影  

```go
type Group struct {
	cancel func()

	wg sync.WaitGroup

	errOnce sync.Once
	err     error
}
```

他保证了，只会记录第一个出错的`goroutine`的错误信息  

### 实现原理

```go
// Once is an object that will perform exactly one action.
type Once struct {
	// 0未执行，1执行了
	done uint32
	// 互斥锁
	m    Mutex
}
```

里面就一个对外的函数

```go
func (o *Once) Do(f func()) {
	// 原子的读取done的值，如果为0代表onec第一次的执行还没有出发
	if atomic.LoadUint32(&o.done) == 0 {
		// 执行
		o.doSlow(f)
	}
}

func (o *Once) doSlow(f func()) {
	// 加锁
	o.m.Lock()
	defer o.m.Unlock()
	// 判断done变量为0表示还没执行第一次
	if o.done == 0 {
		// 机器数器原子的加一
		defer atomic.StoreUint32(&o.done, 1)
		// 执行传入的函数
		f()
	}
}
```

### 总结  

1、总体上也是很简单一个计数器，一把互斥锁，通过`atomic.LoadUint32`的原子读取技术器中的值；  

2、如果计数器中的值为0表示还没有执行；  

3、加锁，执行传入的函数，然后通过`atomic.StoreUint32`原子的对计数器的值进行加一操作；

4、完成。  


