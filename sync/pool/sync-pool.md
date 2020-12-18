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






 




### 参考
【深入Golang之sync.Pool详解】https://www.cnblogs.com/sunsky303/p/9706210.html  
【由浅入深聊聊Golang的sync.Pool】https://juejin.cn/post/6844903903046320136   
【Golang sync.Pool源码阅读与分析】https://jiajunhuang.com/articles/2020_05_05-go_sync_pool.md.html  