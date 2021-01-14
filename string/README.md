<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [go中string是如何实现的呢](#go%E4%B8%ADstring%E6%98%AF%E5%A6%82%E4%BD%95%E5%AE%9E%E7%8E%B0%E7%9A%84%E5%91%A2)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [实现](#%E5%AE%9E%E7%8E%B0)
  - [字符类型](#%E5%AD%97%E7%AC%A6%E7%B1%BB%E5%9E%8B)
    - [什么是字符集](#%E4%BB%80%E4%B9%88%E6%98%AF%E5%AD%97%E7%AC%A6%E9%9B%86)
    - [字库表](#%E5%AD%97%E5%BA%93%E8%A1%A8)
    - [编码字符集(字符集)](#%E7%BC%96%E7%A0%81%E5%AD%97%E7%AC%A6%E9%9B%86%E5%AD%97%E7%AC%A6%E9%9B%86)
    - [字符编码(编码方式)](#%E5%AD%97%E7%AC%A6%E7%BC%96%E7%A0%81%E7%BC%96%E7%A0%81%E6%96%B9%E5%BC%8F)
    - [ASCII 码](#ascii-%E7%A0%81)
    - [Unicode](#unicode)
    - [Unicode和ASCII的区别](#unicode%E5%92%8Cascii%E7%9A%84%E5%8C%BA%E5%88%AB)
  - [UTF-8编码](#utf-8%E7%BC%96%E7%A0%81)
    - [UTF-8和Unicode的关系](#utf-8%E5%92%8Cunicode%E7%9A%84%E5%85%B3%E7%B3%BB)
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








### 参考

【字符串】https://draveness.me/golang/docs/part2-foundation/ch03-datastructure/golang-string/  
【字符串】https://chai2010.gitbooks.io/advanced-go-programming-book/content/ch1-basic/ch1-03-array-string-and-slice.html  
【Golang中的string实现】https://erenming.com/2019/12/11/string-in-golang/    
【字符串】https://gfw.go101.org/article/string.html  
【Go string 实现原理剖析（你真的了解string吗）】https://my.oschina.net/renhc/blog/3019849    
【十分钟搞清字符集和字符编码】http://cenalulu.github.io/linux/character-encoding/  
【字符集和字符编码（Charset & Encoding）】https://www.cnblogs.com/skynet/archive/2011/05/03/2035105.html    
【字符编码笔记：ASCII，Unicode 和 UTF-8】http://www.ruanyifeng.com/blog/2007/10/ascii_unicode_and_utf-8.html  
【字符集和字符编码（Charset & Encoding）】https://cloud.tencent.com/developer/article/1347609   
【你不知道的 字符集和编码（编码字符集与字符集编码）】https://developer.aliyun.com/article/263676  
【字符集详解】https://blog.csdn.net/qq_42068856/article/details/83792174    
【Unicode,ASCII,UTF-8的区别】https://www.jianshu.com/p/8c57d87a76c6    