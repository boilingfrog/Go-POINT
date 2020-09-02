package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan struct{})
	go func() {
		fmt.Println("start working")
		time.Sleep(time.Second * 1)
		ch <- struct{}{}
	}()

	<-ch

	fmt.Println("finished")
}
