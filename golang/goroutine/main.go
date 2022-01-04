package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

//
//func main() {
//	pool2()
//}

//func pool2() {
//	p := pool.NewLimited(10)
//	defer p.Close()
//
//	batch := p.Batch()
//
//	// for max speed Queue in another goroutine
//	// but it is not required, just can't start reading results
//	// until all items are Queued.
//
//	go func() {
//		for i := 0; i < 10; i++ {
//			batch.Queue(sendEmail("email content"))
//		}
//
//		// DO NOT FORGET THIS OR GOROUTINES WILL DEADLOCK
//		// if calling Cancel() it calles QueueComplete() internally
//		batch.QueueComplete()
//	}()
//
//	for email := range batch.Results() {
//
//		if err := email.Error(); err != nil {
//			// handle error
//			// maybe call batch.Cancel()
//		}
//
//		// use return value
//		fmt.Println(email.Value().(bool))
//	}
//}
//
//func sendEmail(email string) pool.WorkFunc {
//	return func(wu pool.WorkUnit) (interface{}, error) {
//
//		// simulate waiting for something, like TCP connection to be established
//		// or connection from pool grabbed
//		time.Sleep(time.Second * 1)
//
//		if wu.IsCancelled() {
//			// return values not used
//			return nil, nil
//		}
//
//		// ready for processing...
//
//		return true, nil // everything ok, send nil, error if not
//	}
//}
//
//func pool1() {
//
//	p := pool.NewLimited(10)
//	defer p.Close()
//
//	user := p.Queue(getUser(13))
//	other := p.Queue(getOtherInfo(13))
//
//	user.Wait()
//	if err := user.Error(); err != nil {
//		// handle error
//	}
//
//	// do stuff with user
//	username := user.Value().(string)
//	fmt.Println(username)
//
//	other.Wait()
//	if err := other.Error(); err != nil {
//		// handle error
//	}
//
//	// do stuff with other
//	otherInfo := other.Value().(string)
//	fmt.Println(otherInfo)
//}
//
//func getUser(id int) pool.WorkFunc {
//
//	return func(wu pool.WorkUnit) (interface{}, error) {
//
//		// simulate waiting for something, like TCP connection to be established
//		// or connection from pool grabbed
//		time.Sleep(time.Second * 1)
//
//		if wu.IsCancelled() {
//			// return values not used
//			return nil, nil
//		}
//		// ready for processing...
//		return "Joeybloggs", nil
//	}
//}
//
//func getOtherInfo(id int) pool.WorkFunc {
//
//	return func(wu pool.WorkUnit) (interface{}, error) {
//		// simulate waiting for something, like TCP connection to be established
//		// or connection from pool grabbed
//		time.Sleep(time.Second * 1)
//
//		if wu.IsCancelled() {
//			// return values not used
//			return nil, nil
//		}
//
//		// ready for processing...
//		return "Other Info", nil
//	}
//}
//
//func demoFunc() {
//	time.Sleep(10 * time.Millisecond)
//	fmt.Println("Hello World!")
//}
//
//var sum int32
//
//func dd() {
//	defer ants.Release()
//
//	runTimes := 1000
//
//	var wg sync.WaitGroup
//	syncCalculateSum := func() {
//		demoFunc()
//		wg.Done()
//	}
//	for i := 0; i < runTimes; i++ {
//		wg.Add(1)
//		_ = ants.Submit(syncCalculateSum)
//	}
//	wg.Wait()
//	fmt.Printf("running goroutines: %d\n", ants.Running())
//	fmt.Printf("finish all tasks.\n")
//}

//func main() {
//	defer ants.Release()
//
//	runTimes := 3000
//
//	// Use the common pool.
//	var wg sync.WaitGroup
//	syncCalculateSum := func() {
//		demoFunc()
//		wg.Done()
//	}
//	for i := 0; i < runTimes; i++ {
//		wg.Add(1)
//		_ = ants.Submit(syncCalculateSum)
//	}
//	wg.Wait()
//	fmt.Printf("running goroutines: %d\n", ants.Running())
//	fmt.Printf("finish all tasks.\n")
//
//	//Use the pool with a function,
//	//set 10 to the capacity of goroutine pool and 1 second for expired duration.duration
//	p, _ := ants.NewPoolWithFunc(10, func(i interface{}) {
//		myFunc(i)
//		wg.Done()
//	})
//	defer p.Release()
//	// Submit tasks one by one.
//	for i := 0; i < runTimes; i++ {
//		wg.Add(1)
//		_ = p.Invoke(int32(i))
//	}
//	wg.Wait()
//	fmt.Printf("running goroutines: %d\n", p.Running())
//	fmt.Printf("finish all tasks, result is %d\n", sum)
//}

