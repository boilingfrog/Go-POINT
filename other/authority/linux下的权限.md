<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [linux中的权限](#linux%E4%B8%AD%E7%9A%84%E6%9D%83%E9%99%90)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [数字权限](#%E6%95%B0%E5%AD%97%E6%9D%83%E9%99%90)
    - [三位数字权限](#%E4%B8%89%E4%BD%8D%E6%95%B0%E5%AD%97%E6%9D%83%E9%99%90)
      - [读(r)](#%E8%AF%BBr)
      - [写(w)](#%E5%86%99w)
      - [执行(x)](#%E6%89%A7%E8%A1%8Cx)
      - [无权限(-)](#%E6%97%A0%E6%9D%83%E9%99%90-)
      - [三位数字权限的转换](#%E4%B8%89%E4%BD%8D%E6%95%B0%E5%AD%97%E6%9D%83%E9%99%90%E7%9A%84%E8%BD%AC%E6%8D%A2)
      - [如何设置权限](#%E5%A6%82%E4%BD%95%E8%AE%BE%E7%BD%AE%E6%9D%83%E9%99%90)
      - [最高位的含义](#%E6%9C%80%E9%AB%98%E4%BD%8D%E7%9A%84%E5%90%AB%E4%B9%89)
    - [四位数字权限](#%E5%9B%9B%E4%BD%8D%E6%95%B0%E5%AD%97%E6%9D%83%E9%99%90)
      - [SUID](#suid)
      - [SGID](#sgid)
      - [SBIT](#sbit)
      - [四位数字权限的转换](#%E5%9B%9B%E4%BD%8D%E6%95%B0%E5%AD%97%E6%9D%83%E9%99%90%E7%9A%84%E8%BD%AC%E6%8D%A2)
      - [如何设置权限](#%E5%A6%82%E4%BD%95%E8%AE%BE%E7%BD%AE%E6%9D%83%E9%99%90-1)
    - [如何改变文件属性](#%E5%A6%82%E4%BD%95%E6%94%B9%E5%8F%98%E6%96%87%E4%BB%B6%E5%B1%9E%E6%80%A7)
      - [改变所属群组, chgrp](#%E6%94%B9%E5%8F%98%E6%89%80%E5%B1%9E%E7%BE%A4%E7%BB%84-chgrp)
      - [改变文件拥有者, chown](#%E6%94%B9%E5%8F%98%E6%96%87%E4%BB%B6%E6%8B%A5%E6%9C%89%E8%80%85-chown)
      - [改变权限, chmod](#%E6%94%B9%E5%8F%98%E6%9D%83%E9%99%90-chmod)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## linux中的权限

### 前言

最近使用`ansible`搭建了mongo的`replica set`。发现很多次的失败都是文件权限不对造成的。那么就总结下吧。  

### 数字权限

数字权限，我们常见的是三位的数字表示的权限，当然四位的我们也难免会遇到，接下我们看下他们直接的区别。  

#### 三位数字权限

linux中的文件权限分为读、写、执行.对应的字母就是r、w、x.  

##### 读(r)
   
Read对文件而言，具有读取文件内容的权限；对目录来说，具有浏览该目录信息的权限。  

##### 写(w)
 
Write对文件而言，具有修改文件内容的权限；对于目录来说，具有删除移动目录内文件的权限。  

##### 执行(x)  

Execute对文件而言，具有执行文件的权限；对目录来说，具有进入目录的权限。  

##### 无权限(-)  

`-`表示不具有该项选项  

linux下的权限粒度有 `拥有者`、`群组`、`其他组`三种。每个文件都可以针对这三个粒度，设置不同的(r,w,x)读写执行权限。通常情况下，一个文件只能归属于一个用户和组， 如果其它的用户想有这个文件的权限，则可以将该用户加入具备权限的群组，一个用户可以同时归属于多个组。  

Linux上通常使用chmod命令对文件的权限进行设置和更改。  

我们规定数字4、2、1表示读、写、执行的权限。即r=4，w=2，x=1。  

栗如:  

````
rwx = 4 + 2 + 1 = 7
rw = 4 + 2 = 6
rx = 4 +1 = 5
````

即  

若要同时设置 rwx (可读写运行） 权限则将该权限位 设置 为 4 + 2 + 1 = 7  

若要同时设置 rw- （可读写不可运行）权限则将该权限位 设置 为 4 + 2 = 6  

若要同时设置 r-x （可读可运行不可写）权限则将该权限位 设置 为 4 +1 = 5  

##### 三位数字权限的转换

读，写，执行对应的权限二进制表示为

```
r-- = 100
-w- = 010
--x = 001
--- = 000
```

将二进制转换成十进制就是

```
r-- = 100 4
-w- = 010 2
--x = 001 1
--- = 000 0
```

故

```
rwx 权限就是 4 + 2 + 1 = 7  
rw- 权限就是 4 + 2 = 6 
...
```

##### 如何设置权限

每个文件都有`拥有者`、`群组`、`其他组`三种权限的粒度。我们在设置权限的时候也是针对这几个粒度进行配置的。即我们可以使用三个数字来配置这三个权限粒度的权限。  

> chmod <abc> file...

````
其中
a,b,c各为一个数字，分别代表User、Group、及Other的权限。
相当于简化版的
chmod u=权限,g=权限,o=权限 file...
而此处的权限将用8进制的数字来表示User、Group、及Other的读、写、执行权限
````

栗子：  

设置所有人可以读写及执行

```asciidoc
chmod 777 file  (等价于  chmod u=rwx,g=rwx,o=rwx file 或  chmod a=rwx file)
```

设置拥有者可读写，其他人不可读写执行

```
chmod 600 file (等价于  chmod u=rw,g=---,o=--- file 或 chmod u=rw,go-rwx file )
```

一些常见的权限

````
-rw------- (600)      只有拥有者有读写权限。
-rw-r--r-- (644)      只有拥有者有读写权限；而属组用户和其他用户只有读权限。
-rwx------ (700)     只有拥有者有读、写、执行权限。
-rwxr-xr-x (755)    拥有者有读、写、执行权限；而属组用户和其他用户只有读、执行权限。
-rwx--x--x (711)    拥有者有读、写、执行权限；而属组用户和其他用户只有执行权限。
-rw-rw-rw- (666)   所有用户都有文件读、写权限。
-rwxrwxrwx (777)  所有用户都有读、写、执行权限
````

##### 最高位的含义

关于第一位最高位的解释： 上面我们说到了权限表示中后九位的含义，剩下的第一位代表的是文件的类型，类型可以是下面几个中的一个：  

```
d代表的是目录(directroy)
-代表的是文件(regular file)
s代表的是套字文件(socket)
p代表的管道文件(pipe)或命名管道文件(named pipe)
l代表的是符号链接文件(symbolic link)
b代表的是该文件是面向块的设备文件(block-oriented device file)
```

#### 四位数字权限

linux除了设置正常的读写操作权限外，还有关于一类设置也是涉及到权限，叫做Linxu附加权限。包括 SET位权限（suid，sgid）和粘滞位权限（sticky）。  

##### SUID 

让本来没有相应权限的用户运行这个程序时，可以访问他没有权限访问的资源。  

```
1、SUID权限仅对二进制程序有效。
2、执行者对于该程序需要具有x的可执行权限。
3、本权限仅在执行该程序的过程中有效。
4、执行者将具有该程序拥有者的权限。
```

栗子：  

```
$ ls -la /bin/su
-rwsr-xr-x 1 root root 63568 7月  16 06:52 /bin/su
```

对于su这个命令，拥有者是root,但是实际上是无论任何人，都可以执行，并且拥有的执行权限和root一样。因为这个设置了SUID，并且设置了拥有者有可执行权限。从上面可以看到，不管是文件拥有者，文件拥有者所属组，还是其他人，都是具有x权限的，所以都可以执行该程序，执行之后就将具有该程序拥有者的权限，即root的权限，这也就是su命令能够切换用户权限的实现原理。    

##### SGID
 
和SUID类似，不过SGID是针对所属用户组权限的。  

```
1、SGID对二进制程序有用；
2、程序执行者对于该程序来说需具备x的权限；
3、SGID主要用在目录上；
```

SET位权限表示形式：  

如果一个文件被设置了suid或sgid位，会分别表现在所有者或同组用户的权限的可执行位上；如果文件设置了suid还设置了x（执行）位，则相应的执行位表示为s(小写)。但是，如果没有设置x位，它将表示为S(大写)。如：  

````
1、-rwsr-xr-x 表示设置了suid，且拥有者有可执行权限
2、-rwSr--r-- 表示suid被设置，但拥有者没有可执行权限
3、-rwxr-sr-x 表示sgid被设置，且群组用户有可执行权限
4、-rw-r-Sr-- 表示sgid被设置，但群组用户没有可执行权限
````

设置方式:  

SET位权限可以可以使用chmod命令设置，栗子：  

```
chmod u+s filename 	设置suid位
chmod u-s filename 	去掉suid设置
chmod g+s filename 	设置sgid位
chmod g-s filename 	去掉sgid设置
``` 

##### SBIT 

目前只针对目录有效，对于目录的作用是：当用户在该目录下建立文件或目录时，仅有自己与 root才有权力删除。  

栗子：  

最具有代表的就是/tmp目录，任何人都可以在/tmp内增加、修改文件（因为权限全是rwx），但仅有该文件/目录建立者与 root能够删除自己的目录或文件  

**注：SBIT对文件不起作用。**  

粘滞位权限表示形式:  

一个文件或目录被设置了粘滞位权限，会表现在其他组用户的权限的可执行位上。如果文件设置了`sticky`还设置了x（执行）位，其他组用户的权限的可执行位为t(小写)。但是，如果没有设置x位，它将表示为T(大写)。如：  

```asciidoc
1、-rwsr-xr-t 表示设置了粘滞位且其他用户组有可执行权限
2、-rwSr--r-T 表示设置了粘滞位但其他用户组没有可执行权限
```

##### 四位数字权限的转换

同样使用二进制表示下:  

SGT分别表示SUID权限、SGID权限、和 粘滞位权限  

```
S-- = 100
-G- = 010
--T = 001
--- = 000
```

和前面的rwx的转换一样，转换成是十进制就是 

```
S-- = 100 4
-G- = 010 2
--T = 001 1
--- = 000 0
```

- suid   位代表数字是 4
- sgid   位代表数字是 2
- sticky 位代表数字是 1

##### 如何设置权限

使用chmod进行权限的设置  

```
chmod <abcd> file
```

栗子：  

使用一个栗子对比下吧   

设置 netlogin 的权限为拥有者可读写执行，群组和其他权限为可读可执行 

```
chmod 755 netlogin
``` 
 
设置 netlogin 的权限为拥有者可读写执行，群组和其他权限为可读可执行，并且设置SUID

```
chmod 4755 netlogin
```

chmod 4755与chmod 755对比多了附加权限值4，这个4表示其他用户执行文件时，具有与所有者同样的权限（设置了SUID）。


**为什么要设置4755 而不是 755？**

假设netlogin是root用户创建的一个上网认证程序，如果其他用户要上网也要用到这个程序，那就需要root用户运行chmod 755 netlogin命令使其他用户也能运行netlogin。但假如netlogin执行时需要访问一些只有root用户才有权访问的文件，那么其他用户执行netlogin时可能因为权限不够还是不能上网。这种情况下，就可以用 chmod 4755 netlogin 设置其他用户在执行netlogin也有root用户的权限，从而顺利上网。  

#### 如何改变文件属性

我们先介绍几个常用于群组、拥有者、各种身份的权限之修改的指令，如下所示：  

chgrp ：改变文件所属群组  
chown ：改变文件拥有者  
chmod ：改变文件的权限, SUID, SGID, SBIT等等的特性  

##### 改变所属群组, chgrp

改变一个文件的群组真是很简单的，直接以chgrp来改变即可。chgrp的拼写也是有点意思的，这个指令就是change group的缩写。  

不过需要注意的是，我们要改变的用户组，必须在在/etc/group文件内存在才行。  

假设你是以root的身份登入Linux系统的，那么在你的家目录内有一个install.log的文件， 如何将该文件的群组改变一下呢？假设你已经知道在/etc/group里面已经存在一个名为users的群组， 但是testing这个群组名字就不存在/etc/group当中了，此时改变群组成为users与testing分别会有什么现象发生呢？  

![chgrp](https://github.com/zhan-liz/Go-POINT/blob/master/img/au_1.png?raw=true)

发现了吗？文件的群组被改成users了，但是要改成testing的时候， 就会发生错误  

##### 改变文件拥有者, chown

chown的拼写也是change owner的拼写。同样用户必须是已经存在系统中的账号，也就是在/etc/passwd 这个文件中有纪录的用户名称才能改变。  

chown同时也可以修改所属群组  

![chgrp](https://github.com/zhan-liz/Go-POINT/blob/master/img/au_2.png?raw=true)

同时chown也可以简单的修改用户的组  

```
chown .sshd install.log
```

其中的sshd代表的就是用户的群组

##### 改变权限, chmod

文章开头已经描述的很清楚了  

### 参考
【第六章、Linux 的文件权限与目录配置】http://cn.linux.vbird.org/linux_basic/0210filepermission.php  
【Linux权限详解（chmod、600、644、666、700、711、755、777、4755、6755、7755）】https://blog.csdn.net/u013197629/article/details/73608613  
【深入Linux文件权限 SUID/SGID/SBIT】https://blog.csdn.net/imkelt/article/details/53054309?utm_medium=distribute.pc_relevant.none-task-blog-BlogCommendFromMachineLearnPai2-2.channel_param&depth_1-utm_source=distribute.pc_relevant.none-task-blog-BlogCommendFromMachineLearnPai2-2.channel_param  
【第六章、Linux 的文件权限与目录配置】http://cn.linux.vbird.org/linux_basic/0210filepermission.php#filepermission_ch