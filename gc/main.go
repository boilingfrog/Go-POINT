package main

import (
	"fmt"
	"sync"
	"unsafe"
)

type s struct {
	Bool  bool
	Int64 int64
}

func main() {
	fmt.Printf("bool size: %d\n", unsafe.Sizeof(bool(true)))
	fmt.Printf("int32 size: %d\n", unsafe.Sizeof(int32(0)))
	fmt.Printf("int8 size: %d\n", unsafe.Sizeof(int8(0)))
	fmt.Printf("int64 size: %d\n", unsafe.Sizeof(int64(0)))
	fmt.Printf("byte size: %d\n", unsafe.Sizeof(byte(0)))
	fmt.Printf("string size: %d\n", unsafe.Sizeof("E"))

	part1 := s{}
	fmt.Printf("part1 size: %d, align: %d\n", unsafe.Sizeof(part1), unsafe.Alignof(part1))
}

func test1() {
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(wg *sync.WaitGroup) {
			var counter int
			for i := 0; i < 1e10; i++ {
				counter++
			}
			wg.Done()
		}(&wg)
	}

	wg.Wait()
}
