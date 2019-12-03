# unsafe


## 指针类型
首先我们先来了解下，GO里面的指针类型。

为什么需要指针类型呢？参考文献 go101.org 里举了这样一个例子：

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

double函数的作用是将3翻倍，但是实际上却没有做到，为什么呢？
因为go语言的函数操作都是值传递。double函数里面的x只是a的一个拷贝，
在函数内部对x的操作不能反馈到实参a。

其实在实际的编写代码的过程中我们会使用一个指针进行解决。

````

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
````
x=nil
````

这个操作没有对我们的结果产生丝毫的影响。
其实也是很好理解的，因为我们知道go里面的函数中使用的都是值传递
x=nil，只是对&a的一个拷贝。

### 我们知道slice 和 map 包含指向底层数据的指针。
我们对它们的操作是会影响到，原参数的值。

````
func change(sl []int64) {
	sl[0] = 2
}

func main() {

	var sl = make([]int64, 2)
	change(sl)
	fmt.Println(sl)  // [2 0]
}
````
我们而已看到输出的值已经是[2 0]

这时候我们可以使用一个copy来操作
````
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

````
错误
a := 5
p := &a

p++
p = &a + 3
````
限制二：不同类型的指针不能互相转换

````
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

前面讨论的指针是类型安全的，但它有很多的限制。go还有非类型安全的指针，就是unsafe
包提供的 unsafe.Pointer。在某些情况下，它会使代码更高效，当然，也更危险。

unsafe 包用于 Go 编译器，在编译阶段使用。从名字就可以看出来，它是不安全的，官方并不建议使用。
它可以绕过 Go 语言的类型系统，直接操作内存。


## 为什么会有unsafe

Go 语言类型系统是为了安全和效率设计的，有时，安全会导致效率低下。有了 unsafe 包，高阶的程序员
就可以利用它绕过类型系统的低效。因此，它就有了存在的意义，阅读 Go 源码，会发现有大量使用 unsafe 包的例子。

## unsafe实现原理

我们来看源码：

````
type ArbitraryType int
type Pointer *ArbitraryType
````

从命名来看， Arbitrary 是任意的意思，也就是说 Pointer 可以指向任意类型，实际上它类似于 C 语言里的 void*。

unsafe包还有其他三个函数：
````
func Sizeof(x ArbitraryType) uintptr
func Offsetof(x ArbitraryType) uintptr
func Alignof(x ArbitraryType) uintptr
````
size返回类型x所占的字节数，单不包含x所指向的内容的大小。例如，对于一个指针，函数返回
的大小为8字节（64位机器上），一个slice的大小则为slice header的大小。

offsetof返回结构体在内存中的位置离结构体起始处的字节数，所传参数必须是结构体的成员。

Alignof 返回 m，m 是指当类型进行内存对齐时，它分配到的内存地址能整除 m。

上面三个函数的返回结果都是uintptr类型，这和 unsafe.Pointer 可以相互转换。三个函数都是在编译期间执
行，它们的结果可以直接赋给 const型变量。另外，因为三个函数执行的结果和操作系统、编译器相关，所以是不可
可移值的。

综上，unsafe包提供了2点重要的能力：

````
1、任何类型的指针和unsafe.Point可以相互转换。
2、uintptr类型和unsafe.Point可以相互转换
````
pointer不能直接进行数学运算，但可以把它转换成uintptr,对uintptr类型进行数学运算，在转换成pointer
类型。












## unsafe.Pointer && uintptr类型

### unsafe.Pointer

这个类型比较重要，它是实现定位欲读写的内存的基础。官方文档对该类型有四个重要描述：

````
（1）任何类型的指针都可以被转化为Pointer
（2）Pointer可以被转化为任何类型的指针
（3）uintptr可以被转化为Pointer
（4）Pointer可以被转化为uintptr
````
大多数指针类型都会写成T，表示是“一个指向T类型变量的指针”。unsafe.Pointer是特别定义的一种指
针类型，它可以包含任何类型变量的地址。当然，我们不可以直接通过*p来获取unsafe.Pointer指针指
向的真是变量的值，因为我们并不知道变量的具体类型。和人普通指针一样，unsafe.Pointer指针是可以
比较的，并且支持和nil常量比较判断是否为空指针。

****
一个普通的的T类型指针可以被转换成unsafe.Pointer类型指针，并且一个unsafe.Pointer类型指针也可以
被转换成普通类型的指针，被转换回普通的指针类型并不需要和原始的T类型相同。
****

通过将float64类型指针转化为uint64类型指针，我们可以查看一个浮点数变量的位模式。

````
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
````


func main() {
``	v1 := uint(12)
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
````
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

````
// uintptr is an integer type that is large enough to hold the bit pattern of
// any pointer.
type uintptr uintptr
````
uintptr是golang的内置类型，是能存储指针的整型，在64位平台上底层的数据类型是，

````
typedef unsigned long long int  uint64;
typedef uint64          uintptr;
````
一个unsafe.Pointer指针也可以被转化成uintptr类型，然后保存到指针类型数值变量中（注：这只是和
当前指针相同的一个数字值，并不是一个指针），然后用以做必要的指针数值运算。（uintptr是一个无符号
的整型数，足以保存一个地址）这种转换虽然是可逆的，但是将uintptr转为unsafe.Pointer指针可能破坏
类型系统，因为并不是所有的数字都是有效的内存地址。

许多将unsafe.Pointer指针转化成原生数字，然后再转换成unsafe.Pointer类型指针的操作也是不安全的
。比如下面的例子需要将变量x的地址加上b字段地址偏移量转化为*int16类型指针，然后通过该指针更新x.b：

````
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
上面的写法尽管很繁琐，但在这里并不是一件坏事，因为这些功能应该很谨慎地使用。不要试图引入一个uintptr类型的临时变
量，因为它可能会破坏代码的安全性（注：这是真正可以体会unsafe包为何不安全的例子）。

#### 下面的这段代码是错误的
````
// NOTE: subtly incorrect!
tmp := uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)
pb := (*int16)(unsafe.Pointer(tmp))
*pb = 42
````
产生错误的原因很微妙。有时候垃圾回收器会移动一些变量以降低内存碎片等问题。这类垃圾回收
器被称为移动GC。当一个变量被移动，所有的保存改变量旧地址的指针必须同时被更新为变量移动
后的地址。从垃圾收集器的角度看，一个unsafe.Pointer是一个指向变量的指针，因此当变量被
移动是对应的指针也必须被更新；但是uintptr类型的临时变量只是一个普通的数字，所以其值
不应该被改变。上面错误的代码因引入一个非指针的临时变量temp，导致垃圾收集器无法正确识别
这个是一个指向变量x的指针。当第二个语句执行是，变量X可能被转移，这时候临时变量tmp也就是
不再是现在&x.b地址。第三个指向之前无效地址空间的赋值将摧毁整个系统。




















### 参考
- 【深度解密Go语言之Slice】 https://mp.weixin.qq.com/s/MTZ0C9zYsNrb8wyIm2D8BA    
- 【Go之unsafe.Pointer && uintptr类型】 https://my.oschina.net/xinxingegeya/blog/729673
















































