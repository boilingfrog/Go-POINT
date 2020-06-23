package main

import (
	"bufio"
	"fmt"
	"strings"
	"time"
)

func main() {
	scanner := bufio.NewScanner(
		strings.NewReader("ABCDEFG\nHIJKELM"),
	)
	scanner.Split(bufio.ScanWords /*四种方式之一，你也可以自定义, 实现SplitFunc方法*/)
	for scanner.Scan() {
		fmt.Println(scanner.Text()) // scanner.Bytes()
	}
}

func Peek(reader *bufio.Reader) {
	line, _ := reader.Peek(5)
	fmt.Printf("%s\n", line)
	time.Sleep(1)
	fmt.Printf("%s\n", line)
}
