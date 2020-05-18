## docker配合Nginx部署go应用 

### Nginx

![](https://img2020.cnblogs.com/blog/1237626/202005/1237626-20200518092739061-1669908595.png)  

什么是`Nginx`？  

`Nginx` 是一个 `Web Server`，可以用作反向代理、负载均衡、邮件代理、`TCP / UDP、HTT`P 服务器等等，它拥有很多吸引人的特性，例如：  

- 以较低的内存占用率处理 10,000 多个并发连接（每10k非活动HTTP保持活动连接约2.5 MB ）
- 静态服务器（处理静态文件）
- 正向、反向代理
- 负载均衡
- 通过OpenSSL 对 TLS / SSL 与 SNI 和 OCSP 支持
- FastCGI、SCGI、uWSGI 的支持
- WebSockets、HTTP/1.1 的支持
- Nginx + Lua

### 名词解释

#### 正向代理

![](https://img2020.cnblogs.com/blog/1237626/202005/1237626-20200518100228680-455394848.png)


#### 反向代理

![](https://img2020.cnblogs.com/blog/1237626/202005/1237626-20200518113332214-234901294.png)



### 参考 

【正向代理与反向代理【总结】】https://www.cnblogs.com/anker/p/6056540.html  