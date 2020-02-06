## centos7下yum安装MariaDB

CentOS 7下mysql下替换成MariaDB了。  
mariadb简介MariaDB数据库管理系统是MySQL的一个分支，主要由开源社区在维护，采用GPL授权
许可 MariaDB的目的是完全兼容MySQL，包括API和命令行，使之能轻松成为MySQL的代替品。

### 使用yum快速安装  

#### 安装
````
]# yum install -y mariadb-server
已加载插件：fastestmirror
Loading mirror speeds from cached hostfile
 * base: mirrors.aliyun.com
 * extras: mirrors.aliyun.com
 * updates: mirror.bit.edu.cn
正在解决依赖关系
--> 正在检查事务
---> 软件包 mariadb-server.x86_64.1.5.5.64-1.el7 将被 安装
--> 正在处理依赖关系 mariadb-libs(x86-64) = 1:5.5.64-1.el7，它被软件包 1:mariadb-server-5.5.64-1.el7.x86_64 需要
--> 正在处理依赖关系 mariadb(x86-64) = 1:5.5.64-1.el7，它被软件包 1:mariadb-server-5.5.64-1.el7.x86_64 需要
--> 正在处理依赖关系 perl-DBI，它被软件包 1:mariadb-server-5.5.64-1.el7.x86_64 需要
--> 正在处理依赖关系 perl-DBD-MySQL，它被软件包 1:mariadb-server-5.5.64-1.el7.x86_64 需要
--> 正在处理依赖关系 perl(Data::Dumper)，它被软件包 1:mariadb-server-5.5.64-1.el7.x86_64 需要
--> 正在处理依赖关系 perl(DBI)，它被软件包 1:mariadb-server-5.5.64-1.el7.x86_64 需要
--> 正在检查事务
---> 软件包 mariadb.x86_64.1.5.5.64-1.el7 将被 安装
---> 软件包 mariadb-libs.x86_64.1.5.5.64-1.el7 将被 安装
---> 软件包 perl-DBD-MySQL.x86_64.0.4.023-6.el7 将被 安装
---> 软件包 perl-DBI.x86_64.0.1.627-4.el7 将被 安装
--> 正在处理依赖关系 perl(RPC::PlServer) >= 0.2001，它被软件包 perl-DBI-1.627-4.el7.x86_64 需要
--> 正在处理依赖关系 perl(RPC::PlClient) >= 0.2000，它被软件包 perl-DBI-1.627-4.el7.x86_64 需要
---> 软件包 perl-Data-Dumper.x86_64.0.2.145-3.el7 将被 安装
--> 正在检查事务
---> 软件包 perl-PlRPC.noarch.0.0.2020-14.el7 将被 安装
--> 正在处理依赖关系 perl(Net::Daemon) >= 0.13，它被软件包 perl-PlRPC-0.2020-14.el7.noarch 需要
--> 正在处理依赖关系 perl(Net::Daemon::Test)，它被软件包 perl-PlRPC-0.2020-14.el7.noarch 需要
--> 正在处理依赖关系 perl(Net::Daemon::Log)，它被软件包 perl-PlRPC-0.2020-14.el7.noarch 需要
--> 正在处理依赖关系 perl(Compress::Zlib)，它被软件包 perl-PlRPC-0.2020-14.el7.noarch 需要
--> 正在检查事务
---> 软件包 perl-IO-Compress.noarch.0.2.061-2.el7 将被 安装
--> 正在处理依赖关系 perl(Compress::Raw::Zlib) >= 2.061，它被软件包 perl-IO-Compress-2.061-2.el7.noarch 需要
--> 正在处理依赖关系 perl(Compress::Raw::Bzip2) >= 2.061，它被软件包 perl-IO-Compress-2.061-2.el7.noarch 需要
---> 软件包 perl-Net-Daemon.noarch.0.0.48-5.el7 将被 安装
--> 正在检查事务
---> 软件包 perl-Compress-Raw-Bzip2.x86_64.0.2.061-3.el7 将被 安装
---> 软件包 perl-Compress-Raw-Zlib.x86_64.1.2.061-4.el7 将被 安装
--> 解决依赖关系完成

