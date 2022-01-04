package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"golang.org/x/sync/semaphore"
)

func main() {
	fdg := [3]uint32{
		2, 2, 1,
	}

	state := atomic.AddUint64((*uint64)(unsafe.Pointer(&fdg)), uint64(-1)<<32)
	fmt.Println("wait:", fdg[0])
	fmt.Println("count:", fdg[1])
	count := int32(state >> 32)
	wait := uint32(state)
	fmt.Println("count:", count)
	fmt.Println("wait:", wait)

	//semaphoreTest()

}

func doSomething(u string) { // 模拟抓取任务的执行
	fmt.Println(u)
	time.Sleep(2 * time.Second)
}

const (
	Limit  = 3 // 同時并行运行的goroutine上限
	Weight = 1 // 每个goroutine获取信号量资源的权重
)

func semaphoreTest() {
	urls := []string{
		"http://www.example.com",
		"http://www.example.net",
		"http://www.example.net/foo",
		"http://www.example.net/bar",
		"http://www.example.net/baz",
	}
	s := semaphore.NewWeighted(Limit)
	var w sync.WaitGroup
	for _, u := range urls {
		w.Add(1)
		go func(u string) {
			s.Acquire(context.Background(), Weight)
			doSomething(u)
			s.Release(Weight)
			w.Done()
		}(u)
	}
	w.Wait()

	fmt.Println("All Done")
}

func waitGroup() {
	var wg sync.WaitGroup

	wg.Add(4)
	go func() {
		defer wg.Done()
		fmt.Println(1)
	}()

	go func() {
		defer wg.Done()
		fmt.Println(2)
	}()

	go func() {
		defer wg.Done()
		fmt.Println(3)
	}()

	go func() {
		defer wg.Done()
		fmt.Println(4)
	}()

	wg.Wait()
	fmt.Println("1 2 3 4 end")
}

func channa() {
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

}
