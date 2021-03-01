<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [unsafe](#unsafe)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [什么是unsafe,为什么需要unsafe](#%E4%BB%80%E4%B9%88%E6%98%AFunsafe%E4%B8%BA%E4%BB%80%E4%B9%88%E9%9C%80%E8%A6%81unsafe)
  - [unsafe实现原理](#unsafe%E5%AE%9E%E7%8E%B0%E5%8E%9F%E7%90%86)
  - [unsafe.Pointer && uintptr类型](#unsafepointer--uintptr%E7%B1%BB%E5%9E%8B)
    - [unsafe.Pointer](#unsafepointer)
    - [uintptr](#uintptr)
  - [总结](#%E6%80%BB%E7%BB%93)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## unsafe

### 前言

在阅读go源码的时候，发现很多地方使用了`unsafe.Pointer`来处理指针类型的转换，这次来深入的探究下。   

### 什么是unsafe,为什么需要unsafe

Go语言在设计的时候，为了编写方便、效率高以及降低复杂度，被设计成为一门强类型的静态语言。强类型意味着一旦定义了，它的类型就不能改变了；静态意味着类型检查在运行前就做了。  

例如go中的指针存在的使用限制  

1、go指针不支持算术运算  

2、一个指针类型的值不能被随意转换为另一个指针类型  

3、一个指针值不能和其它任一指针类型的值进行比较  

4、一个指针值不能被赋值给其它任意类型的指针值  

`unsafe`可以打破这些限制  

> Package unsafe contains operations that step around the type safety of Go programs.

`unsafe`可以绕过go类型的安全检查，直接操控内存，我们可以写出高效的代码。  

但是正如他的名字一样`unsafe`，不安全。我们应该尽可能少的使用它，比如内存的操纵，这是绕过Go本身设计的安全机制的，不当的操作，可能会破坏一块内存，而且这种问题非常不好定位。   

### unsafe实现原理

unsafe主要包含下面三个函数  

```go
// Arbitrary 是任意的意思，也就是说 Pointer 可以指向任意类型
type ArbitraryType int
type Pointer *ArbitraryType

// 返回类型 x 所占据的字节数，但不包含 x 所指向的内容的大小。
// 例如，对于一个指针，函数返回的大小为 8 字节（64位机上），一个 slice 的大小则为 slice header 的大小。
func Sizeof(x ArbitraryType) uintptr

// 返回结构体中某个field的偏移量
// 所传参数必须是结构体的成员
func Offsetof(x ArbitraryType) uintptr

// 对应参数的内存对齐系数
func Alignof(x ArbitraryType) uintptr
```

什么是内存对齐，可参考[什么是内存对齐，go中内存对齐分析](https://www.cnblogs.com/ricklz/p/14455135.html)  

来个简单的例子看下  

```go
type People struct {
	age  uint8
	name string
}

func main() {
	h := People{
		30,
		"xiaobai",
	}

	i := unsafe.Sizeof(h)
	j := unsafe.Alignof(h)
	k := unsafe.Offsetof(h.name)
	fmt.Println("字节大小：", i)
	fmt.Println("对齐系数：", j)
	fmt.Println("偏移量：", k)
	fmt.Printf("直接获取地址：%p\n", &h)

	var p unsafe.Pointer
	p = unsafe.Pointer(&h)
	fmt.Println("使用unsafe获取地址：", p)
}
```

简单看下输出  

```
字节大小： 24
对齐系数： 8
偏移量： 8
直接获取地址：0xc00007c020
使用unsafe获取地址： 0xc00007c020
```

### unsafe.Pointer && uintptr类型

#### unsafe.Pointer

这个类型比较重要，它是实现定位欲读写的内存的基础。官方文档对该类型有四个重要描述：

- 1、任何类型的指针都可以被转化为Pointer  

- 2、Pointer可以被转化为任何类型的指针  

- 3、uintptr可以被转化为Pointer  

- 4、Pointer可以被转化为uintptr     

大多数指针类型都会写成T，表示是“一个指向T类型变量的指针”。`unsafe.Pointer`是特别定义的一种指针类型，它可以包含任何类型变量的地址。当然，我们不可以直接通过*p来获取`unsafe.Pointer`指针指向的真是变量的值，因为我们并不知道变量的具体类型。和人普通指针一样，`unsafe.Pointer`指针是可以比较的，并且支持和nil常量比较判断是否为空指针。  

<img src="../img/unsafe_2.png" width = "742.772" height = "87" alt="unsafe" align=center />

**一个普通的的T类型指针可以被转换成`unsafe.Pointer`类型指针，并且一个`unsafe.Pointer`类型指针也可以被转换成普通类型的指针，被转换回普通的指针类型并不需要和原始的T类型相同。**

举几个栗子来分析下  

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

当再次把v2的指针赋值给p时，会发生错误`cannot use &v2 (type *int) as type *uint in assignment`，也就是说类型不同，一个是`*int`，一个是`*uint`。  

可以使用`unsafe.Pointer`进行转换，如下，

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

#### uintptr

````go
// uintptr is an integer type that is large enough to hold the bit pattern of
// any pointer.
type uintptr uintptr
````

uintptr 的底层实现如下，在`$GOROOT/src/pkg/runtime/runtime.h`中：  

```go
#ifdef _64BIT
typedef uint64          uintptr;
typedef int64           intptr;
typedef int64           intgo; // Go's int
typedef uint64          uintgo; // Go's uint
#else
typedef uint32          uintptr;
typedef int32           intptr;
typedef int32           intgo; // Go's int
typedef uint32          uintgo; // Go's uint
#endif
```

`uintptr`和`intptr`是无符号和有符号的指针类型，并且确保在64位平台上是8个字节，在32位平台上是4个字节，`uintptr`主要用于golang中的指针运算。  

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

`下面的这段代码是错误的`

````go
// NOTE: subtly incorrect!
tmp := uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)
pb := (*int16)(unsafe.Pointer(tmp))
*pb = 42
````

产生错误的原因很微妙。有时候垃圾回收器会移动一些变量以降低内存碎片等问题。这类垃圾回收器被称为移动GC。当一个变量被移动，所有的保存改变量旧地址的指针必须同时被更新为变量移动后的地址。从垃圾收集器的角度看，一个`unsafe.Pointer`是一个指向变量的指针，因此当变量被移动是对应的指针也必须被更新；但是`uintptr`类型的临时变量只是一个普通的数字，所以其值不应该被改变。上面错误的代码因引入一个非指针的临时变量`temp`，导致垃圾收集器无法正确识别这个是一个指向变量x的指针。当第二个语句执行是，变量X可能被转移，这时候临时变量tmp也就是不再是现在`&x.b`地址。第三个指向之前无效地址空间的赋值将摧毁整个系统。  

uintptr 和 unsafe.Pointer 的互相转换  

```go
func main() {
	a := [4]int{0, 1, 2, 3}
	p := &a[1] // 内存地址
	p1 := unsafe.Pointer(p) 
	p2 := uintptr(p1)
	p3 := unsafe.Pointer(p2)
	fmt.Println(p1) // 0xc420014208
	fmt.Println(p2) // 842350543368
	fmt.Println(p3) // 0xc420014208
}
```

### 总结

1、`unsafe`包绕过了GO的类型系统，达到直接操作内存的目的，使用它是有一定风险的。但是在某些场景下，使用`unsafe`包函数会提升代码的效率，GO源码中也是大量使用`unsafe`包。  

2、`uintptr`可以和`unsafe.Pointer`进行相互的转换，`uintptr`可以进行数学运算。这样，通过`uintptr`和`unsafe.Pointer`的结合就解决了Go指针不能进行数学运算的限制。  

3、通过`unsafe`相关函数，可以获取结构体私有成员的地址，进而对其做进一步的读写操作，突破Go的类型安全限制。  

4、`uintptr`并没有指针的含义，意思是`uintptr`所指向的对象会被gc给回收掉。而`unsafe.Pointer`有指针语义，可以保护它所指向的对象在“有用”的时候不会被垃圾回收。  

### 参考

【Go之unsafe.Pointer && uintptr类型】 https://my.oschina.net/xinxingegeya/blog/729673   
【Go unsafe包】https://my.oschina.net/xinxingegeya/blog/841058  
【unsafe包】https://wizardforcel.gitbooks.io/go42/content/content/42_28_unsafe.html  
【非类型安全指针】https://gfw.go101.org/article/unsafe.html  
【Go unsafe 包的使用】https://segmentfault.com/a/1190000021625500   
【Go unsafe Pointer】https://www.flysnow.org/2017/07/06/go-in-action-unsafe-pointer.html   
【指针】https://gfw.go101.org/article/pointer.html  
【深度解密Go语言之unsafe】 https://mp.weixin.qq.com/s/OO-kwB4Fp_FnCaNXwGJoEw     
【golang中的unsafe详解】https://studygolang.com/articles/18436   
