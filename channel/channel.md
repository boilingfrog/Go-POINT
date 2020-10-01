## channel

### 前言

`channel`作为go中并发的一神器，深入研究下吧。  

### 设计的原理

在早期，CPU都是以单核的形式顺序执⾏机器指令。Go语⾔的祖先C语⾔正是这种顺序编程语⾔的代 表。顺序编程语⾔中的顺序是指：所有的指
令都是以串⾏的⽅式执⾏，在相同的时刻有且仅有⼀个 CPU在顺序执⾏程序的指令。

随着处理器技术的发展，单核时代以提升处理器频率来提⾼运⾏效率的⽅式遇到了瓶颈，⽬前各种主 流的CPU频率基本被锁定在了3GHZ附近。单核CPU的
发展的停滞，给多核CPU的发展带来了机遇。 相应地，编程语⾔也开始逐步向并⾏化的⽅向发展。Go语⾔正是在多核和⽹络化的时代背景下诞⽣的 原⽣
⽀持并发的编程语⾔。  

常⻅的并⾏编程有多种模型，主要有多线程、消息传递等。从理论上来看，多线程和基于消息的并发 编程是等价的。由于多线程并发模型可以⾃然对应
到多核的处理器，主流的操作系统因此也都提供了系统级的多线程⽀持，同时从概念上讲多线程似乎也更直观，因此多线程编程模型逐步被吸纳到主流
的编程语⾔特性或语⾔扩展库中。⽽主流编程语⾔对基于消息的并发编程模型⽀持则相⽐较少，Erlang语⾔是⽀持基于消息传递并发编程模型的代表者，
它的并发体之间不共享内存。Go语⾔是基于 消息并发模型的集⼤成者，它将基于CSP模型的并发编程内置到了语⾔中，通过⼀个go关键字就可以 轻易地
启动⼀个Goroutine，与Erlang不同的是Go语⾔的Goroutine之间是共享内存的。  

#### 共享内存

多线程共享内存。其实就是Java或者C++等语言中的多线程开发。单个的goutine代码是顺序执行，而并发编程时，创建多个goroutine，但我们并
不能确定不同的goroutine之间的执行顺序，多个goroutine之间大部分情况是代码交叉执行，在执行过程中，可能会修改或读取共享内存变量，这样
就会产生数据竞争,但是我们可以用锁去消除数据的竞争。  

当然这种在go中是不推荐的  

#### csp

Go语⾔最吸引⼈的地⽅是它内建的并发⽀持。Go语⾔并发体系的理论是C.A.R	Hoare在1978年提出的 CSP（Communicating	Sequential	Process，
通讯顺序进程）。CSP有着精确的数学模型，并实际应 ⽤在了Hoare参与设计的T9000通⽤计算机上。从NewSqueak、Alef、Limbo到现在的Go语⾔，对于
对 CSP有着20多年实战经验的Rob	Pike来说，他更关注的是将CSP应⽤在通⽤编程语⾔上产⽣的潜⼒。 作为Go并发编程核⼼的CSP理论的核⼼概念只有
⼀个：同步通信。  

⾸先要明确⼀个概念：并发不是并⾏。并发更关注的是程序的设计层⾯，并发的程序完全是可以顺序 执⾏的，只有在真正的多核CPU上才可能真正地同时运⾏。
并⾏更关注的是程序的运⾏层⾯，并⾏⼀ 般是简单的⼤量重复，例如GPU中对图像处理都会有⼤量的并⾏运算。为更好的编写并发程序，从设 计之初Go语⾔
就注重如何在编程语⾔层级上设计⼀个简洁安全⾼效的抽象模型，让程序员专注于分解 问题和组合⽅案，⽽且不⽤被线程管理和信号互斥这些繁琐的操作分
散精⼒。  

