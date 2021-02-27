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