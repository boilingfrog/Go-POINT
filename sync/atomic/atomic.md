<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [atomic](#atomic)
  - [原子操作](#%E5%8E%9F%E5%AD%90%E6%93%8D%E4%BD%9C)
  - [Go中原子操作的支持](#go%E4%B8%AD%E5%8E%9F%E5%AD%90%E6%93%8D%E4%BD%9C%E7%9A%84%E6%94%AF%E6%8C%81)
    - [CAS](#cas)
    - [增加或减少](#%E5%A2%9E%E5%8A%A0%E6%88%96%E5%87%8F%E5%B0%91)
    - [读取或写入](#%E8%AF%BB%E5%8F%96%E6%88%96%E5%86%99%E5%85%A5)
  - [原子操作与互斥锁的区别](#%E5%8E%9F%E5%AD%90%E6%93%8D%E4%BD%9C%E4%B8%8E%E4%BA%92%E6%96%A5%E9%94%81%E7%9A%84%E5%8C%BA%E5%88%AB)
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

#### CAS

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



调用函数后，`CompareAndSwap`函数会先判断参数addr指向的操作值与参数old的值是否相等，仅当此判断得到的结果是true之后，才会用参数new代表的新值替换掉原先的旧值，否则操作就会被忽略。  

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

这一系列的函数需要比较后再进行交换，也有不需要进行比较就进行交换的原子操作。  

```go
func SwapInt32(addr *int32, new int32) (old int32)
func SwapInt64(addr *int64, new int64) (old int64)
func SwapPointer(addr *unsafe.Pointer, new unsafe.Pointer) (old unsafe.Pointer)
func SwapUint32(addr *uint32, new uint32) (old uint32)
func SwapUint64(addr *uint64, new uint64) (old uint64)
func SwapUintptr(addr *uintptr, new uintptr) (old uintptr)
```

竞争条件是由于异步的访问共享资源，并试图同时读写该资源而导致的，使用互斥锁和通道的思路都是在线程获得到访问权后阻塞其他线程对共享内存的访问，
而使用原子操作解决数据竞争问题则是利用了其不可被打断的特性。

#### 增加或减少

对一个数值进行增加或者减少的行为也需要保证是原子的，它对应于atomic包的函数就是

```go
func AddInt32(addr *int32, delta int32) (new int32)
func AddInt64(addr *int64, delta int64) (new int64)
func AddUint32(addr *uint32, delta uint32) (new uint32)
func AddUint64(addr *uint64, delta uint64) (new uint64)
func AddUintptr(addr *uintptr, delta uintptr) (new uintptr)
```

#### 读取或写入

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

原子操作也有劣势。还是以CAS操作为例，使用CAS操作的做法趋于乐观，总是假设被操作值未曾被改变（即与旧值相等），并一旦确认这个假设的真实性就立即
进行值替换，那么在被操作值被频繁变更的情况下，CAS操作并不那么容易成功。而使用互斥锁的做法则趋于悲观，我们总假设会有并发的操作要修改被操作的值，
并使用锁将相关操作放入临界区中加以保护。  

下面是几点区别：  

- 互斥锁是一种数据结构，用来让一个线程执行程序的关键部分，完成互斥的多个操作
- 原子操作是无锁的，常常直接通过CPU指令直接实现  
- 原子操作中的cas趋于乐观锁，CAS操作并不那么容易成功，需要判断，然后尝试处理
- 可以把互斥锁理解为悲观锁，共享资源每次只给一个线程使用，其它线程阻塞，用完后再把资源转让给其它线程  

atomic包提供了底层的原子性内存原语，这对于同步算法的实现很有用。这些函数一定要非常小心地使用，使用不当反而会增加系统资源的开销，对于应用层
来说，最好使用通道或sync包中提供的功能来完成同步操作。  

针对atomic包的观点在Google的邮件组里也有很多讨论，其中一个结论解释是：  

> 应避免使用该包装。或者，阅读C ++ 11标准的“原子操作”一章；如果您了解如何在C ++中安全地使用这些操作，那么你才能有安全地使用Go的sync/atomic包的能力。

### 参考
【Go并发编程之美-CAS操作】https://zhuanlan.zhihu.com/p/56733484  
【sync/atomic - 原子操作】https://docs.kilvn.com/The-Golang-Standard-Library-by-Example/chapter16/16.02.html  
【Go语言的原子操作和互斥锁的区别】https://studygolang.com/articles/29240  
【Package atomic】https://go-zh.org/pkg/sync/atomic/  
【Go 语言标准库中 atomic.Value 的前世今生】https://blog.betacat.io/post/golang-atomic-value-exploration/   
【原子操作】https://golang.design/under-the-hood/zh-cn/part4lib/ch15sync/atomic/   
【关于Go语言中的go:linkname】https://blog.csdn.net/IT_DREAM_ER/article/details/103590944  
【原子操作使用】https://www.kancloud.cn/digest/batu-go/153537   
