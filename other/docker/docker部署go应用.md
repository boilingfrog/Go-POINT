<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [使用docker部署一个go应用](#%E4%BD%BF%E7%94%A8docker%E9%83%A8%E7%BD%B2%E4%B8%80%E4%B8%AAgo%E5%BA%94%E7%94%A8)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [直接部署](#%E7%9B%B4%E6%8E%A5%E9%83%A8%E7%BD%B2)
    - [构建镜像](#%E6%9E%84%E5%BB%BA%E9%95%9C%E5%83%8F)
    - [创建并运行容器](#%E5%88%9B%E5%BB%BA%E5%B9%B6%E8%BF%90%E8%A1%8C%E5%AE%B9%E5%99%A8)
  - [使用docker-compose部署](#%E4%BD%BF%E7%94%A8docker-compose%E9%83%A8%E7%BD%B2)
  - [上传到docker-hub，然后拉取镜像，部署](#%E4%B8%8A%E4%BC%A0%E5%88%B0docker-hub%E7%84%B6%E5%90%8E%E6%8B%89%E5%8F%96%E9%95%9C%E5%83%8F%E9%83%A8%E7%BD%B2)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 使用docker部署一个go应用

### 前言

![](https://img2020.cnblogs.com/blog/1237626/202005/1237626-20200507152147225-1977540183.jpg)

使用docker部署应用已经成为现在的主流。Docker 是一个开源的轻量级容器技术，让开发者可以打包他们的应用以及应用运行的上下文环境到一个可移植的镜像中，然后发布到任何支持Docker的系统上运行。 通过容器技术，在几乎没有性能开销的情况下，Docker 为应用提供了一个隔离运行环境。  

- 简化配置
- 代码流水线管理
- 提高开发效率
- 隔离应用
- 快速、持续部署

### 直接部署
 
首先准备好go项目,使用了一段简单的代码来进行测试

````go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}
func main() {
	http.HandleFunc("/", sayHello) //注册URI路径与相应的处理函数
	log.Println("【默认项目】服务启动成功 监听端口 8000")
	er := http.ListenAndServe("0.0.0.0:8000", nil)
	if er != nil {
		log.Fatal("ListenAndServe: ", er)
	}
}
````

服务器需要配置`go`环境。我的`gopath`是在root下面的。

````
GOPATH="/root/go"
````

然后上传代码到`src`目录中。我的项目名用的`test`。  

在项目根目录创建`Dockerfile`文件然后写入内容。  

````
FROM golang:latest

WORKDIR $GOPATH/src/test
COPY . $GOPATH/src/test
RUN go build .

EXPOSE 8000 
ENTRYPOINT ["./test"]
````

Dockerfile中常见命令的分析，详见[Dockerfile常见的命令](https://github.com/boilingfrog/Go-POINT/blob/master/docker/dockerfile%E4%B8%AD%E5%B8%B8%E7%94%A8%E7%9A%84%E5%91%BD%E4%BB%A4.md)

#### 构建镜像

在项目目录下面执行
````
 docker build -t test-docker-go .
````
我们来看下这条命令
````
Usage:  docker image build [OPTIONS] PATH | URL | -
Options:
-t, --tag list     # 镜像名称
-f, --file string  # 指定Dockerfile文件位置

示例：
docker build . 
docker build -t shykes/myapp .
docker build -t shykes/myapp -f /path/Dockerfile /path
````
执行命令，然后打包镜像
````
# docker build -t test-docker-go .
Sending build context to Docker daemon  14.34kB
Step 1/6 : FROM golang:latest
 ---> 2421885b04da
Step 2/6 : WORKDIR $GOPATH/src/test
 ---> Running in f372c7f2e310
Removing intermediate container f372c7f2e310
 ---> bdedf88480c9
Step 3/6 : COPY . $GOPATH/src/test
 ---> 4e8b7f1a47b9
Step 4/6 : RUN go build .
 ---> Running in 851d5c682f76
Removing intermediate container 851d5c682f76
 ---> 3d5ae3a19f94
Step 5/6 : EXPOSE 8000
 ---> Running in 9ed63b8df046
Removing intermediate container 9ed63b8df046
 ---> 40f1958f50a8
Step 6/6 : ENTRYPOINT ["./test"]
 ---> Running in d505df7ce50c
Removing intermediate container d505df7ce50c
 ---> 7c834b14f69a
Successfully built 7c834b14f69a
Successfully tagged test-docker-go:latest
````

#### 创建并运行容器

执行命令运行并创建容器
````
# docker run -p 8000:8010 test-docker-go
2020/05/09 02:55:43 【默认项目】服务启动成功 监听端口 8010
````
![](https://img2020.cnblogs.com/blog/1237626/202005/1237626-20200509105727786-1309538238.png)


### 使用docker-compose部署

上面成功创建并启动了`go`容器，下面尝试用`docker-composer`，创建并启动容器。

在项目的根目录创建`docker-compose.yml`文件。
````
version: '2'

networks:
  basic:

services:

  test-docker:
    container_name: test-docker1
    image: test-docker-go
    ports:
      - "8020:8000"
    networks:
      - basic
````
然后启动`docker-composer`
````
# docker-compose up
Recreating test-docker1 ... done
Attaching to test-docker1
test-docker1   | 2020/05/09 06:33:47 【默认项目】服务启动成功 监听端口 8010
````

### 上传到docker-hub，然后拉取镜像，部署

下面尝试把镜像上传到`hub.docker`，然后通过拉取镜像，启动容器。  

首先打包镜像到镜像仓库，同理先打包成镜像，为了区分上面的，新打了一个镜像。  
````
# docker build -t test-docker-go-hub .
Sending build context to Docker daemon  14.34kB
Step 1/6 : FROM golang:latest
 ---> 2421885b04da
Step 2/6 : WORKDIR $GOPATH/src/test
 ---> Using cache
 ---> bdedf88480c9
Step 3/6 : COPY . $GOPATH/src/test
 ---> Using cache
 ---> 4e8b7f1a47b9
Step 4/6 : RUN go build .
 ---> Using cache
 ---> 3d5ae3a19f94
Step 5/6 : EXPOSE 8000
 ---> Using cache
 ---> 40f1958f50a8
Step 6/6 : ENTRYPOINT ["./test"]
 ---> Using cache
 ---> 7c834b14f69a
Successfully built 7c834b14f69a
Successfully tagged test-docker-go-hub:latest
````

然后登录`hub.docker`。
````
# docker login
Login with your Docker ID to push and pull images from Docker Hub. If you don't have a Docker ID, head over to https://hub.docker.com to create one.
Username: liz2019
Password: 
````

之后把打包的镜像`push`到仓库中。需要注意的是，需要将镜像打个`tag`，不然`push`会报错
````
denied: requested access to the resource is denied
````

打`tag`
````
# docker tag test-docker-go-hub liz2019/test-docker-go-hub
````

然后`push`
````
# docker push liz2019/test-docker-go-hub
The push refers to repository [docker.io/liz2019/test-docker-go-hub]
12a132dad8d5: Pushed 
16b18b49dbb5: Pushed 
1ffec8d4838f: Pushed 
6e69dbdef94b: Pushed 
f0c38edb3fff: Pushed 
ef234633eec2: Pushed 
8967306e673e: Pushed 
9794a3b3ed45: Pushed 
5f77a51ade6a: Pushed 
e40d297cf5f8: Pushed 
latest: digest: sha256:0ec0fa83015614135357629a433a7d9d19ea57c9f6e5d774772c644509884fa8 size: 2421
````

然后创新创建`docker-compose.yml`文件。
````
version: '3'

networks:
  basic:

services:

  test-docker:
    container_name: test-docker2
    image: liz2019/test-docker-go-hub
    ports:
      - "8020:8000"
    networks:
      - basic
````
然后启动
````
# docker-compose up
Creating network "go_basic" with the default driver
Creating test-docker2 ... done
Attaching to test-docker2
test-docker2   | 2020/05/09 09:03:15 【默认项目】服务启动成功 监听端口 8000
^CGracefully stopping... (press Ctrl+C again to force)
Stopping test-docker2 ... done
````


### 参考
【Gin实践 连载九 将Golang应用部署到Docker】https://segmentfault.com/a/1190000013960558   
【Docker三剑客——Compose】https://blog.csdn.net/Anumbrella/article/details/80877643  