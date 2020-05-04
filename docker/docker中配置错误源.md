## docker学习－配置错误的源


### 问题点剖析

使用docker安装了`nginx`，编写`Dockerfile`,映射端口，终于跑起来了。但是，当我重启服务器，再次查看`docker`容器的状态，发现报错了。

````
# docker ps -a
Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?
````

然后重启docker

````
#  sudo service docker start
Redirecting to /bin/systemctl start docker.service
Job for docker.service failed because start of the service was attempted too often. See "systemctl status docker.service" and "journalctl -xe" for details.
To force a start use "systemctl reset-failed docker.service" followed by "systemctl start docker.service" again.

````

发现还是不行，根据提示查看`docker.service`的`status`

````
# systemctl status docker.service
● docker.service - Docker Application Container Engine
   Loaded: loaded (/usr/lib/systemd/system/docker.service; enabled; vendor preset: disabled)
   Active: failed (Result: start-limit) since 二 2020-05-05 01:13:34 CST; 9s ago
     Docs: https://docs.docker.com
  Process: 1880 ExecStart=/usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock (code=exited, status=1/FAILURE)
 Main PID: 1880 (code=exited, status=1/FAILURE)

5月 05 01:13:32 10.0.2.8 systemd[1]: docker.service: main process exited, code=exited, sta...URE
5月 05 01:13:32 10.0.2.8 systemd[1]: Failed to start Docker Application Container Engine.
5月 05 01:13:32 10.0.2.8 systemd[1]: Unit docker.service entered failed state.
5月 05 01:13:32 10.0.2.8 systemd[1]: docker.service failed.
5月 05 01:13:34 10.0.2.8 systemd[1]: docker.service holdoff time over, scheduling restart.
5月 05 01:13:34 10.0.2.8 systemd[1]: Stopped Docker Application Container Engine.
5月 05 01:13:34 10.0.2.8 systemd[1]: start request repeated too quickly for docker.service
5月 05 01:13:34 10.0.2.8 systemd[1]: Failed to start Docker Application Container Engine.
5月 05 01:13:34 10.0.2.8 systemd[1]: Unit docker.service entered failed state.
5月 05 01:13:34 10.0.2.8 systemd[1]: docker.service failed.
Hint: Some lines were ellipsized, use -l to show in full.
````

分析发现原因是docker不能启动  

````
Failed to start Docker Application Container Engine.
````

`daemon.json`如果包含格式不正确的`JSON`，`Docker`将无法启动  

检查了一下`daemon.json`,果真，少了一个`“`。  

修改`daemon.json`，然后重启解决了。


### 参考

【docker安装完了以后，服务启动不了】http://www.docker.org.cn/thread/72.html