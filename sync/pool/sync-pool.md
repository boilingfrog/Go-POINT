## sync.pool

### sync.pool作用

有时候我们为了优化GC的场景，减少并复用内存，我们可以使用 sync.Pool 来复用需要频繁创建临时对象。 

`sync.Pool` 是一个临时对象池。一句话来概括，sync.Pool 管理了一组临时对象， 当需要时从池中获取，使用完毕后从再放回池中，以供他人使用。  

#### 使用

一个小demo

```go
func main() {
	// 初始化一个pool
	pool := &sync.Pool{
		// 默认的返回值设置，不写这个参数，默认是nil
		New: func() interface{} {
			return 0
		},
	}

	// 看一下初始的值，这里是返回0，如果不设置New函数，默认返回nil
	init := pool.Get()
	fmt.Println("初始值", init)

	// 设置一个参数1
	pool.Put(1)

	// 获取查看结果
	num := pool.Get()
	fmt.Println("put之后取值", num)

	// 再次获取，会发现，已经是空的了，只能返回默认的值。
	num = pool.Get()
	fmt.Println("put之后再次取值", num)
}
```

输出

```
初始值 0
put之后取值 1
put之后再次取值 0
```

### 源码解读

pool结构

```go
type Pool struct {
    // 用来标记，当前的 struct 是不能够被 copy 的
    noCopy noCopy
    // P 个固定大小的 poolLocal 数组，每个 P 拥有一个空间
    local     unsafe.Pointer // local fixed-size per-P pool, actual type is [P]poolLocal
    // 上面数组的大小，即 P 的个数
    localSize uintptr        // size of the local array

    // 同 local 和 localSize，只是在 gc 的过程中保留一次
    victim     unsafe.Pointer // local from previous cycle
    victimSize uintptr        // size of victims array

    // 自定义一个 New 函数，然后可以在 Get 不到东西时，自动创建一个
    New func() interface{}
}
```

```go
// unsafe.Sizeof(poolLocal{})  // 128 byte(1byte = 8 bits)
// unsafe.Sizeof(poolLocalInternal{})  // 32 byte(1byte = 8 bits)
type poolLocal struct {
	// 每个P对应的pool
	poolLocalInternal

	// 将 poolLocal 补齐至两个缓存行的倍数，防止 false sharing,
	// 每个缓存行具有 64 bytes，即 512 bit
	// 目前我们的处理器一般拥有 32 * 1024 / 64 = 512 条缓存行
	// 伪共享，仅占位用，防止在 cache line 上分配多个 poolLocalInternal
	pad [128 - unsafe.Sizeof(poolLocalInternal{})%128]byte
}

// Local per-P Pool appendix.
type poolLocalInternal struct {
	// private 存储一个 Put 的数据，pool.Put() 操作优先存入 private，如果private有信息，才会存入 shared
	private interface{} // Can be used only by the respective P.
	// 存储一个链表，用来维护 pool.Put() 操作加入的数据，每个 P 可以操作自己 shared 链表中的头部，而其他的 P 在用完自己的 shared 时，可能会来偷数据，从而操作链表的尾部
	shared  poolChain   // Local P can pushHead/popHead; any P can popTail.
}
```

#### GET 

当从池中获取对象时，会先从 per-P 的 poolLocal slice 中选取一个 poolLocal，选择策略遵循：  

1、优先从 private 中选择对象  
2、若取不到，则尝试从 shared 队列的队头进行读取  
3、若取不到，则尝试从其他的 P 中进行偷取 getSlow  
4、若还是取不到，则使用 New 方法新建  

```go
func (p *Pool) Get() interface{} {
	if race.Enabled {
		race.Disable()
	}
	// 获取一个 poolLocal
	l, pid := p.pin()
	// 先从 private 获取对象,有则立即返回
	x := l.private
	l.private = nil
	if x == nil {
		// 尝试从 localPool 的 shared 队列队头读取，
		// 因为队头的内存局部性比队尾更好。
		x, _ = l.shared.popHead()
		if x == nil {
			// 如果取不到，则获取新的缓存对象
			x = p.getSlow(pid)
		}
	}
	runtime_procUnpin()
	if race.Enabled {
		race.Enable()
		if x != nil {
			race.Acquire(poolRaceAddr(x))
		}
	}
	// 如果取不到，则获取新的缓存对象
	if x == nil && p.New != nil {
		x = p.New()
	}
	return x
}
```

