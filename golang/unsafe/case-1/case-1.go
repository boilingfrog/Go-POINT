package main

import "fmt"

func double(x int) {
	fmt.Println(x)
	x += x
	fmt.Println(x)
}

func main() {
	var a = 3
	double(a)
	fmt.Println(a)

}
