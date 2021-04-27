## vagrant搭建centos

### 什么是vagrant

> Vagrant 是一个简单易用的部署工具，用英文说应该是 Orchestration Tool 。它能帮助开发人员迅速的构建一个开发环境，帮助测试人员构建测试环境，Vagrant 基于 Ruby 开发，使用开源 VirtualBox 作为虚拟化支持，可以轻松的跨平台部署。

### 如何使用

1、构建本地的目录

```go
 /Users/yj/vagrant/centos7
```

2、去官方下载对应的镜像文件，[官方下载地址](http://www.vagrantbox.es/)  

```go
MacBook-Pro-3:centos7 yj$ wget https://github.com/CommanderK5/packer-centos-template/releases/download/0.7.2/vagrant-centos-7.2.box
```

3、导入刚刚下载的镜像(box文件)

```go
MacBook-Pro-3:centos7 yj$ vagrant box add centos7.2 /Users/yj/vagrant/centos7/vagrant-centos-7.2.box 
==> vagrant: A new version of Vagrant is available: 2.2.15 (installed version: 2.2.14)!
==> vagrant: To upgrade visit: https://www.vagrantup.com/downloads.html

==> box: Box file was not detected as metadata. Adding it directly...
==> box: Adding box 'centos7.2' (v0) for provider: 
    box: Unpacking necessary files from: file:///Users/yj/vagrant/centos7/vagrant-centos-7.2.box
==> box: Successfully added box 'centos7.2' (v0) for 'virtualbox'!
```

4、初始化

```go
MacBook-Pro-3:centos7 yj$ vagrant init
```

这时候当前目录会生成一个`Vagrantfile`文件

5、修改Vagrantfile中的box名称

```go
config.vm.box = "centos7-1"
```

6、启动

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

7、登入

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

<img src="img/vagrant_1.jpg" alt="vagrant" align=center />

账号:root  

密码:vagrant