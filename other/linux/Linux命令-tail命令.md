<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Linux命令-tail](#linux%E5%91%BD%E4%BB%A4-tail)
  - [命令分析](#%E5%91%BD%E4%BB%A4%E5%88%86%E6%9E%90)
    - [命令格式](#%E5%91%BD%E4%BB%A4%E6%A0%BC%E5%BC%8F)
    - [参数](#%E5%8F%82%E6%95%B0)
    - [例子](#%E4%BE%8B%E5%AD%90)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

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
# tail -f /var/log/jenkins/jenkins.log
2020-06-08 00:34:16.777+0000 [id=27]	INFO	o.s.c.s.AbstractApplicationContext#obtainFreshBeanFactory: Bean factory for application context [org.springframework.web.context.support.StaticWebApplicationContext@67401a8e]: org.springframework.beans.factory.support.DefaultListableBeanFactory@4eb5cf66
2020-06-08 00:34:16.779+0000 [id=27]	INFO	o.s.b.f.s.DefaultListableBeanFactory#preInstantiateSingletons: Pre-instantiating singletons in org.springframework.beans.factory.support.DefaultListableBeanFactory@4eb5cf66: defining beans [filter,legacy]; root of factory hierarchy
2020-06-08 00:34:17.027+0000 [id=27]	INFO	jenkins.InitReactorRunner$1#onAttained: Completed initialization
2020-06-08 00:34:17.111+0000 [id=20]	INFO	hudson.WebAppMain$3#run: Jenkins is fully up and running
2020-06-08 00:34:51.566+0000 [id=41]	INFO	h.m.DownloadService$Downloadable#load: Obtained the updated data file for hudson.tasks.Maven.MavenInstaller
2020-06-08 00:34:52.639+0000 [id=41]	INFO	h.m.DownloadService$Downloadable#load: Obtained the updated data file for hudson.plugins.gradle.GradleInstaller
2020-06-08 00:34:53.453+0000 [id=41]	INFO	h.m.DownloadService$Downloadable#load: Obtained the updated data file for hudson.tasks.Ant.AntInstaller
2020-06-08 00:34:55.783+0000 [id=41]	INFO	h.m.DownloadService$Downloadable#load: Obtained the updated data file for hudson.tools.JDKInstaller
2020-06-08 00:34:55.783+0000 [id=41]	INFO	hudson.util.Retrier#start: Performed the action check updates server successfully at the attempt #1
2020-06-08 00:34:55.787+0000 [id=41]	INFO	hudson.model.AsyncPeriodicWork#lambda$doRun$0: Finished Download metadata. 42,704 ms
````