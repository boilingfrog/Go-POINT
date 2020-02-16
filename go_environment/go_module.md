## go module了解  


### 前言
Go 的包管理方式是逐渐演进的， 最初是 monorepo 模式，所有的包都放在 GOPATH 里面，使用类似命名
空间的包路径区分包，不过这种包管理显然是有问题，由于包依赖可能会引入破坏性更新，生产环境和测试环
境会出现运行不一致的问题。  

从 v1.5 开始开始引入 vendor 包模式，如果项目目录下有 vendor 目录，那么 go 工具链会优先使
用 vendor 内的包进行编译、测试等，这之后第三方的包管理思路都是通过这种方式来实现，比如说由社区
维护准官方包管理工具 dep。  

不过官方并不认同这种方式，在 v1.11 中加入了 Go Module 作为官方包管理形式，就这样 dep 无奈的
结束了使命。最初的 Go Module 提案的名称叫做 vgo，下面为了介绍简称为 gomod。不过在 v1.11 和 
v1.12 的 Go 版本中 gomod 是不能直接使用的。可以通过 go env 命令返回值的 GOMOD 字段是否为空
来判断是否已经开启了 gomod，如果没有开启，可以通过设置环境变量 export GO111MODULE=on 开启。  

哈哈，是时候开始使用go module了  

### 开启go mod

GO111MODULE 有三个值：
- on    打开
- off   关闭
- auto  默认值

linux:  
````  
export GO111MODULE=on // 开启
export GO111MODULE=of // 关闭
````
 
### 简单使用  
需要注意的是：我们使用go mod不要在GOPATH目录下面

我们在GOPATH以为的目录建立一个项目hello，然后在根目录下面执行
````
liz@liz-PC:~/goWork/src/hello$ go mod init hello
go: creating new go.mod: module hello
````
可以看到go.mod已经生成在了项目的根木下面，然后我们去看下里面的内容
````
liz@liz-PC:~/goWork/src/hello$ cat go.mod
module hello

go 1.13
````
发现里面还没有任何的依赖，因外我们的项目代码还没有创建，那么接下来去创建代码，看看mod如何进行
依赖包的管理，创建hello.go
````
package hello

import "rsc.io/quote"

func Hello() string {
	return quote.Hello()
}
````
然后创建testHello.go来测试
````
package hello

import "testing"

func TestHello(t *testing.T) {
	want := "Hello, world."
	if got := Hello(); got != want {
		t.Errorf("Hello() = %q, want %q", got, want)
	}
}
````
因为我们在GOPATH外，所以rsc.io/quote是找不到的。正常执行test找不到这些包，是会报错的，但是我们
使用go mod这些就不会发生了，那我们执行go test看下
````
go test
go: finding rsc.io/quote v1.5.2
go: downloading rsc.io/quote v1.5.2
go: extracting rsc.io/quote v1.5.2
go: downloading rsc.io/sampler v1.3.0
go: extracting rsc.io/sampler v1.3.0
go: downloading golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c
go: extracting golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c
go: finding rsc.io/sampler v1.3.0
go: finding golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c
--- FAIL: TestHello (0.00s)
    hello_test.go:8: Hello() = "你好，世界。", want "Hello, world."
FAIL
exit status 1
FAIL    hello   0.003s

````
发现会主动帮我们下载所需要的依赖包。(go run/test)，或者 build(执行命令 go build)的时候都会帮助
我们下载所需要的依赖包。并且这些包都会下载到 $GOPATH/pkg/mod 下面，不同版本并存。  

这时候有些包是不能下载的，比如golang.org/x下的包，这时候可以使用代理的方式解决，go也提供了GoProxy
来帮助我们做这些事情，具体看下文GoProxy。



