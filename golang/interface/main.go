package main

import (
	"fmt"
)

type Sortable interface {
	Len() int
	Less(int, int) bool
	Swap(int, int)
}

func bubbleSort(array Sortable) {
	for i := 0; i < array.Len(); i++ {
		for j := 0; j < array.Len()-i-1; j++ {
			if array.Less(j+1, j) {
				array.Swap(j, j+1)
			}
		}
	}
}

//实现接口的整型切片
type IntArr []int

func (array IntArr) Len() int {
	return len(array)
}

func (array IntArr) Less(i int, j int) bool {
	return array[i] < array[j]
}

func (array IntArr) Swap(i int, j int) {
	array[i], array[j] = array[j], array[i]
}

//实现接口的字符串，按照长度排序
type StringArr []string

func (array StringArr) Len() int {
	return len(array)
}

func (array StringArr) Less(i int, j int) bool {
	return len(array[i]) < len(array[j])
}

func (array StringArr) Swap(i int, j int) {
	array[i], array[j] = array[j], array[i]
}

type Coder interface {
	code()
}

type CoderName struct {
	name string
}

func (g CoderName) code() {
	fmt.Printf("%s is coding\n", g.name)
}

func main() {
	var f interface{}
	fmt.Println("+++动态类型和动态值都是nil+++")
	fmt.Println(f == nil)
	fmt.Printf("f: %T, %v\n", f, f)

	var g *string
	f = g
	fmt.Println("+++类型为 *string+++")
	fmt.Println(f == nil)
	fmt.Printf("f: %T, %v\n", f, f)
}
