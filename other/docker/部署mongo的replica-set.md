- [通过docker-compose搭建mongo的replica set高可用](#%E9%80%9A%E8%BF%87docker-compose%E6%90%AD%E5%BB%BAmongo%E7%9A%84replica-set%E9%AB%98%E5%8F%AF%E7%94%A8)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [备份数据](#%E5%A4%87%E4%BB%BD%E6%95%B0%E6%8D%AE)
    - [备份数据到本地](#%E5%A4%87%E4%BB%BD%E6%95%B0%E6%8D%AE%E5%88%B0%E6%9C%AC%E5%9C%B0)
    - [数据恢复](#%E6%95%B0%E6%8D%AE%E6%81%A2%E5%A4%8D)
  - [集群搭建](#%E9%9B%86%E7%BE%A4%E6%90%AD%E5%BB%BA)
    - [生成keyFile](#%E7%94%9F%E6%88%90keyfile)
    - [创建yml文件](#%E5%88%9B%E5%BB%BAyml%E6%96%87%E4%BB%B6)
    - [初始化副本集](#%E5%88%9D%E5%A7%8B%E5%8C%96%E5%89%AF%E6%9C%AC%E9%9B%86)
  - [增加副本集](#%E5%A2%9E%E5%8A%A0%E5%89%AF%E6%9C%AC%E9%9B%86)
  - [测试连接](#%E6%B5%8B%E8%AF%95%E8%BF%9E%E6%8E%A5)

## 通过docker-compose搭建mongo的replica set高可用

### 前言 

搭建一个mongo的集群，同时原来单机mongo的数据需要迁移到集群中。  

处理思路：单机mongo的数据通过`mongodump`备份，然后集群搭建起来了，在`mongorestore`导入到集群中，实现数据的迁移。  

### 备份数据

备份(mongodump)与恢复(mongorestore)  

#### 备份数据到本地

```
mongodump -h 192.168.1.11 -u handle -p jimeng2017 -o /home/liz/Desktop/mongo-bei
```

#### 数据恢复

新的集群安装完成之后，恢复数据到Primary节点，集群会自动同步到副本集中(注意需要确认下Primary节点的端口)

````
mongorestore -h 192.168.1.11:37017 -u handle -p jimeng2017  /home/liz/Desktop/mongo-bei
````

注意：更换自己服务器上面的ip和mongo对应的账号密码

### 集群搭建

#### 生成keyFile

- MongoDB使用keyfile认证，副本集中的每个mongod实例使用keyfile内容作为认证其他成员的共享密码。mongod实例只有拥有正确的keyfile才可以加入副本集。
- keyFile的内容必须是6到1024个字符的长度，且副本集所有成员的keyFile内容必须相同。
- 有一点要注意是的：在UNIX系统中，keyFile必须没有组权限或完全权限（也就是权限要设置成X00的形式）。Windows系统中，keyFile权限没有被检查。
- 可以使用任意方法生成keyFile。例如，如下操作使用openssl生成复杂的随机的1024个字符串。然后使用chmod修改文件权限，只给文件拥有者提供读权限。

````
# 400权限是要保证安全性，否则mongod启动会报错
openssl rand -base64 756 > mongodb.key
chmod 400 mongodb.key
````

每一个副本集成员都要使用相同的keyFile文件

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

文件详解  
`chown 999:999 /data/mongodb.key` 999用户是容器中的mongod用户，通过chown修改文件用户权限  
`mongod --replSet mongos --keyFile /data/mongodb.key` 启动命令，`--replSet mongos` 以副本集形式启动并将副本集名字命名为 mongos ，`--keyFile /data/mongodb.key` 设置keyFile，用于副本集通信，文件通过 volumes 映射到容器内  

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
... ...     _id: "mongos",
... ...     members: [
... ...         { _id : 0, host : "192.168.1.11:37017" },
... ...         { _id : 1, host : "192.168.1.11:37018" },
... ...         { _id : 2, host : "192.168.1.11:37019" }
... ...     ]
... ... });
{ "ok" : 1 }

````
上面提示ok就是表示成功了，这时候会选举出Primary节点。重新通过`rs.status()`查看状态就能看到了。

```
rs.status()
{
	"set" : "mongos",
	"date" : ISODate("2020-07-04T13:02:44.676Z"),
	"myState" : 1,
	"term" : NumberLong(1),
	"syncingTo" : "",
	"syncSourceHost" : "",
	"syncSourceId" : -1,
	"heartbeatIntervalMillis" : NumberLong(2000),
	"majorityVoteCount" : 2,
	"writeMajorityCount" : 2,
	"optimes" : {
		"lastCommittedOpTime" : {
			"ts" : Timestamp(1593867760, 1),
			"t" : NumberLong(1)
		},
		"lastCommittedWallTime" : ISODate("2020-07-04T13:02:40.809Z"),
		"readConcernMajorityOpTime" : {
			"ts" : Timestamp(1593867760, 1),
			"t" : NumberLong(1)
		},
		"readConcernMajorityWallTime" : ISODate("2020-07-04T13:02:40.809Z"),
		"appliedOpTime" : {
			"ts" : Timestamp(1593867760, 1),
			"t" : NumberLong(1)
		},
		"durableOpTime" : {
			"ts" : Timestamp(1593867760, 1),
			"t" : NumberLong(1)
		},
		"lastAppliedWallTime" : ISODate("2020-07-04T13:02:40.809Z"),
		"lastDurableWallTime" : ISODate("2020-07-04T13:02:40.809Z")
	},
	"lastStableRecoveryTimestamp" : Timestamp(1593867720, 5),
	"lastStableCheckpointTimestamp" : Timestamp(1593867720, 5),
	"electionCandidateMetrics" : {
		"lastElectionReason" : "electionTimeout",
		"lastElectionDate" : ISODate("2020-07-04T13:02:00.300Z"),
		"termAtElection" : NumberLong(1),
		"lastCommittedOpTimeAtElection" : {
			"ts" : Timestamp(0, 0),
			"t" : NumberLong(-1)
		},
		"lastSeenOpTimeAtElection" : {
			"ts" : Timestamp(1593867709, 1),
			"t" : NumberLong(-1)
		},
		"numVotesNeeded" : 2,
		"priorityAtElection" : 1,
		"electionTimeoutMillis" : NumberLong(10000),
		"numCatchUpOps" : NumberLong(-29631936),
		"newTermStartDate" : ISODate("2020-07-04T13:02:00.787Z"),
		"wMajorityWriteAvailabilityDate" : ISODate("2020-07-04T13:02:01.528Z")
	},
	"members" : [
		{
			"_id" : 0,
			"name" : "192.168.1.11:37017",
			"ip" : "192.168.1.11",
			"health" : 1,
			"state" : 1,
			"stateStr" : "PRIMARY",
			"uptime" : 155,
			"optime" : {
				"ts" : Timestamp(1593867760, 1),
				"t" : NumberLong(1)
			},
			"optimeDate" : ISODate("2020-07-04T13:02:40Z"),
			"syncingTo" : "",
			"syncSourceHost" : "",
			"syncSourceId" : -1,
			"infoMessage" : "could not find member to sync from",
			"electionTime" : Timestamp(1593867720, 1),
			"electionDate" : ISODate("2020-07-04T13:02:00Z"),
			"configVersion" : 1,
			"self" : true,
			"lastHeartbeatMessage" : ""
		},
		{
			"_id" : 1,
			"name" : "192.168.1.11:37018",
			"ip" : "192.168.1.11",
			"health" : 1,
			"state" : 2,
			"stateStr" : "SECONDARY",
			"uptime" : 54,
			"optime" : {
				"ts" : Timestamp(1593867760, 1),
				"t" : NumberLong(1)
			},
			"optimeDurable" : {
				"ts" : Timestamp(1593867760, 1),
				"t" : NumberLong(1)
			},
			"optimeDate" : ISODate("2020-07-04T13:02:40Z"),
			"optimeDurableDate" : ISODate("2020-07-04T13:02:40Z"),
			"lastHeartbeat" : ISODate("2020-07-04T13:02:44.402Z"),
			"lastHeartbeatRecv" : ISODate("2020-07-04T13:02:43.511Z"),
			"pingMs" : NumberLong(0),
			"lastHeartbeatMessage" : "",
			"syncingTo" : "192.168.1.11:37017",
			"syncSourceHost" : "192.168.1.11:37017",
			"syncSourceId" : 0,
			"infoMessage" : "",
			"configVersion" : 1
		},
		{
			"_id" : 2,
			"name" : "192.168.1.11:37019",
			"ip" : "192.168.1.11",
			"health" : 1,
			"state" : 2,
			"stateStr" : "SECONDARY",
			"uptime" : 54,
			"optime" : {
				"ts" : Timestamp(1593867760, 1),
				"t" : NumberLong(1)
			},
			"optimeDurable" : {
				"ts" : Timestamp(1593867760, 1),
				"t" : NumberLong(1)
			},
			"optimeDate" : ISODate("2020-07-04T13:02:40Z"),
			"optimeDurableDate" : ISODate("2020-07-04T13:02:40Z"),
			"lastHeartbeat" : ISODate("2020-07-04T13:02:44.402Z"),
			"lastHeartbeatRecv" : ISODate("2020-07-04T13:02:43.533Z"),
			"pingMs" : NumberLong(0),
			"lastHeartbeatMessage" : "",
			"syncingTo" : "192.168.1.11:37017",
			"syncSourceHost" : "192.168.1.11:37017",
			"syncSourceId" : 0,
			"infoMessage" : "",
			"configVersion" : 1
		}
	],
	"ok" : 1,
	"$clusterTime" : {
		"clusterTime" : Timestamp(1593867760, 1),
		"signature" : {
			"hash" : BinData(0,"Zr91geLBqJM7xN3vVVMWnJYAexk="),
			"keyId" : NumberLong("6845609731550085123")
		}
	},
	"operationTime" : Timestamp(1593867760, 1)
}
```

通过`rs.status()`的输出我们就能分出那个是`PRIMARY`节点了。


### 增加副本集

之前已经初始化了三个mongo，我们在原来的基础上再增加一个副本集。  

首先编写`docker-compose-set.yml`增加mongodb4

````
mongodb4:
    image: mongo:4.2.1
    volumes:
      - /data/mongo/data/mongo4:/data/db
      - ./mongodb.key:/data/mongodb.key
    user: root
    environment:
      - MONGO_INITDB_ROOT_USERNAME=handle
      - MONGO_INITDB_ROOT_PASSWORD=jimeng2017
      - MONGO_INITDB_DATABASE=handle
    container_name: mongodb4
    ports:
      - 37020:27017
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

启动

````
docker-compose -f  docker-compose-set.yml up -d
````

然后进入到`PRIMARY`,执行

````
> rs.add("192.168.1.11:37020")
{
	"ok" : 1,
	"$clusterTime" : {
		"clusterTime" : Timestamp(1593868399, 1),
		"signature" : {
			"hash" : BinData(0,"BZJ2tCwFE1NvE22/LwGzFTWy+1M="),
			"keyId" : NumberLong("6845609731550085123")
		}
	},
	"operationTime" : Timestamp(1593868399, 1)
}
````

再次查看

````
> rs.status()
{
	"set" : "mongos",
	"date" : ISODate("2020-07-04T13:13:27.944Z"),
	"myState" : 1,
	"term" : NumberLong(3),
	"syncingTo" : "",
	"syncSourceHost" : "",
	"syncSourceId" : -1,
	"heartbeatIntervalMillis" : NumberLong(2000),
	"majorityVoteCount" : 3,
	"writeMajorityCount" : 3,
	"optimes" : {
		"lastCommittedOpTime" : {
			"ts" : Timestamp(1593868399, 1),
			"t" : NumberLong(3)
		},
		"lastCommittedWallTime" : ISODate("2020-07-04T13:13:19.243Z"),
		"readConcernMajorityOpTime" : {
			"ts" : Timestamp(1593868399, 1),
			"t" : NumberLong(3)
		},
		"readConcernMajorityWallTime" : ISODate("2020-07-04T13:13:19.243Z"),
		"appliedOpTime" : {
			"ts" : Timestamp(1593868399, 1),
			"t" : NumberLong(3)
		},
		"durableOpTime" : {
			"ts" : Timestamp(1593868399, 1),
			"t" : NumberLong(3)
		},
		"lastAppliedWallTime" : ISODate("2020-07-04T13:13:19.243Z"),
		"lastDurableWallTime" : ISODate("2020-07-04T13:13:19.243Z")
	},
	"lastStableRecoveryTimestamp" : Timestamp(1593868366, 1),
	"lastStableCheckpointTimestamp" : Timestamp(1593868366, 1),
	"electionCandidateMetrics" : {
		"lastElectionReason" : "electionTimeout",
		"lastElectionDate" : ISODate("2020-07-04T13:06:06.408Z"),
		"termAtElection" : NumberLong(3),
		"lastCommittedOpTimeAtElection" : {
			"ts" : Timestamp(0, 0),
			"t" : NumberLong(-1)
		},
		"lastSeenOpTimeAtElection" : {
			"ts" : Timestamp(1593867860, 1),
			"t" : NumberLong(1)
		},
		"numVotesNeeded" : 2,
		"priorityAtElection" : 1,
		"electionTimeoutMillis" : NumberLong(10000),
		"numCatchUpOps" : NumberLong(37018),
		"newTermStartDate" : ISODate("2020-07-04T13:06:06.753Z"),
		"wMajorityWriteAvailabilityDate" : ISODate("2020-07-04T13:06:07.255Z")
	},
	"members" : [
		{
			"_id" : 0,
			"name" : "192.168.1.11:37017",
			"ip" : "192.168.1.11",
			"health" : 1,
			"state" : 2,
			"stateStr" : "SECONDARY",
			"uptime" : 452,
			"optime" : {
				"ts" : Timestamp(1593868399, 1),
				"t" : NumberLong(3)
			},
			"optimeDurable" : {
				"ts" : Timestamp(1593868399, 1),
				"t" : NumberLong(3)
			},
			"optimeDate" : ISODate("2020-07-04T13:13:19Z"),
			"optimeDurableDate" : ISODate("2020-07-04T13:13:19Z"),
			"lastHeartbeat" : ISODate("2020-07-04T13:13:27.642Z"),
			"lastHeartbeatRecv" : ISODate("2020-07-04T13:13:27.493Z"),
			"pingMs" : NumberLong(0),
			"lastHeartbeatMessage" : "",
			"syncingTo" : "",
			"syncSourceHost" : "",
			"syncSourceId" : -1,
			"infoMessage" : "",
			"configVersion" : 2
		},
		{
			"_id" : 1,
			"name" : "192.168.1.11:37018",
			"ip" : "192.168.1.11",
			"health" : 1,
			"state" : 2,
			"stateStr" : "SECONDARY",
			"uptime" : 452,
			"optime" : {
				"ts" : Timestamp(1593868399, 1),
				"t" : NumberLong(3)
			},
			"optimeDurable" : {
				"ts" : Timestamp(1593868399, 1),
				"t" : NumberLong(3)
			},
			"optimeDate" : ISODate("2020-07-04T13:13:19Z"),
			"optimeDurableDate" : ISODate("2020-07-04T13:13:19Z"),
			"lastHeartbeat" : ISODate("2020-07-04T13:13:27.622Z"),
			"lastHeartbeatRecv" : ISODate("2020-07-04T13:13:27.493Z"),
			"pingMs" : NumberLong(13),
			"lastHeartbeatMessage" : "",
			"syncingTo" : "",
			"syncSourceHost" : "",
			"syncSourceId" : -1,
			"infoMessage" : "",
			"configVersion" : 2
		},
		{
			"_id" : 2,
			"name" : "192.168.1.11:37019",
			"ip" : "192.168.1.11",
			"health" : 1,
			"state" : 1,
			"stateStr" : "PRIMARY",
			"uptime" : 454,
			"optime" : {
				"ts" : Timestamp(1593868399, 1),
				"t" : NumberLong(3)
			},
			"optimeDate" : ISODate("2020-07-04T13:13:19Z"),
			"syncingTo" : "",
			"syncSourceHost" : "",
			"syncSourceId" : -1,
			"infoMessage" : "",
			"electionTime" : Timestamp(1593867966, 1),
			"electionDate" : ISODate("2020-07-04T13:06:06Z"),
			"configVersion" : 2,
			"self" : true,
			"lastHeartbeatMessage" : ""
		},
		{
			"_id" : 3,
			"name" : "192.168.1.11:37020",
			"ip" : "192.168.1.11",
			"health" : 1,
			"state" : 2,
			"stateStr" : "SECONDARY",
			"uptime" : 8,
			"optime" : {
				"ts" : Timestamp(1593868399, 1),
				"t" : NumberLong(3)
			},
			"optimeDurable" : {
				"ts" : Timestamp(1593868399, 1),
				"t" : NumberLong(3)
			},
			"optimeDate" : ISODate("2020-07-04T13:13:19Z"),
			"optimeDurableDate" : ISODate("2020-07-04T13:13:19Z"),
			"lastHeartbeat" : ISODate("2020-07-04T13:13:27.418Z"),
			"lastHeartbeatRecv" : ISODate("2020-07-04T13:13:27.823Z"),
			"pingMs" : NumberLong(0),
			"lastHeartbeatMessage" : "",
			"syncingTo" : "",
			"syncSourceHost" : "",
			"syncSourceId" : -1,
			"infoMessage" : "",
			"configVersion" : 2
		}
	],
	"ok" : 1,
	"$clusterTime" : {
		"clusterTime" : Timestamp(1593868399, 1),
		"signature" : {
			"hash" : BinData(0,"BZJ2tCwFE1NvE22/LwGzFTWy+1M="),
			"keyId" : NumberLong("6845609731550085123")
		}
	},
	"operationTime" : Timestamp(1593868399, 1)
}
````

可以看到最新的节点已经加进去了

### 测试连接

使用go测试下是否能正常连接和读写数据。

show dbs;：查看数据库

````
> show dbs
admin         0.000GB
bs-cx         0.000GB
config        0.000GB
handle        0.001GB
handle-core4  0.019GB
local         0.020GB
stud          0.000GB
````

use stud;：切换到指定数据库，如果不存在该数据库就创建。

````
> use stud
switched to db stud
````

db;：显示当前所在数据库。

````
> db
stud
````

db.dropDatabase()：删除当前数据库

````
> db.dropDatabase();
{ "ok" : 1 }
````

测试写入数据，使用的包`https://github.com/mongodb/mongo-go-driver`

````go
package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Student struct {
	Name string
	Age  int
}

func main() {
	uri := "mongodb://handle:jimeng2017@192.168.1.11:37017,192.168.1.11:37018,192.168.1.11:37019,192.168.1.11:37020/admin?replicaSet=mongos"
	// 设置客户端连接配置
	clientOptions := options.Client().ApplyURI(uri)

	// 连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	collection := client.Database("stud").Collection("student")

	s1 := Student{"小红", 12}
	insertResult, err := collection.InsertOne(context.TODO(), s1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	var result Student
	filter := bson.D{{"name", "小红"}}

	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found a single document: %+v\n", result)

	// 断开连接
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}
````

打印下输出

```go
Connected to MongoDB!
Inserted a single document:  ObjectID("5f00b6092aef651857151754")
Found a single document: {Name:小红 Age:12}
Connection to MongoDB closed.
```
