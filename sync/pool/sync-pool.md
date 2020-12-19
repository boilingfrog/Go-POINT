## sync.pool

### sync.pool作用

有时候我们为了优化GC的场景，减少并复用内存，我们可以使用 sync.Pool 来复用需要频繁创建临时对象。 

#### 使用

一个小demo

```go
func main() {
	// 初始化一个pool
	pool := &sync.Pool{
		// 默认的返回值设置，不写这个参数，默认是nil
		New: func() interface{} {
			return 0
		},
	}

	// 看一下初始的值，这里是返回0，如果不设置New函数，默认返回nil
	init := pool.Get()
	fmt.Println("初始值", init)

	// 设置一个参数1
	pool.Put(1)

	// 获取查看结果
	num := pool.Get()
	fmt.Println("put之后取值", num)

	// 再次获取，会发现，已经是空的了，只能返回默认的值。
	num = pool.Get()
	fmt.Println("put之后再次取值", num)
}
```

输出

```
初始值 0
put之后取值 1
put之后再次取值 0
```

go本身也用到了sync.pool，例如`fmt.Printf`


```
// Sprintf formats according to a format specifier and returns the resulting string.
func Sprintf(format string, a ...interface{}) string {
	p := newPrinter()
	p.doPrintf(format, a)
	s := string(p.buf)
	p.free()
	return s
}

var ppFree = sync.Pool{
	New: func() interface{} { return new(pp) },
}

// newPrinter allocates a new pp struct or grabs a cached one.
func newPrinter() *pp {
	p := ppFree.Get().(*pp)
	p.panicking = false
	p.erroring = false
	p.wrapErrs = false
	p.fmt.init(&p.buf)
	return p
}
```

### 源码解读

pool结构

```go
type Pool struct {
    // 用来标记，当前的 struct 是不能够被 copy 的
    noCopy noCopy
    // P 个固定大小的 poolLocal 数组，每个 P 拥有一个空间
    local     unsafe.Pointer // local fixed-size per-P pool, actual type is [P]poolLocal
    // 上面数组的大小，即 P 的个数
    localSize uintptr        // size of the local array

    // 同 local 和 localSize，只是在 gc 的过程中保留一次
    victim     unsafe.Pointer // local from previous cycle
    victimSize uintptr        // size of victims array

    // 自定义一个 New 函数，然后可以在 Get 不到东西时，自动创建一个
    New func() interface{}
}
```

```go
// Local per-P Pool appendix.
type poolLocalInternal struct {
    // private 存储一个 Put 的数据，pool.Put() 操作优先存入 private，如果private有信息，才会存入 shared
    private interface{} // Can be used only by the respective P.
    // 存储一个链表，用来维护 pool.Put() 操作加入的数据，每个 P 可以操作自己 shared 链表中的头部，而其他的 P 在用完自己的 shared 时，可能会来偷数据，从而操作链表的尾部
    shared  poolChain   // Local P can pushHead/popHead; any P can popTail.
}
```

```go
// unsafe.Sizeof(poolLocal{})  // 128 byte(1byte = 8 bits)
// unsafe.Sizeof(poolLocalInternal{})  // 32 byte(1byte = 8 bits)
type poolLocal struct {
    // 每个P对应的pool
    poolLocalInternal

     // 将 poolLocal 补齐至两个缓存行的倍数，防止 false sharing,
    // 每个缓存行具有 64 bytes，即 512 bit
    // 目前我们的处理器一般拥有 32 * 1024 / 64 = 512 条缓存行
    // 伪共享，仅占位用，防止在 cache line 上分配多个 poolLocalInternal
    pad [128 - unsafe.Sizeof(poolLocalInternal{})%128]byte
}
```


### 参考
【深入Golang之sync.Pool详解】https://www.cnblogs.com/sunsky303/p/9706210.html  
【由浅入深聊聊Golang的sync.Pool】https://juejin.cn/post/6844903903046320136   
【Golang sync.Pool源码阅读与分析】https://jiajunhuang.com/articles/2020_05_05-go_sync_pool.md.html  
【【Go夜读】sync.Pool 源码阅读及适用场景分析】https://www.jianshu.com/p/f61bfe89e473  
【golang的对象池sync.pool源码解读】https://zhuanlan.zhihu.com/p/99710992  
【15.5 缓存池】https://golang.design/under-the-hood/zh-cn/part4lib/ch15sync/pool/  
