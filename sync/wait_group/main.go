package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"
)

func main() {
	fdg := [2]uint32{}

	state := atomic.AddUint64((*uint64)(unsafe.Pointer(&fdg)), uint64(1)<<32)
	fmt.Println("wait:", fdg[0])
	fmt.Println("count:", fdg[1])
	count := int32(state >> 32)
	wait := uint32(state)
	fmt.Println("count:", count)
	fmt.Println("wait:", wait)
}

func waitGroup() {
	var wg sync.WaitGroup

	wg.Add(4)
	go func() {
		defer wg.Done()
		fmt.Println(1)
	}()

	go func() {
		defer wg.Done()
		fmt.Println(2)
	}()

	go func() {
		defer wg.Done()
		fmt.Println(3)
	}()

	go func() {
		defer wg.Done()
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
