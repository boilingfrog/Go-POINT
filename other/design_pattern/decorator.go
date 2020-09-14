package main

import (
	"fmt"
)

/*
装饰者模式(Decorator Pattern)：

动态地给一个对象增加一些额外的职责，增加对象功能来说，装饰模式比生成子类实现更为灵活。装饰模式是一种对象结构型模式。

在装饰者模式中，为了让系统具有更好的灵活性和可扩展性，我们通常会定义一个抽象
*/

type Base interface {
	Eat()
}

type Cat struct {
}

func (c *Cat) Eat() {
	fmt.Println("吃东西")
}

type BigCat struct {
	eat Base
}

func (b *BigCat) Eat() {
	b.eat.Eat()
	fmt.Println("+++++++++++++++,吃两份")
}
