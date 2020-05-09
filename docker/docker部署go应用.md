## 使用docker部署一个go应用

### 前言

![](https://img2020.cnblogs.com/blog/1237626/202005/1237626-20200507152147225-1977540183.jpg)

使用docker部署应用已经成为现在的主流。Docker 是一个开源的轻量级容器技术，让开发者可以打包他们的应用以及应用运行的上下文环境到一个可移植的镜像中，然后发布到任何支持Docker的系统上运行。 通过容器技术，在几乎没有性能开销的情况下，Docker 为应用提供了一个隔离运行环境。  

- 简化配置
- 代码流水线管理
- 提高开发效率
- 隔离应用
- 快速、持续部署

### 部署
 
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
	log.Println("【默认项目】服务启动成功 监听端口 8010")
	er := http.ListenAndServe("0.0.0.0:8010", nil)
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
我们来分析下`Dockerfile`中的几个命令。  

#### FROM 

构建的新镜像是基于哪个镜像   
例如：  
````
FROM centos:6 
````

#### MAINTAINER

镜像维护者的信息  
例如：  
````
MAINTAINER liz
````

#### RUN

构建镜像时运行的`shell`命令  
例如：  
````
RUN [“yum”, “install”, “httpd”]  
RUN yum install httpd  
````

#### CMD 

运行容器时执行的Shell命令  
例如：  
````
CMD [“-c”, “/start.sh”]  
CMD ["/usr/sbin/sshd", "-D"]  
CMD /usr/sbin/sshd –D  
````

#### EXPOSE

声明容器运行的服务端口  
例如：  
````
EXPOSE 80
````

#### ENV

设置容器内环境变量  
例如：  
````
ENV JAVA_HOME /usr/local/jdk1.8.0_45
````

#### ADD

拷贝文件或目录到镜像，如果是URL或压缩包会自动下载或自动解压  
例如：  
````
ADD <src>… <dest>
ADD [“<src>”,… “<dest>”]
ADD https://xxx.com/html.tar.gz /var/www/html
ADD html.tar.gz /var/www/html
```` 

#### COPY 

拷贝文件或目录到镜像，用法同上  
例如：  
````
COPY ./start.sh /start.sh 
````

#### ENTRYPOINT

运行容器时执行的`Shell`命令  
例如：  
````
ENTRYPOINT [“/bin/bash", “-c", “/start.sh"]
ENTRYPOINT /bin/bash -c ‘/start.sh’
````

#### VOLUME

指定容器挂载点到宿主机自动生成的目录或其他容器  
例如:  
````
VOLUME ["/var/lib/mysql"]
````

#### USER

为RUN、CMD和ENTRYPOINT执行命令指定运行用户  
USER <user>[:<group>] or USER <UID>[:<GID>]  
例如:
````
USER liz
````

#### WORKDIR

为RUN、CMD、ENTRYPOINT、COPY和ADD设置工作目录  
例如：
````
WORKDIR /data
````

#### HEALTHCHECK

健康检查
````
HEALTHCHECK --interval=5m --timeout=3s --retries=3 \
  CMD curl -f http://localhost/ || exit 1
````

#### ARG

在构建镜像是指定的一些参数  
例如：
````
FROM centos:6
ARG user     # ARG user=root
USER $user
# docker build --build-arg user=liz Dockerfile .
````

### 构建镜像

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

### 创建并运行容器

执行命令运行并创建容器
````
# docker run -p 8000:8010 test-docker-go
2020/05/09 02:55:43 【默认项目】服务启动成功 监听端口 8010
````
![](https://img2020.cnblogs.com/blog/1237626/202005/1237626-20200509105727786-1309538238.png)

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
      - "8020:8010"
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









### 参考
【Gin实践 连载九 将Golang应用部署到Docker】https://segmentfault.com/a/1190000013960558   
【Docker三剑客——Compose】https://blog.csdn.net/Anumbrella/article/details/80877643  