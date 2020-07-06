<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [login](#login)
  - [Cookie和Session](#cookie%E5%92%8Csession)
  - [Token](#token)
    - [总结下token的优点](#%E6%80%BB%E7%BB%93%E4%B8%8Btoken%E7%9A%84%E4%BC%98%E7%82%B9)
    - [Acesss Token和Refresh Token](#acesss-token%E5%92%8Crefresh-token)
  - [JWT](#jwt)
  - [OAuth2.0](#oauth20)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## login

### Cookie和Session

`http`是无状态的协议，所以服务端需要记录用户的状态时，就需要用某种机制来识具体的用户，这个机制就是`Session`。举个例子，当我们在网上购买商品，然后支付。由于`http`是无状态的，所以支付的时候是不知道到底是那个用户支付的。所以服务端创建了`Session`，用来标示当前操作的用户。
在服务端保存Session的方法很多，内存、数据库、文件都有。集群的时候也要考虑Session的转移，在大型的网站，一般会有专门的`Session`服务器集群，用来保存用户会话，这个时候 `Session` 信息都是放在内存的，使用一些缓存服务比如`Memcached`之类的来放 `Session`。   

对于服务端是如何识别用户呢（客户端）？这时候就需要介绍一下`Cookie`了，实际上大多数的应用都是用`Cookie`来实现`Session`跟踪的。那么`Cookie`好`Session`是如何交互的呢？  、

第一次创建`Session`的时候，服务端会在`HTTP`协议中告诉客户端，需要在`Cookie`里面记录一个`Session ID`，以后每次请求把这个会话ID发送到服务器，我就知道你是谁了。   

![login](images/login-session.png?raw=true)

总结下：  
`Session`是在服务端保存的一个数据结构，用来跟踪用户的状态，这个数据可以保存在集群、数据库、文件中；  
`Cookie`是客户端保存用户信息的一种机制，用来记录用户的一些信息，也是实现`Session`的一种方式。  

两者的区别：  

- 安全性：`Session` 比 `Cookie` 安全，`Session` 是存储在服务器端的，`Cookie` 是存储在客户端的。  
- 存取值的类型不同：`Cookie` 只支持存字符串数据，想要设置其他类型的数据，需要将其转换成字符串，`Session` 可以存任意数据类型。  
- 有效期不同：`Cookie` 可设置为长时间保持，比如我们经常使用的默认登录功能，`Session` 一般失效时间较短，客户端关闭（默认情况下）或者 `Session` 超时都会失效。  
- 存储大小不同：单个`Cookie`保存的数据不能超过`4K`，`Session` 可存储数据远高于 `Cookie`，但是当访问量过多，会占用过多的服务器资源。  

### Token

在基于`token`实现用户登录的方案中：用户第一次登录，服务端根据用户的信息生成`Token`。然后发放给客户端，客户端保存好`Token`，然后在之后的访问中，都带上`Token`的信息。服务端根据`Token`解析，判断当前的用户信息。  

![login](images/login-token-process.png?raw=true)

当然，裸漏的`Token`是很容易被非法持有和篡改的。如何优化呢？  

可以对数据做个签名。 比如说我用`HMAC-SHA256`算法，加上一个只有我才知道的密钥，对数据做一个签名，把这个签名和数据一起作为`token`，由于密钥别人不知道，就无法伪造`token`了。  

![login](images/login-token.png?raw=true)

客户端在第二次请求时候，带上之前服务端返回的`token`信息。服务端就用这个`token`来校验用户的状态。拿到`token`使用签发时相同的算法和秘钥，重新生成签名，和用户传递的`token`里面的签名作比较，如果匹配的上，表示`token`正确这个用户当前处于登陆状态，匹配不上，就提示重新登陆。因为秘钥和签名的算法是用户是不可见的，这样就避免了`token`被伪造了。  

![login](images/login-token-sign.png?raw=true)

当然，如果一个`token`被劫持了，里面的信息没有被更改，那么还是可以验证通过的，这是不可避免的。  

#### 总结下token的优点

- 无状态、可扩展
- 支持移动设备
- 跨程序调用
- 安全

#### Acesss Token和Refresh Token

常用的还是`Acesss Token`配合`Refresh Token`来实现token的认证。  

`Acesss Token`就是上面介绍的token。不过一般`Acesss Token`的过期时间比较短，所以一般配合`Refresh Token`。当`Acesss Token`过期了，通过`Refresh Token`来换取新的`Acesss Token`。避免用户频繁的登录操作。  

![login](images/login-session-refresh.png?raw=true)

一般`Acesss Token`和`Refresh Token`在签发的时候会存储在服务端，当通过`Refresh Token`换取`Acesss Token`的时候也是会更新`Acesss Token`的值。 

### JWT

前面的两种读需要在服务端存储，JWT是服务端不存储任何信息，在每次请求的时候进行校验。  

详情参见`https://www.ruanyifeng.com/blog/2018/07/json_web_token-tutorial.html`这个已经描述的很清晰了。  

总结下：大致的认证流程如上图，生成`JWT`的`Signature`需要配合一个密匙，这个密匙只有服务器才知道，通过`header`中的算法生成签名。 当校验`JWT`的正确性也是，通过秘钥和算法对传递进来的`JWT`中的`Header`和`Payload`部分重新生成签名和之前的签名作比对，来匹配`JWT`的正确性。  

![login](images/login-jwt.png?raw=true)

当然`JWT`也并不是完美的：  

- JWT 的最大缺点是，由于服务器不保存 session 状态，因此无法在使用过程中废止某个 token，或者更改 token 的权限。也就是说，一旦 JWT 签发了，在到期之前就会始终有效，除非服务器部署额外的逻辑。
- JWT 本身包含了认证信息，一旦泄露，任何人都可以获得该令牌的所有权限。为了减少盗用，JWT 的有效期应该设置得比较短。对于一些比较重要的权限，使用时应该再次对用户进行认证。


### OAuth2.0

什么是`OAuth2.0`？  

OAuth 协议为用户资源的授权提供了一个安全的、开放而又简易的标准。与以往的授权方式不同之处是 OAuth的授权不会使第三方触及到用户的帐号信息（如用户名与密码），即第三方无需使用用户的用户名与密码就可以申请获得该用户资源的授权，因此 OAuth是安全的。OAuth 是 `Open Authorization` 的简写。  

目前版本是2.0。所以我们都喊OAuth2.0  

`OAuth 2.0`对定了四种令牌交互的方式：

- 授权码（authorization-code）
- 隐藏式（implicit）
- 密码式（password）
- 客户端凭证（client credentials）  

#### 授权码方式

授权码方式是最完善，流程最严密的授权模式。  

上面是具体的流程图，为了描述流程