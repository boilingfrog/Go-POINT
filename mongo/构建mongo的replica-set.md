<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [搭建mongo的replica set](#%E6%90%AD%E5%BB%BAmongo%E7%9A%84replica-set)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [安装](#%E5%AE%89%E8%A3%85)
    - [构建副本集](#%E6%9E%84%E5%BB%BA%E5%89%AF%E6%9C%AC%E9%9B%86)
    - [加入认证](#%E5%8A%A0%E5%85%A5%E8%AE%A4%E8%AF%81)
  - [备份数据](#%E5%A4%87%E4%BB%BD%E6%95%B0%E6%8D%AE)
    - [备份数据到本地](#%E5%A4%87%E4%BB%BD%E6%95%B0%E6%8D%AE%E5%88%B0%E6%9C%AC%E5%9C%B0)
    - [数据恢复](#%E6%95%B0%E6%8D%AE%E6%81%A2%E5%A4%8D)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 搭建mongo的replica set

### 前言

准备三台机器，相互可以访问的。处理思路，先构建无需认证的集群，然后进入主节点，初始化集群的账号密码。然后开启所有机器的认证。  

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
// 存放数据
# mkdir data/db
// 存放log
# mkdir data/logs
# touch data/mongodb.log
# cd data
# mv mongdb.log logs
````

写入配置文件，`vi mongodb.cof`

````
// 端口
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
#keyFile = /usr/local/mongodb/data/mongodb.key
````

关于ip

- bind_ip = 192.168.0.136     #如果修改成本机Ip，那除了本机外的机器都可以连接（就是自己连不了、哈哈、蛋疼）  
- bind_ip = 0.0.0.0           #改成0，那么大家都可以访问（共赢）  
- bind_ip = 127.0.0.1         #改成127，那就只能自己练了（独吞）  

所以为了方便其他服务器和自己连接，就把`bind_ip`改成`0.0.0.0`  

启动 

````
# ./bin/mongod -f  ./data/mongodb.conf
````

三台机器都安装下

#### 构建副本集

进去其中一台机器

````
./bin/mongo 192.168.56.101:27018
````

查看副本集的状态

````
rs.status()
{
	"operationTime" : Timestamp(0, 0),
	"ok" : 0,
	"errmsg" : "no replset config has been received",
	"code" : 94,
	"codeName" : "NotYetInitialized",
	"$clusterTime" : {
		"clusterTime" : Timestamp(0, 0),
		"signature" : {
			"hash" : BinData(0,"AAAAAAAAAAAAAAAAAAAAAAAAAAA="),
			"keyId" : NumberLong(0)
		}
	}
}
````

发现没有初始化，然后我们来初始化

````
 rs.initiate({
     _id: "mongos",
     members: [
        { _id : 0, host : "192.168.56.101:27018" },
        { _id : 1, host : "192.168.56.102:27018" },
        { _id : 2, host : "192.168.56.103:27018" }
      ]
 });
````

再次查看状态

````
rs.status()
{
	"set" : "mongos",
	"date" : ISODate("2020-07-19T14:49:52.387Z"),
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
			"ts" : Timestamp(1595170188, 2),
			"t" : NumberLong(1)
		},
		"lastCommittedWallTime" : ISODate("2020-07-19T14:49:48.999Z"),
		"readConcernMajorityOpTime" : {
			"ts" : Timestamp(1595170188, 2),
			"t" : NumberLong(1)
		},
		"readConcernMajorityWallTime" : ISODate("2020-07-19T14:49:48.999Z"),
		"appliedOpTime" : {
			"ts" : Timestamp(1595170188, 2),
			"t" : NumberLong(1)
		},
		"durableOpTime" : {
			"ts" : Timestamp(1595170188, 2),
			"t" : NumberLong(1)
		},
		"lastAppliedWallTime" : ISODate("2020-07-19T14:49:48.999Z"),
		"lastDurableWallTime" : ISODate("2020-07-19T14:49:48.999Z")
	},
	"lastStableRecoveryTimestamp" : Timestamp(1595170187, 3),
	"lastStableCheckpointTimestamp" : Timestamp(1595170187, 3),
	"electionCandidateMetrics" : {
		"lastElectionReason" : "electionTimeout",
		"lastElectionDate" : ISODate("2020-07-19T14:49:47.525Z"),
		"termAtElection" : NumberLong(1),
		"lastCommittedOpTimeAtElection" : {
			"ts" : Timestamp(0, 0),
			"t" : NumberLong(-1)
		},
		"lastSeenOpTimeAtElection" : {
			"ts" : Timestamp(1595170176, 1),
			"t" : NumberLong(-1)
		},
		"numVotesNeeded" : 2,
		"priorityAtElection" : 1,
		"electionTimeoutMillis" : NumberLong(10000),
		"numCatchUpOps" : NumberLong(269553681),
		"newTermStartDate" : ISODate("2020-07-19T14:49:47.939Z"),
		"wMajorityWriteAvailabilityDate" : ISODate("2020-07-19T14:49:48.897Z")
	},
	"members" : [
		{
			"_id" : 0,
			"name" : "192.168.56.101:27018",
			"ip" : "192.168.56.101",
			"health" : 1,
			"state" : 1,
			"stateStr" : "PRIMARY",
			"uptime" : 1492,
			"optime" : {
				"ts" : Timestamp(1595170188, 2),
				"t" : NumberLong(1)
			},
			"optimeDate" : ISODate("2020-07-19T14:49:48Z"),
			"syncingTo" : "",
			"syncSourceHost" : "",
			"syncSourceId" : -1,
			"infoMessage" : "could not find member to sync from",
			"electionTime" : Timestamp(1595170187, 1),
			"electionDate" : ISODate("2020-07-19T14:49:47Z"),
			"configVersion" : 1,
			"self" : true,
			"lastHeartbeatMessage" : ""
		},
		{
			"_id" : 1,
			"name" : "192.168.56.102:27018",
			"ip" : "192.168.56.102",
			"health" : 1,
			"state" : 2,
			"stateStr" : "SECONDARY",
			"uptime" : 15,
			"optime" : {
				"ts" : Timestamp(1595170188, 2),
				"t" : NumberLong(1)
			},
			"optimeDurable" : {
				"ts" : Timestamp(1595170188, 2),
				"t" : NumberLong(1)
			},
			"optimeDate" : ISODate("2020-07-19T14:49:48Z"),
			"optimeDurableDate" : ISODate("2020-07-19T14:49:48Z"),
			"lastHeartbeat" : ISODate("2020-07-19T14:49:51.552Z"),
			"lastHeartbeatRecv" : ISODate("2020-07-19T14:49:50.817Z"),
			"pingMs" : NumberLong(1),
			"lastHeartbeatMessage" : "",
			"syncingTo" : "192.168.56.101:27018",
			"syncSourceHost" : "192.168.56.101:27018",
			"syncSourceId" : 0,
			"infoMessage" : "",
			"configVersion" : 1
		},
		{
			"_id" : 2,
			"name" : "192.168.56.103:27018",
			"ip" : "192.168.56.103",
			"health" : 1,
			"state" : 2,
			"stateStr" : "SECONDARY",
			"uptime" : 15,
			"optime" : {
				"ts" : Timestamp(1595170188, 2),
				"t" : NumberLong(1)
			},
			"optimeDurable" : {
				"ts" : Timestamp(1595170188, 2),
				"t" : NumberLong(1)
			},
			"optimeDate" : ISODate("2020-07-19T14:49:48Z"),
			"optimeDurableDate" : ISODate("2020-07-19T14:49:48Z"),
			"lastHeartbeat" : ISODate("2020-07-19T14:49:51.555Z"),
			"lastHeartbeatRecv" : ISODate("2020-07-19T14:49:50.830Z"),
			"pingMs" : NumberLong(0),
			"lastHeartbeatMessage" : "",
			"syncingTo" : "192.168.56.101:27018",
			"syncSourceHost" : "192.168.56.101:27018",
			"syncSourceId" : 0,
			"infoMessage" : "",
			"configVersion" : 1
		}
	],
	"ok" : 1,
	"$clusterTime" : {
		"clusterTime" : Timestamp(1595170188, 2),
		"signature" : {
			"hash" : BinData(0,"AAAAAAAAAAAAAAAAAAAAAAAAAAA="),
			"keyId" : NumberLong(0)
		}
	},
	"operationTime" : Timestamp(1595170188, 2)
}
````

发现已经初始化好了 

#### 加入认证

进去到PRIMARY节点初始化集群的登录账号和密码信息

````
# ./bin/mongo 192.168.56.101:27018
# use admin
# db.createUser({user: 'handle', pwd: 'jimeng2017', roles: ['root']})
````

生成`keyfile`  

- MongoDB使用keyfile认证，副本集中的每个mongod实例使用keyfile内容作为认证其他成员的共享密码。mongod实例只有拥有正确的keyfile才可以加入副本集。
- keyFile的内容必须是6到1024个字符的长度，且副本集所有成员的keyFile内容必须相同。
- 有一点要注意是的：在UNIX系统中，keyFile必须没有组权限或完全权限（也就是权限要设置成X00的形式）。Windows系统中，keyFile权限没有被检查。
- 可以使用任意方法生成keyFile。例如，如下操作使用openssl生成复杂的随机的1024个字符串。然后使用chmod修改文件权限，只给文件拥有者提供读权限。

````
# 400权限是要保证安全性，否则mongod启动会报错
openssl rand -base64 756 > mongodb.key
chmod 400 mongodb.key
````

然后放到mongodb中的data目录，三台机器`keyfile`要一致。我是在一台中生成，然后传到其他的服务器中。  

打开认证，三台机器都要执行  

````
#开启认证
auth = true
#安全文件地址
keyFile = /usr/local/mongodb/data/mongodb.key
````

之后重启mongo  

````
# ./bin/mongo 192.168.56.102:27018
rs.status()
{
	"operationTime" : Timestamp(1595212820, 1),
	"ok" : 0,
	"errmsg" : "command replSetGetStatus requires authentication",
	"code" : 13,
	"codeName" : "Unauthorized",
	"$clusterTime" : {
		"clusterTime" : Timestamp(1595212820, 1),
		"signature" : {
			"hash" : BinData(0,"7N/yM+2RounFZVTIzjbW+rEcZNs="),
			"keyId" : NumberLong("6851374857561571330")
		}
	}
}
````

不带账号密码登录，提示`Unauthorized`

````
# ./bin/mongo 192.168.56.102:27018 -u "handle" -p "jimeng2017" --authenticationDatabase "admin"
rs.status()
{
	"set" : "mongos",
	"date" : ISODate("2020-07-20T02:44:20.836Z"),
	"myState" : 1,
	"term" : NumberLong(3),
	"syncingTo" : "",
	"syncSourceHost" : "",
	"syncSourceId" : -1,
	"heartbeatIntervalMillis" : NumberLong(2000),
	"majorityVoteCount" : 2,
	"writeMajorityCount" : 2,
	"optimes" : {
		"lastCommittedOpTime" : {
			"ts" : Timestamp(1595213060, 1),
			"t" : NumberLong(3)
		},
		"lastCommittedWallTime" : ISODate("2020-07-20T02:44:20.154Z"),
		"readConcernMajorityOpTime" : {
			"ts" : Timestamp(1595213060, 1),
			"t" : NumberLong(3)
		},
		"readConcernMajorityWallTime" : ISODate("2020-07-20T02:44:20.154Z"),
		"appliedOpTime" : {
			"ts" : Timestamp(1595213060, 1),
			"t" : NumberLong(3)
		},
		"durableOpTime" : {
			"ts" : Timestamp(1595213060, 1),
			"t" : NumberLong(3)
		},
		"lastAppliedWallTime" : ISODate("2020-07-20T02:44:20.154Z"),
		"lastDurableWallTime" : ISODate("2020-07-20T02:44:20.154Z")
	},
	"lastStableRecoveryTimestamp" : Timestamp(1595213010, 1),
	"lastStableCheckpointTimestamp" : Timestamp(1595213010, 1),
	"electionCandidateMetrics" : {
		"lastElectionReason" : "electionTimeout",
		"lastElectionDate" : ISODate("2020-07-20T01:56:39.170Z"),
		"termAtElection" : NumberLong(3),
		"lastCommittedOpTimeAtElection" : {
			"ts" : Timestamp(0, 0),
			"t" : NumberLong(-1)
		},
		"lastSeenOpTimeAtElection" : {
			"ts" : Timestamp(1595210160, 1),
			"t" : NumberLong(2)
		},
		"numVotesNeeded" : 2,
		"priorityAtElection" : 1,
		"electionTimeoutMillis" : NumberLong(10000),
		"numCatchUpOps" : NumberLong(808464432),
		"newTermStartDate" : ISODate("2020-07-20T01:56:39.976Z"),
		"wMajorityWriteAvailabilityDate" : ISODate("2020-07-20T01:56:40.667Z")
	},
	"members" : [
		{
			"_id" : 0,
			"name" : "192.168.56.101:27018",
			"ip" : "192.168.56.101",
			"health" : 1,
			"state" : 2,
			"stateStr" : "SECONDARY",
			"uptime" : 2871,
			"optime" : {
				"ts" : Timestamp(1595213050, 1),
				"t" : NumberLong(3)
			},
			"optimeDurable" : {
				"ts" : Timestamp(1595213050, 1),
				"t" : NumberLong(3)
			},
			"optimeDate" : ISODate("2020-07-20T02:44:10Z"),
			"optimeDurableDate" : ISODate("2020-07-20T02:44:10Z"),
			"lastHeartbeat" : ISODate("2020-07-20T02:44:19.053Z"),
			"lastHeartbeatRecv" : ISODate("2020-07-20T02:44:20.779Z"),
			"pingMs" : NumberLong(1),
			"lastHeartbeatMessage" : "",
			"syncingTo" : "192.168.56.102:27018",
			"syncSourceHost" : "192.168.56.102:27018",
			"syncSourceId" : 1,
			"infoMessage" : "",
			"configVersion" : 1
		},
		{
			"_id" : 1,
			"name" : "192.168.56.102:27018",
			"ip" : "192.168.56.102",
			"health" : 1,
			"state" : 1,
			"stateStr" : "PRIMARY",
			"uptime" : 2873,
			"optime" : {
				"ts" : Timestamp(1595213060, 1),
				"t" : NumberLong(3)
			},
			"optimeDate" : ISODate("2020-07-20T02:44:20Z"),
			"syncingTo" : "",
			"syncSourceHost" : "",
			"syncSourceId" : -1,
			"infoMessage" : "",
			"electionTime" : Timestamp(1595210199, 1),
			"electionDate" : ISODate("2020-07-20T01:56:39Z"),
			"configVersion" : 1,
			"self" : true,
			"lastHeartbeatMessage" : ""
		},
		{
			"_id" : 2,
			"name" : "192.168.56.103:27018",
			"ip" : "192.168.56.103",
			"health" : 1,
			"state" : 2,
			"stateStr" : "SECONDARY",
			"uptime" : 2831,
			"optime" : {
				"ts" : Timestamp(1595213050, 1),
				"t" : NumberLong(3)
			},
			"optimeDurable" : {
				"ts" : Timestamp(1595213050, 1),
				"t" : NumberLong(3)
			},
			"optimeDate" : ISODate("2020-07-20T02:44:10Z"),
			"optimeDurableDate" : ISODate("2020-07-20T02:44:10Z"),
			"lastHeartbeat" : ISODate("2020-07-20T02:44:19.069Z"),
			"lastHeartbeatRecv" : ISODate("2020-07-20T02:44:19.689Z"),
			"pingMs" : NumberLong(0),
			"lastHeartbeatMessage" : "",
			"syncingTo" : "192.168.56.101:27018",
			"syncSourceHost" : "192.168.56.101:27018",
			"syncSourceId" : 0,
			"infoMessage" : "",
			"configVersion" : 1
		}
	],
	"ok" : 1,
	"$clusterTime" : {
		"clusterTime" : Timestamp(1595213060, 1),
		"signature" : {
			"hash" : BinData(0,"vWZEW8RyQOU7IwSaLMmrEancUio="),
			"keyId" : NumberLong("6851374857561571330")
		}
	},
	"operationTime" : Timestamp(1595213060, 1)
}
````

账号密码设置成功了

### 备份数据

备份(mongodump)与恢复(mongorestore)  

#### 备份数据到本地

``
mongodump -h 192.168.56.101:27018 -u handle -p jimeng2017 -o /home/liz/Desktop/mongo-bei
``

#### 数据恢复

新的集群安装完成之后，恢复数据到Primary节点，集群会自动同步到副本集中

````
mongorestore -h 192.168.56.101:27018 -u handle -p jimeng2017  /home/liz/Desktop/mongo-bei
````

注意：更换自己服务器上面的ip和mongo对应的账号密码  


