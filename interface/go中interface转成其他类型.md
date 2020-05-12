## go中interface转换成原来的类型

### 首先了解下interface

#### 什么是`interface`?

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

如果一个类型实现了一个 `interface` 中所有方法，必须是所有的方法，我们说类型实现了该 `interface`，所以所有类型都实现了 `empty interface`，因为任何一种类型至少实现了 0 个方法。`go` 没有显式的关键字用来实现 `interface`，只需要实现 `interface` 包含的方法即可。  

`interface`还可以作为返回值使用。  

#### 如何判断 interface 变量存储的是哪种类型

日常中使用`interface`,有时候需要判断原来是什么类型的值转成了`interface`。一般有以下几种方式：

##### fmt

````go
import "fmt"
func main() {
    v := "hello world"
    fmt.Println(typeof(v))
}
func typeof(v interface{}) string {
    return fmt.Sprintf("%T", v)
}
````

##### 反射

````
import (
    "reflect"
    "fmt"
)
func main() {
    v := "hello world"
    fmt.Println(typeof(v))
}
func typeof(v interface{}) string {
    return reflect.TypeOf(v).String()
}
````  

##### 断言

`Go`语言里面有一个语法，可以直接判断是否是该类型的变量： `value, ok = element.(T)`，这里`value`就是变量的值，`ok`是一个`bool`类型，`element`是`interface`变量，`T`是断言的类型。  

如果`element`里面确实存储了`T`类型的数值，那么ok返回`true`，否则返回`false`。  

让我们通过一个例子来更加深入的理解。  


类型不确定可以使用:

````go
func main() {
    v := "hello world"
    fmt.Println(typeof(v))
}
func typeof(v interface{}) string {
    switch t := v.(type) {
    case int:
        return "int"
    case float64:
        return "float64"
    //... etc
    default:
        _ = t
        return "unknown"
    }
}
````




### 参考
【理解 Go interface 的 5 个关键点】https://sanyuesha.com/2017/07/22/how-to-understand-go-interface/ 
【深入理解 Go Interface】https://zhuanlan.zhihu.com/p/32926119   
【GO如何支持泛型】https://zhuanlan.zhihu.com/p/74525591  
【Golang面向对象编程】https://code.tutsplus.com/zh-hans/tutorials/lets-go-object-oriented-programming-in-golang--cms-26540  
【深度解密Go语言之关于 interface 的10个问题】https://www.cnblogs.com/qcrao-2018/p/10766091.html  
【golang如何获取变量的类型：反射，类型断言】https://ieevee.com/tech/2017/07/29/go-type.html    