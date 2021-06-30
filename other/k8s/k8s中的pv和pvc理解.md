<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [pv和pvc](#pv%E5%92%8Cpvc)
  - [什么是pv和PVC](#%E4%BB%80%E4%B9%88%E6%98%AFpv%E5%92%8Cpvc)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## pv和pvc

### 什么是pv和PVC

`PersistentVolume（PV）`是集群中由管理员配置的一段网络存储。 它是集群中的资源，就像节点是集群资源一样。 PV是容量插件，如`Volumes`，但其生命周期独立于使用PV的任何单个pod。 此API对象捕获存储实现的详细信息，包括NFS，`iSCSI`或特定于云提供程序的存储系统。  

`PersistentVolumeClaim（PVC）`是由用户进行存储的请求。 它类似于pod。 Pod消耗节点资源，PVC消耗PV资源。Pod可以请求特定级别的资源（CPU和内存）。声明可以请求特定的大小和访问模式（例如，可以一次读/写或多次只读）。  

<img src="/img/pv_pvc_1.png" alt="pv_pvc" align=center />


### 参考

【kubernetes系列11—PV和PVC详解】https://www.cnblogs.com/along21/p/10342788.html  
【PV、PVC和StorageClass】https://support.huaweicloud.com/basics-cce/kubernetes_0030.html  