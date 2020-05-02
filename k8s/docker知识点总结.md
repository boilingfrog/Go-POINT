## docker知识点总结

### docker



### docker中的网络

docker中支持5种网络模式  

#### bridge

默认网络，Docker启动后默认创建一个docker0网桥，默认创建的容器也是添加到这个网桥中。 

#### host 

容器不会获得一个独立的network namespace，而是与宿主机共用一个。

#### none

获取独立的network namespace，但不为容器进行任何网络配置。

#### container

与指定的容器使用同一个network namespace，网卡配置也都是相同的。

#### 自定义
  
自定义网桥，默认与bridge网络一样。
