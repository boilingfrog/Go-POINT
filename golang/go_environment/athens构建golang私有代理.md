<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [go中私有代理搭建](#go%E4%B8%AD%E7%A7%81%E6%9C%89%E4%BB%A3%E7%90%86%E6%90%AD%E5%BB%BA)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [为什么选择 athens](#%E4%B8%BA%E4%BB%80%E4%B9%88%E9%80%89%E6%8B%A9-athens)
  - [使用 docker-compose 部署](#%E4%BD%BF%E7%94%A8-docker-compose-%E9%83%A8%E7%BD%B2)
    - [配置私有仓库的认证信息](#%E9%85%8D%E7%BD%AE%E7%A7%81%E6%9C%89%E4%BB%93%E5%BA%93%E7%9A%84%E8%AE%A4%E8%AF%81%E4%BF%A1%E6%81%AF)
    - [部署到海外机器中](#%E9%83%A8%E7%BD%B2%E5%88%B0%E6%B5%B7%E5%A4%96%E6%9C%BA%E5%99%A8%E4%B8%AD)
    - [国内机器的部署](#%E5%9B%BD%E5%86%85%E6%9C%BA%E5%99%A8%E7%9A%84%E9%83%A8%E7%BD%B2)
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

Athens 支持 disk, mongo, gcs, s3, minio, 外部存储/自定义，不过一般建议使用 disk。  

### 使用 docker-compose 部署

官方网站已经，提供了通过 docker 和 二进制部署的方案，这里秉着好记性不如烂笔头的原则，这里自己也做了记录  

#### 配置私有仓库的认证信息  

通过 `.netrc` 文件来配置，里面可以放自己的私有仓库的地址，以及用户，密码认证信息  

```
# cat .netrc
machine gitlab.test.com login test-name password test-pass
```

有几个私有仓库，配置几个就可以了  

#### 部署到海外机器中  

如果有一个海外的服务器，那么部署私有代理仓库就简单了，直接部署在上面就可以了  

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
        - /data/athens-storage:/var/lib/athens
        - ./filter_file:/root/filter_file
    environment:
        - ATHENS_NETRC_PATH=/root/.netrc
        - ATHENS_GONOSUM_PATTERNS=gitlab.test.com
        - ATHENS_STORAGE_TYPE=disk
        - ATHENS_FILTER_FILE=/root/filter_file
        - ATHENS_GOGET_WORKERS=100
        - ATHENS_DISK_STORAGE_ROOT=/var/lib/athens
```

ATHENS_GONOSUM_PATTERNS：配置为私库地址, 作用避免私库地址流入公网，支持通配符，多个可以使用`,`分割。  

通过 ATHENS_FILTER_FILE 配置访问的策略  

- `-` 表示禁止下载此软件包,来屏蔽一些有安全隐患的包，若请求，报403；      

```
# cat filter_file

- github.com/gogo
```

栗如：配置了`- github.com/gogo`    

`go: github.com/gogo/googleapis@v1.2.0: reading http://127.0.0.1:3000/github.com/gogo/googleapis/@v/v1.2.0.mod: 403 Forbidden
`

启动 `docker-compose up -d`    

客户端设置代理 `export GOPROXY=http://xxxx:3000`  

这样就能使用我们的代理服务了    

因为选择的 ATHENS_STORAGE_TYPE 为 disk，athens 服务会在拉取资源包的同时，也会下载资源包到配置的 ATHENS_DISK_STORAGE_ROOT 中。  

#### 国内机器的部署

有时候我们可能没有一台可以访问外网的机器，这时候 athens 提供了 ATHENS_GLOBAL_ENDPOINT ，可配置一些国内的公共代理。  

```yaml
version: '2'
services:
  athens:
    image: gomods/athens:v0.11.0
    restart: always
    container_name: athens_proxy
    ports:
      - "3000:3000"
    volumes:
      - ./.netrc:/root/.netrc
      - ./athens-storage:/var/lib/athens
      - ./filter_file:/root/filter_file
    environment:
      - ATHENS_GLOBAL_ENDPOINT=https://goproxy.cn
      - ATHENS_NETRC_PATH=/root/.netrc
      - ATHENS_GONOSUM_PATTERNS=gitlab.test.com
      - ATHENS_STORAGE_TYPE=disk
      - ATHENS_DISK_STORAGE_ROOT=/var/lib/athens
      - ATHENS_FILTER_FILE=/root/filter_file
      - ATHENS_GOGET_WORKERS=100
```

如果配置了 ATHENS_GLOBAL_ENDPOINT，需要配置过滤策略  

athens 可配置软件包的过滤策略，来决定那些包可以存放到本地  

athens 中用 `D、-、+` 这三种方式来确定认证策略  

- `D` 需要放在第一行，如果没放,配置的 ATHENS_GLOBAL_ENDPOINT 不生效；    

- `-` 表示禁止下载此软件包,来屏蔽一些有安全隐患的包，若请求，报403；      

- `+` 表示不需要通过 GlobalEndpoint 代理，而是直接访问的资源包，通过 GlobalEndpoint 访问的资源不会下载到本地，直接访问的包会下载到本地，我们私有仓库的需要配置到这里面，因为通过公共代理是找不到的；   

`-` 与 `+` 对软件包的策略可指定至版本，多个版本用,号分隔，甚至可使用版本修饰符`(~, ^, >)`。此外 # 开头的行表示注释，会被忽略。    

```
# cat filter_file
D
# 表示不需要通过 GlobalEndpoint 代理，而是直接访问的资源包； 
+ gitlab.test.com
```

修饰符  

`~`：`~1.2.3`表示激活所有大于等于3的 patch 版本，在语义化版本方案中，最后一位的3表示补丁版本。如 `1.2.3, 1.2.4, 1.2.5`；  

`^`：`^1.2.3`表示激活所有大于等于2的 minor 与大于等于3的 patch 版本。如 `1.2.3, 1.3.0`；  

`<`：`<1.2.3`表示激活所有小于 `1.2.3` 的版本。如 `1.2.2, 1.0.0, 0.1.1`。  

### 参考

【介绍 ATHENS】https://gomods.io/zh/intro/    