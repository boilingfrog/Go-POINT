## centos7下面操作mysql添加，授权，删除用户

### 添加用户
以root用户登录数据库，运行以下命令：
````
create user test identified by '123456789';
````
上面创建了用户test,密码是123456789。我们在mysql.user表里面可以看到新增的用户信息
````MariaDB [mysql]> select user,host,password from user where user='test';
    +------+----------------+-------------------------------------------+
    | user | host           | password                                  |
    +------+----------------+-------------------------------------------+
    | test | %              | *CC67043C7BCFF5EEA5566BD9B1F3C74FD9A5CF5D |
    +------+----------------+-------------------------------------------+
````
### 授权  

命令格式：grant privilegesCode on dbName.tableName to username@host identified by "password";
````
grant all privileges on test.* to 'test'@'%' identified by '123456789';
flush privileges;
````
上面的语句将test表的所有操作权限都给了用户test，并且密码是123456789
同样我们查看mysql.user的信息
````
select user,host,password, Select_priv,Insert_priv, Update_priv ,Delete_priv,Create_priv ,Drop_priv from user where user='test';
+------+----------------+-------------------------------------------+-------------+-------------+-------------+-------------+-------------+-----------+
| user | host           | password                                  | Select_priv | Insert_priv | Update_priv | Delete_priv | Create_priv | Drop_priv |
+------+----------------+-------------------------------------------+-------------+-------------+-------------+-------------+-------------+-----------+
| test | %              | *CC67043C7BCFF5EEA5566BD9B1F3C74FD9A5CF5D | Y           | Y           | Y           | Y           | Y           | Y         |
+------+----------------+-------------------------------------------+-------------+-------------+-------------+-------------+-------------+-----------+
1 rows in set (0.00 sec)
````
也可以使用show grants命令查看授权的权限信息
````
show grants for 'test';
+--------------------------------------------------------------------------------------------------------------+
| Grants for test@%                                                                                            |
+--------------------------------------------------------------------------------------------------------------+
| GRANT ALL PRIVILEGES ON *.* TO 'test'@'%' IDENTIFIED BY PASSWORD '*CC67043C7BCFF5EEA5566BD9B1F3C74FD9A5CF5D' |
| GRANT ALL PRIVILEGES ON `test`.* TO 'test'@'%'                                                               |
+--------------------------------------------------------------------------------------------------------------+
2 rows in set (0.00 sec)
````

#### privilegesCode表示授予的权限类型，常用的有以下几种类型[1]：
- all privileges：所有权限。
- select：读取权限。
- delete：删除权限。
- update：更新权限。
- create：创建权限。
- drop：删除数据库、数据表权限。

### dbName.tableName表示授予权限的具体库或表，常用的有以下几种选项：
- .：授予该数据库服务器所有数据库的权限。
- dbName.*：授予dbName数据库所有表的权限。
- dbName.dbTable：授予数据库dbName中dbTable表的权限。
- username@host表示授予的用户以及允许该用户登录的IP地址。其中Host有以下几种类型：
- localhost：只允许该用户在本地登录，不能远程登录。
- %：允许在除本机之外的任何一台机器远程登录。
- 192.168.52.32：具体的IP表示只允许该用户从特定IP登录。

### password指定该用户登录时的面。

### flush privileges表示刷新权限变更。

### 修改密码
````
update mysql.user set password = password('123') where user = 'test' and host = '%';
flush privileges;
````
### 删除用户
````
drop user test@'%';
````
drop user命令会删除用户以及对应的权限，执行命令后你会发现mysql.user表和mysql.db表的相应记录都消失了。

