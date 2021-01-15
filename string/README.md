<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [go中string是如何实现的呢](#go%E4%B8%ADstring%E6%98%AF%E5%A6%82%E4%BD%95%E5%AE%9E%E7%8E%B0%E7%9A%84%E5%91%A2)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [实现](#%E5%AE%9E%E7%8E%B0)
  - [字符串的拼接](#%E5%AD%97%E7%AC%A6%E4%B8%B2%E7%9A%84%E6%8B%BC%E6%8E%A5)
  - [字符串](#%E5%AD%97%E7%AC%A6%E4%B8%B2)
  - [字符类型](#%E5%AD%97%E7%AC%A6%E7%B1%BB%E5%9E%8B)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## go中string是如何实现的呢


### 前言

go中的string可谓是用到的最频繁的关键词之一了，如何实现，我们来探究下  

### 实现

```go
// string is the set of all strings of 8-bit bytes, conventionally but not
// necessarily representing UTF-8-encoded text. A string may be empty, but
// not nil. Values of string type are immutable.
type string string
```

string我们看起来是一个整体，但是本质上是一片连续的内存空间，我们也可以将它理解成一个由字符组成的数组，相比于切片仅仅少了一个Cap属性。  

`src/reflect/value.go`  
```go
type StringHeader struct {
	Data uintptr
	Len  int
}
```

切片的数据结构

```go
type SliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}
```

1、相比于切片少了一个容量的cap字段，就意味着string是不能发生地址空间扩容；  

2、可以把string当成一个只读的切片类型；  

3、string本身的切片是只读的，所以不会直接向字符串直接追加元素改变其本身的内存空间，所有在字符串上的写入操作都是通过拷贝实现的。  

### 字符串的拼接



### 字符串


### 字符类型

我们在go中经常遇到rune和byte两种字符串类型，作为go中字符串的两种类型：  

- byte 也叫 uint8。代表了 ASCII 码的一个字符。  

- rune 等价于 int32 类型，代表一个 UTF-8 字符，当需要处理中文、日文或者其他复合字符时，则需要用到 rune 类型。  

```go
// byte is an alias for uint8 and is equivalent to uint8 in all ways. It is
// used, by convention, to distinguish byte values from 8-bit unsigned
// integer values.
type byte = uint8

// rune is an alias for int32 and is equivalent to int32 in all ways. It is
// used, by convention, to distinguish character values from integer values.
type rune = int32
```

### 参考

【字符串】https://draveness.me/golang/docs/part2-foundation/ch03-datastructure/golang-string/  
【字符串】https://chai2010.gitbooks.io/advanced-go-programming-book/content/ch1-basic/ch1-03-array-string-and-slice.html  
【Golang中的string实现】https://erenming.com/2019/12/11/string-in-golang/    
【字符串】https://gfw.go101.org/article/string.html  
【Go string 实现原理剖析（你真的了解string吗）】https://my.oschina.net/renhc/blog/3019849    
【Go语言字符串高效拼接（一）】https://cloud.tencent.com/developer/article/1367934    
【Go语言字符类型（byte和rune）】http://c.biancheng.net/view/18.html    
【go 的 [] rune 和 [] byte 区别】https://learnku.com/articles/23411/the-difference-between-rune-and-byte-of-go  
