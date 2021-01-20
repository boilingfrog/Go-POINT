package main

import (
	"bytes"
	"fmt"
)

func main() {
	var b bytes.Buffer
	b.WriteString("hello")
	b.WriteString("world")
	fmt.Println(b.String())
}
