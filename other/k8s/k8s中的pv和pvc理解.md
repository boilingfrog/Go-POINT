<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [pv和pvc](#pv%E5%92%8Cpvc)
  - [什么是pv和PVC](#%E4%BB%80%E4%B9%88%E6%98%AFpv%E5%92%8Cpvc)
  - [生命周期](#%E7%94%9F%E5%91%BD%E5%91%A8%E6%9C%9F)
  - [PV创建的流程](#pv%E5%88%9B%E5%BB%BA%E7%9A%84%E6%B5%81%E7%A8%8B)
    - [1、创建一个远程块存储，相当于创建了一个磁盘，称为Attach](#1%E5%88%9B%E5%BB%BA%E4%B8%80%E4%B8%AA%E8%BF%9C%E7%A8%8B%E5%9D%97%E5%AD%98%E5%82%A8%E7%9B%B8%E5%BD%93%E4%BA%8E%E5%88%9B%E5%BB%BA%E4%BA%86%E4%B8%80%E4%B8%AA%E7%A3%81%E7%9B%98%E7%A7%B0%E4%B8%BAattach)
    - [2、将这个磁盘设备挂载到宿主机的挂载点，称为Mount](#2%E5%B0%86%E8%BF%99%E4%B8%AA%E7%A3%81%E7%9B%98%E8%AE%BE%E5%A4%87%E6%8C%82%E8%BD%BD%E5%88%B0%E5%AE%BF%E4%B8%BB%E6%9C%BA%E7%9A%84%E6%8C%82%E8%BD%BD%E7%82%B9%E7%A7%B0%E4%B8%BAmount)
    - [3、绑定](#3%E7%BB%91%E5%AE%9A)
  - [持久化卷声明的保护](#%E6%8C%81%E4%B9%85%E5%8C%96%E5%8D%B7%E5%A3%B0%E6%98%8E%E7%9A%84%E4%BF%9D%E6%8A%A4)
  - [PV类型](#pv%E7%B1%BB%E5%9E%8B)
  - [PV卷阶段状态](#pv%E5%8D%B7%E9%98%B6%E6%AE%B5%E7%8A%B6%E6%80%81)
  - [基本的使用](#%E5%9F%BA%E6%9C%AC%E7%9A%84%E4%BD%BF%E7%94%A8)
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

pv和pvc遵循以下生命周期：  

1、 供应准备。通过集群外的存储系统或者云平台来提供存储持久化支持。  

- 静态提供：管理员手动创建多个PV，供PVC使用。

- 动态提供：动态创建PVC特定的PV，并绑定。

2、绑定。用户创建pvc并指定需要的资源和访问模式。在找到可用pv之前，pvc会保持未绑定状态。  

3、使用。用户可在pod中像volume一样使用pvc。  

4、释放。用户删除pvc来回收存储资源，pv将变成“released”状态。由于还保留着之前的数据，这些数据需要根据不同的策略来处理，否则这些存储资源无法被其他pvc使用。  

5、回收(Reclaiming)。pv可以设置三种回收策略：保留（Retain），回收（Recycle）和删除（Delete）。  

- 保留策略：允许人工处理保留的数据。  

- 删除策略：将删除pv和外部关联的存储资源，需要插件支持。  

- 回收策略：将执行清除操作，之后可以被新的pvc使用，需要插件支持。

目前只有NFS和HostPath类型卷支持回收策略，`AWS EBS,GCE PD,Azure Disk`和`Cinder`支持删除`(Delete)`策略。  

### PV创建的流程

大多数情况，持久化`Volume`的实现，依赖于远程存储服务，如远程文件存储（NFS、GlusterFS）、远程块存储（公有云提供的远程磁盘）等。  

K8s需要使用这些存储服务，来为容器准备一个持久化的宿主机目录，以供以后挂载使用，创建这个目录分为两阶段：  

#### 1、创建一个远程块存储，相当于创建了一个磁盘，称为Attach

由`Volume Controller`负责维护，不断地检查 每个Pod对应的PV和所在的宿主机的挂载情况。可以理解为创建了一块NFS磁盘，相当于执行  

```
gcloud compute instances attach-disk < 虚拟机名字 > --disk < 远程磁盘名字 >
```

为了使用这块磁盘，还需要挂载操作  

#### 2、将这个磁盘设备挂载到宿主机的挂载点，称为Mount

将远程磁盘挂载到宿主机上，发生在Pod对应的宿主机上，是kubelet组件一部分，利用goroutine执行  

相当于执行  

```
mount -t nfs <NFS 服务器地址 >:/ /var/lib/kubelet/pods/<Pod 的 ID>/volumes/kubernetes.io~<Volume 类型 >/<Volume 名字 > 
```

通过这个挂载操作，Volume的宿主机目录就成为了一个远程NFS目录的挂载点，以后写入的所有文件，都会被保存在NFS服务器上  

如果是已经有NFS磁盘，第一步可以省略.  

同样，删除PV的时候，也需要Umount和Dettach两个阶段处理   

#### 3、绑定

master中的控制环路监视新的PVC，寻找匹配的PV（如果可能），并将它们绑定在一起。如果为新的PVC动态调配PV，则该环路将始终将该PV绑定到PVC。否则，用户总会得到他们所请求的存储，但是容量可能超出要求的数量，一旦PV和PVC绑定后，Persistent Volume Claim绑定是排它性的，不管它们是如何绑定的，PVC和PV绑定是一对一映射的。  

### 持久化卷声明的保护

PVC保护的目的是确保由Pod正在使用的PVC不会冲系统中移除，因为如果被移除的话可能会导致数据的丢失，当启用PVC保护alpha功能时，如果用户删除了一个Pod正在使用的PVC，则该PVC不会被立即删除，PVC的删除将会被推迟，直到PVC不再被任何Pod使用  

注意：当Pod状态为Pending，并且Pod已经分配给节点或Pod为Running状态时，PVC处于活动状态。  

### PV类型

```
awsElasticBlockStore - AWS 弹性块存储（EBS）
azureDisk - Azure Disk
azureFile - Azure File
cephfs - CephFS volume
cinder - Cinder （OpenStack 块存储） (弃用)
csi - 容器存储接口 (CSI)
fc - Fibre Channel (FC) 存储
flexVolume - FlexVolume
flocker - Flocker 存储
gcePersistentDisk - GCE 持久化盘
glusterfs - Glusterfs 卷
hostPath - HostPath 卷 （仅供单节点测试使用；不适用于多节点集群； 请尝试使用 local 卷作为替代）
iscsi - iSCSI (SCSI over IP) 存储
local - 节点上挂载的本地存储设备
nfs - 网络文件系统 (NFS) 存储
photonPersistentDisk - Photon 控制器持久化盘。 （这个卷类型已经因对应的云提供商被移除而被弃用）。
portworxVolume - Portworx 卷
quobyte - Quobyte 卷
rbd - Rados 块设备 (RBD) 卷
scaleIO - ScaleIO 卷 (弃用)
storageos - StorageOS 卷
vsphereVolume - vSphere VMDK 卷
 ```

### PV卷阶段状态

 - Available: 资源尚未被claim使用
 
 - Bound: 卷已经被绑定到claim了
 
 - Released: claim被删除，卷处于释放状态，但未被集群回收
 
 - Failed: 卷自动回收失败
 
### 基本的使用  

定义NFS PV 资源(静态):  

```yaml
#pv定义如下:
apiVersion: v1
kind: PersistentVolume
metadata:
  name: nfs
spec:
  storageClassName: manual
  capacity:
    storage: 1Gi
  accessModes:
    - ReadWriteMany
  nfs:
    server: 10.244.1.4
    path: "/"
```

定义pvc资源:  

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: nfs
spec:
  accessModes:
    - ReadWriteMany
  storageClassName: manual
  resources:
    requests:
      storage: 1Gi
```

pvc和pv匹配规则：  

- PV 和 PVC 的 spec 字段。比如，PV 的存储（storage）大小，必须满足 PVC的要求。

- PV 和 PVC 的 storageClassName 字段必须一样。


### 参考

【kubernetes系列11—PV和PVC详解】https://www.cnblogs.com/along21/p/10342788.html  
【PV、PVC和StorageClass】https://support.huaweicloud.com/basics-cce/kubernetes_0030.html  
【持久卷】https://kubernetes.io/zh/docs/concepts/storage/persistent-volumes/  
【存储类】https://kubernetes.io/zh/docs/concepts/storage/storage-classes/  
【持久化存储之 PV、PVC、StorageClass】https://www.cnblogs.com/menkeyi/p/10903647.html  
 