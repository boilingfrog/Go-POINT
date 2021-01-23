package main

import "fmt"

var s0 string // 一个包级变量

func main() {
	s := "212"
	fmt.Println(&s)
	fmt.Println(len(s))
	s = "33hhhhhhhhhnihihfn你皇帝是滴是滴四十多"
	fmt.Println(&s)
	fmt.Println(len(s))

}

func f(s1 string) {
	s0 = s1[:50]
	// 目前，s0和s1共享着承载它们的字节序列的同一个内存块。
	// 虽然s1到这里已经不再被使用了，但是s0仍然在使用中，
	// 所以它们共享的内存块将不会被回收。虽然此内存块中
	// 只有50字节被真正使用，而其它字节却无法再被使用。
}

func bigString() string {
	var buf []byte
	for i := 0; i < 10; i++ {
		buf = append(buf, make([]byte, 1024*1024)...)
	}
	return string(buf)

}
