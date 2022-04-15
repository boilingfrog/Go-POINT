<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [深入了解下 go 中的 select](#%E6%B7%B1%E5%85%A5%E4%BA%86%E8%A7%A3%E4%B8%8B-go-%E4%B8%AD%E7%9A%84-select)
  - [前言](#%E5%89%8D%E8%A8%80)
    - [1、栗子一](#1%E6%A0%97%E5%AD%90%E4%B8%80)
    - [2、栗子二](#2%E6%A0%97%E5%AD%90%E4%BA%8C)
    - [3、栗子三](#3%E6%A0%97%E5%AD%90%E4%B8%89)
  - [看下源码实现](#%E7%9C%8B%E4%B8%8B%E6%BA%90%E7%A0%81%E5%AE%9E%E7%8E%B0)
    - [1、不存在 case](#1%E4%B8%8D%E5%AD%98%E5%9C%A8-case)
    - [2、select 中仅存在一个 case](#2select-%E4%B8%AD%E4%BB%85%E5%AD%98%E5%9C%A8%E4%B8%80%E4%B8%AA-case)
    - [3、select 中存在两个 case，其中一个是 default](#3select-%E4%B8%AD%E5%AD%98%E5%9C%A8%E4%B8%A4%E4%B8%AA-case%E5%85%B6%E4%B8%AD%E4%B8%80%E4%B8%AA%E6%98%AF-default)
      - [发送值](#%E5%8F%91%E9%80%81%E5%80%BC)
      - [接收值](#%E6%8E%A5%E6%94%B6%E5%80%BC)
    - [4、多个 case 的场景](#4%E5%A4%9A%E4%B8%AA-case-%E7%9A%84%E5%9C%BA%E6%99%AF)
      - [具体的实现逻辑](#%E5%85%B7%E4%BD%93%E7%9A%84%E5%AE%9E%E7%8E%B0%E9%80%BB%E8%BE%91)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 深入了解下 go 中的 select

### 前言

这里借助于几个经常遇到的 select 的使用 demo 来作为开始，先来看看，下面几个 demo 的输出情况  

#### 1、栗子一  

```go
func main() {
	chan1 := make(chan int)
	chan2 := make(chan int)

	go func() {
		chan1 <- 1
	}()

	go func() {
		chan2 <- 1
	}()

	select {
	case <-chan1:
		fmt.Println("chan1 ready.")
	case <-chan2:
		fmt.Println("chan2 ready.")
	default:
		fmt.Println("default")
	}
}
```

select 中的 case 执行是随机的，所以当 case 监听的 channel 有数据传入，就执行相应的流程并退出 select，如果对应的 case 没有收到 channel 的数据，就执行 default 语句，然后退出 select。  

上面的协程启动时间是无法预估的，所以上面的两个 case 和 default ，都有机会执行。  

可能的输出  

可能输出1、
```
chan1 ready.
```

可能输出2、
```
chan2 ready.
```

可能输出3、
```
default
```

#### 2、栗子二

```go
func main() {
	chan1 := make(chan int)
	chan2 := make(chan int)

	go func() {
		close(chan1)
	}()

	go func() {
		close(chan2)
	}()

	select {
	case <-chan1:
		fmt.Println("chan1 ready.")
	case <-chan2:
		fmt.Println("chan2 ready.")
	default:
		fmt.Println("default")
	}
}
```

已经关闭的 channel ，使用 select 是可以从中读出对应的零值，同时两面关闭 channel 的协程的执行实际也是不可控的，原则上，上面两个 case 和 default 都有可能被执行。   

可能的输出  

可能输出1、
```
chan1 ready.
```

可能输出2、
```
chan2 ready.
```

可能输出3、
```
default
```

#### 3、栗子三

```go
func main() {
	select {}
}
```

上面这个，应为没有机会退出，所以会发生死锁  

### 看下源码实现

select 中的多个 case 是随机触发执行的，一次只有一个 case 得到执行。如果我们按照顺序依次判断，那么后面的条件永远都会得不到执行，而随机的引入就是为了避免饥饿问题的发生。  

1、如果没有 default 分支

如果没有 default 分支，select 将会一直处于阻塞状态，直到其中的一个 case 就绪；  

2、如果有 default 分支

如果有 default 分支，随机将 case 分支遍历一遍，如果有 case 分支可执行，处理对应的 case 分支；  

如果 遍历完 case 分支，没有可执行的分支，执行 default 分支。  

源码版本 `go version go1.16.13 darwin/amd64`

源码包 `src/runtime/select.go` 定义了表示case语句的数据结构：  

```
type scase struct {
	c    *hchan         // chan
	elem unsafe.Pointer // data element
}
```

c为当前 case 语句所操作的 channel 指针，这也说明了一个 case 语句只能操作一个 channel。  


编译阶段，select 对应的 opType 是 OSELECT，select 语句在编译期间会被转换成 OSELECT 节点。  

```
// https://github.com/golang/go/blob/release-branch.go1.16/src/cmd/compile/internal/gc/syntax.go#L922
OSELECT // select { List } (List is list of OCASE)
```

如果是 OSELECT 就会调用 `walkselect()`,然后 `walkselect()` 最后调用 `walkselectcases()`   

```go
// https://github.com/golang/go/blob/release-branch.go1.16/src/cmd/compile/internal/gc/walk.go#L104
// The result of walkstmt MUST be assigned back to n, e.g.
// 	n.Left = walkstmt(n.Left)
func walkstmt(n *Node) *Node {
	if n == nil {
		return n
	}

	setlineno(n)

	walkstmtlist(n.Ninit.Slice())

	switch n.Op {
		...
	case OSELECT:
		walkselect(n)

	case OSWITCH:
		walkswitch(n)

	case ORANGE:
		n = walkrange(n)
	}

	if n.Op == ONAME {
		Fatalf("walkstmt ended up with name: %+v", n)
	}
	return n
}

// https://github.com/golang/go/blob/release-branch.go1.16/src/cmd/compile/internal/gc/select.go#L90
func walkselect(sel *Node) {
	lno := setlineno(sel)
	if sel.Nbody.Len() != 0 {
		Fatalf("double walkselect")
	}

	init := sel.Ninit.Slice()
	sel.Ninit.Set(nil)
    // 调用walkselectcases
	init = append(init, walkselectcases(&sel.List)...)
	sel.List.Set(nil)

	sel.Nbody.Set(init)
	walkstmtlist(sel.Nbody.Slice())

	lineno = lno
}
```

上面的调用逻辑，select 的逻辑是在 `walkselectcases()` 函数中完成的，这里来重点看下  

`walkselectcases()` 在处理中会分成下面几种情况来处理  

1、select 中不存在 case, 直接堵塞；  

2、select 中仅存在一个 case；  

3、select 中存在两个 case，其中一个是 default；  

4、其他 select 情况如: 包含多个 case 并且有 default 等。  

```go
// src/cmd/compile/internal/gc/select.go
func walkselectcases(cases *Nodes) []*Node {
	// 获取 case 分支的数量
	n := cases.Len()

	// 优化: 没有 case 的情况
	if n == 0 {
		// 翻译为：block()
		...
		return
	}

	// 优化: 只有一个 case 的情况
	if n == 1 {
		// 翻译为：if ch == nil { block() }; n;
		...
		return
	}

	// 优化: select 中存在两个 case，其中一个是 default 的情况
	if n == 2 {
		// 翻译为：发送或接收
		// if selectnbsend(c, v) { body } else { default body }

		// 接收 
		// if selectnbrecv(&v, &received, c) { body } else { default body }
		return
	}

	// 一般情况，调用 selecggo
	...
}
```

#### 1、不存在 case

如果不存在 case ，空的 select 语句会直接阻塞当前 Goroutine，导致 Goroutine 进入无法被唤醒的永久休眠状态。  

```go
// https://github.com/golang/go/blob/release-branch.go1.16/src/cmd/compile/internal/gc/walk.go#L104
func walkselectcases(cases *Nodes) []*Node {
	n := cases.Len()

	if n == 0 {
		return []*Node{mkcall("block", nil, nil)}
	}
	...
}

// 调用 runtime.gopark 让出当前 Goroutine 对处理器的使用权并传入等待原因 waitReasonSelectNoCases。
func block() {
	gopark(nil, nil, waitReasonSelectNoCases, traceEvGoStop, 1)
}
```

如果没有 case，导致 Goroutine 进入无法被唤醒的永久休眠状态，会触发 `deadlock!`  

#### 2、select 中仅存在一个 case

如果只有一个 case ，编译器会将 select 改写成 if 条件语句。  

```go
// 改写前
select {
case v, ok <-ch: // case ch <- v
    ...    
}

// 改写后
if ch == nil {
    block()
}
v, ok := <-ch // case ch <- v
...
```

如果只有一个 case ，walkselectcases 会将 select 根据收发情况装换成 if 语句，如果 case 中的 Channel 是空指针时，会直接挂起当前 Goroutine 并陷入永久休眠。  

#### 3、select 中存在两个 case，其中一个是 default  

##### 发送值  

在 walkselectcases 中 OSEND，对应的就是向 channel 中发送数据，如果是发送的话，会翻译成下面的语句  

```go
select {
case ch <- i:
    ...
default:
    ...
}

if selectnbsend(ch, i) {
    ...
} else {
    // default body
    ...
}

func selectnbsend(c *hchan, elem unsafe.Pointer) (selected bool) {
	return chansend(c, elem, false, getcallerpc())
}
```

如果是发送，这里翻译之后最终调用 chansend 向 channel 中发送数据  

```go
// 这里提供了一个 block，参数设置成 true，那么表示当前发送操作是阻塞的
func chansend(c *hchan, ep unsafe.Pointer, block bool, callerpc uintptr) bool {
	...
	// 对于不阻塞的 send，快速检测失败场景
	//
	// 如果 channel 未关闭且 channel 没有多余的缓冲空间。这可能是：
	// 1. channel 是非缓冲型的，且等待接收队列里没有 goroutine
	// 2. channel 是缓冲型的，但循环数组已经装满了元素
	if !block && c.closed == 0 && ((c.dataqsiz == 0 && c.recvq.first == nil) ||
		(c.dataqsiz > 0 && c.qcount == c.dataqsiz)) {
		return false
	}
	...
}
```

总结下  

- 1、如果 block 为 true 表示当前向 channel 中的数据发送是阻塞的。这里可以看到 selectnbsend 中传入的是 false,说明 channel 的发送不会阻塞 select。  

- 2、对于不阻塞的发送，会进行下面的检测，如果 channel 未关闭且 channel 没有多余的缓冲空间，就会发送失败，然后跳出当前的 case,走到 default 的逻辑。  

如果 channel 未关闭且 channel 没有多余的缓冲空间。这可能是： 

1、channel 是非缓冲型的，且等待接收队列里没有 goroutine；  

2、channel 是缓冲型的，但循环数组已经装满了元素；  

##### 接收值

在 walkselectcases 函数中可以看到，接收方式会有两个，分别是 OSELRECV 和 OSELRECV2  

```go
// https://github.com/golang/go/blob/release-branch.go1.16/src/cmd/compile/internal/gc/walk.go#L104
func walkselectcases(cases *Nodes) []*Node {
	...
	// optimization: two-case select but one is default: single non-blocking op.
	if ncas == 2 && dflt != nil {
		switch n.Op {
		default:
			Fatalf("select %v", n.Op)

		case OSELRECV:
			// if selectnbrecv(&v, c) { body } else { default body }
			...
			r.Left = mkcall1(chanfn("selectnbrecv", 2, ch.Type), types.Types[TBOOL], &r.Ninit, elem, ch)

		case OSELRECV2:
			// if selectnbrecv2(&v, &received, c) { body } else { default body }
			...
			r.Left = mkcall1(chanfn("selectnbrecv2", 2, ch.Type), types.Types[TBOOL], &r.Ninit, elem, receivedp, ch)
		}

		r.Left = typecheck(r.Left, ctxExpr)
		r.Nbody.Set(cas.Nbody.Slice())
		r.Rlist.Set(append(dflt.Ninit.Slice(), dflt.Nbody.Slice()...))
		return []*Node{r, nod(OBREAK, nil, nil)}
	}
	...
}
```

walkselectcases 对这两种情况的改写  

selectnbrecv  

```go
select {
case v = <-c:
	...
default:
	...
}

// 改写后

if selectnbrecv(&v, c) {
	...
} else {
    // default body
	...
}
```

selectnbrecv2  

```go
select {
case v, ok = <-c:
	... foo
default:
	... bar
}

// 改写后

if c != nil && selectnbrecv2(&v, &ok, c) {
	... foo
} else {
    // default body
	... bar
}
```

selectnbrecv 和 selectnbrecv2 有什么区别呢？  

```go
// https://github.com/golang/go/blob/release-branch.go1.16/src/runtime/chan.go#L707
func selectnbrecv(elem unsafe.Pointer, c *hchan) (selected bool) {
	selected, _ = chanrecv(c, elem, false)
	return
}

func selectnbrecv2(elem unsafe.Pointer, received *bool, c *hchan) (selected bool) {
	// TODO(khr): just return 2 values from this function, now that it is in Go.
	selected, *received = chanrecv(c, elem, false)
	return
}
```

可以发现只是针对返回的值处理不同，selectnbrecv2 多了一个是否 received 的 bool 值  

总结下：  

对于接收值的 case 会有两种处理方式，这两种，区别在于是否将 received 的 bool 值传送给调用方  

#### 4、多个 case 的场景  

多个 case 的场景  

1、会将其中所有 case 转化为 scase 结构体；  

2、调用运行时函数 selectgo 选取触发的 scase 结构体；  

3、通过 for 循环生成一组 if 语句,来判断是否选中 case；  

这里来看下 selectgo 的实现     

`这里看下参数`    

- cas0：为 scase 数组的首地址，selectgo() 就是从这些 scase 中找出一个返回；  

- order0：为一个两倍 cas0 数组长度的 buffer，保存 scase 随机序列 pollorder 和 scase 中 channel 地址序列 lockorder,数组前一半是 pollorder,后一半用来 lockorder；  

pollorder：每次 selectgo 执行都会把 scase 序列打乱，以达到随机检测 case 的目的；  

lockorder：所有 case 语句中 channel 序列，以达到去重防止对 channel 加锁时重复加锁的目的；  
     
- nsends: 发送的 case 的个数；   

- nrecvs: 接收的 case 的个数； 

- block: 表示是否存在 default,没有 default 就表示 select 是阻塞的。  

`看下返回的数据`  

- int： 选中case的编号，这个case编号跟代码一致；  

- bool: 是否成功从channle中读取了数据，如果选中的case是从channel中读数据，则该返回值表示是否读取成功。  

##### 具体的实现逻辑  

- 1、打乱 scase 的顺序，将锁定scase语句中所有的channel；  

- 2、按照随机顺序检测scase中的channel是否ready；  

2.1 如果case可读，则读取channel中数据，解锁所有的channel，然后返回(case index, true)  

2.2 如果case可写，则将数据写入channel，解锁所有的channel，然后返回(case index, false)  

2.3 所有case都未ready，并且有default语句，则解锁所有的channel，然后返回（default index, false）  

- 3、所有case都未ready，且没有default语句  

3.1 将当前协程加入到所有channel的等待队列  

3.2 当将协程转入阻塞，等待被唤醒  

- 4、唤醒后返回channel对应的case index  

4.1 如果是读操作，解锁所有的channel，然后返回(case index, true)  

4.2 如果是写操作，解锁所有的channel，然后返回(case index, false)  

这里来分析下 selectgo 的具体实现  

```go
func selectgo(cas0 *scase, order0 *uint16, pc0 *uintptr, nsends, nrecvs int, block bool) (int, bool) {
	if debugSelect {
		print("select: cas0=", cas0, "\n")
	}

	// NOTE: In order to maintain a lean stack size, the number of scases
	// is capped at 65536.
	cas1 := (*[1 << 16]scase)(unsafe.Pointer(cas0))
	order1 := (*[1 << 17]uint16)(unsafe.Pointer(order0))

	ncases := nsends + nrecvs
	scases := cas1[:ncases:ncases]
	pollorder := order1[:ncases:ncases]
	lockorder := order1[ncases:][:ncases:ncases]
	// NOTE: pollorder/lockorder's underlying array was not zero-initialized by compiler.

	// Even when raceenabled is true, there might be select
	// statements in packages compiled without -race (e.g.,
	// ensureSigM in runtime/signal_unix.go).
	var pcs []uintptr
	if raceenabled && pc0 != nil {
		pc1 := (*[1 << 16]uintptr)(unsafe.Pointer(pc0))
		pcs = pc1[:ncases:ncases]
	}
	casePC := func(casi int) uintptr {
		if pcs == nil {
			return 0
		}
		return pcs[casi]
	}

	var t0 int64
	if blockprofilerate > 0 {
		t0 = cputicks()
	}

	// The compiler rewrites selects that statically have
	// only 0 or 1 cases plus default into simpler constructs.
	// The only way we can end up with such small sel.ncase
	// values here is for a larger select in which most channels
	// have been nilled out. The general code handles those
	// cases correctly, and they are rare enough not to bother
	// optimizing (and needing to test).

	// 生成随机顺序
	norder := 0
	for i := range scases {
		cas := &scases[i]

		// 忽略轮询和锁定命令中没有通道的情况
		if cas.c == nil {
			cas.elem = nil // allow GC
			continue
		}

		j := fastrandn(uint32(norder + 1))
		pollorder[norder] = pollorder[j]
		pollorder[j] = uint16(i)
		norder++
	}
	pollorder = pollorder[:norder]
	lockorder = lockorder[:norder]

	// 根据 channel 地址进行排序,决定获取锁的顺序
	for i := range lockorder {
		j := i
		// Start with the pollorder to permute cases on the same channel.
		c := scases[pollorder[i]].c
		for j > 0 && scases[lockorder[(j-1)/2]].c.sortkey() < c.sortkey() {
			k := (j - 1) / 2
			lockorder[j] = lockorder[k]
			j = k
		}
		lockorder[j] = pollorder[i]
	}
	...

	// 锁定选中的 channel
	sellock(scases, lockorder)

	var (
		gp     *g
		sg     *sudog
		c      *hchan
		k      *scase
		sglist *sudog
		sgnext *sudog
		qp     unsafe.Pointer
		nextp  **sudog
	)

	// pass 1 - 遍历所有 scase,确定已经准备好的 scase
	var casi int
	var cas *scase
	var caseSuccess bool
	var caseReleaseTime int64 = -1
	var recvOK bool
	for _, casei := range pollorder {
		casi = int(casei)
		cas = &scases[casi]
		c = cas.c
		// 接收数据
		if casi >= nsends {
			// 有 goroutine 等待发送数据
			sg = c.sendq.dequeue()
			if sg != nil {
				goto recv
			}
			// 缓冲区有数据
			if c.qcount > 0 {
				goto bufrecv
			}
			// 通道关闭
			if c.closed != 0 {
				goto rclose
			}
			// 发送数据
		} else {
			if raceenabled {
				racereadpc(c.raceaddr(), casePC(casi), chansendpc)
			}
			// 判断通道的关闭情况
			if c.closed != 0 {
				goto sclose
			}
			// 接收等待队列有 goroutine
			sg = c.recvq.dequeue()
			if sg != nil {
				goto send
			}
			// 缓冲区有空位置
			if c.qcount < c.dataqsiz {
				goto bufsend
			}
		}
	}

	// 如果不阻塞，意味着有 default,准备退出select
	if !block {
		selunlock(scases, lockorder)
		casi = -1
		goto retc
	}

	// pass 2 - 所有 channel 入队，等待处理
	gp = getg()
	if gp.waiting != nil {
		throw("gp.waiting != nil")
	}
	nextp = &gp.waiting
	for _, casei := range lockorder {
		casi = int(casei)
		// 获取一个 scase
		cas = &scases[casi]
		// 监听的 channel
		c = cas.c
		sg := acquireSudog()
		sg.g = gp
		sg.isSelect = true
		// No stack splits between assigning elem and enqueuing
		// sg on gp.waiting where copystack can find it.
		sg.elem = cas.elem
		sg.releasetime = 0
		if t0 != 0 {
			sg.releasetime = -1
		}
		sg.c = c
		// 按锁定顺序构造等待列表。
		*nextp = sg
		nextp = &sg.waitlink

		if casi < nsends {
			c.sendq.enqueue(sg)
		} else {
			c.recvq.enqueue(sg)
		}
	}

	// goroutine 陷入睡眠,等待某一个 channel 唤醒 gooutine
	gp.param = nil
	// Signal to anyone trying to shrink our stack that we're about
	// to park on a channel. The window between when this G's status
	// changes and when we set gp.activeStackChans is not safe for
	// stack shrinking.
	atomic.Store8(&gp.parkingOnChan, 1)
	gopark(selparkcommit, nil, waitReasonSelect, traceEvGoBlockSelect, 1)
	gp.activeStackChans = false

	sellock(scases, lockorder)

	gp.selectDone = 0
	sg = (*sudog)(gp.param)
	gp.param = nil

	// pass 3 - 删除队列中没有触发的 channels
	// 如果不删除的话,他们会触发 channel.我们按锁的顺序单向链接 sudog
	casi = -1
	cas = nil
	caseSuccess = false
	sglist = gp.waiting
	// 在从 gp.waiting 取消链接之前清除所有元素。
	for sg1 := gp.waiting; sg1 != nil; sg1 = sg1.waitlink {
		sg1.isSelect = false
		sg1.elem = nil
		sg1.c = nil
	}
	gp.waiting = nil

	for _, casei := range lockorder {
		k = &scases[casei]
		if sg == sglist {
			// sg has already been dequeued by the G that woke us up.
			casi = int(casei)
			cas = k
			caseSuccess = sglist.success
			if sglist.releasetime > 0 {
				caseReleaseTime = sglist.releasetime
			}
		} else {
			c = k.c
			if int(casei) < nsends {
				c.sendq.dequeueSudoG(sglist)
			} else {
				c.recvq.dequeueSudoG(sglist)
			}
		}
		sgnext = sglist.waitlink
		sglist.waitlink = nil
		releaseSudog(sglist)
		sglist = sgnext
	}

	if cas == nil {
		throw("selectgo: bad wakeup")
	}

	c = cas.c

	if debugSelect {
		print("wait-return: cas0=", cas0, " c=", c, " cas=", cas, " send=", casi < nsends, "\n")
	}

	if casi < nsends {
		if !caseSuccess {
			goto sclose
		}
	} else {
		recvOK = caseSuccess
	}

	if raceenabled {
		if casi < nsends {
			raceReadObjectPC(c.elemtype, cas.elem, casePC(casi), chansendpc)
		} else if cas.elem != nil {
			raceWriteObjectPC(c.elemtype, cas.elem, casePC(casi), chanrecvpc)
		}
	}
	if msanenabled {
		if casi < nsends {
			msanread(cas.elem, c.elemtype.size)
		} else if cas.elem != nil {
			msanwrite(cas.elem, c.elemtype.size)
		}
	}

	selunlock(scases, lockorder)
	goto retc

bufrecv:
	// 可以从 buffer 接收 
	if raceenabled {
		if cas.elem != nil {
			raceWriteObjectPC(c.elemtype, cas.elem, casePC(casi), chanrecvpc)
		}
		racenotify(c, c.recvx, nil)
	}
	if msanenabled && cas.elem != nil {
		msanwrite(cas.elem, c.elemtype.size)
	}
	recvOK = true
	qp = chanbuf(c, c.recvx)
	if cas.elem != nil {
		typedmemmove(c.elemtype, cas.elem, qp)
	}
	typedmemclr(c.elemtype, qp)
	c.recvx++
	if c.recvx == c.dataqsiz {
		c.recvx = 0
	}
	c.qcount--
	selunlock(scases, lockorder)
	goto retc

bufsend:
	// 可以发送到 buffer
	if raceenabled {
		racenotify(c, c.sendx, nil)
		raceReadObjectPC(c.elemtype, cas.elem, casePC(casi), chansendpc)
	}
	if msanenabled {
		msanread(cas.elem, c.elemtype.size)
	}
	typedmemmove(c.elemtype, chanbuf(c, c.sendx), cas.elem)
	c.sendx++
	if c.sendx == c.dataqsiz {
		c.sendx = 0
	}
	c.qcount++
	selunlock(scases, lockorder)
	goto retc

recv:
	// 可以从一个休眠的发送方 (sg)直接接收
	recv(c, sg, cas.elem, func() { selunlock(scases, lockorder) }, 2)
	if debugSelect {
		print("syncrecv: cas0=", cas0, " c=", c, "\n")
	}
	recvOK = true
	goto retc

rclose:
	// 在已经关闭的 channel 末尾进行读
	selunlock(scases, lockorder)
	recvOK = false
	if cas.elem != nil {
		typedmemclr(c.elemtype, cas.elem)
	}
	if raceenabled {
		raceacquire(c.raceaddr())
	}
	goto retc

send:
	// 可以向一个休眠的接收方 (sg) 发送
	if raceenabled {
		raceReadObjectPC(c.elemtype, cas.elem, casePC(casi), chansendpc)
	}
	if msanenabled {
		msanread(cas.elem, c.elemtype.size)
	}
	send(c, sg, cas.elem, func() { selunlock(scases, lockorder) }, 2)
	if debugSelect {
		print("syncsend: cas0=", cas0, " c=", c, "\n")
	}
	goto retc

retc:
	if caseReleaseTime > 0 {
		blockevent(caseReleaseTime-t0, 1)
	}
	return casi, recvOK

sclose:
	// 向已关闭的 channel 进行发送
	selunlock(scases, lockorder)
	panic(plainError("send on closed channel"))
}
```

### 参考

【Select 语句的本质】https://golang.design/under-the-hood/zh-cn/part1basic/ch03lang/chan/#select-    
【GO专家编程】https://book.douban.com/subject/35144587/  



