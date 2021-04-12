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
      - [设计思路](#%E8%AE%BE%E8%AE%A1%E6%80%9D%E8%B7%AF)
    - [go-playground/pool](#go-playgroundpool)
      - [workUnit](#workunit)
      - [limitedPool](#limitedpool)
      - [batch](#batch)
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

##### 设计思路

整体的设计思路  

梳理下思路：  

1、先初始化缓存池的大小，然后处理任务事件的时候，一个`task`分配一个`goWorker`；  

2、在拿goWorker的过程中会存在下面集中情况；  

- 本地的缓存中有空闲的`goWorker`，直接取出；  

- 本地缓存没有就去`sync.Pool`，拿一个`goWorker`。   

3、如果缓存池满了，就循环去拿直到成功拿出一个；  

4、同时也会定期清理掉过期的`goWorker`,通过`sync.Cond`唤醒其的阻塞等待；  

5、对于使用完成的`goWorker`在使用完成之后重新归还到`free pool`。     

具体的设计细节可参考，作者的文章[Goroutine 并发调度模型深度解析之手撸一个高性能 goroutine 池](HTTPS://strikefreedom.top/high-performance-implementation-of-goroutine-pool)

#### go-playground/pool

先放几个使用的demo  

**Per Unit Work**

```go
package main

import (
	"fmt"
	"time"

	"gopkg.in/go-playground/pool.v3"
)

func main() {

	p := pool.NewLimited(10)
	defer p.Close()

	user := p.Queue(getUser(13))
	other := p.Queue(getOtherInfo(13))

	user.Wait()
	if err := user.Error(); err != nil {
		// handle error
	}

	// do stuff with user
	username := user.Value().(string)
	fmt.Println(username)

	other.Wait()
	if err := other.Error(); err != nil {
		// handle error
	}

	// do stuff with other
	otherInfo := other.Value().(string)
	fmt.Println(otherInfo)
}

func getUser(id int) pool.WorkFunc {

	return func(wu pool.WorkUnit) (interface{}, error) {

		// simulate waiting for something, like TCP connection to be established
		// or connection from pool grabbed
		time.Sleep(time.Second * 1)

		if wu.IsCancelled() {
			// return values not used
			return nil, nil
		}

		// ready for processing...

		return "Joeybloggs", nil
	}
}

func getOtherInfo(id int) pool.WorkFunc {

	return func(wu pool.WorkUnit) (interface{}, error) {

		// simulate waiting for something, like TCP connection to be established
		// or connection from pool grabbed
		time.Sleep(time.Second * 1)

		if wu.IsCancelled() {
			// return values not used
			return nil, nil
		}

		// ready for processing...

		return "Other Info", nil
	}
}
```

**Batch Work**

```go
package main

import (
	"fmt"
	"time"

	"gopkg.in/go-playground/pool.v3"
)

func main() {

	p := pool.NewLimited(10)
	defer p.Close()

	batch := p.Batch()

	// for max speed Queue in another goroutine
	// but it is not required, just can't start reading results
	// until all items are Queued.

	go func() {
		for i := 0; i < 10; i++ {
			batch.Queue(sendEmail("email content"))
		}

		// DO NOT FORGET THIS OR GOROUTINES WILL DEADLOCK
		// if calling Cancel() it calles QueueComplete() internally
		batch.QueueComplete()
	}()

	for email := range batch.Results() {

		if err := email.Error(); err != nil {
			// handle error
			// maybe call batch.Cancel()
		}

		// use return value
		fmt.Println(email.Value().(bool))
	}
}

func sendEmail(email string) pool.WorkFunc {
	return func(wu pool.WorkUnit) (interface{}, error) {

		// simulate waiting for something, like TCP connection to be established
		// or connection from pool grabbed
		time.Sleep(time.Second * 1)

		if wu.IsCancelled() {
			// return values not used
			return nil, nil
		}

		// ready for processing...

		return true, nil // everything ok, send nil, error if not
	}
}
```

来看下实现  

##### workUnit

`workUnit`作为`channel`信息进行传递，用来给`work`传递当前需要执行的任务信息。   

```go
// WorkUnit contains a single uint of works values
type WorkUnit interface {

	// 阻塞直到当前任务被完成或被取消
	Wait()

	// 执行函数返回的结果
	Value() interface{}

	// Error returns the Work Unit's error
	Error() error

	// 取消当前的可执行任务
	Cancel()

	// 判断当前的可执行单元是否被取消了
	IsCancelled() bool
}

var _ WorkUnit = new(workUnit)

// workUnit contains a single unit of works values
type workUnit struct {
	// 任务执行的结果
	value      interface{}
	// 错误信息
	err        error
	// 通知任务完成
	done       chan struct{}
	// 需要执行的任务函数
	fn         WorkFunc
	// 任务是会否被取消
	cancelled  atomic.Value
	// 是否正在取消任务
	cancelling atomic.Value
	// 任务是否正在执行
	writing    atomic.Value
}
```

##### limitedPool

```go
var _ Pool = new(limitedPool)

// limitedPool contains all information for a limited pool instance.
type limitedPool struct {
	// 并发量
	workers uint
	// work的channel
	work    chan *workUnit
	// 通知结束的channel
	cancel  chan struct{}
	// 是否关闭的标识
	closed  bool
	// 读写锁
	m       sync.RWMutex
}

// 初始化一个pool
func NewLimited(workers uint) Pool {

	if workers == 0 {
		panic("invalid workers '0'")
	}
	// 初始化pool的work数量
	p := &limitedPool{
		workers: workers,
	}
	// 初始化pool的操作
	p.initialize()

	return p
}

func (p *limitedPool) initialize() {
	// channel的长度为work数量的两倍
	p.work = make(chan *workUnit, p.workers*2)
	p.cancel = make(chan struct{})
	p.closed = false

	// fire up workers here
	for i := 0; i < int(p.workers); i++ {
		p.newWorker(p.work, p.cancel)
	}
}

// 将工作传递并取消频道到newWorker（）以避免任何潜在的竞争状况
// 在p.work读写之间
func (p *limitedPool) newWorker(work chan *workUnit, cancel chan struct{}) {
	go func(p *limitedPool) {

		var wu *workUnit

		defer func(p *limitedPool) {
			// 捕获异常，结束掉异常的工作单元，并将其再次作为新的任务启动
			if err := recover(); err != nil {

				trace := make([]byte, 1<<16)
				n := runtime.Stack(trace, true)

				s := fmt.Sprintf(errRecovery, err, string(trace[:int(math.Min(float64(n), float64(7000)))]))

				iwu := wu
				iwu.err = &ErrRecovery{s: s}
				close(iwu.done)

				// 重新启动
				p.newWorker(p.work, p.cancel)
			}
		}(p)

		var value interface{}
		var err error
		// 监听channel，读取内容
		for {
			select {
			// channel中取出数据
			case wu = <-work:

				// 防止channel 被关闭后读取到零值
				if wu == nil {
					continue
				}

				// 单个和批量的cancellation这个都支持
				if wu.cancelled.Load() == nil {
					// 执行我们的业务函数
					value, err = wu.fn(wu)

					wu.writing.Store(struct{}{})

					// 如果WorkFunc取消了此工作单元，则需要再次检查
					// 防止产生竞争条件
					if wu.cancelled.Load() == nil && wu.cancelling.Load() == nil {
						wu.value, wu.err = value, err

						// 执行完成，关闭当前channel
						close(wu.done)
					}
				}
				// 如果取消了，就退出
			case <-cancel:
				return
			}
		}

	}(p)
}

// 放置一个执行的task到channel，并返回channel
func (p *limitedPool) Queue(fn WorkFunc) WorkUnit {
	// 初始化一个workUnit类型的channel
	w := &workUnit{
		done: make(chan struct{}),
		// 具体的执行函数
		fn:   fn,
	}

	go func() {
		p.m.RLock()
		// 如果pool关闭的时候通知channel关闭
		if p.closed {
			w.err = &ErrPoolClosed{s: errClosed}
			if w.cancelled.Load() == nil {
				close(w.done)
			}
			p.m.RUnlock()
			return
		}
		// 将channel传递给pool的work
		p.work <- w

		p.m.RUnlock()
	}()

	return w
}
```

梳理下流程：  

1、首先初始化`pool`的大小；  

2、然后根据`pool`的大小启动对应数量的`worker`，阻塞等待`channel`被塞入可执行函数；  

3、然后可执行函数会被放入`workUnit`，然后通过`channel`传递给阻塞的`worker`。  


同样这里也提供了批量执行的方法  

##### batch

```go
// batch contains all information for a batch run of WorkUnits
type batch struct {
	pool    Pool
	m       sync.Mutex
    // WorkUnit的切片
	units   []WorkUnit
    // 结果集,执行完后的workUnit会更新其value,error,可以从结果集channel中读取
	results chan WorkUnit
    // 通知batch是否完成
	done    chan struct{}
	closed  bool
	wg      *sync.WaitGroup
}

// 初始化Batch
func newBatch(p Pool) Batch {
	return &batch{
		pool:    p,
		units:   make([]WorkUnit, 0, 4),
		results: make(chan WorkUnit),
		done:    make(chan struct{}),
		wg:      new(sync.WaitGroup),
	}
}


// 将WorkFunc放入到WorkUnit中并保留取消和输出结果的参考。
func (b *batch) Queue(fn WorkFunc) {

	b.m.Lock()

	if b.closed {
		b.m.Unlock()
		return
	}
    // 返回一个WorkUnit
	wu := b.pool.Queue(fn)

    // 放到WorkUnit的切片中
	b.units = append(b.units, wu) // keeping a reference for cancellation purposes
	b.wg.Add(1)
	b.m.Unlock()

    // 执行任务
	go func(b *batch, wu WorkUnit) {
		wu.Wait()
		b.results <- wu
		b.wg.Done()
	}(b, wu)
}


// QueueComplete让批处理知道不再有排队的工作单元
// 以便在所有工作完成后可以关闭结果渠道。
// 警告：如果未调用此函数，则结果通道将永远不会耗尽，
// 但会永远阻止以获取更多结果。
func (b *batch) QueueComplete() {
	b.m.Lock()
	b.closed = true
	close(b.done)
	b.m.Unlock()
}

// 取消批次的任务
func (b *batch) Cancel() {

	b.QueueComplete() 

	b.m.Lock()

	// 一个个取消units，倒叙的取消
	for i := len(b.units) - 1; i >= 0; i-- {
		b.units[i].Cancel()
	}

	b.m.Unlock()
}

// 输出执行完成的结果集
func (b *batch) Results() <-chan WorkUnit {
	go func(b *batch) {
		<-b.done
		b.m.Lock()
		b.wg.Wait()
		b.m.Unlock()
		close(b.results)
	}(b)

	return b.results
}
```


### 参考

【Golang 开发需要协程池吗？】https://www.zhihu.com/question/302981392  
【来，控制一下 Goroutine 的并发数量】https://segmentfault.com/a/1190000017956396  
【golang协程池设计】https://segmentfault.com/a/1190000018193161  
【fasthttp中的协程池实现】https://segmentfault.com/a/1190000009133154    
【panjf2000/ants】https://github.com/panjf2000/ants   
【golang协程池设计】https://segmentfault.com/a/1190000018193161  