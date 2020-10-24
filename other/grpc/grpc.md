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
  - [证书认证](#%E8%AF%81%E4%B9%A6%E8%AE%A4%E8%AF%81)
    - [使用根证书](#%E4%BD%BF%E7%94%A8%E6%A0%B9%E8%AF%81%E4%B9%A6)
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

### 证书认证

gRPC建立在HTTP/2协议之上，对TLS提供了很好的支持。我们使用公钥，私钥，实现最近本的token认证  

首先看下代码目录结构

```go
gRPC_cert
├── cert
│   ├── cert.conf
│   ├── server.crt
│   └── server.key
├── client
│   └── main.go
├── hello.pb.go
├── hello.proto
└── server
    └── main.go
```

创建`cert.conf`  

```go
[ req ]
default_bits = 2048
prompt = no
default_md = sha256
req_extensions = req_ext
distinguished_name = dn
 
[ dn ]
C = UK
ST = London
L = London
O = liz Ltd.
OU = Information Technologies
emailAddress = email@email.com
CN = localhost
 
[ req_ext ]
subjectAltName = @alt_names
 
[ alt_names ]
DNS.1 = localhost
```

生成`private key`

```go
$ openssl genrsa -out cert/server.key 2048
 
$ openssl req -nodes -new -x509 -sha256 -days 1825 -config cert/cert.conf -extensions 'req_ext' -key cert/server.key -out cert/server.crt
```

启动gRPC服务端的时候传入证书的参数选项  

```go
package main

import (
	"context"
	"daily-test/gRPC_cert"
	"log"
	"net"

	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
)

type HelloServiceImpl struct{}

func (p *HelloServiceImpl) Hello(
	ctx context.Context, args *gRPC_cert.String,
) (*gRPC_cert.String, error) {
	reply := &gRPC_cert.String{Value: "hello:" + args.GetValue()}
	return reply, nil
}

func main() {

	creds, err := credentials.NewServerTLSFromFile("./gRPC_cert/cert/server.crt", "./gRPC_cert/cert/server.key")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	gRPC_cert.RegisterHelloServiceServer(grpcServer, new(HelloServiceImpl))

	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer.Serve(lis)
}
```

credentials.NewServerTLSFromFile函数是从文件为服务器构造证书对象，然后通过grpc.Creds(creds)函数将证书包装为选项后作为参数传入grpc.NewServer函数。  

客户端基于服务端的证书和服务端名字进行对服务端的认证校验  

```go
package main

import (
	"context"
	"daily-test/gRPC_cert"
	"fmt"
	"log"

	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
)

func main() {
	// 带入证书的信息
	creds, err := credentials.NewClientTLSFromFile(
		"./gRPC_cert/cert/server.crt", "localhost",
	)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial("localhost:1234",
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := gRPC_cert.NewHelloServiceClient(conn)
	reply, err := client.Hello(context.Background(), &gRPC_cert.String{Value: "hello"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply.GetValue())
}
```

demo地址`https://github.com/boilingfrog/daily-test/tree/master/gRPC_cert`  

这种需要讲服务端的证书进行下发，这样客户在交互的时候每次都需要带过来，但是这样是不安全的。在传输的过程中证书存在被
监听和替换的可能性。  

可以引入根证书，通过对服务端和客户端来进行签名，来保证安全。  

#### 使用根证书

为了避免证书的传递过程中被篡改，可以通过一个安全可靠的根证书分别对服务器和客户端的证书进行签名。这样客户端或服务器在收到对方的证书后可以通过根
证书进行验证证书的有效性。  

首先看下我的项目目录  

```go
gRPC_cert_ca
├── cert
│   ├── ca.key
│   ├── ca.pem
│   ├── cert.conf
│   ├── client
│   │   ├── client.csr
│   │   ├── client.key
│   │   └── client.pem
│   ├── server
│   │   ├── server.csr
│   │   ├── server.key
│   │   └── server.pem
│   ├── server.crt
│   └── server.key
├── client
│   └── main.go
├── hello.pb.go
├── hello.proto
└── server
    └── main.go
```

**生成根证书**  

公钥

````
openssl genrsa -out ca.key 2048
````

秘钥

```go
openssl req -new -x509 -days 7200 -key ca.key -out ca.pem
```

生成秘钥的时候需要填写信息  

```go
Country Name (2 letter code) []:
State or Province Name (full name) []:
Locality Name (eg, city) []:
Organization Name (eg, company) []:
Organizational Unit Name (eg, section) []:
Common Name (eg, fully qualified host name) []:localhost
Email Address []:
```

`Common Name`这个我们需要注意下，这个是主机名，测试的我就放了`localhost`  

**server端**

生成 Key  

```go
openssl ecparam -genkey -name secp384r1 -out server.key
```

生成CSR

```go
openssl req -new -key server.key -out server.csr
```

需要填写信息

```go
Country Name (2 letter code) []:
State or Province Name (full name) []:
Locality Name (eg, city) []:
Organization Name (eg, company) []:
Organizational Unit Name (eg, section) []:
Common Name (eg, fully qualified host name) []:localhost
Email Address []:

Please enter the following 'extra' attributes
to be sent with your certificate request
A challenge password []:
```

使用ca进行证书的签发

```go
openssl x509 -req -sha256 -CA ca.pem -CAkey ca.key -CAcreateserial -days 3650 -in server.csr -out server.pem
```

注意下自己的ca公钥秘钥的路径问题  

**client**

生成key

```go
openssl ecparam -genkey -name secp384r1 -out client.key
```

生成csr

```go
openssl req -new -key client.key -out client.csr
```

使用ca进行证书的签发  

```go
openssl x509 -req -sha256 -CA ca.pem -CAkey ca.key -CAcreateserial -days 3650 -in client.csr -out client.pem
```

到此完成证书的生成   

server端的代码  

```go
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"daily-test/gRPC_cert_ca"
	"io/ioutil"
	"log"
	"net"

	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
)

type HelloServiceImpl struct{}

func (p *HelloServiceImpl) Hello(
	ctx context.Context, args *gRPC_cert_ca.String,
) (*gRPC_cert_ca.String, error) {
	reply := &gRPC_cert_ca.String{Value: "hello:" + args.GetValue()}
	return reply, nil
}

func main() {

	cert, err := tls.LoadX509KeyPair("./gRPC_cert_ca/cert/server/server.pem", "./gRPC_cert_ca/cert/server/server.key")
	if err != nil {
		log.Fatalf("tls.LoadX509KeyPair err: %v", err)
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("./gRPC_cert_ca/cert/ca.pem")
	if err != nil {
		log.Fatalf("ioutil.ReadFile err: %v", err)
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("certPool.AppendCertsFromPEM err")
	}

	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	})

	grpcServer := grpc.NewServer(grpc.Creds(c))
	gRPC_cert_ca.RegisterHelloServiceServer(grpcServer, new(HelloServiceImpl))

	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer.Serve(lis)
}
```

服务器端使用credentials.NewTLS函数生成证书，通过ClientCAs选择CA根证书，并通过ClientAuth选项启用对客户端进行验证。  

客户端代码  

```go
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"daily-test/gRPC_cert_ca"
	"fmt"
	"io/ioutil"
	"log"

	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
)

func main() {
	cert, err := tls.LoadX509KeyPair("./gRPC_cert_ca/cert/client/client.pem", "./gRPC_cert_ca/cert/client/client.key")
	if err != nil {
		log.Fatalf("tls.LoadX509KeyPair err: %v", err)
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("./gRPC_cert_ca/cert/ca.pem")
	if err != nil {
		log.Fatalf("ioutil.ReadFile err: %v", err)
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("certPool.AppendCertsFromPEM err")
	}

	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   "localhost",
		RootCAs:      certPool,
	})

	conn, err := grpc.Dial(":1234",
		grpc.WithTransportCredentials(c),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := gRPC_cert_ca.NewHelloServiceClient(conn)
	reply, err := client.Hello(context.Background(), &gRPC_cert_ca.String{Value: "hello"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply.GetValue())

}
```

在`credentials.NewTLS`函数调用中，客户端通过引入一个CA根证书和服务器的名字来实现对服务器进行验证。客户端在链接服务器时会首先请求服务器的证书，
然后使用CA根证书对收到的服务器端证书进行验证。  

ca认证的demo`https://github.com/boilingfrog/daily-test/tree/master/gRPC_cert_ca`


### 参考

【https://www.cnblogs.com/yilezhu/p/10645804.html】https://www.cnblogs.com/yilezhu/p/10645804.html  
【gRPC 官方文档中文版】https://doc.oschina.net/grpc?t=60133    
【HTTP和RPC的优缺点】https://cloud.tencent.com/developer/article/1353110  
【gRPC入门】https://chai2010.gitbooks.io/advanced-go-programming-book/content/ch4-rpc/ch4-04-grpc.html  
【Using TLS/SSL certificates for gRPC client and server communications in Golang - Updated】http://www.inanzzz.com/index.php/post/jo4y/using-tls-ssl-certificates-for-grpc-client-and-server-communications-in-golang-updated  