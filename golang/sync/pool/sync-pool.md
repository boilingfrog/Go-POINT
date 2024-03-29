<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [sync.pool](#syncpool)
  - [sync.pool作用](#syncpool%E4%BD%9C%E7%94%A8)
    - [使用](#%E4%BD%BF%E7%94%A8)
  - [适用场景](#%E9%80%82%E7%94%A8%E5%9C%BA%E6%99%AF)
  - [案例](#%E6%A1%88%E4%BE%8B)
  - [源码解读](#%E6%BA%90%E7%A0%81%E8%A7%A3%E8%AF%BB)
    - [GET](#get)
    - [pin](#pin)
    - [pinSlow](#pinslow)
    - [getSlow](#getslow)
    - [Put](#put)
    - [poolChain](#poolchain)
    - [popHead](#pophead)
    - [pushHead](#pushhead)
    - [pack/unpack](#packunpack)
    - [popTail](#poptail)
  - [缓存的回收](#%E7%BC%93%E5%AD%98%E7%9A%84%E5%9B%9E%E6%94%B6)
  - [总结](#%E6%80%BB%E7%BB%93)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## sync.pool

### sync.pool作用

有时候我们为了优化GC的场景，减少并复用内存，我们可以使用 sync.Pool 来复用需要频繁创建临时对象。 

`sync.Pool` 是一个临时对象池。一句话来概括，sync.Pool 管理了一组临时对象， 当需要时从池中获取，使用完毕后从再放回池中，以供他人使用。  

本次探究的go版本`go version go1.13.15 darwin/amd64`   

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

### 适用场景

1、Pool 里对象的生命周期受 GC 影响，不适合于做连接池，因为连接池需要自己管理对象的生命周期  

放入本地池中的值有可能会在任何时候被删除，但是不通知调用者,也有可能被其他的goroutine偷走  

2、适用于储存一些会在goroutine间分享的临时对象。主要作用是减少GC，提高性能  

3、适用于已经申请了内存，目前未使用，接下来会使用的内存，来缓解GC     
 
4、一些生命周期比较短的不适合使用sync.pool来维护  

### 案例

go中的fmt就是使用到了sync.pool

```go

// Use simple []byte instead of bytes.Buffer to avoid large dependency.
type buffer []byte

func (b *buffer) write(p []byte) {
	*b = append(*b, p...)
}

func (b *buffer) writeString(s string) {
	*b = append(*b, s...)
}

func (b *buffer) writeByte(c byte) {
	*b = append(*b, c)
}

func (bp *buffer) writeRune(r rune) {
	if r < utf8.RuneSelf {
		*bp = append(*bp, byte(r))
		return
	}

	b := *bp
	n := len(b)
	for n+utf8.UTFMax > cap(b) {
		b = append(b, 0)
	}
	w := utf8.EncodeRune(b[n:n+utf8.UTFMax], r)
	*bp = b[:n+w]
}

// pp对象
// pp is used to store a printer's state and is reused with sync.Pool to avoid allocations.
type pp struct {
	buf buffer

	// arg holds the current item, as an interface{}.
	arg interface{}

	// value is used instead of arg for reflect values.
	value reflect.Value

	// fmt is used to format basic items such as integers or strings.
	fmt fmt

	// reordered records whether the format string used argument reordering.
	reordered bool
	// goodArgNum records whether the most recent reordering directive was valid.
	goodArgNum bool
	// panicking is set by catchPanic to avoid infinite panic, recover, panic, ... recursion.
	panicking bool
	// erroring is set when printing an error string to guard against calling handleMethods.
	erroring bool
	// wrapErrs is set when the format string may contain a %w verb.
	wrapErrs bool
	// wrappedErr records the target of the %w verb.
	wrappedErr error
}

// 定义sync.pool
var ppFree = sync.Pool{
	New: func() interface{} { return new(pp) },
}

// 在pool取出一个对象，然后初始化
func newPrinter() *pp {
	p := ppFree.Get().(*pp)
	p.panicking = false
	p.erroring = false
	p.wrapErrs = false
	p.fmt.init(&p.buf)
	return p
}

// Sprintln formats using the default formats for its operands and returns the resulting string.
// Spaces are always added between operands and a newline is appended.
func Sprintln(a ...interface{}) string {
	p := newPrinter()
	p.doPrintln(a)
	s := string(p.buf)
	p.free()
	return s
}

// 归还sync.pool
func (p *pp) free() {
	// Proper usage of a sync.Pool requires each entry to have approximately
	// the same memory cost. To obtain this property when the stored type
	// contains a variably-sized buffer, we add a hard limit on the maximum buffer
	// to place back in the pool.
	//
	// See https://golang.org/issue/23199
	if cap(p.buf) > 64<<10 {
		return
	}

	p.buf = p.buf[:0]
	p.arg = nil
	p.value = reflect.Value{}
	p.wrappedErr = nil
	ppFree.Put(p)
}
```

使用的时候在pool中取出一个pp，然后使用完毕之后，归还到的pool中。  

### 源码解读

**pool结构**

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
	// 本地 P 可以 pushHead/popHead；其他 P 则只能 popTail
	shared  poolChain   // Local P can pushHead/popHead; any P can popTail.
}
```

**noCopy**

意思就是不让copy,是如何实现的呢？

> Go中没有原生的禁止拷贝的方式，所以如果有的结构体，你希望使用者无法拷贝，只能指针传递保证全局唯一的话，可以这么干，定义 一个结构体叫 noCopy，实现如下的接口，然后嵌入到你想要禁止拷贝的结构体中，这样go vet就能检测出来。

```go
// noCopy may be embedded into structs which must not be copied
// after the first use.
//
// See https://golang.org/issues/8005#issuecomment-190753527
// for details.
type noCopy struct{}

// Lock is a no-op used by -copylocks checker from `go vet`.
func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}
```

测试下

```go
type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

type Person struct {
	noCopy noCopy
	name   string
}

// go中的函数传参都是值拷贝
func test(person Person) {
	fmt.Println(person)
}

func main() {
	var person Person
	test(person)
}
```

go vet main.go

```go
$ go vet main.go
# command-line-arguments
./main.go:18:18: test passes lock by value: command-line-arguments.Person contains command-line-arguments.noCopy
./main.go:19:14: call of fmt.Println copies lock value: command-line-arguments.Person contains command-line-arguments.noCopy
./main.go:24:7: call of test copies lock value: command-line-arguments.Person contains command-line-arguments.noCopy
```

使用vet检测到了不能copy的错误  

**伪共享**

```go
[128 - unsafe.Sizeof(poolLocalInternal{})%128]byte
```

这是来处理伪共享  

什么意思呢？  

我们知道从处理访问到内存，中间有好几级缓存，而这些缓存的存储单位是cacheline，也就是说每次从内存中加载数据，是以cacheline为单位的，这样就会存在一个问题，如果代码中的变量A和B被分配在了一个cacheline，但是处理器a要修改变量A，处理器b要修改B变量。此时，这个cacheline会被分别加载到a处理器的cache和b处理器的cache，当a修改A时，缓存系统会强制使b处理器的cacheline置为无效，同样当b要修改B时，会强制使得a处理器的cacheline失效，这样会导致cacheline来回无效，来回从低级的缓存加载数据，影响性能。  

增加一个 pad，补齐缓存行，让相关的字段能独立地加载到缓存行就不会出现 false sharding 了。  

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
    // private只保存了一个对象
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

<img src="/img/k8s/pool_1.png" width = "537" height = "990" alt="gc" align="center" />

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

绑定是通过procPin实现的   

```go
// src/runtime/proc.go

func procPin() int {
	_g_ := getg()
	mp := _g_.m

	mp.locks++
	return int(mp.p.ptr().id)
}
```

procPin函数实际上就是先获取当前goroutine，然后对当前协程绑定的线程（即为m）加锁，即mp.locks++，然后返回m目前绑定的p的id。  

系统线程对协程协调调度，会涉及到协程之间的抢占调度，有时候会抢占当前协程所属的P，原因就是不能让一个协程一直占用资源，抢占的时候如何判断是否可以抢占，
一个重要的条件就是判断`m.locks==0`。`procPin`所做的就是禁止当前P被强占。   

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

Put的逻辑就相对简单了

1、首先调用p.pin()抢占p  
2、优先放归还到private中  
3、如果private有值，则放到shared中  

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

总结下这个过程  

<img src="/img/k8s/pool_2.png" width = "353" height = "465" alt="gc" align="center" />

其实还有点疑惑，上面的victim是如何被保留下来的呢？来分析下pool的垃圾回收机制  

#### poolChain

```go
type poolChain struct {
	// 只有生产者会 push to，不用加锁
	head *poolChainElt

	// 读写需要原子控制。 
	tail *poolChainElt
}

type poolChainElt struct {
	poolDequeue

	// next 被 producer 写，consumer 读。所以只会从 nil 变成 non-nil
	// prev 被 consumer 写，producer 读。所以只会从 non-nil 变成 nil
	next, prev *poolChainElt
}

// poolDequeue is a lock-free fixed-size single-producer,
// multi-consumer queue. The single producer can both push and pop
// from the head, and consumers can pop from the tail.
//
// It has the added feature that it nils out unused slots to avoid
// unnecessary retention of objects. This is important for sync.Pool,
// but not typically a property considered in the literature.
type poolDequeue struct {
	// headTail packs together a 32-bit head index and a 32-bit
	// tail index. Both are indexes into vals modulo len(vals)-1.
	// headTail  包含一个 32 位的 head 和一个 32 位的 tail 指针。这两个值都和 len(vals)-1 取模过。
	//
	// tail = index of oldest data in queue 
	// 是队列中最老的数据
	// head = index of next slot to fill 
	// 指向下一个将要填充的 slot
	//
	// Slots in the range [tail, head) are owned by consumers. slots 
	// Slots的有效范围是 [tail, head)，由 consumers 持有
	// A consumer continues to own a slot outside this range until
	// it nils the slot, at which point ownership passes to the
	// producer.
	//
	// The head index is stored in the most-significant bits so
	// that we can atomically add to it and the overflow is
	// harmless.
	headTail uint64

	// vals is a ring buffer of interface{} values stored in this
	// dequeue. The size of this must be a power of 2.
	//
	// vals[i].typ is nil if the slot is empty and non-nil
	// otherwise. A slot is still in use until *both* the tail
	// index has moved beyond it and typ has been set to nil. This
	// is set to nil atomically by the consumer and read
	// atomically by the producer.
	// vals 是一个存储 interface{} 的环形队列，它的 size 必须是 2 的幂
	// 如果 slot 为空，则 vals[i].typ 为空；否则，非空。
	// 一个 slot 在这时宣告无效：tail 不指向它了，vals[i].typ 为 nil
	// 由 consumer 设置成 nil，由 producer 读
	vals []eface
}

type eface struct {
	typ, val unsafe.Pointer
}
```

`poolDequeue`被实现为单生产者、多消费者的固定大小的无锁（atomic 实现） Ring 式队列（底层存储使用数组，使用两个指针标记 head、tail）。生产者可以从 head 插入、head 删除，而消费者仅可从 tail 删除。  

headTail 指向队列的头和尾，通过位运算将 head 和 tail 存入 headTail 变量中。  


#### popHead

发生在从本地 shared 队列中消费并获取对象（消费者）  

```go
func (c *poolChain) popHead() (interface{}, bool) {
	d := c.head
	// d 是一个 poolDequeue，如果 d.popHead 是并发安全的，
	// 那么这里取 val 也是并发安全的。若 d.popHead 失败，则
	// 说明需要重新尝试。这个过程会持续到整个链表为空。
	for d != nil {
		if val, ok := d.popHead(); ok {
			return val, ok
		}
		d = loadPoolChainElt(&d.prev)
	}
	return nil, false
}
```

拿到头结点：c.head，是一个 poolDequeue。如果过头结点不为空调用，poolDequeue的popHead方法。  

```go
// popHead removes and returns the element at the head of the queue.
// It returns false if the queue is empty. It must only be called by a
// single producer.
func (d *poolDequeue) popHead() (interface{}, bool) {
	var slot *eface
	for {
		ptrs := atomic.LoadUint64(&d.headTail)
		head, tail := d.unpack(ptrs)
		if tail == head {
			// 队列是空的
			return nil, false
		}

		// head 位置是队头的前一个位置，所以此处要先退一位。
		// 在读出 slot 的 value 之前就把 head 值减 1，取消对这个 slot 的控制
		head--
		ptrs2 := d.pack(head, tail)
		if atomic.CompareAndSwapUint64(&d.headTail, ptrs, ptrs2) {
			// 拿到slot
			slot = &d.vals[head&uint32(len(d.vals)-1)]
			break
		}
	}

	val := *(*interface{})(unsafe.Pointer(slot))
	if val == dequeueNil(nil) {
		val = nil
	}
	// 重置 slot，typ 和 val 均为 nil
	// 这里清空的方式与 popTail 不同，与 pushHead 没有竞争关系，所以不用太小心
	*slot = eface{}
	return val, true
}
```

1、判断是否是空队列，通过unpack 函数分离出 head 和 tail 指针，如果 head 和 tail 相等，即首尾相等，那么这个队列就是空的，直接就返回 nil，false。  
2、队列不为空，就循环队列直到拿到slot。  
3、得到相应 slot 的元素后，经过类型转换并判断是否是 dequeueNil，如果是，说明没取到缓存的对象，返回 nil。  
4、返回val，清空slot。  

#### pushHead

发生在向本地 shared 队列中放置对象（生产者）  

```go
const (
	dequeueBits  = 32
	dequeueLimit = (1 << dequeueBits) / 4
)

func (c *poolChain) pushHead(val interface{}) {
	d := c.head
	// 如果链表为空，创建新的链表
	if d == nil {
		// 初始化链表
		const initSize = 8
		// 固定长度为 8，必须为 2 的指数
		d = new(poolChainElt)
		d.vals = make([]eface, initSize)
		c.head = d
		storePoolChainElt(&c.tail, d)
	}

	// 队列满了，pushHead就返回false
	if d.pushHead(val) {
		return
	}

	// 如果满了，就新创建一个是原来两倍大小的队列
	newSize := len(d.vals) * 2
	// 最大的限制
	if newSize >= dequeueLimit {
		// Can't make it any bigger.
		newSize = dequeueLimit
	}

	d2 := &poolChainElt{prev: d}
	d2.vals = make([]eface, newSize)
	c.head = d2
	storePoolChainElt(&d.next, d2)
	d2.pushHead(val)
}

// 将 val 添加到双端队列头部。如果队列已满，则返回 false。此函数只能被一个生产者调用
func (d *poolDequeue) pushHead(val interface{}) bool {
	ptrs := atomic.LoadUint64(&d.headTail)
	head, tail := d.unpack(ptrs)
	if (tail+uint32(len(d.vals)))&(1<<dequeueBits-1) == head {
		// 队列满了
		return false
	}
	slot := &d.vals[head&uint32(len(d.vals)-1)]

	// 检测这个 slot 是否被 popTail 释放
	typ := atomic.LoadPointer(&slot.typ)
	if typ != nil {
		// 另一个 goroutine 正在 popTail 这个 slot，说明队列仍然是满的
		return false
	}

	// The head slot is free, so we own it.
	if val == nil {
		val = dequeueNil(nil)
	}
	*(*interface{})(unsafe.Pointer(slot)) = val

	// Increment head. This passes ownership of slot to popTail
	// and acts as a store barrier for writing the slot.
	atomic.AddUint64(&d.headTail, 1<<dequeueBits)
	return true
}
```

#### pack/unpack

实现对head和tail` 的读写  

```go
// 将 head 和 tail 指针从 d.headTail 中分离开来
func (d *poolDequeue) unpack(ptrs uint64) (head, tail uint32) {
	const mask = 1<<dequeueBits - 1
	head = uint32((ptrs >> dequeueBits) & mask)
	tail = uint32(ptrs & mask)
	return
}
// 将 head 和 tail 指针打包到 d.headTail 一个 64bit 的变量中
func (d *poolDequeue) pack(head, tail uint32) uint64 {
	const mask = 1<<dequeueBits - 1
	return (uint64(head) << dequeueBits) |
		uint64(tail&mask)
}
```

#### popTail

```go
func (c *poolChain) popTail() (interface{}, bool) {
	d := loadPoolChainElt(&c.tail)
	if d == nil {
		return nil, false
	}

	for {
		// d可能是暂时为空，但如果next不为null并且popTail失败，则d为永久为空，这是唯一的天剑可以安全地将d从链中删除。
		d2 := loadPoolChainElt(&d.next)

		if val, ok := d.popTail(); ok {
			return val, ok
		}

		if d2 == nil {
			// 这是唯一的出队。它现在是空的，但是将来可能会被推倒。
			return nil, false
		}

		// 双向链表的尾节点里的双端队列被“掏空”，所以继续看下一个节点。
		// 并且由于尾节点已经被“掏空”，所以要甩掉它。这样，下次 popHead 就不会再看它有没有缓存对象了。
		if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&c.tail)), unsafe.Pointer(d), unsafe.Pointer(d2)) {
			// 甩掉尾节点
			storePoolChainElt(&d2.prev, nil)
		}
		d = d2
	}
}

func (d *poolDequeue) popTail() (interface{}, bool) {
	var slot *eface
	for {
		ptrs := atomic.LoadUint64(&d.headTail)
		head, tail := d.unpack(ptrs)
		// 队列满
		if tail == head {
			return nil, false
		}
		ptrs2 := d.pack(head, tail+1)
		if atomic.CompareAndSwapUint64(&d.headTail, ptrs, ptrs2) {
			slot = &d.vals[tail&uint32(len(d.vals)-1)]
			break
		}
	}

	val := *(*interface{})(unsafe.Pointer(slot))
	if val == dequeueNil(nil) {
		val = nil
	}

	// 注意：此处可能与 pushHead 发生竞争，解决方案是：
	// 1. 让 pushHead 先读取 typ 的值，如果 typ 值不为 nil，则说明 popTail 尚未清理完 slot
	// 2. 让 popTail 先清理掉 val 中的内容，在清理掉 typ，从而确保不会与 pushHead 对 slot 的写行为发生竞争
	slot.val = nil
	atomic.StorePointer(&slot.typ, nil)
	return val, true
}
```

### 缓存的回收

```go
// 将缓存清理函数注册到运行时 GC 时间段
func init() {
	runtime_registerPoolCleanup(poolCleanup)
}

// 由运行时实现
func runtime_registerPoolCleanup(cleanup func())
```

在 src/runtime/mgc.go 中:  

```go
// 开始 GC
func gcStart(trigger gcTrigger) {
	...
	clearpools()
	...
}

// 实现缓存清理
func clearpools() {
	// clear sync.Pools
	if poolcleanup != nil {
		poolcleanup()
	}
	...
}

var poolcleanup func()

// 利用编译器标志将 sync 包中的清理注册到运行时
//go:linkname sync_runtime_registerPoolCleanup sync.runtime_registerPoolCleanup
func sync_runtime_registerPoolCleanup(f func()) {
	poolcleanup = f
}
```

来看下具体的实现  

```go
var (
	// allPools 是一组 pool 的集合，具有非空主缓存。
	// 有两种形式来保护它的读写：1. allPoolsMu 锁; 2. STW.
	allPools   []*Pool

	// oldPools 是一组 pool 的集合，具有非空 victim 缓存。由 STW 保护
	oldPools []*Pool
)


func poolCleanup() {
	// 该函数会注册到运行时 GC 阶段(前)，此时为 STW 状态，不需要加锁
	// 它必须不处理分配且不调用任何运行时函数。

	// 由于此时是 STW，不存在用户态代码能尝试读取 localPool，进而所有的 P 都已固定（与 goroutine 绑定）

	// 从所有的 oldPols 中删除 victim
	for _, p := range oldPools {
		p.victim = nil
		p.victimSize = 0
	}

	// 将主缓存移动到 victim 缓存
	for _, p := range allPools {
		p.victim = p.local
		p.victimSize = p.localSize

		p.local = nil
		p.localSize = 0
	}

	// 具有非空主缓存的池现在具有非空的 victim 缓存，并且没有任何 pool 具有主缓存。
	oldPools, allPools = allPools, nil
}
```

`poolCleanup` 会在stw阶段的时候被调用。主要是删除oldPools中内容，然后将allPools里面的内容放入到victim中。最后标记allPools为oldPools。这样把pool里内容回收了，有victim进行兜底。   

### 总结

`sync.pool`拿取缓存的过程  

一个goroutine会抢占p，然后炒年糕当前p的private 中选择对象，如果里面没找到，然后尝试本地p的shared 队列的队头进行读取，若还是取不到，则尝试从其他 P 的 shared 队列队尾中偷取。 若偷不到，则尝试从上一个 GC 周期遗留到 victim 缓存中取，否则调用 New 创建一个新的对象。  

对于回收而言，池中所有临时对象在一次 GC 后会被放入 victim 缓存中， 而前一个周期被放入 victim 的缓存则会被清理掉。  

在加入 victim 机制前，sync.Pool 里对象的最⼤缓存时间是一个 GC 周期，当 GC 开始时，没有被引⽤的对象都会被清理掉；加入 victim 机制后，最大缓存时间为两个 GC 周期。  

当get一个对象使用完成之后，调用put归还的时候，需要注意将里面的内容清除  

Pool 不可以指定⼤⼩，⼤⼩只受制于 GC 临界值。

### 参考
【深入Golang之sync.Pool详解】https://www.cnblogs.com/sunsky303/p/9706210.html  
【由浅入深聊聊Golang的sync.Pool】https://juejin.cn/post/6844903903046320136   
【Golang sync.Pool源码阅读与分析】https://jiajunhuang.com/articles/2020_05_05-go_sync_pool.md.html  
【【Go夜读】sync.Pool 源码阅读及适用场景分析】https://www.jianshu.com/p/f61bfe89e473  
【golang的对象池sync.pool源码解读】https://zhuanlan.zhihu.com/p/99710992  
【15.5 缓存池】https://golang.design/under-the-hood/zh-cn/part4lib/ch15sync/pool/  
【请问sync.Pool有什么缺点？】https://mp.weixin.qq.com/s?__biz=MzA4ODg0NDkzOA==&mid=2247487149&idx=1&sn=f38f2d72fd7112e19e97d5a2cd304430&source=41#wechat_redirect  
【七分钟读懂 Go 的临时对象池pool及其应用场景】https://segmentfault.com/a/1190000016987629    
【伪共享(False Sharing)】http://ifeve.com/falsesharing/    
