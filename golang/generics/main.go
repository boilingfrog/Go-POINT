package main

import "fmt"

func Pair[A, B any](a A, b B) (A, B) {
	return a, b
}

func main() {
	x, y := Pair(10, "Go") // 推断 A=int，B=string
	fmt.Printf("%T(%v), %T(%v)\n", x, x, y, y)

	a, b := Pair[int](20, "泛型") // A 显式指定，B 继续推断
	fmt.Printf("%T(%v), %T(%v)\n", a, a, b, b)
}