依赖关系解决

=============================================================================================================
 Package                             架构               版本                          源                大小
=============================================================================================================
正在安装:
 mariadb-server                      x86_64             1:5.5.64-1.el7                base              11 M
为依赖而安装:
 mariadb                             x86_64             1:5.5.64-1.el7                base             8.7 M
 mariadb-libs                        x86_64             1:5.5.64-1.el7                base             759 k
 perl-Compress-Raw-Bzip2             x86_64             2.061-3.el7                   base              32 k
 perl-Compress-Raw-Zlib              x86_64             1:2.061-4.el7                 base              57 k
 perl-DBD-MySQL                      x86_64             4.023-6.el7                   base             140 k
 perl-DBI                            x86_64             1.627-4.el7                   base             802 k
 perl-Data-Dumper                    x86_64             2.145-3.el7                   base              47 k
 perl-IO-Compress                    noarch             2.061-2.el7                   base             260 k
 perl-Net-Daemon                     noarch             0.48-5.el7                    base              51 k
 perl-PlRPC                          noarch             0.2020-14.el7                 base              36 k

事务概要
=============================================================================================================
安装  1 软件包 (+10 依赖软件包)

总下载量：22 M
安装大小：115 M
Downloading packages:
(1/11): perl-Compress-Raw-Bzip2-2.061-3.el7.x86_64.rpm                                |  32 kB  00:00:00     
(2/11): perl-Compress-Raw-Zlib-2.061-4.el7.x86_64.rpm                                 |  57 kB  00:00:00     
(3/11): perl-DBD-MySQL-4.023-6.el7.x86_64.rpm                                         | 140 kB  00:00:00     
(4/11): mariadb-libs-5.5.64-1.el7.x86_64.rpm                                          | 759 kB  00:00:01     
(5/11): perl-Data-Dumper-2.145-3.el7.x86_64.rpm                                       |  47 kB  00:00:00     
(6/11): perl-IO-Compress-2.061-2.el7.noarch.rpm                                       | 260 kB  00:00:00     
(7/11): perl-Net-Daemon-0.48-5.el7.noarch.rpm                                         |  51 kB  00:00:00     
(8/11): perl-PlRPC-0.2020-14.el7.noarch.rpm                                           |  36 kB  00:00:00     
(9/11): mariadb-5.5.64-1.el7.x86_64.rpm                                               | 8.7 MB  00:00:02     
(10/11): perl-DBI-1.627-4.el7.x86_64.rpm                                              | 802 kB  00:00:01     
(11/11): mariadb-server-5.5.64-1.el7.x86_64.rpm                                       |  11 MB  00:00:15     
-------------------------------------------------------------------------------------------------------------
总计                                                                         1.5 MB/s |  22 MB  00:00:15     
Running transaction check
Running transaction test
Transaction test succeeded
Running transaction
  正在安装    : 1:mariadb-libs-5.5.64-1.el7.x86_64                                                      1/11 
  正在安装    : perl-Data-Dumper-2.145-3.el7.x86_64                                                     2/11 
  正在安装    : 1:mariadb-5.5.64-1.el7.x86_64                                                           3/11 
  正在安装    : perl-Compress-Raw-Bzip2-2.061-3.el7.x86_64                                              4/11 
  正在安装    : 1:perl-Compress-Raw-Zlib-2.061-4.el7.x86_64                                             5/11 
  正在安装    : perl-IO-Compress-2.061-2.el7.noarch                                                     6/11 
  正在安装    : perl-Net-Daemon-0.48-5.el7.noarch                                                       7/11 
  正在安装    : perl-PlRPC-0.2020-14.el7.noarch                                                         8/11 
  正在安装    : perl-DBI-1.627-4.el7.x86_64                                                             9/11 
  正在安装    : perl-DBD-MySQL-4.023-6.el7.x86_64                                                      10/11 
  正在安装    : 1:mariadb-server-5.5.64-1.el7.x86_64                                                   11/11 
  验证中      : 1:mariadb-libs-5.5.64-1.el7.x86_64                                                      1/11 
  验证中      : perl-Net-Daemon-0.48-5.el7.noarch                                                       2/11 
  验证中      : perl-Data-Dumper-2.145-3.el7.x86_64                                                     3/11 
  验证中      : 1:mariadb-5.5.64-1.el7.x86_64                                                           4/11 
  验证中      : perl-DBD-MySQL-4.023-6.el7.x86_64                                                       5/11 
  验证中      : perl-IO-Compress-2.061-2.el7.noarch                                                     6/11 
  验证中      : 1:perl-Compress-Raw-Zlib-2.061-4.el7.x86_64                                             7/11 
  验证中      : 1:mariadb-server-5.5.64-1.el7.x86_64                                                    8/11 
  验证中      : perl-DBI-1.627-4.el7.x86_64                                                             9/11 
  验证中      : perl-Compress-Raw-Bzip2-2.061-3.el7.x86_64                                             10/11 
  验证中      : perl-PlRPC-0.2020-14.el7.noarch                                                        11/11 

