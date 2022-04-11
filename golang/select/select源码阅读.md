<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [深入了解下 go 中的 select](#%E6%B7%B1%E5%85%A5%E4%BA%86%E8%A7%A3%E4%B8%8B-go-%E4%B8%AD%E7%9A%84-select)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [看下源码实现](#%E7%9C%8B%E4%B8%8B%E6%BA%90%E7%A0%81%E5%AE%9E%E7%8E%B0)
    - [1、不存在 case](#1%E4%B8%8D%E5%AD%98%E5%9C%A8-case)
    - [2、select 中仅存在一个 case](#2select-%E4%B8%AD%E4%BB%85%E5%AD%98%E5%9C%A8%E4%B8%80%E4%B8%AA-case)
    - [3、select 中存在两个 case，其中一个是 default](#3select-%E4%B8%AD%E5%AD%98%E5%9C%A8%E4%B8%A4%E4%B8%AA-case%E5%85%B6%E4%B8%AD%E4%B8%80%E4%B8%AA%E6%98%AF-default)
      - [发送值](#%E5%8F%91%E9%80%81%E5%80%BC)
      - [接收值](#%E6%8E%A5%E6%94%B6%E5%80%BC)
    - [4、多个 case 的场景](#4%E5%A4%9A%E4%B8%AA-case-%E7%9A%84%E5%9C%BA%E6%99%AF)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 深入了解下 go 中的 select

### 前言

这里借助于几个经常遇到的 select 的使用 demo 来作为开始，先来看看，下面几个 demo 的输出情况  

1、栗子一  

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

2、栗子二

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

3、栗子三

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

如果有 default 分支，当所有的 case 分支

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




### 参考

【Select 语句的本质】https://golang.design/under-the-hood/zh-cn/part1basic/ch03lang/chan/#select-    



