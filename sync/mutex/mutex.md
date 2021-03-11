## 互斥锁

### 前言

本次的代码是基于`go version go1.13.15 darwin/amd64`   

### 什么是sync.Mutex

`sync.Mutex`是Go标准库中常用的一个排外锁。当一个`goroutine`获得了这个锁的拥有权后， 其它请求锁的`goroutine`就会阻塞在`Lock`方法的调用上，直到锁被释放。  

```go
var (
	mu      sync.Mutex
	balance int
)

func main() {
	Deposit(1)
	fmt.Println(Balance())
}

func Deposit(amount int) {
	mu.Lock()
	balance = balance + amount
	mu.Unlock()
}

func Balance() int {
	mu.Lock()
	b := balance
	mu.Unlock()
	return b
}
```

使用起来很简单，对需要锁定的资源，前面加`Lock()`锁定,完成的时候加`Unlock()`解锁就好了。  

### 分下源码

```go
const (
   // mutex is locked
// 是否加锁的标识
	mutexLocked = 1 << iota 
	mutexWoken
	mutexStarving
	mutexWaiterShift = iota

// 公平锁
//
// 锁有两种模式：正常模式和饥饿模式。
// 在正常模式下，所有的等待锁的goroutine都会存在一个先进先出的队列中（轮流被唤醒）
// 但是一个被唤醒的goroutine并不是直接获得锁，而是仍然需要和那些新请求锁的（new arrivial）
// 的goroutine竞争，而这其实是不公平的，因为新请求锁的goroutine有一个优势——它们正在CPU上
// 运行，并且数量可能会很多。所以一个被唤醒的goroutine拿到锁的概率是很小的。在这种情况下，
// 这个被唤醒的goroutine会加入到队列的头部。如果一个等待的goroutine有超过1ms（写死在代码中）
// 都没获取到锁，那么就会把锁转变为饥饿模式。
//
// 在饥饿模式中，锁的所有权会直接从释放锁(unlock)的goroutine转交给队列头的goroutine，
// 新请求锁的goroutine就算锁是空闲状态也不会去获取锁，并且也不会尝试自旋。它们只是排到队列的尾部。
//
// 如果一个goroutine获取到了锁之后，它会判断以下两种情况：
// 1. 它是队列中最后一个goroutine；
// 2. 它拿到锁所花的时间小于1ms；
// 以上只要有一个成立，它就会把锁转变回正常模式。

// 正常模式会有比较好的性能，因为即使有很多阻塞的等待锁的goroutine，
// 一个goroutine也可以尝试请求多次锁。
// 饥饿模式对于防止尾部延迟来说非常的重要。
	starvationThresholdNs = 1e6
)

// A Mutex is a mutual exclusion lock.
// The zero value for a Mutex is an unlocked mutex.
//
// A Mutex must not be copied after first use.
type Mutex struct {
// mutex锁当前的状态
	state int32
// 信号量，用于唤醒goroutine
	sema  uint32
}
```

重点开看下`state`的几种状态：  

大神写代码的思路就是惊奇，这里`state`又运用到了位移的操作  

- mutexLocked 对应右边低位第一个bit 1 代表锁被占用  0代表锁空闲  

- mutexWoken 对应右边低位第二个bit 1 表示已唤醒  0表示未唤醒  

- mutexStarving 对应右边低位第三个bit 1 代表锁处于饥饿模式  0代表锁处于正常模式

- mutexWaiterShift 值为3，根据 `mutex.state >> mutexWaiterShift` 得到当前阻塞的`goroutine`数目，最多可以阻塞`2^29`个`goroutine`。

- starvationThresholdNs 值为1e6纳秒，也就是1毫秒，当等待队列中队首g`oroutine`等待时间超过`starvationThresholdNs`也就是1毫秒，mutex进入饥饿模式。  

<img src="/img/sync_mutex_state.png" width = "568" height = "173" alt="sync_mutex" align=center />

### Lock

