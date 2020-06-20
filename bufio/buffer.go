package main

import (
	"bufio"
	"bytes"
	"fmt"
)

func main() {
	var s = bytes.NewBuffer([]byte{})
	var w = bufio.NewWriter(s)
	w.WriteString("你好")
	fmt.Println(s.String())
}
