package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type Program struct {
	name     string
	age      int
	language string
}

func Float64bits(f float64) uint64 {
	fmt.Println(reflect.TypeOf(unsafe.Pointer(&f)))            //unsafe.Pointer
	fmt.Println(reflect.TypeOf((*uint64)(unsafe.Pointer(&f)))) //*uint64
	return *(*uint64)(unsafe.Pointer(&f))
}

func main() {
	fmt.Printf("%#016x\n", Float64bits(1.0)) // "0x3ff0000000000000"
}
