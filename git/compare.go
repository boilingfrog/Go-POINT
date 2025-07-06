package main

import "log"

func CompareNumbers(a, b int) int {
	log.Printf("Comparing numbers: a=%f, b=%f", a, b)

	if a > b {
		log.Printf("Result: %f > %f", a, b)
		return a
	} else if a < b {
		log.Printf("Result: %f < %f", a, b)
		return b
	} else {
		log.Printf("Result: %f == %f", a, b)
		return a
	}
}
