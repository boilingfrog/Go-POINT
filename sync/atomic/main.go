package main

import (
	"fmt"
	"sync/atomic"
)

var value int32 = 1

func main() {
	var a, b int32 = 13, 13
	var c int32 = 9
	res := atomic.CompareAndSwapInt32(&a, b, c)
	fmt.Println("swapped:", res)
	fmt.Println("替换的值:", c)
	fmt.Println("替换之后a的值:", a)
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
