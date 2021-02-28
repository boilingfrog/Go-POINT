<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [unsafe](#unsafe)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [指针类型](#%E6%8C%87%E9%92%88%E7%B1%BB%E5%9E%8B)
    - [我们知道slice 和 map 包含指向底层数据的指针](#%E6%88%91%E4%BB%AC%E7%9F%A5%E9%81%93slice-%E5%92%8C-map-%E5%8C%85%E5%90%AB%E6%8C%87%E5%90%91%E5%BA%95%E5%B1%82%E6%95%B0%E6%8D%AE%E7%9A%84%E6%8C%87%E9%92%88)
  - [什么是 unsafe](#%E4%BB%80%E4%B9%88%E6%98%AF-unsafe)
  - [为什么会有unsafe](#%E4%B8%BA%E4%BB%80%E4%B9%88%E4%BC%9A%E6%9C%89unsafe)
  - [unsafe实现原理](#unsafe%E5%AE%9E%E7%8E%B0%E5%8E%9F%E7%90%86)
  - [unsafe.Pointer && uintptr类型](#unsafepointer--uintptr%E7%B1%BB%E5%9E%8B)
    - [unsafe.Pointer](#unsafepointer)
    - [uintptr](#uintptr)
      - [下面的这段代码是错误的](#%E4%B8%8B%E9%9D%A2%E7%9A%84%E8%BF%99%E6%AE%B5%E4%BB%A3%E7%A0%81%E6%98%AF%E9%94%99%E8%AF%AF%E7%9A%84)
  - [总结](#%E6%80%BB%E7%BB%93)
    - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# unsafe

## 前言

最近关注了一个大佬的文章，文章写的非常好，大家可以去关注下。  
微信公众号【码农桃花源】  

## 指针类型

首先我们先来了解下，`GO`里面的指针类型。

为什么需要指针类型呢？参考文献 `go101.org` 里举了这样一个例子：

````
func double(x int) {
	fmt.Println(x)
	x += x
	fmt.Println(x)
}

func main() {
	var a = 3
	double(a)
	fmt.Println(a)

}
````

`double`函数的作用是将3翻倍，但是实际上却没有做到，为什么呢？  

因为go语言的函数操作都是值传递。`double`函数里面的`x`只是`a`的一个拷贝，在函数内部对x的操作不能反馈到实参`a`。  

其实在实际的编写代码的过程中我们会使用一个指针进行解决。  

````go
func double1(x *int) {
	*x += *x
	x = nil
}

func main() {
	var a = 3
	double1(&a)

	fmt.Println(a)

	p := &a

	double1(p)

	fmt.Println(*p)

}
````

其中有一个操作

````go
x=nil
````

这个操作没有对我们的结果产生丝毫的影响。其实也是很好理解的，因为我们知道go里面的函数中使用的都是值传递`x=nil`，只是对`&a`的一个拷贝。  

### 我们知道slice 和 map 包含指向底层数据的指针

我们对它们的操作是会影响到，原参数的值。  

````go
func change(sl []int64) {
	sl[0] = 2
}

func main() {

	var sl = make([]int64, 2)
	change(sl)
	fmt.Println(sl)  // [2 0]
}
````

我们而已看到输出的值已经是`[2 0]`

这时候我们可以使用一个copy来操作

````go
func change(sl []int64) {
	sl[0] = 2
}

func changeNo(sl []int64) {
	s2 := make([]int64, 2)
	copy(sl, s2)
	s2[0] = 2
}

func main() {

	var sl = make([]int64, 2)
	change(sl)
	fmt.Println(sl)

	changeNo(sl)
	fmt.Println(sl)
}
````

限制一：GO里面的指针不能进行数学的运算

````go
错误
a := 5
p := &a

p++
p = &a + 3
````

限制二：不同类型的指针不能互相转换

````go
错误的
func main(){
   a:=int(100)
   var f *float64

    f=&a
}
````
限制三：不同类型的指针不能使用==或!=比较。

限制四：不能类型的指针变量不能相互赋值。

## 什么是 unsafe

前面讨论的指针是类型安全的，但它有很多的限制。go还有非类型安全的指针，就是`unsafe`包提供的`unsafe.Pointer`。在某些情况下，它会使代码更高效，当然，也更危险。  

`unsafe` 包用于Go编译器，在编译阶段使用。从名字就可以看出来，它是不安全的，官方并不建议使用。它可以绕过Go语言的类型系统，直接操作内存。  

## 为什么会有unsafe

Go 语言类型系统是为了安全和效率设计的,有时,安全会导致效率低下。有了`unsafe`包，高阶的程序员就可以利用它绕过类型系统的低效。因此，它就有了存在的意义，阅读Go源码，会发现有大量使用`unsafe`包的例子。  

## unsafe实现原理

我们来看源码：

````go
type ArbitraryType int
type Pointer *ArbitraryType
````

从命名来看，`Arbitrary`是任意的意思，也就是说`Pointer`可以指向任意类型，实际上它类似于C语言里的`void*`。  

`unsafe`包还有其他三个函数：  

````go
func Sizeof(x ArbitraryType) uintptr
func Offsetof(x ArbitraryType) uintptr
func Alignof(x ArbitraryType) uintptr
````

`size`返回类型x所占的字节数，单不包含x所指向的内容的大小。例如，对于一个指针，函数返回的大小为8字节（64位机器上），一个slice的大小则为`slice header`的大小。  

`offsetof`返回结构体在内存中的位置离结构体起始处的字节数，所传参数必须是结构体的成员。  

`Alignof` 返回m,m是指当类型进行内存对齐时，它分配到的内存地址能整除m。  

上面三个函数的返回结果都是`uintptr`类型，这和`unsafe.Pointer`可以相互转换。三个函数都是在编译期间执行，它们的结果可以直接赋给`const`型变量。另外，因为三个函数执行的结果和操作系统、编译器相关，所以是不可可移值的。

综上，`unsafe`包提供了2点重要的能力：  

````go
1、任何类型的指针和unsafe.Point可以相互转换。
2、uintptr类型和unsafe.Point可以相互转换
````

<img src="/img/unsafe_1.png" alt="gc" align="unsafe" />

`pointer`不能直接进行数学运算，但可以把它转换成`uintptr`,对uintptr类型进行数学运算，在转换成pointer类型。  

````go
// uintptr 是一个整数类型，它足够大，可以存储
type uintptr uintptr
````

还有一点需要注意的是，`uintptr`并没有指针的含义，意思是`uintptr`所指向的对象会被gc给回收掉。而`unsafe.Pointer`有指针语义，可以保护它所指向的对象在“有用”的时候不会被垃圾回收。  

## unsafe.Pointer && uintptr类型

### unsafe.Pointer

这个类型比较重要，它是实现定位欲读写的内存的基础。官方文档对该类型有四个重要描述：

````go
（1）任何类型的指针都可以被转化为Pointer
（2）Pointer可以被转化为任何类型的指针
（3）uintptr可以被转化为Pointer
（4）Pointer可以被转化为uintptr
````

大多数指针类型都会写成T，表示是“一个指向T类型变量的指针”。`unsafe.Pointer`是特别定义的一种指针类型，它可以包含任何类型变量的地址。当然，我们不可以直接通过*p来获取`unsafe.Pointer`指针指向的真是变量的值，因为我们并不知道变量的具体类型。和人普通指针一样，`unsafe.Pointer`指针是可以比较的，并且支持和nil常量比较判断是否为空指针。  

****
一个普通的的T类型指针可以被转换成`unsafe.Pointer`类型指针，并且一个`unsafe.Pointer`类型指针也可以
被转换成普通类型的指针，被转换回普通的指针类型并不需要和原始的T类型相同。
****

通过将float64类型指针转化为uint64类型指针，我们可以查看一个浮点数变量的位模式。

````go
func Float64bits(f float64) uint64 {
	fmt.Println(reflect.TypeOf(unsafe.Pointer(&f)))            //unsafe.Pointer
	fmt.Println(reflect.TypeOf((*uint64)(unsafe.Pointer(&f)))) //*uint64
	return *(*uint64)(unsafe.Pointer(&f))
}

func main() {
	fmt.Printf("%#016x\n", Float64bits(1.0)) // "0x3ff0000000000000"
}
````

再看一个例子

````go
func main() {
	v1 := uint(12)
	v2 := int(12)

	fmt.Println(reflect.TypeOf(v1)) //uint
	fmt.Println(reflect.TypeOf(v2)) //int

	fmt.Println(reflect.TypeOf(&v1)) //*uint
	fmt.Println(reflect.TypeOf(&v2)) //*int

	p := &v1

	//两个变量的类型不同,不能赋值
	//p = &v2 //cannot use &v2 (type *int) as type *uint in assignment

	fmt.Println(reflect.TypeOf(p)) // *unit
}
````

````
当再次把 v2 的指针赋值给p时，会发生错误cannot use &v2 (type *int) as type *uint in assignment，也就是说类型不同，一个是*int，一个是*uint。
````

可以使用unsafe.Pointer进行转换，如下，

````go
func main() {

	v1 := uint(12)
	v2 := int(13)

	fmt.Println(reflect.TypeOf(v1)) //uint
	fmt.Println(reflect.TypeOf(v2)) //int

	fmt.Println(reflect.TypeOf(&v1)) //*uint
	fmt.Println(reflect.TypeOf(&v2)) //*int

	p := &v1

	p = (*uint)(unsafe.Pointer(&v2)) //使用unsafe.Pointer进行类型的转换

	fmt.Println(reflect.TypeOf(p)) // *unit
	fmt.Println(*p)                //13
}
````

### uintptr

````go
// uintptr is an integer type that is large enough to hold the bit pattern of
// any pointer.
type uintptr uintptr
````

`uintptr`是golang的内置类型，是能存储指针的整型，在64位平台上底层的数据类型是，

````go
typedef unsigned long long int  uint64;
typedef uint64          uintptr;
````

一个`unsafe.Pointer`指针也可以被转化成`uintptr`类型，然后保存到指针类型数值变量中（注：这只是和当前指针相同的一个数字值，并不是一个指针），然后用以做必要的指针数值运算。（uintptr是一个无符号的整型数，足以保存一个地址）这种转换虽然是可逆的，但是将`uintptr`转为`unsafe.Pointer`指针可能破坏类型系统，因为并不是所有的数字都是有效的内存地址。 
 

许多将`unsafe.Pointer`指针转化成原生数字，然后再转换成`unsafe.Pointer`类型指针的操作也是不安全的。比如下面的例子需要将变量x的地址加上b字段地址偏移量转化为`*int16`类型指针，然后通过该指针更新`x.b`：  

````go
func main() {

	var x struct {
		a bool
		b int16
		c []int
	}

	/**
	unsafe.Offsetof 函数的参数必须是一个字段 x.f, 然后返回 f 字段相对于 x 起始地址的偏移量, 包括可能的空洞.
	*/

	/**
	uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)
	指针的运算
	*/
	// 和 pb := &x.b 等价
	pb := (*int16)(unsafe.Pointer(uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)))
	*pb = 42
	fmt.Println(x.b) // "42"
}
````

上面的写法尽管很繁琐，但在这里并不是一件坏事，因为这些功能应该很谨慎地使用。不要试图引入一个`uintptr`类型的临时变量，因为它可能会破坏代码的安全性（注：这是真正可以体会`unsafe`包为何不安全的例子）。  

#### 下面的这段代码是错误的

````go
// NOTE: subtly incorrect!
tmp := uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)
pb := (*int16)(unsafe.Pointer(tmp))
*pb = 42
````
产生错误的原因很微妙。有时候垃圾回收器会移动一些变量以降低内存碎片等问题。这类垃圾回收器被称为移动GC。当一个变量被移动，所有的保存改变量旧地址的指针必须同时被更新为变量移动后的地址。从垃圾收集器的角度看，一个`unsafe.Pointer`是一个指向变量的指针，因此当变量被移动是对应的指针也必须被更新；但是`uintptr`类型的临时变量只是一个普通的数字，所以其值不应该被改变。上面错误的代码因引入一个非指针的临时变量`temp`，导致垃圾收集器无法正确识别这个是一个指向变量x的指针。当第二个语句执行是，变量X可能被转移，这时候临时变量tmp也就是不再是现在`&x.b`地址。第三个指向之前无效地址空间的赋值将摧毁整个系统。  

## 总结

`unsafe`包绕过了GO的类型系统，达到直接操作内存的目的，使用它是有一定风险的。但是在某些场景下，使用`unsafe`包函数会提升代码的效率，GO源码中也是大量使用`unsafe`包。  

`unsafe`包定义了`Pointer`和三个函数：  

````go
func Sizeof(x ArbitraryType) uintptr
func Offsetof(x ArbitraryType) uintptr
func Alignof(x ArbitraryType) uintptr
````

通过三个函数可以获取变量的大小，偏移，对齐等信息。  

`uintptr`可以和`unsafe.Pointer`进行相互的转换，`uintptr`可以进行数学运算。这样，通过`uintptr`和`unsafe.Pointer`的结合就解决了Go指针不能进行数学运算的限制。  

通过`unsafe`相关函数，可以获取结构体私有成员的地址，进而对其做进一步的读写操作，突破Go的类型安全限制。  

### 参考
【深度解密Go语言之unsafe】 https://mp.weixin.qq.com/s/OO-kwB4Fp_FnCaNXwGJoEw    
【Go之unsafe.Pointer && uintptr类型】 https://my.oschina.net/xinxingegeya/blog/729673
















































