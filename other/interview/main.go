package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

//func main() {
//	ch := make(chan struct{})
//	go func() {
//		fmt.Println("start working")
//		time.Sleep(time.Second * 1)
//		ch <- struct{}{}
//	}()
//
//	<-ch
//
//	fmt.Println("finished")
//}

func gen(ctx context.Context) <-chan int {
	ch := make(chan int)
	go func() {
		var n int
		for {
			select {
			case <-ctx.Done():
				return
			case ch <- n:
				n++
				time.Sleep(time.Second)
			}
		}
	}()
	return ch
}

func main() {
	jobsCount := 100
	group := sync.WaitGroup{}
	var jobsChan = make(chan int, 3)
	// a) 生成指定数目的 goroutine，每个 goroutine 消费 jobsChan 中的数据
	poolCount := 3
	for i := 1; i < poolCount; i++ {
		go func() {
			for j := range jobsChan {
				fmt.Printf("hello %d\n", j)
				time.Sleep(time.Second)
				group.Done()
			}
		}()
	}
	// b) 把 job 依次推送到 jobsChan 供 goroutine 消费
	for i := 0; i < jobsCount; i++ {
		jobsChan <- i
		group.Add(1)
		fmt.Printf("index: %d,goroutine Num: %d\n", i, runtime.NumGoroutine())
	}
	group.Wait()
	fmt.Println("done!")
}
