package main

import (
	_ "net/http/pprof"
)

var (
	poolCount      = 5
	goroutineCount = 10
)

func main() {
	//jobsChan := make(chan int, poolCount)
	//
	//// workers
	//var wg sync.WaitGroup
	//for i := 0; i < goroutineCount; i++ {
	//	wg.Add(1)
	//	go func() {
	//		defer wg.Done()
	//		for item := range jobsChan {
	//			// ...
	//			fmt.Println(item)
	//		}
	//	}()
	//}
	//
	//// senders
	//for i := 0; i < 1000; i++ {
	//	jobsChan <- i
	//}
	//
	//// 关闭channel，上游的goroutine在读完channel的内容，就会通过wg的done退出
	//close(jobsChan)
	//wg.Wait()

}

var base string

//func main() {
//	// 开启pprof，监听请求
//	ip := "127.0.0.1:6069"
//	// 开启pprof
//	go func() {
//		if err := http.ListenAndServe(ip, nil); err != nil {
//			fmt.Printf("start pprof failed on %s\n", ip)
//			os.Exit(1)
//		}
//	}()
//
//	// 路由，访问，触发内存泄露的代码判断
//	http.HandleFunc("/test", handler)
//
//	// 阻塞
//	c := make(chan os.Signal, 1)
//	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
//	for {
//		s := <-c
//		switch s {
//		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
//			time.Sleep(time.Second)
//			return
//		case syscall.SIGHUP:
//		default:
//			return
//		}
//	}
//}
//
//func handler(w http.ResponseWriter, r *http.Request) {
//	var (
//		jobsChan = make(chan int, 10)
//		rw       sync.RWMutex
//		num      int
//	)
//
//	wg := sync.WaitGroup{}
//	wg.Add(3)
//	for i := 0; i < 3; i++ {
//		go func() {
//			defer wg.Done()
//			for redPacket := range jobsChan {
//				_ = redPacket
//				rw.Lock()
//				num++
//				rw.Unlock()
//			}
//		}()
//	}
//
//	for i := 0; i < 10; i++ {
//		jobsChan <- i + 1
//	}
//	close(jobsChan)
//	wg.Wait()
//	fmt.Println(num)
//}

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
