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
	const NumSenders = 1000

	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(1)

	// 发送数据的channel
	dataCh := make(chan int, 100)

	// 无缓冲的channel作为信号量，通知senders的推出
	stopCh := make(chan struct{})

	// 启动个NumSenders个sender
	for i := 0; i < NumSenders; i++ {
		go func() {
			for {
				value := rand.Intn(MaxRandomNumber)

				// 监测到退出信号，马上退出goroutine
				// 否则正常写入dataCh，数据
				select {
				case <-stopCh:
					return
				case dataCh <- value:
				}
			}
		}()
	}

	// 消费者
	go func() {
		defer wgReceivers.Done()

		for value := range dataCh {
			// 某个场景下发出退出的信号量
			if value == MaxRandomNumber-1 {
				close(stopCh)
				return
			}

			log.Println(value)
		}
	}()

	wgReceivers.Wait()
}
