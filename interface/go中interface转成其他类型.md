- [go中interface转换成原来的类型](#go%E4%B8%ADinterface%E8%BD%AC%E6%8D%A2%E6%88%90%E5%8E%9F%E6%9D%A5%E7%9A%84%E7%B1%BB%E5%9E%8B)
  - [首先了解下interface](#%E9%A6%96%E5%85%88%E4%BA%86%E8%A7%A3%E4%B8%8Binterface)
    - [什么是interface?](#%E4%BB%80%E4%B9%88%E6%98%AFinterface)
    - [如何判断interface变量存储的是哪种类型](#%E5%A6%82%E4%BD%95%E5%88%A4%E6%96%ADinterface%E5%8F%98%E9%87%8F%E5%AD%98%E5%82%A8%E7%9A%84%E6%98%AF%E5%93%AA%E7%A7%8D%E7%B1%BB%E5%9E%8B)
      - [fmt](#fmt)
      - [反射](#%E5%8F%8D%E5%B0%84)
      - [断言](#%E6%96%AD%E8%A8%80)
  - [参考](#%E5%8F%82%E8%80%83)

## go中interface转换成原来的类型

### 首先了解下interface

#### 什么是interface?

首先 `interface` 是一种类型，从它的定义可以看出来用了 `type` 关键字，更准确的说 `interface` 是一种具有一组方法的类型，这些方法定义了 `interface` 的行为。

````go
type I interface {
    Get() int
}
````
`interface`是一组`method`的集合，是`duck-type programming`的一种体现（不关心属性（数据），只关心行为（方法））。我们可以自己定义`interface`类型的`struct`,并提供方法。

````go
type MyInterface interface{
    Print()
}

func TestFunc(x MyInterface) {}
type MyStruct struct {}
func (me MyStruct) Print() {}

func main() {
    var me MyStruct
    TestFunc(me)
}
````

`go` 允许不带任何方法的 `interface` ，这种类型的 `interface` 叫 `empty interface`。  

如果一个类型实现了一个 `interface` 中所有方法，必须是所有的方法，我们说类型实现了该 `interface`，所以所有类型都实现了 `empty interface`，因为任何一种类型至少实现了 0 个方法。`go` 没有显式的关键字用来实现 `interface`，只需要实现 `interface` 包含的方法即可。  

`interface`还可以作为返回值使用。  

#### 如何判断interface变量存储的是哪种类型

日常中使用`interface`,有时候需要判断原来是什么类型的值转成了`interface`。一般有以下几种方式：

##### fmt

````go
import "fmt"
func main() {
    v := "hello world"
    fmt.Println(typeof(v))
}
func typeof(v interface{}) string {
    return fmt.Sprintf("%T", v)
}
````

##### 反射

````
import (
    "reflect"
    "fmt"
)
func main() {
    v := "hello world"
    fmt.Println(typeof(v))
}
func typeof(v interface{}) string {
    return reflect.TypeOf(v).String()
}
````  

##### 断言

`Go`语言里面有一个语法，可以直接判断是否是该类型的变量： `value, ok = element.(T)`，这里`value`就是变量的值，`ok`是一个`bool`类型，`element`是`interface`变量，`T`是断言的类型。  

如果`element`里面确实存储了`T`类型的数值，那么ok返回`true`，否则返回`false`。  

让我们通过一个例子来更加深入的理解。  

````go
value, ok := v.(string)

if ok {
    return value
}
````

类型不确定可以配合`switch`:

````go
func main() {
    v := "hello world"
    fmt.Println(typeof(v))
}
func typeof(v interface{}) string {
    switch t := v.(type) {
    case int:
        return "int"
    case float64:
        return "float64"
    //... etc
    default:
        _ = t
        return "unknown"
    }
}
````

对于fmt也是用了反射的,同时里面也用到了断言:

````go
func (p *pp) printArg(arg interface{}, verb rune) {
	p.arg = arg
	p.value = reflect.Value{}

	if arg == nil {
		switch verb {
		case 'T', 'v':
			p.fmt.padString(nilAngleString)
		default:
			p.badVerb(verb)
		}
		return
	}

	// Special processing considerations.
	// %T (the value's type) and %p (its address) are special; we always do them first.
	switch verb {
	case 'T':
		p.fmt.fmtS(reflect.TypeOf(arg).String())
		return
	case 'p':
		p.fmtPointer(reflect.ValueOf(arg), 'p')
		return
	}

	// Some types can be done without reflection.
	switch f := arg.(type) {
	case bool:
		p.fmtBool(f, verb)
	case float32:
		p.fmtFloat(float64(f), 32, verb)
	case float64:
		p.fmtFloat(f, 64, verb)
	case complex64:
		p.fmtComplex(complex128(f), 64, verb)
	case complex128:
		p.fmtComplex(f, 128, verb)
	case int:
		p.fmtInteger(uint64(f), signed, verb)
	case int8:
		p.fmtInteger(uint64(f), signed, verb)
	case int16:
		p.fmtInteger(uint64(f), signed, verb)
	case int32:
		p.fmtInteger(uint64(f), signed, verb)
	case int64:
		p.fmtInteger(uint64(f), signed, verb)
	case uint:
		p.fmtInteger(uint64(f), unsigned, verb)
	case uint8:
		p.fmtInteger(uint64(f), unsigned, verb)
	case uint16:
		p.fmtInteger(uint64(f), unsigned, verb)
	case uint32:
		p.fmtInteger(uint64(f), unsigned, verb)
	case uint64:
		p.fmtInteger(f, unsigned, verb)
	case uintptr:
		p.fmtInteger(uint64(f), unsigned, verb)
	case string:
		p.fmtString(f, verb)
	case []byte:
		p.fmtBytes(f, verb, "[]byte")
	case reflect.Value:
		// Handle extractable values with special methods
		// since printValue does not handle them at depth 0.
		if f.IsValid() && f.CanInterface() {
			p.arg = f.Interface()
			if p.handleMethods(verb) {
				return
			}
		}
		p.printValue(f, verb, 0)
	default:
		// If the type is not simple, it might have methods.
		if !p.handleMethods(verb) {
			// Need to use reflection, since the type had no
			// interface methods that could be used for formatting.
			p.printValue(reflect.ValueOf(f), verb, 0)
		}
	}
}
````

下面来简单探究下反射是如何判断`interface`  

````go
// TypeOf returns the reflection Type that represents the dynamic type of i.
// If i is a nil interface value, TypeOf returns nil.
func TypeOf(i interface{}) Type {
	eface := *(*emptyInterface)(unsafe.Pointer(&i))
	return toType(eface.typ)
}
````

`eface := *(*emptyInterface)(unsafe.Pointer(&i))`用到了一个`emptyInterface`，我们来看下这个结构的信息:

````go
// emptyInterface is the header for an interface{} value.
type emptyInterface struct {
	typ  *rtype
	word unsafe.Pointer
}
````

其中`typ`指向一个`rtype`实体， 它表示`interface`的类型以及赋给这个`interface`的实体类型。`word`则指向`interface`具体的值，一般而言是一个指向堆内存的指针。  

`TypeOf`看到的是空接口`interface{}`，它将变量的地址转换为空接口，然后将得到的`rtype`转为`Type`接口返回。需要注意，当调用`reflect.TypeOf`的之前，已经发生了一次隐式的类型转换，即将具体类型的向空接口转换。这个过程比较简单，只要拷贝`typ *rtype`和`word unsafe.Pointer`就可以了。  

### 来看下interface的底层源码

我的go版本是`go version go1.13.7`  

`iface`和`eface`都是`Go`中描述接口的底层结构体，区别在于`iface`描述的接口包含方法，而`eface`则是不包含任何方法的空接口：`interface{}`。  

### eface

代码在runtime/runtime2.go:  

````go
type eface struct {
	_type *_type
	data  unsafe.Pointer
}
````

`eface`有两个字段，`_type`指向对象的类型信息，`data`数据指针。指针指向的数据地址，一般是在堆上的。  

我们来看下`_type`  

````
// src/rumtime/runtime2.go
type _type struct {
    size       uintptr     // 类型的大小
    ptrdata    uintptr     // size of memory prefix holding all pointers
    hash       uint32      // 类型的Hash值
    tflag      tflag       // 类型的Tags 
    align      uint8       // 结构体内对齐
    fieldalign uint8       // 结构体作为field时的对齐
    kind       uint8       // 类型编号 定义于runtime/typekind.go
    alg        *typeAlg    // 类型元方法 存储hash和equal两个操作。map key便使用key的_type.alg.hash(k)获取hash值
    gcdata    *byte        // GC相关信息
    str       nameOff      // 类型名字的偏移    
    ptrToThis typeOff    
}
````

`_type`是`go`中类型的公共描述，里面包含`GC`，反射等需要的细节，它决定`data`应该如何解释和操作。对于不同的数据类型它的描述信息是不一样的，在`_type`的基础之上配合一些额外的描述信息，来进行区分。  

````go
// src/runtime/type.go
// ptrType represents a pointer type.
type ptrType struct {
   typ     _type   // 指针类型 
   elem  *_type // 指针所指向的元素类型
}
type chantype struct {
    typ  _type        // channel类型
    elem *_type     // channel元素类型
    dir  uintptr
}
type maptype struct {
    typ           _type
    key           *_type
    elem          *_type
    bucket        *_type // internal type representing a hash bucket
    hmap          *_type // internal type representing a hmap
    keysize       uint8  // size of key slot
    indirectkey   bool   // store ptr to key instead of key itself
    valuesize     uint8  // size of value slot
    indirectvalue bool   // store ptr to value instead of value itself
    bucketsize    uint16 // size of bucket
    reflexivekey  bool   // true if k==k for all keys
    needkeyupdate bool   // true if we need to update key on an overwrite
}
````

这些类型信息的第一个字段都是`_type`(类型本身的信息)，接下来是一堆类型需要的其它详细信息(如子类型信息)，这样在进行类型相关操作时，可通过一个字`(typ *_type)`即可表述所有类型，然后再通过`_type.kind`可解析出其具体类型，最后通过地址转换即可得到类型完整的”_type树”，参考`reflect.Type.Elem()`函数:

````go
// reflect/type.go
// reflect.rtype结构体定义和runtime._type一致  type.kind定义也一致(为了分包而重复定义)
// Elem()获取rtype中的元素类型，只针对复合类型(Array, Chan, Map, Ptr, Slice)有效
func (t *rtype) Elem() Type {
   switch t.Kind() {
   case Array:
      tt := (*arrayType)(unsafe.Pointer(t))
      return toType(tt.elem)
   case Chan:
      tt := (*chanType)(unsafe.Pointer(t))
      return toType(tt.elem)
   case Map:
      // 对Map来讲，Elem()得到的是其Value类型
      // 可通过rtype.Key()得到Key类型
      tt := (*mapType)(unsafe.Pointer(t))
      return toType(tt.elem)
   case Ptr:
      tt := (*ptrType)(unsafe.Pointer(t))
      return toType(tt.elem)
   case Slice:
      tt := (*sliceType)(unsafe.Pointer(t))
      return toType(tt.elem)
   }
   panic("reflect: Elem of invalid type")
}
````


#### iface  

表示的是非空的接口:  

````go
type iface struct {
	tab  *itab
	data unsafe.Pointer
}

// layout of Itab known to compilers
// allocated in non-garbage-collected memory
// Needs to be in sync with
// ../cmd/compile/internal/gc/reflect.go:/^func.dumptypestructs.
type itab struct {
	inter *interfacetype  // 接口定义的类型信息
	_type *_type          // 接口实际指向值的类型信息
	hash  uint32 // copy of _type.hash. Used for type switches.
	_     [4]byte
	fun   [1]uintptr     // 接口方法实现列表，即函数地址列表，按字典序排序 variable sized
}
// runtime/type.go
// 非空接口类型，接口定义，包路径等。
type interfacetype struct {
   typ     _type
   pkgpath name
   mhdr    []imethod      // 接口方法声明列表，按字典序排序
}

// 接口的方法声明 
type imethod struct {
   name nameOff              // 方法名
   ityp typeOff              // 描述方法参数返回值等细节
}
````

`iface`同样也是有两个指针，`tab`指向一个`itab`实体， 它表示接口的类型以及赋给这个接口的实体类型。`data`则指向接口具体的值，一般而言是一个指向堆内存的指针。  

`fun`表示`interface`中`method`的具体实现。比如`interfacetype`包含了两个`method`分别是`A`和`B`。但是有一点很奇怪，这个`fun`是长度为1的`uintptr`数组，那么是怎么表示多个的呢？  
其实上面源码的注释已经能给到我们答案了，`variable sized`，这是个是可变大小的。go中的`uintptr`一般用来存放指针的值，那这里对应的就是函数指针的值（也就是函数的调用地址）。但是这里的 fun 是一个长度为 1 的 uintptr 数组。我们看一下 runtime 包的 additab 函数。

````go
func additab(m *itab, locked, canfail bool) {
    ...
    *(*unsafe.Pointer)(add(unsafe.Pointer(&m.fun[0]), uintptr(k)*sys.PtrSize)) = ifn
    ...
}
````

也就是在`fun[0]`后面一次写入其他`method`对应的函数指针。





`iface`是


### 参考
【理解 Go interface 的 5 个关键点】https://sanyuesha.com/2017/07/22/how-to-understand-go-interface/  
【深入理解 Go Interface】https://zhuanlan.zhihu.com/p/32926119   
【GO如何支持泛型】https://zhuanlan.zhihu.com/p/74525591  
【Golang面向对象编程】https://code.tutsplus.com/zh-hans/tutorials/lets-go-object-oriented-programming-in-golang--cms-26540  
【深度解密Go语言之关于 interface 的10个问题】https://www.cnblogs.com/qcrao-2018/p/10766091.html  
【golang如何获取变量的类型：反射，类型断言】https://ieevee.com/tech/2017/07/29/go-type.html    
【Go接口详解】https://zhuanlan.zhihu.com/p/27055513  
