## 记一次http超时引发的事故

### 前言

我们使用的是golang标准库的`http client`，对于一些http请求，我们在处理的时候，会考虑加上超时时间，防止http请求一直在请求，导致业务长时间阻塞等待。  

### 分下下http client超时的处理方法







### 参考

【[译]Go net/http 超时机制完全手册】https://colobu.com/2016/07/01/the-complete-guide-to-golang-net-http-timeouts/  