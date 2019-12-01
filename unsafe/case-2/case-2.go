package main

import "fmt"

func double1(x *int) {
	*x += *x
	x = nil
}

func main() {
	var a = 3
	double1(&a)

	fmt.Println(a)

	p := &a

	double1(p)

	fmt.Println(*p)

}
