<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Bazel使用了解](#bazel%E4%BD%BF%E7%94%A8%E4%BA%86%E8%A7%A3)
  - [什么是Bazel](#%E4%BB%80%E4%B9%88%E6%98%AFbazel)
  - [Bazel产生的背景](#bazel%E4%BA%A7%E7%94%9F%E7%9A%84%E8%83%8C%E6%99%AF)
  - [Bazel的作用](#bazel%E7%9A%84%E4%BD%9C%E7%94%A8)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Bazel使用了解

### 什么是Bazel

`Bazel`是一个支持多语言、跨平台的构建工具。`Bazel`支持任意大小的构建目标，并支持跨多个仓库的构建，是`Google`主推的一种构建工具。  

`bazel`优点很多，主要有  

- 构建快。支持增量编译。对依赖关系进行了优化，从而支持并发执行。

- 可构建多种语言。bazel可用来构建Java C++ Android ios等很多语言和框架，并支持mac windows linux等不同平台

- 可伸缩。可处理任意大小的代码库，可处理多个库，也可以处理单个库

- 可扩展。使用bazel扩展语言可支持新语言和新平台。

### Bazel产生的背景

1、开源成为当前软件开发的主旋律。哪怕你是商业软件，也逃离不了社区的包围。如何方便地获取依赖，并做到平滑升级很重要。如果构建工具能够很方便地获取源代码，那就太好了。  

2、混合多语言编程成为一种选择。每种语言都有自己适用的场景，但是构建多语言的软件系统非常具有挑战性。例如，Python社区很喜欢搭配C/C++，高性能计算扔个Ｃ/C++，Python提供编程接口。如果构建工具能够无缝支持多语言构建，真的很方便。  

3、代码复用。我只想复用第三方的一个头文件，而不是整个系统。拒绝拷贝是优秀程序员的基本素养，如果构建工具能帮我方便地获取到所依赖的组件，剔除不必要的依赖，那就太完美了。  

4、增量构建。当只修改了一行代码，构建系统能够准确计算需要构建的依赖目标，而不是全构建；否则生命都浪费在编译上了。  

5、云构建。大型软件公司，复用计算资源，可以带来巨大的收益。  

### Bazel的作用



### 参考
【带你深入AI（6）- 详解bazel】https://blog.csdn.net/u013510838/article/details/80102438   
【bazel文档】https://docs.bazel.build/versions/4.1.0/skylark/concepts.html  
【Bazel 构建 golang 项目】https://zhuanlan.zhihu.com/p/95998597  
【如何评价 Google 开源的 Bazel ？】https://www.zhihu.com/question/29025960  