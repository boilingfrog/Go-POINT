## Elasticsearch

### 最近工作中用到了`Elasticsearch`，但是遇到几个挺坑的点，还是记录下。

### 深度分页的问题

es中的普通的查询`from+size`,存在查询数量的10000条限制。  


> index.max_result_window  
> The maximum value of from + size for searches to this index. Defaults to 10000. Search requests take heap memory and time proportional to from + size and this limits that memory. See Scroll or Search After for a more efficient alternative to raising this.
