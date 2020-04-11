## go中的err

go中的错误处理，是通过返回值的形式来出来，要么你忽略，要么你处理（处理也可以是继续返回给调用者），对于golang这种设计方式，我们会在代码中写大量的if判断，以便做出决定。

````go
func main() {
	conent,err:=ioutil.ReadFile("filepath")
	if err !=nil{
		//错误处理
	}else {
		fmt.Println(string(conent))
	}
}
````
对于err如果是nil就代表没有错误，如果不是nil就代表程序出问题了，需要对错误进行处理了。

### err接口

go中的err是一个接口类型  

````go
// The error built-in interface type is the conventional interface for
// representing an error condition, with the nil value representing no error.
type error interface {
	Error() string
}
````

``errors.New()``是我们会经常使用的，我们来探究下这个函数  

````go
// src/errors/errors.go

func New(text string) error {
	return &errorString{text}
}

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}
````

使用 `New` 函数创建出来的 `error` 类型实际上是 `errors` 包里未导出的 `errorString` 类型，它包含唯一的一个字段 s，并且实现了唯一的方法：`Error() string。`

举个使用的栗子  

````go
func Sqrt(f float64) (float64, error) {
    if f < 0 {
        return 0, errors.New("math: square root of negative number")
    }
    // implementation
}
````

我们可以使用`errors.New`来定制我们需要的错误信息  

但是对于上面的报错，我们知道是不知道报错的上下文信息的，我们就知道程序出错了，不利于我们错误的排查。我们可以使用`fmt.Errorf`来输出上下文信息。

````go
if f < 0 {
    return 0, fmt.Errorf("math: square root of negative number %g", f)
}
````
通过`fmt.Errorf`我们不仅能打印错误，同时还能看到具体什么数数值引起的错误。它会先将字符串格式化，然后再调用`errors.New`来创建错误。

当我们想知道错误类型，并且打印错误的时候，直接打印 `error`：  

````go
fmt.Println(err)
````

或者：  

````go
fmt.Println(err.Error)
````
``fmt`` 包会自动调用 ``err.Error()`` 函数来打印字符串。  

****注意：对于err我们都是将err放在函数返回值的最后一个，同时对于会出错的函数我们都会返回一个err，当然对于一些函数，我们可能不确定之后是否会有错误的产生，所以一般也是预留一个err的返回。****

### go中err的困局

在go中，`err`的是通过返回值的形式返回。编程人员，要不处理，要不忽略。所以我们的代码就会大量的出现对错误的`if`判断。  
出于对代码的健壮性考虑，我们对于每一个错误，都是不能忽略的。因为出错的同时，很可能会返回一个 nil 类型的对象。如果不对错误进行判断，那下一行对 `nil` 对象的操作百分之百会引发一个 `panic`。
所以就造成了`err`满天飞。

还有比如，我们想对返回的`error`附加更多的信息后再返回，比如以上的例子，我们怎么做呢？我们只能先通过`Error`方法，取出原来的错误信息，然后自己再拼接，再使用`errors.New`函数生成新错误返回。


### 推荐方法

`github.com/pkg/errors`这个包给除了解决的方案。  

它的使用非常简单，如果我们要新生成一个错误，可以使用New函数,生成的错误，自带调用堆栈信息。 

````go
func New(message string) error
````

如果有一个现成的error，我们需要对他进行再次包装处理，这时候有三个函数可以选择。

````go
//只附加新的信息
func WithMessage(err error, message string) error

//只附加调用堆栈信息
func WithStack(err error) error

//同时附加堆栈和信息
func Wrap(err error, message string) error
````

这个错误处理库为我们提供了`Cause`函数让我们可以获得最根本的错误原因。

````go
func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}
````

使用`for`循环一直找到最根本（最底层）的那个`error`。

以上的错误我们都包装好了，也收集好了，那么怎么把他们里面存储的堆栈、错误原因等这些信息打印出来呢？其实，这个错误处理库的错误类型，都实现了`Formatter`接口，我们可以通过`fmt.Printf`函数输出对应的错误信息。

````
%s,%v //功能一样，输出错误信息，不包含堆栈
%q //输出的错误信息带引号，不包含堆栈
%+v //输出错误信息和堆栈
````

以上如果有循环包装错误类型的话，会递归的把这些错误都会输出。

### 总结

对于err我们都是将err放在函数返回值的最后一个，同时对于会出错的函数我们都会返回一个err，当然对于一些函数，我们可能不确定之后是否会有错误的产生，所以一般也是预留一个err的返回。  

Go 中的 error 过于简单，以至于无法记录太多的上下文信息，对于错误包裹也没有比较好的办法。当然，这些可以通过第三方库来解决。



### 参考

【Go语言(golang)的错误(error)处理的推荐方案】https://www.flysnow.org/2019/01/01/golang-error-handle-suggestion.html  
【Golang error 的突围】https://www.cnblogs.com/qcrao-2018/p/11538387.html  
【The Go Blog Error handling and Go】https://blog.golang.org/error-handling-and-go  
【Go 1.13 errors 基本用法】https://segmentfault.com/a/1190000020398774  
【Go与Error的前世今生】https://zhuanlan.zhihu.com/p/55975116  