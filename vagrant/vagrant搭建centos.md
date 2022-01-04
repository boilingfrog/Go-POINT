<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [vagrant搭建centos](#vagrant%E6%90%AD%E5%BB%BAcentos)
  - [什么是vagrant](#%E4%BB%80%E4%B9%88%E6%98%AFvagrant)
  - [如何使用](#%E5%A6%82%E4%BD%95%E4%BD%BF%E7%94%A8)
    - [1、构建本地的目录](#1%E6%9E%84%E5%BB%BA%E6%9C%AC%E5%9C%B0%E7%9A%84%E7%9B%AE%E5%BD%95)
    - [2、官方下载对应的镜像文件，官方下载地址](#2%E5%AE%98%E6%96%B9%E4%B8%8B%E8%BD%BD%E5%AF%B9%E5%BA%94%E7%9A%84%E9%95%9C%E5%83%8F%E6%96%87%E4%BB%B6%E5%AE%98%E6%96%B9%E4%B8%8B%E8%BD%BD%E5%9C%B0%E5%9D%80)
    - [3、导入刚刚下载的镜像(box文件)](#3%E5%AF%BC%E5%85%A5%E5%88%9A%E5%88%9A%E4%B8%8B%E8%BD%BD%E7%9A%84%E9%95%9C%E5%83%8Fbox%E6%96%87%E4%BB%B6)
    - [4、初始化](#4%E5%88%9D%E5%A7%8B%E5%8C%96)
    - [5、修改Vagrantfile中的box名称](#5%E4%BF%AE%E6%94%B9vagrantfile%E4%B8%AD%E7%9A%84box%E5%90%8D%E7%A7%B0)
    - [6、启动](#6%E5%90%AF%E5%8A%A8)
    - [7、登入](#7%E7%99%BB%E5%85%A5)
  - [同时构建多台](#%E5%90%8C%E6%97%B6%E6%9E%84%E5%BB%BA%E5%A4%9A%E5%8F%B0)
    - [修改Vagrantfile](#%E4%BF%AE%E6%94%B9vagrantfile)
    - [启动](#%E5%90%AF%E5%8A%A8)
  - [vagrant中的网络](#vagrant%E4%B8%AD%E7%9A%84%E7%BD%91%E7%BB%9C)
    - [私有网络](#%E7%A7%81%E6%9C%89%E7%BD%91%E7%BB%9C)
    - [公有网络](#%E5%85%AC%E6%9C%89%E7%BD%91%E7%BB%9C)
  - [常用的命令](#%E5%B8%B8%E7%94%A8%E7%9A%84%E5%91%BD%E4%BB%A4)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## vagrant搭建centos

### 什么是vagrant

> Vagrant 是一个简单易用的部署工具，用英文说应该是 Orchestration Tool 。它能帮助开发人员迅速的构建一个开发环境，帮助测试人员构建测试环境，Vagrant 基于 Ruby 开发，使用开源 VirtualBox 作为虚拟化支持，可以轻松的跨平台部署。

### 如何使用

#### 1、构建本地的目录

```go
 /Users/yj/vagrant/centos7
```

#### 2、官方下载对应的镜像文件，[官方下载地址](http://www.vagrantbox.es/)  

```go
MacBook-Pro-3:centos7 yj$ wget https://github.com/CommanderK5/packer-centos-template/releases/download/0.7.2/vagrant-centos-7.2.box
```

#### 3、导入刚刚下载的镜像(box文件)

```go
MacBook-Pro-3:centos7 yj$ vagrant box add centos7.2 /Users/yj/vagrant/centos7/vagrant-centos-7.2.box 
==> vagrant: A new version of Vagrant is available: 2.2.15 (installed version: 2.2.14)!
==> vagrant: To upgrade visit: https://www.vagrantup.com/downloads.html

==> box: Box file was not detected as metadata. Adding it directly...
==> box: Adding box 'centos7.2' (v0) for provider: 
    box: Unpacking necessary files from: file:///Users/yj/vagrant/centos7/vagrant-centos-7.2.box
==> box: Successfully added box 'centos7.2' (v0) for 'virtualbox'!
```

#### 4、初始化

```go
MacBook-Pro-3:centos7 yj$ vagrant init
```

这时候当前目录会生成一个`Vagrantfile`文件

#### 5、修改Vagrantfile中的box名称

```go
config.vm.box = "centos7-1"
```

#### 6、启动

```go
MacBook-Pro-3:centos7 yj$ vagrant up
Bringing machine 'default' up with 'virtualbox' provider...
==> default: Importing base box 'centos7-1'...
==> default: Matching MAC address for NAT networking...
==> default: Setting the name of the VM: centos7_default_1619487768038_36727
==> default: Fixed port collision for 22 => 2222. Now on port 2200.
==> default: Clearing any previously set network interfaces...
==> default: Preparing network interfaces based on configuration...
    default: Adapter 1: nat
==> default: Forwarding ports...
    default: 22 (guest) => 2200 (host) (adapter 1)
==> default: Booting VM...
==> default: Waiting for machine to boot. This may take a few minutes...
    default: SSH address: 127.0.0.1:2200
    default: SSH username: vagrant
    default: SSH auth method: private key
    default: Warning: Remote connection disconnect. Retrying...
    default: Warning: Connection reset. Retrying...
    default: 
    default: Vagrant insecure key detected. Vagrant will automatically replace
    default: this with a newly generated keypair for better security.
    default: 
    default: Inserting generated public key within guest...
    default: Removing insecure key from the guest if it's present...
    default: Key inserted! Disconnecting and reconnecting using new SSH key...
==> default: Machine booted and ready!
==> default: Checking for guest additions in VM...
    default: The guest additions on this VM do not match the installed version of
    default: VirtualBox! In most cases this is fine, but in rare cases it can
    default: prevent things such as shared folders from working properly. If you see
    default: shared folder errors, please make sure the guest additions within the
    default: virtual machine match the version of VirtualBox you have installed on
    default: your host and reload your VM.
    default: 
    default: Guest Additions Version: 5.0.14
    default: VirtualBox Version: 6.1
==> default: Mounting shared folders...
    default: /vagrant => /Users/yj/vagrant/centos7
```

#### 7、登入

可直接只用`vagrant ssh`登入

```go
MacBook-Pro-3:centos7 yj$ vagrant ssh
Last failed login: Mon Apr 26 22:52:26 BRT 2021 from 10.0.2.2 on ssh:notty
There were 5 failed login attempts since the last successful login.
Last login: Mon Apr 26 22:50:07 2021 from 10.0.2.2
[vagrant@localhost ~]$ 
```

也可以使用ssh

```go
$ ssh -p 2200 root@127.0.0.1
```

上面启动的时候已经告诉我们地址和端口了

<img src="/img/vagrant_1.jpg" alt="vagrant" align=center />

```
账号:root  
密码:vagrant
```

### 同时构建多台

#### 修改Vagrantfile

修改之前产生的`Vagrantfile`文件为  

```go
Vagrant.configure("2") do |config|
  
  config.vm.define "centos7-1" do |vb|
      config.vm.provider "virtualbox" do |v|
      v.memory = 1024
      v.cpus = 1
    end
  vb.vm.host_name = "centos7-1"
  vb.vm.network "private_network", ip: "192.168.56.111"
  vb.vm.box = "centos7.2"
  end

  config.vm.define "centos7-2" do |vb1|
      config.vm.provider "virtualbox" do |v|
      v.memory = 1024
      v.cpus = 1
    end
  vb1.vm.host_name = "centos7-2"
  vb1.vm.network "private_network", ip: "192.168.56.112"
  vb1.vm.box = "centos7.2"
  end

  config.vm.define "centos7-3" do |vb2|
      config.vm.provider "virtualbox" do |v|
      v.memory = 1024
      v.cpus = 1
    end
  vb2.vm.host_name = "centos7-3"
  vb2.vm.network "private_network", ip: "192.168.56.113"
  vb2.vm.box = "centos7.2"
  end
end
``` 

网络使用的是私有网络，私有网络和公有网络区别可以看下文  

#### 启动  

```go
MacBook-Pro-3:centos7 yj$ vagrant up
```

默认的账号还是`root`，密码还是`vagrant`  

这里设置了静态的`ip`,我们就可以通过静态`ip`直接访问虚拟机了  

```go
$ ssh root@192.168.44.113
``` 

### vagrant中的网络

#### 私有网络

`private_network`  

私有网络，对应于`virtualbox`的`host-only`网络模型，这种模型下，虚拟机之间和宿主机(的虚拟网卡)之间可以互相通信，但不在该网络内的设备无法访问虚拟机  

如果私有网络的虚机不在一个网络，`vagrant`为这些`private_network`网络配置的IP地址并不在同一个网段。`vagrant`会自动为不同网段创建对应的`host-only`网络。  

所以使用`private_network`如果没有外部机器(虚拟机宿主机之外的机器)连接，使用这种方式设置的静态`ip`，能够摆脱主机网络变换的限制。

PS:比如`public_network`如果更换了`wefi`连接，之前设置的静态`ip`可能就不可用了，因为网段不一样了。   

```go
vb1.vm.network "private_network", ip: "192.168.56.112"
```

#### 公有网络

`public_network`  

公有网络，对应于`virtualbox`的桥接模式，这种模式下，虚拟机的网络和宿主机的物理网卡是平等的，它们在同一个网络内，虚拟机可以访问外网，外界网络(特指能访问物理网卡的设备)也能访问虚拟机  

`vagrant`为`virtualbox`配置的`public_network`，其本质是将虚拟机加入到了`virtualbox`的桥接网络内。  

`vagrant`在将虚拟机的网卡加入桥接网络时，默认会交互式地询问用户要和哪个宿主机上的网卡进行桥接，一般来说，应该选择可以上外网的物理设备进行桥接。  

由于需要非交互式选择或者需要先指定要桥接的设备名，而且不同用户的网络环境不一样，因此如非必要，一般不在`vagrant`中为虚拟机配置`public_network`。  

公有网络的`iP`网络要和主机的网段一致。  

<img src="/img/vagrant_2.jpg" alt="vagrant" align=center />

```go
vb.vm.network "public_network", ip: "192.168.44.111",bridge: "en0: Wi-Fi (AirPort)"
```

### 常用的命令

|  子命令        | 功能说明  |
|  ----         | ----   |
| box	        |管理box镜像(box是创建虚拟机的模板)|
| init	        |初始化项目目录，将在当前目录下生成Vagrantfile文件|
| up	        |启动虚拟机，第一次执行将创建并初始化并启动虚拟机|
| reload	    |重启虚拟机|
| halt	        |将虚拟机关机|
| destroy	    |删除虚拟机(包括虚拟机文件)|
| suspend	    |暂停(休眠、挂起)虚拟机|
| resume	    |恢复已暂停(休眠、挂起)的虚拟机|
| snapshot	    |管理虚拟机快照(hyperv中叫检查点)|
| status	    |列出当前目录(Vagrantfile所在目录)下安装的虚拟机列表及它们的状态|
| global-status	|列出全局已安装虚拟机列表及它们的状态|
| ssh	        |通过ssh连接虚拟机|
| ssh-config	|输出ssh连接虚拟机时使用的配置项|
| port	        |查看各虚拟机映射的端口列表(hyperv不支持该功能)|


### 参考

【熟练使用vagrant(11)：vagrant配置虚拟机网络】https://www.junmajinlong.com/virtual/vagrant/vagrant_network/    

