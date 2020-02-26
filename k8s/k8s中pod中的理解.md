## k8s中pod的理解  

- [基本概念](#%e5%9f%ba%e6%9c%ac%e6%a6%82%e5%bf%b5)
- [pod存在的意义](#pod%e5%ad%98%e5%9c%a8%e7%9a%84%e6%84%8f%e4%b9%89)
- [实现机制](#%e5%ae%9e%e7%8e%b0%e6%9c%ba%e5%88%b6)
    - [共享存储](#%e5%85%b1%e4%ba%ab%e5%ad%98%e5%82%a8)
    - [共享网络](#%e5%85%b1%e4%ba%ab%e7%bd%91%e7%bb%9c)

#### 基本概念

Pod 是 Kubernetes 集群中能够被创建和管理的最小部署单元,它是虚拟存在的。pod是一组容器的集合，pod里面的容器、
共享网络和存储空间，pod是短暂的。  
关键词：   
- 1、最小的部署单元  
- 2、一组容器的集合  
- 3、一个pod里面的容器共享网络和存储空间  
- 4、pod是短暂的  

#### pod存在的意义

为亲密应用而存在   

亲密性应用场景   
1、两个应用之间发生文件交换  
2、两个应用需要通过127.0.0.1或者socket通信  
3、两个应用需要发生频繁的调用  

#### 实现机制

共享网络  
共享存储  

##### 共享存储
在同一个pod中的多个容器能够共享pod级别的存储卷Volume.Volume可以被多个容器进行挂载的操作。  
为什么要共享存储呢?  
pod的生命周期短暂的，随时可能被删除和重启，当一个pod被删除了，又启动一个pod，共享公共的存储卷，以至于信息不会丢失。
![Aaron Swartz](https://github.com/zhan-liz/Go-POINT/blob/master/img/pod_1.png?raw=true)

##### 共享网络


#### pod容器分类与设计模式

Infrastructure Container：基础容器  
维护整个Pod网络空间  
InitContainers：初始化容器   
 先于业务容器开始执行  
Containers：业务容器  
 并行  