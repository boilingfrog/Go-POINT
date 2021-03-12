<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [互斥锁](#%E4%BA%92%E6%96%A5%E9%94%81)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [什么是sync.Mutex](#%E4%BB%80%E4%B9%88%E6%98%AFsyncmutex)
  - [分下源码](#%E5%88%86%E4%B8%8B%E6%BA%90%E7%A0%81)
  - [Lock](#lock)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

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
	// 是否处于饥饿模式
	starving := false
	// 用来存当前goroutine是否已唤醒
	awoke := false
	// 用来存当前goroutine的循环次数
	iter := 0
	// 记录下当前的状态
	old := m.state
	for {
		// 在饥饿模式下，锁不需要自旋了
		// 锁的所有权会直接转交给队列头的goroutine
		// 如果 已经获取了锁，并且不是饥饿状态，并且可以自旋
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
			// 循环次数加一
			iter++
			old = m.state
			continue
		}
		new := old
		// 不要试图获取饥饿的互斥体，新的goroutine必须排队
		// 如果锁是饥饿状态，新到的goroutine乖乖排队去
		// 非饥饿状态，把期望状态设置为mutexLocked(获取锁)
		if old&mutexStarving == 0 {
			// 伪代码：newState = locked
			new |= mutexLocked
		}
		// 如果锁是被获取状态，或者饥饿状态
		// 就把期望状态中的等待队列的等待者数量+1(实际上是new + 8)
		if old&(mutexLocked|mutexStarving) != 0 {
			new += 1 << mutexWaiterShift
		}
		// 当前的goroutine将互斥锁切换到饥饿模式。
		// 但是，如果互斥锁当前已解锁，请不要进行切换。
		// Unlock期望一个饥饿的锁会有一些等待拿锁的goroutine，而不只是一个
		// 这种情况下不会成立
		if starving && old&mutexLocked != 0 {
			// 设置为饥饿状态
			new |= mutexStarving
		}
		if awoke {
			// goroutine已从睡眠中唤醒，
			// 因此，无论哪种情况，我们都需reset
			if new&mutexWoken == 0 {
				throw("sync: inconsistent mutex state")
			}
			// 设置new设置为非唤醒状态
			// &^的意思是and not
			new &^= mutexWoken
		}
		// 原子(cas)更新state的状态
		if atomic.CompareAndSwapInt32(&m.state, old, new) {
			// 如果说old状态不是饥饿状态也不是被获取状态
			// 那么代表当前goroutine已经通过CAS成功获取了锁
			if old&(mutexLocked|mutexStarving) == 0 {
				// 直接break
				break // locked the mutex with CAS
			}
			// 如果我们之前已经在等了，那就排在队伍前面。
			queueLifo := waitStartTime != 0
			// 如果说之前没有等待过，就初始化设置现在的等待时间
			if waitStartTime == 0 {
				waitStartTime = runtime_nanotime()
			}
			// queueLifo为true，也就是之前已经在等了
			// runtime_SemacquireMutex中的lifo为true，则将等待服务程序放在等待队列的开头。
			// 会被阻塞
			runtime_SemacquireMutex(&m.sema, queueLifo, 1)
			// 阻塞被唤醒

			// 如果当前goroutine已经是饥饿状态了
			// 或者当前goroutine已经等待了1ms（在上面定义常量）以上
			// 就把当前goroutine的状态设置为饥饿
			starving = starving || runtime_nanotime()-waitStartTime > starvationThresholdNs
			old = m.state
			// 如果是饥饿模式
			if old&mutexStarving != 0 {
				// 如果goroutine被唤醒，互斥锁处于饥饿模式
				// 锁的所有权转移给当前goroutine，但是锁处于不一致的状态中：mutexLocked没有设置
				// 并且我们将仍然被认为是waiter。这个状态需要被修复。
				if old&(mutexLocked|mutexWoken) != 0 || old>>mutexWaiterShift == 0 {
					throw("sync: inconsistent mutex state")
				}
				// 当前goroutine获取锁，waiter数量-1
				delta := int32(mutexLocked - 1<<mutexWaiterShift)
				// 如果当前goroutine非饥饿状态，或者说当前goroutine是队列中最后一个goroutine
				// 那么就退出饥饿模式，把状态设置为正常
				if !starving || old>>mutexWaiterShift == 1 {
					// 退出饥饿模式
					// 在这里这么做至关重要，还要考虑等待时间。
					// 饥饿模式是非常低效率的，一旦两个goroutine将互斥锁切换为饥饿模式，它们便可以无限锁。
					delta -= mutexStarving
				}
				// 原子的加上更新的值
				atomic.AddInt32(&m.state, delta)
				break
			}
			// 不是饥饿模式,就把当前的goroutine设为被唤醒
			awoke = true
			// 重置循环的次数
			iter = 0
		} else {
			// 如果CAS不成功，也就是说没能成功获得锁，锁被别的goroutine获得了或者锁一直没被释放
			// 那么就更新状态，重新开始循环尝试拿锁
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
【一份详细注释的go Mutex源码】http://cbsheng.github.io/posts/%E4%B8%80%E4%BB%BD%E8%AF%A6%E7%BB%86%E6%B3%A8%E9%87%8A%E7%9A%84go-mutex%E6%BA%90%E7%A0%81/  
【源码剖析 golang 中 sync.Mutex】https://www.purewhite.io/2019/03/28/golang-mutex-source/    