package main

import (
	"fmt"
	"unsafe"
)

type People struct {
	age  uint8
	name string
}

func main() {
	h := People{
		30,
		"xiaobai",
	}

	i := unsafe.Sizeof(h)
	j := unsafe.Alignof(h)
	k := unsafe.Offsetof(h.name)
	fmt.Println("字节大小：", i)
	fmt.Println("对齐系数：", j)
	fmt.Println("偏移量：", k)
	fmt.Printf("直接获取地址：%p\n", &h)

	var p unsafe.Pointer
	p = unsafe.Pointer(&h)
	fmt.Println("使用unsafe获取地址：", p)
}
