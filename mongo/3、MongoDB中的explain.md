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


### 参考

【MongoDB简介】https://docs.mongoing.com/mongo-introduction      
【MySQL 中的索引】https://www.cnblogs.com/ricklz/p/17262747.html   
【performance-best-practices-indexing】https://www.mongodb.com/blog/post/performance-best-practices-indexing  
【tune_page_size_and_comp】https://source.wiredtiger.com/3.0.0/tune_page_size_and_comp.html       
【equality-sort-range-rule】https://www.mongodb.com/docs/manual/tutorial/equality-sort-range-rule/  
【使用索引来排序查询结果】https://mongoing.com/docs/tutorial/sort-results-with-indexes.html      