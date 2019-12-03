package main

import (
	"fmt"
	"reflect"
)

func main() {

	v1 := uint(12)
	v2 := int(12)

	fmt.Println(reflect.TypeOf(v1)) //uint
	fmt.Println(reflect.TypeOf(v2)) //int

	fmt.Println(reflect.TypeOf(&v1)) //*uint
	fmt.Println(reflect.TypeOf(&v2)) //*int

	p := &v1

	//两个变量的类型不同,不能赋值
	//p = &v2 //cannot use &v2 (type *int) as type *uint in assignment

	fmt.Println(reflect.TypeOf(p)) // *unit
}
