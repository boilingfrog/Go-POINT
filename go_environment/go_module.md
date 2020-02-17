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

#### 1、初始化

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
	want := "你好，世界。"
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

#### 2、依赖升级（降级）

可以使用如下命令来查看当前项目依赖的所有的包。
````
$ go list -m all
hello
golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c
rsc.io/quote v1.5.2
rsc.io/sampler v1.3.0
````
如果我们想要升级（降级）某个package则只需要go get就可以了，比如：
````
go get package@version
````
同时我们可以可以查看这个包所支持的版本，选择合适的进行升级
````
go list -m -versions package
````
栗子：  
版本查看：
````
$  go list -m -versions rsc.io/sampler
rsc.io/sampler v1.0.0 v1.2.0 v1.2.1 v1.3.0 v1.3.1 v1.99.99
````
查看当前的版本：
````
$ go list -m all
hello
golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c
rsc.io/quote v1.5.2
rsc.io/sampler v1.3.1
````
可以发现　rsc.io/sampler　现在的版本是v1.3.1，那么他最高支持的是v1.99.99  
版本升级：
````
$ go get rsc.io/sampler
go: finding rsc.io/sampler v1.99.99
go: downloading rsc.io/sampler v1.99.99
go: extracting rsc.io/sampler v1.99.99
````
当我们没有选择升级的版本，默认的版本就是最高的  
升级完成之后我们再来查看当前的版本
````
$ go list -m all
hello
golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c
rsc.io/quote v1.5.2
rsc.io/sampler v1.99.99
````
rsc.io/sampler的版本已经变成了v1.99.99，同样我们也可以对它的版本进行降级，方法和升级的方法一样  
````
$ go get rsc.io/sampler@v1.3.1
$ go list -m all
hello
golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c
rsc.io/quote v1.5.2
rsc.io/sampler v1.3.1
````
哈哈，已经变成了v1.3.1，降级成功了。


### 3、更改使用的pkg
让我们完成使用rsc的转换。rsc.io/quote 只使用rsc.io/quote/v3。首先我们看下 rsc.io/quote/v3支持的api，因为相比于rsc.io/quote，
rsc.io/quote/v3所支持的api可能已经发生了改变。
````
$ go doc rsc.io/quote/v3
package quote // import "rsc.io/quote"

Package quote collects pithy sayings.

func Concurrency() string
func GlassV3() string
func GoV3() string
func HelloV3() string
func OptV3() string
````
我们来更新我们的代码 hello.go
````
package hello

import "rsc.io/quote/v3"

func Hello() string {
	return quote.HelloV3()
}
````
然后执行go test
````
$ go test
go: downloading rsc.io/quote/v3 v3.1.0
go: extracting rsc.io/quote/v3 v3.1.0
PASS
ok      hello   0.003s
````
发现下载了rsc.io/quote/v3，并且成功运行

#### 4、清除不需要的依赖包
上面我们把rsc.io/quote换成了rsc.io/quote/v3，那么rsc.io/quote就已经不在需要了，这时候我们可以选择清除掉这些不需要的包
````
$ go list -m all
hello
golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c
rsc.io/quote v1.5.2
rsc.io/quote/v3 v3.1.0
rsc.io/sampler v1.3.1
$ cat go.mod
module hello

go 1.13

require (
        rsc.io/quote v1.5.2
        rsc.io/quote/v3 v3.1.0
        rsc.io/sampler v1.3.1 // indirect
)
````
我们知道go build(test)只会在运行的时候检查少了那些包，然后进行加载。但是不能确定何时可以安全地删除某些东西。 仅在检查模块中的所有软
件包以及这些软件包的所有可能的构建标记组合之后，才能删除依赖项。普通的build命令不会加载此信息，因此它不能安全地删除依赖项。  
可以使用go mod tidy清除不需要的包  
````
$ go mod tidy
$ go list -m all
hello
golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c
rsc.io/quote/v3 v3.1.0
rsc.io/sampler v1.3.1
$ cat go.mod
module hello

go 1.13

require (
        rsc.io/quote/v3 v3.1.0
        rsc.io/sampler v1.3.1 // indirect
)
$ go test
PASS
ok      hello   0.003s
````
### GoProxy
proxy 顾名思义，代理服务器。对于国内的网络环境，有些包是下载不下来的，当然有的人会用梯子去解决。go官方也意识到了这一点，提供了GOPROXY
的方法让我们下载包。要使用 GoProxy 只需要设置环境变量 GOPROXY 即可。目前公开的 GOPROXY 有：
  
- goproxy.io
- goproxy.cn: 由七牛云提供，这是一个应届生发起的项目，好强

当然你也可以实现自己的 GoProxy 服务，比如项目中的依赖包含外部依赖和内部依赖的时候，那么只需要实现 module proxy protocal 协议即可。

值得注意的是，在最新 release 的 Go 1.13 版本中默认将 GOPROXY 设置为 https://proxy.golang.org，这个对于国内的开发者是无法直接
使用的。所以如果升级了 Go 1.13 版本一定要把 GOPROXY 手动改掉。
````
export GOPROXY=https://goproxy.io
````