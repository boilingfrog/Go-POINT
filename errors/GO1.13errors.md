## go1.13errors的用法

### 前言

go 1.13发布了`error`的一些新的特性，那么就来探究学习下。

### 基本用法

#### fmt.Errorf

使用 `fmt.Errorf` 加上 `%w` 格式符来生成一个嵌套的 `error`，它并没有像 `pkg/errors` 那样使用一个 `Wrap` 函数来嵌套` error`，非常简洁。  

````go
err1 := errors.New("new error")
err2 := fmt.Errorf("err2: [%w]", err1)
err3 := fmt.Errorf("err3: [%w]", err2)
fmt.Println(err3)
````
输出
````
// output
err3: [err2: [new error]]
````
`err2` 就是一个合法的被包装的 `error`，同样地，`err3` 也是一个被包装的 `error`，如此可以一直套下去。  







### 参考
 
【Go 1.13 errors 基本用法】https://segmentfault.com/a/1190000020398774