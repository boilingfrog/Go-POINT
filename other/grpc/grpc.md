<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [gRPC](#grpc)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [为什么使用gRPC](#%E4%B8%BA%E4%BB%80%E4%B9%88%E4%BD%BF%E7%94%A8grpc)
    - [传输协议](#%E4%BC%A0%E8%BE%93%E5%8D%8F%E8%AE%AE)
    - [传输效率](#%E4%BC%A0%E8%BE%93%E6%95%88%E7%8E%87)
    - [性能消耗](#%E6%80%A7%E8%83%BD%E6%B6%88%E8%80%97)
  - [gRPC入门](#grpc%E5%85%A5%E9%97%A8)
  - [gRPC流](#grpc%E6%B5%81)
  - [发布和订阅模式](#%E5%8F%91%E5%B8%83%E5%92%8C%E8%AE%A2%E9%98%85%E6%A8%A1%E5%BC%8F)
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

### gRPC入门

创建hello.proto文件，定义HelloService接口  

```go
syntax = "proto3";

package main;

message String {
    string value = 1;
}

service HelloService {
    rpc Hello (String) returns (String);
}
```

使用protoc-gen-go内置的gRPC插件生成gRPC代码：

```go
protoc --go_out=plugins=grpc:. hello.proto
```

然后会生成`hello.pb.go`文件，服务端和客户端生成的代码都包含在内  

运行服务端代码  

```go
package main

import (
	"context"
	"daily-test/gRPC"
	"log"
	"net"

	"google.golang.org/grpc"
)

type HelloServiceImpl struct{}

func (p *HelloServiceImpl) Hello(
	ctx context.Context, args *gRPC.String,
) (*gRPC.String, error) {
	reply := &gRPC.String{Value: "hello:" + args.GetValue()}
	return reply, nil
}

func main() {

	grpcServer := grpc.NewServer()
	gRPC.RegisterHelloServiceServer(grpcServer, new(HelloServiceImpl))

	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer.Serve(lis)
}
```

首先是通过grpc.NewServer()构造一个gRPC服务对象，然后通过gRPC插件生成的RegisterHelloServiceServer函数注册我们实现的HelloServiceImpl服务。然后通过grpcServer.Serve(lis)在一个监听端口上提供gRPC服务。  

定义客户端代码  

```go
package main

import (
	"context"
	"daily-test/gRPC"
	"fmt"
	"log"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:1234", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := gRPC.NewHelloServiceClient(conn)
	reply, err := client.Hello(context.Background(), &gRPC.String{Value: "hello"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply.GetValue())
}
```

其中`grpc.Dial`负责和gRPC服务建立链接，然后`NewHelloServiceClient`函数基于已经建立的链接构造`HelloServiceClient`对象。返回的client其实是
一个`HelloServiceClient`接口对象，通过接口定义的方法就可以调用服务端对应的gRPC服务提供的方法。  

### gRPC流

RPC是远程函数调用，因此每次调用的函数参数和返回值不能太大，否则将严重影响每次调用的响应时间。因此传统的RPC方法调用对于上传和下载较大数据量场景
并不适合。同时传统RPC模式也不适用于对时间不确定的订阅和发布模式。为此，gRPC框架针对服务器端和客户端分别提供了流特性。  

在HelloService增加一个支持双向流的Channel方法  

```go
syntax = "proto3";

package main;

message String {
    string value = 1;
}

service HelloService {
    rpc Hello (String) returns (String);

    rpc Channel (stream String) returns (stream String);
}
```

关键字stream指定启用流特性，参数部分是接收客户端参数的流，返回值是返回给客户端的流。  

使用protoc-gen-go内置的gRPC插件生成gRPC代码：

```go
protoc --go_out=plugins=grpc:. hello.proto
```

重新生成代码可以看到接口中新增加的Channel方法的定义：  

```go
type HelloServiceServer interface {
    Hello(context.Context, *String) (*String, error)
    Channel(HelloService_ChannelServer) error
}
type HelloServiceClient interface {
    Hello(ctx context.Context, in *String, opts ...grpc.CallOption) (
        *String, error,
    )
    Channel(ctx context.Context, opts ...grpc.CallOption) (
        HelloService_ChannelClient, error,
    )
}
```

实现流服务  

```go
func (p *HelloServiceImpl) Channel(stream gRPC_stream.HelloService_ChannelServer) error {
	for {
		args, err := stream.Recv()
		if err != nil {
			// io.EOF表示客户端流关闭
			if err == io.EOF {
				return nil
			}
			return err
		}

		reply := &gRPC_stream.String{Value: "hello:" + args.GetValue()}

		err = stream.Send(reply)
		if err != nil {
			return err
		}
	}
}
```

实现服务端代码：  

```go
package main

import (
	"context"
	"daily-test/gRPC_stream"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

type HelloServiceImpl struct{}

func (p *HelloServiceImpl) Hello(
	ctx context.Context, args *gRPC_stream.String,
) (*gRPC_stream.String, error) {
	reply := &gRPC_stream.String{Value: "hello:" + args.GetValue()}
	return reply, nil
}

func (p *HelloServiceImpl) Channel(stream gRPC_stream.HelloService_ChannelServer) error {
	for {
		args, err := stream.Recv()
		if err != nil {
			// io.EOF表示客户端流关闭
			if err == io.EOF {
				return nil
			}
			return err
		}

		reply := &gRPC_stream.String{Value: "hello:" + args.GetValue()}

		err = stream.Send(reply)
		if err != nil {
			return err
		}
	}
}

func main() {

	grpcServer := grpc.NewServer()
	gRPC_stream.RegisterHelloServiceServer(grpcServer, new(HelloServiceImpl))

	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer.Serve(lis)
}
```

客户端代码  

```go
package main

import (
	"context"
	"daily-test/gRPC_stream"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:1234", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := gRPC_stream.NewHelloServiceClient(conn)

	stream, err := client.Channel(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 发送
	go func() {
		for {
			if err := stream.Send(&gRPC_stream.String{Value: "hi"}); err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Second)
		}
	}()

	// 接收
	for {
		reply, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		fmt.Println(reply.GetValue())
	}

}
```

### 发布和订阅模式

### 参考

【https://www.cnblogs.com/yilezhu/p/10645804.html】https://www.cnblogs.com/yilezhu/p/10645804.html  
【gRPC 官方文档中文版】https://doc.oschina.net/grpc?t=60133    
【HTTP和RPC的优缺点】https://cloud.tencent.com/developer/article/1353110  
【gRPC入门】https://chai2010.gitbooks.io/advanced-go-programming-book/content/ch4-rpc/ch4-04-grpc.html  