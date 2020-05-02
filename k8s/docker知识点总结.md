## docker知识点总结

### docker

### 镜像 

什么是镜像？   
简单说，Docker镜像是一个不包含Linux内核而又精简的Linux操作系统。  

镜像从哪里来？  
Docker Hub是由Docker公司负责维护的公共注册中心，包含大量的容器镜像，Docker工具默认从这个公共镜像库下载镜像。  
https://hub.docker.com/explore  
默认是国外的源，下载会慢，建议配置国内镜像仓库：   
````
# vi /etc/docker/daemon.json 
{
  "registry-mirrors": [ "https://registry.docker-cn.com"]
}
````
  

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

### Dockerfile

#### Dockerfile指令
  

