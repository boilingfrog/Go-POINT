## go中interface转换成原来的类型

### 首先了解下interface

什么是`interface`?

首先 `interface` 是一种类型，从它的定义可以看出来用了 `type` 关键字，更准确的说 `interface` 是一种具有一组方法的类型，这些方法定义了 `interface` 的行为。

````go
type I interface {
    Get() int
}
````
`interface`是一组`method`的集合，是`duck-type programming`的一种体现（不关心属性（数据），只关心行为（方法））。我们可以自己定义`interface`类型的`struct`,并提供方法。

````go
type MyInterface interface{
    Print()
}

func TestFunc(x MyInterface) {}
type MyStruct struct {}
func (me MyStruct) Print() {}

func main() {
    var me MyStruct
    TestFunc(me)
}
````

`go` 允许不带任何方法的 `interface` ，这种类型的 `interface` 叫 `empty interface`。  

如果一个类型实现了一个 `interface` 中所有方法，我们说类型实现了该 `interface`，所以所有类型都实现了 `empty interface`，因为任何一种类型至少实现了 0 个方法。`go` 没有显式的关键字用来实现 `interface`，只需要实现 `interface` 包含的方法即可。


### 参考
【理解 Go interface 的 5 个关键点】https://sanyuesha.com/2017/07/22/how-to-understand-go-interface/ 
【深入理解 Go Interface】https://zhuanlan.zhihu.com/p/32926119   