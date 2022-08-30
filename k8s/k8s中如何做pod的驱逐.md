<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [对 k8S 中的 node 进行迁移](#%E5%AF%B9-k8s-%E4%B8%AD%E7%9A%84-node-%E8%BF%9B%E8%A1%8C%E8%BF%81%E7%A7%BB)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [处理思路](#%E5%A4%84%E7%90%86%E6%80%9D%E8%B7%AF)
    - [1、设置节点不可调度](#1%E8%AE%BE%E7%BD%AE%E8%8A%82%E7%82%B9%E4%B8%8D%E5%8F%AF%E8%B0%83%E5%BA%A6)
    - [2、驱逐节点上的 pod](#2%E9%A9%B1%E9%80%90%E8%8A%82%E7%82%B9%E4%B8%8A%E7%9A%84-pod)
    - [3、node 迁移结束，设置 node 为可调度状态](#3node-%E8%BF%81%E7%A7%BB%E7%BB%93%E6%9D%9F%E8%AE%BE%E7%BD%AE-node-%E4%B8%BA%E5%8F%AF%E8%B0%83%E5%BA%A6%E7%8A%B6%E6%80%81)
    - [4、pod 回迁](#4pod-%E5%9B%9E%E8%BF%81)
    - [驱逐中的 pod 保护](#%E9%A9%B1%E9%80%90%E4%B8%AD%E7%9A%84-pod-%E4%BF%9D%E6%8A%A4)
      - [什么是 PDB](#%E4%BB%80%E4%B9%88%E6%98%AF-pdb)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 对 k8S 中的 node 进行迁移

### 前言

k8s 机房中的机器需要更换机房，因为用到了多台物理机，所以处理思路，分批次对这几台物理机进行迁移操作。  

每台 node 机器中，都有 pod 部署，所以先对 node 中的 pod 进行驱逐，然后迁移 node 对应的物理机，迁移完成，使 node 恢复 pod 的调度。    

依次完成所有机器的迁移。    

### 处理思路  

#### 1、设置节点不可调度     

命令  

```
$ kubectl cordon node1
```

设置完不可调度之后，对应的 node 状态已经变成了 `SchedulingDisabled`。  

`SchedulingDisabled` 状态表示后续创建的 pod,不会发布到当前的 node 节点。已经在该节点上创建的 pod ,还是能够被正常调度。     

```
$ kubectl get node
NAME              STATUS                     ROLES    AGE    VERSION
node1   Ready,SchedulingDisabled   <none>   468d   v1.19.9
node2   Ready                      <none>   468d   v1.19.9
master  Ready                      master   468d   v1.19.9
```

#### 2、驱逐节点上的 pod 

命令

```
$ kubectl drain node1 --delete-local-data --ignore-daemonsets --force
```

执行之后对应 node 上的 pod 就会被驱逐出当前 node，然后部署到其他的 node 中。  

```
$ kubectl drain node1 --delete-local-data --ignore-daemonsets --force
node/node1 already cordoned
WARNING: ignoring DaemonSet-managed Pods: kube-system/calico-node-v27bp, kube-system/kube-proxy-fvfbk, kube-system/nodelocaldns-qbwvp

evicting pod test/demo3-main-test-6897ff9496-mhm5n
evicting pod default/secret-env-pod
evicting pod kubernetes-dashboard/dashboard-metrics-scraper-79c5968bdc-hg56n
evicting pod test/demo-main-test-cff4cd8f-c2pf4
evicting pod kube-system/coredns-7677f9bb54-2shn6
pod/secret-env-pod evicted
pod/coredns-7677f9bb54-2shn6 evicted
pod/dashboard-metrics-scraper-79c5968bdc-hg56n evicted
pod/demo-main-test-cff4cd8f-c2pf4 evicted
pod/demo3-main-test-6897ff9496-mhm5n evicted
node/node1 evicted
```

状态 evicted 就表示被驱赶完成。   

参数说明：  
  
- delete-local-data ：即使 pod 使用了 emptyDir 也删除；  

- ignore-daemonsets ：忽略 deamonset 控制器的 pod，如果不忽略，deamonset 控制器控制的 pod 被删除后可能马上又在此节点上启动起来,会成为死循环；  

- force ：不加 force 参数只会删除该 node 上由 `ReplicationController, ReplicaSet,DaemonSet,StatefulSet or Job` 创建的 Pod，加了后还会删除’裸奔的 pod’(没有绑定到任何replication controller)。  

#### 3、node 迁移结束，设置 node 为可调度状态  

命令  

```
$ kubectl uncordon node1
```

查看 node 状态已经正常了   

```
$ kubectl get node
NAME              STATUS   ROLES    AGE    VERSION
node1   Ready    <none>   468d   v1.19.9
node2   Ready    <none>   468d   v1.19.9
master  Ready    master   468d   v1.19.9
```

#### 4、pod 回迁

重新部署发布即可。   

#### 驱逐中的 pod 保护

有两种方法  

1、为了防止 驱逐节点上的 pod，可用 pod 数量过少，使用 pdb 来保护； 

2、也可以，执行完成步骤 1，设置节点不可用之后，重新部署 pod,让 pod 部署到其他的 node 中，来实现 node 中 pod 的驱逐。  

##### 什么是 PDB

对于 K8s 集群，我们可能会有下面的操作  

1、排空（drain）节点进行修复或升级；  

2、从集群中排空节点以缩小集群；  

3、从节点中移除一个 pod，以允许其他 pod 使用该节点。   

当发生上面的这几种操作的时候，当前存活的 pod 的数量可能很少或者为0，我们的业务就可能发生中断，所以我们需要控制宕机 pod 的数量，保证有足够的 pod 保障业务的正常运行。    

PDB（PodDisruptionBudget）意为主动驱逐保护，作用就是当 pod 发生自愿中断，也就是上面的三种情况的时候，保证在节点维护时有不低于一定数量的 pod 运行，从而保持服务质量。    

### 参考

【安全地清空一个节点】https://kubernetes.io/zh-cn/docs/tasks/administer-cluster/safely-drain-node/  
【k8s之Pod驱逐迁移和Node节点维护】https://blog.csdn.net/weixin_44729138/article/details/112603786  