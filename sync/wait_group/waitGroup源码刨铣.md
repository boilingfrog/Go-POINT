<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [waitGroup源码刨铣](#waitgroup%E6%BA%90%E7%A0%81%E5%88%A8%E9%93%A3)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [WaitGroup实现](#waitgroup%E5%AE%9E%E7%8E%B0)
    - [noCopy](#nocopy)
    - [state1](#state1)
  - [Add](#add)
  - [Wait](#wait)
  - [总结](#%E6%80%BB%E7%BB%93)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## waitGroup源码刨铣

### 前言

学习下waitGroup的实现  

本文是在`go version go1.13.15 darwin/amd64`上进行的  

### WaitGroup实现

看一个小demo  

```go
func waitGroup() {
	var wg sync.WaitGroup

	wg.Add(4)
	go func() {
		defer wg.Done()
		fmt.Println(1)
	}()

	go func() {
		defer wg.Done()
		fmt.Println(2)
	}()

	go func() {
		defer wg.Done()
		fmt.Println(3)
	}()

	go func() {
		defer wg.Done()
		fmt.Println(4)
	}()

	wg.Wait()
	fmt.Println("1 2 3 4 end")
}
```

1、启动goroutine前将计数器通过`Add(4)`将计数器设置为待启动的`goroutine`个数。  
2、启动`goroutine`后，使用`Wait()`方法阻塞主协程，等待计数器变为0。  
3、每个`goroutine`执行结束通过`Done()`方法将计数器减1。  
4、计数器变为0后，阻塞的`goroutine`被唤醒。   

看下具体的实现  

```go
// WaitGroup 不能被copy
type WaitGroup struct {
	noCopy noCopy

	state1 [3]uint32
}
```

#### noCopy

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

#### state1

```go
	// 64 位值: 高 32 位用于计数，低 32 位用于等待计数
	// 64 位的原子操作要求 64 位对齐，但 32 位编译器无法保证这个要求
	// 因此分配 12 字节然后将他们对齐，其中 8 字节作为状态，其他 4 字节用于存储原语
	state1 [3]uint32
```

这点是`wait_group`很巧妙的一点，大神写代码的思路就是惊奇  

这个设计很奇妙，通过内存对齐来处理`wait_group`中的waiter数、计数值、信号量。什么是内存对齐可参考[什么是内存对齐，go中内存对齐分析](https://www.cnblogs.com/ricklz/p/14455135.html)  

来分析下`state1`是如何内存对齐来处理几个计数值的存储  

计算机为了加快内存的访问速度，会对内存进行对齐处理，CPU把内存当作一块一块的，块的大小可以是2、4、8、16字节大小，因此CPU读取内存是一块一块读取的。  

合理的内存对齐可以提高内存读写的性能，并且便于实现变量操作的原子性。  

在不同平台上的编译器都有自己默认的 “对齐系数”，可通过预编译命令#pragma pack(n)进行变更，n就是代指 “对齐系数”。  

一般来讲，我们常用的平台的系数如下：  

- 32 位：4

- 64 位：8

`state1`这块就兼容了两种平台的对齐系数  

对于64未系统来讲。内存访问的步长是8。也就是cpu一次访问8位偏移量的内存空间。当时对于32未的系统，内存的对齐系数是4，也就是访问的步长是4个偏移量。  

所以为了兼容这两种模式，这里采用了`uint32`结构的数组，保证在不同类型的机器上都是`12`个字节，一个`uint32`是4字节。这样对于32位的4步长访问是没有问题了，64位的好像也没有解决，8步长的访问会一次读入两个`uint32`的长度。  

所以，下面的读取也进行了操作，将两个`uint32`的内存放到一个`uint64`中返回，这样就同时解决了32位和64位访问步长的问题。  

所以，64位系统和32位系统，`state1`中`counter，waiter，semaphore`的内存布局是不一样的。  

|          | state [0] | state [1] | state [2] |
| :------: | :------:  | :------:  | :------:  |
| 32位     | waiter    | counter   | semaphore |
| 32位     | semaphore | waiter    | counter     |

`counter`位于高地址位，`waiter`位于地址位  

<img src="/img/waitgroup_state1_1.png" width = "455" height = "375.5" alt="waitgroup" align=center />

下面是state的代码  
```go
func (wg *WaitGroup) state() (statep *uint64, semap *uint32) {
	// 判断是否是64位
	if uintptr(unsafe.Pointer(&wg.state1))%8 == 0 {
		return (*uint64)(unsafe.Pointer(&wg.state1)), &wg.state1[2]
	} else {
		return (*uint64)(unsafe.Pointer(&wg.state1[1])), &wg.state1[0]
	}
}
```

对于count和wait在高低地址位的体现，在add中的代码可体现

```go
        // // 将 delta 加到 statep 的前 32 位上，即加到计数器上
	state := atomic.AddUint64(statep, uint64(delta)<<32)
	v := int32(state >> 32) // count
	w := uint32(state)  // wait
```

这里面也用到了信号量  

信号量是Unix系统提供的一种保护共享资源的机制，用于防止多个线程同时访问某个资源。  

可简单理解为信号量为一个数值：  

- 当信号量>0时，表示资源可用，获取信号量时系统自动将信号量减1；  

- 当信号量==0时，表示资源暂不可用，获取信号量时，当前线程会进入睡眠，当信号量为正时被唤醒。

`WaitGroup`中的实现就用到了这个，在下面的代码实现就能看到  

### Add

```go
// Add将增量（可能为负）添加到WaitGroup计数器中。
// 如果计数器为零，则释放等待时阻塞的所有goroutine。
// 如果计数器变为负数，请添加恐慌。
//
// 请注意，当计数器为 0 时发生的带有正的 delta 的调用必须在 Wait 之前。
// 当计数器大于 0 时，带有负 delta 的调用或带有正 delta 调用可能在任何时候发生。
// 通常，这意味着对Add的调用应在语句之前执行创建要等待的goroutine或其他事件。
// 如果将WaitGroup重用于等待几个独立的事件集，新的Add调用必须在所有先前的Wait调用返回之后发生。
func (wg *WaitGroup) Add(delta int) {
	// 获取counter,waiter,以及semaphore对应的指针
	statep, semap := wg.state()
	...
	// 将 delta 加到 statep 的前 32 位上，即加到计数器上
	state := atomic.AddUint64(statep, uint64(delta)<<32)
	// 高地址位counter
	v := int32(state >> 32)
	// 低地址为waiter
	w := uint32(state)
	...
	// 计数器不允许为负数
	if v < 0 {
		panic("sync: negative WaitGroup counter")
	}
	// wait不等于0说明已经执行了Wait，此时不容许Add
	if w != 0 && delta > 0 && v == int32(delta) {
		panic("sync: WaitGroup misuse: Add called concurrently with Wait")
	}
	// 计数器的值大于或者没有waiter在等待,直接返回
	if v > 0 || w == 0 {
		return
	}
	// 运行到这里只有一种情况 v == 0 && w != 0

	// 这时 Goroutine 已经将计数器清零，且等待器大于零（并发调用导致）
	// 这时不允许出现并发使用导致的状态突变，否则就应该 panic
	// - Add 不能与 Wait 并发调用
	// - Wait 在计数器已经归零的情况下，不能再继续增加等待器了
	// 仍然检查来保证 WaitGroup 不会被滥用

	// 这一点很重要，这段代码同时也保证了这是最后的一个需要等待阻塞的goroutine
	// 然后在下面通过runtime_Semrelease，唤醒被信号量semap阻塞的waiter
	if *statep != state {
		panic("sync: WaitGroup misuse: Add called concurrently with Wait")
	}
	// 结束后将等待器清零
	*statep = 0
	for ; w != 0; w-- {
		// 释放信号量，通过runtime_Semacquire唤醒被阻塞的waiter
		runtime_Semrelease(semap, false, 0)
	}
}
```

梳理下流程  

1、首先获取存储在`state1`中对应的几个变量的指针；  

2、`counter`存储在高位，增加的时候需要左移32位；   

3、counter的数量不能小于0，小于0抛出panic;  

4、同样也会判断，已经的执行`wait`之后，不能在增加`counter`;  

5、（这点很重要，我自己看了好久才明白）计数器的值大于或者没有`waiter`在等待,直接返回  

```go
	// 计数器的值大于或者没有waiter在等待,直接返回
	if v > 0 || w == 0 {
		return
	}
```

因为`waiter`的值只会被执行一次+1操作，所以这段代码保证了只有在`v == 0 && w != 0`，也就是最后一个`Done()`操作的时候，走到下面的代码，释放信号量，唤醒被信号量阻塞的`Wait()`，结束整个`WaitGroup`。   

### Wait

```go
// Wait blocks until the WaitGroup counter is zero.
func (wg *WaitGroup) Wait() {
	// 获取counter,waiter,以及semaphore对应的指针
	statep, semap := wg.state()
	...
	for {
		// 获取对应的counter和waiter数量
		state := atomic.LoadUint64(statep)
		v := int32(state >> 32)
		w := uint32(state)
		// Counter为0，不需要等待
		if v == 0 {
			if race.Enabled {
				race.Enable()
				race.Acquire(unsafe.Pointer(wg))
			}
			return
		}
		// 原子（cas）增加waiter的数量（只会被+1操作一次）
		if atomic.CompareAndSwapUint64(statep, state, state+1) {
			...
			// 这块用到了，我们上文讲的那个信号量
			// 等待被runtime_Semrelease释放的信号量唤醒
			// 如果 *semap > 0 则会减 1,等于0则被阻塞
			runtime_Semacquire(semap)

			// 在这种情况下，如果 *statep 不等于 0 ，则说明使用失误，直接 panic
			if *statep != 0 {
				panic("sync: WaitGroup is reused before previous Wait has returned")
			}
			...
			return
		}
	}
}
```

梳理下流程  

1、首先获取存储在`state1`中对应的几个变量的指针；  

2、一个for循环，来阻塞等待所有的`goroutine`退出；  

3、如果`counter`为0，不需要等待，直接退出即可；  

4、原子（cas）增加`waiter`的数量（只会被+1操作一次）；  

5、整个`Wait()`会被`runtime_Semacquire`阻塞，直到等到退出的信号量；  

6、`Done()`会在最后一次的时候通过`runtime_Semrelease`发出取消阻塞的信号，然后被`runtime_Semacquire`阻塞的`Wait()`就可以退出了；  

7、整个`WaitGroup`执行成功。  

<img src="/img/waitgroup_all_1.png" width = "501" height = "438" alt="waitgroup" align=center />

### 总结

代码中我感到设计比较巧妙的有两个部分：  

1、`state1`的处理，保证内存对齐，设置高低位内存来存储不同的值，同时32位和64位平台的处理方式还不同；  

2、信号量的阻塞退出，这块最后一个`Done`退出的时候，才会触发阻塞信号量，退出`Wait()`，然后结束整个`waitGroup`。再此之前，当Wait()在成功将waiter变量+1操作之后，就会被`runtime_Semacquire`阻塞，直到最后一个`Done`，信号的发出。  

对于`WaitGroup`的使用  

1、计数器的值不能为负数，可能是`Add(-1)`触发的，也可能是`Done()`触发的，否则会panic;  

2、`Add`数量的添加，要发生在`Wait()`之前；  

3、`WaitGroup`是可以重用的，但是需要等上一批的`goroutine` 都调用`Wait`完毕后才能继续重用`WaitGroup`。  

### 参考

【《Go专家编程》Go WaitGroup实现原理】https://my.oschina.net/renhc/blog/2249061  
【Go中由WaitGroup引发对内存对齐思考】https://cloud.tencent.com/developer/article/1776930  
【Golang 之 WaitGroup 源码解析】https://www.linkinstar.wiki/2020/03/15/golang/source-code/sync-waitgroup-source-code/  
【sync.WaitGroup】https://golang.design/under-the-hood/zh-cn/part1basic/ch05sync/waitgroup/    