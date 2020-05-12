## go中interface转换成原来的类型

### 首先了解下interface

#### 什么是`interface`?

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

#### 如何判断 interface 变量存储的是哪种类型

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

对于fmt也是用了反射的:

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


### 参考
【理解 Go interface 的 5 个关键点】https://sanyuesha.com/2017/07/22/how-to-understand-go-interface/ 
【深入理解 Go Interface】https://zhuanlan.zhihu.com/p/32926119   
【GO如何支持泛型】https://zhuanlan.zhihu.com/p/74525591  
【Golang面向对象编程】https://code.tutsplus.com/zh-hans/tutorials/lets-go-object-oriented-programming-in-golang--cms-26540  
【深度解密Go语言之关于 interface 的10个问题】https://www.cnblogs.com/qcrao-2018/p/10766091.html  
【golang如何获取变量的类型：反射，类型断言】https://ieevee.com/tech/2017/07/29/go-type.html    