#### pin

```go
// pin 会将当前的 goroutine 固定到 P 上，禁用抢占，并返回 localPool 池以及当前 P 的 pid。
func (p *Pool) pin() (*poolLocal, int) {
	// 返回当前 P.id
	pid := runtime_procPin()
	// In pinSlow we store to local and then to localSize, here we load in opposite order.
	// Since we've disabled preemption, GC cannot happen in between.
	// Thus here we must observe local at least as large localSize.
	// We can observe a newer/larger local, it is fine (we must observe its zero-initialized-ness).
	s := atomic.LoadUintptr(&p.localSize) // load-acquire
	l := p.local                          // load-consume
	// 因为可能存在动态的 P（运行时调整 P 的个数）procresize/GOMAXPROCS
	// 如果 P.id 没有越界，则直接返回
	if uintptr(pid) < s {
		return indexLocal(l, pid), pid
	}
	return p.pinSlow()
}
```
pin 的作用就是将当前 groutine 和 P 绑定在一起，禁止抢占。并且返回对应的 poolLocal 以及 P 的 id。  

pin() 首先会调用运行时实现获得当前 P 的 id，将 P 设置为禁止抢占，达到固定当前 goroutine 的目的。  

如果 G 被抢占，则 G 的状态从 running 变成 runnable，会被放回 P 的 localq 或 globaq，等待下一次调度。下次再执行时，就不一定是和现在的 P 相结合了。因为之后会用到 pid，如果被抢占了，有可能接下来使用的 pid 与所绑定的 P 并非同一个。  

```go
// src/runtime/proc.go

func procPin() int {
	_g_ := getg()
	mp := _g_.m

	mp.locks++
	return int(mp.p.ptr().id)
}
```

#### pinSlow

```go
var (
	allPoolsMu Mutex
	// allPools 是一组 pool 的集合，具有非空主缓存。
	// 有两种形式来保护它的读写：1. allPoolsMu 锁; 2. STW.
	allPools   []*Pool
)

func (p *Pool) pinSlow() (*poolLocal, int) {
	// 这时取消 P 的禁止抢占，因为使用 mutex 时候 P 必须可抢占
	runtime_procUnpin()
	// 加锁
	allPoolsMu.Lock()
	defer allPoolsMu.Unlock()
	// 当锁住后，再次固定 P 取其 id
	pid := runtime_procPin()
	// 因为 pinSlow 中途可能已经被其他的线程调用，因此这时候需要再次对 pid 进行检查。 如果 pid 在 p.local 大小范围内，则不用创建 poolLocal 切片，直接返回。
	s := p.localSize
	l := p.local
	if uintptr(pid) < s {
		return indexLocal(l, pid), pid
	}

	// 如果数组为空，新建
	// 将其添加到 allPools，垃圾回收器从这里获取所有 Pool 实例
	if p.local == nil {
		allPools = append(allPools, p)
	}
	// 根据 P 数量创建 slice，如果 GOMAXPROCS 在 GC 间发生变化
	// 我们重新分配此数组并丢弃旧的
	size := runtime.GOMAXPROCS(0)
	local := make([]poolLocal, size)
	// 将底层数组起始指针保存到 p.local，并设置 p.localSize
	atomic.StorePointer(&p.local, unsafe.Pointer(&local[0])) // store-release
	atomic.StoreUintptr(&p.localSize, uintptr(size))         // store-release
	return &local[pid], pid
}
```

pinSlow() 会首先取消 P 的禁止抢占，这是因为使用 mutex 时 P 必须为可抢占的状态。 然后使用 allPoolsMu 进行加锁。 当完成加锁后，再重新固定 P 
，取其 pid。注意，因为中途可能已经被其他的线程调用，因此这时候需要再次对 pid 进行检查。 如果 pid 在 p.local 大小范围内，则不再此时创建，直接返回。  

