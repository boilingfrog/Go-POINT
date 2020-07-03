## 通过docker-compose搭建mongo的replica set高可用

搭建一个mongo的集群，同时原来单机mongo的数据需要迁移到集群中。  

处理思路：单机mongo的数据通过`mongodump`备份，然后集群搭建起来了，在`mongorestore`导入到集群中，实现数据的迁移。  


### 备份数据


### 集群搭建

#### 生成keyFile

````
# 400权限是要保证安全性，否则mongod启动会报错
openssl rand -base64 756 > mongodb.key
chmod 400 mongodb.key
````

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

然后进去到第一个容器里面，初始化副本集

````
docker exec -it mongodb2 /bin/bash
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
...         { _id : 2, host : "192.168.56.201:37019" }
...     ]
... });
{ "ok" : 1 }
````
上面提示ok就是表示成功了，这时候会选举出Primary节点。


                            
                                          
