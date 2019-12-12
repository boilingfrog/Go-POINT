
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
这个函数，我们使用go函数启动了3个goroutine。它们的执行是并发的。也就是说，这三个函数
的执行，我们是无法判断它们执行的先后顺序的。
````
GOROOT=/usr/local/go #gosetup
GOPATH=/home/liz/go #gosetup
/usr/local/go/bin/go build -o /tmp/___go_build_Go_POINT_goroutine Go-POINT/goroutine #gosetup
/tmp/___go_build_Go_POINT_goroutine #gosetup

Process finished with exit code 0
````
我们看到结果什么也没有打印，正常应该是 有 1 2 3,虽然它们的顺序我们没有办法确定，但是
至少应该出现在打印的结果里面。

是这样的，main函数也是一个goroutine，当我们main函数调用case2启动后面的3个go函数的
时候，这时候是用4个go函数在并发的执行，所以如果main函数在调度器中先被执行了，马上就
会结束掉，就看不到后面启动的这个3个go函数的打印了。

验证下，可以把main函数阻塞一下
````
func main() {
	// case1()
	case2()
	time.Sleep(time.Millisecond)
}

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
输出的信息
````
GOROOT=/usr/local/go #gosetup
GOPATH=/home/liz/go #gosetup
/usr/local/go/bin/go build -o /tmp/___go_build_Go_POINT_goroutine Go-POINT/goroutine #gosetup
/tmp/___go_build_Go_POINT_goroutine #gosetup
goroutine 1
goroutine 2
goroutine 3

Process finished with exit code 0

````
我们看到1 2 3已经可以看到打印的信息了。，当我们多次执行，可以看到执行的顺序是不确定的
这也验证了上面说的，go函数并发执行，但谁先谁后不确定。

再看一个例子：
````
func case3() {
	name := "小白"
	go func() {
		fmt.Println(name)
	}()
	name = "小李"
	time.Sleep(time.Millisecond)
}
````
上面的输出什么呢，是小李还是小白呢
我们来试下
````
/usr/local/go/bin/go build -o /tmp/___go_build_Go_POINT_goroutine Go-POINT/goroutine #gosetup
/tmp/___go_build_Go_POINT_goroutine #gosetup
小白

Process finished with exit code 0
````
为什么呢？
因为最后的sleep，在把小白赋值给name之后，才sleep。当go函数去执行的发现
name已经变成了小白，然后就打印出了小白。

我们换下位置，如下
````
func case4() {
	name := "小白"
	go func() {
		fmt.Println(name)
	}()
	time.Sleep(time.Millisecond)
	name = "小李"
}
````
我们发现a打印的结果变成了小白。
因为在改变变量之前，就给那个go函数执行的机会了。

接着看，当我们处理多个值的时候。

````
func case5() {
	names := []string{"小白", "小明", "小红", "小张"}
	for _, name := range names {
		go func() {
			fmt.Println("名字", name)
		}()
	}
	time.Sleep(time.Millisecond)
}
````
