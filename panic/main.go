package main

import (
	"fmt"
)

func main() {
	defer fmt.Println("in main")
	defer func() {
		defer func() {
			panic("3 panic again and again")
		}()
		panic("2 panic again")
	}()

	panic("1 panic once")
}
