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
>wqwqw