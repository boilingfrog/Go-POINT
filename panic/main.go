package main

import "fmt"

func main() {
	fmt.Println(1)
	func() {
		fmt.Println(2)
		panic("3")
	}()
	fmt.Println(4)
}
