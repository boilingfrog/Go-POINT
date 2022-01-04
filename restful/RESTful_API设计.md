<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [RESTful API](#restful-api)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [设计的规范](#%E8%AE%BE%E8%AE%A1%E7%9A%84%E8%A7%84%E8%8C%83)
    - [协议](#%E5%8D%8F%E8%AE%AE)
    - [域名](#%E5%9F%9F%E5%90%8D)
    - [版本](#%E7%89%88%E6%9C%AC)
    - [路径](#%E8%B7%AF%E5%BE%84)
    - [HTTP动词](#http%E5%8A%A8%E8%AF%8D)
    - [GET 方法](#get-%E6%96%B9%E6%B3%95)
    - [POST 方法](#post-%E6%96%B9%E6%B3%95)
    - [PUT 方法](#put-%E6%96%B9%E6%B3%95)
    - [DELETE 方法](#delete-%E6%96%B9%E6%B3%95)
    - [PATCH 方法](#patch-%E6%96%B9%E6%B3%95)
  - [过滤信息](#%E8%BF%87%E6%BB%A4%E4%BF%A1%E6%81%AF)
  - [状态码](#%E7%8A%B6%E6%80%81%E7%A0%81)
  - [URL中不能有动词](#url%E4%B8%AD%E4%B8%8D%E8%83%BD%E6%9C%89%E5%8A%A8%E8%AF%8D)
  - [避免层级过深的URI](#%E9%81%BF%E5%85%8D%E5%B1%82%E7%BA%A7%E8%BF%87%E6%B7%B1%E7%9A%84uri)
  - [URL路径中首选小写字母](#url%E8%B7%AF%E5%BE%84%E4%B8%AD%E9%A6%96%E9%80%89%E5%B0%8F%E5%86%99%E5%AD%97%E6%AF%8D)
  - [URL路径名词均为复数](#url%E8%B7%AF%E5%BE%84%E5%90%8D%E8%AF%8D%E5%9D%87%E4%B8%BA%E5%A4%8D%E6%95%B0)
  - [幂等性](#%E5%B9%82%E7%AD%89%E6%80%A7)
    - [什么是幂等性](#%E4%BB%80%E4%B9%88%E6%98%AF%E5%B9%82%E7%AD%89%E6%80%A7)
  - [总结](#%E6%80%BB%E7%BB%93)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## RESTful API

### 前言
一直在使用RESTful API，但是好像概念还是很模糊的，总结下使用到的点

### 设计的规范

#### 协议

API与用户的通信协议，总是使用HTTPs协议。 

#### 域名

应该尽量将API部署在专用域名之下。

````
https://api.example.com
````

#### 版本

应该将API的版本号放入URL。

````
https://api.example.com/v1/
````

也可以将版本信息加入到HTTP头信息中，但不如放入URL方便和直观

#### 路径

在RESTful架构中，每个网址代表一种资源（resource），所以网址中不能有动词，只能有名词，而且所用的名词往往与
数据库的表格名对应。一般来说，数据库中的表都是同种记录的"集合"（collection），所以API中的名词也应该使用复数。

举例来说，有一个API提供动物园（zoo）的信息，还包括各种动物和雇员的信息，则它的路径应该设计成下面这样。

````
https://api.example.com/v1/zoos
https://api.example.com/v1/animals
https://api.example.com/v1/employees
````

#### HTTP动词

对于资源的具体操作类型，由HTTP动词表示。  

常用的HTTP动词有下面五个（括号里是对应的SQL命令）。   

````
GET（SELECT）：从服务器取出资源（一项或多项）。
POST（CREATE）：在服务器新建一个资源。
PUT（UPDATE）：在服务器更新资源（客户端提供改变后的完整资源）。
PATCH（UPDATE）：在服务器更新资源（客户端提供改变的属性）。
DELETE（DELETE）：从服务器删除资源。
````  
#### GET 方法

成功的 GET 方法通常返回 HTTP 状态代码 200（正常）。 如果找不到资源，该方法应返回 404（未找到）。

#### POST 方法

如果 POST 方法创建了新资源，则会返回 HTTP 状态代码 201（已创建）。 新资源的 URI 包含在响应的 Location 标头中。 响应正文包含资源的表示形式。  
如果该方法执行了一些处理但未创建新资源，则可以返回 HTTP 状态代码 200，并在响应正文中包含操作结果。 或者，如果没有可返回的结果，该方法可以返回 HTTP 状态代码 204（无内容）但不返回任何响应正文。  
如果客户端将无效数据放入请求，服务器应返回 HTTP 状态代码 400（错误的请求）。 响应正文可以包含有关错误的其他信息，或包含可提供更多详细信息的 URI 链接。  

#### PUT 方法

与 POST 方法一样，如果 PUT 方法创建了新资源，则会返回 HTTP 状态代码 201（已创建）。 如果该方法更新了现有资源，则会返回 200（正常）或 204（无内容）。 在某些情况下，可能无法更新现有资源。 在这种情况下，可考虑返回 HTTP 状态代码 409（冲突）。  
请考虑实现可批量更新集合中的多个资源的批量 HTTP PUT 操作。 PUT 请求应指定集合的 URI，而请求正文则应指定要修改的资源的详细信息。 此方法可帮助减少交互成本并提高性能。

#### DELETE 方法

如果删除操作成功，Web 服务器应以 HTTP 状态代码 204 做出响应，指示已成功处理该过程，但响应正文不包含其他信息。 如果资源不存在，Web 服务器可以返回 HTTP 404（未找到）。

#### PATCH 方法

PATCH方法是新引入的，是对PUT方法的补充，用来对已知资源进行局部更新

还有常用的HTTP动词  
````
HEAD：获取资源的元数据。
OPTIONS：获取信息，关于资源的哪些属性是客户端可以改变的。
````
下面是一些例子。  
````
GET /zoos：列出所有动物园
POST /zoos：新建一个动物园
GET /zoos/ID：获取某个指定动物园的信息
PUT /zoos/ID：更新某个指定动物园的信息（提供该动物园的全部信息）
PATCH /zoos/ID：更新某个指定动物园的信息（提供该动物园的部分信息）
DELETE /zoos/ID：删除某个动物园
GET /zoos/ID/animals：列出某个指定动物园的所有动物
DELETE /zoos/ID/animals/ID：删除某个指定动物园的指定动物
````
### 过滤信息
如果记录数量很多，服务器不可能都将它们返回给用户。API应该提供参数，过滤返回结果。  

下面是一些常见的参数。  

````
?limit=10：指定返回记录的数量
?offset=10：指定返回记录的开始位置。
?page=2&per_page=100：指定第几页，以及每页的记录数。
?sortby=name&order=asc：指定返回结果按照哪个属性排序，以及排序顺序。
?animal_type_id=1：指定筛选条件
````
参数的设计允许存在冗余，即允许API路径和URL参数偶尔有重复。比如，GET /zoo/ID/animals 与 GET /animals?zoo_id=ID 的含义是相同的。

### 状态码

服务器向用户返回的状态码和提示信息，常见的有以下一些（方括号中是该状态码对应的HTTP动词）。  
````
200 OK - [GET]：服务器成功返回用户请求的数据，该操作是幂等的（Idempotent）。
201 CREATED - [POST/PUT/PATCH]：用户新建或修改数据成功。
202 Accepted - [*]：表示一个请求已经进入后台排队（异步任务）
204 NO CONTENT - [DELETE]：用户删除数据成功。
400 INVALID REQUEST - [POST/PUT/PATCH]：用户发出的请求有错误，服务器没有进行新建或修改数据的操作，该操作是幂等的。
401 Unauthorized - [*]：表示用户没有权限（令牌、用户名、密码错误）。
403 Forbidden - [*] 表示用户得到授权（与401错误相对），但是访问是被禁止的。
404 NOT FOUND - [*]：用户发出的请求针对的是不存在的记录，服务器没有进行操作，该操作是幂等的。
406 Not Acceptable - [GET]：用户请求的格式不可得（比如用户请求JSON格式，但是只有XML格式）。
410 Gone -[GET]：用户请求的资源被永久删除，且不会再得到的。
422 Unprocesable entity - [POST/PUT/PATCH] 当创建一个对象时，发生一个验证错误。
500 INTERNAL SERVER ERROR - [*]：服务器发生错误，用户将无法判断发出的请求是否成功。
````
### URL中不能有动词

在Restful架构中，每个网址代表的是一种资源，所以网址中不能有动词，只能有名词，动词由HTTP的 get、post、put、delete
 四种方法来表示。

### 避免层级过深的URI

过深的导航容易导致url膨胀，不易维护，如`` GET /zoos/1/areas/3/animals/4``，尽量使用查询参数代替路径中的实体
导航，如``GET /animals?zoo=1&area=3``

### URL路径中首选小写字母

RFC 3986将URI定义为区分大小写，但scheme 和 host components除外。

### URL路径名词均为复数

为了保证url格式的一致性，建议使用复数形式。  

### 幂等性

#### 什么是幂等性
幂等性是指一次和多次请求某一个资源应该具有同样的效果。  


方法　   | 是否安全                                  | 是否幂等性
------- | ---------------------------------------- | -------------
GET     | TRUE                                     | True
PUT     | FALSE                                    | True
DELETE  | FALSE                                    | True
POST    | FALSE                                    | False
PATCH   | FALSE                                    | False

### 总结

api的设计还是一门很深的学问，是我们代码设计的一种体现。对于其中的重要几点，使用名词的复数，命名的简洁，api设计的
直观性。牢牢把握上面提到的几点，我们都可以设计出优雅的api。


### 参考  
【Microsoft Web API 设计】https://github.com/Microsoft/api-guidelines/blob/master/Guidelines.md#71-url-structure  
【RESTful API 设计指南】http://www.ruanyifeng.com/blog/2014/05/restful_api.html  
【Restful API 的设计规范】https://novoland.github.io/%E8%AE%BE%E8%AE%A1/2015/08/17/Restful%20API%20%E7%9A%84%E8%AE%BE%E8%AE%A1%E8%A7%84%E8%8C%83.html  