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

需要注意的是，我们一边设置GOPATH的路径一遍和GOROOT的路径不一样。  

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

