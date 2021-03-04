<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [atomic](#atomic)
  - [原子操作](#%E5%8E%9F%E5%AD%90%E6%93%8D%E4%BD%9C)
  - [Go中原子操作的支持](#go%E4%B8%AD%E5%8E%9F%E5%AD%90%E6%93%8D%E4%BD%9C%E7%9A%84%E6%94%AF%E6%8C%81)
    - [CompareAndSwap(CAS)](#compareandswapcas)
    - [Swap(交换)](#swap%E4%BA%A4%E6%8D%A2)
    - [Add(增加或减少)](#add%E5%A2%9E%E5%8A%A0%E6%88%96%E5%87%8F%E5%B0%91)
    - [Load(原子读取)](#load%E5%8E%9F%E5%AD%90%E8%AF%BB%E5%8F%96)
    - [Store(原子写入)](#store%E5%8E%9F%E5%AD%90%E5%86%99%E5%85%A5)
  - [原子操作与互斥锁的区别](#%E5%8E%9F%E5%AD%90%E6%93%8D%E4%BD%9C%E4%B8%8E%E4%BA%92%E6%96%A5%E9%94%81%E7%9A%84%E5%8C%BA%E5%88%AB)
  - [atomic.Value](#atomicvalue)
    - [Load](#load)
    - [Store](#store)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## atomic

### 原子操作

原子操作即是进行过程中不能被中断的操作，针对某个值的原子操作在被进行的过程中，CPU绝不会再去进行其他的针对该值的操作。为了实现这样的严谨性，原
子操作仅会由一个独立的CPU指令代表和完成。原子操作是无锁的，常常直接通过CPU指令直接实现。 事实上，其它同步技术的实现常常依赖于原子操作。  

具体的原子操作在不同的操作系统中实现是不同的。比如在Intel的CPU架构机器上，主要是使用总线锁的方式实现的。 大致的意思就是当一个CPU需要操作一个
内存块的时候，向总线发送一个LOCK信号，所有CPU收到这个信号后就不对这个内存块进行操作了。 等待操作的CPU执行完操作后，发送UNLOCK信号，才结束。 
在AMD的CPU架构机器上就是使用MESI一致性协议的方式来保证原子操作。 所以我们在看atomic源码的时候，我们看到它针对不同的操作系统有不同汇编语言
文件。  

### Go中原子操作的支持

Go语言的`sync/atomic`提供了对原子操作的支持，用于同步访问整数和指针。  

- Go语言提供的原子操作都是非入侵式的
- 原子操作支持的类型包括`int32、int64、uint32、uint64、uintptr、unsafe.Pointer`。  

竞争条件是由于异步的访问共享资源，并试图同时读写该资源而导致的，使用互斥锁和通道的思路都是在线程获得到访问权后阻塞其他线程对共享内存的访问，而使用原子操作解决数据竞争问题则是利用了其不可被打断的特性。  

#### CompareAndSwap(CAS)

go中的Cas操作，是借用了CPU提供的原子性指令来实现。CAS操作修改共享变量时候不需要对共享变量加锁，而是通过类似乐观锁的方式进行检查，本质还是不断的占用CPU 资源换取加锁带来的开销（比如上下文切换开销）。    

原子操作中的CAS(Compare And Swap),在`sync/atomic`包中，这类原子操作由名称以`CompareAndSwap`为前缀的若干个函数提供    

```go
func CompareAndSwapInt32(addr *int32, old, new int32) (swapped bool)
func CompareAndSwapInt64(addr *int64, old, new int64) (swapped bool)
func CompareAndSwapPointer(addr *unsafe.Pointer, old, new unsafe.Pointer) (swapped bool)
func CompareAndSwapUint32(addr *uint32, old, new uint32) (swapped bool)
func CompareAndSwapUint64(addr *uint64, old, new uint64) (swapped bool)
func CompareAndSwapUintptr(addr *uintptr, old, new uintptr) (swapped bool)
```

`CompareAndSwap`函数会先判断参数addr指向的操作值与参数old的值是否相等，仅当此判断得到的结果是true之后，才会用参数new代表的新值替换掉原先的旧值，否则操作就会被忽略。  

查看下源码，这几个代码差不多，以`CompareAndSwapUint32`为例子,golang主要还是依赖汇编来来实现的原子操作，不同的CPU架构是有对应不同的.s汇编文件的。  

`/usr/local/go/src/sync/atomic/asm.s`

```cgo
TEXT ·CompareAndSwapUint32(SB),NOSPLIT,$0
	JMP	runtime∕internal∕atomic·Cas(SB)
```

看下汇编的`Cas`  

```cgo
// bool Casp1(void **val, void *old, void *new)
// Atomically:
//	if(*val == old){
//		*val = new;
//		return 1;
//	} else
//		return 0;
TEXT runtime∕internal∕atomic·Casp1(SB), NOSPLIT, $0-25
        // 首先将 ptr 的值放入 BX
	MOVQ	ptr+0(FP), BX
        // 将假设的旧值放入 AX
	MOVQ	old+8(FP), AX
        // 需要比较的新值放入到CX
	MOVQ	new+16(FP), CX
	LOCK
	CMPXCHGQ	CX, 0(BX)
	SETEQ	ret+24(FP)
	RET
```
> MOV 指令有有好几种后缀 MOVB MOVW MOVL MOVQ 分别对应的是 1 字节 、2 字节 、4 字节、8 字节

`TEXT runtime∕internal∕atomic·Cas(SB),NOSPLIT,$0-17`，`$0-17`表示的意思是这个`TEXT block`运行的时候，需要开辟的栈帧大小是0，而`17 = 8 + 4 + 4 + 1 = sizeof(pointer of int32) + sizeof(int32) + sizeof(int32) + sizeof(bool)`（返回值是 bool ，占据 1 个字节）  

`FP`，是伪寄存器(pseudo) ，里边存的是 `Frame Pointer`, `FP`配合偏移 可以指向函数调用参数或者临时变量  

`MOVQ ptr+0(FP) BX` 这一句话是指把函数的第一个参数`ptr+0(FP)`移动到`BX`寄存器中  

`MOVQ代`表移动的是8个字节,Q 代表`64bit` ，参数的引用是 参数名称+偏移(FP),可以看到这里名称用了`ptr`,并不是v`al`,变量名对汇编不会有什么影响，但是语法上是必须带上的，可读性也会更好些。  

`LOCK`并不是指令，而是一个指令的前缀`(instruction prefix)`，是用来修饰`CMPXCHGL CX,0(BX)` 的  

> The LOCK prefix ensures that the CPU has exclusive ownership of the appropriate cache line for the duration of the operation, and provides certain additional ordering guarantees. This may be achieved by asserting a bus lock, but the CPU will avoid this where possible. If the bus is locked then it is only for the duration of the locked instruction

`CMPXCHGL` 有两个操作数，`CX 和 0(BX)`,`0(BX)`代表的是`val`的地址。  

`CMPXCHGL`指令做的事情，首先会把`0(BX)`里的值和`AX`寄存器里存的值做比较，如果一样的话会把`CX`里边存的值保存到`0(BX)`这块地址里 (虽然这条指令里并没有出现`AX`，但是还是用到了，汇编里还是有不少这样的情况)   

`SETEQ` 会在`AX`和`CX`相等的时候把1写进 `ret+16(FP)`(否则写 0）

看下如何使用  

```go
func main() {
	var a, b int32 = 13, 13
	var c int32 = 9
	res := atomic.CompareAndSwapInt32(&a, b, c)
	fmt.Println("swapped:", res)
	fmt.Println("替换的值:", c)
	fmt.Println("替换之后a的值:", a)
}
```

查看下输出  

```go
swapped: true
替换的值: 9
替换之后a的值: 9
```

`a`值和`b`值作比较，当`a`和`b`相等时，会用`c`的值替换掉`a`的值  

我们使用的`mutex`互斥锁类似悲观锁，总是假设会有并发的操作要修改被操作的值，所以使用锁将相关操作放入到临界区加以保存。而CAS操作做法趋于乐观锁，总是假设被操作的值未曾改变（即与旧值相等），并一旦确认这个假设的真实性就立即进行值替换。在被操作值被频繁变更的情况下，`CAS`操作并不那么容易成功所以需要不断进行尝试，直到成功为止。    

举个栗子  

```go
func main() {
	fmt.Println("======old value=======")
	fmt.Println(value)
	addValue(10)
	fmt.Println("======New value=======")
	fmt.Println(value)

}

//不断地尝试原子地更新value的值,直到操作成功为止
func addValue(delta int32) {
	for {
		v := value
		if atomic.CompareAndSwapInt32(&value, v, v+delta) {
			break
		}
	}
}
```

#### Swap(交换)

上面的`CompareAndSwap`系列的函数需要比较后再进行交换，也有不需要进行比较就进行交换的原子操作。  

```go
func SwapInt32(addr *int32, new int32) (old int32)
func SwapInt64(addr *int64, new int64) (old int64)
func SwapPointer(addr *unsafe.Pointer, new unsafe.Pointer) (old unsafe.Pointer)
func SwapUint32(addr *uint32, new uint32) (old uint32)
func SwapUint64(addr *uint64, new uint64) (old uint64)
func SwapUintptr(addr *uintptr, new uintptr) (old uintptr)
```

几个差不多，来看下`SwapInt32`的源码，也是通过汇编来实现的  

`/usr/local/go/src/sync/atomic/asm.s`

```cgo
TEXT ·SwapUint32(SB),NOSPLIT,$0
	JMP	runtime∕internal∕atomic·Xchg(SB)
```

看下汇编的`Xchg`

```cgo
TEXT runtime∕internal∕atomic·Xchg(SB), NOSPLIT, $0-20
	MOVQ	ptr+0(FP), BX
	MOVL	new+8(FP), AX
        // 原子操作, 把_value的值和newValue交换, 且返回_value原来的值
	XCHGL	AX, 0(BX)
	MOVL	AX, ret+16(FP)
	RET
```

举个栗子  

```go
func main() {
	var a, b int32 = 13, 12
	old := atomic.SwapInt32(&a, b)
	fmt.Println("old的值:", old)
	fmt.Println("替换之后a的值", a)
}
```

查看下输出  

```go
old的值: 13
替换之后a的值 12
```

#### Add(增加或减少)

对一个数值进行增加或者减少的行为也需要保证是原子的，它对应于atomic包的函数就是

```go
func AddInt32(addr *int32, delta int32) (new int32)
func AddInt64(addr *int64, delta int64) (new int64)
func AddUint32(addr *uint32, delta uint32) (new uint32)
func AddUint64(addr *uint64, delta uint64) (new uint64)
func AddUintptr(addr *uintptr, delta uintptr) (new uintptr)
```

举个栗子  

```go
func main() {
	var a int32 = 13
	addValue := atomic.AddInt32(&a, 1)
	fmt.Println("增加之后:", addValue)
	delValue := atomic.AddInt32(&a, -4)
	fmt.Println("减少之后:", delValue)
}
```

查看下输出  

```go
增加之后: 14
减少之后: 10
```

#### Load(原子读取)

当我们要读取一个变量的时候，很有可能这个变量正在被写入，这个时候，我们就很有可能读取到写到一半的数据。 所以读取操作是需要一个原子行为的。
在atomic包中就是Load开头的函数群。  

```go
func LoadInt32(addr *int32) (val int32)
func LoadInt64(addr *int64) (val int64)
func LoadPointer(addr *unsafe.Pointer) (val unsafe.Pointer)
func LoadUint32(addr *uint32) (val uint32)
func LoadUint64(addr *uint64) (val uint64)
func LoadUintptr(addr *uintptr) (val uintptr)
```

#### Store(原子写入)

读取是有原子性的操作的，同样写入atomic包也提供了相关的操作包。

```go
func StoreInt32(addr *int32, val int32)
func StoreInt64(addr *int64, val int64)
func StorePointer(addr *unsafe.Pointer, val unsafe.Pointer)
func StoreUint32(addr *uint32, val uint32)
func StoreUint64(addr *uint64, val uint64)
func StoreUintptr(addr *uintptr, val uintptr)
```

### 原子操作与互斥锁的区别

首先atomic操作的优势是更轻量，比如CAS可以在不形成临界区和创建互斥量的情况下完成并发安全的值替换操作。这可以大大的减少同步对程序性能的损耗。  

原子操作也有劣势。还是以CAS操作为例，使用CAS操作的做法趋于乐观，总是假设被操作值未曾被改变（即与旧值相等），并一旦确认这个假设的真实性就立即进行值替换，那么在被操作值被频繁变更的情况下，CAS操作并不那么容易成功。而使用互斥锁的做法则趋于悲观，我们总假设会有并发的操作要修改被操作的值，并使用锁将相关操作放入临界区中加以保护。    

下面是几点区别：  

- 互斥锁是一种数据结构，用来让一个线程执行程序的关键部分，完成互斥的多个操作
- 原子操作是无锁的，常常直接通过CPU指令直接实现  
- 原子操作中的cas趋于乐观锁，CAS操作并不那么容易成功，需要判断，然后尝试处理
- 可以把互斥锁理解为悲观锁，共享资源每次只给一个线程使用，其它线程阻塞，用完后再把资源转让给其它线程  

`atomic`包提供了底层的原子性内存原语，这对于同步算法的实现很有用。这些函数一定要非常小心地使用，使用不当反而会增加系统资源的开销，对于应用层来说，最好使用通道或sync包中提供的功能来完成同步操作。    

针对`atomic`包的观点在Google的邮件组里也有很多讨论，其中一个结论解释是：  

> 应避免使用该包装。或者，阅读C ++ 11标准的“原子操作”一章；如果您了解如何在C ++中安全地使用这些操作，那么你才能有安全地使用Go的sync/atomic包的能力。

### atomic.Value

此类型的值相当于一个容器，可以被用来“原子地"存储（Store）和加载（Load）任意类型的值。当然这个类型也是原子性的。   

有了`atomic.Value`这个类型，这样用户就可以在不依赖`Go`内部类型`unsafe.Pointer`的情况下使用到atomic提供的原子操作。  

分析下源码  

```cgo
// A Value must not be copied after first use.
type Value struct {
	v interface{}
}
```

里面主要是包含了两个方法  

- `v.Store(c)` - 写操作，将原始的变量c存放到一个`atomic.Value`类型的v里。

- `c = v.Load()` - 读操作，从线程安全的v中读取上一步存放的内容。  

#### Load

```go
// ifaceWords is interface{} internal representation.
type ifaceWords struct {
	// 类型
	typ unsafe.Pointer
	// 数据
	data unsafe.Pointer
}

// 如果没Store将返回nil
func (v *Value) Load() (x interface{}) {
	// 获得 interface 结构的指针
	vp := (*ifaceWords)(unsafe.Pointer(v))
	// 获取类型
	typ := LoadPointer(&vp.typ)
	// 判断，第一次写入还没有开始，或者还没完成，返回nil
	if typ == nil || uintptr(typ) == ^uintptr(0) {
		// First store not yet completed.
		return nil
	}
	// 获得存储值的实际数据
	data := LoadPointer(&vp.data)
	// 将复制得到的 typ 和 data 给到 x
	xp := (*ifaceWords)(unsafe.Pointer(&x))
	xp.typ = typ
	xp.data = data
	return
}
```

1、Load中也是借助于`atomic.LoadPointer`来实现的；  

2、使用了`Go`运行时类型系统中的`interface{}`这一类型本质上由 两段内容组成，一个是类型`typ`区域，另一个是实际数据d`ata`区域；  

3、保证与原子性，加入了一个判断：  

- typ为nil表示还没有写入值  

- `uintptr(typ) == ^uintptr(0)`表示有第一次写入还没有完成  

#### Store

```go
// 如果两次Store的类型不同将会panic
// 如果写入nil，也会panic
func (v *Value) Store(x interface{}) {
	// value不能为nil
	if x == nil {
		panic("sync/atomic: store of nil value into Value")
	}
	// Value存储的指针
	vp := (*ifaceWords)(unsafe.Pointer(v))
	// 写入value的目标指针x
	xp := (*ifaceWords)(unsafe.Pointer(&x))
	for {
		typ := LoadPointer(&vp.typ)
		// 第一次Store
		if typ == nil {
			// 禁止抢占当前 Goroutine 来确保存储顺利完成
			runtime_procPin()
			// 如果typ为nil，设置一个标志位，宣告正在有人操作此值
			if !CompareAndSwapPointer(&vp.typ, nil, unsafe.Pointer(^uintptr(0))) {
				// 如果没有成功，取消不可抢占，下次再试
				runtime_procUnpin()
				continue
			}
			// 如果标志位设置成功，说明其他人都不会向 interface{} 中写入数据
			// 这点细品很巧妙，先写数据，在写类型，应该类型设置了不可写入的表示位
			// 写入数据
			StorePointer(&vp.data, xp.data)
			// 写入类型
			StorePointer(&vp.typ, xp.typ)
			// 存储成功，取消不可抢占，，直接返回
			runtime_procUnpin()
			return
		}
		// 已经有值写入了，或者有正在写入的Goroutine

		// 有其他 Goroutine 正在对 v 进行写操作
		if uintptr(typ) == ^uintptr(0) {
			continue
		}

		// 如果本次存入的类型与前次存储的类型不同
		if typ != xp.typ {
			panic("sync/atomic: store of inconsistently typed value into Value")
		}
		// 类型已经写入，直接保存数据
		StorePointer(&vp.data, xp.data)
		return
	}
}
```

梳理下流程：  

1、首先判断类型如果为nil直接panic；  

2、然后通过有个for循环来连续判断是否可以进行值的写入；  

3、如果是`typ == nil`表示是第一次写入,然后给type设置一个标识位，来表示有goroutine正在写入；  

4、然后写入值，退出；  

5、如果type不为nil，但是等于标识位，表示有正在写入的goroutine，然后继续循环；    

6、最后type不为nil，并且不等于标识位，并且和value里面的type类型一样，写入内容，然后退出。    

注意：其中使用了`runtime_procPin()`方法，它可以将一个`goroutine`死死占用当前使用的`P(P-M-G中的processor)`，不允许其它`goroutine/M`抢占,这样就能保证存储顺利完成，不必担心竞争的问题。释放pin的方法是`runtime_procUnpin`。  

<img src="/img/atomic_store_1.png" width = "584" height = "586" alt="atomic" align=center />


### 参考
【Go并发编程之美-CAS操作】https://zhuanlan.zhihu.com/p/56733484  
【sync/atomic - 原子操作】https://docs.kilvn.com/The-Golang-Standard-Library-by-Example/chapter16/16.02.html  
【Go语言的原子操作和互斥锁的区别】https://studygolang.com/articles/29240  
【Package atomic】https://go-zh.org/pkg/sync/atomic/  
【Go 语言标准库中 atomic.Value 的前世今生】https://blog.betacat.io/post/golang-atomic-value-exploration/   
【原子操作】https://golang.design/under-the-hood/zh-cn/part4lib/ch15sync/atomic/   
【关于Go语言中的go:linkname】https://blog.csdn.net/IT_DREAM_ER/article/details/103590944  
【原子操作使用】https://www.kancloud.cn/digest/batu-go/153537   
【Go源码解析之atomic】https://amazingao.com/posts/2020/11/go-src/sync/atomic/  
【Plan 9 汇编语言】https://golang.design/under-the-hood/zh-cn/part1basic/ch01basic/asm/  
