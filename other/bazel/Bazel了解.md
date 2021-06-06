<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Bazel使用了解](#bazel%E4%BD%BF%E7%94%A8%E4%BA%86%E8%A7%A3)
  - [Bazel产生的背景](#bazel%E4%BA%A7%E7%94%9F%E7%9A%84%E8%83%8C%E6%99%AF)
  - [什么是Bazel](#%E4%BB%80%E4%B9%88%E6%98%AFbazel)
    - [快(Fast)](#%E5%BF%ABfast)
    - [可伸缩(scalable)](#%E5%8F%AF%E4%BC%B8%E7%BC%A9scalable)
    - [跨语言(multi-language)](#%E8%B7%A8%E8%AF%AD%E8%A8%80multi-language)
    - [可扩展(extensible)](#%E5%8F%AF%E6%89%A9%E5%B1%95extensible)
  - [使用Bazel部署go应用](#%E4%BD%BF%E7%94%A8bazel%E9%83%A8%E7%BD%B2go%E5%BA%94%E7%94%A8)
    - [手动通过Bazel部署go应用](#%E6%89%8B%E5%8A%A8%E9%80%9A%E8%BF%87bazel%E9%83%A8%E7%BD%B2go%E5%BA%94%E7%94%A8)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Bazel使用了解

### Bazel产生的背景

1、开源成为当前软件开发的主旋律。哪怕你是商业软件，也逃离不了社区的包围。如何方便地获取依赖，并做到平滑升级很重要。如果构建工具能够很方便地获取源代码，那就太好了。  

2、混合多语言编程成为一种选择。每种语言都有自己适用的场景，但是构建多语言的软件系统非常具有挑战性。例如，Python社区很喜欢搭配C/C++，高性能计算扔个Ｃ/C++，Python提供编程接口。如果构建工具能够无缝支持多语言构建，真的很方便。  

3、代码复用。我只想复用第三方的一个头文件，而不是整个系统。拒绝拷贝是优秀程序员的基本素养，如果构建工具能帮我方便地获取到所依赖的组件，剔除不必要的依赖，那就太完美了。  

4、增量构建。当只修改了一行代码，构建系统能够准确计算需要构建的依赖目标，而不是全构建；否则生命都浪费在编译上了。  

5、云构建。大型软件公司，复用计算资源，可以带来巨大的收益。  

### 什么是Bazel

`Bazel`是一个支持多语言、跨平台的构建工具。`Bazel`支持任意大小的构建目标，并支持跨多个仓库的构建，是`Google`主推的一种构建工具。  

`bazel`优点很多，主要有  

- 构建快。支持增量编译。对依赖关系进行了优化，从而支持并发执行。

- 可构建多种语言。bazel可用来构建Java C++ Android ios等很多语言和框架，并支持mac windows linux等不同平台

- 可伸缩。可处理任意大小的代码库，可处理多个库，也可以处理单个库

- 可扩展。使用bazel扩展语言可支持新语言和新平台。

#### 快(Fast)

`Bazel`的构建过程很快，它集合了之前构建系统的加速的一些常见做法。包括：    

1、增量编译。只重新编译必须的部分，即通过依赖分析，只编译修改过的部分及其影响的路径。    

2、并行编译。将没有依赖的部分进行并行执行，可以通过`--jobs`来指定并行流的个数，一般可以是你机器`CPU`的个数。遇到大项目马力全开时，`Bazel`能把你机器的`CPU`各个核都吃满。    

3、分布式/本地缓存。`Bazel`将构建过程视为函数式的，只要输入给定，那么输出就是一定的。而不会随着构建环境的不同而改变（当然这需要做一些限制），这样就可以分布式的缓存/复用不同模块，这点对于超大项目的速度提升极为明显。  

#### 可伸缩(scalable)

`Bazel`号称无论什么量级的项目都可以应对，无论是超大型单体项目`monorepo`、还是超多库的分布式项目`multirepo`。`Bazel`还可以很方便的集成`CD/CI` ，并在云端利用分布式环境进行构建。  

它使用沙箱机制进行编译，即将所有编译依赖隔绝在一个沙箱中，比如编译`golang`项目时，不会依赖你本机的`GOPATH`，从而做到同样源码、跨环境编译、输出相同，即构建的确定性。   

#### 跨语言(multi-language)

如果一个项目不同模块使用不同的语言，利用`Bazel`可以使用一致的风格来管理项目外部依赖和内部依赖。典型的项目如 Ray。该项目使用`C++`构建`Ray`的核心调度组件、通过`Python/Java`来提供多语言的`API`，并将上述所有模块用单个`repo`进行管理。如此组织使其项目整合相当困难，但`Bazel`在此处理的游刃有余，大家可以去该`repo`一探究竟。  

#### 可扩展(extensible)

`Bazel`使用的语法是基于`Python`裁剪而成的一门语言：`Startlark`。其表达能力强大，往小了说，可以使用户自定义一些`rules`（类似一般语言中的函数）对构建逻辑进行复用；往大了说，可以支持第三方编写适配新的语言或平台的`rules`集，比如`rules go`。 `Bazel`并不原生支持构建`golang`工程，但通过引入`rules go` ，就能以比较一致的风格来管理`golang`工程。  

### 使用Bazel部署go应用

1、安装`Bazel`  

```
$ brew install Bazel
```

这是mac下面的安装，其他平台自行google  

2、安装`gazelle`

```
$ go get github.com/bazelbuild/bazel-gazelle/cmd/gazelle
```

#### 手动通过Bazel部署go应用

查看运行的结果

```
$ bazel run //:test
Starting local Bazel server and connecting to it...
DEBUG: /root/.cache/bazel/_bazel_root/4aeaa44aa45c9ae450c5a2536656d9b5/external/bazel_tools/tools/cpp/lib_cc_configure.bzl:118:5: 
Auto-Configuration Warning: CC with -fuse-ld=gold returned 0, but its -v output didn't contain 'gold', falling back to the default linker.
INFO: Analyzed target //:test (28 packages loaded, 6580 targets configured).
INFO: Found 1 target...
Target //:test up-to-date:
  bazel-bin/linux_amd64_stripped/test
INFO: Elapsed time: 35.144s, Critical Path: 1.19s
INFO: 3 processes: 3 linux-sandbox.
INFO: Build completed successfully, 7 total actions
INFO: Build completed successfully, 7 total actions
hello world
```

创建

#### 使用gazelle自动生成BUILD.bazel文件

在实际的项目中，里面的`BUILD.bazel`我们肯定是使用工具自动生成的，来看下如何自动生成的  



### 参考
【带你深入AI（6）- 详解bazel】https://blog.csdn.net/u013510838/article/details/80102438   
【bazel文档】https://docs.bazel.build/versions/4.1.0/skylark/concepts.html  
【Bazel 构建 golang 项目】https://zhuanlan.zhihu.com/p/95998597  
【如何评价 Google 开源的 Bazel ？】https://www.zhihu.com/question/29025960  
【使用bazel编译go项目】https://juejin.cn/post/6844903892757528590  