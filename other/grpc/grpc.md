## gRPC

### 前言

之前这玩意一直没用到过，最近项目中用到了。好好研究下。  


### 为什么使用gRPC

gRPC是Google公司基于Protobuf开发的跨语言的开源RPC框架。gRPC基于HTTP/2协议设计，可以基于一个HTTP/2链接提供多个服务，对于移动设备更加友好。  

#### 传输协议

RPC: 可以使用TCP协议，也可以使用HTTP协议  
HTTP:基于HTTP协议  

#### 传输效率

RPC:使用自定义的TCP协议，或者使用HTTP协议，都可以减少报文的体积，提高传输的效率  
HTTP:如果是基于HTTP1.1的协议，请求中会包含很多无用的内容，如果是基于HTTP2.0，那么简单的封装以下是可以作为一个RPC来使用的，这时标准RPC框架更多的是服务治理  



### 参考

【https://www.cnblogs.com/yilezhu/p/10645804.html】https://www.cnblogs.com/yilezhu/p/10645804.html  
【gRPC 官方文档中文版】https://doc.oschina.net/grpc?t=60133    
【HTTP和RPC的优缺点】https://cloud.tencent.com/developer/article/1353110  