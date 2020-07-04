<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [通过docker-compose搭建mongo的replica set高可用](#%E9%80%9A%E8%BF%87docker-compose%E6%90%AD%E5%BB%BAmongo%E7%9A%84replica-set%E9%AB%98%E5%8F%AF%E7%94%A8)
  - [备份数据](#%E5%A4%87%E4%BB%BD%E6%95%B0%E6%8D%AE)
    - [使用数据库连接工具](#%E4%BD%BF%E7%94%A8%E6%95%B0%E6%8D%AE%E5%BA%93%E8%BF%9E%E6%8E%A5%E5%B7%A5%E5%85%B7)
    - [备份数据到本地](#%E5%A4%87%E4%BB%BD%E6%95%B0%E6%8D%AE%E5%88%B0%E6%9C%AC%E5%9C%B0)
    - [数据恢复](#%E6%95%B0%E6%8D%AE%E6%81%A2%E5%A4%8D)
  - [集群搭建](#%E9%9B%86%E7%BE%A4%E6%90%AD%E5%BB%BA)
    - [生成keyFile](#%E7%94%9F%E6%88%90keyfile)
    - [创建yml文件](#%E5%88%9B%E5%BB%BAyml%E6%96%87%E4%BB%B6)
    - [初始化副本集](#%E5%88%9D%E5%A7%8B%E5%8C%96%E5%89%AF%E6%9C%AC%E9%9B%86)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 通过docker-compose搭建mongo的replica set高可用

搭建一个mongo的集群，同时原来单机mongo的数据需要迁移到集群中。  

处理思路：单机mongo的数据通过`mongodump`备份，然后集群搭建起来了，在`mongorestore`导入到集群中，实现数据的迁移。  


### 备份数据

#### 使用数据库连接工具

备份(mongodump)与恢复(mongorestore)  

#### 备份数据到本地

``
mongodump -h 192.168.56.201:37018 -u handle -p jimeng2017 -o /home/liz/Desktop/mongo-bei
``

#### 数据恢复

新的集群安装完成之后，恢复数据到Primary节点，集群会自动同步到副本集中

````
mongorestore -h 192.168.56.201:37017 -u handle -p jimeng2017  /home/liz/Desktop/mongo-bei
````

注意：更换自己服务器上面的ip和mongo对应的账号密码

### 集群搭建

#### 生成keyFile

````
# 400权限是要保证安全性，否则mongod启动会报错
openssl rand -base64 756 > mongodb.key
chmod 400 mongodb.key
````

#### 创建yml文件

配置文件：

````
version: '3.3'

services:
  mongodb1:
    image: mongo:4.2.1
    volumes:
      - /data/mongo/data/mongo1:/data/db
      - ./mongodb.key:/data/mongodb.key
    user: root
    environment:
      - MONGO_INITDB_ROOT_USERNAME=handle
      - MONGO_INITDB_ROOT_PASSWORD=jimeng2017
      - MONGO_INITDB_DATABASE=handle
    container_name: mongodb1
    ports:
      - 37017:27017
    command: mongod --replSet mongos --keyFile /data/mongodb.key
    restart: always
    entrypoint:
      - bash
      - -c
      - |
        chmod 400 /data/mongodb.key
        chown 999:999 /data/mongodb.key
        exec docker-entrypoint.sh $$@

  mongodb2:
    image: mongo:4.2.1
    volumes:
      - /data/mongo/data/mongo2:/data/db
      - ./mongodb.key:/data/mongodb.key
    user: root
    environment:
      - MONGO_INITDB_ROOT_USERNAME=handle
      - MONGO_INITDB_ROOT_PASSWORD=jimeng2017
      - MONGO_INITDB_DATABASE=handle
    container_name: mongodb2
    ports:
      - 37018:27017
    command: mongod --replSet mongos --keyFile /data/mongodb.key
    restart: always
    entrypoint:
      - bash
      - -c
      - |
        chmod 400 /data/mongodb.key
        chown 999:999 /data/mongodb.key
        exec docker-entrypoint.sh $$@

  mongodb3:
    image: mongo:4.2.1
    volumes:
      - /data/mongo/data/mongo3:/data/db
      - ./mongodb.key:/data/mongodb.key
    user: root
    environment:
      - MONGO_INITDB_ROOT_USERNAME=handle
      - MONGO_INITDB_ROOT_PASSWORD=jimeng2017
      - MONGO_INITDB_DATABASE=handle
    container_name: mongodb3
    ports:
      - 37019:27017
    command: mongod --replSet mongos --keyFile /data/mongodb.key
    restart: always
    entrypoint:
      - bash
      - -c
      - |
        chmod 400 /data/mongodb.key
        chown 999:999 /data/mongodb.key
        exec docker-entrypoint.sh $$@
````

然后启动

````
docker-compose -f  docker-compose-set.yml up -d
````

#### 初始化副本集

然后进去到第一个容器里面，初始化副本集

````
docker exec -it mongodb1 /bin/bash
````

登录

````
mongo -u 账号 -p 密码
````

登录成功可以查看状态

````
> rs.status()
{
	"ok" : 0,
	"errmsg" : "no replset config has been received",
	"code" : 94,
	"codeName" : "NotYetInitialized"
}
````

配置副本集

````
> rs.initiate({
...     _id: "mongos",
...     members: [
...         { _id : 0, host : "192.168.56.201:37017" },
...         { _id : 1, host : "192.168.56.201:37018" },
...         { _id : 2, host : "192.168.56.201:37019" },
            { _id : 3, host : "192.168.56.201:37020" }
...     ]
... });
{ "ok" : 1 }
````
上面提示ok就是表示成功了，这时候会选举出Primary节点。重新通过`rs.status()`查看状态就能看到了。  


