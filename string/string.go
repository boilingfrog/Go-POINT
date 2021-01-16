package main

import (
	"fmt"
)

func main() {
	var a byte = 65
	// 8进制写法: var c byte = '\101'     其中 \ 是固定前缀
	// 16进制写法: var c byte = '\x41'    其中 \x 是固定前缀

	var b uint8 = 66
	fmt.Printf("a 的值: %c \nb 的值: %c", a, b)

	// 或者使用 string 函数
	// fmt.Println("a 的值: ", string(a)," \nb 的值: ", string(b))
}
