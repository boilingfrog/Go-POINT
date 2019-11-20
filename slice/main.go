package main

import "fmt"

func main() {
	s := make([]int, 0)

	oldCap := cap(s)

	for i := 0; i < 2048; i++ {
		s = append(s, i)

		newCap := cap(s)

		if newCap != oldCap {
			fmt.Printf("[%d -> %4d] cap = %-4d  |  after append %-4d  cap = %-4d\n", 0, i-1, oldCap, i, newCap)
			oldCap = newCap
		}
	}

	slice := *new([]int)

	fmt.Println(slice)

	s1 := []int{0, 1, 2, 3, 8: 100}
	fmt.Println(s1, len(s1), cap(s1))
	s2 := []int{1, 2}
	s2 = append(s2, 4, 5, 6)
	fmt.Printf("len=%d, cap=%d", len(s2), cap(s2))
	fmt.Println(s2)

	s3 := []int{1, 2}
	s3 = append(s3, 4)
	s3 = append(s3, 5)
	s3 = append(s3, 6)

	fmt.Printf("len=%d, cap=%d", len(s3), cap(s3))
	fmt.Println(s3)

	fmt.Println(2 / 4)
}
