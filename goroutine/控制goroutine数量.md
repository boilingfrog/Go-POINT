<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [控制goroutine数量](#%E6%8E%A7%E5%88%B6goroutine%E6%95%B0%E9%87%8F)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [控制goroutine的数量](#%E6%8E%A7%E5%88%B6goroutine%E7%9A%84%E6%95%B0%E9%87%8F)
    - [通过channel+sync](#%E9%80%9A%E8%BF%87channelsync)
    - [使用semaphore](#%E4%BD%BF%E7%94%A8semaphore)
    - [线程池](#%E7%BA%BF%E7%A8%8B%E6%B1%A0)
  - [几个开源的线程池的设计](#%E5%87%A0%E4%B8%AA%E5%BC%80%E6%BA%90%E7%9A%84%E7%BA%BF%E7%A8%8B%E6%B1%A0%E7%9A%84%E8%AE%BE%E8%AE%A1)
    - [fasthttp中的协程池实现](#fasthttp%E4%B8%AD%E7%9A%84%E5%8D%8F%E7%A8%8B%E6%B1%A0%E5%AE%9E%E7%8E%B0)
      - [Start](#start)
      - [Stop](#stop)
      - [clean](#clean)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 控制goroutine数量

### 前言

`goroutine`被无限制的大量创建，造成的后果就不啰嗦了，主要讨论几种如何控制`goroutine`的方法  

### 控制goroutine的数量


#### 通过channel+sync

```go
var (
	// channel长度
	poolCount      = 5
	// 复用的goroutine数量
	goroutineCount = 10
)

func pool() {
	jobsChan := make(chan int, poolCount)

	// workers
	var wg sync.WaitGroup
	for i := 0; i < goroutineCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range jobsChan {
				// ...
				fmt.Println(item)
			}
		}()
	}

	// senders
	for i := 0; i < 1000; i++ {
		jobsChan <- i
	}

	// 关闭channel，上游的goroutine在读完channel的内容，就会通过wg的done退出
	close(jobsChan)
	wg.Wait()
}
```

通过`WaitGroup`启动指定数量的`goroutine`，监听`channel`的通知。发送者推送信息到`channel`，信息处理完了，关闭`channel`,等待`goroutine`依次退出。  

#### 使用semaphore

```go
package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

const (
	// 同时运行的goroutine上限
	Limit = 3
	// 信号量的权重
	Weight = 1
)

func main() {
	names := []string{
		"小白",
		"小红",
		"小明",
		"小李",
		"小花",
	}

	sem := semaphore.NewWeighted(Limit)
	var w sync.WaitGroup
	for _, name := range names {
		w.Add(1)
		go func(name string) {
			sem.Acquire(context.Background(), Weight)
			// ... 具体的业务逻辑
			fmt.Println(name, "-吃饭了")
			time.Sleep(2 * time.Second)
			sem.Release(Weight)
			w.Done()
		}(name)
	}
	w.Wait()

	fmt.Println("ending--------")
}
```

借助于x包中的`semaphore`，也可以进行`goroutine`的数量限制。  

#### 线程池

不过原本go中的协程已经是非常轻量了，对于协程池还是要根据具体的场景分析。   

对于小场景使用`channel+sync`就可以，其他复杂的可以考虑使用第三方的协程池库。  

- [panjf2000/ants](https://github.com/panjf2000/ants)

- [go-playground/pool](https://github.com/go-playground/pool)

- [Jeffail/tunny](https://github.com/Jeffail/tunny)  

### 几个开源的线程池的设计

#### fasthttp中的协程池实现  

`fasthttp`比`net/http`效率高很多倍的重要原因，就是利用了协程池。来看下大佬的设计思路。   

看下实现  

```go
// workerPool通过一组工作池服务传入的连接
// 按照FILO（先进后出）的顺序，即最近停止的工作人员将为下一个工作传入的连接。
//
// 这种方案能够保持cpu的缓存保持高效（理论上）
type workerPool struct {
	// 这个函数用于server的连接
	// It must leave c unclosed.
	WorkerFunc ServeHandler

	MaxWorkersCount int

	LogAllErrors bool

	MaxIdleWorkerDuration time.Duration

	Logger Logger

	lock         sync.Mutex
	workersCount int
	mustStop     bool

	ready []*workerChan

	stopCh chan struct{}

	// workerChan的缓存池，通过sync.Pool实现
	workerChanPool sync.Pool

	connState func(net.Conn, ConnState)
}

// workerChan的结构
type workerChan struct {
	lastUseTime time.Time
	ch          chan net.Conn
}
```

##### Start

```go
func (wp *workerPool) Start() {
	// 判断是否已经Start过了
	if wp.stopCh != nil {
		panic("BUG: workerPool already started")
	}
	// stopCh塞入值
	wp.stopCh = make(chan struct{})
	stopCh := wp.stopCh
	wp.workerChanPool.New = func() interface{} {
		// 如果单核cpu则让workerChan阻塞
		// 否则，使用非阻塞，workerChan的长度为1
		return &workerChan{
			ch: make(chan net.Conn, workerChanCap),
		}
	}
	go func() {
		var scratch []*workerChan
		for {
			wp.clean(&scratch)
			select {
            // 接收到退出信号，退出
			case <-stopCh:
				return
			default:
				time.Sleep(wp.getMaxIdleWorkerDuration())
			}
		}
	}()
}

// 如果单核cpu则让workerChan阻塞
// 否则，使用非阻塞，workerChan的长度为1
var workerChanCap = func() int {
	// 如果GOMAXPROCS=1，workerChan的长度为0，变成一个阻塞的channel
	if runtime.GOMAXPROCS(0) == 1 {
		return 0
	}

	// 如果GOMAXPROCS>1则使用非阻塞的workerChan
	return 1
}()
```

梳理下流程：  

1、首先判断下`stopCh`是否为`nil`，不为`nil`表示已经`started`了；  

2、初始化`wp.stopCh = make(chan struct{})`，`stopCh`是一个标识，用了`struct{}`不用`bool`，因为空结构体变量的内存占用大小为0，而`bool`类型内存占用大小为1，这样可以更加最大化利用我们服务器的内存空间；  

3、设置`workerChanPool`的`New`函数，然后可以在`Get`不到东西时，自动创建一个；如果单核`cpu`则让`workerChan`阻塞，否则，使用非阻塞，`workerChan`的长度设置为1；  

4、启动一个`goroutine`，处理`clean`操作，在接收到退出信号，退出。  

##### Stop

```go
func (wp *workerPool) Stop() {
	// 同start，stop也只能触发一次
	if wp.stopCh == nil {
		panic("BUG: workerPool wasn't started")
	}
	// 关闭stopCh
	close(wp.stopCh)
	// 将stopCh置为nil
	wp.stopCh = nil

	// 停止所有的等待获取连接的workers
	// 正在运行的workers，不需要等待他们退出，他们会在完成connection或mustStop被设置成true退出
	wp.lock.Lock()
	ready := wp.ready
	// 循环将ready的workerChan置为nil
	for i := range ready {
		ready[i].ch <- nil
		ready[i] = nil
	}
	wp.ready = ready[:0]
	// 设置mustStop为true
	wp.mustStop = true
	wp.lock.Unlock()
}
```

梳理下流程：  

1、判断stop只能被关闭一次；  

2、关闭`stopCh`，设置`stopCh`为`nil`；  

3、停止所有的等待获取连接的`workers`，正在运行的`workers`，不需要等待他们退出，他们会在完成`connection`或`mustStop`被设置成`true`退出。  

##### clean

```go
func (wp *workerPool) clean(scratch *[]*workerChan) {
	maxIdleWorkerDuration := wp.getMaxIdleWorkerDuration()

    // 清理掉最近最少使用的workers如果他们过了maxIdleWorkerDuration时间没有提供服务
	criticalTime := time.Now().Add(-maxIdleWorkerDuration)

	wp.lock.Lock()
	ready := wp.ready
	n := len(ready)

    // 使用二分搜索算法找出最近可以被清除的worker
    // 最后使用的workerChan 一定是放回队列尾部的。
	l, r, mid := 0, n-1, 0
	for l <= r {
		mid = (l + r) / 2
		if criticalTime.After(wp.ready[mid].lastUseTime) {
			l = mid + 1
		} else {
			r = mid - 1
		}
	}
	i := r
	if i == -1 {
		wp.lock.Unlock()
		return
	}

    // 将ready中i之前的的全部清除
	*scratch = append((*scratch)[:0], ready[:i+1]...)
	m := copy(ready, ready[i+1:])
	for i = m; i < n; i++ {
		ready[i] = nil
	}
	wp.ready = ready[:m]
	wp.lock.Unlock()

	// Notify obsolete workers to stop.
	// This notification must be outside the wp.lock, since ch.ch
	// may be blocking and may consume a lot of time if many workers
	// are located on non-local CPUs.
	tmp := *scratch
	for i := range tmp {
		tmp[i].ch <- nil
		tmp[i] = nil
	}
}
```

### 参考

【Golang 开发需要协程池吗？】https://www.zhihu.com/question/302981392  
【来，控制一下 Goroutine 的并发数量】https://segmentfault.com/a/1190000017956396  
【golang协程池设计】https://segmentfault.com/a/1190000018193161  
【fasthttp中的协程池实现】https://segmentfault.com/a/1190000009133154    