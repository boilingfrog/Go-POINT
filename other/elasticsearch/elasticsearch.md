<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [Elasticsearch](#elasticsearch)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [深度分页的问题](#%E6%B7%B1%E5%BA%A6%E5%88%86%E9%A1%B5%E7%9A%84%E9%97%AE%E9%A2%98)
    - [如何解决](#%E5%A6%82%E4%BD%95%E8%A7%A3%E5%86%B3)
      - [修改默认值](#%E4%BF%AE%E6%94%B9%E9%BB%98%E8%AE%A4%E5%80%BC)
      - [使用search_after方法](#%E4%BD%BF%E7%94%A8search_after%E6%96%B9%E6%B3%95)
      - [scroll 滚动搜索](#scroll-%E6%BB%9A%E5%8A%A8%E6%90%9C%E7%B4%A2)
  - [es中的近似聚合](#es%E4%B8%AD%E7%9A%84%E8%BF%91%E4%BC%BC%E8%81%9A%E5%90%88)
  - [总结](#%E6%80%BB%E7%BB%93)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Elasticsearch

### 前言
 
最近工作中用到了`Elasticsearch`，但是遇到几个挺坑的点，还是记录下。

### 深度分页的问题

es中的普通的查询`from+size`,存在查询数量的10000条限制。  

> index.max_result_window  
> The maximum value of from + size for searches to this index. Defaults to 10000. Search requests take heap memory and time proportional to from + size and this limits that memory. See Scroll or Search After for a more efficient alternative to raising this.

es为了减少内存的使用，限制了内存中索引数据的加载，默认`10000`。也就是
``
from 10000
size 1
``
这样的查询就是不行的，将会报错

````go
Result window is too large, from + size must be less than or equal to:[10000] but was [10500]. See the scroll api for a more efficient way to requestlarge data sets. This limit can be set by changing the[index.max_result_window] index level parameter
````

#### 如何解决

##### 修改默认值

通过设置index 的设置参数max_result_window的值，来改变查询条数的限制。

```go
curl -XPUT http://127.0.0.1:9200/book/_settings -d '{ "index" : { "max_result_window" : 200000000}}'
```

##### 使用search_after方法

例子可参考官方`https://www.elastic.co/guide/en/elasticsearch/reference/6.8/search-request-search-after.html`  

如何使用呢，就是设置一个全局唯一的字段，然后在查看的时候加上这个字段的排序。这样第二次查询，search_after第一次查询最后一条的对应的唯一值的值。有点绕哈，看例子  

```go
GET twitter/_search
{
    "size": 10,
    "query": {
        "match" : {
            "title" : "elasticsearch"
        }
    },
    "sort": [
        {"accessId": "asc"}      
    ]
}
```

比如我的`accessId`是全局唯一，并且自增的，第二次查询`search_after`最新的`accessId`就好了

```go
GET twitter/_search
{
    "size": 10,
    "query": {
        "match" : {
            "title" : "elasticsearch"
        }
    },
    "search_after": [10],
    "sort": [
        {"accessId": "asc"}
    ]
}
```
后面的查询依次类推

##### scroll 滚动搜索

`scroll` 查询可以用来对`Elasticsearch`有效地执行大批量的文档查询，而又不用付出深度分页那种代价。 

游标查询允许我们 先做查询初始化，然后再批量地拉取结果。 这有点儿像传统数据库中的`cursor` 。  

游标查询会取某个时间点的快照数据。查询初始化之后索引上的任何变化会被它忽略。它通过保存旧的数据文件来实现这个特性，结果就像保留初始化时的索引 视图 一样。  

深度分页的代价根源是结果集全局排序，如果去掉全局排序的特性的话查询结果的成本就会很低。游标查询用字段`_doc`来排序。 这个指令让 `Elasticsearch`仅仅从还有结果的分片返回下一批结果。  

启用游标查询可以通过在查询的时候设置参数` scroll`的值为我们期望的游标查询的过期时间。游标查询的过期时间会在每次做查询的时候刷新，所以这个时间只需要足够处理当前批的结果就可以了，而不是处理查询结果的所有文档的所需时间。 这个过期时间的参数很重要，因为保持这个游标查询窗口需要消耗资源，所以我们期望如果不再需要维护这种资源就该早点儿释放掉。设置这个超时能够让`Elasticsearch`在稍后空闲的时候自动释放这部分资源。   

```go
GET /old_index/_search?scroll=1m  // 设置查询窗口一分钟
{
    "query": { "match_all": {}},
    "sort" : ["_doc"], // 使用_doc字段排序
    "size":  1000
}
```

这个查询的返回结果包括一个字段`_scroll_id`，它是一个base64编码的长字符串 。现在我们能传递字段`_scroll_id`到`_search/scroll`查询接口获取下一批结果：

```go
GET /_search/scroll
{
    "scroll": "1m", // 时间
    "scroll_id" : "cXVlcnlUaGVuRmV0Y2g7NTsxMDk5NDpkUmpiR2FjOFNhNnlCM1ZDMWpWYnRROzEwOTk1OmRSamJHYWM4U2E2eUIzVkMxalZidFE7MTA5OTM6ZFJqYkdhYzhTYTZ5QjNWQzFqVmJ0UTsxMTE5MDpBVUtwN2lxc1FLZV8yRGVjWlI2QUVBOzEwOTk2OmRSamJHYWM4U2E2eUIzVkMxalZidFE7MDs="
} // scroll_id是上次返回的

```

之后的查询依次类推  

这个游标查询返回的下一批结果。 尽管我们指定字段`size`的值为1000，我们有可能取到超过这个值数量的文档。 当查询的时候， 字段 size 作用于单个分片，所以每个批次实际返回的文档数量最大为`size * number_of_primary_shards`。   

### es中的近似聚合

对于es来讲，其中的去重计数。是个近似的值，不像mysql中的是精确值，存在5%的误差，不过可以通过设置`precision_threshold`来解决少量数据的精准度

```go
GET /cars/transactions/_search
{
    "size" : 0,
    "aggs" : {
        "distinct_colors" : {
            "cardinality" : {
              "field" : "color",
              "precision_threshold" : 100  // precision_threshold 接受 0–40,000 之间的数字，更大的值还是会被当作 40,000 来处理。
            }
        }
    }
}
```

示例会确保当字段唯一值在 100 以内时会得到非常准确的结果。尽管算法是无法保证这点的，但如果基数在阈值以下，几乎总是`100%`正确的。高于阈值的基数会开始节省内存而牺牲准确度，同时也会对度量结果带入误差。  

对于指定的阈值,`HLL`的数据结构会大概使用`precision_threshold * 8`字节的内存，所以就必须在牺牲内存和获得额外的准确度间做平衡。  

在实际应用中，`100`的阈值可以在唯一值为百万的情况下仍然将误差维持`5%`以内。  

### 总结

当我们选型es时候，要充分考虑到上面的几点。