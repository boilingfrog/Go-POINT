<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [go中私有代理搭建](#go%E4%B8%AD%E7%A7%81%E6%9C%89%E4%BB%A3%E7%90%86%E6%90%AD%E5%BB%BA)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [为什么选择 athens](#%E4%B8%BA%E4%BB%80%E4%B9%88%E9%80%89%E6%8B%A9-athens)
  - [使用 docker-compose 部署](#%E4%BD%BF%E7%94%A8-docker-compose-%E9%83%A8%E7%BD%B2)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## go中私有代理搭建

### 前言

最近公司的代理出现问题了，刚好借这个机会来学习下，athens 如何构建私有代理  

### 为什么选择 athens

私有化代理的选取标准无非就是下面的几点  

1、托管私有模块；  

2、排除对公有模块的访问；   

3、存储公有模块；  

**athens 的特点：**  

Athens 首先可以配置访问私有仓库；  

Athens 的会存储每次拉取的包，如果该模块之前没有通过 athens，athens 会向目标地址请求数据，在返回给客户端的时候，会存储该模块到本地的磁盘中，这样实现了 `go mod download`永远只会发生一次；  

Athens 处理存储的策略为仅追加，一个模块被保存，它就永远不会改变，即使开发人员对 tag 进行了强推，那么也不会被删除；  

Athens 也可以 filter 文件配置，过滤一些有安全隐患的包。  

Athens 支持 disk, mongo, gcs, s3, minio, 外部存储/自定义，不过一半建议使用 disk。  

### 使用 docker-compose 部署

官方网站已经，提供了通过 docker 和 二进制部署的方案，这里秉着好记性不如烂笔头的原则，这里自己也做了记录  

#### 1、配置私有仓库的认证信息  

通过 .netrc 文件来配置，里面可以放自己的私有仓库的地址，以及用户，密码认证信息  

```
# cat .netrc
machine gitlab.test.com login test-name password test-pass
```

有几个私有仓库，配置几个就可以了  

#### 2、配置过滤策略  

athens 可配置软件包的过滤策略,来决定那些包可以存放到本地  

athens 中用 `D、-、+` 这三种方式来确定认证策略  

- `D` 需要放在第一行，如果没放，配置的过滤策略就不生效；    

- `-` 表示禁止下载此软件包,来屏蔽一些有安全隐患的包，若请求，报403；      

栗如：配置了`- github.com/gogo`    

`go: github.com/gogo/googleapis@v1.2.0: reading http://127.0.0.1:3000/github.com/gogo/googleapis/@v/v1.2.0.mod: 403 Forbidden
`
- `+` 表示这些仓库的包是否可以缓存在本地，如果配置了仓库地址，第一次请求会从源端拉取包，然后就会缓存在本地，之后就不需要请求源数据了，从本地直接获取；      

`-` 与 `+` 对软件包的策略可指定至版本，多个版本用,号分隔，甚至可使用版本修饰符`(~, ^, >)`。此外 # 开头的行表示注释，会被忽略。    

```
# cat filter_file
D
# 以下不需要通过GlobalEndpoint下载
+ gitlab.test.com
```

#### docker-compose 文件  

```yaml
version: '2'
services:
  athens:
    image: gomods/athens:v0.10.0
    restart: always
    container_name: athens_proxy
    ports:
        - "3000:3000"
    volumes:
        - /data/athens/.netrc:/root/.netrc
        - /data/athens/filter_file:/root/filter_file
        - /data/athens-storage:/var/lib/athens
    environment:
        - ATHENS_GLOBAL_ENDPOINT=https://goproxy.cn
        - ATHENS_FILTER_FILE=/root/filter_file
        - ATHENS_NETRC_PATH=/root/.netrc
        - ATHENS_GONOSUM_PATTERNS=gitlab.test.com
        - ATHENS_STORAGE_TYPE=disk
        - ATHENS_DISK_STORAGE_ROOT=/var/lib/athens
```

不过虽然我们部署了athens，墙外的资源还是不能访问，可以通过以下两种方式解决：   

1、我们可以将 athens 部署在海外的机器中；  

2、可以设置一个 ATHENS_GLOBAL_ENDPOINT，使用公共的代理，来完成对公共资源的访问；    

介绍两个重要的配置  

ATHENS_GLOBAL_ENDPOINT：代理地址，配置成 `https://goproxy.cn` ,这样墙外的包可以通过这个去拉取。  

ATHENS_GONOSUM_PATTERNS：配置为私库地址, 作用避免私库地址流入公网，支持通配符，多个可以使用`,`分割。  

### 使用

启动 `docker-compose up -d`    

客户端设置代理 `export GOPROXY=http://127.0.0.1:3000`  

这样就能使用我们的代理服务了  





### 参考

【介绍 ATHENS】https://gomods.io/zh/intro/    