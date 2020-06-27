## login

### Cookie和Session

`http`是无状态的协议，所以服务端需要记录用户的状态时，就需要用某种机制来识具体的用户，这个机制就是`Session`。举个例子，当我们在网上购买商品，然后支付。由于`http`是无状态的，所以支付的时候是不知道到底是那个用户支付的。所以服务端创建了`Session`，用来标示当前操作的用户。
在服务端保存Session的方法很多，内存、数据库、文件都有。集群的时候也要考虑Session的转移，在大型的网站，一般会有专门的`Session`服务器集群，用来保存用户会话，这个时候 `Session` 信息都是放在内存的，使用一些缓存服务比如`Memcached`之类的来放 `Session`。   

对于服务端是如何识别用户呢（客户端）？这时候就需要介绍一下`Cookie`了，实际上大多数的应用都是用`Cookie`来实现`Session`跟踪的。那么`Cookie`好`Session`是如何交互的呢？  、

第一次创建`Session`的时候，服务端会在`HTTP`协议中告诉客户端，需要在`Cookie`里面记录一个`Session ID`，以后每次请求把这个会话ID发送到服务器，我就知道你是谁了。   

![bufio](images/login-session.png?raw=true)
