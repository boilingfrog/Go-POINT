<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [panic源码解读](#panic%E6%BA%90%E7%A0%81%E8%A7%A3%E8%AF%BB)
  - [panic的作用](#panic%E7%9A%84%E4%BD%9C%E7%94%A8)
    - [panic使用场景](#panic%E4%BD%BF%E7%94%A8%E5%9C%BA%E6%99%AF)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## panic源码解读

### panic的作用

- `panic`能够改变程序的控制流，调用`panic`后会立刻停止执行当前函数的剩余代码，并在当前`Goroutine`中递归执行调用方的`defer`；  

- `recover`可以中止`panic`造成的程序崩溃。它是一个只能在`defer`中发挥作用的函数，在其他作用域中调用不会发挥作用；  

举个栗子  

```go
package main

import "fmt"

func main() {
	fmt.Println(1)
	func() {
		fmt.Println(2)
		panic("3")
	}()
	fmt.Println(4)
}
```

输出  

```go
1
2
panic: 3

goroutine 1 [running]:
main.main.func1(...)
        /Users/yj/Go/src/Go-POINT/panic/main.go:9
main.main()
        /Users/yj/Go/src/Go-POINT/panic/main.go:10 +0xee
```

`panic`后会立刻停止执行当前函数的剩余代码，所以4没有打印出来  

**对于recover**

- panic只会触发当前Goroutine的defer；

- recover只有在defer中调用才会生效；

- panic允许在defer中嵌套多次调用；

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(1)

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		fmt.Println(2)
		panic("3")
	}()
	time.Sleep(time.Second)
	fmt.Println(4)
}
```

上面的栗子，因为`recover`和`panic`不在同一个`goroutine`中，所以不会捕获到  

嵌套的demo  

```go
func main() {
	defer fmt.Println("in main")
	defer func() {
		defer func() {
			panic("3 panic again and again")
		}()
		panic("2 panic again")
	}()

	panic("1 panic once")
}
```

输出  

```go
in main
panic: 1 panic once
        panic: 2 panic again
        panic: 3 panic again and again

goroutine 1 [running]:
...
```

多次调用`panic`也不会影响`defer`函数的正常执行，所以使用`defer`进行收尾工作一般来说都是安全的。  

#### panic使用场景

- error：可预见的错误

- panic：不可预见的异常

需要注意的是，你应该尽可能地使用`error`，而不是使用`panic`和`recover`。只有当程序不能继续运行的时候，才应该使用`panic`和`recover`机制。    

`panic`有两个合理的用例。  

1、发生了一个不能恢复的错误，此时程序不能继续运行。 一个例子就是 web 服务器无法绑定所要求的端口。在这种情况下，就应该使用 panic，因为如果不能绑定端口，啥也做不了。  

2、发生了一个编程上的错误。 假如我们有一个接收指针参数的方法，而其他人使用 nil 作为参数调用了它。在这种情况下，我们可以使用panic，因为这是一个编程错误：用 nil 参数调用了一个只能接收合法指针的方法。  

在一般情况下，我们不应通过调用panic函数来报告普通的错误，而应该只把它作为报告致命错误的一种方式。当某些不应该发生的场景发生时，我们就应该调用panic。  

总结下`panic`的使用场景:  

- 1、空指针引用

- 2、下标越界

- 3、除数为0

- 4、不应该出现的分支，比如default

- 5、输入不应该引起函数错误  

### 看下实现



### 参考

【panic 和 recover】https://draveness.me/golang/docs/part2-foundation/ch05-keyword/golang-panic-recover/  