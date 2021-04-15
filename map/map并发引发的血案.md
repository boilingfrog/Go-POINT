## 错误使用map引发的血案

### 前言  

最近业务中，同事使用`map`来接收返回的结果，使用`waitGroup`来并发的处理执行返回的结果，结果上线之后，直接崩了。  

<img src="/img/map_1.jpg"  alt="map" align=center />

日志大量的数据库缓存池连接失败

```go
{"ecode":-500,"message":"timed out while checking out a connection from connection pool"}

{"ecode":-500,"message":"connection(xxxxxxxxx:xxxxx) failed to write: context deadline exceeded"}
```

### 原因

先来看来伪代码

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

var count = 300

func main() {
	var data = make(map[int]string, count)
	var wg sync.WaitGroup

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			time.Sleep(time.Second * 1)
			mockSqlPool()
			data[i] = "test"
		}(i)
	}
	fmt.Println("-----------WaitGroup执行结束了-----------")
	wg.Wait()
}

// 模拟数据库的连接和释放
func mockSqlPool() {
	defer fmt.Println("关闭pool")
	fmt.Println("我是pool")
}
```

一个全局的`map`，然后`WaitGroup`开启一组协程并发的读写数据，写入内容到`map`中。  

上面的错误原因有两个：  

- 1、`map`不是并发安全，并发写的时候会触发`panic`；  

- 2、避免在循环中连接梳理数据库； 

因为map的并发写入，导致触发了`panic`，对于`panic`,当`panic`异常发生时，程序会中断后面代码的执行，
