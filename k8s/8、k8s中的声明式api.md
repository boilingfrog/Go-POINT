<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [声明式API](#%E5%A3%B0%E6%98%8E%E5%BC%8Fapi)
  - [声明式和命令式的对比](#%E5%A3%B0%E6%98%8E%E5%BC%8F%E5%92%8C%E5%91%BD%E4%BB%A4%E5%BC%8F%E7%9A%84%E5%AF%B9%E6%AF%94)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 声明式API

### 声明式和命令式的对比

命令式  

命令式有时也称为指令式，命令式的场景下，计算机只会机械的完成指定的命令操作，执行的结果就取决于执行的命令是否正确。    

声明式  

声明式也称为描述式或者申明式，这种方式告诉计算机想要的，由计算机自己去设定执行的路径，需要计算机有一定的`智能`。    

最常见的声明式栗子就是数据库，查询的 sql 就表示我们想要的结果集，数据库运行查询 sql 的时候，会帮我们处理查询，并且返回查询的结果。数据库在查询的时候，会进行索引匹配，做查询优化等处理，再返回数据结果的时候，同时使用最优的查询路径。如果我们自己去处理这些操作，就需要写很多代码了，而不是仅仅通过一行代码就能解决。     

### 参考

【深入剖析 Kubernetes】https://time.geekbang.org/column/intro/100015201?code=UhApqgxa4VLIA591OKMTemuH1%2FWyLNNiHZ2CRYYdZzY%3D  
【k8s 声明式 API】https://www.51cto.com/article/712066.html     





