package main

import (
	"fmt"
	"runtime"
)

func main() {
	defer fmt.Println("in main")
	defer func() {
		defer func() {
			panic("3 panic again and again")
		}()
		panic("2 panic again")
	}()
	fmt.Println("++++++++++++++++++++++++++++++")
	runtime.Goexit()
	fmt.Println("main")
	select {}
}
