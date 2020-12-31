package main

import (
	"fmt"
)

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

type Person struct {
	noCopy noCopy
	name   string
}

// go中的函数传参都是值拷贝
func test(person Person) {
	fmt.Println(person)
}

func main() {
	var person Person
	test(person)
}
