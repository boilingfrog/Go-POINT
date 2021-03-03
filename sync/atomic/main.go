package main

import (
	"fmt"
	"sync/atomic"
)

var value int32 = 1

func main() {
	var p int32 = 1
	res := atomic.LoadInt32(&p)

	fmt.Println(res)
}