已安装:
  mariadb-server.x86_64 1:5.5.64-1.el7                                                                       

作为依赖被安装:
  mariadb.x86_64 1:5.5.64-1.el7                         mariadb-libs.x86_64 1:5.5.64-1.el7                  
  perl-Compress-Raw-Bzip2.x86_64 0:2.061-3.el7          perl-Compress-Raw-Zlib.x86_64 1:2.061-4.el7         
  perl-DBD-MySQL.x86_64 0:4.023-6.el7                   perl-DBI.x86_64 0:1.627-4.el7                       
  perl-Data-Dumper.x86_64 0:2.145-3.el7                 perl-IO-Compress.noarch 0:2.061-2.el7               
  perl-Net-Daemon.noarch 0:0.48-5.el7                   perl-PlRPC.noarch 0:0.2020-14.el7                   

完毕！

````
#### mariadb相关命令
````
yum install mariadb mariadb-server
systemctl start mariadb   #启动mariadb
systemctl enable mariadb  #设置开机自启动
systemctl stop mariadb    #停止MariaDB
systemctl restart mariadb #重启MariaDB
mysql_secure_installation #设置root密码等相关
mysql -uroot -p           #测试登录   
````

当我启动mariadb的时候出现了如下的错误  
````
Failed to start mariadb.service: Unit not found.
````
出现这个的原因是机器上之前安装了mysql，删除掉就可以了。 



#### 1、查看mysql安装了哪些东西
````
# rpm -qa |grep -i mysql
mysql-community-common-5.7.29-1.el7.x86_64
mysql-community-client-5.7.29-1.el7.x86_64
mysql-community-libs-compat-5.7.29-1.el7.x86_64
mysql-community-libs-5.7.29-1.el7.x86_64
mysql-community-server-5.7.29-1.el7.x86_64
````

#### 2、开始卸载
````
yum remove mysql-community-common-5.7.29-1.el7.x86_64
yum remove mysql-community-client-5.7.29-1.el7.x86_64
yum mysql-community-libs-compat-5.7.29-1.el7.x86_64
yum remove mysql-community-libs-5.7.29-1.el7.x86_64
yum remove mysql-community-server-5.7.29-1.el7.x86_64
````
#### 3、查看是否卸载完成

````
# rpm -qa |grep -i mysql
   
````

#### 4、查找mysql相关目录
````
# find / -name mysql
/usr/share/mysql
````

#### 5、删除相关目录
````
# rm -rf /usr/share/mysql
````

#### 6、删除/etc/my.cnf
````
# rm -rf /etc/my.cnf
````

#### 7、删除/var/log/mysqld.log（如果不删除这个文件，会导致新安装的mysql无法生存新密码，导致无法登陆）
````
# rm -rf /var/log/mysqld.log
````



### 参考
【Centos7 完全卸载mysql】https://www.jianshu.com/p/ef58fb333cd6 