在并发编程中，对共享资源的正确访问需要精确的控制，在⽬前的绝⼤多数语⾔中，都是通过加锁等线程同步⽅案来解决这⼀困难问题，⽽Go语⾔却另辟蹊
径，它将共享的值通过Channel传递(实际上多 个独⽴执⾏的线程很少主动共享资源)。在任意给定的时刻，最好只有⼀个Goroutine能够拥有该资源。数
据竞争从设计层⾯上就被杜绝了。为了提倡这种思考⽅式，Go语⾔将其并发编程哲学化为⼀句⼝号：

> Do	not	communicate	by	sharing	memory;	instead,	share	memory	by	communicating.

不要通过共享内存来通信，⽽应通过通信来共享内存。  
这是更⾼层次的并发编程哲学(通过管道来传值是Go语⾔推荐的做法)。虽然像引⽤计数这类简单的并 发问题通过原⼦操作或互斥锁就能很好地实现，但是通
过Channel来控制访问能够让你写出更简洁正确的程序。  

### channel

Golang中使用 CSP中 channel 这个概念。channel 是被单独创建并且可以在进程之间传递，它的通信模式类似于 boss-worker 模式的，一个实体通
过将消息发送到channel 中，然后又监听这个 channel 的实体处理，两个实体之间是匿名的，这个就实现实体中间的解耦，其中 channel 是同步的一
个消息被发送到 channel 中，最终是一定要被另外的实体消费掉的。  

#### channel的定义

channel 是一个引用类型，所以在它被初始化之前，它的值是 nil，channel 使用 make 函数进行初始化。go中内置的类型，初始化的时候，我们需要初始化
channel的长度。  

指定了长度代表有缓冲
```go
ch := make(chan int, 1)
```

未指定就是无缓冲
```go
ch := make(chan int)
```

有缓冲和无缓冲的差别是什么呢？  

对不带缓冲的 channel 进行的操作实际上可以看作“同步模式”，带缓冲的则称为“异步模式”。  

同步模式下，发送方和接收方要同步就绪，只有在两者都 ready 的情况下，数据才能在两者间传输（后面会看到，实际上就是内存拷贝）。否则，任意一方
先行进行发送或接收操作，都会被挂起，等待另一方的出现才能被唤醒。  

异步模式下，在缓冲槽可用的情况下（有剩余容量），发送和接收操作都可以顺利进行。否则，操作的一方（如写入）同样会被挂起，直到出现相反操作
（如接收）才会被唤醒。  

举个栗子：  

无缓冲的  就是一个送信人去你家门口送信 ，你不在家 他不走，你一定要接下信，他才会走。  

无缓冲保证信能到你手上  

有缓冲的 就是一个送信人去你家仍到你家的信箱 转身就走 ，除非你的信箱满了 他必须等信箱空下来。  

有缓冲的 保证 信能进你家的邮箱  


### 源码剖析

```go
type hchan struct {
    qcount   uint           // buffer 中已放入的元素个数
    dataqsiz uint           // 用户构造 channel 时指定的 buf 大小,也就是底层循环数组的长度
    buf      unsafe.Pointer // 指向底层循环数组的指针 只针对有缓冲的 channel
    elemsize uint16         // buffer 中每个元素的大小
    closed   uint32         // channel 是否关闭，== 0 代表未 closed
    elemtype *_type         // channel 元素的类型信息
    sendx    uint           // 已发送元素在循环数组中的索引
    recvx    uint           // 已接收元素在循环数组中的索引
    recvq    waitq          // 等待接收的 goroutine  list of recv waiters
    sendq    waitq          // 等待发送的 goroutine list of send waiters

    lock mutex              // 保护 hchan 中所有字段
}
```

简单分析下:  

buf指向底层的循环数组，只有缓冲类型的channel才有。  

sendx，recvx 均指向底层循环数组，表示当前可以发送和接收的元素位置索引值（相对于底层数组）。  

sendq，recvq 分别表示被阻塞的 goroutine，这些 goroutine 由于尝试读取 channel 或向 channel 发送数据而被阻塞。读的时候，如果循环
数据为空，那么当前读的goroutine就会加入到recvq，等待有消息写入结束阻塞。同理写入的goroutine,一样，如果队列满了，就加入到sendq,阻塞直到
消息写入。   
  
