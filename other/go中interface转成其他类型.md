## go中interface转换成原来的类型

### 首先了解下interface

什么是`interface`?

首先 `interface` 是一种类型，从它的定义可以看出来用了 `type` 关键字，更准确的说 `interface` 是一种具有一组方法的类型，这些方法定义了 `interface` 的行为。

````go
type I interface {
    Get() int
}
````
 



### 参考
【理解 Go interface 的 5 个关键点】https://sanyuesha.com/2017/07/22/how-to-understand-go-interface/ 
【深入理解 Go Interface】https://zhuanlan.zhihu.com/p/32926119   