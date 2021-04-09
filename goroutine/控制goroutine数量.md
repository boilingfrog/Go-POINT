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
      - [getCh](#getch)
      - [workerFunc](#workerfunc)
    - [panjf2000/ants](#panjf2000ants)
    - [看下实现](#%E7%9C%8B%E4%B8%8B%E5%AE%9E%E7%8E%B0)
      - [NewPoolWithFunc](#newpoolwithfunc)
      - [Release](#release)
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

1、按需增长`goroutine`数量，有一个最大值，同时监听`channel`，`Server`会把`accept`到的`connection`放入到`channel`中，这样监听的`goroutine`就能处理消费。  

2、本地维护了一个待使用的`channel`列表，当本地`channel`列表拿不到`ch`，会在`sync.pool`中取。  

3、对于待使用的`channel`列表，会定期清理掉超过最大空闲时间的`workerChan`。  

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

	// 最大的Workers数量
	MaxWorkersCount int

	LogAllErrors bool

	MaxIdleWorkerDuration time.Duration

	Logger Logger

	lock         sync.Mutex
	// 当前worker的数量
	workersCount int
	// worker停止的标识
	mustStop     bool

	// 等待使用的workerChan
	// 可能会被清理
	ready []*workerChan

	// 用来标识start和stop
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

	// 通知淘汰的workers停止
	// 此通知必须位于wp.lock之外，因为ch.ch
	// 如果有很多workers，可能会阻塞并且可能会花费大量时间
	// 位于非本地CPU上。
	tmp := *scratch
	for i := range tmp {
		tmp[i].ch <- nil
		tmp[i] = nil
	}
}
```

主要是清理掉最近最少使用的`workers`如果他们过了`maxIdleWorkerDuration`时间没有提供服务  

##### getCh 

获取一个`workerChan`

```go
func (wp *workerPool) getCh() *workerChan {
	var ch *workerChan
	createWorker := false

	wp.lock.Lock()
	ready := wp.ready
	n := len(ready) - 1
	// 如果ready为空
	if n < 0 {
		if wp.workersCount < wp.MaxWorkersCount {
			createWorker = true
			wp.workersCount++
		}
	} else {
		// 不为空从ready中取一个
		ch = ready[n]
		ready[n] = nil
		wp.ready = ready[:n]
	}
	wp.lock.Unlock()

	// 如果没拿到ch
	if ch == nil {
		if !createWorker {
			return nil
		}
		// 从缓存中获取一个ch
		vch := wp.workerChanPool.Get()
		ch = vch.(*workerChan)
		go func() {
			// 具体的执行函数
			wp.workerFunc(ch)
			// 再放入到pool中
		}()
	}
	return ch
}
```

梳理下流程：  

1、获取一个可执行的`workerChan`，如果`ready`中为空，并且`workersCount`没有达到最大值，增加`workersCount`数量，并且设置当前操作`createWorker = true`；  

2、`ready`中不为空，直接在`ready`获取一个；  

3、如果没有获取到则在`sync.pool`中获取一个，之后再放回到`pool`中。  

##### workerFunc

```go
func (wp *workerPool) workerFunc(ch *workerChan) {
	var c net.Conn

	var err error
	// 监听workerChan
	for c = range ch.ch {
		if c == nil {
			break
		}

		// 具体的业务逻辑
		...
		c = nil

		// 释放workerChan
		// 在mustStop的时候将会跳出循环
		if !wp.release(ch) {
			break
		}
	}

	wp.lock.Lock()
	wp.workersCount--
	wp.lock.Unlock()
}

// 把Conn放入到channel中
func (wp *workerPool) Serve(c net.Conn) bool {
	ch := wp.getCh()
	if ch == nil {
		return false
	}
	ch.ch <- c
	return true
}

func (wp *workerPool) release(ch *workerChan) bool {
	// 修改 ch.lastUseTime
	ch.lastUseTime = time.Now()
	wp.lock.Lock()
	// 如果需要停止，直接返回
	if wp.mustStop {
		wp.lock.Unlock()
		return false
	}
	// 将ch放到ready中
	wp.ready = append(wp.ready, ch)
	wp.lock.Unlock()
	return true
}
```

梳理下流程：  

1、`workerFunc`会监听`workerChan`，并且在使用完`workerChan`归还到`ready`中；  

2、`Serve`会把`connection`放入到`workerChan`中，这样`workerFunc`就能通过`workerChan`拿到需要处理的连接请求；  

3、当`workerFunc`拿到的`workerChan`为`nil`或`wp.mustStop`被设为了`true`，就跳出`for`循环。 

#### panjf2000/ants

先看下示例  

示例一   

```go
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/panjf2000/ants"
)

func demoFunc() {
	time.Sleep(10 * time.Millisecond)
	fmt.Println("Hello World!")
}

func main() {
	defer ants.Release()

	runTimes := 1000

	var wg sync.WaitGroup
	syncCalculateSum := func() {
		demoFunc()
		wg.Done()
	}
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		_ = ants.Submit(syncCalculateSum)
	}
	wg.Wait()
	fmt.Printf("running goroutines: %d\n", ants.Running())
	fmt.Printf("finish all tasks.\n")
}
```

示例二  

```go
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/panjf2000/ants"
)

var sum int32

func myFunc(i interface{}) {
	n := i.(int32)
	atomic.AddInt32(&sum, n)
	fmt.Printf("run with %d\n", n)
}

func main() {
	var wg sync.WaitGroup
	runTimes := 1000

	// Use the pool with a method,
	// set 10 to the capacity of goroutine pool and 1 second for expired duration.
	p, _ := ants.NewPoolWithFunc(10, func(i interface{}) {
		myFunc(i)
		wg.Done()
	})
	defer p.Release()
	// Submit tasks one by one.
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		_ = p.Invoke(int32(i))
	}
	wg.Wait()
	fmt.Printf("running goroutines: %d\n", p.Running())
	fmt.Printf("finish all tasks, result is %d\n", sum)
	if sum != 499500 {
		panic("the final result is wrong!!!")
	}
}
```

#### 看下实现  

示例一的实现：  

```go
// Pool accepts the tasks from client, it limits the total of goroutines to a given number by recycling goroutines.
type Pool struct {
	// capacity of the pool, a negative value means that the capacity of pool is limitless, an infinite pool is used to
	// avoid potential issue of endless blocking caused by nested usage of a pool: submitting a task to pool
	// which submits a new task to the same pool.
	capacity int32

	// running is the number of the currently running goroutines.
	running int32

	// workers is a slice that store the available workers.
	workers workerArray

	// state is used to notice the pool to closed itself.
	state int32

	// lock for synchronous operation.
	lock sync.Locker

	// cond for waiting to get a idle worker.
	cond *sync.Cond

	// workerCache speeds up the obtainment of the an usable worker in function:retrieveWorker.
	workerCache sync.Pool

	// blockingNum is the number of the goroutines already been blocked on pool.Submit, protected by pool.lock
	blockingNum int

	options *Options
}

// goWorker is the actual executor who runs the tasks,
// it starts a goroutine that accepts tasks and
// performs function calls.
type goWorker struct {
	// pool who owns this worker.
	pool *Pool

	// task is a job should be done.
	task chan func()

	// recycleTime will be update when putting a worker back into queue.
	recycleTime time.Time
}
```

##### Submit

```go
// NewPool generates an instance of ants pool.
func NewPool(size int, options ...Option) (*Pool, error) {
	opts := loadOptions(options...)

	if size <= 0 {
		size = -1
	}

	if expiry := opts.ExpiryDuration; expiry < 0 {
		return nil, ErrInvalidPoolExpiry
	} else if expiry == 0 {
		opts.ExpiryDuration = DefaultCleanIntervalTime
	}

	if opts.Logger == nil {
		opts.Logger = defaultLogger
	}

	p := &Pool{
		capacity: int32(size),
		lock:     internal.NewSpinLock(),
		options:  opts,
	}
	p.workerCache.New = func() interface{} {
		return &goWorker{
			pool: p,
			task: make(chan func(), workerChanCap),
		}
	}
	if p.options.PreAlloc {
		if size == -1 {
			return nil, ErrInvalidPreAllocSize
		}
		p.workers = newWorkerArray(loopQueueType, size)
	} else {
		p.workers = newWorkerArray(stackType, 0)
	}

	p.cond = sync.NewCond(p.lock)

	// Start a goroutine to clean up expired workers periodically.
	go p.purgePeriodically()

	return p, nil
}

// Submit submits a task to this pool.
func (p *Pool) Submit(task func()) error {
	if p.IsClosed() {
		return ErrPoolClosed
	}
	var w *goWorker
	if w = p.retrieveWorker(); w == nil {
		return ErrPoolOverload
	}
	w.task <- task
	return nil
}

// retrieveWorker returns a available worker to run the tasks.
func (p *Pool) retrieveWorker() (w *goWorker) {
	spawnWorker := func() {
        // sync.pool中获取
		w = p.workerCache.Get().(*goWorker)
		w.run()
	}

	p.lock.Lock()

	w = p.workers.detach()
	if w != nil {
		p.lock.Unlock()
     // capacity为负数，不限制池的容量
	} else if capacity := p.Cap(); capacity == -1 {
		p.lock.Unlock()
		spawnWorker()
        // 正在运行的goroutine数量小于pool的容量
	} else if p.Running() < capacity {
		p.lock.Unlock()
		spawnWorker()
	} else {
        // 当Nonblocking为true，Pool.Submit将不会被阻塞
		if p.options.Nonblocking {
			p.lock.Unlock()
			return
		}
	Reentry:
		// 当前阻塞的goroutine大于最大的限制
		if p.options.MaxBlockingTasks != 0 && p.blockingNum >= p.options.MaxBlockingTasks {
			p.lock.Unlock()
			return
		}
		p.blockingNum++
		p.cond.Wait()
		p.blockingNum--
		var nw int
        // 如果正在运行的goroutine数量为0
		if nw = p.Running(); nw == 0 {
			p.lock.Unlock()
			if !p.IsClosed() {
                // sync.pool中取一个
				spawnWorker()
			}
			return
		}
        // 没找到
		if w = p.workers.detach(); w == nil {
            // 正在运行的goroutine数量小于pool的容量
			if nw < capacity {
				p.lock.Unlock()
				spawnWorker()
				return
			}
            // 循环
			goto Reentry
		}

		p.lock.Unlock()
	}
	return
}

// 关闭pool
func (p *Pool) Release() {
    // 写入关闭的标识
	atomic.StoreInt32(&p.state, CLOSED)
	p.lock.Lock()
	p.workers.reset()
	p.lock.Unlock()
    // 唤醒等待的
	p.cond.Broadcast()
}

// purgePeriodically clears expired workers periodically which runs in an individual goroutine, as a scavenger.
func (p *Pool) purgePeriodically() {
	heartbeat := time.NewTicker(p.options.ExpiryDuration)
	defer heartbeat.Stop()

	for range heartbeat.C {
		if p.IsClosed() {
			break
		}

		p.lock.Lock()
		expiredWorkers := p.workers.retrieveExpiry(p.options.ExpiryDuration)
		p.lock.Unlock()

		// Notify obsolete workers to stop.
		// This notification must be outside the p.lock, since w.task
		// may be blocking and may consume a lot of time if many workers
		// are located on non-local CPUs.
		for i := range expiredWorkers {
			expiredWorkers[i].task <- nil
			expiredWorkers[i] = nil
		}

		// There might be a situation that all workers have been cleaned up(no any worker is running)
		// while some invokers still get stuck in "p.cond.Wait()",
		// then it ought to wakes all those invokers.
		if p.Running() == 0 {
			p.cond.Broadcast()
		}
	}
}
```

梳理下思路：  

1、先初始化缓存池的大小，然后处理任务事件的时候，一个task分配一个goWorker；  

2、在拿goWorker的过程中会存在下面集中情况；  

- 本地的缓存中有空闲的goWorker，直接取出；  

- 本地缓存没有，通过`sync.Pool`注册的new事件，生成一个`goWorker`，不过这块有点奇怪对于`sync.Pool`的使用，生成的`goWorker`没有放到`pool`中，只是用到了`sync.Pool`的找不到然后`New`的特点。

对于`sync.pool`的`Get`

1、优先从`private`中选择对象  
2、若取不到，则尝试从`localPool`的`shared`队列的队头进行读取  
3、若取不到，则尝试从其他的P中进行偷取`getSlow`  
4、若还是取不到，则使用`New`方法新建  

<img src="/img/pool_1.png" width = "537" height = "990" alt="gc" align="center" />

每次都经历了几次无效的查询，最后才使用`New`方法新建，感觉这个`sync.pool`的使用有点多余了。   

3、如果缓存池满了，就循环去拿直到成功拿出一个；  

4、同时也会定期清理掉过期的`goWorker`,通过`sync.Cond`唤醒其的阻塞等待；  

示例二的实现：

```go
// PoolWithFunc accepts the tasks from client,
// it limits the total of goroutines to a given number by recycling goroutines.
type PoolWithFunc struct {
	// capacity of the pool.
	// pool的大小
	capacity int32

	// running is the number of the currently running goroutines.
	// 当前正在运行的goroutine的数量
	running int32

	// workers is a slice that store the available workers.
	// 存储可用的goWorkerWithFunc
	workers []*goWorkerWithFunc

	// state is used to notice the pool to closed itself.
	// 用于通知pool关闭的标识
	state int32

	// lock for synchronous operation.
	lock sync.Locker

	// cond for waiting to get a idle worker.
	// 通过Cond找到一个空闲的worker
	cond *sync.Cond

	// poolFunc is the function for processing tasks.
	// 处理具体的任务
	poolFunc func(interface{})

	// workerCache speeds up the obtainment of the an usable worker in function:retrieveWorker.
	workerCache sync.Pool

	// blockingNum is the number of the goroutines already been blocked on pool.Submit, protected by pool.lock
	// 被blocked的goroutine的数量
	blockingNum int

	options *Options
}

// goWorkerWithFunc是执行任务的实际执行者，
// 它启动一个接受任务的goroutine并执行函数调用。
type goWorkerWithFunc struct {
	// pool who owns this worker.
	pool *PoolWithFunc

	// args is a job should be done.
	args chan interface{}

	// recycleTime will be update when putting a worker back into queue.
	recycleTime time.Time
}

// NewPoolWithFunc generates an instance of ants pool with a specific function.
func NewPoolWithFunc(size int, pf func(interface{}), options ...Option) (*PoolWithFunc, error) {
	if size <= 0 {
		return nil, ErrInvalidPoolSize
	}

	if pf == nil {
		return nil, ErrLackPoolFunc
	}

	opts := loadOptions(options...)

	if expiry := opts.ExpiryDuration; expiry < 0 {
		return nil, ErrInvalidPoolExpiry
	} else if expiry == 0 {
		opts.ExpiryDuration = DefaultCleanIntervalTime
	}

	if opts.Logger == nil {
		opts.Logger = defaultLogger
	}

	p := &PoolWithFunc{
		capacity: int32(size),
		poolFunc: pf,
		lock:     internal.NewSpinLock(),
		options:  opts,
	}

	// 这里的sync.pool里面存储的是goWorkerWithFunc，在sync.pool拿不到到的时候创建
	p.workerCache.New = func() interface{} {
		return &goWorkerWithFunc{
			pool: p,
			args: make(chan interface{}, workerChanCap),
		}
	}
	if p.options.PreAlloc {
		p.workers = make([]*goWorkerWithFunc, 0, size)
	}
	p.cond = sync.NewCond(p.lock)

	// 定期清理过期的workers
	go p.purgePeriodically()

	return p, nil
}

// purgePeriodically clears expired workers periodically which runs in an individual goroutine, as a scavenger.
func (p *PoolWithFunc) purgePeriodically() {
	heartbeat := time.NewTicker(p.options.ExpiryDuration)
	defer heartbeat.Stop()

	var expiredWorkers []*goWorkerWithFunc
	for range heartbeat.C {
		if p.IsClosed() {
			break
		}
		currentTime := time.Now()
		p.lock.Lock()
		idleWorkers := p.workers
		n := len(idleWorkers)
		var i int
		// 找到需要清理的workers的范围
		for i = 0; i < n && currentTime.Sub(idleWorkers[i].recycleTime) > p.options.ExpiryDuration; i++ {
		}
		expiredWorkers = append(expiredWorkers[:0], idleWorkers[:i]...)
		if i > 0 {
			m := copy(idleWorkers, idleWorkers[i:])
			for i = m; i < n; i++ {
				idleWorkers[i] = nil
			}
			p.workers = idleWorkers[:m]
		}
		p.lock.Unlock()

		// 通知淘汰的workers停止
		// 此通知必须位于wp.lock之外，因为ch.ch
		// 如果有很多workers，可能会阻塞并且可能会花费大量时间
		// 位于非本地CPU上。
		for i, w := range expiredWorkers {
			w.args <- nil
			expiredWorkers[i] = nil
		}

		// 可能存在清理所有工作程序的情况（没有任何工作程序在运行）
		// 虽然某些调用程序仍然卡在“ p.cond.Wait（）”中，
		// 然后应该唤醒所有这些调用者。
		if p.Running() == 0 {
			p.cond.Broadcast()
		}
	}
}
```

梳理下流程：  

1、在生成特定的执行的实例池的时候，也会定期清理掉那些无用的`workers`；  

2、这里使用`sync.pool`存储`goWorkerWithFunc`，在`sync.pool`拿不到到的时候创建。  

##### Invoke

```go
// 提交任务到任务池
func (p *PoolWithFunc) Invoke(args interface{}) error {
	if p.IsClosed() {
		return ErrPoolClosed
	}
	var w *goWorkerWithFunc
	if w = p.retrieveWorker(); w == nil {
		return ErrPoolOverload
	}
	w.args <- args
	return nil
}

// retrieveWorker returns a available worker to run the tasks.
func (p *PoolWithFunc) retrieveWorker() (w *goWorkerWithFunc) {
	spawnWorker := func() {
		// 从sync.pool中获取一个
		w = p.workerCache.Get().(*goWorkerWithFunc)
		w.run()
	}

	p.lock.Lock()
	idleWorkers := p.workers
	n := len(idleWorkers) - 1
	// 如果workers有，取一个出来
	if n >= 0 {
		w = idleWorkers[n]
		idleWorkers[n] = nil
		p.workers = idleWorkers[:n]
		p.lock.Unlock()
		// 当前正在运行的goroutine的数量小于pool的大小
		// 从sync.pool中取一个
	} else if p.Running() < p.Cap() {
		p.lock.Unlock()
		spawnWorker()
	} else {
		if p.options.Nonblocking {
			p.lock.Unlock()
			return
		}
	Reentry:
		// 当前阻塞的goroutine大于最大的限制
		if p.options.MaxBlockingTasks != 0 && p.blockingNum >= p.options.MaxBlockingTasks {
			p.lock.Unlock()
			return
		}
		p.blockingNum++
		// 阻塞当前
		p.cond.Wait()
		p.blockingNum--
		// 当前正在运行的goroutine数量为0
		// sync.pool中取一个
		if p.Running() == 0 {
			p.lock.Unlock()
			if !p.IsClosed() {
				spawnWorker()
			}
			return
		}
		l := len(p.workers) - 1
		// 可用的workers数量为0
		if l < 0 {
			// 正在运行的goroutine数量小于pool的大小
			// sync.pool中取一个
			if p.Running() < p.Cap() {
				p.lock.Unlock()
				spawnWorker()
				return
			}
			// 循环
			goto Reentry
		}
		w = p.workers[l]
		p.workers[l] = nil
		p.workers = p.workers[:l]
		p.lock.Unlock()
	}
	return
}

// 关闭pool
func (p *PoolWithFunc) Release() {
	// 写入state关闭的标识
	atomic.StoreInt32(&p.state, CLOSED)
	p.lock.Lock()
	idleWorkers := p.workers
	for _, w := range idleWorkers {
		w.args <- nil
	}
	// 清理掉workers
	p.workers = nil
	p.lock.Unlock()
	// 可能有一些调用者在retrieveWorker（）中等待，因此我们需要将其唤醒以防止那些调用者无限阻塞。
	p.cond.Broadcast()
}
```

梳理下流程：  

1、`Invoke`提交任务到任务池，调用`retrieveWorker`取出一个`ch`来，放入当前的`task`；  

2、



 
### 参考

【Golang 开发需要协程池吗？】https://www.zhihu.com/question/302981392  
【来，控制一下 Goroutine 的并发数量】https://segmentfault.com/a/1190000017956396  
【golang协程池设计】https://segmentfault.com/a/1190000018193161  
【fasthttp中的协程池实现】https://segmentfault.com/a/1190000009133154    
【panjf2000/ants】https://github.com/panjf2000/ants   