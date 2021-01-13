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

#### ASCII 码

字节：

计算机内部，所有信息最终都是一个二进制值。每一个二进制位（bit）有0和1两种状态，因此八个二进制位就可以组合出256种状态，这被称为一个字节（byte）。也就是说，一个字节一共可以用来表示256种不同的状态，每一个状态对应一个符号，就是256个符号，从`00000000`到`11111111`。  

ASCII：  

英语字符与二进制位之间的关系。   

ASCII 码一共规定了128个字符的编码，比如空格SPACE是32（二进制00100000），大写的字母A是65（二进制01000001）。这128个符号（包括32个不能打印出来的控制符号），只占用了一个字节的后面7位，最前面的一位统一规定为0。  

ASCII是单字节编码，无法用来表示中文（中文编码至少需要2个字节），所以，中国制定了GB2312编码，用来把中文编进去。但世界上有许多不同的语言，所以需要一种统一的编码。  


#### Unicode

Unicode 只是一个符号集，它只规定了符号的二进制代码，却没有规定这个二进制代码应该如何存储。  

- Unicode把所有语言都统一到一套编码里，这样就不会再有乱码问题了。

- Unicode最常用的是用两个字节表示一个字符（如果要用到非常偏僻的字符，就需要4个字节）。现代操作系统和大多数编程语言都直接支持Unicode。

#### Unicode和ASCII的区别

- ASCII编码是1个字节，而Unicode编码通常是2个字节。  

字母A用ASCII编码是十进制的65，二进制的01000001；而在Unicode中，只需要在前面补0，即为：00000000 01000001。  

统一编码的问题就是会出现空间的浪费。所以引入了UTF-8,字符编码，来解决这一问题。  

### UTF-8编码

UTF-8编码可以把Unicode编码转换为“可变长编码”。  

UTF-8 最大的一个特点，就是它是一种变长的编码方式。它可以使用1~4个字节表示一个符号，根据不同的符号而变化字节长度。  

UTF-8编码把一个Unicode字符根据不同的数字大小编码成1-6个字节，常用的英文字母被编码成1个字节，汉字通常是3个字节，只有很生僻的字符才会被编码成4-6个字节。如果你要传输的文本包含大量英文字符，用UTF-8编码就能节省空间。  

UTF-8 的编码规则：  

1、对于单字节的符号，字节的第一位设为0，后面7位为这个符号的 Unicode 码。因此对于英语字母，UTF-8 编码和 ASCII 码是相同的。  

2、对于n字节的符号（n > 1），第一个字节的前n位都设为1，第n + 1位设为0，后面字节的前两位一律设为10。剩下的没有提及的二进制位，全部为这个符号的 Unicode 码。  

#### UTF-8和Unicode的关系

弄明白上面的几个概念，理解UTF-8和Unicode的关系就比较容易了。  

`Unicode`就是上文中提到的编码字符集。它只规定了符号的二进制代码，却没有规定这个二进制代码应该如何存储。  

`UTF-8`就是字符编码，即Unicode规则字库的一种实现形式。  


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