## go中环境安装

### 前言
最近在工作中需要新配置go环境，每次都要去网上查找教程，浪费时间，那么就自己总结下。

### 环境中一些基本的配置项

##### GOROOT

golang的安装路径。
在Windows中，GOROOT的默认值是C:/go，而在Mac OS或Linux中GOROOT的默认值是usr/loca/go，如果
将Go安装在其他目录中，而需要将GOROOT的值修改为对应的目录。

同时，我们需要将GOROOT/bin则包含Go为我们提供的工具链，因此，应该将GOROOT/bin配置到环境变量
PATH中，方便我们在全局中使用Go工具链。

##### linux下面设置GOROOT
````
export GOROOT=~/go
export PATH=$PATH:$GOROOT/bin
````

#### GOPATH

go的工作目录  

需要注意的是，我们一般设置GOPATH的路径和GOROOT的路径不一样。  

环境变量GOPATH用于指定我们的开发工作区(workspace),是存放源代码、测试文件、库静态文件、可执行文
件的工作。

在类Unix(Mac OS或Linux)操作系统中GOPATH的默认值是$home/go。而在Windows中GOPATH的默认值则
为%USERPROFILE%\go(比如在Admin用户，其值为C:\Users\Admin\go)。

##### linux设置GOPATH
````
export GOPATH=/home/liz/go
````
GOPATH的工作目录是可以设置多个的，比如
````
export GOPATH=/hone/liz/go:/home/liz/gowork
````
但是GOPATH目录里面必须包含三个子目录  
- bin golang编译可执行文件存放路径，可自动生成。
- src 源码路径。按照golang默认约定，go run，go install等命令的当前工作路径（即在此路径下执行上述命令）。
- pkg golang编译的.a中间文件存放路径，可自动生成。  

### GOBIN

环境变量GOBIN表示我们开发程序编译后二进制命令的安装目录。  

当我们使用go install命令编译和打包应用程序时，该命令会将编译后二进制程序打包GOBIN目录，一般
我们将GOBIN设置为GOPATH/bin目录。

不允许设置多个路径。可以为空。为空时则遵循“约定优于配置”原则，可执行文件放在各自GOPATH目录的bin
文件夹中（前提是：package main的main函数文件不能直接放到GOPATH的src下面。

##### linux下面设置GOBIN
````
export GOBIN=$GOPATH/bin
````
上面的代码中，我们都是使用export命令设置环境变量的，这样设置只能在当前shell中有效，如果想一直
有效，如在Linux中，则应该将环境变量添加到/etc/profile等文件当中。

### PATH

这个是环境变量的路径，上面介绍的我们要将GOROOT下面的bin目录，加入到PATH中，同时我们也要注意把
GOPATH下面的bin也放进去，当然GOBIN加进去也行，毕竟GOBIN也是指向这个目录的，不然我们生成的可执
行文件文件就不能全局的被执行。

### 交叉编译

什么是交叉编译？所谓的交叉编译，是指在一个平台上就能生成可以在另一个平台运行的代码，例如，我们
可以32位的Windows操作系统开发环境上，生成可以在64位Linux操作系统上运行的二进制程序。  

在其他编程语言中进行交叉编译可能要借助第三方工具，但在Go语言进行交叉编译非常简单，最简单只需要
设置GOOS和GOARCH这两个环境变量就可以了。  

##### GOOS与GOARCH

GOOS的默认值是我们当前的操作系统， 如果windows，linux,注意mac os操作的上的值是darwin。 GOARCH
则表示CPU架构，如386，amd64,arm等。

##### 获取GOOS和GOARCH的值

````
$ go env GOOS GOARCH
````
##### GOOS和GOARCH的取值范围

GOOS和GOARCH的值成对出现，而且只能是下面列表对应的值。  
````
$GOOS	    $GOARCH
android	    arm
darwin	    386
darwin	    amd64
darwin	    arm
darwin	    arm64
dragonfly   amd64
freebsd	    386
freebsd	    amd64
freebsd	    arm
linux	    386
linux	    amd64
linux	    arm
linux	    arm64
linux	    ppc64
linux	    ppc64le
linux	    mips
linux	    mipsle
linux	    mips64
linux	    mips64le
linux	    s390x
netbsd	    386
netbsd	    amd64
netbsd	    arm
openbsd	    386
openbsd	    amd64
openbsd	    arm
plan9	    386
plan9	    amd64
solaris	    amd64
windows	    386
windows	    amd64
````

