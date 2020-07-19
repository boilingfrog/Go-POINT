<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [搭建mongo的replica set](#%E6%90%AD%E5%BB%BAmongo%E7%9A%84replica-set)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [安装](#%E5%AE%89%E8%A3%85)
    - [构建副本集](#%E6%9E%84%E5%BB%BA%E5%89%AF%E6%9C%AC%E9%9B%86)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

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

