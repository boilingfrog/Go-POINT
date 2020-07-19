## 搭建mongo的replica set

### 前言

准备三台机器，相互可以访问的。

### 安装

mongo1  

安装mongo

````
// 下载
# wget http://downloads.mongodb.org/linux/mongodb-linux-x86_64-rhel70-4.2.1.tgz
# tar zxvf mongodb-linux-x86_64-rhel70-4.2.1.tgz
# mv mongodb-linux-x86_64-rhel70-4.2.1 /usr/local/mongodb

// 创建db的data目录
# cd /usr/local/mongodb/
# mkdir data
# mkdir data/db
# mkdir data/logs
# touch data/mongodb.log
# cd data
# mv mongdb.log logs
````

写入配置文件，`vi mongodb.cof`

````
port = 27018
bind_ip = 0.0.0.0
#集群名字
replSet = mongos
#数据目录(自己刚才设置的位置)
dbpath = /usr/local/mongodb/data/db
#日志目录(自己刚才设置的位置))
logpath = /usr/local/mongodb/data/logs/mongodb.log
#设置后台运行
fork = true
#日志输出方式
logappend = true
#开启认证
#auth = true
#安全文件地址
#keyFile = ./keyFile
````

关于ip

- bind_ip: 192.168.0.136     #如果修改成本机Ip，那除了本机外的机器都可以连接（就是自己连不了、哈哈、蛋疼）  
- bind_ip: 0.0.0.0           #改成0，那么大家都可以访问（共赢）  
- bind_ip: 127.0.0.1         #改成127，那就只能自己练了（独吞）  

所以为了方便其他服务器和自己连接，就把`bind_ip`改成`0.0.0.0`  

启动 

````
# ./bin/mongod -f  ./data/mongodb.conf
````

三台机器du