package main

import (
	"fmt"
	"time"
)

func main() {
	// case1()
	case2()
	time.Sleep(time.Millisecond)
}

func case2() {
	go func() {
		fmt.Println("goroutine", 1)
	}()
	go func() {
		fmt.Println("goroutine", 2)
	}()
	go func() {
		fmt.Println("goroutine", 3)
	}()
}

func case1() {
	names := []string{"小白", "小李", "小张"}
	for _, name := range names {
		go func() {
			fmt.Println(name)
		}()
		time.Sleep(time.Second)
	}
}
