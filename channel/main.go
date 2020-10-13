package main

import (
	"fmt"
)

func main() {
	ch := make(chan int, 5)
	ch <- 18
	close(ch)
	x, ok := <-ch
	if ok {
		fmt.Println("received: ", x)
	}

	x = <-ch
	fmt.Println("channel closed, data invalid.")
	// ok为false
	fmt.Println(x)
}