##### 编译示例
编译在64位Linux操作系统上运行的目标程序
 ````
$ GOOS=linux GOARCH=amd64 go build main.go
`````
编译arm架构Android操作上的目标程序
````
$ GOOS=android GOARCH=arm GOARM=7 go build main.go
````

### go build 
这个命令主要用于编译代码。在包的编译过程中，若有必要，会同时编译与之相关联的包。  
- 普通的包：go build之后不会产生任何的文件，如果想要得到可执行的文件，就需要使用go install
- main包：go build之后会在本目录生成可执行的文件。如果想要放到$GOPATH/bin中，则需要执行
go install或go build -o 路径/a.exe
- 项目里面有多个文件，我们只想执行其中的一个文件，就在go build后面加上文件的名，go build XX.go
go build默认执行是包里面全部的文件
- go build会忽略目录下以“_”或“.”开头的go文件。

参数介绍  

- -o 指定输出的文件名，可以带上路径，例如 go build -o a/b/c
- -i 安装相应的包，编译+go install
- -a 更新全部已经是最新的包的，但是对标准包不适用
- -n 把需要执行的编译命令打印出来，但是不执行，这样就可以很容易的知道底层是如何运行的
- -p n 指定可以并行可运行的编译数目，默认是CPU数目
- -race 开启编译的时候自动检测数据竞争的情况，目前只支持64位的机器
- -v 打印出来我们正在编译的包名
- -work 打印出来编译时候的临时文件夹名称，并且如果已经存在的话就不要删除
- -x 打印出来执行的命令，其实就是和-n的结果类似，只是这个会执行
- -ccflags 'arg list' 传递参数给5c, 6c, 8c 调用
- -compiler name 指定相应的编译器，gccgo还是gc
- -gccgoflags 'arg list' 传递参数给gccgo编译连接调用
- -gcflags 'arg list' 传递参数给5g, 6g, 8g 调用
- -installsuffix suffix 为了和默认的安装包区别开来，采用这个前缀来重新安装那些依赖的包，-race的时候默认已经是-installsuffix race,大家可以通过-n命令来验证
- -ldflags 'flag list' 传递参数给5l, 6l, 8l 调用
- -tags 'tag list' 设置在编译的时候可以适配的那些tag，详细的tag限制参考里面的 Build Constraints

### go get
这个命令是用来动态获取远程代码包的，目前支持的有BitBucket、GitHub、Google Code和Launchpad。这个命令在内部
实际上分成了两步操作：第一步是下载源码包，第二步是执行go install。  

参数介绍：

- -d 只下载不安装
- -f 只有在你包含了-u参数的时候才有效，不让-u去验证import中的每一个都已经获取了，这对于本地fork的包特别有用
- -fix 在获取源码之后先运行fix，然后再去做其他的事情
- -t 同时也下载需要为运行测试所需要的包
- -u 强制使用网络去更新包和它的依赖包
- -v 显示执行的命令

### go install 

这个命令在内部实际上分成了两步操作：第一步是生成结果文件(可执行文件或者.a包)，第二步会把编译好
的结果移到$GOPATH/pkg或者$GOPATH/bin。其中bin下面放的是可执行文件。

参数支持go build的编译参数。

### go run

编译并运行Go程序

### go test

执行这个命令，会自动读取源码目录下面名为*_test.go的文件，生成并运行测试用的可执行文件

### go clean

这个命令是用来移除当前源码包和关联源码包里面编译生成的文件。这些文件包括  
````
_obj/            旧的object目录，由Makefiles遗留
_test/           旧的test目录，由Makefiles遗留
_testmain.go     旧的gotest文件，由Makefiles遗留
test.out         旧的test记录，由Makefiles遗留
build.out        旧的test记录，由Makefiles遗留
*.[568ao]        object文件，由Makefiles遗留

DIR(.exe)        由go build产生
DIR.test(.exe)   由go test -c产生
MAINFILE(.exe)   由go build MAINFILE.go产生
*.so             由 SWIG 产生
````
参数介绍

- -i 清除关联的安装的包和可运行文件，也就是通过go install安装的文件
- -n 把需要执行的清除命令打印出来，但是不执行，这样就可以很容易的知道底层是如何运行的
- -r 循环的清除在import中引入的包
- -x 打印出来执行的详细命令，其实就是-n打印的执行版本