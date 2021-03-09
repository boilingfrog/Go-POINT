package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	locker = new(sync.Mutex)
	cond   = sync.NewCond(locker)
)

func listen(x int) {
	// 获取锁
	cond.L.Lock()
	// 等待通知  暂时阻塞
	cond.Wait()
	fmt.Println(x)
	// 释放锁
	cond.L.Unlock()
}

func main() {
	// 启动40个阻塞的县城
	for i := 1; i <= 40; i++ {
		go listen(i)
	}

	fmt.Println("start all")

	// 3秒之后 下发一个通知给已经获取锁的goroutine	time.Sleep(time.Second * 3)
	fmt.Println("++++++++++++++++++++one Signal")
	cond.Signal()

	// 3秒之后 下发一个通知给已经获取锁的goroutine
	time.Sleep(time.Second * 3)
	fmt.Println("++++++++++++++++++++one Signal")
	cond.Signal()

	// 3秒之后 下发广播给所有等待的goroutine
	time.Sleep(time.Second * 3)
	fmt.Println("++++++++++++++++++++begin broadcast")
	cond.Broadcast()
	// 阻塞直到所有的全部输出
	time.Sleep(time.Second * 60)
}
