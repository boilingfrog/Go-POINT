<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [semaphore](#semaphore)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [semaphore的作用](#semaphore%E7%9A%84%E4%BD%9C%E7%94%A8)
  - [如何使用](#%E5%A6%82%E4%BD%95%E4%BD%BF%E7%94%A8)
  - [分析下原理](#%E5%88%86%E6%9E%90%E4%B8%8B%E5%8E%9F%E7%90%86)
    - [Acquire](#acquire)
    - [TryAcquire](#tryacquire)
    - [Release](#release)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## semaphore

### 前言

本文是在`go version go1.13.15 darwin/amd64`上进行的  

### semaphore的作用

信号量是在并发编程中比较常见的一种同步机制，它会保证持有的计数器在0到初始化的权重之间，每次获取资源时都会将信号量中的计数器减去对应的数值，在释放时重新加回来，当遇到计数器大于信号量大小时就会进入休眠等待其他进程释放信号，我们常常会在控制访问资源的进程数量时用到。  

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

可以使用`semaphore`对控制一下`goroutine`的并发数量  

```go

```

### 分析下原理

```go
type waiter struct {
	n     int64
	ready chan<- struct{} // Closed when semaphore acquired.
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

	if n > s.size {
		// Don't make other Acquire calls block on one that's doomed to fail.
		s.mu.Unlock()
		<-ctx.Done()
		return ctx.Err()
	}

	ready := make(chan struct{})
	w := waiter{n: n, ready: ready}
	elem := s.waiters.PushBack(w)
	s.mu.Unlock()

	select {
	case <-ctx.Done():
		err := ctx.Err()
		s.mu.Lock()
		select {
		case <-ready:
			// Acquired the semaphore after we were canceled.  Rather than trying to
			// fix up the queue, just pretend we didn't notice the cancelation.
			err = nil
		default:
			s.waiters.Remove(elem)
		}
		s.mu.Unlock()
		return err

	case <-ready:
		return nil
	}
}
```

#### TryAcquire

```go
// TryAcquire acquires the semaphore with a weight of n without blocking.
// On success, returns true. On failure, returns false and leaves the semaphore unchanged.
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

#### Release

```go
// Release releases the semaphore with a weight of n.
func (s *Weighted) Release(n int64) {
	s.mu.Lock()
	s.cur -= n
	if s.cur < 0 {
		s.mu.Unlock()
		panic("semaphore: bad release")
	}
	for {
		next := s.waiters.Front()
		if next == nil {
			break // No more waiters blocked.
		}

		w := next.Value.(waiter)
		if s.size-s.cur < w.n {
			// Not enough tokens for the next waiter.  We could keep going (to try to
			// find a waiter with a smaller request), but under load that could cause
			// starvation for large requests; instead, we leave all remaining waiters
			// blocked.
			//
			// Consider a semaphore used as a read-write lock, with N tokens, N
			// readers, and one writer.  Each reader can Acquire(1) to obtain a read
			// lock.  The writer can Acquire(N) to obtain a write lock, excluding all
			// of the readers.  If we allow the readers to jump ahead in the queue,
			// the writer will starve — there is always one token available for every
			// reader.
			break
		}

		s.cur += w.n
		s.waiters.Remove(next)
		close(w.ready)
	}
	s.mu.Unlock()
}
```

### 参考
【Golang并发同步原语之-信号量Semaphor】https://blog.haohtml.com/archives/25563    



