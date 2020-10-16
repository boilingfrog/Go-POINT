<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [gRPC](#grpc)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [为什么使用gRPC](#%E4%B8%BA%E4%BB%80%E4%B9%88%E4%BD%BF%E7%94%A8grpc)
    - [传输协议](#%E4%BC%A0%E8%BE%93%E5%8D%8F%E8%AE%AE)
    - [传输效率](#%E4%BC%A0%E8%BE%93%E6%95%88%E7%8E%87)
    - [性能消耗](#%E6%80%A7%E8%83%BD%E6%B6%88%E8%80%97)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## gRPC

### 前言

之前这玩意一直没用到过，最近项目中用到了。好好研究下。  


### 为什么使用gRPC

gRPC是Google公司基于Protobuf开发的跨语言的开源RPC框架。gRPC基于HTTP/2协议设计，可以基于一个HTTP/2链接提供多个服务，对于移动设备更加友好。  

#### 传输协议

gRPC: 可以使用TCP协议，也可以使用HTTP协议  
HTTP:基于HTTP协议  

#### 传输效率

gRPC:使用自定义的TCP协议，或者使用HTTP协议，都可以减少报文的体积，提高传输的效率  
HTTP:如果是基于HTTP1.1的协议，请求中会包含很多无用的内容，如果是基于HTTP2.0，那么简单的封装以下是可以作为一个RPC来使用的，这时标准RPC框架更多的是服务治理  

#### 性能消耗

gRPC:gRPC消息使用一种有效的二进制消息格式protobuf进行序列化。Protobuf在服务器和客户机上的序列化非常快。Protobuf序列化后的消息体积很小，
能够有效负载，在移动应用程序等有限带宽场景中显得很重要  
HTTP:大部分是通过json来实现的，字节大小和序列化耗时都比Protobuf要更消耗性能  


gRPC主要用于公司内部的服务调用，性能消耗低，传输效率高。HTTP主要用于对外的异构环境，浏览器接口调用，APP接口调用，第三方接口调用等。


### 参考

【https://www.cnblogs.com/yilezhu/p/10645804.html】https://www.cnblogs.com/yilezhu/p/10645804.html  
【gRPC 官方文档中文版】https://doc.oschina.net/grpc?t=60133    
【HTTP和RPC的优缺点】https://cloud.tencent.com/developer/article/1353110  