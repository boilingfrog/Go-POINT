<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [docker镜像打包](#docker%E9%95%9C%E5%83%8F%E6%89%93%E5%8C%85)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [简单栗子](#%E7%AE%80%E5%8D%95%E6%A0%97%E5%AD%90)
    - [1、制作dockerfile](#1%E5%88%B6%E4%BD%9Cdockerfile)
    - [2、打包镜像](#2%E6%89%93%E5%8C%85%E9%95%9C%E5%83%8F)
    - [3、镜像打上tag](#3%E9%95%9C%E5%83%8F%E6%89%93%E4%B8%8Atag)
    - [4、上传到仓库](#4%E4%B8%8A%E4%BC%A0%E5%88%B0%E4%BB%93%E5%BA%93)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## docker镜像打包

### 前言

docker打包镜像

### 简单栗子

使用nginx打包一个静态页面的镜像  

测试的代码地址[https://github.com/boilingfrog/daily-test/tree/master/docker-images/test]  

#### 1、制作dockerfile

```
FROM nginx

COPY test /usr/share/nginx/html
```
#### 2、打包镜像

````
docker build -t test-static ./test
````
结果
````
$ docker build -t test-static ./test
Sending build context to Docker daemon  9.728kB
Step 1/2 : FROM nginx
 ---> 4bb46517cac3
Step 2/2 : COPY test /usr/share/nginx/html
 ---> fc981d4aa54c
Successfully built fc981d4aa54c
Successfully tagged test-static:latest
````

` test-static`表示打包成的镜像名，`./test`打包镜像代码地址

#### 3、镜像打上tag  

如果我要上传的地址是`https://hub.docker.com/`,把tag打成你在`https://hub.docker.com/`注册的用户名加镜像的命名就好了  

````
 docker tag test-static:latest liz2019/test-static
````

当然后面也是可以加上版本，如果不加就是默认的`latest`

````
 docker tag test-static:latest liz2019/test-static:v1
````

如果希望上传到自己的搭建的仓库，那么只需加上自己的仓库地址就好了

````
 docker tag test-static:latest hub.xxx.com/xxx/test-static:v1
````

#### 4、上传到仓库

直接push刚打完tag的镜像就好了,上传到`https://hub.docker.com/`是需要登录的。

````
docker push liz2019/test-static
````

结果

```
$ docker push liz2019/test-static
The push refers to repository [docker.io/liz2019/test-static]
525ddb970a89: Pushed 
550333325e31: Mounted from liz2019/docker-file-image 
22ea89b1a816: Mounted from liz2019/docker-file-image 
a4d893caa5c9: Mounted from liz2019/docker-file-image 
0338db614b95: Mounted from liz2019/docker-file-image 
d0f104dc0a1f: Mounted from liz2019/docker-file-image 
latest: digest: sha256:53e8eb1dc6749f05cd303a13588584f9944b6f66b25b8914c49923a16c1ba6b2 size: 1569
```

成功了