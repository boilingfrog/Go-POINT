package main

import (
	"fmt"
	"sync"
)

func main() {
	waitGroup()
}

func waitGroup() {
	var wg sync.WaitGroup

	wg.Add(3)
	go func() {
		wg.Done()
		fmt.Println(2)
	}()
	go func() {
		wg.Done()
		fmt.Println(3)
	}()
	go func() {
		wg.Done()
		fmt.Println(4)
	}()

	wg.Wait()
	fmt.Println("1 2 3 4 end")
}

func channa() {
	sign := make(chan int, 3)
	go func() {
		sign <- 2
		fmt.Println(2)
	}()
	go func() {
		sign <- 3
		fmt.Println(3)
	}()

	go func() {
		sign <- 4
		fmt.Println(4)
	}()

	for i := 0; i < 3; i++ {
		fmt.Println("执行", <-sign)
	}

}
