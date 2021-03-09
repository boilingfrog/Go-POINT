## sync.Cond

### 前言

本次的代码是基于`go version go1.13.15 darwin/amd64`   

### 什么是sync.Cond

Go语言标准库中的条件变量`sync.Cond`，它可以让一组的`Goroutine`都在满足特定条件时被唤醒。  

每个`Cond`都会关联一个Lock`（*sync.Mutex or *sync.RWMutex）`  

```go
var (
	locker = new(sync.Mutex)
	cond   = sync.NewCond(locker)
)

func listen(x int) {
	// 获取锁
	cond.L.Lock()
	// 等待通知  暂时阻塞
	cond.Wait()
	fmt.Println(x)
	// 释放锁
	cond.L.Unlock()
}

func main() {
	// 启动60个被cond阻塞的线程
	for i := 1; i <= 60; i++ {
		go listen(i)
	}

	fmt.Println("start all")

	// 3秒之后 下发一个通知给已经获取锁的goroutine	time.Sleep(time.Second * 3)
	fmt.Println("++++++++++++++++++++one Signal")
	cond.Signal()

	// 3秒之后 下发一个通知给已经获取锁的goroutine
	time.Sleep(time.Second * 3)
	fmt.Println("++++++++++++++++++++one Signal")
	cond.Signal()

	// 3秒之后 下发广播给所有等待的goroutine
	time.Sleep(time.Second * 3)
	fmt.Println("++++++++++++++++++++begin broadcast")
	cond.Broadcast()
	// 阻塞直到所有的全部输出
	time.Sleep(time.Second * 60)
}
```

上面是个简单的例子，我们启动了60个线程，然后都被`cond`阻塞，主函数通过`Signal()`通知一个`goroutine`接触阻塞，通过`Broadcast()`通知所有被阻塞的全部解除阻塞。  

### 看下源码


```go
type Cond struct {
	// 用于保证结构体不会在编译期间拷贝
	noCopy noCopy
	// 锁
	L Locker
	// goroutine链表，维护等待唤醒的goroutine队列
	notify  notifyList
	// 保证运行期间不会发生copy
	checker copyChecker
}
``` 

重点分析下：`notifyList`和`copyChecker`  

- notify   

```go
type notifyList struct {
	// 正在等待的数量
	wait   uint32
	// 已经通知的数量
	notify uint32
	// 锁
	lock   uintptr
	// 指向链表头部
	head   unsafe.Pointer
	// 指向链表尾部
	tail   unsafe.Pointer
}
```

这个是核心，所有`wait`的`goroutine`都会被加入到这个链表中，然后在通知的时候再从这个链表中获取。  

- copyChecker  

保证运行期间不会发生copy

```go
type copyChecker uintptr
// copyChecker holds back pointer to itself to detect object copying
func (c *copyChecker) check() {
	if uintptr(*c) != uintptr(unsafe.Pointer(c)) &&
		!atomic.CompareAndSwapUintptr((*uintptr)(c), 0, uintptr(unsafe.Pointer(c))) &&
		uintptr(*c) != uintptr(unsafe.Pointer(c)) {
		panic("sync.Cond is copied")
	}
}
```

#### Wait  

```go
func (c *Cond) Wait() {
// 监测是否赋复制
	c.checker.check()
// 
	t := runtime_notifyListAdd(&c.notify)
	c.L.Unlock()
	runtime_notifyListWait(&c.notify, t)
	c.L.Lock()
}
```
