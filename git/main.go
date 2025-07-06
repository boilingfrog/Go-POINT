package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println(GetMaxNumber(10, 100))
}

func GetMaxNumber(a, b int) int {
	log.Printf("Comparing numbers: a=%d, b=%d", a, b)

	if a > b {
		log.Printf("Result: %d > %d", a, b)
		return a
	} else if a < b {
		log.Printf("Result: %d < %d", a, b)
		return b
	} else {
		log.Printf("Result: %d == %d", a, b)
		return a
	}
}
