## docker配合Nginx部署go应用 

### Nginx

什么是`Nginx`？  

`Nginx`是一个高性能的`HTTP`和反向代理服务器，也是一个`IMAP/POP3/SMTP`代理服务器。  

`Nginx`是一款轻量级的`Web`服务器/反向代理服务器以及电子邮件代理服务器，并在一个`BSD-like`协议下发行。由俄罗斯的程序设计师`lgor Sysoev`所开发，供俄国大型的入口网站及搜索引擎Rambler使用。其特点是占有内存少，并发能力强，事实上`nginx`的并发能力确实在同类型的网页服务器中表现较好。  

`Nginx`相较于`Apache\lighttpd`具有占有内存少，稳定性高等优势，并且依靠并发能力强，丰富的模块库以及友好灵活的配置而闻名。在`Linux`操作系统下，`nginx`使用`epoll`事件模型,得益于此，nginx在Linux操作系统下效率相当高。同时`Nginx`在`OpenBSD`或`FreeBSD`操作系统上采用类似于`Epoll`的高效事件模型`kqueue`.  
