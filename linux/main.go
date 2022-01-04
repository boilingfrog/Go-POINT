package main

import "fmt"

func c() (i int) {
	defer func() {
		i++
	}()
	return i
}

func main() {
	fmt.Println(c())
}
