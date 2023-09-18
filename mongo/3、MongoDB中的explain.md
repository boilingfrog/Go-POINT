<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [MongoDB 中如何使用 explain 分析查询计划](#mongodb-%E4%B8%AD%E5%A6%82%E4%BD%95%E4%BD%BF%E7%94%A8-explain-%E5%88%86%E6%9E%90%E6%9F%A5%E8%AF%A2%E8%AE%A1%E5%88%92)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [查询计划 explain](#%E6%9F%A5%E8%AF%A2%E8%AE%A1%E5%88%92-explain)
    - [explain](#explain)
      - [1、queryPlanner](#1queryplanner)
      - [2、executionStats](#2executionstats)
      - [3、allPlansExecution](#3allplansexecution)
    - [indexfilter](#indexfilter)
    - [Stage 参数说明](#stage-%E5%8F%82%E6%95%B0%E8%AF%B4%E6%98%8E)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->


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
		"namespace" : "gleeman.test_explain", // 查询的命名空间，作用于那个库的那个表
		"indexFilterSet" : false, // 针对该query是否有 indexfilter，indexfilter 的作用见下文
		"parsedQuery" : { // 解析查询条件，即过滤条件是什么，此处是 age = 59
			"age" : {
				"$eq" : 59 
			}
		},
		"winningPlan" : { //查询优化器根据该 query 选择的最优的查询计划
			"stage" : "SORT", // 最优执行计划的，这里是 sort 表示在内存中具体的 stage 中的参数含义可见下文  
			"sortPattern" : {
				"_id" : -1
			},
			"inputStage" : { // 用来描述子 stage，并且为其父 stage 提供文档和索引关键字,这里面含有着执行计划中比较主要的信息
				"stage" : "SORT_KEY_GENERATOR", // 表示在内存中发生了排序
				"inputStage" : {
					"stage" : "FETCH", // FETCH 根据索引检索指定的文件
					"inputStage" : {
						"stage" : "IXSCAN", // stage 表示索引扫描
						"keyPattern" : { // 查询命中的索引
							"age" : -1
						},
						"indexName" : "age_-1", // 计算选择的索引的名字
						"isMultiKey" : false, // 是否为多键索引，因为用到的索引是单列索引，这里是 false
						"multiKeyPaths" : {
							"age" : [ ]
						},
						"isUnique" : false,
						"isSparse" : false,
						"isPartial" : false,
						"indexVersion" : 2,
						"direction" : "forward", // 此 query 的查询状态，forward 是升序，降序则是 backward
						"indexBounds" : { // 最优计划所扫描的索引范围
							"age" : [
								"[59.0, 59.0]" // [MinKey,MaxKey]
							]
						}
					}
				}
			}
		},
		"rejectedPlans" : [ // 其他计划，因为不是最优而被查询优化器拒绝(reject)  
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
	"serverInfo" : { // 服务器信息，包括主机名和ip，MongoDB的version等信息
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

上面的查询总结下来就是，`age = 59` 的查询使用到了 age 上面建立的索引，这块的匹配走的是 IXSCAN 也就是索引扫描。  

查询中还有一个 sort 的排序，这个排序动作是在内存中进行的，对应到的 stage 就是 SORT。  

因为 MongoDB 中内存排序对大数据量的排序效率不是很高，所以当有排序需求的时候，一般考虑创建组合索引，让排序在索引中完成。    

##### 2、executionStats

MongoDB 查询优化器会对当前的查询进行评估并且选择一个最佳的查询执行计划进行执行，在执行完毕后返回这个最佳执行计划执行完成时的相关统计信息，对于那些被拒绝的执行计划不返回器统计信息。  

```
$ db.getCollection("test_explain").find({"age" : 59}).sort({_id: -1}).explain("executionStats");  

{
	"queryPlanner" : {
		"plannerVersion" : 1,
		"namespace" : "gleeman.test_explain", // 查询的命名空间，作用于那个库那个表
		"indexFilterSet" : false, // 针对该query是否有 indexfilter，indexfilter 的作用见下文
		"parsedQuery" : { // 解析查询条件，即过滤条件是什么，此处是 age = 59
			"age" : {
				"$eq" : 59
			}
		},
		"winningPlan" : {  // 查询优化器针对该query所返回的最优执行计划的详细内容
			"stage" : "SORT", // 最优执行计划的，这里是 sort 表示在内存中排序，具体的 stage 中的参数含义可见下文  
			"sortPattern" : {
				"_id" : -1
			},
			"inputStage" : { // 用来描述子 stage，并且为其父 stage 提供文档和索引关键字,这里面含有着执行计划中比较主要的信息
				"stage" : "SORT_KEY_GENERATOR", // 表示在内存中发生了排序
				"inputStage" : {
					"stage" : "FETCH", // FETCH 根据索引检索指定的文件
					"inputStage" : { 
						"stage" : "IXSCAN", // 表示执行了索引扫描
						"keyPattern" : { // 查询命中的索引
							"age" : -1
						},
						"indexName" : "age_-1", // 查询选择的的索引的名字
						"isMultiKey" : false, // 是否为多键索引，因为用到的索引是单列索引，这里是 false
						"multiKeyPaths" : {
							"age" : [ ]
						},
						"isUnique" : false,
						"isSparse" : false,
						"isPartial" : false,
						"indexVersion" : 2,
						"direction" : "forward", // 此 query 的查询状态，forward 是升序，降序则是 backward
						"indexBounds" : { // 最优计划所扫描的索引范围
							"age" : [
								"[59.0, 59.0]"
							]
						}
					}
				}
			}
		},
		"rejectedPlans" : [ // 其他计划，因为不是最优而被查询优化器拒绝(reject) 
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
	"executionStats" : { // 已阶段树的形式纤细说明获胜执行计划的完成执行情况，例如，一个阶段可以有一个或者多个inputStage组成，每个阶段有特定于该阶段的执行信息组成。    
		"executionSuccess" : true, // 是否执行成功
		"nReturned" : 4786, // 此 query 匹配到的文档数
		"executionTimeMillis" : 26, // 查询计划选择和查询执行所需的总时间,单位：毫秒
		"totalKeysExamined" : 4786, // 扫描的索引条目数
		"totalDocsExamined" : 4786, // 扫描的文档数
		"executionStages" : { // 最优计划完整的执行信息
			"stage" : "SORT", // sort 表示在内存中发生了排序
			"nReturned" : 4786,
			"executionTimeMillisEstimate" : 20,
			"works" : 9575, // 指定查询执行阶段执行的“工作单元”的数量。查询执行将其工作划分为小单元。
			"advanced" : 4786, // 返回父阶段的结果数
			"needTime" : 4788, // 将中间结果返回给其父级的工作循环数
			"needYield" : 0, // 存储层请求查询系统产生锁定的次数
			"saveState" : 113,
			"restoreState" : 113,
			"isEOF" : 1,
			"invalidates" : 0,
			"sortPattern" : {
				"_id" : -1
			},
			"memUsage" : 362617,
			"memLimit" : 33554432,
			"inputStage" : { // 子执行单元，一个执行计划中，可以有一个或者多个inputStage
				"stage" : "SORT_KEY_GENERATOR",
				"nReturned" : 4786,
				"executionTimeMillisEstimate" : 20,
				"works" : 4788,
				"advanced" : 4786,
				"needTime" : 1,
				"needYield" : 0,
				"saveState" : 113,
				"restoreState" : 113,
				"isEOF" : 1,
				"invalidates" : 0,
				"inputStage" : {
					"stage" : "FETCH",
					"nReturned" : 4786,
					"executionTimeMillisEstimate" : 10,
					"works" : 4787,
					"advanced" : 4786,
					"needTime" : 0,
					"needYield" : 0,
					"saveState" : 113,
					"restoreState" : 113,
					"isEOF" : 1,
					"invalidates" : 0,
					"docsExamined" : 4786,
					"alreadyHasObj" : 0,
					"inputStage" : {
						"stage" : "IXSCAN",
						"nReturned" : 4786,
						"executionTimeMillisEstimate" : 10,
						"works" : 4787,
						"advanced" : 4786,
						"needTime" : 0,
						"needYield" : 0,
						"saveState" : 113,
						"restoreState" : 113,
						"isEOF" : 1,
						"invalidates" : 0,
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
						},
						"keysExamined" : 4786,
						"seeks" : 1,
						"dupsTested" : 0,
						"dupsDropped" : 0,
						"seenInvalidated" : 0
					}
				}
			}
		}
	},
	"serverInfo" : {
		"host" : "host-192-168-61-214",
		"port" : 27017,
		"version" : "4.0.3",
		"gitVersion" : "0377a277ee6d90364318b2f8d581f59c1a7abcd4"
	},
	"ok" : 1,
	"operationTime" : Timestamp(1694656526, 1),
	"$clusterTime" : {
		"clusterTime" : Timestamp(1694656526, 1),
		"signature" : {
			"hash" : BinData(0,"QaGoa0kgXBOfN9LA0px9V7NjU/c="),
			"keyId" : NumberLong("7233287468395528194")
		}
	}
}
```

##### 3、allPlansExecution  

该模式包括上述2种模式的所有信息，即按照最佳的执行计划执行以及列出统计信息，如果存在其他候选计划，也会列出这些候选的执行计划。  

下面的栗子就能看到，allPlansExecution 模式下，包含了上面两种查询模式的索引信息，这里不展开讨论了。    

```
$ db.getCollection("test_explain").find({"age" : 59}).sort({_id: -1}).explain("allPlansExecution");  

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
		"winningPlan" : {
			"stage" : "SORT",
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
	"executionStats" : {
		"executionSuccess" : true,
		"nReturned" : 4786,
		"executionTimeMillis" : 78,
		"totalKeysExamined" : 4786,
		"totalDocsExamined" : 4786,
		"executionStages" : {
			"stage" : "SORT",
			"nReturned" : 4786,
			"executionTimeMillisEstimate" : 60,
			"works" : 9575,
			"advanced" : 4786,
			"needTime" : 4788,
			"needYield" : 0,
			"saveState" : 113,
			"restoreState" : 113,
			"isEOF" : 1,
			"invalidates" : 0,
			"sortPattern" : {
				"_id" : -1
			},
			"memUsage" : 362617,
			"memLimit" : 33554432,
			"inputStage" : {
				"stage" : "SORT_KEY_GENERATOR",
				"nReturned" : 4786,
				"executionTimeMillisEstimate" : 50,
				"works" : 4788,
				"advanced" : 4786,
				"needTime" : 1,
				"needYield" : 0,
				"saveState" : 113,
				"restoreState" : 113,
				"isEOF" : 1,
				"invalidates" : 0,
				"inputStage" : {
					"stage" : "FETCH",
					"nReturned" : 4786,
					"executionTimeMillisEstimate" : 50,
					"works" : 4787,
					"advanced" : 4786,
					"needTime" : 0,
					"needYield" : 0,
					"saveState" : 113,
					"restoreState" : 113,
					"isEOF" : 1,
					"invalidates" : 0,
					"docsExamined" : 4786,
					"alreadyHasObj" : 0,
					"inputStage" : {
						"stage" : "IXSCAN",
						"nReturned" : 4786,
						"executionTimeMillisEstimate" : 0,
						"works" : 4787,
						"advanced" : 4786,
						"needTime" : 0,
						"needYield" : 0,
						"saveState" : 113,
						"restoreState" : 113,
						"isEOF" : 1,
						"invalidates" : 0,
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
						},
						"keysExamined" : 4786,
						"seeks" : 1,
						"dupsTested" : 0,
						"dupsDropped" : 0,
						"seenInvalidated" : 0
					}
				}
			}
		},
		"allPlansExecution" : [
			{
				"nReturned" : 101,
				"executionTimeMillisEstimate" : 60,
				"totalKeysExamined" : 4786,
				"totalDocsExamined" : 4786,
				"executionStages" : {
					"stage" : "SORT",
					"nReturned" : 101,
					"executionTimeMillisEstimate" : 60,
					"works" : 4889,
					"advanced" : 101,
					"needTime" : 4788,
					"needYield" : 0,
					"saveState" : 76,
					"restoreState" : 76,
					"isEOF" : 0,
					"invalidates" : 0,
					"sortPattern" : {
						"_id" : -1
					},
					"memUsage" : 362617,
					"memLimit" : 33554432,
					"inputStage" : {
						"stage" : "SORT_KEY_GENERATOR",
						"nReturned" : 4786,
						"executionTimeMillisEstimate" : 50,
						"works" : 4788,
						"advanced" : 4786,
						"needTime" : 1,
						"needYield" : 0,
						"saveState" : 76,
						"restoreState" : 76,
						"isEOF" : 1,
						"invalidates" : 0,
						"inputStage" : {
							"stage" : "FETCH",
							"nReturned" : 4786,
							"executionTimeMillisEstimate" : 50,
							"works" : 4787,
							"advanced" : 4786,
							"needTime" : 0,
							"needYield" : 0,
							"saveState" : 76,
							"restoreState" : 76,
							"isEOF" : 1,
							"invalidates" : 0,
							"docsExamined" : 4786,
							"alreadyHasObj" : 0,
							"inputStage" : {
								"stage" : "IXSCAN",
								"nReturned" : 4786,
								"executionTimeMillisEstimate" : 0,
								"works" : 4787,
								"advanced" : 4786,
								"needTime" : 0,
								"needYield" : 0,
								"saveState" : 76,
								"restoreState" : 76,
								"isEOF" : 1,
								"invalidates" : 0,
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
								},
								"keysExamined" : 4786,
								"seeks" : 1,
								"dupsTested" : 0,
								"dupsDropped" : 0,
								"seenInvalidated" : 0
							}
						}
					}
				}
			},
			{
				"nReturned" : 51,
				"executionTimeMillisEstimate" : 10,
				"totalKeysExamined" : 4889,
				"totalDocsExamined" : 4889,
				"executionStages" : {
					"stage" : "FETCH",
					"filter" : {
						"age" : {
							"$eq" : 59
						}
					},
					"nReturned" : 51,
					"executionTimeMillisEstimate" : 10,
					"works" : 4889,
					"advanced" : 51,
					"needTime" : 4838,
					"needYield" : 0,
					"saveState" : 113,
					"restoreState" : 113,
					"isEOF" : 0,
					"invalidates" : 0,
					"docsExamined" : 4889,
					"alreadyHasObj" : 0,
					"inputStage" : {
						"stage" : "IXSCAN",
						"nReturned" : 4889,
						"executionTimeMillisEstimate" : 10,
						"works" : 4889,
						"advanced" : 4889,
						"needTime" : 0,
						"needYield" : 0,
						"saveState" : 113,
						"restoreState" : 113,
						"isEOF" : 0,
						"invalidates" : 0,
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
						},
						"keysExamined" : 4889,
						"seeks" : 1,
						"dupsTested" : 0,
						"dupsDropped" : 0,
						"seenInvalidated" : 0
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
	"operationTime" : Timestamp(1694999650, 1),
	"$clusterTime" : {
		"clusterTime" : Timestamp(1694999650, 1),
		"signature" : {
			"hash" : BinData(0,"Ew5W9lYNoAtbo9zyErAbjbrqMlw="),
			"keyId" : NumberLong("7233287468395528194")
		}
	}
}
```

#### indexfilter

我们可以针对某些查询，指定特定的索引(索引必须存在)。如果查询条件吻合，就会使用指定的索引，如果指定了多个索引，会从中选出一个查询计划最优的索引执行。    

```
db.runCommand(
   {
      planCacheSetFilter: <collection>, // 数据表
      query: <query>, // 指定的查询条件
      sort: <sort>, // 指定的排序条件
      projection: <projection>,
      indexes: [ <index1>, <index2>, ...] // 指定查询条件对应的索引
   }
)
```

来个栗子   

```
db.runCommand(
   {
      planCacheSetFilter: "orders",
      query: { status: "A" },
      indexes: [
         { status: 1，cust_id: 1,  }
      ]
   }
)
```

上面对 orders 表建立了一个 indexFilter，指定了查询条件，在查询 orders 表中的 `status= A` 时，就会命中 indexFilter 中指定的 indexes。   

如果执行时通过 hint 指定了其他的 index，查询优化器将会忽略 hint 所设置 index，仍然使用 indexfilter 中设定的查询计划。    


#### Stage 参数说明  

来看下 Stage 中的参数   

|     类型                                       |      描述                                                 |           
| --------------------------------------------- | -------------------------------------------------------- | 
|     COLLSCAN                                  |      全表扫描                                              |           
|     IXSCAN                                    |      索引扫描                                              |           
|     FETCH                                     |      根据索引检索指定的文档                                  |           
|     SHARD_MERGE                               |      将各个分片的返回结果进行merge                           |           
|     SORT                                      |      表示在内存中进行了排序                                  |           
|     LIMIT                                     |      使用 limit 限制返回结果的数量                           |           
|     SKIP                                      |      使用 SKIP 进行跳过                                    |           
|     IDHACK                                    |      针对 _id 字段进行查询                                  |           
|     SHANRDING_FILTER                          |      通过 mongos 对分片数据进行查询                          |           
|     COUNT                                     |      利用 db.coll.explain().count() 进行 count 运算        |           
|     COUNTSCAN                                 |      count 不使用 index进行 count时的 stage 返回            |           
|     COUNT_SCAN                                |      count 使用了 index进行 count时的 stage 返回            |           
|     SUBPLA                                    |      未使用到索引的 $or 查询的 stage 返回                    |           
|     TEXT                                      |      使用全文索引进行查询时的 stage 返回                      |           
|     PROJECTION                                |      限定返回字段时候stage的返回                             |           

### 总结

这里总结了 MongoDB 中使用 explain 来判断我们创建的索引计划是否合理。  

其中常见的 explain 有三种模式,可以作为 explain 的参数进行模式选择

1、queryPlanner(默认模式)；queryPlanner 是 explain 默认的模式，queryPlanner 模式下不会真正的去执行 query 语句查询，查询优化器根据查询语句执行计划分析选出 `winning plan`。  

2、executionStats；MongoDB 查询优化器会对当前的查询进行评估并且选择一个最佳的查询执行计划进行执行，在执行完毕后返回这个最佳执行计划执行完成时的相关统计信息，对于那些被拒绝的执行计划不返回器统计信息。  

3、allPlansExecution；该模式包括上述2种模式的所有信息，即按照最佳的执行计划执行以及列出统计信息，如果存在其他候选计划，也会列出这些候选的执行计划。  

一般使用默认模式就能满足我们的需求了  

其中最核心的指标就是 Stage，通过这个参数，我们基本上就能判断我们创建的索引计划是否合理了。    

### 参考

【MongoDB简介】https://docs.mongoing.com/mongo-introduction      
【MySQL 中的索引】https://www.cnblogs.com/ricklz/p/17262747.html   
【performance-best-practices-indexing】https://www.mongodb.com/blog/post/performance-best-practices-indexing  
【tune_page_size_and_comp】https://source.wiredtiger.com/3.0.0/tune_page_size_and_comp.html       
【equality-sort-range-rule】https://www.mongodb.com/docs/manual/tutorial/equality-sort-range-rule/  
【使用索引来排序查询结果】https://mongoing.com/docs/tutorial/sort-results-with-indexes.html      
【Mongodb problem with sorting & indexes】https://groups.google.com/g/mongodb-user/c/YsY5h4KrwT4     
【MongoDB - 执行计划 】https://www.cnblogs.com/Neeo/articles/14326471.html     