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
#### 修改root的密码
````
update mysql.user set password=PASSWORD('yhb123456') where user='root';
// 更新权限
flush privileges; 
````
#### 新建用户
````
// create user  '用户名'@'主机' identified by '密码'   如果只允许本机访问 @'localhost'  , 或者指定一个ip  @'192.xx.xx.xx' 或者使用通配: @'%'
create user 'read_visa'@'%' identified by '123456';
````
#### 给用户分配权限
````
// grant 操作类型 on 数据库.表 to 用户@'主机'   数据库,表,主机都支持通配符 grant select, insert on *.* to  'read_visa'@'%'
// grant all on visa.* to 'read_visa'@'%'; // all 表示所有权限
grant select on visa.* to 'read_visa'@'%';
````

## 当我启动mariadb的时候出现了如下的错误  
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
【centos7 mariadb安装 MySql】https://www.jianshu.com/p/f55a31ae0cea 