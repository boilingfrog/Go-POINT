package main

import (
	"fmt"
	"strings"
)

func main() {
	strs := []string{"小白", "小红", "小李"}
	res := strings.Join(strs, ",")
	fmt.Println(res)
}
