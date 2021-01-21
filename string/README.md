<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [go中string是如何实现的呢](#go%E4%B8%ADstring%E6%98%AF%E5%A6%82%E4%BD%95%E5%AE%9E%E7%8E%B0%E7%9A%84%E5%91%A2)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [实现](#%E5%AE%9E%E7%8E%B0)
  - [[]byte转string](#byte%E8%BD%ACstring)
  - [string转[]byte](#string%E8%BD%ACbyte)
  - [字符串的拼接](#%E5%AD%97%E7%AC%A6%E4%B8%B2%E7%9A%84%E6%8B%BC%E6%8E%A5)
    - [+方式进行拼接](#%E6%96%B9%E5%BC%8F%E8%BF%9B%E8%A1%8C%E6%8B%BC%E6%8E%A5)
    - [fmt 拼接](#fmt-%E6%8B%BC%E6%8E%A5)
    - [Join 拼接](#join-%E6%8B%BC%E6%8E%A5)
    - [buffer 拼接](#buffer-%E6%8B%BC%E6%8E%A5)
    - [builder 拼接](#builder-%E6%8B%BC%E6%8E%A5)
    - [测试下几种方法的性能](#%E6%B5%8B%E8%AF%95%E4%B8%8B%E5%87%A0%E7%A7%8D%E6%96%B9%E6%B3%95%E7%9A%84%E6%80%A7%E8%83%BD)
  - [字符类型](#%E5%AD%97%E7%AC%A6%E7%B1%BB%E5%9E%8B)
    - [byte](#byte)
    - [rune](#rune)
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

### []byte转string

`src/runtime/string.go`
```go
type stringStruct struct {
	str unsafe.Pointer
	len int
}

func slicebytetostring(buf *tmpBuf, b []byte) (str string) {
	l := len(b)
	if l == 0 {
		return ""
	}
	if l == 1 {
		stringStructOf(&str).str = unsafe.Pointer(&staticbytes[b[0]])
		stringStructOf(&str).len = 1
		return
	}
	var p unsafe.Pointer
	// 判断传入的缓冲区大小，决定是否重新分配内存
	if buf != nil && len(b) <= len(buf) {
		p = unsafe.Pointer(buf)
	} else {
		// 重新分配内存
		p = mallocgc(uintptr(len(b)), nil, false)
	}
	// 将输出的str转化成stringStruct结构
	// 并且赋值
	stringStructOf(&str).str = p
	stringStructOf(&str).len = len(b)
	// 将[]byte中的内容，复制到内存空间p中
	memmove(p, (*(*slice)(unsafe.Pointer(&b))).array, uintptr(len(b)))
	return
}

// 转换成
func stringStructOf(sp *string) *stringStruct {
	return (*stringStruct)(unsafe.Pointer(sp))
}
```

总结下流程：  

1、根据传入的内存大小，判断是否需要分配重新分配内存；  
2、构建stringStruct，分类长度和内存空间;    
3、赋值[]byte里面的数据到新构建stringStruct的内存空间中。

### string转[]byte

`src/runtime/string.go`
```go
func stringtoslicebyte(buf *tmpBuf, s string) []byte {
	var b []byte
	if buf != nil && len(s) <= len(buf) {
		*buf = tmpBuf{}
		b = buf[:len(s)]
	} else {
		b = rawbyteslice(len(s))
	}
	copy(b, s)
	return b
}

// 为[]byte重新分配一段内存
func rawbyteslice(size int) (b []byte) {
	cap := roundupsize(uintptr(size))
	p := mallocgc(cap, nil, false)
	if cap != uintptr(size) {
		memclrNoHeapPointers(add(p, uintptr(size)), cap-uintptr(size))
	}

	*(*slice)(unsafe.Pointer(&b)) = slice{p, size, int(cap)}
	return
}
```

1、判断传入的缓存区大小，如果内存够用就使用传入的缓冲区存储 []byte；  
2、传入的缓存区的大小不够，调用 `runtime.rawbyteslice`创建指定大小的[]byte；  
3、将string拷贝到切片。  

### 字符串的拼接

#### +方式进行拼接

```go
func main() {
	s := "hai~"
	s += "hello world"
	fmt.Println(s)
}
```

#### fmt 拼接

```go
func main() {
	s := fmt.Sprintf("%s%s%d", "hello", "world", 2021)
	fmt.Println(s)
}
```

fmt可以方便对各种类型的数据进行拼接，转换成string,具体详见[printf的用法](https://github.com/boilingfrog/Go-POINT/blob/master/printf/main.go)  

#### Join 拼接

```go
func main() {
	var s []string

	s = append(s, "hello")
	s = append(s, "world")

	fmt.Println(strings.Join(s, ""))
}
```

#### buffer 拼接

```go
func main() {
	var b bytes.Buffer
	b.WriteString("hello")
	b.WriteString("world")
	fmt.Println(b.String())
}
```

#### builder 拼接

```go
func main() {
	var b bytes.Buffer
	b.WriteString("hello")
	b.WriteString("world")
	fmt.Println(b.String())
}
```

#### 测试下几种方法的性能

压力测试

```go
package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func String() string {
	var s string
	s += "hello" + "\n"
	s += "world" + "\n"
	s += "今天的天气很不错的"
	return s
}

func BenchmarkString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		String()
	}
}

func StringFmt() string {
	return fmt.Sprintf("%s %s %s", "hello", "world", "今天的天气很不错的")
}

func BenchmarkStringFmt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StringFmt()
	}
}

func StringJoin() string {
	var s []string
	s = append(s, "hello ")
	s = append(s, "world ")
	s = append(s, "今天的天气很不错的 ")

	return strings.Join(s, "")
}

func BenchmarkStringJoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StringJoin()
	}
}

func StringBuffer() string {
	var s bytes.Buffer
	s.WriteString("hello ")
	s.WriteString("world ")
	s.WriteString("今天的天气很不错的 ")

	return s.String()
}

func BenchmarkStringBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StringBuffer()
	}
}

func StringBuilder() string {
	var s strings.Builder
	s.WriteString("hello ")
	s.WriteString("world ")
	s.WriteString("今天的天气很不错的 ")

	return s.String()
}

func BenchmarkStringBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StringBuilder()
	}
}
```

看下执行的结果

```go
$ go test string_test.go  -bench=. -benchmem -benchtime=3s
goos: darwin
goarch: amd64
BenchmarkString-4               28862838               115 ns/op              64 B/op          2 allocs/op
BenchmarkStringFmt-4            20697505               169 ns/op              48 B/op          1 allocs/op
BenchmarkStringJoin-4           11304583               293 ns/op             160 B/op          4 allocs/op
BenchmarkStringBuffer-4         31151836               104 ns/op             112 B/op          2 allocs/op
BenchmarkStringBuilder-4        29142151               120 ns/op              72 B/op          3 allocs/op
PASS
ok      command-line-arguments  17.740s
```

`ns/op` 平均一次执行的时间
`B/op` 平均一次真申请的内存大小
`allocs/op` 平均一次，申请的内存次数  

从上面我们就能直观的看出差距，不过差距不大，当然具体的性能信息要结合当前go版本，具体讨论。  

看上去很low的+拼接方式，在性能上倒是还不错。  

我的版本是`go version go1.13.15 darwin/amd64`  

### 字符类型

我们在go中经常遇到rune和byte两种字符串类型，作为go中字符串的两种类型：  

#### byte

byte 也叫 uint8。代表了 ASCII 码的一个字符。  

对于英文，一个ASCII表示一个字符，根据[ASCII表](https://baike.baidu.com/item/ASCII?fr=aladdin#reference-[1]-15482-wrap)。我们知道
A对应的十进制编码是`65`。我们看看byte打印的结果  

```go
func main() {
	s := "A"
	fmt.Print("打印下[]byte(s)，结果十进制：")
	fmt.Println([]byte(s))

	fmt.Print("打印下[]byte(s)中存储的类型，存储的是十六进制：")
	fmt.Printf("%#v\n", []byte(s))

	s1 := "世界"
	fmt.Print("打印下[]byte(s1)，结果十进制：")
	fmt.Println([]byte(s1))

	fmt.Print("打印下[]byte(s1)中存储的类型，存储的是十六进制：")
	fmt.Printf("%#v\n", []byte(s1))

	fmt.Print("打印下s1的十六进制：")
	fmt.Printf("%x\n", s1)
}
```

结果

```
打印下[]byte(s)，结果十进制：[65]
打印下[]byte(s)中存储的类型，存储的是十六进制：[]byte{0x41}
打印下[]byte(s1)，结果十进制：[228 184 150 231 149 140]
打印下[]byte(s1)中存储的类型，存储的是十六进制：[]byte{0xe4, 0xb8, 0x96, 0xe7, 0x95, 0x8c}
打印下s1的十六进制：e4b896e7958c
```

对于ASCII一个索引位表示一个字符（也就是英文）  

对于非ASCII，索引更新的步长将超过1个字节，中文的是三个字节表示一个中文。  

#### rune

rune 等价于 int32 类型，UTF8编码的Unicode码点。  

```go
func main() {
	s := "哈哈"
	fmt.Println([]rune(s))

	s1 := "A"
	fmt.Println([]rune(s1))
}
```

打印下结果

```
[21704 21704]
[65]
```

我们可以看到里面对应的是UTF-8的十进制数字。对于英文来讲UTF-8的码点，就是对应的ASCII。  

<img src="/img/string_6.jpg"  alt="string" align="center" />

```go
// byte is an alias for uint8 and is equivalent to uint8 in all ways. It is
// used, by convention, to distinguish byte values from 8-bit unsigned
// integer values.
type byte = uint8

// rune is an alias for int32 and is equivalent to int32 in all ways. It is
// used, by convention, to distinguish character values from integer values.
type rune = int32
```

关于Unicode和UTF8的区别和联系，以及ASCII码的联系，参考[字符编码-字库表,字符集,字符编码](https://www.cnblogs.com/ricklz/p/14271477.html#utf-8%E5%92%8Cunicode%E7%9A%84%E5%85%B3%E7%B3%BB)

### 参考

【字符串】https://draveness.me/golang/docs/part2-foundation/ch03-datastructure/golang-string/  
【字符串】https://chai2010.gitbooks.io/advanced-go-programming-book/content/ch1-basic/ch1-03-array-string-and-slice.html  
【Golang中的string实现】https://erenming.com/2019/12/11/string-in-golang/    
【字符串】https://gfw.go101.org/article/string.html  
【Go string 实现原理剖析（你真的了解string吗）】https://my.oschina.net/renhc/blog/3019849    
【Go语言字符串高效拼接（一）】https://cloud.tencent.com/developer/article/1367934    
【Go语言字符类型（byte和rune）】http://c.biancheng.net/view/18.html    
【go 的 [] rune 和 [] byte 区别】https://learnku.com/articles/23411/the-difference-between-rune-and-byte-of-go  
【go语言圣经】http://books.studygolang.com/gopl-zh/ch3/ch3-05.html  
【【Go语言踩坑系列（二）】字符串】https://www.mdeditor.tw/pl/pCg8  
