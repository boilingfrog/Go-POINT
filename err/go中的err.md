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










### 参考

【Go语言(golang)的错误(error)处理的推荐方案】https://www.flysnow.org/2019/01/01/golang-error-handle-suggestion.html  
【Golang error 的突围】https://www.cnblogs.com/qcrao-2018/p/11538387.html  
【The Go Blog Error handling and Go】https://blog.golang.org/error-handling-and-go