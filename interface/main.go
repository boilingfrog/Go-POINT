package main

import "fmt"

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

//测试
func main() {
	intArray1 := IntArr{3, 4, 2, 6, 10, 1}
	bubbleSort(intArray1)
	fmt.Println(intArray1)

	stringArray1 := StringArr{"hello", "i", "am", "go", "lang"}
	bubbleSort(stringArray1)
	fmt.Println(stringArray1)

	v := "hello world"
	fmt.Println(typeof(v))
}

func typeof(v interface{}) string {
	return fmt.Sprintf("%T", v)
}
