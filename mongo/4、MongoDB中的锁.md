<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [MongoDB 中的锁](#mongodb-%E4%B8%AD%E7%9A%84%E9%94%81)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [MongoDB 中锁的类型](#mongodb-%E4%B8%AD%E9%94%81%E7%9A%84%E7%B1%BB%E5%9E%8B)
  - [锁的让渡释放](#%E9%94%81%E7%9A%84%E8%AE%A9%E6%B8%A1%E9%87%8A%E6%94%BE)
  - [常见操作使用的锁类型](#%E5%B8%B8%E8%A7%81%E6%93%8D%E4%BD%9C%E4%BD%BF%E7%94%A8%E7%9A%84%E9%94%81%E7%B1%BB%E5%9E%8B)
  - [如果定位 MongoDB 中锁操作](#%E5%A6%82%E6%9E%9C%E5%AE%9A%E4%BD%8D-mongodb-%E4%B8%AD%E9%94%81%E6%93%8D%E4%BD%9C)
    - [1、查询运行超过20S 的请求](#1%E6%9F%A5%E8%AF%A2%E8%BF%90%E8%A1%8C%E8%B6%85%E8%BF%8720s-%E7%9A%84%E8%AF%B7%E6%B1%82)
    - [2、批量删除请求大于 20s 的请求](#2%E6%89%B9%E9%87%8F%E5%88%A0%E9%99%A4%E8%AF%B7%E6%B1%82%E5%A4%A7%E4%BA%8E-20s-%E7%9A%84%E8%AF%B7%E6%B1%82)
    - [3、kill 掉特定 client 端 ip 的请求](#3kill-%E6%8E%89%E7%89%B9%E5%AE%9A-client-%E7%AB%AF-ip-%E7%9A%84%E8%AF%B7%E6%B1%82)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## MongoDB 中的锁

### 前言

MongoDB 是一种常见的文档型数据库，因为其高性能、高可用、高扩展性等特点，被广泛应用于各种场景。   

在多线程的访问下，可能会出现多线程同时操作一个集合的情况，进而出现数据冲突的情况，为了保证数据的一致性，MongoDB 采用了锁机制来保证数据的一致性。   

下面来看看 MongoDB 中的锁机制。   

### MongoDB 中锁的类型

MongoDB 中使用多粒度锁定，它允许操作锁定在全局，数据库或集合级别，同时允许各个存储引擎在集合级别一下实现自己的并发控制(例如，WiredTiger 中的文档级别)。

MongoDB 中使用一个 readers-writer 锁，它允许并发多个读操作访问数据库，但是只提供唯一写操作访问。    

当一个读锁存在时，其它的读操作可以继续，不会被阻塞，如果一个写锁占有这个资源的时候，其它所有的读操作和写操作都会被阻塞。也就是读读不阻塞，读写阻塞，写写阻塞。    

MongoDB 中的锁首先提供了读写锁，即共享锁（Shared, S）（读锁）以及排他锁（Exclusive, X）（写锁），同时，为了解决多层级资源之间的互斥关系，提高多层级资源请求的效率，还在此基础上提供了意向锁（Intent Lock）。即锁可以划分为4中类型：   

1、共享锁，读锁（S），允许多个线程同时读取一个集合，读读不互斥；   

2、排他锁，写锁（X），允许一个线程写入数据，写写互斥，读写互斥；   

3、意向共享锁，IS，表示意向读取；  

4、意向排他锁，IX，表示意向写入；  

什么是意向锁呢？

如果另一个任务企图在某表级别上应用共享或排他锁，则受由第一个任务控制的表级别意向锁的阻塞，第二个任务在锁定该表前不需要检查各个页或行锁，而只需检查表上的意向锁。   

简单的讲就是意向锁是为了快速判断，表里面是否有记录被加锁。如果没有意向锁，大更新操作要判断是否有更小的操作在进行，

意向锁之间是不会产生冲突的，也不和 AUTO_INC 表锁冲突，它只会阻塞表级读锁或表级写锁，另外，意向锁也不会和行锁冲突，行锁只会和行锁冲突。

例如，当以写入模式（X模式）锁定集合时，相应的数据库锁和全局锁都必须以意向独占（IX）模式锁定。一个数据库可以同时以IS和IX模式进行锁定，但独占（X）锁无法与其他模式并存，共享（S）锁只能与意向共享（IS）锁并存。     

MongoDB 中的锁是公平的，所有的请求都会排队获取相应的锁。但是 MongoDB 为了优化吞吐量，在执行某个请求的时候，会同时执行和它兼容的其它请求。比如一个队列一个请求队列需要的锁如下，执行IS请求的同时，会同时执行和它相容的其他S和IS请求。等这一批请求的S锁释放后，再执行X锁的请求。   

```
IS → IS → X → X → S → IS
```

这种处理机制保证了在相对公平的前提下，提高了吞吐量，不会让某一类的请求长时间的等待。    

对于长时间的读或者写操作，某些条件下，mongodb 会临时的 让渡 锁，以防止长时间的阻塞。

### 锁的让渡释放   

对于大多数读取和写入操作，WiredTiger 使用乐观并发控制。WiredTiger 仅在全局、数据库和集合级别使用意向锁。当存储引擎检测到两个操作之间的冲突时，其中一个将导致写冲突，从而使 MongoDB 可以在不可见的情况下重新尝试该操作。  

在某些情况下，读写操作可以释放它们持有的锁。以防止长时间的阻塞。   

长时间运行的读取和写入操作，比如查询、更新和删除，在许多情况下都会释放锁。MongoDB操作也可以在影响多个文档的写入操作中，在单个文档修改之间释放锁。  

对于支持文档级并发控制的存储引擎，比如 WiredTiger，在访问存储时通常不需要释放锁，因为在全局、数据库和集合级别保持的意向锁不会阻塞其他读取者和写入者。然而，操作会定期释放锁，以便：

1、避免长时间存储事务，因为这些事务可能需要在内存中保存大量数据；  

2、充当中断点，以便可以终止长时间运行的操作；  

3、允许需要对集合进行独占访问的操作，比如索引/集合的删除和创建。  

### 常见操作使用的锁类型

下面列列举一下 MongoDB 中的常见操作对应的锁类型      

select：库级别的意向读锁(r)，表级别的意向读锁(r)，文档级别的读锁(R)；  

update/insert：库级别的意向写锁(w)，表级别的意向写锁(w)，文档级别的写锁(W)；  

foreground 方式创建索引：库级别的写锁(W)；当为一个集合创建索引时，因为是库级别的写锁，这个操作将阻塞其他的所有操作，任意基于所有数据库申请读或写锁都将等待直到前台完成索引创建操作；        

background 方式创建索引：库级别的意向写锁(w)，表级别的意向写锁(w)；  

### 如果定位 MongoDB 中锁操作  

当查询有慢查询出现的时候，有时候会出现锁的阻塞等待，紧急情况，需要我们快速定位并且结束当前操作。   

使用 `db.currentOp()` 就能查看当前数据库正在执行的操作。   

```
db.currentOp()

{
    "inprog" : [
        {
            "opid" : 6222,   #进程号
            "active" : true, #是否活动状态
            "secs_running" : 3,#操作运行了多少秒
            "microsecs_running" : NumberLong(3662328),#操作持续时间(以微秒为单位)。MongoDB通过从操作开始时减去当前时间来计算这个值。
            "op" : "getmore",#操作类型，包括(insert/query/update/remove/getmore/command)
            "ns" : "local.oplog.rs",#命名空间
            "query" : {#如果op是查询操作，这里将显示查询内容；也有说这里显示具体的操作语句的
                 
            },
            "client" : "192.168.91.132:45745",#连接的客户端信息
            "desc" : "conn5",#数据库的连接信息
            "threadId" : "0x7f1370cb4700",#线程ID
            "connectionId" : 5,#数据库的连接ID
            "waitingForLock" : false,#是否等待获取锁
            "numYields" : 0,
            "lockStats" : {
                "timeLockedMicros" : {#持有的锁时间微秒
                    "r" : NumberLong(141),#整个MongoDB实例的全局读锁
                    "w" : NumberLong(0)#整个MongoDB实例的全局写锁
                },
                "timeAcquiringMicros" : {#为了获得锁，等待的微秒时间
                    "r" : NumberLong(16),#整个MongoDB实例的全局读锁
                    "w" : NumberLong(0)#整个MongoDB实例的全局写锁
                }
            }
        }
    ]
}
```

来看下几个主要的字段含义    

- client：发起请求的客户端；  

- opid: 操作的唯一标识；  

- secs_running：该操作已经执行的时间，单位：微妙。如果该字段的返回值很大，就需要查询请求是否合理；   
 
- op：操作类型。通常是query、insert、update、delete、command中的一种；  

- query/ns：这个字段能看出是对哪个集合正在执行什么操作。    

当发现一个语句执行时间很久，影响到了整个数据库的运行，这时候我们 可以考虑中断这条语句的执行。   

使用 `db.killOp(opid)` 命令终止该请求。    

来个试验的栗子  

对表里面一个大表创建索引，不添加 backend。   

```
db.notifications.createIndex({userId: -1});
```

#### 1、查询运行超过20S 的请求  

```
$ db.currentOp({"active" : true, "secs_running":{ "$gt" : 20 }})

{
	"inprog" : [
		{
			"host" : "host-192-168-61-214:27017",
			"desc" : "conn50774156",
			"connectionId" : 50774156,
			"client" : "172.18.91.66:52088",
			"appName" : "Navicat",
			"clientMetadata" : {
				"application" : {
					"name" : "Navicat"
				},
				"driver" : {
					"name" : "mongoc",
					"version" : "1.16.2"
				},
				"os" : {
					"type" : "Darwin",
					"name" : "macOS",
					"version" : "20.6.0",
					"architecture" : "x86_64"
				},
				"platform" : "cfg=0x0000d6a0e9 posix=200112 stdc=201112 CC=clang 8.0.0 (clang-800.0.42.1) CFLAGS=\"\" LDFLAGS=\"\""
			},
			"active" : true,
			"currentOpTime" : "2023-10-24T01:32:00.615+0000",
			"opid" : -1782291565,
			"lsid" : {
				"id" : UUID("fff3c45d-b6ac-4a30-b83f-5a565ba166ef"),
				"uid" : BinData(0,"EJF4gS8MLpU7cuurTHswrdjF5hInXITH3796necT7PU=")
			},
			"secs_running" : NumberLong(103),
			"microsecs_running" : NumberLong(103025729),
			"op" : "command",
			"ns" : "gleeman.$cmd",
			"command" : {
				"createIndexes" : "notifications",
				"indexes" : [
					{
						"key" : {
							"userId" : -1
						},
						"name" : "userId_-1"
					}
				],
				"$db" : "gleeman",
				"lsid" : {
					"id" : UUID("fff3c45d-b6ac-4a30-b83f-5a565ba166ef")
				},
				"$clusterTime" : {
					"clusterTime" : Timestamp(1698111011, 1),
					"signature" : {
						"hash" : BinData(0,"iKilM1hvvIJC4hrTgu3FebYNhEw="),
						"keyId" : NumberLong("7233287468395528194")
					}
				}
			},
			"msg" : "Index Build (background) Index Build (background): 27288147/34043394 80%",
			"progress" : {
				"done" : 27288148,
				"total" : 34043394
			},
			"numYields" : 213205,
			"locks" : {
				"Global" : "w",
				"Database" : "w",
				"Collection" : "w"
			},
			"waitingForLock" : false,
			"lockStats" : {
				"Global" : {
					"acquireCount" : {
						"r" : NumberLong(213208),
						"w" : NumberLong(213208)
					}
				},
				"Database" : {
					"acquireCount" : {
						"w" : NumberLong(213209),
						"W" : NumberLong(1)
					}
				},
				"Collection" : {
					"acquireCount" : {
						"w" : NumberLong(213207)
					}
				},
				"oplog" : {
					"acquireCount" : {
						"w" : NumberLong(1)
					}
				}
			}
		}
	],
	"ok" : 1,
	"operationTime" : Timestamp(1698111118, 1),
	"$clusterTime" : {
		"clusterTime" : Timestamp(1698111118, 1),
		"signature" : {
			"hash" : BinData(0,"oLzIpVSGpZ213BW4x/jY6ESKvdA="),
			"keyId" : NumberLong("7233287468395528194")
		}
	}
}
```

#### 2、批量删除请求大于 20s 的请求

```
var ops = db.currentOp(
    {
        "active": true,
        "secs_running": {
            "$gt": 20
        }
    }
).inprog

for (i = 0; i < ops.length; i++) {
    ns = ops[i].ns;
    op = ops[i].op;
    if (ns.startsWith("system.") || ns.startsWith("local.oplog.") || ns.length === 0 || op == "none" || ns.command == "" || ns in["admin", "local", "config"]) {
        continue;
    }
    var opid = ops[i].opid;
    db.killOp(opid);
    print("Stopping op #" + opid)
}
```

####  3、kill 掉特定 client 端 ip 的请求

```
var clientIp="172.18.91.66";
var currOp = db.currentOp();

for (op in currOp.inprog) {
    if (clientIp == currOp.inprog[op].client.split(":")[0]) {
        db.killOp(currOp.inprog[op].opid)
    }
}
```


### 参考

【mongodb锁表命令-相关文档】https://www.volcengine.com/theme/900385-M-7-1   
【mongo 中的锁】https://www.jinmuinfo.com/community/MongoDB/docs/15-faq/03-concurrency.html  
【FAQ: Concurrency】https://www.mongodb.com/docs/manual/faq/concurrency/      