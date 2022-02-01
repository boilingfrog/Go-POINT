<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [错误使用map引发的血案](#%E9%94%99%E8%AF%AF%E4%BD%BF%E7%94%A8map%E5%BC%95%E5%8F%91%E7%9A%84%E8%A1%80%E6%A1%88)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [场景复原](#%E5%9C%BA%E6%99%AF%E5%A4%8D%E5%8E%9F)
  - [原因](#%E5%8E%9F%E5%9B%A0)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 错误使用map引发的血案

### 前言  

最近业务中，同事使用`map`来接收返回的结果，使用`waitGroup`来并发的处理执行返回的结果，结果上线之后，直接崩了。  

<img src="/img/golang/map_1.jpg"  alt="map" align=center />

日志大量的数据库缓存池连接失败

```go
{"ecode":-500,"message":"timed out while checking out a connection from connection pool"}

{"ecode":-500,"message":"connection(xxxxxxxxx:xxxxx) failed to write: context deadline exceeded"}
```

### 场景复原

先来看来伪代码  

一个全局的`map`，然后`WaitGroup`开启一组协程并发的读写数据，写入内容到`map`中。  

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

goroutine 192 [semacquire]:
internal/poll.runtime_Semacquire(0xc00009e06c)
        /usr/local/go/src/runtime/sema.go:61 +0x42
internal/poll.(*fdMutex).rwlock(0xc00009e060, 0x10fae00, 0xc00023ad00)
        /usr/local/go/src/internal/poll/fd_mutex.go:154 +0xe9
internal/poll.(*FD).writeLock(...)
        /usr/local/go/src/internal/poll/fd_mutex.go:239
internal/poll.(*FD).Write(0xc00009e060, 0xc000246100, 0xb, 0x10, 0x0, 0x0, 0x0)
        /usr/local/go/src/internal/poll/fd_unix.go:255 +0x6f
os.(*File).write(...)
        /usr/local/go/src/os/file_unix.go:276
os.(*File).Write(0xc00009c008, 0xc000246100, 0xb, 0x10, 0xc000124580, 0x40, 0x0)
        /usr/local/go/src/os/file.go:153 +0xa7
fmt.Fprintln(0x1158520, 0xc00009c008, 0xc00014d728, 0x1, 0x1, 0x107e3e6, 0xc0000d8100, 0x16)
        /usr/local/go/src/fmt/print.go:265 +0xb3
fmt.Println(...)
        /usr/local/go/src/fmt/print.go:274
main.mockSqlPool()
        /Users/yj/Go/src/Go-POINT/map/main.go:35 +0x129
main.main.func1(0xc0000a0004, 0xc000088180, 0x8f)
        /Users/yj/Go/src/Go-POINT/map/main.go:21 +0x75
created by main.main
        /Users/yj/Go/src/Go-POINT/map/main.go:18 +0x102

goroutine 193 [semacquire]:
internal/poll.runtime_Semacquire(0xc00009e06c)
        /usr/local/go/src/runtime/sema.go:61 +0x42
internal/poll.(*fdMutex).rwlock(0xc00009e060, 0x10fae00, 0xc000286410)
        /usr/local/go/src/internal/poll/fd_mutex.go:154 +0xe9
internal/poll.(*FD).writeLock(...)
        /usr/local/go/src/internal/poll/fd_mutex.go:239
internal/poll.(*FD).Write(0xc00009e060, 0xc0000a01a0, 0xb, 0x10, 0x0, 0x0, 0x0)
        /usr/local/go/src/internal/poll/fd_unix.go:255 +0x6f
os.(*File).write(...)
        /usr/local/go/src/os/file_unix.go:276
os.(*File).Write(0xc00009c008, 0xc0000a01a0, 0xb, 0x10, 0xc0001245c0, 0x40, 0x0)
        /usr/local/go/src/os/file.go:153 +0xa7
fmt.Fprintln(0x1158520, 0xc00009c008, 0xc00014df28, 0x1, 0x1, 0x107e3e6, 0xc0000d8100, 0x17)
        /usr/local/go/src/fmt/print.go:265 +0xb3
fmt.Println(...)
        /usr/local/go/src/fmt/print.go:274
main.mockSqlPool()
        /Users/yj/Go/src/Go-POINT/map/main.go:35 +0x129
main.main.func1(0xc0000a0004, 0xc000088180, 0x90)
        /Users/yj/Go/src/Go-POINT/map/main.go:21 +0x75
created by main.main
        /Users/yj/Go/src/Go-POINT/map/main.go:18 +0x102

goroutine 194 [semacquire]:
internal/poll.runtime_Semacquire(0xc00009e06c)
        /usr/local/go/src/runtime/sema.go:61 +0x42
internal/poll.(*fdMutex).rwlock(0xc00009e060, 0x10fae00, 0xc00023add0)
        /usr/local/go/src/internal/poll/fd_mutex.go:154 +0xe9
internal/poll.(*FD).writeLock(...)
        /usr/local/go/src/internal/poll/fd_mutex.go:239
internal/poll.(*FD).Write(0xc00009e060, 0xc000246110, 0xb, 0x10, 0x0, 0x0, 0x0)
        /usr/local/go/src/internal/poll/fd_unix.go:255 +0x6f
os.(*File).write(...)
        /usr/local/go/src/os/file_unix.go:276
os.(*File).Write(0xc00009c008, 0xc000246110, 0xb, 0x10, 0xc000124600, 0x40, 0x0)
        /usr/local/go/src/os/file.go:153 +0xa7
fmt.Fprintln(0x1158520, 0xc00009c008, 0xc000146728, 0x1, 0x1, 0x107e3e6, 0xc0000d8100, 0x18)
        /usr/local/go/src/fmt/print.go:265 +0xb3
fmt.Println(...)
        /usr/local/go/src/fmt/print.go:274
main.mockSqlPool()
        /Users/yj/Go/src/Go-POINT/map/main.go:35 +0x129
main.main.func1(0xc0000a0004, 0xc000088180, 0x91)
        /Users/yj/Go/src/Go-POINT/map/main.go:21 +0x75
created by main.main
        /Users/yj/Go/src/Go-POINT/map/main.go:18 +0x102

```

会发现很多`goroutine`处于`semacquire`状态，说明这些`goroutine`正在等待被信号量唤醒。但是这时候`waitGroup`已经因为`panic`退出了，这些`goroutine`不会在通过`waitGroup.Done()`退出，造成这些`goroutine`一直阻塞到这，最后的结果就是这些`goroutine`占用的数据库连接不能被释放。  

关于`waitGroup`的信号量  

整个`Wait()`会被`runtime_Semacquire`阻塞，直到等到全部退出的信号量；

`Done()`会在最后一次的时候通过`runtime_Semrelease`发出取消阻塞的信号，然后被`runtime_Semacquire`阻塞的`Wait()`就可以退出了；  

上面涉及到的几种状态  

- semacquire 状态，这个状态表示等待调用  
- Waiting 等待状态。线程在等待某件事的发生。例如等待网络数据、硬盘；调用操作系统 API；等待内存同步访问条件 ready，如 atomic, mutexes
- Runnable 就绪状态。只要给 CPU 资源我就能运行

### 原因

上面的错误原因有两个：  

- 1、`map`不是并发安全，并发写的时候会触发`panic`；  

- 2、避免在循环中连接数据库； 

### 参考

【map 并发崩溃一例】https://xargin.com/map-concurrent-throw/  