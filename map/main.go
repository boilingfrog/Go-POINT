package main

import (
	"fmt"
	"sync"
	"time"
)

var count = 500

func main() {
	var data = make(map[int]string, count)
	_ = data
	var wg sync.WaitGroup

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			time.Sleep(time.Second * 1)
			mockSqlPool()
			// panic("+++")
			data[i] = "test"
		}(i)
	}
	fmt.Println("-----------WaitGroup执行-----------")
	wg.Wait()
	fmt.Println("哈哈，你好吗")
	//fmt.Println(data)
}

// 模拟数据库的连接和释放
func mockSqlPool() {
	defer fmt.Println("关闭pool")
	fmt.Println("我是pool")
}
