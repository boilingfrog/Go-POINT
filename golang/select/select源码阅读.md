## 深入了解下 go 中的 select

### 前言

这里借助于几个经常遇到的 select 的使用 demo 来作为开始，先来看看，下面几个 demo 的输出情况  

1、栗子一  

```go
func main() {
	chan1 := make(chan int)
	chan2 := make(chan int)

	go func() {
		chan1 <- 1
	}()

	go func() {
		chan2 <- 1
	}()

	select {
	case <-chan1:
		fmt.Println("chan1 ready.")
	case <-chan2:
		fmt.Println("chan2 ready.")
	default:
		fmt.Println("default")
	}
}
```

select 中的 case 执行是随机的，所以当 case 监听的 channel 有数据传入，就执行相应的流程并退出 select，如果对应的 case 没有收到 channel 的数据，就执行 default 语句，然后退出 select。  

上面的协程启动时间是无法预估的，所以上面的两个 case 和 default ，都有机会执行。  

可能的输出  

可能输出1、
```
chan1 ready.
```

可能输出2、
```
chan2 ready.
```

可能输出3、
```
default
```

2、栗子二

```go
func main() {
	chan1 := make(chan int)
	chan2 := make(chan int)

	go func() {
		close(chan1)
	}()

	go func() {
		close(chan2)
	}()

	select {
	case <-chan1:
		fmt.Println("chan1 ready.")
	case <-chan2:
		fmt.Println("chan2 ready.")
	default:
		fmt.Println("default")
	}
}
```

已经关闭的 channel ，使用 select 是可以从中读出对应的零值，同时两面关闭 channel 的协程的执行实际也是不可控的，原则上，上面两个 case 和 default 都有可能被执行。   

可能的输出  

可能输出1、
```
chan1 ready.
```

可能输出2、
```
chan2 ready.
```

可能输出3、
```
default
```