package main

import (
	"context"
	"errors"
	"fmt"
	_ "net/http/pprof"
	"time"
)

func main() {
	ch := make(chan int, 1)
	go func() {
		ch <- 1
	}()

	go func() {
		ch <- 2
	}()

	close(ch)

	select {
	case item := <-ch:
		fmt.Println(item)
	}
}

func test1() error {
	return errors.New("just err1")
}

func test2() error {
	return errors.New("just err2")
}

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

func main1() {

	//// 开启pprof，监听请求
	//ip := "0.0.0.0:5555"
	//if err := http.ListenAndServe(ip, nil); err != nil {
	//	fmt.Printf("start pprof failed on %s\n", ip)
	//}

	var cs = []int{2, 3, 5, 7, 8, 9, 3}
	fmt.Println(cs)
	fmt.Println(cap(cs))

	//ch := make(chan int, 3)
	//
	//ch <- 1
	//ch <- 2
	//ch <- 3
	//
	//close(ch)
	//for value := range ch {
	//	fmt.Println("value:", value)
	//}

	//rand.Seed(time.Now().UnixNano())
	//log.SetFlags(0)
	//
	//// ...
	//const MaxRandomNumber = 100000
	//const NumReceivers = 10
	//const NumSenders = 1000
	//
	//wgReceivers := sync.WaitGroup{}
	//wgReceivers.Add(NumReceivers)
	//
	//// ...
	//dataCh := make(chan int, 100)
	//stopCh := make(chan struct{})
	//// stopCh is an additional signal channel.
	//// Its sender is the moderator goroutine shown below.
	//// Its reveivers are all senders and receivers of dataCh.
	//toStop := make(chan string, 1)
	//// the channel toStop is used to notify the moderator
	//// to close the additional signal channel (stopCh).
	//// Its senders are any senders and receivers of dataCh.
	//// Its reveiver is the moderator goroutine shown below.
	//
	//var stoppedBy string
	//
	//// moderator
	//go func() {
	//	stoppedBy = <-toStop // part of the trick used to notify the moderator
	//	// to close the additional signal channel.
	//	close(stopCh)
	//}()
	//
	//// senders
	//for i := 0; i < NumSenders; i++ {
	//	go func(id string) {
	//		for {
	//			value := rand.Intn(MaxRandomNumber)
	//			if value == 0 {
	//				// here, a trick is used to notify the moderator
	//				// to close the additional signal channel.
	//				select {
	//				case toStop <- "sender#" + id:
	//				default:
	//				}
	//				return
	//			}
	//
	//			// the first select here is to try to exit the
	//			// goroutine as early as possible.
	//			select {
	//			case <-stopCh:
	//				return
	//			default:
	//			}
	//
	//			select {
	//			case <-stopCh:
	//				return
	//			case dataCh <- value:
	//			}
	//		}
	//	}(strconv.Itoa(i))
	//}
	//
	//// receivers
	//for i := 0; i < NumReceivers; i++ {
	//	go func(id string) {
	//		defer wgReceivers.Done()
	//
	//		for {
	//			// same as senders, the first select here is to
	//			// try to exit the goroutine as early as possible.
	//			select {
	//			case <-stopCh:
	//				return
	//			default:
	//			}
	//
	//			select {
	//			case <-stopCh:
	//				return
	//			case value := <-dataCh:
	//				if value == MaxRandomNumber-1 {
	//					// the same trick is used to notify the moderator
	//					// to close the additional signal channel.
	//					select {
	//					case toStop <- "receiver#" + id:
	//					default:
	//					}
	//					return
	//				}
	//
	//				log.Println(value)
	//			}
	//		}
	//	}(strconv.Itoa(i))
	//}
	//
	//// ...
	//wgReceivers.Wait()
	//log.Println("stopped by", stoppedBy)
}
