package main

import (
	"fmt"
	"sync"
)

func main() {
	ch2 := make(chan int, 1)
	ch3 := make(chan int, 1)
	ch4 := make(chan int, 1)

	go func() {
		for i := 1; i < 100; i++ {
			ch4 <- 1
			if i%3 == 0 {
				ch2 <- i
			} else {
				ch3 <- i
			}
		}
		close(ch2)
		close(ch3)
		close(ch4)
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for item := range ch2 {
			fmt.Println(item)

			<-ch4
		}

	}()

	go func() {
		defer wg.Done()
		for item := range ch3 {
			fmt.Println(item)
			<-ch4
		}
	}()

	wg.Wait()

}
