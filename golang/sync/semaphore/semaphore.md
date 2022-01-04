<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [semaphore](#semaphore)
  - [semaphore的作用](#semaphore%E7%9A%84%E4%BD%9C%E7%94%A8)
  - [如何使用](#%E5%A6%82%E4%BD%95%E4%BD%BF%E7%94%A8)
  - [分析下原理](#%E5%88%86%E6%9E%90%E4%B8%8B%E5%8E%9F%E7%90%86)
    - [Acquire](#acquire)
    - [TryAcquire](#tryacquire)
    - [Release](#release)
  - [总结](#%E6%80%BB%E7%BB%93)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## semaphore

### semaphore的作用

信号量是在并发编程中比较常见的一种同步机制，它会保证持有的计数器在0到初始化的权重之间，每次获取资源时都会将信号量中的计数器减去对应的数值，在释放时重新加回来，当遇到计数器大于信号量大小时就会进入休眠等待其他进程释放信号。  

go中的`semaphore`，提供`sleep`和`wakeup`原语，使其能够在其它同步原语中的竞争情况下使用。当一个`goroutine`需要休眠时，将其进行集中存放，当需要`wakeup`时，再将其取出，重新放入调度器中。   

go中本身提供了`semaphore`的相关方法，不过只能在内部调用  

```go
// go/src/sync/runtime.go
func runtime_Semacquire(s *uint32)

func runtime_SemacquireMutex(s *uint32, lifo bool, skipframes int)

func runtime_Semrelease(s *uint32, handoff bool, skipframes int)
```

扩展包`golang.org/x/sync/semaphore`提供了一种带权重的信号量实现方式，我们可以按照不同的权重对资源的访问进行管理。   

### 如何使用

通过信号量来限制并行的`goroutine`数量，达到最大的`maxWorkers`数量，`Acquire`将会阻塞，直到其中一个`goroutine`执行完成，释放出信号量。  
```go
// Example_workerPool演示如何使用信号量来限制
// 用于并行任务的goroutine。
func main() {
	ctx := context.Background()

	var (
		maxWorkers = runtime.GOMAXPROCS(0)
		sem        = semaphore.NewWeighted(int64(maxWorkers))
		out        = make([]int, 32)
	)

	// Compute the output using up to maxWorkers goroutines at a time.
	for i := range out {
		// When maxWorkers goroutines are in flight, Acquire blocks until one of the
		// workers finishes.
		if err := sem.Acquire(ctx, 1); err != nil {
			log.Printf("Failed to acquire semaphore: %v", err)
			break
		}

		go func(i int) {
			defer sem.Release(1)
			// doSomething
			out[i] = i + 1
		}(i)
	}

	// Acquire all of the tokens to wait for any remaining workers to finish.
	//
	// If you are already waiting for the workers by some other means (such as an
	// errgroup.Group), you can omit this final Acquire call.
	if err := sem.Acquire(ctx, int64(maxWorkers)); err != nil {
		log.Printf("Failed to acquire semaphore: %v", err)
	}

	fmt.Println(out)
}
```

### 分析下原理

```go
type waiter struct {
	// 信号量的权重
	n     int64
	// 获得信号量后关闭
	ready chan<- struct{}
}

// NewWeighted使用给定的值创建一个新的加权信号量
// 并发访问的最大组合权重。
func NewWeighted(n int64) *Weighted {
	w := &Weighted{size: n}
	return w
}

// 加权提供了一种方法来绑定对资源的并发访问。
// 呼叫者可以请求以给定的权重进行访问。
type Weighted struct {
	// 表示最大资源数量，取走时会减少，释放时会增加
	size    int64
	// 计数器，记录当前已使用资源数，值范围[0 - size]
	cur     int64
	mu      sync.Mutex
	// 等待队列，表示申请资源时由于可使用资源不够而陷入阻塞等待的调用者列表
	waiters list.List
```

#### Acquire

阻塞的获取指定权种的资源，如果没有空闲的资源，会进去休眠等待。  

```go
// Acquire获取权重为n的信号量，阻塞直到资源可用或ctx完成。
// 成功时，返回nil。失败时返回 ctx.Err（）并保持信号量不变。
// 如果ctx已经完成，则Acquire仍然可以成功执行而不会阻塞
func (s *Weighted) Acquire(ctx context.Context, n int64) error {
	s.mu.Lock()
	// 如果资源足够，并且没有排队等待的waiters
	// cur+n,直接返回
	if s.size-s.cur >= n && s.waiters.Len() == 0 {
		s.cur += n
		s.mu.Unlock()
		return nil
	}
	// 资源不够，err返回
	if n > s.size {
		// 不要其他的Acquire，阻塞在此
		s.mu.Unlock()
		<-ctx.Done()
		return ctx.Err()
	}

	ready := make(chan struct{})
	// 组装waiter
	w := waiter{n: n, ready: ready}
	// 插入waiters中
	elem := s.waiters.PushBack(w)
	s.mu.Unlock()

	// 阻塞等待，直到资源可用或ctx完成
	select {
	case <-ctx.Done():
		err := ctx.Err()
		s.mu.Lock()
		select {
		case <-ready:
			// 在canceled之后获取了信号量，不要试图去修复队列，假装没看到取消
			err = nil
		default:
			s.waiters.Remove(elem)
		}
		s.mu.Unlock()
		return err
		// 等待者被唤醒了
	case <-ready:
		return nil
	}
}
```

梳理下流程：  

1、如果资源够用并且没有等待队列，添加已经使用的资源数;  

2、如果超过资源数，抛出`err`;   

3、资源够用，并且等待队列，将之后的加入到等待队列中;  

4、阻塞直到资源可用或`ctx`完成。  

#### TryAcquire

非阻塞地获取指定权重的资源，如果当前没有空闲资源，会直接返回`false`。  

```go
// TryAcquire获取权重为n的信号量而不阻塞。
// 成功时返回true。 失败时，返回false并保持信号量不变。
func (s *Weighted) TryAcquire(n int64) bool {
	s.mu.Lock()
	success := s.size-s.cur >= n && s.waiters.Len() == 0
	if success {
		s.cur += n
	}
	s.mu.Unlock()
	return success
}
```

`TryAcquire`获取权重为n的信号量而不阻塞，相比`Acquire`少了等待队列的处理。  

#### Release

用于释放指定权重的资源，如果有`waiters`则尝试去一一唤醒`waiter`。   

```go
// Release释放权值为n的信号量。
func (s *Weighted) Release(n int64) {
	s.mu.Lock()
	s.cur -= n
	// cur的范围在[0 - size]
	if s.cur < 0 {
		s.mu.Unlock()
		panic("semaphore: bad release")
	}
	s.notifyWaiters()
	s.mu.Unlock()
}

func (s *Weighted) notifyWaiters() {
	// 如果有阻塞的waiters，尝试去进行一一唤醒 
	// 唤醒的时候，先进先出，避免被资源比较大的waiter被饿死
	for {
		next := s.waiters.Front()
		// 已经没有waiter了
		if next == nil {
			break
		}

		w := next.Value.(waiter)
		// waiter需要的资源不足
		if s.size-s.cur < w.n {
			// 没有足够的令牌供下一个waiter使用。我们可以继续（尝试
			// 查找请求较小的waiter），但在负载下可能会导致
			// 饥饿的大型请求；相反，我们留下所有剩余的waiter阻塞
			//
			// 考虑一个用作读写锁的信号量，带有N个令牌，N个reader和一位writer
			// 每个reader都可以通过Acquire（1）获取读锁。
			// writer写入可以通过Acquire（N）获得写锁定，但不包括所有的reader。
			// 如果我们允许读者在队列中前进，writer将会饿死-总是有一个令牌可供每个读者。
			break
		}

		s.cur += w.n
		s.waiters.Remove(next)
		close(w.ready)
	}
}
```

对于`waiters`的唤醒，遵循的原则总是先进先出。当有10个资源可以被使用，第一个`waiter`需要100个资源，第二个`waiter`需要1个资源。不会让第二个先释放，必须等待第一个`waiter`被释放。这样避免需要资源比较大`waiter`的被饿死，因为这样需要资源数比较小的`waiter`，总是可以被最先释放，需要资源比较大的`waiter`，就没有获取资源的机会了。  

### 总结

`Acquire`和`TryAcquire`都可用于获取资源，`Acquire`是可以阻塞的获取资源，`TryAcquire`只能非阻塞的获取资源;  

`Release`对于`waiters`的唤醒原则，总是先进先出，避免资源需求比较大的`waiter`被饿死;   

### 参考

【Golang并发同步原语之-信号量Semaphor】https://blog.haohtml.com/archives/25563    
【Go并发编程实战--信号量的使用方法和其实现原理】https://juejin.cn/post/6906677772479889422  



