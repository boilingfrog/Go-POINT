## virtual box设置网络，使用nat网络和仅主机(Host Only)网络进行连接

### 前言
作为程序员难免要在本机电脑安装虚拟机，最近在用virtual box安装虚拟机的时候遇到了点问题。
对于虚拟机的网络设置最简单的就是使用桥接网卡的方式，所有的网络场景都能连通。但是也有几个
缺点:1 网络ip不固定，2 当虚拟机的网段和宿主机的网段不在同一个网段的时候就不能使用了。
我也遇到了这些问题，所以就换了一种方式，使用at网络和仅主机(Host Only)网络组合的方式进
行连接。

### 网络设置

首先我们下来了解下，集中网络的应用场景  
![Aaron Swartz](https://github.com/zhan-liz/Go-POINT/blob/master/img/virtualbox_1.png?raw=true)
