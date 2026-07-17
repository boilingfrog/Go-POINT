package main

import "log"

func CompareNumbers(a, b int) int {
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
