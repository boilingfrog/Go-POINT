package main

import "fmt"

func main() {
	ch := make(chan interface{}, 2)
	x := 1
	ch <- x
	select {
	case ch <- x:
		println("send success") // 如果初始化为有缓存 channel，则会发送成功
	default:
		fmt.Println("111")
	}
	return
}
