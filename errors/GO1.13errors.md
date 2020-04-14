- [go1.13errors的用法](#go113errors%E7%9A%84%E7%94%A8%E6%B3%95)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [基本用法](#%E5%9F%BA%E6%9C%AC%E7%94%A8%E6%B3%95)
    - [fmt.Errorf](#fmterrorf)
    - [Unwrap](#unwrap)
    - [errors.Is](#errorsis)
    - [As](#as)
    - [扩展](#%E6%89%A9%E5%B1%95)
  - [参考](#%E5%8F%82%E8%80%83)


## go1.13errors的用法

### 前言

go 1.13发布了`error`的一些新的特性，那么就来探究学习下。

### 基本用法

#### fmt.Errorf

使用 `fmt.Errorf` 加上 `%w` 格式符来生成一个嵌套的 `error`，它并没有像 `pkg/errors` 那样使用一个 `Wrap` 函数来嵌套` error`，非常简洁。  

````go
err1 := errors.New("new error")
err2 := fmt.Errorf("err2: [%w]", err1)
err3 := fmt.Errorf("err3: [%w]", err2)
fmt.Println(err3)
````
输出
````
// output
err3: [err2: [new error]]
````
`err2` 就是一个合法的被包装的 `error`，同样地，`err3` 也是一个被包装的 `error`，如此可以一直套下去。  

#### Unwrap

拆开一个被包装的 `error`  

````go
func Unwrap(err error) error
````
将嵌套的 error 解析出来，多层嵌套需要调用 Unwrap 函数多次，才能获取最里层的 error。  

源码如下：  

````go
func Unwrap(err error) error {
    // 判断是否实现了 Unwrap 方法
	u, ok := err.(interface {
		Unwrap() error
	})
	// 如果不是，返回 nil
	if !ok {
		return nil
	}
	// 调用 Unwrap 方法返回被嵌套的 error
	return u.Unwrap()
}
````

对 `err` 进行断言，看它是否实现了 `Unwrap` 方法，如果是，调用它的 `Unwrap` 方法。否则，返回 `nil`。

````go
err1 := errors.New("new error")
err2 := fmt.Errorf("err2: [%w]", err1)
err3 := fmt.Errorf("err3: [%w]", err2)

fmt.Println(errors.Unwrap(err3))
fmt.Println(errors.Unwrap(errors.Unwrap(err3)))
````

输出

````go
// output
err2: [new error]
new error
````

#### errors.Is

判断被包装的error是是否含有指定错误。  

当多层调用返回的错误被一次次地包装起来，我们在调用链上游拿到的错误如何判断是否是底层的某个错误呢？  

它递归调用 Unwrap 并判断每一层的 err 是否相等，如果有任何一层 err 和传入的目标错误相等，则返回 true。  

源码如下：
````go
func Is(err, target error) bool {
	if target == nil {
		return err == target
	}

	isComparable := reflectlite.TypeOf(target).Comparable()
	
	// 无限循环，比较 err 以及嵌套的 error
	for {
		if isComparable && err == target {
			return true
		}
		// 调用 error 的 Is 方法，这里可以自定义实现
		if x, ok := err.(interface{ Is(error) bool }); ok && x.Is(target) {
			return true
		}
		// 返回被嵌套的下一层的 error
		if err = Unwrap(err); err == nil {
			return false
		}
	}
}
````
通过一个无限循环，使用 `Unwrap` 不断地将 `err` 里层嵌套的 `error` 解开，再看被解开的 `error` 是否实现了 `Is` 方法，并且调用它的 `Is` 方法，当两者都返回 `true` 的时候，整个函数返回 `true`。

举个栗子

````go
err1 := errors.New("new error")
err2 := fmt.Errorf("err2: [%w]", err1)
err3 := fmt.Errorf("err3: [%w]", err2)

fmt.Println(errors.Is(err3, err2))
fmt.Println(errors.Is(err3, err1))
````
输出
````
// output
true
true
````

#### As

这个和上面的 `errors.Is` 大体上是一样的，区别在于 `Is` 是严格判断相等，即两个 `error` 是否相等。而 `As` 则是判断类型是否相同，并提取第一个符合目标类型的错误，用来统一处理某一类错误。

````go
func As(err error, target interface{}) bool
````

源码如下：

````go
func As(err error, target interface{}) bool {
    // target 不能为 nil
	if target == nil {
		panic("errors: target cannot be nil")
	}
	
	val := reflectlite.ValueOf(target)
	typ := val.Type()
	
	// target 必须是一个非空指针
	if typ.Kind() != reflectlite.Ptr || val.IsNil() {
		panic("errors: target must be a non-nil pointer")
	}
	
	// 保证 target 是一个接口类型或者实现了 Error 接口
	if e := typ.Elem(); e.Kind() != reflectlite.Interface && !e.Implements(errorType) {
		panic("errors: *target must be interface or implement error")
	}
	targetType := typ.Elem()
	for err != nil {
	    // 使用反射判断是否可被赋值，如果可以就赋值并且返回true
		if reflectlite.TypeOf(err).AssignableTo(targetType) {
			val.Elem().Set(reflectlite.ValueOf(err))
			return true
		}
		
		// 调用 error 自定义的 As 方法，实现自己的类型断言代码
		if x, ok := err.(interface{ As(interface{}) bool }); ok && x.As(target) {
			return true
		}
		// 不断地 Unwrap，一层层的获取嵌套的 error
		err = Unwrap(err)
	}
	return false
}
````

举个栗子

````go
type ErrorString struct {
    s string
}

func (e *ErrorString) Error() string {
    return e.s
}

var targetErr *ErrorString
err := fmt.Errorf("new error:[%w]", &ErrorString{s:"target err"})
fmt.Println(errors.As(err, &targetErr))
````
输出
````
// output
true
````

#### 扩展

`Is As` 两个方法已经预留了口子，可以由自定义的 `error struct` 实现并覆盖调用。


### 参考
 
【Go 1.13 errors 基本用法】https://segmentfault.com/a/1190000020398774  
【Go语言(golang)新发布的1.13中的Error Wrapping深度分析】https://www.flysnow.org/2019/09/06/go1.13-error-wrapping.html