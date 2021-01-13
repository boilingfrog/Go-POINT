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



### 字符类型

#### 什么是字符集

字符集就规定了某个文字对应的二进制数字存放方式（编码）和某串二进制数值代表了哪个文字（解码）的转换关系。  

简单的讲就是一个规则集合的名字，比如我们不同的国家有不同的语言。  

为什么会有多套呢？  

很多规范和标准在最初制定时并不会意识到这将会是以后全球普适的准则，或者处于组织本身利益就想从本质上区别于现有标准。于是，就产生了那么多具有相同效果但又不相互兼容的标准了。  

字符集转换需要的三个关键元素  

#### 字库表

字库表（character repertoire）  

相当于所有可读或者可显示字符的数据库，字库表决定了整个字符集能够展现表示的所有字符的范围。  

#### 编码字符集(字符集)

编码字符集（coded character set）（字符集）  

编码字符集也就是我们经常提到的字符集，在对应的字库表中，每个字符都有对应的二进制地址，编码字符集就是这些地址的集合。  

例如：在`ASCII`中A在表中排第65位，而编码后A的数值是`0100 0001`，也就是十进制的65的二进制转换结果。编码字符集就是用来存储这些二进制数的。而这个二进制数就是编码字符集中的一个元素，同时它也是字库表中字母A的地址。我们根据这个地址就可以显示出字母A。  

#### 字符编码(编码方式)

字符编码（character encoding form）  

编码字符集和实际存储数值之间的转换关系，就是把字符集中的字符编码为指定集合中的某一对象，以便在计算机中存储和通过网络传递。  

知道了字库表和编码字符集后，我们是可以直接通过二进制来得到字符串的，但是还是引入了字符编码。  

原因很简单：统一字库表的目的是为了能够涵盖世界上所有的字符，但实际使用过程中会发现真正用的上的字符相对整个字库表来说比例非常低。例如中文地区的程序几乎不会需要日语字符，而一些英语国家甚至简单的ASCII字库表就能满足基本需求。而如果把每个字符都用字库表中的序号来存储的话，每个字符就需要3个字节（这里以Unicode字库为例），这样对于原本用仅占一个字符的ASCII编码的英语地区国家显然是一个额外成本（存储体积是原来的三倍）。这就导致一个可以用1G保存的文件，现在需要3G才能保存，这是极其浪费的做法。  

所以制定了字符编码来处理这种问题。每种不同的算法被称为一种编码方式。一套编码规范可以有不同的编码方式，不同的编码方式有不同的适应场景。  

例如：UTF-8这样的变长编码。在UTF-8编码中原本只需要一个字节的ASCII字符，仍然只占一个字节。而像中文及日语这样的复杂字符就需要2个到3个字节来存储。  


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