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


   

首先安装虚拟机  
可参考安装教程：https://www.jianshu.com/p/18207167b1e7
虚拟机网络的理解：https://www.cnblogs.com/ricklz/p/12216715.html  
virtual box中安装centos，当安装成功之后选择重启，发现又重新安装了，原因是安装成功之后
没有移除虚拟盘  