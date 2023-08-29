<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

## MongoDB 中如何使用 explain 分析查询计划

### 前言

创建完索引，如何分析索引的执行情况呢，MongoDB 中同样提供了 explain 来帮助我们进行分析，这里来看下 explain 的使用细节。   

### 查询计划 explain

为了更好的了解 explain 的使用，这里首先来准备一个测试的表结构，并且预先插入数据，方便后面进行查询计划的分析。   

```
for (var i = 1; i <1000000; i++) {
  db.test_explain.insert({name:'user'+i,age:parseInt(Math.random()*100+1),sn:Math.floor(Math.random()*10000000)})
}
```

创建索引  

```
db.test_explain.createIndex({age:-1},{background: true});
```

#### explain  

explain 有三种模式,可以作为 explain 的参数进行模式选择  

1、queryPlanner(默认模式)；  

2、executionStats；  

3、allPlansExecution；  

使用下面的查询条件来分析下    

```
db.getCollection("test_explain").find({"age" : 59}).sort({_id: -1})
```

##### 1、queryPlanner

queryPlanner 是 explain 默认的模式，queryPlanner 模式下不会真正的去执行 query 语句查询，查询优化器根据查询语句执行计划分析选出 `winning plan`。  

```
db.getCollection("test_explain").find({"age" : 59}).sort({_id: -1}).explain()

{
	"queryPlanner" : {
		"plannerVersion" : 1,
		"namespace" : "gleeman.test_explain",
		"indexFilterSet" : false,
		"parsedQuery" : {
			"age" : {
				"$eq" : 59
			}
		},
		"winningPlan" : { //查询优化器根据该 query 选择的最优的查询计划
			"stage" : "SORT", // 最优执行计划的 stage
			"sortPattern" : {
				"_id" : -1
			},
			"inputStage" : {
				"stage" : "SORT_KEY_GENERATOR",
				"inputStage" : {
					"stage" : "FETCH",
					"inputStage" : {
						"stage" : "IXSCAN",
						"keyPattern" : {
							"age" : -1
						},
						"indexName" : "age_-1",
						"isMultiKey" : false,
						"multiKeyPaths" : {
							"age" : [ ]
						},
						"isUnique" : false,
						"isSparse" : false,
						"isPartial" : false,
						"indexVersion" : 2,
						"direction" : "forward",
						"indexBounds" : {
							"age" : [
								"[59.0, 59.0]"
							]
						}
					}
				}
			}
		},
		"rejectedPlans" : [
			{
				"stage" : "FETCH",
				"filter" : {
					"age" : {
						"$eq" : 59
					}
				},
				"inputStage" : {
					"stage" : "IXSCAN",
					"keyPattern" : {
						"_id" : 1
					},
					"indexName" : "_id_",
					"isMultiKey" : false,
					"multiKeyPaths" : {
						"_id" : [ ]
					},
					"isUnique" : true,
					"isSparse" : false,
					"isPartial" : false,
					"indexVersion" : 2,
					"direction" : "backward",
					"indexBounds" : {
						"_id" : [
							"[MaxKey, MinKey]"
						]
					}
				}
			}
		]
	},
	"serverInfo" : {
		"host" : "host-192-168-61-214",
		"port" : 27017,
		"version" : "4.0.3",
		"gitVersion" : "0377a277ee6d90364318b2f8d581f59c1a7abcd4"
	},
	"ok" : 1,
	"operationTime" : Timestamp(1693278617, 5),
	"$clusterTime" : {
		"clusterTime" : Timestamp(1693278617, 5),
		"signature" : {
			"hash" : BinData(0,"1fdWNqFjbdgCwBm3/NFEpD9yXjI="),
			"keyId" : NumberLong("7233287468395528194")
		}
	}
}
```

#### Stage 参数说明  

来看下 Stage 中的参数   

- COLLSCAN：表示执行是全表扫描

- IXSCAN #索引扫描

- FETCH #根据索引去检索指定document

SHARD_MERGE #将各个分片返回数据进行merge

SORT #表明在内存中进行了排序（与老版本的scanAndOrder:true一致）

LIMIT #使用limit限制返回数

SKIP #使用skip进行跳过

IDHACK #针对_id进行查询

SHARDING_FILTER #通过mongos对分片数据进行查询

COUNT #利用db.coll.explain().count()之类进行count运算

COUNTSCAN #count不使用Index进行count时的stage返回

COUNT_SCAN #count使用了Index进行count时的stage返回

SUBPLA #未使用到索引的$or查询的stage返回

TEXT #使用全文索引进行查询时候的stage返回

PROJECTION #限定返回字段时候stage的返回





### 参考

【MongoDB简介】https://docs.mongoing.com/mongo-introduction      
【MySQL 中的索引】https://www.cnblogs.com/ricklz/p/17262747.html   
【performance-best-practices-indexing】https://www.mongodb.com/blog/post/performance-best-practices-indexing  
【tune_page_size_and_comp】https://source.wiredtiger.com/3.0.0/tune_page_size_and_comp.html       
【equality-sort-range-rule】https://www.mongodb.com/docs/manual/tutorial/equality-sort-range-rule/  
【使用索引来排序查询结果】https://mongoing.com/docs/tutorial/sort-results-with-indexes.html      