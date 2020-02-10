## virtual box设置网络，使用nat网络和仅主机(Host Only)网络进行连接

### 前言
作为程序员难免要在本机电脑安装虚拟机，最近在用virtual box安装虚拟机的时候遇到了点问题。
对于虚拟机的网络设置最简单的就是使用桥接网卡的方式，所有的网络场景都能连通。但是也有几个
缺点:1 网络ip不固定，2 当虚拟机的网段和宿主机的网段不在同一个网段的时候就不能使用了。
我也遇到了这些问题，所以就换了一种方式，使用nat网络和仅主机(Host Only)网络组合的方式进
行连接。

### 网络设置

首先我们下来了解下，集中网络的应用场景  
![Aaron Swartz](https://github.com/zhan-liz/Go-POINT/blob/master/img/virtualbox_1.png?raw=true)

##### 我的装机环境  
电脑的系统环境：  
````
$ head -n 1 /etc/issue  
Deepin GNU/Linux 15.11 \n \l
````
软件的版本信息：  
````
Oracle® VM VirtualBox®
User Manual
Version 6.0.8 Edition
````
安装的虚拟机：
````
centos7
````
需求：各个网络的场景全部支持  

#### 全局设置Nat网络

选择管理->全局设定->网络->添加Nat网络
![Aaron Swartz](https://github.com/zhan-liz/Go-POINT/blob/master/img/virtualbox_2.png?raw=true)

#### 添加主机网络管理器

管理->主机网络管理器->新建主机网络
![Aaron Swartz](https://github.com/zhan-liz/Go-POINT/blob/master/img/virtualbox_3.png?raw=true)
注意：DHCP服务不要勾选，我们去添加静态的ip，这样ip就是固定的

#### 设置虚拟机的网络

对应的虚拟机->设置->网络->网卡１设置(选择nat网络)->网卡
2(选择Host Only网络)  
网卡1
![Aaron Swartz](https://github.com/zhan-liz/Go-POINT/blob/master/img/virtualbox_4.png?raw=true)
网卡2  
![Aaron Swartz](https://github.com/zhan-liz/Go-POINT/blob/master/img/virtualbox_5.png?raw=true)

#### 进去虚拟机修改设置Host-only静态IP

![Aaron Swartz](https://github.com/zhan-liz/Go-POINT/blob/master/img/virtualbox_6.png?raw=true)

