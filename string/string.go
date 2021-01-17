package main

import "fmt"

func main() {
	s := "xiaobai"

	fmt.Println([]byte(s))

	s1 := "哈哈"
	fmt.Println([]byte(s1))

	fmt.Printf("%#v\n", []byte(s1))

}