type poolLimit struct {
	poolCount      int
	goroutineCount int
}

func newPool(poolCount, goroutineCount int) *poolLimit {
	return &poolLimit{
		poolCount:      poolCount,
		goroutineCount: goroutineCount,
	}
}

func (pool *poolLimit) Go(f func() error) {
	go func() {
		if err := f(); err != nil {

		}
	}()

}

var (
	// channel长度
	poolCount = 5
	// 复用的goroutine数量
	goroutineCount = 10
)

func limit() {
	jobsChan := make(chan int, poolCount)

	// workers
	var wg sync.WaitGroup
	for i := 0; i < goroutineCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range jobsChan {
				// ...
				fmt.Println(item)
			}
		}()
	}

	// senders
	for i := 0; i < 1000; i++ {
		jobsChan <- i
	}

	// 关闭channel，上游的goroutine在读完channel的内容，就会通过wg的done退出
	close(jobsChan)
	wg.Wait()
}

var base string

func main() {
	// 开启pprof，监听请求
	ip := "127.0.0.1:6069"
	// 开启pprof
	go func() {
		if err := http.ListenAndServe(ip, nil); err != nil {
			fmt.Printf("start pprof failed on %s\n", ip)
			os.Exit(1)
		}
	}()

	// 路由，访问，触发内存泄露的代码判断
	http.HandleFunc("/test", handler)

	// 阻塞
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	var count = 1000

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
}

// 模拟数据库的连接和释放
func mockSqlPool() {
	defer fmt.Println("关闭pool")
	fmt.Println("我是pool")
}

func query() int {
	n := rand.Intn(100)
	time.Sleep(time.Duration(n) * time.Millisecond)
	return n
}

func TimeoutCancelChannel() {
	done := make(chan struct{}, 1)

	go func() {
		// 执行业务逻辑
		done <- struct{}{}
	}()

	select {
	case <-done:
		fmt.Println("call successfully!!!")
		return
	case <-time.After(time.Duration(800 * time.Millisecond)):
		fmt.Println("timeout!!!")
		// 使用独立的协程处理超时，需求添加return退出协程，否则会导致当前协程被通知channel阻塞，进而导致内存泄露
		return
	}
}

func TimeoutCancelContext() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*800))
	go func() {
		// 具体的业务逻辑
		// 取消超时
		defer cancel()
	}()

	select {
	case <-ctx.Done():
		fmt.Println("call successfully!!!")
		return
	}
}

func channelNoReceiver() {
	ch := make(chan int)
	go func() {
		ch <- 1
		fmt.Println(111)
	}()
}

func channelNoProducer() {
	ch := make(chan int, 1)
	go func() {
		<-ch
		fmt.Println(111)
	}()
}

func case6() {
	names := []string{"小白", "小明", "小红", "小张"}
	for _, name := range names {
		go func(who string) {
			fmt.Println("名字", who)
		}(name)
	}
	time.Sleep(time.Millisecond)
}

func case5() {
	names := []string{"小白", "小明", "小红", "小张"}
	for _, name := range names {
		go func() {
			fmt.Println("名字", name)
		}()
	}
	time.Sleep(time.Millisecond)
}

func case4() {
	name := "小白"
	go func() {
		fmt.Println(name)
	}()
	time.Sleep(time.Millisecond)
	name = "小李"
}

func case3() {
	name := "小白"
	go func() {
		fmt.Println(name)
	}()
	name = "小李"
	time.Sleep(time.Millisecond)

}

func case2() {
	go func() {
		fmt.Println("goroutine", 1)
	}()
	go func() {
		fmt.Println("goroutine", 2)
	}()
	go func() {
		fmt.Println("goroutine", 3)
	}()
}

func case1() {
	names := []string{"小白", "小李", "小张"}
	for _, name := range names {
		go func() {
			fmt.Println(name)
		}()
		time.Sleep(time.Second)
	}
}
