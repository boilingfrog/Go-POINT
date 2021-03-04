package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"
)

func main() {
	// waitGroup()

	fdg := [2]uint32{}

	fmt.Println((*uint32)(unsafe.Pointer(&fdg)))
	fmt.Println((*uint32)(unsafe.Pointer(&fdg[0])))
	//
	//fmt.Println((*uint32)(unsafe.Pointer(&fdg[1])))
	//fmt.Println((*uint32)(unsafe.Pointer(&fdg[2])))
	// fmt.Printf("uint64 size: %d\n", unsafe.Sizeof(uint64(0)))

	state := atomic.AddUint64((*uint64)(unsafe.Pointer(&fdg)), uint64(10)<<32)
	fmt.Println(state)
	v := int32(state >> 32)
	w := uint32(state)
	fmt.Println(v)
	fmt.Println(w)
	fmt.Println(fmt.Printf("%b\n", 1)) // 1110
	fmt.Printf("%b\n", v)              // 1110

	//state = atomic.AddUint64((*uint64)(unsafe.Pointer(&fdg)), uint64(int(-1))<<32)
	//
	//fmt.Println(state)
	//v = int32(state >> 32)
	//w = uint32(state)
	//fmt.Println(v)
	//fmt.Println(w)

	//	fmt.Println(uint64(-1))
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
