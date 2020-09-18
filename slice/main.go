package main

import "fmt"

func main() {
	s1 := []int{2, 3, 6, 2, 4, 5, 6, 7}
	fmt.Println("原切片")
	fmt.Println("cap", cap(s1), "len", len(s1))

	s2 := s1[6:7]
	fmt.Println("新切片")
	fmt.Println("cap", cap(s2), "len", len(s2))
	fmt.Println(s2)

	s2 = append(s2, 100)
	fmt.Println("append之后的新切片")
	fmt.Println("cap", cap(s2), "len", len(s2))
	fmt.Println(cap(s2), len(s2))
	fmt.Println(s2)

	fmt.Println("老切片")
	fmt.Println(s1)

	s2 = append(s2, 888)
	fmt.Println("append之后的新切片，发生扩容")
	fmt.Println("cap", cap(s2), "len", len(s2))
	fmt.Println(cap(s2), len(s2))
	fmt.Println(s2)

	fmt.Println("老切片")
	fmt.Println(s1)

}
