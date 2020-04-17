## centos7下安装pgsql10.3

### 前言

在centos7上面安装pgsql-10.3,在网上找了很多的文章，试了好久才成功．那就总结下,安装的过程吧，避免下次浪费时间．

### 下载pgsql-10.3

系统版本`centos7`  

下载pgsql-10.3:`https://www.postgresql.org/ftp/source/v10.3/`  

上传tar包到服务器

````
$ scp postgresql-10.3.tar.gz root@192.168.56.189:~
The authenticity of host '192.168.56.189 (192.168.56.189)' can't be established.
ECDSA key fingerprint is SHA256:q2lore2LaeBsH4j3jmEVg0VYbfudDDR4LkmF/rt+Zp0.
Are you sure you want to continue connecting (yes/no)? yes
Failed to add the host to the list of known hosts (/home/liz/.ssh/known_hosts).
root@192.168.56.189's password: 
postgresql-10.3.tar.gz                                                     100%   25MB  37.5MB/s   00:00  
````

### 安装

#### 解压
````
# tar -xzvf postgresql-10.3.tar.gz
````

#### 安装基本的工具

````
yum install -y vim lrzsz tree wget gcc gcc-c++ readline-devel zlib-devel
````

#### 编译

进入到刚刚解压的文件夹，执行命令

````
./configure --prefix=/usr/local/pgsql
````

后面的`/usr/local/pgsql`表示的是要编译安装的具体位置，可以自己定义
![](https://img2020.cnblogs.com/blog/1237626/202004/1237626-20200418002144034-1841400582.png)