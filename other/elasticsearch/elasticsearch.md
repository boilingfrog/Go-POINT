## Elasticsearch

### 最近工作中用到了`Elasticsearch`，但是遇到几个挺坑的点，还是记录下。

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