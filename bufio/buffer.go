package main

import (
	"bufio"
	"fmt"
	"strings"
)

func main() {
	reader := bufio.NewReader(strings.NewReader("hello \n world"))
	line1, _ := reader.ReadSlice('\n')
	fmt.Printf("the line1:%s\n", line1)

	line2, _ := reader.ReadSlice('\n')
	fmt.Printf("the line2:%s\n", line2)

}
