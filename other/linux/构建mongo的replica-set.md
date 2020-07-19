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

启动

````
# ./bin/mongod -f  ./data/mongodb.conf
````
