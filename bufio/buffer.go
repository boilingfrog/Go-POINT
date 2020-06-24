package main

import (
	"bufio"
	"bytes"
	"fmt"
)

func main() {

	var s = bytes.NewBuffer([]byte{})
	var w = bufio.NewWriter(s)
	w.WriteString("hello world")
	w.WriteString("你好")
	fmt.Printf("string--%s", s.String())
	fmt.Println()
	w.Flush()
	fmt.Printf("string--%s", s.String())
}
