<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [记一次http超时引发的事故](#%E8%AE%B0%E4%B8%80%E6%AC%A1http%E8%B6%85%E6%97%B6%E5%BC%95%E5%8F%91%E7%9A%84%E4%BA%8B%E6%95%85)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [服务设置超时](#%E6%9C%8D%E5%8A%A1%E8%AE%BE%E7%BD%AE%E8%B6%85%E6%97%B6)
  - [客户端设置超时](#%E5%AE%A2%E6%88%B7%E7%AB%AF%E8%AE%BE%E7%BD%AE%E8%B6%85%E6%97%B6)
    - [http.client](#httpclient)
    - [context](#context)
    - [http.Transport](#httptransport)
  - [问题](#%E9%97%AE%E9%A2%98)
    - [如果在客户端在超时的临界点，触发了超时机制，这时候服务端刚好也接收到了，http的请求](#%E5%A6%82%E6%9E%9C%E5%9C%A8%E5%AE%A2%E6%88%B7%E7%AB%AF%E5%9C%A8%E8%B6%85%E6%97%B6%E7%9A%84%E4%B8%B4%E7%95%8C%E7%82%B9%E8%A7%A6%E5%8F%91%E4%BA%86%E8%B6%85%E6%97%B6%E6%9C%BA%E5%88%B6%E8%BF%99%E6%97%B6%E5%80%99%E6%9C%8D%E5%8A%A1%E7%AB%AF%E5%88%9A%E5%A5%BD%E4%B9%9F%E6%8E%A5%E6%94%B6%E5%88%B0%E4%BA%86http%E7%9A%84%E8%AF%B7%E6%B1%82)
  - [分下下http client超时的处理方法](#%E5%88%86%E4%B8%8B%E4%B8%8Bhttp-client%E8%B6%85%E6%97%B6%E7%9A%84%E5%A4%84%E7%90%86%E6%96%B9%E6%B3%95)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 记一次http超时引发的事故

### 前言

我们使用的是golang标准库的`http client`，对于一些http请求，我们在处理的时候，会考虑加上超时时间，防止http请求一直在请求，导致业务长时间阻塞等待。  

### 服务设置超时

`http.Server`有两个设置超时的方法:  

- ReadTimeout 

- WriteTimeout




### 客户端设置超时

#### http.client

最简单的我们通过`http.Client`的`Timeout`字段，就可以实现客户端的超时控制  

`http.client`超时是超时的高层实现，包含了从`Dial`到`Response Body`的整个请求流程。`http.client`的实现提供了一个结构体类型可以接受一个额外的`time.Duration`类型的`Timeout`属性。这个参数定义了从请求开始到响应消息体被完全接收的时间限制。    

```go
func httpClientTimeout() {
	c := &http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := c.Get("http://127.0.0.1:8081/test")
	fmt.Println(resp)
	fmt.Println(err)
}
```

#### context  

`net/http`中的`request`实现了`context`,所以我们可以借助于context本身的超时机制，实现`http`中`request`的超时处理  

```go
func contextTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequest("GET", "http://127.0.0.1:8081/test", nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	fmt.Println(resp)
	fmt.Println(err)
}
```

使用`context`的优点就是，当父`context`被取消时，子`context`就会层层退出。   

#### http.Transport

通过`Transport`还可以进行一些更小维度的超时设置  

- net.Dialer.Timeout 限制建立TCP连接的时间
- http.Transport.TLSHandshakeTimeout 限制 TLS握手的时间
- http.Transport.ResponseHeaderTimeout 限制读取response header的时间
- http.Transport.ExpectContinueTimeout 限制client在发送包含 Expect: 100-continue的header到收到继续发送body的response之间的时间等待。注意在1.6中设置这个值会禁用HTTP/2(DefaultTransport自1.6.2起是个特例)


```go
func transportTimeout() {
	transport := &http.Transport{
		DialContext:           (&net.Dialer{}).DialContext,
		ResponseHeaderTimeout: 3 * time.Second,
	}

	c := http.Client{Transport: transport}

	resp, err := c.Get("http://127.0.0.1:8081/test")
	fmt.Println(resp)
	fmt.Println(err)
}
```

### 问题

#### 如果在客户端在超时的临界点，触发了超时机制，这时候服务端刚好也接收到了，http的请求  



### 分下下http client超时的处理方法









### 参考

【[译]Go net/http 超时机制完全手册】https://colobu.com/2016/07/01/the-complete-guide-to-golang-net-http-timeouts/  
【Go 语言 HTTP 请求超时入门】https://studygolang.com/articles/14405  