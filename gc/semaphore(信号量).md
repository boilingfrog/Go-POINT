<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [运行时信号量机制 semaphore](#%E8%BF%90%E8%A1%8C%E6%97%B6%E4%BF%A1%E5%8F%B7%E9%87%8F%E6%9C%BA%E5%88%B6-semaphore)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [作用是什么](#%E4%BD%9C%E7%94%A8%E6%98%AF%E4%BB%80%E4%B9%88)
  - [几个主要的方法](#%E5%87%A0%E4%B8%AA%E4%B8%BB%E8%A6%81%E7%9A%84%E6%96%B9%E6%B3%95)
  - [如何实现](#%E5%A6%82%E4%BD%95%E5%AE%9E%E7%8E%B0)
  - [sudog 缓存](#sudog-%E7%BC%93%E5%AD%98)
    - [acquireSudog](#acquiresudog)
    - [releaseSudog](#releasesudog)
  - [semaphore](#semaphore)
    - [poll_runtime_Semacquire](#poll_runtime_semacquire)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 运行时信号量机制 semaphore

### 前言

最近在看源码，发现好多地方用到了这个`semaphore`。  

本文是在`go version go1.13.15 darwin/amd64`上进行的  

### 作用是什么  

下面是官方的描述  

```go
// Semaphore implementation exposed to Go.
// Intended use is provide a sleep and wakeup
// primitive that can be used in the contended case
// of other synchronization primitives.
// Thus it targets the same goal as Linux's futex,
// but it has much simpler semantics.
//
// That is, don't think of these as semaphores.
// Think of them as a way to implement sleep and wakeup
// such that every sleep is paired with a single wakeup,
// even if, due to races, the wakeup happens before the sleep.

// 具体的用法是提供 sleep 和 wakeup 原语
// 以使其能够在其它同步原语中的竞争情况下使用
// 因此这里的 semaphore 和 Linux 中的 futex 目标是一致的
// 只不过语义上更简单一些
//
// 也就是说，不要认为这些是信号量
// 把这里的东西看作 sleep 和 wakeup 实现的一种方式
// 每一个 sleep 都会和一个 wakeup 配对
// 即使在发生 race 时，wakeup 在 sleep 之前时也是如此  
```

上面提到了和`futex`作用一样，关于`futex`  

> futex（快速用户区互斥的简称）是一个在Linux上实现锁定和构建高级抽象锁如信号量和POSIX互斥的基本工具

> Futex 由一块能够被多个进程共享的内存空间（一个对齐后的整型变量）组成；这个整型变量的值能够通过汇编语言调用CPU提供的原子操作指令来增加或减少，并且一个进程可以等待直到那个值变成正数。Futex 的操作几乎全部在用户空间完成；只有当操作结果不一致从而需要仲裁时，才需要进入操作系统内核空间执行。这种机制允许使用 futex 的锁定原语有非常高的执行效率：由于绝大多数的操作并不需要在多个进程之间进行仲裁，所以绝大多数操作都可以在应用程序空间执行，而不需要使用（相对高代价的）内核系统调用。

go中的`semaphore`作用和`futex`目标一样，提供`sleep`和`wakeup`原语，使其能够在其它同步原语中的竞争情况下使用。当一个`goroutine`需要休眠时，将其进行集中存放，当需要`wakeup`时，再将其取出，重新放入调度器中。   

例如在读写锁的实现中，读锁和写锁之前的相互阻塞唤醒，就是通过`sleep`和`wakeup`实现，当有读锁存在的时候，新加入的写锁通过`semaphore`阻塞自己，当前面的读锁完成，在通过`semaphore`唤醒被阻塞的写锁。    

写锁

```go
// 获取互斥锁
// 阻塞等待所有读操作结束（如果有的话）
func (rw *RWMutex) Lock() {
	...
	// 原子的修改readerCount的值，直接将readerCount减去rwmutexMaxReaders
	// 说明，有写锁进来了，这在上面的读锁中也有体现
	r := atomic.AddInt32(&rw.readerCount, -rwmutexMaxReaders) + rwmutexMaxReaders
	// 当r不为0说明，当前写锁之前有读锁的存在
	// 修改下readerWait，也就是当前写锁需要等待的读锁的个数  
	if r != 0 && atomic.AddInt32(&rw.readerWait, r) != 0 {
		// 阻塞当前写锁
		runtime_SemacquireMutex(&rw.writerSem, false, 0)
	}
	...
}
```

通过`runtime_SemacquireMutex`对当前写锁进行`sleep`  

读锁释放  

```go
// 减少读操作计数，即readerCount--
// 唤醒等待写操作的协程（如果有的话）
func (rw *RWMutex) RUnlock() {
	...
	// 首先通过atomic的原子性使readerCount-1
	// 1.若readerCount大于0, 证明当前还有读锁, 直接结束本次操作
	// 2.若readerCount小于0, 证明已经没有读锁, 但是还有因为读锁被阻塞的写锁存在
	if r := atomic.AddInt32(&rw.readerCount, -1); r < 0 {
		// 尝试唤醒被阻塞的写锁
		rw.rUnlockSlow(r)
	}
	...
}

func (rw *RWMutex) rUnlockSlow(r int32) {
	...
	// readerWait--操作，如果readerWait--操作之后的值为0，说明，写锁之前，已经没有读锁了
	// 通过writerSem信号量，唤醒队列中第一个阻塞的写锁
	if atomic.AddInt32(&rw.readerWait, -1) == 0 {
		// 唤醒一个写锁
		runtime_Semrelease(&rw.writerSem, false, 1)
	}
}
```

写锁处理完之后，调用`runtime_Semrelease`来唤醒`sleep`的写锁

### 几个主要的方法

在`go/src/sync/runtime.go`中，定义了这几个方法  

```go
// Semacquire等待*s > 0，然后原子递减它。
// 它是一个简单的睡眠原语，用于同步
// library and不应该直接使用。
func runtime_Semacquire(s *uint32)

// SemacquireMutex类似于Semacquire,用来阻塞互斥的对象
// 如果lifo为true，waiter将会被插入到队列的头部
// skipframes是跟踪过程中要省略的帧数，从这里开始计算
// runtime_SemacquireMutex's caller.
func runtime_SemacquireMutex(s *uint32, lifo bool, skipframes int)

// Semrelease会自动增加*s并通知一个被Semacquire阻塞的等待的goroutine
// 它是一个简单的唤醒原语，用于同步
// library and不应该直接使用。
// 如果handoff为true, 传递信号到队列头部的waiter
// skipframes是跟踪过程中要省略的帧数，从这里开始计算
// runtime_Semrelease's caller.
func runtime_Semrelease(s *uint32, handoff bool, skipframes int)
```

具体的实现是在`go/src/runtime/sema.go`中

```go
//go:linkname sync_runtime_Semacquire sync.runtime_Semacquire
func sync_runtime_Semacquire(addr *uint32) {
	semacquire1(addr, false, semaBlockProfile, 0)
}

//go:linkname sync_runtime_Semrelease sync.runtime_Semrelease
func sync_runtime_Semrelease(addr *uint32, handoff bool, skipframes int) {
	semrelease1(addr, handoff, skipframes)
}

//go:linkname sync_runtime_SemacquireMutex sync.runtime_SemacquireMutex
func sync_runtime_SemacquireMutex(addr *uint32, lifo bool, skipframes int) {
	semacquire1(addr, lifo, semaBlockProfile|semaMutexProfile, skipframes)
}
```
### 如何实现

### sudog 缓存

`semaphore`的实现使用到了`sudog`，我们先来看下  

sudog 是运行时用来存放处于阻塞状态的`goroutine`的一个上层抽象，是用来实现用户态信号量的主要机制之一。 例如当一个`goroutine`因为等待`channel`的数据需要进行阻塞时，`sudog`会将`goroutine`及其用于等待数据的位置进行记录， 并进而串联成一个等待队列，或二叉平衡树。  

```go
// sudogs are allocated from a special pool. Use acquireSudog and
// releaseSudog to allocate and free them.
type sudog struct {
	// 以下字段受hchan保护
	g *g

	// isSelect 表示 g 正在参与一个 select, so
	// 因此 g.selectDone 必须以 CAS 的方式来获取wake-up race.
	isSelect bool
	next     *sudog
	prev     *sudog
	elem     unsafe.Pointer // 数据元素（可能指向栈）

	// 以下字段不会并发访问。
	// 对于通道，waitlink只被g访问。
	// 对于信号量，所有字段(包括上面的字段)
	// 只有当持有一个semroot锁时才被访问。
	acquiretime int64
	releasetime int64
	ticket      uint32
	parent      *sudog //semaRoot 二叉树
	waitlink    *sudog // g.waiting 列表或 semaRoot
	waittail    *sudog // semaRoot
	c           *hchan // channel
}
```

#### acquireSudog

`go/src/runtime/proc.go`

```go
//go:nosplit
func acquireSudog() *sudog {
	// Delicate dance: 信号量的实现调用acquireSudog，然后acquireSudog调用new(sudog)
	// new调用malloc, malloc调用垃圾收集器，垃圾收集器在stopTheWorld调用信号量
	// 通过在new(sudog)周围执行acquirem/releasem来打破循环
	// acquirem/releasem在new(sudog)期间增加m.locks，防止垃圾收集器被调用。

	// 获取当前 g 所在的 m
	mp := acquirem()
	// 获取p的指针
	pp := mp.p.ptr()
	if len(pp.sudogcache) == 0 {
		lock(&sched.sudoglock)
		// 首先，尝试从中央缓存获取一批数据。
		for len(pp.sudogcache) < cap(pp.sudogcache)/2 && sched.sudogcache != nil {
			s := sched.sudogcache
			sched.sudogcache = s.next
			s.next = nil
			pp.sudogcache = append(pp.sudogcache, s)
		}
		unlock(&sched.sudoglock)
		// 如果中央缓存中没有，新分配
		if len(pp.sudogcache) == 0 {
			pp.sudogcache = append(pp.sudogcache, new(sudog))
		}
	}
	// 取缓存中最后一个
	n := len(pp.sudogcache)
	s := pp.sudogcache[n-1]
	pp.sudogcache[n-1] = nil
	// 将刚取出的在缓存中移除
	pp.sudogcache = pp.sudogcache[:n-1]
	if s.elem != nil {
		throw("acquireSudog: found s.elem != nil in cache")
	}
	releasem(mp)
	return s
}
```

#### releaseSudog

```go
//go:nosplit
func releaseSudog(s *sudog) {
	if s.elem != nil {
		throw("runtime: sudog with non-nil elem")
	}
	if s.isSelect {
		throw("runtime: sudog with non-false isSelect")
	}
	if s.next != nil {
		throw("runtime: sudog with non-nil next")
	}
	if s.prev != nil {
		throw("runtime: sudog with non-nil prev")
	}
	if s.waitlink != nil {
		throw("runtime: sudog with non-nil waitlink")
	}
	if s.c != nil {
		throw("runtime: sudog with non-nil c")
	}
	gp := getg()
	if gp.param != nil {
		throw("runtime: releaseSudog with non-nil gp.param")
	}
	mp := acquirem() // avoid rescheduling to another P
	pp := mp.p.ptr()
	if len(pp.sudogcache) == cap(pp.sudogcache) {
		// Transfer half of local cache to the central cache.
		var first, last *sudog
		for len(pp.sudogcache) > cap(pp.sudogcache)/2 {
			n := len(pp.sudogcache)
			p := pp.sudogcache[n-1]
			pp.sudogcache[n-1] = nil
			pp.sudogcache = pp.sudogcache[:n-1]
			if first == nil {
				first = p
			} else {
				last.next = p
			}
			last = p
		}
		lock(&sched.sudoglock)
		last.next = sched.sudogcache
		sched.sudogcache = first
		unlock(&sched.sudoglock)
	}
	pp.sudogcache = append(pp.sudogcache, s)
	releasem(mp)
}
```


### semaphore

`go/src/runtime/sema.go`  

```go
// 用于sync.Mutex的异步信号量。

// semaRoot拥有一个具有不同地址（s.elem）的sudog平衡树。
// 每个sudog都可以依次（通过s.waitlink）指向一个列表，在相同地址上等待的其他sudog。
// 对具有相同地址的sudog内部列表进行的操作全部为O（1）。顶层semaRoot列表的扫描为O（log n），
// 其中，n是阻止goroutines的不同地址的数量，通过他们散列到给定的semaRoot。
type semaRoot struct {
	lock  mutex
	// waiters的平衡树的根节点
	treap *sudog
	// waiters的数量，读取的时候无所
	nwait uint32
}

// Prime to not correlate with any user patterns.
const semTabSize = 251

var semtable [semTabSize]struct {
	root semaRoot
	pad  [cpu.CacheLinePadSize - unsafe.Sizeof(semaRoot{})]byte
}
```

#### poll_runtime_Semacquire

```go
//go:linkname poll_runtime_Semacquire internal/poll.runtime_Semacquire
func poll_runtime_Semacquire(addr *uint32) {
	semacquire1(addr, false, semaBlockProfile, 0)
}

func semacquire1(addr *uint32, lifo bool, profile semaProfileFlags, skipframes int) {
    // 判断当前的
	gp := getg()
	if gp != gp.m.curg {
		throw("semacquire not on the G stack")
	}

	// *addr -= 1
	if cansemacquire(addr) {
		return
	}

	// 增加等待计数
	// 再试一次 cansemacquire 如果成功则直接返回
	// 将自己作为等待者入队
	// 休眠
	// (等待器描述符由出队信号产生出队行为)
	s := acquireSudog()
	root := semroot(addr)
	t0 := int64(0)
	s.releasetime = 0
	s.acquiretime = 0
	s.ticket = 0
	if profile&semaBlockProfile != 0 && blockprofilerate > 0 {
		t0 = cputicks()
		s.releasetime = -1
	}
	if profile&semaMutexProfile != 0 && mutexprofilerate > 0 {
		if t0 == 0 {
			t0 = cputicks()
		}
		s.acquiretime = t0
	}
	for {
		lock(&root.lock)
		// Add ourselves to nwait to disable "easy case" in semrelease.
		atomic.Xadd(&root.nwait, 1)
		// Check cansemacquire to avoid missed wakeup.
		if cansemacquire(addr) {
			atomic.Xadd(&root.nwait, -1)
			unlock(&root.lock)
			break
		}
		// Any semrelease after the cansemacquire knows we're waiting
		// (we set nwait above), so go to sleep.
		root.queue(addr, s, lifo)
		goparkunlock(&root.lock, waitReasonSemacquire, traceEvGoBlockSync, 4+skipframes)
		if s.ticket != 0 || cansemacquire(addr) {
			break
		}
	}
	if s.releasetime > 0 {
		blockevent(s.releasetime-t0, 3+skipframes)
	}
	releaseSudog(s)
}

func cansemacquire(addr *uint32) bool {
	for {
		v := atomic.Load(addr)
		if v == 0 {
			return false
		}
		if atomic.Cas(addr, v, v-1) {
			return true
		}
	}
}
```





### 参考

【同步原语】https://golang.design/under-the-hood/zh-cn/part2runtime/ch06sched/sync/  
【Go并发编程实战--信号量的使用方法和其实现原理】https://juejin.cn/post/6906677772479889422  
【Semaphore】https://github.com/cch123/golang-notes/blob/master/semaphore.md    