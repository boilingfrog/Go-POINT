package main

import (
	"fmt"
	"sync"
)

// 全局变量
var counter int

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(n int) {
			wg.Done()
			fmt.Println(n)
		}(i)
	}
	wg.Wait()
}
