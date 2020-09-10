package main

import (
	"fmt"
	"sync/atomic"
)

var value int32 = 1

func main() {
	fmt.Println("======old value=======")
	fmt.Println(value)
	addValue(10)
	fmt.Println("======New value=======")
	fmt.Println(value)

}

//不断地尝试原子地更新value的值,直到操作成功为止
func addValue(delta int32) {
	for {
		v := value
		if atomic.CompareAndSwapInt32(&value, v, (v + delta)) {
			break
		}
	}
}
