<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [深入了解下 go 中的 select](#%E6%B7%B1%E5%85%A5%E4%BA%86%E8%A7%A3%E4%B8%8B-go-%E4%B8%AD%E7%9A%84-select)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [看下源码实现](#%E7%9C%8B%E4%B8%8B%E6%BA%90%E7%A0%81%E5%AE%9E%E7%8E%B0)
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

如果是 OSELECT 就会调用 `walkselect()`,然后 `walkselect()` 会调用 `walkselectcases()`   

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

上面的调用逻辑，如果是 select 的逻辑是在 `walkselectcases()` 函数中完成的，这里来重点看下  

walkselectcases() 在处理中会分成下面几种情况来处理  

1、select 中不存在 case, 直接堵塞；  

2、select 中仅存在一个 case；  

3、select 中存在两个 case，其中一个是 default；  

4、其他 select 情况如: 包含多个 case 并且有 default 等。  











### 参考


