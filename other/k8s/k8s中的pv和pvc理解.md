<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [pv和pvc](#pv%E5%92%8Cpvc)
  - [什么是pv和PVC](#%E4%BB%80%E4%B9%88%E6%98%AFpv%E5%92%8Cpvc)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## pv和pvc

### 什么是pv和PVC

持久卷（PersistentVolume，PV）是集群中由管理员配置的一段网络存储。它是集群中的资源，就像节点是集群资源一样。PV持久卷和普通的Volume一样，也是使用卷插件来实现的，只是它们拥有独立于任何使用PV的Pod的生命周期。此API对象捕获存储实现的详细信息，包括NFS，`iSCSI`或特定于云提供程序的存储系统。  

持久卷申领（PersistentVolumeClaim，PVC）表达的是用户对存储的请求。概念上与Pod类似。Pod会耗用节点资源，而PVC申领会耗用PV资源。Pod可以请求特定数量的资源（CPU 和内存）；同样 PVC 申领也可以请求特定的大小和访问模式 （例如，可以要求PV卷能够以 ReadWriteOnce、ReadOnlyMany 或 ReadWriteMany 模式之一来挂载）。  

<img src="/img/pv_pvc_1.png" alt="pv_pvc" align=center />


虽然`PersistentVolumeClaims`允许用户使用抽象存储资源，但是PersistentVolumes对于不同的问题，用户通常需要具有不同属性（例如性能）。群集管理员需要能够提供各种`PersistentVolumes`不同的方式，而不仅仅是大小和访问模式，而不会让用户了解这些卷的实现方式。对于这些需求，有StorageClass 资源。  

`StorageClass`为管理员提供了一种描述他们提供的存储的“类”的方法。 不同的类可能映射到服务质量级别，或备份策略，或者由群集管理员确定的任意策略。 `Kubernetes`本身对于什么类别代表是不言而喻的。 这个概念有时在其他存储系统中称为“配置文件”。


总结下来就是

- PVC：描述 Pod想要使用的持久化属性，比如存储大小、读写权限等  

- PV：描述一个具体的Volume属性，比如Volume的类型、挂载目录、远程存储服务器地址等  

- StorageClass：充当PV的模板，自动为PVC创建PV

<img src="/img/pv_pvc-2.png" alt="pv_pvc" align=center />

### 生命周期



### 参考

【kubernetes系列11—PV和PVC详解】https://www.cnblogs.com/along21/p/10342788.html  
【PV、PVC和StorageClass】https://support.huaweicloud.com/basics-cce/kubernetes_0030.html  
【持久卷】https://kubernetes.io/zh/docs/concepts/storage/persistent-volumes/  
【存储类】https://kubernetes.io/zh/docs/concepts/storage/storage-classes/  
【持久化存储之 PV、PVC、StorageClass】https://www.cnblogs.com/menkeyi/p/10903647.html  
 