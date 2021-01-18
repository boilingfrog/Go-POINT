package main

import "fmt"

func main() {
	s := "A"
	fmt.Print("打印下[]byte(s)，结果十进制：")
	fmt.Println([]byte(s))

	fmt.Print("打印下[]byte(s)中存储的类型，存储的是十六进制：")
	fmt.Printf("%#v\n", []byte(s))

	s1 := "世界"
	fmt.Print("打印下[]byte(s1)，结果十进制：")
	fmt.Println([]byte(s1))

	fmt.Print("打印下[]byte(s1)中存储的类型，存储的是十六进制：")
	fmt.Printf("%#v\n", []byte(s1))

	fmt.Print("打印下s1的十六进制：")
	fmt.Printf("%x\n", s1)

	fmt.Println([]rune(s1))
}
