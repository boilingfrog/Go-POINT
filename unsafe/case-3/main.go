package main

import "fmt"

func change(sl []int64) {
	sl[0] = 2
}

func changeNo(sl []int64) {
	s2 := make([]int64, 2)
	copy(sl, s2)
	s2[0] = 2
}

func main() {

	var sl = make([]int64, 2)
	change(sl)
	fmt.Println(sl)

	changeNo(sl)
	fmt.Println(sl)
}
