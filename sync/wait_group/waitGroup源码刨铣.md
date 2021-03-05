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

**state1**

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
| 32位     | count     | wait      | semaphore |
| 32位     | semaphore | count     | wait      |

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

#### Add

```go
func (wg *WaitGroup) Add(delta int) {
	statep, semap := wg.state()
	if race.Enabled {
		_ = *statep // trigger nil deref early
		if delta < 0 {
			// Synchronize decrements with Wait.
			race.ReleaseMerge(unsafe.Pointer(wg))
		}
		race.Disable()
		defer race.Enable()
	}
	state := atomic.AddUint64(statep, uint64(delta)<<32)
	v := int32(state >> 32)
	w := uint32(state)
	if race.Enabled && delta > 0 && v == int32(delta) {
		// The first increment must be synchronized with Wait.
		// Need to model this as a read, because there can be
		// several concurrent wg.counter transitions from 0.
		race.Read(unsafe.Pointer(semap))
	}
	if v < 0 {
		panic("sync: negative WaitGroup counter")
	}
	if w != 0 && delta > 0 && v == int32(delta) {
		panic("sync: WaitGroup misuse: Add called concurrently with Wait")
	}
	if v > 0 || w == 0 {
		return
	}
	// This goroutine has set counter to 0 when waiters > 0.
	// Now there can't be concurrent mutations of state:
	// - Adds must not happen concurrently with Wait,
	// - Wait does not increment waiters if it sees counter == 0.
	// Still do a cheap sanity check to detect WaitGroup misuse.
	if *statep != state {
		panic("sync: WaitGroup misuse: Add called concurrently with Wait")
	}
	// Reset waiters count to 0.
	*statep = 0
	for ; w != 0; w-- {
		runtime_Semrelease(semap, false, 0)
	}
}
```





### 参考

【《Go专家编程》Go WaitGroup实现原理】https://my.oschina.net/renhc/blog/2249061  
【Go中由WaitGroup引发对内存对齐思考】https://cloud.tencent.com/developer/article/1776930  
【Golang 之 WaitGroup 源码解析】https://www.linkinstar.wiki/2020/03/15/golang/source-code/sync-waitgroup-source-code/  
【sync.WaitGroup】https://golang.design/under-the-hood/zh-cn/part1basic/ch05sync/waitgroup/    