如果 p.local 为空，则将 p 扔给 allPools 并在垃圾回收阶段回收所有 Pool 实例。 最后再完成对 p.local 的创建（彻底丢弃旧数组）  

#### getSlow 

在取对象的过程中，本地p中的private没有，并且shared没有，这时候就会调用`Pool.getSlow()`，尝试从其他的p中获取。  

```go
func (p *Pool) getSlow(pid int) interface{} {
	// See the comment in pin regarding ordering of the loads.
	size := atomic.LoadUintptr(&p.localSize) // load-acquire
	locals := p.local                        // load-consume
	// 尝试熊其他的p中偷取
	for i := 0; i < int(size); i++ {
		l := indexLocal(locals, (pid+i+1)%int(size))
		if x, _ := l.shared.popTail(); x != nil {
			return x
		}
	}

	// 尝试从victim cache中取对象。这发生在尝试从其他 P 的 poolLocal 偷去失败后，
	// 因为这样可以使 victim 中的对象更容易被回收。
	size = atomic.LoadUintptr(&p.victimSize)
	if uintptr(pid) >= size {
		return nil
	}
	locals = p.victim
	l := indexLocal(locals, pid)
	if x := l.private; x != nil {
		l.private = nil
		return x
	}
	for i := 0; i < int(size); i++ {
		l := indexLocal(locals, (pid+i)%int(size))
		if x, _ := l.shared.popTail(); x != nil {
			return x
		}
	}

	// 清空 victim cache。下次就不用再从这里找了
	atomic.StoreUintptr(&p.victimSize, 0)

	return nil
}
```

从索引 pid+1 的 poolLocal 处开始，尝试调用 shared.popTail() 获取缓存对象。如果没有找到就从victim去找。  

如果没有找到，清空 victim cache。对于get来讲，如何其他的p也没找到，就new一个出来。   







#### Put

```go
// Put adds x to the pool.
func (p *Pool) Put(x interface{}) {
	if x == nil {
		return
	}
	if race.Enabled {
		if fastrand()%4 == 0 {
			// Randomly drop x on floor.
			return
		}
		race.ReleaseMerge(poolRaceAddr(x))
		race.Disable()
	}
	// 获得一个 localPool
	l, _ := p.pin()
	// 优先写入 private 变量
	if l.private == nil {
		l.private = x
		x = nil
	}
	// 如果 private 有值，则写入 shared poolChain 链表
	if x != nil {
		l.shared.pushHead(x)
	}
	runtime_procUnpin()
	if race.Enabled {
		race.Enable()
	}
}
```





go本身也用到了sync.pool，例如`fmt.Printf`


```
// Sprintf formats according to a format specifier and returns the resulting string.
func Sprintf(format string, a ...interface{}) string {
	p := newPrinter()
	p.doPrintf(format, a)
	s := string(p.buf)
	p.free()
	return s
}

var ppFree = sync.Pool{
	New: func() interface{} { return new(pp) },
}

// newPrinter allocates a new pp struct or grabs a cached one.
func newPrinter() *pp {
	p := ppFree.Get().(*pp)
	p.panicking = false
	p.erroring = false
	p.wrapErrs = false
	p.fmt.init(&p.buf)
	return p
}
```


### 参考
【深入Golang之sync.Pool详解】https://www.cnblogs.com/sunsky303/p/9706210.html  
【由浅入深聊聊Golang的sync.Pool】https://juejin.cn/post/6844903903046320136   
【Golang sync.Pool源码阅读与分析】https://jiajunhuang.com/articles/2020_05_05-go_sync_pool.md.html  
【【Go夜读】sync.Pool 源码阅读及适用场景分析】https://www.jianshu.com/p/f61bfe89e473  
【golang的对象池sync.pool源码解读】https://zhuanlan.zhihu.com/p/99710992  
【15.5 缓存池】https://golang.design/under-the-hood/zh-cn/part4lib/ch15sync/pool/  
【请问sync.Pool有什么缺点？】https://mp.weixin.qq.com/s?__biz=MzA4ODg0NDkzOA==&mid=2247487149&idx=1&sn=f38f2d72fd7112e19e97d5a2cd304430&source=41#wechat_redirect  
