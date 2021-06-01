<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [记一次http超时引发的事故](#%E8%AE%B0%E4%B8%80%E6%AC%A1http%E8%B6%85%E6%97%B6%E5%BC%95%E5%8F%91%E7%9A%84%E4%BA%8B%E6%95%85)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [分析下具体的代码实现](#%E5%88%86%E6%9E%90%E4%B8%8B%E5%85%B7%E4%BD%93%E7%9A%84%E4%BB%A3%E7%A0%81%E5%AE%9E%E7%8E%B0)
  - [服务设置超时](#%E6%9C%8D%E5%8A%A1%E8%AE%BE%E7%BD%AE%E8%B6%85%E6%97%B6)
  - [客户端设置超时](#%E5%AE%A2%E6%88%B7%E7%AB%AF%E8%AE%BE%E7%BD%AE%E8%B6%85%E6%97%B6)
    - [http.client](#httpclient)
    - [context](#context)
    - [http.Transport](#httptransport)
  - [问题](#%E9%97%AE%E9%A2%98)
    - [如果在客户端在超时的临界点，触发了超时机制，这时候服务端刚好也接收到了，http的请求](#%E5%A6%82%E6%9E%9C%E5%9C%A8%E5%AE%A2%E6%88%B7%E7%AB%AF%E5%9C%A8%E8%B6%85%E6%97%B6%E7%9A%84%E4%B8%B4%E7%95%8C%E7%82%B9%E8%A7%A6%E5%8F%91%E4%BA%86%E8%B6%85%E6%97%B6%E6%9C%BA%E5%88%B6%E8%BF%99%E6%97%B6%E5%80%99%E6%9C%8D%E5%8A%A1%E7%AB%AF%E5%88%9A%E5%A5%BD%E4%B9%9F%E6%8E%A5%E6%94%B6%E5%88%B0%E4%BA%86http%E7%9A%84%E8%AF%B7%E6%B1%82)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 记一次http超时引发的事故

### 前言

我们使用的是golang标准库的`http client`，对于一些http请求，我们在处理的时候，会考虑加上超时时间，防止http请求一直在请求，导致业务长时间阻塞等待。  

最近同事写了一个超时的组件，这几天访问量上来了，网络也出现了波动，造成了接口在报错超时的情况下，还是出现了请求结果的成功。

### 分析下具体的代码实现

```go
type request struct {
	method string
	url    string
	value  string
	ps     *params
}

type params struct {
	timeout     int //超时时间
	retry       int //重试次数
	headers     map[string]string
	contentType string
}

func (req *request) Do(result interface{}) ([]byte, error) {
	res, err := asyncCall(doRequest, req)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return res, nil
	}

	switch req.ps.contentType {
	case "application/xml":
		if err := xml.Unmarshal(res, result); err != nil {
			return nil, err
		}
	default:
		if err := json.Unmarshal(res, result); err != nil {
			return nil, err
		}
	}

	return res, nil
}
type timeout struct {
	data []byte
	err  error
}


func doRequest(request *request) ([]byte, error) {
	var (
		req    *http.Request
		errReq error
	)
	if request.value != "null" {
		buf := strings.NewReader(request.value)
		req, errReq = http.NewRequest(request.method, request.url, buf)
		if errReq != nil {
			return nil, errReq
		}
	} else {
		req, errReq = http.NewRequest(request.method, request.url, nil)
		if errReq != nil {
			return nil, errReq
		}
	}
    // 这里的client没有设置超时时间
    // 所以当下面检测到一次超时的时候，会重新又发起一次请求
    // 但是老的请求其实没有被关闭，一直在执行
	client := http.Client{}
	res, err := client.Do(req)
    ...
}

// 重试调用请求
// 当超时的时候发起一次新的请求
func asyncCall(f func(request *request) ([]byte, error), req *request) ([]byte, error) {
	p := req.ps
	ctx := context.Background()
	done := make(chan *timeout, 1)

	for i := 0; i < p.retry; i++ {
		go func(ctx context.Context) {
			// 发送HTTP请求
			res, err := f(req)
			done <- &timeout{
				data: res,
				err:  err,
			}
		}(ctx)
        // 错误主要在这里
        // 如果超时重试为3，第一次超时了，马上又发起了一次新的请求，但是这里错误使用了超时的退出
        // 具体看上面
		select {
		case res := <-done:
			return res.data, res.err
		case <-time.After(time.Duration(p.timeout) * time.Millisecond):
		}
	}
	return nil, ecode.TimeoutErr
}
```

错误的原因  

1、超时重试，之后过了一段时间没有拿到结果就认为是超超时了，但是http请求没有被关闭；  

2、错误使用了`http`的超时，具体的做法要通过`context`或`http.client`去实现，见下文；  

修改之后的代码

```go
func doRequest(request *request) ([]byte, error) {
	var (
		req    *http.Request
		errReq error
	)
	if request.value != "null" {
		buf := strings.NewReader(request.value)
		req, errReq = http.NewRequest(request.method, request.url, buf)
		if errReq != nil {
			return nil, errReq
		}
	} else {
		req, errReq = http.NewRequest(request.method, request.url, nil)
		if errReq != nil {
			return nil, errReq
		}
	}

   // 这里通过http.Client设置超时时间
	client := http.Client{
		Timeout: time.Duration(request.ps.timeout) * time.Millisecond,
	}
	res, err := client.Do(req)
    ...
}

func asyncCall(f func(request *request) ([]byte, error), req *request) ([]byte, error) {
	p := req.ps
    // 重试的时候只有上一个http请求真的超时了，之后才会发起一次新的请求
	for i := 0; i < p.retry; i++ {
		// 发送HTTP请求
		res, err := f(req)
		// 判断超时
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			continue
		}

		return res, err

	}
	return nil, ecode.TimeoutErr
}
```

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

这种服务端还是可以拿到请求的数据，所以对于超时时间的设置我们需要根据实际情况进行权衡，同时我们要考虑接口的幂等性。  

### 参考

【[译]Go net/http 超时机制完全手册】https://colobu.com/2016/07/01/the-complete-guide-to-golang-net-http-timeouts/  
【Go 语言 HTTP 请求超时入门】https://studygolang.com/articles/14405  