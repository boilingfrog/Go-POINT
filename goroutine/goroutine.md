
- [什么是goroutine](#%e4%bb%80%e4%b9%88%e6%98%afslice)
   - [go和goroutine](#slice%e7%9a%84%e5%88%9b%e5%bb%ba%e4%bd%bf%e7%94%a8)
 



## 什么是goroutine

go中的go关键字是用来启动goroutine的唯一途径。

### go和goroutine

一条go语句意味着一个函数的并发执行。当我们需要启动一个goroutine的时候非常简单只需要添加一个go
关键字就可以了。

````
go fmt.Println("goroutine")
````

如果是匿名函数的话，就像下面的这样

````
go func() {
   fmt.Println("goroutine")
}()
````
我们需要注意的是，无论我们是否选择给匿名函数传递参数，后面的圆括号，我们不能忘记，当然编译器也
会给我们错误的语法提示。它们代表了对函数的调用，也是调用表达式的重要组成部分。

Go运行的时候对于go函数是并发执行的。也就是当go语句执行的时候，go函数会被单独放到一个goroutine
中。之后这个go函数的执行会独立于当前goroutine的运行。

````
func case2() {
	go func() {
		fmt.Println("goroutine", 1)
	}()
	go func() {
		fmt.Println("goroutine", 2)
	}()
	go func() {
		fmt.Println("goroutine", 3)
	}()
}
````
这个函数，我们使用go函数启动了3个goroutine。它们的执行是并发的。