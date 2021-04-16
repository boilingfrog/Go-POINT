## 错误使用map引发的血案

### 前言  

最近业务中，同事使用`map`来接收返回的结果，使用`waitGroup`来并发的处理执行返回的结果，结果上线之后，直接崩了。  

<img src="/img/map_1.jpg"  alt="map" align=center />

日志大量的数据库缓存池连接失败

```go
{"ecode":-500,"message":"timed out while checking out a connection from connection pool"}

{"ecode":-500,"message":"connection(xxxxxxxxx:xxxxx) failed to write: context deadline exceeded"}
```

### 场景复原

先来看来伪代码

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

var count = 300

func main() {
	var data = make(map[int]string, count)
	var wg sync.WaitGroup

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			time.Sleep(time.Second * 1)
			mockSqlPool()
			data[i] = "test"
		}(i)
	}
	fmt.Println("-----------WaitGroup执行结束了-----------")
	wg.Wait()
}

// 模拟数据库的连接和释放
func mockSqlPool() {
	defer fmt.Println("关闭pool")
	fmt.Println("我是pool")
}
```

运行的输出  

```go
...
我是pool
关闭pool
我是pool
fatal error: 关闭pool
concurrent map writes
我是pool

goroutine 56 [running]:
runtime.throw(0x10d3923, 0x15)
        /usr/local/go/src/runtime/panic.go:774 +0x72 fp=0xc00023cf20 sp=0xc00023cef0 pc=0x10298d2
runtime.mapassign_fast64(0x10b29e0, 0xc000066180, 0x16, 0x0)
        /usr/local/go/src/runtime/map_fast64.go:101 +0x350 fp=0xc00023cf60 sp=0xc00023cf20 pc=0x100f620
main.main.func1(0xc00008c004, 0xc000066180, 0x16)
        /Users/yj/Go/src/Go-POINT/map/main.go:23 +0x87 fp=0xc00023cfc8 sp=0xc00023cf60 pc=0x109a297
runtime.goexit()
        /usr/local/go/src/runtime/asm_amd64.s:1357 +0x1 fp=0xc00023cfd0 sp=0xc00023cfc8 pc=0x1053a51
created by main.main
        /Users/yj/Go/src/Go-POINT/map/main.go:18 +0xbb

goroutine 1 [semacquire]:
sync.runtime_Semacquire(0xc00008c004)
        /usr/local/go/src/runtime/sema.go:56 +0x42
sync.(*WaitGroup).Wait(0xc00008c004)
        /usr/local/go/src/sync/waitgroup.go:130 +0x64
main.main()
        /Users/yj/Go/src/Go-POINT/map/main.go:27 +0x138

goroutine 22 [semacquire]:
internal/poll.runtime_Semacquire(0xc00008606c)
        /usr/local/go/src/runtime/sema.go:61 +0x42
internal/poll.(*fdMutex).rwlock(0xc000086060, 0xc000030500, 0x1097137)
        /usr/local/go/src/internal/poll/fd_mutex.go:154 +0xad
internal/poll.(*FD).writeLock(...)
        /usr/local/go/src/internal/poll/fd_mutex.go:239
internal/poll.(*FD).Write(0xc000086060, 0xc000226030, 0xb, 0x10, 0x0, 0x0, 0x0)
        /usr/local/go/src/internal/poll/fd_unix.go:255 +0x5e
os.(*File).write(...)
        /usr/local/go/src/os/file_unix.go:276
os.(*File).Write(0xc000084008, 0xc000226030, 0xb, 0x10, 0xc0000306b0, 0x103d37e, 0xc00000c060)
        /usr/local/go/src/os/file.go:153 +0x77
fmt.Fprintln(0x10ec4e0, 0xc000084008, 0xc000030730, 0x1, 0x1, 0x10459b6, 0xc00000c060, 0x3)
        /usr/local/go/src/fmt/print.go:265 +0x8b
fmt.Println(...)
        /usr/local/go/src/fmt/print.go:274
main.mockSqlPool()
        /Users/yj/Go/src/Go-POINT/map/main.go:35 +0x104
main.main.func1(0xc00008c004, 0xc000066180, 0x4)
        /Users/yj/Go/src/Go-POINT/map/main.go:21 +0x63
created by main.main
        /Users/yj/Go/src/Go-POINT/map/main.go:18 +0xbb
...
```

一个全局的`map`，然后`WaitGroup`开启一组协程并发的读写数据，写入内容到`map`中。  

#### WaitGroup中某个goroutine发生panic会如何？

```go
func waitGroupPanic() {
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		time.Sleep(time.Second)
		panic("just panic")
	}()

	go func() {
		defer wg.Done()
		fmt.Println("run 1")
	}()

	go func() {
		defer wg.Done()
		time.Sleep(time.Second * 2)
		fmt.Println("run 2")
	}()

	fmt.Println("执行了吗")
	wg.Wait()
	fmt.Println("完美退出了")
}
```

对于`panic`能够改变程序的控制流，调用`panic`后会立刻停止执行当前函数的剩余代码，并在当前Goroutine中递归执行调用方的`defer`  

上面的测试代码中，当第一个`goroutine`发生`panic`的时候，`panic`会向上传递，直到`wg.Wait()`导致整个wg.Wait()

上面的错误原因有两个：  

- 1、`map`不是并发安全，并发写的时候会触发`panic`；  

- 2、避免在循环中连接梳理数据库； 

因为map的并发写入，导致触发了`panic`，对于`panic`,当`panic`异常发生时，程序会中断后面代码的执行，
