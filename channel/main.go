package main

import (
	"fmt"
	"unsafe"
)

func main() {

	fmt.Println(121)

	ch := make(chan int, 1)
	_ = ch
}