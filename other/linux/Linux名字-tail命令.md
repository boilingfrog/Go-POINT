## Linux命令-tail

### 命令分析

`tail`命令可用于查看文件的内容，通常用来查看日志，加上`-f`参数就可以查看最新的日志并且不断刷新。

#### 命令格式

````
tail [参数] [文件]  
````

#### 参数

- -f 循环读取
- -q 不显示处理信息
- -v 显示详细的处理信息
- -c<数目> 显示的字节数
- -n<行数> 显示文件的尾部 n 行内容
- --pid=PID 与-f合用,表示在进程ID,PID死掉之后结束
- -q, --quiet, --silent 从不输出给出文件名的首部
- -s, --sleep-interval=S 与-f合用,表示在每次反复的间隔休眠S秒

#### 例子

实时查看`jenkins`的日志

````
# tail -f  /var/log/jenkins/jenkins.log
Please use the following password to proceed to installation:

600ef7d47db54ff2bfc518a6c2f2a857

This may also be found at: /var/lib/jenkins/secrets/initialAdminPassword

*************************************************************
*************************************************************
*************************************************************

2020-06-07 14:10:31.384+0000 [id=44]	INFO	h.m.DownloadService$Downloadable#load: Obtained the updated data file for hudson.tasks.Maven.MavenInstaller
2020-06-07 14:10:31.385+0000 [id=44]	INFO	hudson.util.Retrier#start: Performed the action check updates server successfully at the attempt #1
2020-06-07 14:10:31.390+0000 [id=44]	INFO	hudson.model.AsyncPeriodicWork#lambda$doRun$0: Finished Download metadata. 739,411 ms
````