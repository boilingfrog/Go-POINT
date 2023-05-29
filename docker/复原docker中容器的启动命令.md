<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [复原 docker 容器的启动命令](#%E5%A4%8D%E5%8E%9F-docker-%E5%AE%B9%E5%99%A8%E7%9A%84%E5%90%AF%E5%8A%A8%E5%91%BD%E4%BB%A4)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [查看 docker 容器的启动命令](#%E6%9F%A5%E7%9C%8B-docker-%E5%AE%B9%E5%99%A8%E7%9A%84%E5%90%AF%E5%8A%A8%E5%91%BD%E4%BB%A4)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 复原 docker 容器的启动命令

### 前言

不规范的操作，在启动 docker 容器，没有留命令脚本，或者没有使用 docker-compose, 这时候再次编辑重启，因为不知道启动的命令，这时候编辑操作就会变的困难了。     

所以如何查看 docker 容器的启动命令呢？  

### 查看 docker 容器的启动命令

使用 `get_command_4_run_container`    

这是一个不知道启动命令的 docker 容器  

```
$ docker ps | grep nginx
 
45d33e955017     nginx     "/docker-entrypoint.…"   2 years ago    Up 3 months    0.0.0.0:111->111/tcp, 0.0.0.0:222->222/tcp, 0.0.0.0:333->333/tcp, 0.0.0.0:444->444/tcp   nginx-doc
```

使用 `get_command_4_run_container` 来获取启动命令  

1、get_command_4_run_container 本身是个 docker 镜像，首先下载镜像；  

```
docker pull cucker/get_command_4_run_container
```

2、通过命令获取容器启动的命令；   

```
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock cucker/get_command_4_run_container [容器名称]/[容器ID]
```

操作下上面的栗子  

```
$ docker run --rm -v /var/run/docker.sock:/var/run/docker.sock cucker/get_command_4_run_container 45d33e955017

docker run -d \
 --name nginx-doc \
 --ipc shareable \
 --log-opt max-file=100 \
 --log-opt max-size=10M \
 -p 111:111/tcp \
 -p 222:222/tcp \
 -p 333:333/tcp \
 -p 444:444/tcp \
 --stop-signal SIGQUIT \
 -v /var/log/nginx:/var/log/nginx \
 -v /data/gleeman-docs:/usr/share/nginx/html/docs:ro \
 -v /data/gleeman-blog/public:/usr/share/nginx/html/blog:ro \
 -v /data/node-monorepo-docs:/usr/share/nginx/html/monorepo:ro \
 -v /data/purchase-doc:/usr/share/purchase:ro \
 -v /data/reborn-doc:/usr/share/reborn:ro \
 -v /data/node-mirrors:/usr/share/mirrors:ro \
 -v /data/google-storage-cache:/var/cache/google-storage:z \
 -v /data/nginx-setup/conf.d:/etc/nginx/conf.d:ro \
 -v /data/nginx-setup/nginx.conf:/etc/nginx/nginx.conf:ro \
 nginx
```

### 参考

【get_command_4_run_container】https://hub.docker.com/r/cucker/get_command_4_run_container