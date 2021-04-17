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
func pool1() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer func() {
			fmt.Println("7777777777")
			wg.Done()
		}()
		time.Sleep(time.Second)

		mockSqlPool()
		panic("++++")
	}()

	go func() {
		defer wg.Done()
		mockSqlPool()

		time.Sleep(time.Second * 2)

		fmt.Println("run 1")
	}()

	fmt.Println("执行了吗")
	wg.Wait()
	fmt.Println("完美退出了")
}
