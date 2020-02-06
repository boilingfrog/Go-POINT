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