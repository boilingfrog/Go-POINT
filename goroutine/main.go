package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	// 开启pprof，监听请求
	ip := "127.0.0.1:6069"
	if err := http.ListenAndServe(ip, nil); err != nil {
		fmt.Printf("start pprof failed on %s\n", ip)
	}
}

//func query() int {
//	n := rand.Intn(100)
//	time.Sleep(time.Duration(n) * time.Millisecond)
//	return n
//}
//
//func TimeoutCancelChannel() {
//	done := make(chan struct{}, 1)
//
//	go func() {
//		// 执行业务逻辑
//		done <- struct{}{}
//	}()
//
//	select {
//	case <-done:
//		fmt.Println("call successfully!!!")
//		return
//	case <-time.After(time.Duration(800 * time.Millisecond)):
//		fmt.Println("timeout!!!")
//		// 使用独立的协程处理超时，需求添加return退出协程，否则会导致当前协程被通知channel阻塞，进而导致内存泄露
//		return
//	}
//}
//
//func TimeoutCancelContext() {
//	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*800))
//	go func() {
//		// 具体的业务逻辑
//		// 取消超时
//		defer cancel()
//	}()
//
//	select {
//	case <-ctx.Done():
//		fmt.Println("call successfully!!!")
//		return
//	}
//}
//
//func channelNoReceiver() {
//	ch := make(chan int)
//	go func() {
//		ch <- 1
//		fmt.Println(111)
//	}()
//}
//
//func channelNoProducer() {
//	ch := make(chan int, 1)
//	go func() {
//		<-ch
//		fmt.Println(111)
//	}()
//}
//
//func case6() {
//	names := []string{"小白", "小明", "小红", "小张"}
//	for _, name := range names {
//		go func(who string) {
//			fmt.Println("名字", who)
//		}(name)
//	}
//	time.Sleep(time.Millisecond)
//}
//
//func case5() {
//	names := []string{"小白", "小明", "小红", "小张"}
//	for _, name := range names {
//		go func() {
//			fmt.Println("名字", name)
//		}()
//	}
//	time.Sleep(time.Millisecond)
//}
//
//func case4() {
//	name := "小白"
//	go func() {
//		fmt.Println(name)
//	}()
//	time.Sleep(time.Millisecond)
//	name = "小李"
//}
//
//func case3() {
//	name := "小白"
//	go func() {
//		fmt.Println(name)
//	}()
//	name = "小李"
//	time.Sleep(time.Millisecond)
//
//}
//
//func case2() {
//	go func() {
//		fmt.Println("goroutine", 1)
//	}()
//	go func() {
//		fmt.Println("goroutine", 2)
//	}()
//	go func() {
//		fmt.Println("goroutine", 3)
//	}()
//}
//
//func case1() {
//	names := []string{"小白", "小李", "小张"}
//	for _, name := range names {
//		go func() {
//			fmt.Println(name)
//		}()
//		time.Sleep(time.Second)
//	}
//}
