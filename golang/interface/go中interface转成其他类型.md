<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [go中interface转换成原来的类型](#go%E4%B8%ADinterface%E8%BD%AC%E6%8D%A2%E6%88%90%E5%8E%9F%E6%9D%A5%E7%9A%84%E7%B1%BB%E5%9E%8B)
  - [首先了解下interface](#%E9%A6%96%E5%85%88%E4%BA%86%E8%A7%A3%E4%B8%8Binterface)
    - [什么是interface?](#%E4%BB%80%E4%B9%88%E6%98%AFinterface)
    - [如何判断interface变量存储的是哪种类型](#%E5%A6%82%E4%BD%95%E5%88%A4%E6%96%ADinterface%E5%8F%98%E9%87%8F%E5%AD%98%E5%82%A8%E7%9A%84%E6%98%AF%E5%93%AA%E7%A7%8D%E7%B1%BB%E5%9E%8B)
      - [fmt](#fmt)
      - [反射](#%E5%8F%8D%E5%B0%84)
      - [断言](#%E6%96%AD%E8%A8%80)
  - [来看下interface的底层源码](#%E6%9D%A5%E7%9C%8B%E4%B8%8Binterface%E7%9A%84%E5%BA%95%E5%B1%82%E6%BA%90%E7%A0%81)
  - [eface](#eface)
  - [iface](#iface)
  - [接口的动态类型和动态值](#%E6%8E%A5%E5%8F%A3%E7%9A%84%E5%8A%A8%E6%80%81%E7%B1%BB%E5%9E%8B%E5%92%8C%E5%8A%A8%E6%80%81%E5%80%BC)
  - [interface如何支持泛型](#interface%E5%A6%82%E4%BD%95%E6%94%AF%E6%8C%81%E6%B3%9B%E5%9E%8B)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

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


### iface  

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
其实上面源码的注释已经能给到我们答案了，`variable sized`，这是个是可变大小的。go中的`uintptr`一般用来存放指针的值，那这里对应的就是函数指针的值（也就是函数的调用地址）。如果有更多的方法，在它之后的内存空间里继续存储。也就是在`fun[0]`后面一次写入其他`method`对应的函数指针。  

接口的类型转换是怎么实现的呢？  

举个例子: 

````go
type coder interface {
	code()
	run()
}

type runner interface {
	run()
}

type Gopher struct {
	language string
}

func (g Gopher) code() {
	return
}

func (g Gopher) run() {
	return
}

func main() {
	var c coder = Gopher{}

	var r runner
	r = c
	fmt.Println(c, r)
}
````

定义了两个` interface: coder` 和 `runner`。定义了一个实体类型 `Gopher`，类型 `Gopher` 实现了两个方法，分别是 `run()` 和 `code()`。`main` 函数里定义了一个接口变量 `c`，绑定了一个 `Gopher` 对象，之后将 `c` 赋值给另外一个接口变量` r` 。赋值成功的原因是 `c` 中包含 `run()` 方法。这样，两个接口变量完成了转换。  

上面的转换调用了下面的函数实现的

````go
func convI2I(inter *interfacetype, i iface) (r iface) {
	tab := i.tab
	if tab == nil {
		return
	}
	if tab.inter == inter {
		r.tab = tab
		r.data = i.data
		return
	}
	r.tab = getitab(inter, tab._type, false)
	r.data = i.data
	return
}
````

关于`conv`的函数定义，其中E代表eface，I代表iface，T代表编译器已知类型，即静态类型。  

`inter`表示转换之后的接口类型，`i`表示转换之前的实体类型接口，`r`表示转换之后的实体类型接口。  
这个函数先做了判断，如果两个转换之前和转换之后的接口类型是一样的，就直接把转换之前的接口信息赋值给r就可以了。如果不一样，就调用`getitab`

````go
func getitab(inter *interfacetype, typ *_type, canfail bool) *itab {
	if len(inter.mhdr) == 0 {
		throw("internal error - misuse of itab")
	}

	// easy case
	if typ.tflag&tflagUncommon == 0 {
		if canfail {
			return nil
		}
		name := inter.typ.nameOff(inter.mhdr[0].name)
		panic(&TypeAssertionError{nil, typ, &inter.typ, name.name()})
	}

	var m *itab

	// First, look in the existing table to see if we can find the itab we need.
	// This is by far the most common case, so do it without locks.
	// Use atomic to ensure we see any previous writes done by the thread
	// that updates the itabTable field (with atomic.Storep in itabAdd).
	t := (*itabTableType)(atomic.Loadp(unsafe.Pointer(&itabTable)))
	if m = t.find(inter, typ); m != nil {
		goto finish
	}

	// Not found.  Grab the lock and try again.
	lock(&itabLock)
	if m = itabTable.find(inter, typ); m != nil {
		unlock(&itabLock)
		goto finish
	}

	// Entry doesn't exist yet. Make a new entry & add it.
	m = (*itab)(persistentalloc(unsafe.Sizeof(itab{})+uintptr(len(inter.mhdr)-1)*sys.PtrSize, 0, &memstats.other_sys))
	m.inter = inter
	m._type = typ
	m.init()
	itabAdd(m)
	unlock(&itabLock)
finish:
	if m.fun[0] != 0 {
		return m
	}
	if canfail {
		return nil
	}
	// this can only happen if the conversion
	// was already done once using the , ok form
	// and we have a cached negative result.
	// The cached result doesn't record which
	// interface function was missing, so initialize
	// the itab again to get the missing function name.
	panic(&TypeAssertionError{concrete: typ, asserted: &inter.typ, missingMethod: m.init()})
}
````
简单总结一下：`getitab` 函数会根据 `interfacetype` 和 `_type` 去全局的 `itab` 哈希表中查找，如果能找到，则直接返回；否则，会根据给定的 `interfacetype` 和 `_type` 新生成一个 `itab`，并插入到 `itab` 哈希表，这样下一次就可以直接拿到 `itab`。  
第一次去查询的时候如果查找到，直接返回

````go
if m = t.find(inter, typ); m != nil {
		goto finish
	}
````

如果在`hash`表中没有找到，这时候锁住`itabLock`，然后去重新写入`itab`到哈希表，当写入之后，上游的查询拿到值了，解除锁的阻塞，然后返回。

````go
if m = itabTable.find(inter, typ); m != nil {
		unlock(&itabLock)
		goto finish
	}
````

再来看一下 `itabAdd` 函数的代码：  

````go
// itabAdd adds the given itab to the itab hash table.
// itabLock must be held.
func itabAdd(m *itab) {
	// Bugs can lead to calling this while mallocing is set,
	// typically because this is called while panicing.
	// Crash reliably, rather than only when we need to grow
	// the hash table.
	if getg().m.mallocing != 0 {
		throw("malloc deadlock")
	}

	t := itabTable
	if t.count >= 3*(t.size/4) { // 75% load factor
		// Grow hash table.
		// t2 = new(itabTableType) + some additional entries
		// We lie and tell malloc we want pointer-free memory because
		// all the pointed-to values are not in the heap.
		t2 := (*itabTableType)(mallocgc((2+2*t.size)*sys.PtrSize, nil, true))
		t2.size = t.size * 2

		// Copy over entries.
		// Note: while copying, other threads may look for an itab and
		// fail to find it. That's ok, they will then try to get the itab lock
		// and as a consequence wait until this copying is complete.
		iterate_itabs(t2.add)
		if t2.count != t.count {
			throw("mismatched count during itab table copy")
		}
		// Publish new hash table. Use an atomic write: see comment in getitab.
		atomicstorep(unsafe.Pointer(&itabTable), unsafe.Pointer(t2))
		// Adopt the new table as our own.
		t = itabTable
		// Note: the old table can be GC'ed here.
	}
	t.add(m)
}
````

最后总结下:  

- 1、具体类型转空接口时，_type 字段直接复制源类型的 _type；调用 mallocgc 获得一块新内存，把值复制进去，data 再指向这块新内存。
- 2、具体类型转非空接口时，入参 tab 是编译器在编译阶段预先生成好的，新接口 tab 字段直接指向入参 tab 指向的 itab；调用 mallocgc 获得一块新内存，把值复制进去，data 再指向这块新内存。
- 3、而对于接口转接口，itab 调用 getitab 函数获取。只用生成一次，之后直接从 hash 表中获取。

### 接口的动态类型和动态值

````go
type iface struct {
	tab  *itab
	data unsafe.Pointer
}
````

`iface`我们可以看到，是有一个`tab`接口指针，指向数据类型，`data`数据指针，指向具体的数据。他们也被称为`动态类型`和`动态值`。  
因为两个都是指针，所以默认值都是`nil`。所以当两者都是`nil`的时候这个接口值才是`nil`,也就是`接口值 == nil`。  

````go
func main() {
	var f interface{}
	fmt.Println("+++动态类型和动态值都是nil+++")
	fmt.Println(f == nil)
	fmt.Printf("f: %T, %v\n", f, f)

	var g *string
	f = g
	fmt.Println("+++类型为 *string+++")
	fmt.Println(f == nil)
	fmt.Printf("f: %T, %v\n", f, f)
}
````

打印下输出:  

````go
+++动态类型和动态值都是nil+++
true
f: <nil>, <nil>
+++类型为 *string+++
false
f: *string, <nil> 
````

### interface如何支持泛型

严格来说，在 Golang 中并不支持泛型编程。在 C++ 等高级语言中使用泛型编程非常的简单，所以泛型编程一直是 `Golang` 诟病最多的地方。但是使用 `interface` 我们可以实现“泛型编程”，为什么？因为 `interface` 是一种抽象类型，任何具体类型（int, string）和抽象类型（user defined）都可以封装成 `interface`。以标准库的 sort 为例。

````go
package sort

// A type, typically a collection, that satisfies sort.Interface can be
// sorted by the routines in this package.  The methods require that the
// elements of the collection be enumerated by an integer index.
type Interface interface {
    // Len is the number of elements in the collection.
    Len() int
    // Less reports whether the element with
    // index i should sort before the element with index j.
    Less(i, j int) bool
    // Swap swaps the elements with indexes i and j.
    Swap(i, j int)
}

...

// Sort sorts data.
// It makes one call to data.Len to determine n, and O(n*log(n)) calls to
// data.Less and data.Swap. The sort is not guaranteed to be stable.
func Sort(data Interface) {
    // Switch to heapsort if depth of 2*ceil(lg(n+1)) is reached.
    n := data.Len()
    maxDepth := 0
    for i := n; i > 0; i >>= 1 {
        maxDepth++
    }
    maxDepth *= 2
    quickSort(data, 0, n, maxDepth)
}
````

`Sort` 函数的形参是一个 `interface`，包含了三个方法：`Len()，Less(i,j int)，Swap(i, j int)`。使用的时候不管数组的元素类型是什么类型`（int, float, string…）`，只要我们实现了这三个方法就可以使用 `Sort` 函数，这样就实现了“泛型编程”。有一点比较麻烦的是，我们需要自己封装一下。下面是一个例子。

````go
type Person struct {
    Name string
    Age  int
}

func (p Person) String() string {
    return fmt.Sprintf("%s: %d", p.Name, p.Age)
}

// ByAge implements sort.Interface for []Person based on
// the Age field.
type ByAge []Person //自定义

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }

func main() {
    people := []Person{
        {"Bob", 31},
        {"John", 42},
        {"Michael", 17},
        {"Jenny", 26},
    }

    fmt.Println(people)
    sort.Sort(ByAge(people))
    fmt.Println(people)
}
````
具体一点来说，也就是如果是在实现一个服务时，对于不同场景，可以将其共同特征抽象出来，在一个`interface`中声明，然后给不同的场景定义其特定的`struct`，上层的逻辑可以通过传入`interface`来执行，特化则通过`struct`实现对应的方法，从而达到一定程度的泛型。 


### 参考
【理解 Go interface 的 5 个关键点】https://sanyuesha.com/2017/07/22/how-to-understand-go-interface/  
【深入理解 Go Interface】https://zhuanlan.zhihu.com/p/32926119   
【GO如何支持泛型】https://zhuanlan.zhihu.com/p/74525591  
【Golang面向对象编程】https://code.tutsplus.com/zh-hans/tutorials/lets-go-object-oriented-programming-in-golang--cms-26540  
【深度解密Go语言之关于 interface 的10个问题】https://www.cnblogs.com/qcrao-2018/p/10766091.html  
【golang如何获取变量的类型：反射，类型断言】https://ieevee.com/tech/2017/07/29/go-type.html    
【Go接口详解】https://zhuanlan.zhihu.com/p/27055513  
