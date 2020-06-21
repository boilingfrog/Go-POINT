package main

import (
	"bufio"
	"bytes"
	"fmt"
)

func main() {
	var s = bytes.NewBuffer([]byte{})
	//var w = bufio.NewWriter(s)
	//w.WriteString("你好")
	//fmt.Println(s.String())

	r := bufio.NewReader(s)

	s.WriteString("哈哈哈1\n")
	s.WriteString("哈哈哈2\n")
	s.WriteString("哈哈哈3\n")
	fmt.Println(r.Size())
	//var dd byte

	//fmt.Println(r.ReadString(dd))
	p, _ := r.Peek(100)

	fmt.Printf("%s\n", p)

}
