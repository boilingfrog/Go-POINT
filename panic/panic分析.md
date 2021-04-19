<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [panic源码解读](#panic%E6%BA%90%E7%A0%81%E8%A7%A3%E8%AF%BB)
  - [panic的作用](#panic%E7%9A%84%E4%BD%9C%E7%94%A8)
    - [panic使用场景](#panic%E4%BD%BF%E7%94%A8%E5%9C%BA%E6%99%AF)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## panic源码解读

### panic的作用

- `panic`能够改变程序的控制流，调用`panic`后会立刻停止执行当前函数的剩余代码，并在当前`Goroutine`中递归执行调用方的`defer`；  

- `recover`可以中止`panic`造成的程序崩溃。它是一个只能在`defer`中发挥作用的函数，在其他作用域中调用不会发挥作用；  

举个栗子  

```go
package main

import "fmt"

func main() {
	fmt.Println(1)
	func() {
		fmt.Println(2)
		panic("3")
	}()
	fmt.Println(4)
}
```

输出  

```go
1
2
panic: 3

goroutine 1 [running]:
main.main.func1(...)
        /Users/yj/Go/src/Go-POINT/panic/main.go:9
main.main()
        /Users/yj/Go/src/Go-POINT/panic/main.go:10 +0xee
```

`panic`后会立刻停止执行当前函数的剩余代码，所以4没有打印出来  

**对于recover**

- panic只会触发当前Goroutine的defer；

- recover只有在defer中调用才会生效；

- panic允许在defer中嵌套多次调用；

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(1)

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		fmt.Println(2)
		panic("3")
	}()
	time.Sleep(time.Second)
	fmt.Println(4)
}
```

上面的栗子，因为`recover`和`panic`不在同一个`goroutine`中，所以不会捕获到  

嵌套的demo  

```go
func main() {
	defer fmt.Println("in main")
	defer func() {
		defer func() {
			panic("3 panic again and again")
		}()
		panic("2 panic again")
	}()

	panic("1 panic once")
}
```

输出  

```go
in main
panic: 1 panic once
        panic: 2 panic again
        panic: 3 panic again and again

goroutine 1 [running]:
...
```

多次调用`panic`也不会影响`defer`函数的正常执行，所以使用`defer`进行收尾工作一般来说都是安全的。  

#### panic使用场景

- error：可预见的错误

- panic：不可预见的异常

需要注意的是，你应该尽可能地使用`error`，而不是使用`panic`和`recover`。只有当程序不能继续运行的时候，才应该使用`panic`和`recover`机制。    

`panic`有两个合理的用例。  

1、发生了一个不能恢复的错误，此时程序不能继续运行。 一个例子就是 web 服务器无法绑定所要求的端口。在这种情况下，就应该使用 panic，因为如果不能绑定端口，啥也做不了。  

2、发生了一个编程上的错误。 假如我们有一个接收指针参数的方法，而其他人使用 nil 作为参数调用了它。在这种情况下，我们可以使用panic，因为这是一个编程错误：用 nil 参数调用了一个只能接收合法指针的方法。  

在一般情况下，我们不应通过调用panic函数来报告普通的错误，而应该只把它作为报告致命错误的一种方式。当某些不应该发生的场景发生时，我们就应该调用panic。  

总结下`panic`的使用场景:  

- 1、空指针引用

- 2、下标越界

- 3、除数为0

- 4、不应该出现的分支，比如default

- 5、输入不应该引起函数错误  

### 看下实现

先来看下`_panic`的结构  

```go
// _panic 保存了一个活跃的 panic
//
// 这个标记了 go:notinheap 因为 _panic 的值必须位于栈上
//
// argp 和 link 字段为栈指针，但在栈增长时不需要特殊处理：因为他们是指针类型且
// _panic 值只位于栈上，正常的栈指针调整会处理他们。
//
//go:notinheap
type _panic struct {
	argp      unsafe.Pointer // panic 期间 defer 调用参数的指针; 无法移动 - liblink 已知
	arg       interface{}    // panic的参数
	link      *_panic        // link 链接到更早的 panic
	recovered bool           // panic是否结束
	aborted   bool           // panic是否被忽略
}
```

`link`指向了保存在`goroutine`链表中先前的`panic`链表  

#### gopanic 和 gorecover

```go
// 预先声明的函数 panic 的实现
func gopanic(e interface{}) {
	gp := getg()
	// 判断在系统栈上还是在用户栈上
	// 如果执行在系统或信号栈时，getg() 会返回当前 m 的 g0 或 gsignal
	// 因此可以通过 gp.m.curg == gp 来判断所在栈
	// 系统栈上的 panic 无法恢复
	if gp.m.curg != gp {
		print("panic: ")
		printany(e)
		print("\n")
		throw("panic on system stack")
	}
	// 如果正在进行 malloc 时发生 panic 也无法恢复
	if gp.m.mallocing != 0 {
		print("panic: ")
		printany(e)
		print("\n")
		throw("panic during malloc")
	}
	// 在禁止抢占时发生 panic 也无法恢复
	if gp.m.preemptoff != "" {
		print("panic: ")
		printany(e)
		print("\n")
		print("preempt off reason: ")
		print(gp.m.preemptoff)
		print("\n")
		throw("panic during preemptoff")
	}
	// 在 g 锁在 m 上时发生 panic 也无法恢复
	if gp.m.locks != 0 {
		print("panic: ")
		printany(e)
		print("\n")
		throw("panic holding locks")
	}

	// 下面是可以恢复的
	var p _panic
	p.arg = e
	// panic 保存了对应的消息，并指向了保存在 goroutine 链表中先前的 panic 链表
	p.link = gp._panic
	gp._panic = (*_panic)(noescape(unsafe.Pointer(&p)))

	atomic.Xadd(&runningPanicDefers, 1)

	for {
		// 开始逐个取当前 goroutine 的 defer 调用
		d := gp._defer
		// 没有defer，退出循环
		if d == nil {
			break
		}

		// 如果 defer 是由早期的 panic 或 Goexit 开始的（并且，因为我们回到这里，这引发了新的 panic），
		// 则将 defer 带离链表。更早的 panic 或 Goexit 将无法继续运行。
		if d.started {
			if d._panic != nil {
				d._panic.aborted = true
			}
			d._panic = nil
			d.fn = nil
			gp._defer = d.link
			freedefer(d)
			continue
		}

		// 将deferred标记为started
		// 如果栈增长或者垃圾回收在 reflectcall 开始执行 d.fn 前发生
		// 标记 defer 已经开始执行，但仍将其保存在列表中，从而 traceback 可以找到并更新这个 defer 的参数帧
		d.started = true

		// 记录正在运行的延迟的panic。
		// 如果在延迟调用期间有新的panic，那么这个panic
		// 将在列表中找到d，并将标记d._panic(此panic)中止。
		d._panic = (*_panic)(noescape(unsafe.Pointer(&p)))

		p.argp = unsafe.Pointer(getargp(0))
        // 通过reflectcall函数调用defered函数
        // 如果defered函数再次发生panic而且并未被该defered函数recover，则reflectcall永远不会返回
        // 如果defered函数并没有发生过panic或者发生了panic但该defered函数成功recover了新发生的panic，
        // 则此函数会返回继续执行后面的代码。
		reflectcall(nil, unsafe.Pointer(d.fn), deferArgs(d), uint32(d.siz), uint32(d.siz))
		p.argp = nil

		// reflectcall没有panic。删除d
		if gp._defer != d {
			throw("bad defer entry in panic")
		}
		d._panic = nil
		d.fn = nil
		gp._defer = d.link

		// trigger shrinkage to test stack copy. See stack_test.go:TestStackPanic
		//GC()

		pc := d.pc
		sp := unsafe.Pointer(d.sp) // must be pointer so it gets adjusted during stack copy
		freedefer(d)
		if p.recovered {
			atomic.Xadd(&runningPanicDefers, -1)

			gp._panic = p.link
			// 忽略的 panic 会被标记，但仍然保留在 g.panic 列表中
			// 这里将它们移出列表
			for gp._panic != nil && gp._panic.aborted {
				gp._panic = gp._panic.link
			}
			if gp._panic == nil { // 必须由 signal 完成
				gp.sig = 0
			}
			// 传递关于恢复帧的信息
			gp.sigcode0 = uintptr(sp)
			gp.sigcode1 = pc
			// 调用 recover，并重新进入调度循环，不再返回
			mcall(recovery)
			// 如果无法重新进入调度循环，则无法恢复错误
			throw("recovery failed") // mcall should not return
		}
	}

	// 消耗完所有的 defer 调用，保守地进行 panic
	// 因为在冻结之后调用任意用户代码是不安全的，所以我们调用 preprintpanics 来调用
	// 所有必要的 Error 和 String 方法来在 startpanic 之前准备 panic 字符串。
	preprintpanics(gp._panic)

	fatalpanic(gp._panic) // 不应该返回
	*(*int)(nil) = 0      // 无法触及
}

// reflectcall 使用 arg 指向的 n 个参数字节的副本调用 fn。
// fn 返回后，reflectcall 在返回之前将 n-retoffset 结果字节复制回 arg+retoffset。
// 如果重新复制结果字节，则调用者应将参数帧类型作为 argtype 传递，以便该调用可以在复制期间执行适当的写障碍。
// reflect 包传递帧类型。在 runtime 包中，只有一个调用将结果复制回来，即 cgocallbackg1，
// 并且它不传递帧类型，这意味着没有调用写障碍。参见该调用的页面了解相关理由。
//
// 包 reflect 通过 linkname 访问此符号
func reflectcall(argtype *_type, fn, arg unsafe.Pointer, argsize uint32, retoffset uint32)
```

梳理下流程  

1、在处理`panic`期间，会先判断当前`panic`的类型，确定`panic`是否可恢复;  

- 系统栈上的panic无法恢复
- 如果正在进行malloc时发生panic也无法恢复
- 在禁止抢占时发生panic也无法恢复
- 在g锁在m上时发生panic也无法恢复

2、可恢复的`panic`，`panic`的`link`指向`goroutine`链表中先前的`panic`链表；   

3、循环逐个获取当前`goroutine`的`defer`调用；  

- 如果defer是由早期panic或Goexit开始的，则将defer带离链表，更早的panic或Goexit将无法继续运行；  

- 

- 

### 参考

【panic 和 recover】https://draveness.me/golang/docs/part2-foundation/ch05-keyword/golang-panic-recover/  
【恐慌与恢复内建函数】https://golang.design/under-the-hood/zh-cn/part1basic/ch03lang/panic/  
【Go语言panic/recover的实现】https://zhuanlan.zhihu.com/p/72779197  
【panic and recover】https://eddycjy.gitbook.io/golang/di-6-ke-chang-yong-guan-jian-zi/panic-and-recover#yuan-ma  