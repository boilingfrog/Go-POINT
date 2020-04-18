<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [centos7下安装pgsql10.3](#centos7%E4%B8%8B%E5%AE%89%E8%A3%85pgsql103)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [下载pgsql-10.3](#%E4%B8%8B%E8%BD%BDpgsql-103)
  - [安装](#%E5%AE%89%E8%A3%85)
    - [解压](#%E8%A7%A3%E5%8E%8B)
    - [安装基本的工具](#%E5%AE%89%E8%A3%85%E5%9F%BA%E6%9C%AC%E7%9A%84%E5%B7%A5%E5%85%B7)
    - [编译](#%E7%BC%96%E8%AF%91)
    - [安装](#%E5%AE%89%E8%A3%85-1)
    - [创建目录 data、log](#%E5%88%9B%E5%BB%BA%E7%9B%AE%E5%BD%95-datalog)
    - [加入系统环境变量](#%E5%8A%A0%E5%85%A5%E7%B3%BB%E7%BB%9F%E7%8E%AF%E5%A2%83%E5%8F%98%E9%87%8F)
    - [增加用户 postgres 并赋权](#%E5%A2%9E%E5%8A%A0%E7%94%A8%E6%88%B7-postgres-%E5%B9%B6%E8%B5%8B%E6%9D%83)
    - [初始化数据库](#%E5%88%9D%E5%A7%8B%E5%8C%96%E6%95%B0%E6%8D%AE%E5%BA%93)
    - [编辑配置文件](#%E7%BC%96%E8%BE%91%E9%85%8D%E7%BD%AE%E6%96%87%E4%BB%B6)
    - [启动服务](#%E5%90%AF%E5%8A%A8%E6%9C%8D%E5%8A%A1)
    - [查看版本](#%E6%9F%A5%E7%9C%8B%E7%89%88%E6%9C%AC)
    - [设置开机启动](#%E8%AE%BE%E7%BD%AE%E5%BC%80%E6%9C%BA%E5%90%AF%E5%8A%A8)
    - [关掉防火墙](#%E5%85%B3%E6%8E%89%E9%98%B2%E7%81%AB%E5%A2%99)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

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

#### 安装

````
make && make install
````
然后等待安装........
![](https://img2020.cnblogs.com/blog/1237626/202004/1237626-20200418002533956-1995789683.png)
直到出现`PostgreSQL installation complete.`表示安装成功了
![](https://img2020.cnblogs.com/blog/1237626/202004/1237626-20200418003122587-1027120284.png)

#### 创建目录 data、log

````
# mkdir /usr/local/pgsql/data
# mkdir /usr/local/pgsql/log
````

#### 加入系统环境变量

````
 vim /etc/profile
````

在最后写入

````
PGHOME=/usr/local/pgsql
export PGHOME
PGDATA=/usr/local/pgsql/data
export PGDATA
PATH=$PATH:$HOME/.local/bin:$HOME/bin:$PGHOME/bin
````

注意：`/usr/local/pgsql`需要修改为自己的安装目录  

使配置文件生效

````
# source /etc/profile
````

#### 增加用户 postgres 并赋权

````
# adduser postgres
# chown -R postgres:root /usr/local/pgsql/
````

修改密码（在root）

````
# passwd postgres 
更改用户 postgres 的密码 。
新的 密码：
无效的密码： 密码未通过字典检查 - 过于简单化/系统化
重新输入新的 密码：
passwd：所有的身份验证令牌已经成功更新。
````

#### 初始化数据库

注意：需要在`postgres`用户下初始化  

切换用户 postgres  

````
# su postgres
````

然后初始化数据库

````
# /usr/local/pgsql/bin/initdb -D /usr/local/pgsql/data/
````
![](https://img2020.cnblogs.com/blog/1237626/202004/1237626-20200418004700573-758488394.png)

#### 编辑配置文件

````
# vim /usr/local/ppgsql/data/postgresql.conf
````

修改

````
listen_addresses = '*'
port = 5432
````

![](https://img2020.cnblogs.com/blog/1237626/202004/1237626-20200418005154166-390939545.png)

同样修改

````i
# vim /usr/local/pgsql/data/pg_hba.conf
````
在最后面添加
![](https://img2020.cnblogs.com/blog/1237626/202004/1237626-20200418005728518-1098452113.png)

说明：

TYPE：pg的连接方式，local：本地unix套接字，host：tcp/ip连接  
DATABASE：指定数据库  
USER：指定数据库用户  
ADDRESS：ip地址，可以定义某台主机或某个网段，32代表检查整个ip地址，相当于固定的ip，24代表只检查前三位，最后一 位是0~255之间的任何一个  
METHOD：认证方式，常用的有ident，md5，password，trust，reject。  
md5是常用的密码认证方式。  
password是以明文密码传送给数据库，建议不要在生产环境中使用。  
trust是只要知道数据库用户名就能登录，建议不要在生产环境中使用。  
reject是拒绝认证。  

#### 启动服务

````
$ pg_ctl start -l /usr/local/pgsql/log/pg_server.log
could not change directory to "/root/postgresql-10.3": 权限不够
waiting for server to start.... done
server started
````

启动，停止，重启
````
./pg_ctl start\stop\restart -D /usr/local/pgsql/data/
````

#### 查看版本

````
# psql -V
psql (PostgreSQL) 10.3
````

#### 设置开机启动

将pgsql安装包中的linux文件复制到/etc/init.d或者/etc/rc.d

````
[root@10 postgresql-10.3]# cp contrib/start-scripts/linux /etc/init.d/pgsql
````

根据安装路径修改pgsql文件中的配置项  

````
## EDIT FROM HERE

# Installation prefix （安装路径）
prefix=/usr/local/pgsql

# Data directory （data路径）
PGDATA="/usr/local/pgsql/data"
````

修改pgsql文件权限  

````
# chmod +x /etc/init.d/pgsql
````

开机执行pgsql文件

````
# chkconfig --add pgsql
````

#### 关掉防火墙

````
$ systemctl stop firewalld
$ systemctl disable firewalld
````