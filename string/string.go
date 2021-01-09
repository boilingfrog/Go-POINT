package main

import (
	"fmt"
	"sync"
)

func main() {
	var (
		jobsChan = make(chan int, 10)
		rw       sync.RWMutex
		num      int
	)

	wg := sync.WaitGroup{}
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()
			for redPacket := range jobsChan {
				_ = redPacket
				rw.Lock()
				num++
				rw.Unlock()
			}
		}()
	}

	for i := 0; i < 10; i++ {
		jobsChan <- i + 1
	}
	close(jobsChan)
	wg.Wait()
	fmt.Println(num)
}
