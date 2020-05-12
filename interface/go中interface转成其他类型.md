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

如果一个类型实现了一个 `interface` 中所有方法，必须是所有的方法，我们说类型实现了该 `interface`，所以所有类型都实现了 `empty interface`，因为任何一种类型至少实现了 0 个方法。`go` 没有显式的关键字用来实现 `interface`，只需要实现 `interface` 包含的方法即可。  

`interface`还可以作为返回值使用。  

为什么要使用`interface`

Francesc 给出了下面上个理由：  

1、writing generic algorithm  
2、hiding implementation detail  
3、providing interception points  

#### writing generic algorithm

首先我们应该清楚一点。`go`中是没有泛型的。

>Why does Go not have generic types? Generics may well be added at some point. We don’t feel an urgency for them.Generics are convenient but they come at a cost in complexity in the type system and run-time… Meanwhile, Go’s built-in maps and slices, plus the ability to use the empty interface to construct containers mean in many cases it is possible to write code that does what generics would enable, if less smoothly.

但是通过使用 `interface` 我们可以实现“泛型编程”，为什么？因为 `interface` 是一种抽象类型，任何具体类型（int, string）和抽象类型（user defined）都可以封装成 `interface`。

我们来看下标准库中的`sort`

````go
package sort

// A type, typically a collection, that satisfies sort.Interface can be
// sorted by the routines in this package.  The methods require that the
// elements of the collection be enumerated by an integer index.
type Interface interface {
    // Len is the number of elements in the collection.
    Len() int
    // Less reports whether the element with
    // index i should sort before the element with index j.
    Less(i, j int) bool
    // Swap swaps the elements with indexes i and j.
    Swap(i, j int)
}



// Sort sorts data.
// It makes one call to data.Len to determine n, and O(n*log(n)) calls to
// data.Less and data.Swap. The sort is not guaranteed to be stable.
func Sort(data Interface) {
    // Switch to heapsort if depth of 2*ceil(lg(n+1)) is reached.
    n := data.Len()
    maxDepth := 0
    for i := n; i > 0; i >>= 1 {
        maxDepth++
    }
    maxDepth *= 2
    quickSort(data, 0, n, maxDepth)
}
````
// TODO

### hiding implement detail







### 参考
【理解 Go interface 的 5 个关键点】https://sanyuesha.com/2017/07/22/how-to-understand-go-interface/ 
【深入理解 Go Interface】https://zhuanlan.zhihu.com/p/32926119   
【GO如何支持泛型】https://zhuanlan.zhihu.com/p/74525591  
【Golang面向对象编程】https://code.tutsplus.com/zh-hans/tutorials/lets-go-object-oriented-programming-in-golang--cms-26540  
【深度解密Go语言之关于 interface 的10个问题】https://www.cnblogs.com/qcrao-2018/p/10766091.html  