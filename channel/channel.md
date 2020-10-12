<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [channel](#channel)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [设计的原理](#%E8%AE%BE%E8%AE%A1%E7%9A%84%E5%8E%9F%E7%90%86)
    - [共享内存](#%E5%85%B1%E4%BA%AB%E5%86%85%E5%AD%98)
    - [csp](#csp)
  - [channel](#channel-1)
    - [channel的定义](#channel%E7%9A%84%E5%AE%9A%E4%B9%89)
  - [源码剖析](#%E6%BA%90%E7%A0%81%E5%89%96%E6%9E%90)
    - [环形队列](#%E7%8E%AF%E5%BD%A2%E9%98%9F%E5%88%97)
  - [创建](#%E5%88%9B%E5%BB%BA)
  - [写入数据](#%E5%86%99%E5%85%A5%E6%95%B0%E6%8D%AE)
  - [读取数据](#%E8%AF%BB%E5%8F%96%E6%95%B0%E6%8D%AE)
  - [channel的关闭](#channel%E7%9A%84%E5%85%B3%E9%97%AD)
  - [channel的关闭](#channel%E7%9A%84%E5%85%B3%E9%97%AD-1)
    - [M个receivers，一个sender](#m%E4%B8%AAreceivers%E4%B8%80%E4%B8%AAsender)
    - [一个receiver，N个sender](#%E4%B8%80%E4%B8%AAreceivern%E4%B8%AAsender)
    - [M个receiver，N个sender](#m%E4%B8%AAreceivern%E4%B8%AAsender)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

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

recvq和sendq，它们是 waitq 结构体，而waitq实际上就是一个双向链表，链表的元素是sudog，里面包含 g 字段，g 表示一个 goroutine，所以 sudog 可以看成一个 goroutine。但是两个还是有区别的。    

lock通过互斥锁保证数据安全。  

设计思路：  

对于无缓冲的是没有buf,有缓冲的buf是有buf的，长度也就是创建channel制定的长度。  

有缓冲channel的buf是循环使用的，已经读取过的，会被后面新写入的消息覆盖，通过sendx，recvx这两个指向底层数据的指针的滑动，实现对buf的复用。  

具体的消息写入读读取，以及goroutine的阻塞，请看下面  

#### 环形队列  

chan内部实现了一个环形队列作为其缓冲区，队列的长度是创建chan时指定的。  

看下实现的图片：  

![channel](/img/channel.jpg?raw=true)

- dataqsiz指示了队列长度为6，即可缓存6个元素；
- buf指向队列的内存，队列中还剩余两个元素；
- qcount表示队列中还有两个元素,也就是[1,3}；
- sendx指示后续写入的数据存储的位置，取值[0, 6)；
- recvx指示从该位置读取数据, 取值[0, 6)；

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
func chansend(c *hchan, ep unsafe.Pointer, block bool, callerpc uintptr) bool {
	// 如果 channel 是 nil
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
	// 对于不阻塞的 send，快速检测失败场景
	//
	// 如果 channel 未关闭且 channel 没有多余的缓冲空间。这可能是：
	// 1. channel 是非缓冲型的，且等待接收队列里没有 goroutine
	// 2. channel 是缓冲型的，但循环数组已经装满了元素
	if !block && c.closed == 0 && ((c.dataqsiz == 0 && c.recvq.first == nil) ||
		(c.dataqsiz > 0 && c.qcount == c.dataqsiz)) {
		return false
	}

	var t0 int64
	if blockprofilerate > 0 {
		t0 = cputicks()
	}
	// 加锁
	lock(&c.lock)

	// 如果channel关闭了
	if c.closed != 0 {
		unlock(&c.lock)
		panic(plainError("send on closed channel"))
	}

	// 如果接收队列里有 goroutine，直接将要发送的数据拷贝到接收 goroutine
	if sg := c.recvq.dequeue(); sg != nil {
		// Found a waiting receiver. We pass the value we want to send
		// directly to the receiver, bypassing the channel buffer (if any).
		send(c, sg, ep, func() { unlock(&c.lock) }, 3)
		return true
	}

	// 缓冲型的channel，buffer 中已放入的元素个数小于循环数组的长度
	if c.qcount < c.dataqsiz {
		// qp 指向 buf 的 sendx 位置
		qp := chanbuf(c, c.sendx)
		if raceenabled {
			raceacquire(qp)
			racerelease(qp)
		}
		// 将数据从 ep 处拷贝到 qp
		typedmemmove(c.elemtype, qp, ep)
		// 发送的游标加1
		c.sendx++
		// 如果发送的游标值等于容量值，游标值归0
		if c.sendx == c.dataqsiz {
			c.sendx = 0
		}
		// 缓冲区的数量加1
		c.qcount++
		// 解锁
		unlock(&c.lock)
		return true
	}
	// buff空间已经满了
	// 如果不需要阻塞，则直接返回错误
	if !block {
		unlock(&c.lock)
		return false
	}

	// 否则，阻塞该 goroutine.
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
	// 将该 goroutine 的结构放入 sendq 队列
	c.sendq.enqueue(mysg)
	// 休眠
	// 等待 goready 唤醒
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

简单的流程图如下：  

![channel](/img/channel_send1.png?raw=true)

### 读取数据

从一个channel读取数据的流程如下：  

1、如果等待发送的`goroutine list`,也就是sendq不为空。并且没有缓存区。直接从sendq中取出一个goroutine，读出当前goroutine中的消息，唤醒goroutine,
结束读取的过程。  
2、如果等待发送的`goroutine list`,也就是sendq不为空。说明缓冲区已经满了，移动recvx指针的位置，取出一个数据。同时在sendq中取出一个goroutine,
读取里面的数据到buf中，结束当前读取。  
3、如果等待发送的`goroutine list`,也就是sendq为空。并且缓冲区，有数据。直接在缓冲区取出数据，完成本次读取。  
4、如果等待发送的`goroutine list`,也就是sendq为空。并且缓冲区，没有数据。将当前goroutine加入recvq，进入睡眠，等待被写goroutine唤醒。    

```go
// chanrecv在通道c上接收并将接收到的数据写入ep。
// 如果 ep 是 nil，说明忽略了接收值。
// 如果 block == false，即非阻塞型接收，在没有数据可接收的情况下，返回 (false, false)
// 否则，如果 c 处于关闭状态，将 ep 指向的地址清零，返回 (true, false)
// 否则，用返回值填充 ep 指向的内存地址。返回 (true, true)
// 如果 ep 非空，则应该指向堆或者函数调用者的栈
func chanrecv(c *hchan, ep unsafe.Pointer, block bool) (selected, received bool) {

	if debugChan {
		print("chanrecv: chan=", c, "\n")
	}
	// 如果channel为nil
	if c == nil {
		// block == false，即非阻塞型接收，在没有数据可接收的情况下，返回 (false, false)
		if !block {
			return
		}
		// 接收一个 nil 的 channel，goroutine 挂起
		gopark(nil, nil, waitReasonChanReceiveNilChan, traceEvGoStop, 2)
		throw("unreachable")
	}

	// 在非阻塞模式下，快速检测到失败，不用获取锁，快速返回
	// 当我们观察到 channel 没准备好接收：
	// 1. 非缓冲型，等待发送列队 sendq 里没有 goroutine 在等待
	// 2. 缓冲型，但 buf 里没有元素
	// 之后，又观察到 closed == 0，即 channel 未关闭。
	// 因为 channel 不可能被重复打开，所以前一个观测的时候 channel 也是未关闭的，
	// 因此在这种情况下可以直接宣布接收失败，返回 (false, false)
	if !block && (c.dataqsiz == 0 && c.sendq.first == nil ||
		c.dataqsiz > 0 && atomic.Loaduint(&c.qcount) == 0) &&
		atomic.Load(&c.closed) == 0 {
		return
	}

	var t0 int64
	if blockprofilerate > 0 {
		t0 = cputicks()
	}

	// 加锁
	lock(&c.lock)

	// channel 已关闭，并且循环数组 buf 里没有元素
	// 这里可以处理非缓冲型关闭 和 缓冲型关闭但 buf 无元素的情况
	// 也就是说即使是关闭状态，但在缓冲型的 channel，
	// buf 里有元素的情况下还能接收到元素
	if c.closed != 0 && c.qcount == 0 {
		if raceenabled {
			raceacquire(c.raceaddr())
		}
		// 解锁
		unlock(&c.lock)
		if ep != nil {
			// 从一个已关闭的 channel 执行接收操作，且未忽略返回值
			// 那么接收的值将是一个该类型的零值
			// typedmemclr 根据类型清理相应地址的内存
			typedmemclr(c.elemtype, ep)
		}
		// 从一个已关闭的 channel 接收，selected 会返回true
		return true, false
	}

	if sg := c.sendq.dequeue(); sg != nil {
		// 发现一个等待的发送者。如果缓冲区大小为0，则接收值
		// 直接来自发件人。否则，从队列头接收
		// 并将发送方的值添加到队列的尾部(两者都映射到相同的缓冲区槽，因为队列已满)。
		recv(c, sg, ep, func() { unlock(&c.lock) }, 3)
		return true, true
	}
	// 缓冲类型
	if c.qcount > 0 {
		// 直接从循环数组里找到要接收的元素
		qp := chanbuf(c, c.recvx)
		if raceenabled {
			raceacquire(qp)
			racerelease(qp)
		}
		// 代码里，没有忽略要接收的值，不是 "<- ch"，而是 "val <- ch"，ep 指向 val
		if ep != nil {
			typedmemmove(c.elemtype, ep, qp)
		}
		// 清除掉循环数组里相应位置的值
		typedmemclr(c.elemtype, qp)
		// 接收游标向前移动
		c.recvx++
		// 达到数据的长度，下标重新计算
		if c.recvx == c.dataqsiz {
			c.recvx = 0
		}
		// buf 数组里的元素个数减 1
		c.qcount--
		// 解锁
		unlock(&c.lock)
		return true, true
	}
	// 非阻塞接收，解锁。selected 返回 false，因为没有接收到值
	if !block {
		unlock(&c.lock)
		return false, false
	}

	// 阻塞
	// 构建recvq的阻塞队列
	gp := getg()
	mysg := acquireSudog()
	mysg.releasetime = 0
	if t0 != 0 {
		mysg.releasetime = -1
	}

	mysg.elem = ep
	mysg.waitlink = nil
	gp.waiting = mysg
	mysg.g = gp
	mysg.isSelect = false
	mysg.c = c
	gp.param = nil
	c.recvq.enqueue(mysg)
	// 将当前 goroutine 挂起
	goparkunlock(&c.lock, waitReasonChanReceive, traceEvGoBlockRecv, 3)

	// someone woke us up
	if mysg != gp.waiting {
		throw("G waiting list is corrupted")
	}
	gp.waiting = nil
	if mysg.releasetime > 0 {
		blockevent(mysg.releasetime-t0, 2)
	}
	closed := gp.param == nil
	gp.param = nil
	mysg.c = nil
	releaseSudog(mysg)
	return true, !closed
}
```

梳理下流程图： 

![channel](/img/channel_read.png?raw=true)

### channel的关闭

chanbel的关闭，对于其中的recvq和sendq也就是阻塞的发送者和接收者，对于等待接收者而言，会收到一个相应类型的零值。对于等待发送者，会直接 panic。  
  
channel的关闭不当会出现panic的场景：  

1、关闭值为nil的channel  
2、关闭已经关闭的channel  
3、向已经关闭的channel写入数据  

```go
func closechan(c *hchan) {
	// 关闭值为nil的channel,报错panic
	if c == nil {
		panic(plainError("close of nil channel"))
	}
	// 加锁
	lock(&c.lock)
	// 关闭已经关闭的channel，报错panic
	if c.closed != 0 {
		unlock(&c.lock)
		panic(plainError("close of closed channel"))
	}

	if raceenabled {
		callerpc := getcallerpc()
		racewritepc(c.raceaddr(), callerpc, funcPC(closechan))
		racerelease(c.raceaddr())
	}
	// 修改关闭饿状态
	c.closed = 1

	var glist gList

	// 释放recvq中的sudog
	for {
		// 接收一个sudog
		sg := c.recvq.dequeue()
		// 全部接收完毕了
		if sg == nil {
			break
		}
		// 如果 elem 不为空，说明此 receiver 未忽略接收数据
		// 给它赋一个相应类型的零值
		if sg.elem != nil {
			typedmemclr(c.elemtype, sg.elem)
			sg.elem = nil
		}
		if sg.releasetime != 0 {
			sg.releasetime = cputicks()
		}
		// 取出goroutine
		gp := sg.g
		gp.param = nil
		if raceenabled {
			raceacquireg(gp, c.raceaddr())
		}
		glist.push(gp)
	}

	// 将 channel 等待发送队列里的 sudog 释放
	// 如果存在，这些 goroutine 将会 panic
	for {
		// 取出
		sg := c.sendq.dequeue()
		if sg == nil {
			break
		}
		// 发送者会 panic
		sg.elem = nil
		if sg.releasetime != 0 {
			sg.releasetime = cputicks()
		}
		gp := sg.g
		gp.param = nil
		if raceenabled {
			raceacquireg(gp, c.raceaddr())
		}
		glist.push(gp)
	}
	unlock(&c.lock)

	// Ready all Gs now that we've dropped the channel lock.
	for !glist.empty() {
		// 取出一个
		gp := glist.pop()
		gp.schedlink = 0
		// 唤醒相应 goroutine
		goready(gp, 3)
	}
}
```

### channel的关闭

对于channel的关闭，我们需要注意下：  

1、在不能更改channel状态的情况下，没有简单普遍的方式来检查channel是否已经关闭了  

2、关闭已经关闭的channel会导致panic，所以在closer(关闭者)不知道channel是否已经关闭的情况下去关闭channel是很危险的  

3、发送值到已经关闭的channel会导致panic，所以如果sender(发送者)在不知道channel是否已经关闭的情况下去向channel发送值是很危险的  

对于channel的关闭有这样的一个原则：  

> don't close a channel from the receiver side and don't close a channel if the channel has multiple concurrent senders.

不要从一个receiver侧关闭channel，也不要在有多个sender时，关闭channel。 

向已经关闭的channel发送数据会导致panic，所以在receiver侧关闭，sender是不知道channel是否关闭的，多个sender的情况下，某一个sender关闭了channel,
其他的sender是不知道这个channel是否关闭的，再次写入数据和关闭，都会导致panic。   

#### M个receivers，一个sender

这个是最简单的场景，当需要关闭的时候sender直接关闭就好了

```go
package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)

	// ...
	const MaxRandomNumber = 100000
	const NumReceivers = 100

	// 使用WaitGroup来阻塞查看打印的效果
	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(NumReceivers)

	// 设置channel的长度为10
	dataCh := make(chan int, 100)

	// the sender
	go func() {
		for {
			if value := rand.Intn(MaxRandomNumber); value == 0 {
				// 需要关闭的时候直接关闭就好了，是很安全的
				close(dataCh)
				return
			} else {
				dataCh <- value
			}
		}
	}()

	// receivers
	for i := 0; i < NumReceivers; i++ {
		go func() {
			defer wgReceivers.Done()

			// 监听dataCh，接收里面的值
			for value := range dataCh {
				log.Println(value)
			}
		}()
	}

	wgReceivers.Wait()

}
```

#### 一个receiver，N个sender

这种情况下可以给个`signal channel`，然后通知senders去停止向channel发送数据。因为receiver不能去关闭channel,这样senders将会触发panic,
但是我们可以让receiver通知`signal channel`来告诉senders来停止发送。  

```go
package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)

	// ...
	const MaxRandomNumber = 100000
	const NumSenders = 1000

	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(1)

	// 发送数据的channel
	dataCh := make(chan int, 100)

	// 无缓冲的channel作为信号量，通知senders的推出
	stopCh := make(chan struct{})

	// 启动个NumSenders个sender
	for i := 0; i < NumSenders; i++ {
		go func() {
			for {
				value := rand.Intn(MaxRandomNumber)

				// 监测到退出信号，马上退出goroutine
				// 否则正常写入dataCh，数据
				select {
				case <-stopCh:
					return
				case dataCh <- value:
				}
			}
		}()
	}

	// 消费者
	go func() {
		defer wgReceivers.Done()

		for value := range dataCh {
			// 某个场景下发出退出的信号量
			if value == MaxRandomNumber-1 {
				close(stopCh)
				return
			}

			log.Println(value)
		}
	}()

	wgReceivers.Wait()
}
```

通过stopCh来，作为信号量，来通知发送者的goroutine退出。 

#### M个receiver，N个sender

这是比较复杂的一个，相比上面的一个receiver，M个receiver中任意一个receiver发出关闭的信息，需要同步到其他的receiver，防止其他的receiver，
再次发出关闭的请求，出发panic。  

```go
package main

import (
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)

	const MaxRandomNumber = 100000
	const NumReceivers = 10
	const NumSenders = 1000

	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(NumReceivers)

	// 数据的channel
	dataCh := make(chan int, 100)
	// 关闭的channel的信号
	stopCh := make(chan struct{})
	// toStop通知关闭stopCh，同时作为receiver退出的信息
	toStop := make(chan string, 1)

	var stoppedBy string

	// 当收到toStop的信号，关闭stopCh
	go func() {
		stoppedBy = <-toStop
		close(stopCh)
	}()

	// 发送端
	for i := 0; i < NumSenders; i++ {
		go func(id string) {
			for {
				value := rand.Intn(MaxRandomNumber)
				// 满足条件发出关闭的请求到toStop
				if value == 0 {
					select {
					case toStop <- "sender#" + id:
					default:
					}
					return
				}

				select {
				// 检测的关闭的stopCh，退出发送者
				case <-stopCh:
					return
				case dataCh <- value:
				}
			}
		}(strconv.Itoa(i))
	}

	// 接收端
	for i := 0; i < NumReceivers; i++ {
		go func(id string) {
			defer wgReceivers.Done()

			for {
				select {
				// 检测的关闭的stopCh，退出接收者
				case <-stopCh:
					return
				case value := <-dataCh:
					// 满足条件发出关闭的请求到toStop
					if value == MaxRandomNumber-1 {
						select {
						case toStop <- "receiver#" + id:
						default:
						}
						return
					}

					log.Println(value)
				}
			}
		}(strconv.Itoa(i))
	}

	wgReceivers.Wait()
	log.Println("stopped by", stoppedBy)
}
```

这样的设计就很好了，可以在sender和receiver两端发出关闭的请求。保证了sender和receiver都能够退出。   

### 控制goroutine的数量  

go中在大量并发的情况下会产生很多的goroutine，而goroutine使用之后，是不会被完全回收的，大概会有2kb的空间，所以我们希望控制下goroytine的
并发数量。  




### 参考

【Go的CSP并发模型】https://www.jianshu.com/p/a3c9a05466e1  
【goroutine, channel 和 CSP】http://www.moye.me/2017/05/05/go-concurrency-patterns/  
【通过同步和加锁解决多线程的线程安全问题】https://blog.ailemon.me/2019/05/15/solving-multithreaded-thread-safety-problems-by-synchronization-and-locking/  
【golang channel 有缓冲 与 无缓冲 的重要区别】https://my.oschina.net/u/157514/blog/149192  
【Golang channel 源码深度剖析】https://www.cyhone.com/articles/analysis-of-golang-channel/   
【《Go专家编程》Go channel实现原理剖析】https://my.oschina.net/renhc/blog/2246871  
【深度解密Go语言之channel】https://www.cnblogs.com/qcrao-2018/p/11220651.html  
【如何优雅地关闭Go channel】https://www.jianshu.com/p/d24dfbb33781  