<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Dockerfile中常用的命令](#dockerfile%E4%B8%AD%E5%B8%B8%E7%94%A8%E7%9A%84%E5%91%BD%E4%BB%A4)
    - [FROM](#from)
    - [MAINTAINER](#maintainer)
    - [RUN](#run)
    - [CMD](#cmd)
    - [EXPOSE](#expose)
    - [ENV](#env)
    - [ADD](#add)
    - [COPY](#copy)
    - [ENTRYPOINT](#entrypoint)
    - [VOLUME](#volume)
    - [USER](#user)
    - [WORKDIR](#workdir)
    - [HEALTHCHECK](#healthcheck)
    - [ARG](#arg)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Dockerfile中常用的命令

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



### 参考
【Gin实践 连载九 将Golang应用部署到Docker】https://segmentfault.com/a/1190000013960558   
【Docker三剑客——Compose】https://blog.csdn.net/Anumbrella/article/details/80877643  