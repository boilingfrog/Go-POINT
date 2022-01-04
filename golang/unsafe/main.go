package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type People struct {
	age  uint8
	name string
}

func Float64bits(f float64) uint64 {
	fmt.Println(reflect.TypeOf(unsafe.Pointer(&f)))            //unsafe.Pointer
	fmt.Println(reflect.TypeOf((*uint64)(unsafe.Pointer(&f)))) //*uint64
	return *(*uint64)(unsafe.Pointer(&f))
}

func main() {
	fmt.Printf("%#016x\n", float64(1.0)) // "0x3ff0000000000000"

	fmt.Printf("%#016x\n", Float64bits(1.0)) // "0x3ff0000000000000"
}
