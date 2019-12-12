## goroutine

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
