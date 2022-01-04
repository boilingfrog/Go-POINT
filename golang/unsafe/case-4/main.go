package main

import (
	"fmt"
)

type Program struct {
	name     string
	age      int
	language string
}

func main() {
	p := Program{
		name:     "小明",
		age:      18,
		language: "golang",
	}
	fmt.Println(p)
}