```go
// Lock locks m.
// 如果锁正在使用中，新的goroutine请求，将被阻塞，直到锁被释放
func (m *Mutex) Lock() {
	// 原子的(cas)来判断是否加锁
// 未加锁，直接加锁，返回
	if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
		if race.Enabled {
			race.Acquire(unsafe.Pointer(m))
		}
		return
	}
    // 这把锁，已经被别的goroutine持有
	m.lockSlow()
}

func (m *Mutex) lockSlow() {
	var waitStartTime int64
	starving := false
	awoke := false
	iter := 0
	old := m.state
	for {
// 在饥饿模式下，锁不需要自旋了
// 锁的所有权会直接转交给队列头的goroutine
// TODO
		if old&(mutexLocked|mutexStarving) == mutexLocked && runtime_canSpin(iter) {
        // 主动旋转是有意义的。
        // 尝试设置MutexWoken标志来通知解锁
        // 不唤醒其他被阻止的goroutine。
			if !awoke && old&mutexWoken == 0 && old>>mutexWaiterShift != 0 &&
				atomic.CompareAndSwapInt32(&m.state, old, old|mutexWoken) {
				awoke = true
			}
// 主动自旋
			runtime_doSpin()
			iter++
			old = m.state
			continue
		}
		new := old
		// 不要试图获取饥饿的互斥体，新的goroutine必须排队

		if old&mutexStarving == 0 {
			new |= mutexLocked
		}
		if old&(mutexLocked|mutexStarving) != 0 {
			new += 1 << mutexWaiterShift
		}
		// The current goroutine switches mutex to starvation mode.
		// But if the mutex is currently unlocked, don't do the switch.
		// Unlock expects that starving mutex has waiters, which will not
		// be true in this case.
		if starving && old&mutexLocked != 0 {
			new |= mutexStarving
		}
		if awoke {
			// The goroutine has been woken from sleep,
			// so we need to reset the flag in either case.
			if new&mutexWoken == 0 {
				throw("sync: inconsistent mutex state")
			}
			new &^= mutexWoken
		}
		if atomic.CompareAndSwapInt32(&m.state, old, new) {
			if old&(mutexLocked|mutexStarving) == 0 {
				break // locked the mutex with CAS
			}
			// If we were already waiting before, queue at the front of the queue.
			queueLifo := waitStartTime != 0
			if waitStartTime == 0 {
				waitStartTime = runtime_nanotime()
			}
			runtime_SemacquireMutex(&m.sema, queueLifo, 1)
			starving = starving || runtime_nanotime()-waitStartTime > starvationThresholdNs
			old = m.state
			if old&mutexStarving != 0 {
				// If this goroutine was woken and mutex is in starvation mode,
				// ownership was handed off to us but mutex is in somewhat
				// inconsistent state: mutexLocked is not set and we are still
				// accounted as waiter. Fix that.
				if old&(mutexLocked|mutexWoken) != 0 || old>>mutexWaiterShift == 0 {
					throw("sync: inconsistent mutex state")
				}
				delta := int32(mutexLocked - 1<<mutexWaiterShift)
				if !starving || old>>mutexWaiterShift == 1 {
					// Exit starvation mode.
					// Critical to do it here and consider wait time.
					// Starvation mode is so inefficient, that two goroutines
					// can go lock-step infinitely once they switch mutex
					// to starvation mode.
					delta -= mutexStarving
				}
				atomic.AddInt32(&m.state, delta)
				break
			}
			awoke = true
			iter = 0
		} else {
			old = m.state
		}
	}

	if race.Enabled {
		race.Acquire(unsafe.Pointer(m))
	}
}
```



### 参考

【sync.Mutex 源码分析】https://reading.hidevops.io/articles/sync/sync_mutex_source_code_analysis/  
【http://cbsheng.github.io/posts/%E4%B8%80%E4%BB%BD%E8%AF%A6%E7%BB%86%E6%B3%A8%E9%87%8A%E7%9A%84go-mutex%E6%BA%90%E7%A0%81/】http://cbsheng.github.io/posts/%E4%B8%80%E4%BB%BD%E8%AF%A6%E7%BB%86%E6%B3%A8%E9%87%8A%E7%9A%84go-mutex%E6%BA%90%E7%A0%81/  