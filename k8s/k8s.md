## k8s学习

### 学习笔记记录  
   
安全：集群的认证，鉴全，访问控制，原理及流程  

高可用的节点数最好　用奇数>=3  

APISERVER:所有服务访问统一入口  
CrontrollerManager:维持副本期望数目  
Scheduler：负责介绍任务，选择合适的节点进行分配任务。  
Etcd:键值对数据库　　　存储k8s集群所有重要信息（持久化）  
Kubelet:直接跟容器引擎交互实现容器的生命周期管理  
Kube-proxy：负责写入规则至　 IPTABLES,IPVS实现服务映射访问  
COREDNS:可以为集群中的svc创建一个域名ip对应的关系解析  
DASHBOARD:给k8s集群提供一个b/s结构访问的体系  
INGRESS　controller：官方只能实现四层代理，INGRESS可以实现七层代理  
FEDETATION：可以提供一个夸集群中心多k8s统一挂你功能  
PROMETHEUS:提供一个k8s集群的监控能力   
ELK：提供k8s几集群日志统一分析接入平台  

#### pod  

里面的端口号不能重复  


DaemonSet 确保全部（或者一些）node上运行一个pod的副本。当有node加入到集群时，也会为他
们新增一个pod。当有node从集群移除时，这些pod也会被回收。删除DaemonSet　将会删除它创建的
所有pod。


使用DaemonSet的一些典型用法  

- 运行集群存储daemon，例如在每个node上运行　glusterd,ceph  
- 在每个node上运行日志收集danmon,例如　fluented,logstash  
- 在每个node上运行监控　daemon 例如　promenthens node exporter  

job 负责批处理任务，即仅执行一次的任务，它保证处理任务的一个或者多个pod成功结束  

Cron Job　管理基于时间的job，即:     
- 在给定的时间只运行一次  
- 周期性的给定时间运行  


#### 网络的通讯方式　

k8s的网络模型假定了所有的pod都是一个可以直接连通的扁平的网络空间中，这在gce(google 
computer engine)里面是现成的网络模型，k8s假定这个网络已经存在。而在私有云搭建k8s集群
，就不能假定这个网络已经存在了。我们需要自己实现这个网络假设，将不同节点上的docker容器
之间的互相访问先打通，然后运行k8s

同一个pod内的多个容器之间：lo
 









   

首先安装虚拟机  
可参考安装教程：https://www.jianshu.com/p/18207167b1e7
虚拟机网络的理解：https://www.cnblogs.com/ricklz/p/12216715.html  
virtual box中安装centos，当安装成功之后选择重启，发现又重新安装了，原因是安装成功之后
没有移除虚拟盘  