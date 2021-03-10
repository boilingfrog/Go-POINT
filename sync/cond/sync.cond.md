<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [sync.Cond](#synccond)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [什么是sync.Cond](#%E4%BB%80%E4%B9%88%E6%98%AFsynccond)
  - [看下源码](#%E7%9C%8B%E4%B8%8B%E6%BA%90%E7%A0%81)
    - [Wait](#wait)
    - [Signal](#signal)
    - [Broadcast](#broadcast)
  - [总结](#%E6%80%BB%E7%BB%93)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

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

<img src="/img/sync_cond_1.png" width = "446" height = "246" alt="sync_cond" align=center />

### 看下源码

```go
// Wait 原子式的 unlock c.L， 并暂停执行调用的 goroutine。
// 在稍后执行后，Wait 会在返回前 lock c.L. 与其他系统不同，
// 除非被 Broadcast 或 Signal 唤醒，否则等待无法返回。
//
// 因为等待第一次 resume 时 c.L 没有被锁定，所以当 Wait 返回时，
// 调用者通常不能认为条件为真。相反，调用者应该在循环中使用 Wait()：
//
//    c.L.Lock()
//    for !condition() {
//        c.Wait()
//    }
//    ... make use of condition ...
//    c.L.Unlock()
//
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
	// 总共需要等待的数量
	wait   uint32
	// 已经通知的数量
	notify uint32
	// 锁
	lock   uintptr
	// 指向链表头部
	head   *sudog
	// 指向链表尾部
	tail   *sudog
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
	// 监测是否复制
	c.checker.check()
	// 更新 notifyList中需要等待的wait的数量
	// 返回当前需要插入链表节点ticket
	t := runtime_notifyListAdd(&c.notify)
	c.L.Unlock()
	// 为当前的加入的waiter构建一个链表的节点，插入链表的尾部
	runtime_notifyListWait(&c.notify, t)
	c.L.Lock()
}

// go/src/runtime/sema.go
// 更新 notifyList中需要等待的wait的数量
// 同时返回当前的加入的 waiter 的 ticket 编号，从0开始  
//go:linkname notifyListAdd sync.runtime_notifyListAdd
func notifyListAdd(l *notifyList) uint32 {
    // 使用atomic原子的对wait字段进行加一操作
	return atomic.Xadd(&l.wait, 1) - 1
}

// go/src/runtime/sema.go
// 为当前的加入的waiter构建一个链表的节点，插入链表的尾部
//go:linkname notifyListWait sync.runtime_notifyListWait
func notifyListWait(l *notifyList, t uint32) {
	lock(&l.lock)

	// 当t小于notifyList中的notify，说明当前节点已经被通知了  
	if less(t, l.notify) {
		unlock(&l.lock)
		return
	}

	// 构建当前节点
	s := acquireSudog()
	s.g = getg()
	s.ticket = t
	s.releasetime = 0
	t0 := int64(0)
	if blockprofilerate > 0 {
		t0 = cputicks()
		s.releasetime = -1
	}
	// 头结点没构建，插入头结点
	if l.tail == nil {
		l.head = s
	} else {
		// 插入到尾节点
		l.tail.next = s
	}
	l.tail = s
	// 将当前goroutine置于等待状态并解锁
	// 通过调用goready（gp），可以使goroutine再次可运行。
	// 也就是将 M/P/G 解绑，并将 G 调整为等待状态，放入 sudog 等待队列中
	goparkunlock(&l.lock, waitReasonSyncCondWait, traceEvGoBlockCond, 3)
	if t0 != 0 {
		blockevent(s.releasetime-t0, 2)
	}
	releaseSudog(s)
}
```

梳理流程  

1、首先检测对象的复制行为，如果有复制发生直接抛出panic；  

2、然后调用`runtime_notifyListAdd`对`notifynotifyListList`中的`wait`(需要等待的数量)进行加一操作，同时返回一个`ticket`，用来作为当前`wait`的编号，这个编号，会和`notifyList`中的`notify`对应起来；  

3、然后调用`runtime_notifyListWait`把当前的`wait`封装成链表的一个节点，插入到`notifyList`维护的链表的尾部。  

<img src="/img/sync_cond_wait.png" width = "299" height = "314" alt="sync_cond" align=center />

#### Signal

```go
// 唤醒一个被wait的goroutine
func (c *Cond) Signal() {
	// 监测是否复制
	c.checker.check()
	runtime_notifyListNotifyOne(&c.notify)
}

// go/src/runtime/sema.go
// 通知链表中的第一个
//go:linkname notifyListNotifyOne sync.runtime_notifyListNotifyOne
func notifyListNotifyOne(l *notifyList) {
	// wait和notify，说明已经全部通知到了
	if atomic.Load(&l.wait) == atomic.Load(&l.notify) {
		return
	}

	lock(&l.lock)

	// 这里做了二次的确认
	// wait和notify，说明已经全部通知到了
	t := l.notify
	if t == atomic.Load(&l.wait) {
		unlock(&l.lock)
		return
	}

	// 原子的对notify执行+1操作
	atomic.Store(&l.notify, t+1)

	// 尝试找到需要被通知的 g
	// 如果目前还没来得及入队，是无法找到的
	// 但是，当它看到通知编号已经发生改变是不会被 park 的
	//
	// 这个查找过程看起来是线性复杂度，但实际上很快就停了
	// 因为 g 的队列与获取编号不同，因而队列中会出现少量重排，但我们希望找到靠前的 g
	// 而 g 只有在不再 race 后才会排在靠前的位置，因此这个迭代也不会太久，
	// 同时，即便找不到 g，这个情况也成立：
	// 它还没有休眠，并且已经失去了我们在队列上找到的（少数）其他 g 的 race。
	for p, s := (*sudog)(nil), l.head; s != nil; p, s = s, s.next {
		// 顺序拿到一个节点的ticket，会和上面会和notifyList中的notify做比较，相同才进行后续的操作
		// 这个我们分析了，notifyList中的notify和链表节点中的ticket是一一对应的  
		if s.ticket == t {
			n := s.next
			if p != nil {
				p.next = n
			} else {
				l.head = n
			}
			if n == nil {
				l.tail = p
			}
			unlock(&l.lock)
			s.next = nil
			// 通过goready掉起在上面通过goparkunlock挂起的goroutine
			readyWithTime(s, 4)
			return
		}
	}
	unlock(&l.lock)
}
```

梳理下流程：  

1、首先检测对象的复制行为，如果有复制发生直接抛出`panic`；  

2、判断`wait`和`notify`，如果两者相同说明已经已经全部通知到了；  

3、调用`notifyListNotifyOne`，通过for循环，依次遍历这个链表，直到找到和`notifyList`中的`notify`，相匹配的`ticket`的节点；  

4、掉起`goroutine`，完成通知。  

<img src="/img/sync_cond_sin.png" width = "299" height = "413" alt="sync_cond" align=center />

#### Broadcast

```go
// 唤醒所有被wait的goroutine
func (c *Cond) Broadcast() {
	c.checker.check()
	runtime_notifyListNotifyAll(&c.notify)
}

// go/src/runtime/sema.go
// notifyListNotifyAll notifies all entries in the list.
//go:linkname notifyListNotifyAll sync.runtime_notifyListNotifyAll
func notifyListNotifyAll(l *notifyList) {
	// wait和notify，说明已经全部通知到了
	if atomic.Load(&l.wait) == atomic.Load(&l.notify) {
		return
	}

	// 加锁
	lock(&l.lock)
	s := l.head
	l.head = nil
	l.tail = nil

	// 这个很粗暴，直接将notify的值置换成wait
	atomic.Store(&l.notify, atomic.Load(&l.wait))
	unlock(&l.lock)

	// 循环链表，一个个唤醒goroutine
	for s != nil {
		next := s.next
		s.next = nil
		readyWithTime(s, 4)
		s = next
	}
}
```

梳理下流程：  

1、首先检测对象的复制行为，如果有复制发生直接抛出panic；  

2、判断`wait`和`notify`，如果两者相同说明已经已经全部通知到了；  

3、`notifyListNotifyAll`，就相对简单了，直接将`notify`的值置为`wait`，标注这个已经全部通知了；  

4、循环链表，一个个唤醒`goroutine`。  

<img src="/img/sync_cond_pro.png" width = "408" height = "434" alt="sync_cond" align=center />

### 总结

`sync.Cond`不是一个常用的同步机制，但是在条件长时间无法满足时，与使用`for {}`进行忙碌等待相比,`sync.Cond`能够让出处理器的使用权，提供`CPU`的利用率。使用时我们也需要注意以下问题：  

1、`sync.Cond.Wait`在调用之前一定要使用获取互斥锁，否则会触发程序崩溃；  

2、`sync.Cond.Signal` 唤醒的 `Goroutine`都是队列最前面、等待最久的`Goroutine`；  

3、`sync.Cond.Broadcast`会按照一定顺序广播通知等待的全部 `Goroutine`。