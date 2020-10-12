package main

import (
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)

	const MaxRandomNumber = 100000
	const NumReceivers = 10
	const NumSenders = 1000

	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(NumReceivers)

	// 数据的channel
	dataCh := make(chan int, 100)
	// 关闭的channel的信号
	stopCh := make(chan struct{})
	// toStop通知关闭stopCh，同时作为receiver退出的信息
	toStop := make(chan string, 1)

	var stoppedBy string

	// 当收到toStop的信号，关闭stopCh
	go func() {
		stoppedBy = <-toStop
		close(stopCh)
	}()

	// 发送端
	for i := 0; i < NumSenders; i++ {
		go func(id string) {
			for {
				value := rand.Intn(MaxRandomNumber)
				// 满足条件发出关闭的请求到toStop
				if value == 0 {
					select {
					case toStop <- "sender#" + id:
					default:
					}
					return
				}

				select {
				// 检测的关闭的stopCh，退出发送者
				case <-stopCh:
					return
				case dataCh <- value:
				}
			}
		}(strconv.Itoa(i))
	}

	// 接收端
	for i := 0; i < NumReceivers; i++ {
		go func(id string) {
			defer wgReceivers.Done()

			for {
				select {
				// 检测的关闭的stopCh，退出接收者
				case <-stopCh:
					return
				case value := <-dataCh:
					// 满足条件发出关闭的请求到toStop
					if value == MaxRandomNumber-1 {
						select {
						case toStop <- "receiver#" + id:
						default:
						}
						return
					}

					log.Println(value)
				}
			}
		}(strconv.Itoa(i))
	}

	wgReceivers.Wait()
	log.Println("stopped by", stoppedBy)
}