waitq 相关的属性，可以理解为是一个 FIFO 的标准队列。其中 recvq 中是正在等待接收数据的 goroutine，sendq 中是等待发送数据的 goroutine。
waitq 使用双向链表实现。  

lock通过互斥锁保证数据安全。  

设计思路：  

对于无缓冲的是没有buf,有缓冲的buf是有buf的，长度也就是创建channel制定的长度。  

有缓冲channel的buf是循环使用的，已经读取过的，会被后面新写入的消息覆盖，通过sendx，recvx这两个指向底层数据的指针的滑动，实现对buf的复用。  

具体的消息写入读读取，以及goroutine的阻塞，请看下面  

### 创建

```go
func makechan(t *chantype, size int) *hchan {
	elem := t.elem

	// 做的一些检查
	if elem.size >= 1<<16 {
		throw("makechan: invalid channel element type")
	}
	if hchanSize%maxAlign != 0 || elem.align > maxAlign {
		throw("makechan: bad alignment")
	}

	mem, overflow := math.MulUintptr(elem.size, uintptr(size))
	if overflow || mem > maxAlloc-hchanSize || size < 0 {
		panic(plainError("makechan: size out of range"))
	}

	// Hchan does not contain pointers interesting for GC when elements stored in buf do not contain pointers.
	// buf points into the same allocation, elemtype is persistent.
	// SudoG's are referenced from their owning thread so they can't be collected.
	// TODO(dvyukov,rlh): Rethink when collector can move allocated objects.
	var c *hchan
	switch {
	case mem == 0:
		// 队列或元素大小为零
        // 当前 Channel 中不存在缓冲区，为 runtime.hchan 分配一段内存空间
		c = (*hchan)(mallocgc(hchanSize, nil, true))
		c.buf = c.raceaddr()
	case elem.ptrdata == 0:
		// 类型不是指针
		// 一次性给channel和buf（也就是底层数组）分类一块连续的空间
		c = (*hchan)(mallocgc(hchanSize+mem, nil, true))
		c.buf = add(unsafe.Pointer(c), hchanSize)
	default:
		// 默认情况下会单独为 runtime.hchan 和缓冲区分配内存
		c = new(hchan)
		c.buf = mallocgc(mem, elem, true)
	}

	// 最后更新几个字段的值
	c.elemsize = uint16(elem.size)
	c.elemtype = elem
	c.dataqsiz = uint(size)

	if debugChan {
		print("makechan: chan=", c, "; elemsize=", elem.size, "; elemalg=", elem.alg, "; dataqsiz=", size, "\n")
	}
	return c
}
```

### 写入数据

向一个channel写入数据的流程  

1、首先判断recvq等待接收队列是否为空，不为空说明缓冲区中没有内容或者是一个无缓冲channel。直接从recvq中取出一个goroutine，然后写入数据，接着唤醒
goroutine，最后结束发发送过程。   
2、如果缓冲区有空余的位置，写入数据到缓冲区，完成发送。  
3、如果缓冲区满了，那么就把写入数据的goroutine放到sendq中，进入睡眠，最后等待goroutine被唤醒。  

