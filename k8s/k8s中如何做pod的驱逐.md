<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [对 k8S 中的 NODE 进行迁移](#%E5%AF%B9-k8s-%E4%B8%AD%E7%9A%84-node-%E8%BF%9B%E8%A1%8C%E8%BF%81%E7%A7%BB)
  - [处理思路](#%E5%A4%84%E7%90%86%E6%80%9D%E8%B7%AF)
    - [1、设置节点不可调度](#1%E8%AE%BE%E7%BD%AE%E8%8A%82%E7%82%B9%E4%B8%8D%E5%8F%AF%E8%B0%83%E5%BA%A6)
    - [2、驱逐节点上的POD](#2%E9%A9%B1%E9%80%90%E8%8A%82%E7%82%B9%E4%B8%8A%E7%9A%84pod)
    - [3、node 迁移结束，设置 node 为可调度状态](#3node-%E8%BF%81%E7%A7%BB%E7%BB%93%E6%9D%9F%E8%AE%BE%E7%BD%AE-node-%E4%B8%BA%E5%8F%AF%E8%B0%83%E5%BA%A6%E7%8A%B6%E6%80%81)
    - [4、pod 回迁](#4pod-%E5%9B%9E%E8%BF%81)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 对 k8S 中的 NODE 进行迁移

### 处理思路  

#### 1、设置节点不可调度     

命令  

```
$ kubectl cordon node1
```

设置完不可调度发现，对应的 node 状态已经变成了 `SchedulingDisabled`，该状态表示后续创建的 pod,就不会向改节点调度。已经在该节点上创建的 pod ,还是能够被正常调度。     

```
$ kubectl get node
NAME              STATUS                     ROLES    AGE    VERSION
node1   Ready,SchedulingDisabled   <none>   468d   v1.19.9
node2   Ready                      <none>   468d   v1.19.9
master  Ready                      master   468d   v1.19.9
```

#### 2、驱逐节点上的POD  

命令

```
$ kubectl drain node1 --delete-local-data --ignore-daemonsets --force
```

执行之后对应 node 上的 pod 就从改 node 中被驱赶，然后部署到其他的 Node 中。  

```go
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
  
- delete-local-data ：即使pod使用了emptyDir也删除；  

- ignore-daemonsets ：忽略deamonset控制器的pod，如果不忽略，deamonset控制器控制的pod被删除后可能马上又在此节点上启动起来,会成为死循环；  

- force ：不加force参数只会删除该NODE上由ReplicationController, ReplicaSet,DaemonSet,StatefulSet or Job创建的Pod，加了后还会删除’裸奔的pod’(没有绑定到任何replication controller)。  

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

### 参考

【安全地清空一个节点】https://kubernetes.io/zh-cn/docs/tasks/administer-cluster/safely-drain-node/  
【k8s之Pod驱逐迁移和Node节点维护】https://blog.csdn.net/weixin_44729138/article/details/112603786  