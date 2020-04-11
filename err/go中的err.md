## err

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

使用 New 函数创建出来的 error 类型实际上是 errors 包里未导出的 errorString 类型，它包含唯一的一个字段 s，并且实现了唯一的方法：Error() string。

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

当我们想知道错误类型，并且打印错误的时候，直接打印 error：  

````go
fmt.Println(err)
````

或者：  

````go
fmt.Println(err.Error)
````
``fmt`` 包会自动调用 ``err.Error()`` 函数来打印字符串。  

****注意：对于err我们都是将err放在函数返回值的最后一个，同时对于会出错的函数我们都会返回一个err，当然对于一些函数，我们可能不确定之后是否会有错误的产生，所以一般也是预留一个err的返回。****


















### 参考

【Go语言(golang)的错误(error)处理的推荐方案】https://www.flysnow.org/2019/01/01/golang-error-handle-suggestion.html  
【Golang error 的突围】https://www.cnblogs.com/qcrao-2018/p/11538387.html  
【The Go Blog Error handling and Go】https://blog.golang.org/error-handling-and-go  
【Go 1.13 errors 基本用法】https://segmentfault.com/a/1190000020398774  