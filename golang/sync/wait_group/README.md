## wait_group

sync.WaitGroup 类型是并发安全的，也是开箱就能用的。
该类型有三个指针方法，即：Add,Done和Wait.

sync.WaitGroup是一个结构体类型。其中一个代表计数的字节数组类型的字段，该
字段用4字节表示给定计数，另4字节表示等待计数。当一个sync.WaitGroup类型的变
量被声明之后，其中的这两个计数都会是0。可以通过add方法增大或减少其中给定计数，
例如：
````
wg.Add(3)
````
或
````
wg.Add(-3)
````
需要注意的是，我们不能让这个计数值变成一个负数。
````
wg.Done()
````
相当于
````
wg.Add(-1)
````

我们知道add和done方法可以变更计数器的值，但是变更之后具体有什么
作用呢？


当调用sync.WaitGroup类型值的Wait方法时，它会去检查给定计数。如果
该计数为0，那么该方法会立即返回，且不会对程序产生任何影响。但是，
如果这个计数器大于0，该方法调用所在的那个goroutine就会阻塞，同时
等待计数器会加1。直到在该值的add或done方法被调用时给定计数变回0.
该值才会去唤醒因此等待而阻塞的所有goroutine，同时清零等待计数。


现在我们有一个案例：
假设程序启用了4个goroutine，分别是g1,g2,g3,g4。其中g2,g3,g4是由代码g1
启用的，g1启用之后并且要等待这些特殊任务的完成。

使用通道来进行阻塞

````
    sign := make(chan int, 3)
	go func() {
		sign <- 2
		fmt.Println(2)
	}()
	go func() {
		sign <- 3
		fmt.Println(3)
	}()

	go func() {
		sign <- 4
		fmt.Println(4)
	}()

	for i := 0; i < 3; i++ {
		fmt.Println("执行", <-sign)
	}
````

使用通道的有过于繁重了，原则上，我们不应该把通道当做互斥锁或信号量来使用。

使用waitGroup

````
	var wg sync.WaitGroup

	wg.Add(3)
	go func() {
		wg.Done()
		fmt.Println(2)
	}()
	go func() {
		wg.Done()
		fmt.Println(3)
	}()
	go func() {
		wg.Done()
		fmt.Println(4)
	}()

	wg.Wait()
	fmt.Println("1 2 3 4 end")
````



















