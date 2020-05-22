## linux连接远程的服务器报错WARNING: REMOTE HOST IDENTIFICATION HAS CHANGED!

### 车祸现场
首先来看下车祸现场
````
$ ssh root@192.168.56.109
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@    WARNING: REMOTE HOST IDENTIFICATION HAS CHANGED!     @
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
IT IS POSSIBLE THAT SOMEONE IS DOING SOMETHING NASTY!
Someone could be eavesdropping on you right now (man-in-the-middle attack)!
It is also possible that a host key has just been changed.
The fingerprint for the ECDSA key sent by the remote host is
SHA256:L3GF07vcBbprWBZYgN5vebKs4crGIuNllDg0S0UPaMA.
Please contact your system administrator.
Add correct host key in /home/liz/.ssh/known_hosts to get rid of this message.
Offending ECDSA key in /home/liz/.ssh/known_hosts:18
  remove with:
  ssh-keygen -f "/home/liz/.ssh/known_hosts" -R 192.168.56.109
ECDSA host key for 192.168.56.109 has changed and you have requested strict checking.
Host key verification failed.

````
### 解决方法

出现这个的原因是电脑上保存，我们电脑保存的远程服务器的公钥信息是老的，可能这个远程服务器已经重装，或者更换了密匙。所以匹配不上了。
很简单只需要将当前远程ip对应的信息清除掉就行了
````
ssh-keygen -R XXXXXXXXXXX
````
下面是我电脑上的执行
````
$ ssh-keygen -R 192.168.56.109
# Host 192.168.56.109 found: line 18
/home/liz/.ssh/known_hosts updated.
Original contents retained as /home/liz/.ssh/known_hosts.old
$ ssh root@192.168.56.109
The authenticity of host '192.168.56.109 (192.168.56.109)' can't be established.
ECDSA key fingerprint is SHA256:L3GF07vcBbprWBZYgN5vebKs4crGIuNllDg0S0UPaMA.
Are you sure you want to continue connecting (yes/no)? yes
````
成功了