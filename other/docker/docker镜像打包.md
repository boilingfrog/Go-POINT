## docker镜像打包

### 前言

docker打包镜像

### 简单栗子

使用nginx打包一个静态页面的镜像  

1、制作`dockerfile`  

```
FROM nginx

COPY test /usr/share/nginx/html
```
2、打包镜像

````
docker build -t test-static ./test
````
` test-static`表示打包成的镜像名，`./test`打包镜像代码地址