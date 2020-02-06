## centos7下yum安装MariaDB

CentOS 7下mysql下替换成MariaDB了（MariaDB是从MySQL fork过来的，和MySQL有很好的兼容性 ）

### 使用yum快速安装  
````
yum install mariadb mariadb-server
systemctl start mariadb   #启动mariadb
systemctl enable mariadb  #设置开机自启动
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