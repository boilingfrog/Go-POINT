package main

import (
	"fmt"
	"sync"
)

type s struct {
	Bool  bool
	Byte  byte
	In32  int32
	Int64 int64
	//Int8  int8 //Int64 int64
	//
	//Byte byte
	//In32 int32
	//Int8 int8
}

//func main() {
//	fmt.Printf("bool size: %d\n", unsafe.Sizeof(bool(true)))
//	fmt.Printf("int32 size: %d\n", unsafe.Sizeof(int32(0)))
//	fmt.Printf("int8 size: %d\n", unsafe.Sizeof(int8(0)))
//	fmt.Printf("int64 size: %d\n", unsafe.Sizeof(int64(0)))
//	fmt.Printf("byte size: %d\n", unsafe.Sizeof(byte(0)))
//	fmt.Printf("string size: %d\n", unsafe.Sizeof("E"))
//
//	part1 := s{}
//	fmt.Printf("part1 size: %d, align: %d\n", unsafe.Sizeof(part1), unsafe.Alignof(part1))
//}

func main() {
	fmt.Println(12)
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
