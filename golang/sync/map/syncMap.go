package main

import (
	"fmt"
	"sync"
)

func syncMapDemo() {

	var smp sync.Map

	// 数据写入
	smp.Store("name", "小红")
	smp.Store("age", 18)

	// 数据读取
	name, _ := smp.Load("name")
	fmt.Println(name)

	age, _ := smp.Load("age")
	fmt.Println(age)

	// 遍历
	smp.Range(func(key, value interface{}) bool {
		fmt.Println(key, value)
		return true
	})

	// 删除
	smp.Delete("age")
	age, ok := smp.Load("age")
	fmt.Println("删除后的查询")
	fmt.Println(age, ok)

	// 读取或写入,存在就读取，不存在就写入
	smp.LoadOrStore("age", 100)
	age, _ = smp.Load("age")
	fmt.Println("不存在")
	fmt.Println(age)

	smp.LoadOrStore("age", 99)
	age, _ = smp.Load("age")
	fmt.Println("存在")
	fmt.Println(age)
}
