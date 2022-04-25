<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [go module了解](#go-module%E4%BA%86%E8%A7%A3)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [开启go mod](#%E5%BC%80%E5%90%AFgo-mod)
  - [简单使用](#%E7%AE%80%E5%8D%95%E4%BD%BF%E7%94%A8)
    - [1、初始化](#1%E5%88%9D%E5%A7%8B%E5%8C%96)
    - [2、依赖升级（降级）](#2%E4%BE%9D%E8%B5%96%E5%8D%87%E7%BA%A7%E9%99%8D%E7%BA%A7)
    - [3、更改使用的pkg](#3%E6%9B%B4%E6%94%B9%E4%BD%BF%E7%94%A8%E7%9A%84pkg)
    - [4、清除不需要的依赖包](#4%E6%B8%85%E9%99%A4%E4%B8%8D%E9%9C%80%E8%A6%81%E7%9A%84%E4%BE%9D%E8%B5%96%E5%8C%85)
  - [GoProxy](#goproxy)
  - [拉取私有仓库](#%E6%8B%89%E5%8F%96%E7%A7%81%E6%9C%89%E4%BB%93%E5%BA%93)
  - [安全性问题](#%E5%AE%89%E5%85%A8%E6%80%A7%E9%97%AE%E9%A2%98)
  - [总结](#%E6%80%BB%E7%BB%93)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## go module了解  


### 前言
Go 的包管理方式是逐渐演进的， 最初是 monorepo 模式，所有的包都放在 GOPATH 里面，使用类似命名空间的包路径区分包，不过这种包管理显然是有问题，由于包依赖可能会引入破坏性更新，生产环境和测试环境会出现运行不一致的问题。  

从 v1.5 开始开始引入 vendor 包模式，如果项目目录下有 vendor 目录，那么 go 工具链会优先使用 vendor 内的包进行编译、测试等，这之后第三方的包管理思路都是通过这种方式来实现，比如说由社区维护准官方包管理工具 dep。  

不过官方并不认同这种方式，在 v1.11 中加入了 Go Module 作为官方包管理形式，就这样 dep 无奈的结束了使命。最初的 `Go Module` 提案的名称叫做 vgo，下面为了介绍简称为 gomod。不过在 v1.11 和 v1.12 的 Go 版本中 gomod 是不能直接使用的。可以通过 go env 命令返回值的 GOMOD 字段是否为空来判断是否已经开启了 gomod，如果没有开启，可以通过设置环境变量 `export GO111MODULE=on` 开启。  

哈哈，是时候开始使用`go module`了  

### 开启go mod

GO111MODULE 有三个值：  

- on    打开，不会去 GOPATH 下面查找依赖包。  

- off   关闭  

- auto  Golang 自己检测是不是使用 modules 功能。  

linux:    

````  
export GO111MODULE=on // 开启
export GO111MODULE=of // 关闭
````
 
### 简单使用  

#### 1、初始化

需要注意的是：我们使用 `go mod` 不要在 GOPATH 目录下面

我们在 GOPATH 以外的目录建立一个项目 hello，然后在根目录下面执行  

````
liz@liz-PC:~/goWork/src/hello$ go mod init hello
go: creating new go.mod: module hello
````

可以看到 `go.mod` 已经生成在了项目的根木下面，然后我们去看下里面的内容  

````
liz@liz-PC:~/goWork/src/hello$ cat go.mod
module hello

go 1.13
````

发现里面还没有任何的依赖，因外我们的项目代码还没有创建，那么接下来去创建代码，看看 mod 如何进行依赖包的管理  
 
创建 `hello.go`  

````go
package hello

import "rsc.io/quote"

func Hello() string {
	return quote.Hello()
}
````  

然后创建 `testHello.go` 来测试  

````go
package hello

import "testing"

func TestHello(t *testing.T) {
	want := "你好，世界。"
	if got := Hello(); got != want {
		t.Errorf("Hello() = %q, want %q", got, want)
	}
}
````  

因为我们在 GOPATH 外，所以 `rsc.io/quote` 是找不到的。正常执行test找不到这些包，是会报错的，但是我们使用 `go mod` 这些就不会发生了，那我们执行 `go test` 看下    

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

发现会主动帮我们下载所需要的依赖包。`(go run/test)`，或者 build(执行命令 go build)的时候都会帮助我们下载所需要的依赖包。并且这些包都会下载到 `$GOPATH/pkg/mod` 下面，不同版本并存。    

这时候有些包是不能下载的，比如 `golang.org/x` 下的包，这时候可以使用代理的方式解决，go 也提供了 GoProxy 来帮助我们做这些事情，具体看下文 GoProxy。

#### 2、依赖升级（降级）

可以使用如下命令来查看当前项目依赖的所有的包。  

````
$ go list -m all
hello
golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c
rsc.io/quote v1.5.2
rsc.io/sampler v1.3.0
````

如果我们想要升级（降级）某个 package 则只需要 `go get` 就可以了，比如：  

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
  
可以发现 `rsc.io/sampler` 现在的版本是 `v1.3.1`，那么他最高支持的是 `v1.99.99 ` 

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

`rsc.io/sampler` 的版本已经变成了 `v1.99.99`，同样我们也可以对它的版本进行降级，方法和升级的方法一样    

````
$ go get rsc.io/sampler@v1.3.1
$ go list -m all
hello
golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c
rsc.io/quote v1.5.2
rsc.io/sampler v1.3.1
````

哈哈，已经变成了v1.3.1，降级成功了。  

#### 3、更改使用的pkg  

让我们完成使用 rsc 的转换。`rsc.io/quote` 只使用 `rsc.io/quote/v3`。首先我们看下 `rsc.io/quote/v3`支持的api，因为相比于 `rsc.io/quote`，`rsc.io/quote/v3` 所支持的 api 可能已经发生了改变。

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

然后执行 `go test`  

````
$ go test
go: downloading rsc.io/quote/v3 v3.1.0
go: extracting rsc.io/quote/v3 v3.1.0
PASS
ok      hello   0.003s
````

发现下载了 `rsc.io/quote/v3`，并且成功运行  

#### 4、清除不需要的依赖包

上面我们把 `rsc.io/quote` 换成了 `rsc.io/quote/v3`，那么 `rsc.io/quote` 就已经不在需要了，这时候我们可以选择清除掉这些不需要的包    

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

我们知道 `go build(test)` 只会在运行的时候检查少了那些包，然后进行加载。但是不能确定何时可以安全地删除某些东西。 仅在检查模块中的所有软件包以及这些软件包的所有可能的构建标记组合之后，才能删除依赖项。普通的 build 命令不会加载此信息，因此它不能安全地删除依赖项。    

可以使用 `go mod tidy` 清除不需要的包    

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

proxy 顾名思义，代理服务器。对于国内的网络环境，有些包是下载不下来的，当然有的人会用梯子去解决。go官方也意识到了这一点，提供了 GOPROXY 的方法让我们下载包。要使用 GoProxy 只需要设置环境变量 GOPROXY 即可。目前公开的 GOPROXY 有：
  
- goproxy.io  

- goproxy.cn: 由七牛云提供，这是一个应届生发起的项目，好强  

当然你也可以实现自己的 GoProxy 服务，比如项目中的依赖包含外部依赖和内部依赖的时候，那么只需要实现 module proxy protocal 协议即可。

值得注意的是，在最新 release 的 Go 1.13 版本中默认将 GOPROXY 设置为 `https://proxy.golang.org`，这个对于国内的开发者是无法直接使用的。所以如果升级了 Go 1.13 版本一定要把 GOPROXY 手动改掉。  

````
export GOPROXY=https://goproxy.io
````

有时候会看到 `export GOPROXY=https://goproxy.io,direct` 后面带了一个 direct  

作用就是如果代理中没有找到，`go get`会接着中从 VCS（例如github.com）源地址下载依赖。  

栗子  

比如我们的私有仓库的资源，公共代理是找不到的，如果配置了 direct ，公共代理中找不到，就会向我们的私有地址直接发起请求，这样就能拉取我们私有仓库的资源了。  

### 拉取私有仓库

对于私有仓库资源的拉取，除了上面讲的配置 direct，go 中也专门提供了 GOPRIVATE 这个环境变量，来配置我们的私有仓库地址  

设置完 GOPROXY ，我们就能拉取墙外的包了，但是私有仓库的包就拉取不到了，因为配置的 GOPROXY,是不能访问我们的私有仓库的，当然 GOPROXY 的代理如果是自己搭建的私有代理，这种除外  

可以设置 GOPRIVATE 来拉取我们的私有仓库代码  

```
export GOPRIVATE="gitlab.test.com"
```

设置了 GOPRIVATE 也可以跳过安全性检查，具体见下文  

### 安全性问题

go 处于安全性考虑，为了保证开发者的依赖库不被人恶意劫持篡改，所以引入了 GOSUMDB 环境变量来设置校验服务器   

当你在本地对依赖进行变动（更新/添加）操作时，Go 会自动去这个服务器进行数据校验，保证你下的这个代码库和世界上其他人下的代码库是一样的。如果有问题，会有个大大的安全提示。当然背后的这些操作都已经集成在 Go 里面了，开发者不需要进行额外的操作。    

golang 的服务器由 Google 托管默认地址是 `sum.golang.org`，国内没有梯子的情况下是不能访问的，不过国内可以使用 `https://gosum.io/`  

``
export GOSUMDB=gosum.io+ce6e7565+AY5qEHUk/qmHc5btzW45JVoENfazw8LielDsaI+lEbq6
``  

对于我们的私有仓库，去公共安全校验库校验，肯定是不能通过校验的，我们可以通过 GONOSUMDB 这个环境变量来设置不做校验的代码仓库， 它可以设置多个匹配路径，用逗号相隔。  

```
export GONOSUMDB=gitlab.test.com
```

如果对我们的我们的私有仓库设置了 GOPRIVATE ，那么也会跳过安全检查，这时候就不用配置 GONOSUMDB 了。  

为了更灵活的控制那些依赖软件包经过 `proxy server` 和 sumdb 校验，可以通过 GONOPROXY 和 GONOSUMDB 来单独进行控制，这两个环境变量的被设置后将覆盖 GOPRIVATE 环境变量，同样这两个变量也支持逗号分隔。  

### 总结  

Go 1.11之后的版本都支持`go module`了， modules 在 Go 1.13 的版本下是默认开启的。国内使用需要修改 GoProxy 的代理。    
 
下面是经常使用到的命令：  
 
 ````
go mod init 创建初始化go module  
go build, go test 运行的时候就会加载所需要的依赖项
go list -m all 打印档当前模块所有的依赖
go get 更改模块依赖的版本
go mod tidy 删除不适用的依赖
`````

对于我们的私有仓库，我们除了配置 GOPROXY 之外，同时也会配置一个 GOPRIVATE ；  

配置了 GOPRIVATE 之后，除了在拉取我们私有仓库资源的时候会跳过 GOPROXY，同时也会跳过安全性检查。  

### 参考  

【Using Go Modules】https://blog.golang.org/using-go-modules  
【Go Modules 不完全教程】https://mp.weixin.qq.com/s/v-NdYEJBgKbiKsdoQaRsQg  
【GOSUMDB 环境变量】https://goproxy.io/zh/docs/GOSUMDB-env.html    
【如何使用go mod】https://boilingfrog.github.io/2022/04/25/go%E4%B8%ADgomod%E4%BD%BF%E7%94%A8%E5%B0%8F%E7%BB%93/  