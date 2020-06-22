package main

import (
	"bufio"
	"fmt"
	"strings"
	"time"
)

func main() {
	reader := bufio.NewReaderSize(strings.NewReader("hello world"), 12)
	go Peek(reader)
	go reader.ReadBytes('d')
	time.Sleep(1e8)
}

func Peek(reader *bufio.Reader) {
	line, _ := reader.Peek(5)
	fmt.Printf("%s\n", line)
	time.Sleep(1)
	fmt.Printf("%s\n", line)
}
