package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

var count = 500
var a = make(map[int]struct{})

func main() {
	v := struct{}{}

	for i := 0; i < 10000; i++ {
		a[i] = v
	}

	runtime.GC()
	printMemStats("添加1万个键值对后")
	fmt.Println("删除前Map长度：", len(a))

	for i := 0; i < 2; i++ {
		delete(a, i)
	}
	fmt.Println("删除后Map长度：", len(a))

	// 再次进行手动GC回收
	runtime.GC()
	printMemStats("删除1万个键值对后")

	// 设置为nil进行回收
	a = nil
	runtime.GC()
	printMemStats("设置为nil后")
}

func printMemStats(mag string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%v：分配的内存 = %vKB, GC的次数 = %v\n", mag, m.Alloc/1024, m.NumGC)
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

func test2() {
	var mapp = make(map[int]string, 100)
	for i := 0; i < 1000; i++ {
		go func(i int) {
			mapp[i] = "xxx"
			fmt.Println(i)
		}(i)
	}

}
