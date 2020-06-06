<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [服务器重装ip未更改连不上](#%E6%9C%8D%E5%8A%A1%E5%99%A8%E9%87%8D%E8%A3%85ip%E6%9C%AA%E6%9B%B4%E6%94%B9%E8%BF%9E%E4%B8%8D%E4%B8%8A)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [原因](#%E5%8E%9F%E5%9B%A0)
  - [解决方法](#%E8%A7%A3%E5%86%B3%E6%96%B9%E6%B3%95)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 服务器重装ip未更改连不上

### 前言

重装了虚拟机，ip还保留了，但是发现连不上了

````
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@    WARNING: REMOTE HOST IDENTIFICATION HAS CHANGED!     @
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
IT IS POSSIBLE THAT SOMEONE IS DOING SOMETHING NASTY!
Someone could be eavesdropping on you right now (man-in-the-middle attack)!
It is also possible that a host key has just been changed.
The fingerprint for the ECDSA key sent by the remote host is
SHA256:+fniGCBsOCVoDRcIJrBc5x46D13paqRLuQwNDF5zECY.
Please contact your system administrator.
Add correct host key in /root/.ssh/known_hosts to get rid of this message.
Offending ECDSA key in /root/.ssh/known_hosts:1
ECDSA host key for 192.168.56.203 has changed and you have requested strict checking.
Host key verification failed.
fatal: Could not read from remote repository.

Please make sure you have the correct access rights
and the repository exists.
````

### 原因

当远程系统上的主机密钥发生更改时，或者因为它们是手动重新生成的，或者因为重新安装了ssh，新的主机密钥将与用户已知主机文件中存储的密钥不匹配，ssh将报告错误，然后退出。

### 解决方法

如果确定是创装服务器造成的，把当前ip的密匙删除就好了

````
# ssh-keygen -R 192.168.56.203
````


