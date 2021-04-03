<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [控制goroutine数量](#%E6%8E%A7%E5%88%B6goroutine%E6%95%B0%E9%87%8F)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [控制goroutine的数量](#%E6%8E%A7%E5%88%B6goroutine%E7%9A%84%E6%95%B0%E9%87%8F)
    - [通过channel+sync](#%E9%80%9A%E8%BF%87channelsync)
    - [使用semaphore](#%E4%BD%BF%E7%94%A8semaphore)
    - [线程池](#%E7%BA%BF%E7%A8%8B%E6%B1%A0)
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


### 参考

【Golang 开发需要协程池吗？】https://www.zhihu.com/question/302981392  
【来，控制一下 Goroutine 的并发数量】https://segmentfault.com/a/1190000017956396  