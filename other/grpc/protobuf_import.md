<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [protobuf使用import导入包找不到](#protobuf%E4%BD%BF%E7%94%A8import%E5%AF%BC%E5%85%A5%E5%8C%85%E6%89%BE%E4%B8%8D%E5%88%B0)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [解决方案](#%E8%A7%A3%E5%86%B3%E6%96%B9%E6%A1%88)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## protobuf使用import导入包找不到

### 前言

使用`protobuf`生成go代码，发现`protobuf`中一个import引用找不到  

protobuf代码  

```go
syntax = "proto3";

package main;

import "github.com/mwitkow/go-proto-validators/validator.proto";

message Message {
    string important_string = 1 [
        (validator.field) = {regex: "^[a-z]{2,5}$"}
    ];
    int32 age = 2 [
        (validator.field) = {int_gt: 0, int_lt: 100}
    ];
}
```

生成的时候报错

```go
$ protoc --govalidators_out=. --go_out=plugins=grpc:. hello.proto
github.com/mwitkow/go-proto-validators/validator.proto: File not found.
hello.proto:5:1: Import "github.com/mwitkow/go-proto-validators/validator.proto" was not found or had errors.
```

### 解决方案

我们要弄明白protoc中`proto_path`参数的含义  

- proto_path: 指定了在哪个目录中搜索import中导入的和要编译为.go的proto文件，可以定义多个

所以添加proto_path就可以了，指定两个地址，一个是import的地址，一个是要编译为.go的proto文件的地址  

```go
protoc --proto_path=. --proto_path=/Users/yj/Go/src/ --govalidators_out=. --go_out=plugins=grpc:. hello.proto
```