```go
/*
 * generic single channel send/recv
 * If block is not nil,
 * then the protocol will not
 * sleep but return if it could
 * not complete.
 *
 * sleep can wake up with g.param == nil
 * when a channel involved in the sleep has
 * been closed.  it is easiest to loop and re-run
 * the operation; we'll see that it's now closed.
 */
func chansend(c *hchan, ep unsafe.Pointer, block bool, callerpc uintptr) bool {
	if c == nil {
		if !block {
			return false
		}
		gopark(nil, nil, waitReasonChanSendNilChan, traceEvGoStop, 2)
		throw("unreachable")
	}

	if debugChan {
		print("chansend: chan=", c, "\n")
	}

	if raceenabled {
		racereadpc(c.raceaddr(), callerpc, funcPC(chansend))
	}

	// Fast path: check for failed non-blocking operation without acquiring the lock.
	//
	// After observing that the channel is not closed, we observe that the channel is
	// not ready for sending. Each of these observations is a single word-sized read
	// (first c.closed and second c.recvq.first or c.qcount depending on kind of channel).
	// Because a closed channel cannot transition from 'ready for sending' to
	// 'not ready for sending', even if the channel is closed between the two observations,
	// they imply a moment between the two when the channel was both not yet closed
	// and not ready for sending. We behave as if we observed the channel at that moment,
	// and report that the send cannot proceed.
	//
	// It is okay if the reads are reordered here: if we observe that the channel is not
	// ready for sending and then observe that it is not closed, that implies that the
	// channel wasn't closed during the first observation.
	if !block && c.closed == 0 && ((c.dataqsiz == 0 && c.recvq.first == nil) ||
		(c.dataqsiz > 0 && c.qcount == c.dataqsiz)) {
		return false
	}

	var t0 int64
	if blockprofilerate > 0 {
		t0 = cputicks()
	}

	lock(&c.lock)

	if c.closed != 0 {
		unlock(&c.lock)
		panic(plainError("send on closed channel"))
	}

	if sg := c.recvq.dequeue(); sg != nil {
		// Found a waiting receiver. We pass the value we want to send
		// directly to the receiver, bypassing the channel buffer (if any).
		send(c, sg, ep, func() { unlock(&c.lock) }, 3)
		return true
	}

	if c.qcount < c.dataqsiz {
		// Space is available in the channel buffer. Enqueue the element to send.
		qp := chanbuf(c, c.sendx)
		if raceenabled {
			raceacquire(qp)
			racerelease(qp)
		}
		typedmemmove(c.elemtype, qp, ep)
		c.sendx++
		if c.sendx == c.dataqsiz {
			c.sendx = 0
		}
		c.qcount++
		unlock(&c.lock)
		return true
	}

	if !block {
		unlock(&c.lock)
		return false
	}

	// Block on the channel. Some receiver will complete our operation for us.
	gp := getg()
	mysg := acquireSudog()
	mysg.releasetime = 0
	if t0 != 0 {
		mysg.releasetime = -1
	}
	// No stack splits between assigning elem and enqueuing mysg
	// on gp.waiting where copystack can find it.
	mysg.elem = ep
	mysg.waitlink = nil
	mysg.g = gp
	mysg.isSelect = false
	mysg.c = c
	gp.waiting = mysg
	gp.param = nil
	c.sendq.enqueue(mysg)
	goparkunlock(&c.lock, waitReasonChanSend, traceEvGoBlockSend, 3)
	// Ensure the value being sent is kept alive until the
	// receiver copies it out. The sudog has a pointer to the
	// stack object, but sudogs aren't considered as roots of the
	// stack tracer.
	KeepAlive(ep)

	// someone woke us up.
	if mysg != gp.waiting {
		throw("G waiting list is corrupted")
	}
	gp.waiting = nil
	if gp.param == nil {
		if c.closed == 0 {
			throw("chansend: spurious wakeup")
		}
		panic(plainError("send on closed channel"))
	}
	gp.param = nil
	if mysg.releasetime > 0 {
		blockevent(mysg.releasetime-t0, 2)
	}
	mysg.c = nil
	releaseSudog(mysg)
	return true
}
```














### 参考

【Go的CSP并发模型】https://www.jianshu.com/p/a3c9a05466e1  
【goroutine, channel 和 CSP】http://www.moye.me/2017/05/05/go-concurrency-patterns/  
【通过同步和加锁解决多线程的线程安全问题】https://blog.ailemon.me/2019/05/15/solving-multithreaded-thread-safety-problems-by-synchronization-and-locking/  
【golang channel 有缓冲 与 无缓冲 的重要区别】https://my.oschina.net/u/157514/blog/149192  
【Golang channel 源码深度剖析】https://www.cyhone.com/articles/analysis-of-golang-channel/   
【《Go专家编程》Go channel实现原理剖析】https://my.oschina.net/renhc/blog/2246871  