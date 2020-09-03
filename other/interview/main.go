package main

import (
	"context"
	"fmt"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 避免其他地方忘记 cancel，且重复调用不影响

	for n := range gen(ctx) {
		fmt.Println(n)
		if n == 5 {
			cancel()
			break
		}
	}
	// ……
}
