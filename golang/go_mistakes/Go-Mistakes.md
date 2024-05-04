<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [<<100 Go Mistakes and How to Avoid Them>> 阅读记录](#100-go-mistakes-and-how-to-avoid-them-%E9%98%85%E8%AF%BB%E8%AE%B0%E5%BD%95)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [slice](#slice)
    - [slice 的长度和容量理解有误](#slice-%E7%9A%84%E9%95%BF%E5%BA%A6%E5%92%8C%E5%AE%B9%E9%87%8F%E7%90%86%E8%A7%A3%E6%9C%89%E8%AF%AF)
    - [不高效的 slice 初始化](#%E4%B8%8D%E9%AB%98%E6%95%88%E7%9A%84-slice-%E5%88%9D%E5%A7%8B%E5%8C%96)
    - [使用 nil 表示空 slice](#%E4%BD%BF%E7%94%A8-nil-%E8%A1%A8%E7%A4%BA%E7%A9%BA-slice)
    - [通过判断长度来确定 slice 是否为空](#%E9%80%9A%E8%BF%87%E5%88%A4%E6%96%AD%E9%95%BF%E5%BA%A6%E6%9D%A5%E7%A1%AE%E5%AE%9A-slice-%E6%98%AF%E5%90%A6%E4%B8%BA%E7%A9%BA)
    - [切片的复制](#%E5%88%87%E7%89%87%E7%9A%84%E5%A4%8D%E5%88%B6)
    - [切线的截取](#%E5%88%87%E7%BA%BF%E7%9A%84%E6%88%AA%E5%8F%96)
    - [内存](#%E5%86%85%E5%AD%98)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## <<100 Go Mistakes and How to Avoid Them>> 阅读记录

### 前言

最近阅读了 <<100 Go Mistakes and How to Avoid Them>> 这本书，总结下，其中几个自己印象深刻的常见错误。 

### slice

#### slice 的长度和容量理解有误

slice 长度是 slice 已经存储的元素的数量，容量指的是 slice 当前底层开辟的数组最多能容纳的元素的数量。  

这里简单看下切片的数据结构  

```go
type slice struct {
    array unsafe.Pointer // 指针
    len   int // 长度
    cap   int // 容量
}
```

切片一共三个属性：指针，指向底层的数组；长度，表示切片可用元素的个数。  

也就是说使用下标 对元素进行访问的时候，下标不能超过的长度；容量，底层数组的元素个数，容量>=长度。  

#### 不高效的 slice 初始化

当创建一个 slice 时，如果其长度可以预先确定，那么可以在定义时指定它的长度和容量。这可以改善后期 append 时一次或者多次的内存分配操作，从而改善性能。对于 map 的初始化也是这样的。  

当切片的长度超过其容量时，切片会自动扩容。这通常发生在使用 append 函数向切片中添加元素时。  

扩容时，Go 运行时会分配一个新的底层数组，并将原始切片中的元素复制到新数组中。然后，原始切片将指向新数组，并更新其长度和容量。  

需要注意的是，由于扩容会分配新数组并复制元素，因此可能会影响性能。所以如果能知道 slice 的长度，我们可以提前进行初始化，避免后面的扩容发生。  

#### 使用 nil 表示空 slice

使用 nil 表示空 slice 的优点：nil 不分配内存，它的长度和容量都为0 。对于一个函数的返回值而言，返回 `nil slice` 比 `emtpy slice` 要更好。    

在 marshal 时，`nil slice` 是 null，而 `empty slice` 是 []。因此在使用相关库函数时，要特别注意这两者的区别。    

空切片不 nil 不相等。  

#### 通过判断长度来确定 slice 是否为空

检查一个slice的是否包含任何元素，可以检查其长度，不管slice是nil还是empty，检查长度都是有效的。这个检查方法也适用于map。  

为了设计更明确的API，API不应区分nil和空切片。  

#### 切片的复制

Go 语言的内置函数 copy() 可以将一个数组切片复制到另一个数组切片中，如果加入的两个数组切片不一样大，就会按照其中较小的那个数组切片的元素个数进行复制。  

> 内置函数copy从源切片复制元素到目标切片。（作为特殊情况，它还可以将字符串中的字节复制到字节切片中。）
> 源和目标可能会有重叠。Copy返回复制的元素数量，这将是len(src)和len(dst)的最小值。

```go
func main()  {
	slice1 := []int{1, 2, 3, 4, 5}
	slice2 := []int{5, 4, 3}
	copy(slice2, slice1)
	fmt.Println(slice2)
}

输出 [1 2 3]

func main()  {
	slice1 := []int{1, 2, 3, 4, 5}
	slice2 := []int{5, 4, 3}
	copy(slice1, slice2)
	fmt.Println(slice1)
}

输出 [5 4 3 4 5]
```

#### 切线的截取

基于已有 slice 创建新 slice 对象，被称为 reslice。新 slice 和老 slice 共用底层数组，新老 slice 对底层数组的更改都会影响到彼此。基于数 组创建的新 slice 对象也是同样的效果：对数组或 slice 元素作的更改都会影响到彼此。  

如果新截取的切片发生了扩容了，就会重新申请新的内存空间，这样就新截取的切片就不指向原来的地址空间了。  

#### 内存

### 参考

【100 Go Mistakes and How to Avoid Them】https://100go.co/zh/#1     