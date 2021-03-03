package main

import (
	"sync/atomic"
	"unsafe"
)

var value int32 = 1

func main() {
	var p unsafe.Pointer
	newP := 42
	atomic.CompareAndSwapPointer(&p, nil, unsafe.Pointer(&newP))

	v := (*int)(p)
	println(*v)
}
