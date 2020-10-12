package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)

	// ...
	const MaxRandomNumber = 100000
	const NumReceivers = 100

	// 使用WaitGroup来阻塞查看打印的效果
	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(NumReceivers)

	// 设置channel的长度为10
	dataCh := make(chan int, 100)

	// the sender
	go func() {
		for {
			if value := rand.Intn(MaxRandomNumber); value == 0 {
				// 需要关闭的时候直接关闭就好了，是很安全的
				close(dataCh)
				return
			} else {
				dataCh <- value
			}
		}
	}()

	// receivers
	for i := 0; i < NumReceivers; i++ {
		go func() {
			defer wgReceivers.Done()

			// 监听dataCh，接收里面的值
			for value := range dataCh {
				log.Println(value)
			}
		}()
	}

	wgReceivers.Wait()

}
