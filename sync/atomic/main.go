package main

import (
	"fmt"
	"sync/atomic"
)

var value int32 = 1

func main() {
	var a, b int32 = 13, 13
	var c int32
	res := atomic.CompareAndSwapInt32(&a, b, c)
	fmt.Println(res)
	fmt.Println(c)
}

//不断地尝试原子地更新value的值,直到操作成功为止
func addValue(delta int32) {
	for {
		v := value
		if atomic.CompareAndSwapInt32(&value, v, v+delta) {
			break
		}
	}
}
