<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [互斥锁](#%E4%BA%92%E6%96%A5%E9%94%81)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [什么是sync.Mutex](#%E4%BB%80%E4%B9%88%E6%98%AFsyncmutex)
  - [分析下源码](#%E5%88%86%E6%9E%90%E4%B8%8B%E6%BA%90%E7%A0%81)
    - [Lock](#lock)
      - [位运算](#%E4%BD%8D%E8%BF%90%E7%AE%97)
    - [Unlock](#unlock)
  - [总结](#%E6%80%BB%E7%BB%93)
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

### 分析下源码

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

#### Lock

加锁基本上就这三种情况：  

1、可直接获取锁，直接加锁，返回；  

2、有冲突 首先自旋，如果其他`goroutine`在这段时间内释放了该锁，直接获得该锁；如果没有就走到下面3；  

3、有冲突，且已经过了自旋阶段，通过信号量进行阻塞；  

- 1、刚被唤醒的 加入到等待队列首部；  

- 2、新加入的 加入到等待队列的尾部。  

4、有冲突，根据不同的模式做处理； 

- 1、饥饿模式 获取锁  

- 2、正常模式 唤醒，继续循环，回到2  

```go
// Lock locks m.
// 如果锁正在使用中，新的goroutine请求，将被阻塞，直到锁被释放
func (m *Mutex) Lock() {
	// 原子的(cas)来判断是否加锁
	// 如果可以获取锁，直接加锁，返回
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
		// 第一个条件是state已被锁，但是不是饥饿状态。如果时饥饿状态，自旋时没有用的，锁的拥有权直接交给了等待队列的第一个。
		// 第二个条件是还可以自旋，多核、压力不大并且在一定次数内可以自旋， 具体的条件可以参考`sync_runtime_canSpin`的实现。
		// 如果满足这两个条件，不断自旋来等待锁被释放、或者进入饥饿状态、或者不能再自旋。
		if old&(mutexLocked|mutexStarving) == mutexLocked && runtime_canSpin(iter) {
			// 自旋的过程中如果发现state还没有设置woken标识，则设置它的woken标识， 并标记自己为被唤醒。
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

		// 到了这一步， state的状态可能是：
		// 1. 锁还没有被释放，锁处于正常状态
		// 2. 锁还没有被释放， 锁处于饥饿状态
		// 3. 锁已经被释放， 锁处于正常状态
		// 4. 锁已经被释放， 锁处于饥饿状态

		// new 复制 state的当前状态， 用来设置新的状态
		// old 是锁当前的状态
		new := old
		// 如果old state状态不是饥饿状态, new state 设置锁， 尝试通过CAS获取锁,
		// 如果old state状态是饥饿状态, 则不设置new state的锁，因为饥饿状态下锁直接转给等待队列的第一个.
		if old&mutexStarving == 0 {
			// 伪代码：newState = locked
			new |= mutexLocked
		}
		// 如果锁是被获取状态，或者饥饿状态
		// 就把期望状态中的等待队列的等待者数量+1(实际上是new + 8)
		if old&(mutexLocked|mutexStarving) != 0 {
			new += 1 << mutexWaiterShift
		}

		// 如果当前goroutine已经处于饥饿状态， 并且old state的已被加锁,
		// 将new state的状态标记为饥饿状态, 将锁转变为饥饿状态.
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
		// 注意new的锁标记不一定是true, 也可能只是标记一下锁的state是饥饿状态.
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

const (
	active_spin     = 4
)

// src/runtime/proc.go
// Active spinning for sync.Mutex.
// go:linkname sync_runtime_canSpin sync.runtime_canSpin
// go:nosplit
func sync_runtime_canSpin(i int) bool {
    // sync.Mutex是会被多个goroutine竞争的，所以自旋的次数需要控制
    // active_spin的值为4
    // 满足下面的添加才会发生自旋
    // 1、自旋的次数小于active_spin也就是4
    // 2、如果在单核的cpu是不能自旋的
    // 3、 GOMAXPROCS> 1，并且至少有一个其他正在运行的P，并且本地runq为空。
    // 4、当前P没有其它等待运行的G 
	if i >= active_spin || ncpu <= 1 || gomaxprocs <= int32(sched.npidle+sched.nmspinning)+1 {
		return false
	}
	if p := getg().m.p.ptr(); !runqempty(p) {
		return false
	}
	return true
}

// src/runtime/proc.go
// go:linkname sync_runtime_doSpin sync.runtime_doSpin
// go:nosplit
// procyield的实现是用汇编实现的
func sync_runtime_doSpin() {
	procyield(active_spin_cnt)
}

// src/runtime/asm_amd64.s
TEXT runtime·procyield(SB),NOSPLIT,$0-0
	MOVL	cycles+0(FP), AX
again:
        // 让加锁失败时cpu睡眠30个（about）clock，从而使得读操作的频率低很多。流水线重排的代价也会小很多
	PAUSE
	SUBL	$1, AX
	JNZ	again
	RET
```

梳理下流程  

1、原子的(cas)来判断是否加锁,如果之前锁没有被使用，当前`goroutine`获取锁，结束本次`Lock`操作；   

2、如果已经被别的`goroutine`持有了，启动一个for循环去抢占锁；  

会存在两种状态的切换 饥饿状态和正常状态  

如果一个等待的goroutine有超过1ms（写死在代码中）都没获取到锁，那么就会把锁转变为饥饿模式  

如果一个goroutine获取到了锁之后，它会判断以下两种情况：

- 1、它是队列中最后一个goroutine；

- 2、它拿到锁所花的时间小于1ms；

以上只要有一个成立，它就会把锁转变回正常模式。  

3、如果锁已经被锁了，并且不是饥饿状态，并且满足自旋的条件，当前goroutine会不断的进行自旋，等待锁被释放；  

4、不满足锁自旋的条件，然后结束自旋，这是当前锁的状态可能有下面几种情况：  

- 1、锁还没有被释放，锁处于正常状态 

- 2、锁还没有被释放， 锁处于饥饿状态

- 3、锁已经被释放， 锁处于正常状态

- 4、锁已经被释放， 锁处于饥饿状态  

5、如果`old.state`不是饥饿状态，新的`goroutine`尝试去获锁，如果是饥饿状态，就直接将锁直接转给等待队列的第一个；  

6、如果锁是被获取或饥饿状态，等待者的数量加一；  

7、当本`goroutine`被唤醒了，要么获得了锁，要么进入休眠；  

8、如果`old state`的状态是未被锁状态，并且锁不处于饥饿状态,那么当前`goroutine`已经获取了锁的拥有权，结束`Lock`;  

9、判断一下当前`goroutine`是新来的还是刚被唤醒的，新来的加入到等待队列的尾部，刚被唤醒的加入到等待队列的头部，然后通过信号量阻塞，直到当前`goroutine`被唤醒；  

10、判断如果当前`state`是否是饥饿状态，不是的唤醒本次`goroutine`,继续循环，是饥饿状态继续往下面走；  

11、饥饿状态，当前`goroutine`来设置锁，等待者减一，如果当前`goroutine`是队列中最后一个`goroutine`设置饥饿状态为正常，拿到锁结束`Lock`。  

<img src="/img/sync_mutex_lock.png" width = "625" height = "846" alt="mutex" align=center />

##### 位运算

上面有很多关于&和|的运算和判断，下面来具体的分析下  

```go
&      位运算 AND
|      位运算 OR
^      位运算 XOR
&^     位清空（AND NOT）
<<     左移
>>     右移
```

**&**  

参与运算的两数各对应的二进位相与，两个二进制位都为1时，结果才为1  

```go
    0101
AND 0011
  = 0001
```

**|**

参与运算的两数各对应的二进位相或，两个二进制位都为1时，结果才为0  

```go
   0101（十进制5）
OR 0011（十进制3）
 = 0111（十进制7）
```

**^** 

按位异或运算，对等长二进制模式或二进制数的每一位执行逻辑异或操作。操作的结果是如果某位不同则该位为1，否则该位为0。  

```go
    0101
XOR 0011
  = 0110
```

**&^**

将运算符左边数据相异的位保留，相同位清零  

```go
   0001 0100 
&^ 0000 1111 
 = 0001 0000  
```

**<<**

各二进位全部左移若干位，高位丢弃，低位补0  

```go
   0001（十进制1）
<<    3（左移3位）
 = 1000（十进制8）
```

**>>**

各二进位全部右移若干位，对无符号数，高位补0，有符号数，各编译器处理方法不一样，有的补符号位（算术右移），有的补0  

```go
   1010（十进制10）
>>    2（右移2位）
 = 0010（十进制2）
```

#### Unlock

```go
// Unlock unlocks m.
// 如果没有lock就去unlocak是会报错的  
//
//一个锁定的互斥锁与一个特定的goroutine没有关联。
// 它允许一个goroutine锁定一个互斥锁然后
// 安排另一个goroutine解锁它。
func (m *Mutex) Unlock() {
	if race.Enabled {
		_ = m.state
		race.Release(unsafe.Pointer(m))
	}

	// 修改state的状态
	new := atomic.AddInt32(&m.state, -mutexLocked)
	if new != 0 {
		// 不为0，说明没有成功解锁
		m.unlockSlow(new)
	}
}

func (m *Mutex) unlockSlow(new int32) {
	if (new+mutexLocked)&mutexLocked == 0 {
		throw("sync: unlock of unlocked mutex")
	}
	if new&mutexStarving == 0 {
		old := new
		for {
			// 如果说锁没有等待拿锁的goroutine
			// 或者锁被获取了(在循环的过程中被其它goroutine获取了)
			// 或者锁是被唤醒状态(表示有goroutine被唤醒，不需要再去尝试唤醒其它goroutine)
			// 或者锁是饥饿模式(会直接转交给队列头的goroutine)
			// 那么就直接返回，啥都不用做了

			// 也就是没有等待的goroutine, 或者锁不处于空闲的状态，直接返回.
			if old>>mutexWaiterShift == 0 || old&(mutexLocked|mutexWoken|mutexStarving) != 0 {
				return
			}

			// 走到这一步的时候，说明锁目前还是空闲状态，并且没有goroutine被唤醒且队列中有goroutine等待拿锁
			// 将等待的goroutine数减一，并设置woken标识
			new = (old - 1<<mutexWaiterShift) | mutexWoken
			// 设置新的state, 这里通过信号量会唤醒一个阻塞的goroutine去获取锁.
			if atomic.CompareAndSwapInt32(&m.state, old, new) {
				runtime_Semrelease(&m.sema, false, 1)
				return
			}
			old = m.state
		}
	} else {
		// 饥饿模式下， 直接将锁的拥有权传给等待队列中的第一个.
		// 注意此时state的mutexLocked还没有加锁，唤醒的goroutine会设置它。
		// 在此期间，如果有新的goroutine来请求锁， 因为mutex处于饥饿状态， mutex还是被认为处于锁状态，
		// 新来的goroutine不会把锁抢过去.
		runtime_Semrelease(&m.sema, true, 1)
	}
}
```

梳理下流程：  

1、首先判断如果之前是锁的状态是未加锁，`Unlock`将会触发`panic`；  

2、如果当前锁是正常模式，一个for循环，去不断尝试解锁；  

3、饥饿模式下，通过信号量，唤醒在饥饿模式下面`Lock`操作下队列中第一个`goroutine`。  

<img src="/img/sync_mutex_unlock.png" width = "331" height = "486" alt="mutex" align=center />

### 总结

1、加锁的过程会存在正常模式和互斥模式的转换；  

2、饥饿模式就是保证锁的公平性，正常模式下的互斥锁能够提供更好地性能，饥饿模式的能避免 Goroutine 由于陷入等待无法获取锁而造成的高尾延时；  

3、锁的状态的转换，也使用到了位运算；  

4、一个已经锁定的互斥锁，允许其他协程进行解锁，不过只能被解锁一次；  

### 参考

【sync.Mutex 源码分析】https://reading.hidevops.io/articles/sync/sync_mutex_source_code_analysis/  
【一份详细注释的go Mutex源码】http://cbsheng.github.io/posts/%E4%B8%80%E4%BB%BD%E8%AF%A6%E7%BB%86%E6%B3%A8%E9%87%8A%E7%9A%84go-mutex%E6%BA%90%E7%A0%81/  
【源码剖析 golang 中 sync.Mutex】https://www.purewhite.io/2019/03/28/golang-mutex-source/   
【sync.mutex 源代码分析】https://colobu.com/2018/12/18/dive-into-sync-mutex/    
【源码剖析 golang 中 sync.Mutex】https://www.purewhite.io/2019/03/28/golang-mutex-source/  

>**本文作者**：liz  
>**本文链接**：https://boilingfrog.github.io/2021/03/14/sync.Mutex/    
>**版权声明**：本文为博主原创文章，遵循 [CC 4.0 BY-SA](https://creativecommons.org/licenses/by-sa/4.0/) 版权协议，转载请附上原文出处链接和本声明。    