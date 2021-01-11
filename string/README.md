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



### 字符类型

#### 什么是字符集

字符集就规定了某个文字对应的二进制数字存放方式（编码）和某串二进制数值代表了哪个文字（解码）的转换关系。  

简单的讲就是一个规则集合的名字，比如我们不同的国家有不同的语言。  

为什么会有多套呢？  

很多规范和标准在最初制定时并不会意识到这将会是以后全球普适的准则，或者处于组织本身利益就想从本质上区别于现有标准。于是，就产生了那么多具有相同效果但又不相互兼容的标准了。  

字符集转换需要的三个关键元素  

**字库表（character repertoire）**





### 参考

【字符串】https://draveness.me/golang/docs/part2-foundation/ch03-datastructure/golang-string/  
【字符串】https://chai2010.gitbooks.io/advanced-go-programming-book/content/ch1-basic/ch1-03-array-string-and-slice.html  
【Golang中的string实现】https://erenming.com/2019/12/11/string-in-golang/    
【字符串】https://gfw.go101.org/article/string.html  
【Go string 实现原理剖析（你真的了解string吗）】https://my.oschina.net/renhc/blog/3019849    
【http://cenalulu.github.io/linux/character-encoding/】http://cenalulu.github.io/linux/character-encoding/  
【字符集和字符编码（Charset & Encoding）】https://www.cnblogs.com/skynet/archive/2011/05/03/2035105.html    
【字符编码笔记：ASCII，Unicode 和 UTF-8】http://www.ruanyifeng.com/blog/2007/10/ascii_unicode_and_utf-8.html  