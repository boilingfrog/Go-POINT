## k8s中pod的理解  

- [基本概念](#%e5%9f%ba%e6%9c%ac%e6%a6%82%e5%bf%b5)
- [pod存在的意义](#pod%e5%ad%98%e5%9c%a8%e7%9a%84%e6%84%8f%e4%b9%89)
- [实现机制](#%e5%ae%9e%e7%8e%b0%e6%9c%ba%e5%88%b6)
    - [共享存储](#%e5%85%b1%e4%ba%ab%e5%ad%98%e5%82%a8)
    - [共享网络](#%e5%85%b1%e4%ba%ab%e7%bd%91%e7%bb%9c)

#### 基本概念

Pod 是 Kubernetes 集群中能够被创建和管理的最小部署单元,它是虚拟存在的。pod是一组容器的集合，并且部署在同一个pod里面
的容器是亲密性很强的一组容器，pod里面的容器，共享网络和存储空间，pod是短暂的。  
关键词：   
- 1、最小的部署单元  
- 2、一组容器的集合 （亲密性强）
- 3、一个pod里面的容器共享网络和存储空间  
- 4、pod是短暂的  

#### pod存在的意义

为亲密应用而存在  
 
亲密性应用场景   
1、两个应用之间发生文件交换  
2、两个应用需要通过127.0.0.1或者socket通信  
3、两个应用需要发生频繁的调用  

#### pod容器分类与设计模式

- Infrastructure Container：基础容器,也叫sandbox或pause容器
维护整个Pod网络空间 
- InitContainers：初始化容器   
 先于业务容器开始执行，让一些有依赖的容器有个先后执行的顺序
- Containers：业务容器  
 并行  
 
 ````
 // SyncPod syncs the running pod into the desired pod by executing following steps:
 //
 //  1. Compute sandbox and container changes.
 //  2. Kill pod sandbox if necessary.
 //  3. Kill any containers that should not be running.
 //  4. Create sandbox if necessary.
 //  5. Create ephemeral containers.
 //  6. Create init containers.
 //  7. Create normal containers.
 func (m *kubeGenericRuntimeManager) SyncPod(pod *v1.Pod, podStatus *kubecontainer.PodStatus, pullSecrets []v1.Secret, backOff *flowcontrol.Backoff) (result kubecontainer.PodSyncResult) {
        ...
 }
 ````
#### 实现机制

共享网络  
共享存储  

##### 共享存储
在同一个pod中的多个容器能够共享pod级别的存储卷Volume.Volume可以被多个容器进行挂载的操作。  
为什么要共享存储呢?  
pod的生命周期短暂的，随时可能被删除和重启，当一个pod被删除了，又启动一个pod，共享公共的存储卷，以至于信息不会丢失。

<img src="../img/pod_1.png" width = "100%" height = "300" alt="图片名称" align=center />

##### 共享网络
同一个pod的多个容器，会被共同分配到同一个Host上共享网络栈。所以pod里面的容器通过localhost就可以通信了。当然这也从侧面
说明了，为什么pod是为了亲密性的应用而生的。   
docker的4中网络模式，其中有一种模式是container模式，它能够让很多容器共享一个网络名称空间， 具体的原理是先使用briage模式启动第一个容器， 之后
启动的其他容器纷纷使用container模式将网络环境绑定到这第一个容器上。这样这些容器的网络就连接到了一起，他们互相可以使用
localhost这种方式进行网络通信。  
  
<img src="../img/pod_2.png" width = "500" height = "300" alt="图片名称" align=center />

如上图所示，这个 Pod 里有两个用户容器 A 和 B，还有一个infra container， 它也叫做pause容器，也被称为sandbox， 意思是沙箱，这个沙箱为其他
容器提供共享的网络和文件挂载资源。 pod在启动的时候Infrastructure Container是第一个启动的容器，也叫做pause容器，也被称为sandbox。之后才启动
InitContainers初始化容器和Containers业务容器。而当这个容器被创建出来并hold住Network Namespace之后，其他由用户自己定义的容器就可以通过
container模式加入到这个容器的Network Namespace中。这也就意味着，对于在一个POD中的容器A和容器B来说，他们拥有相同的IP地址，可以通过
localhost进行互相通信。   

````
// RunPodSandbox creates and starts a pod-level sandbox. Runtimes should ensure
// the sandbox is in ready state.
// For docker, PodSandbox is implemented by a container holding the network
// namespace for the pod.
// Note: docker doesn't use LogDirectory (yet).
func (ds *dockerService) RunPodSandbox(ctx context.Context, r *runtimeapi.RunPodSandboxRequest) (*runtimeapi.RunPodSandboxResponse, error) {
	config := r.GetConfig()

	// Step 1: Pull the image for the sandbox.
	image := defaultSandboxImage
	podSandboxImage := ds.podSandboxImage
	if len(podSandboxImage) != 0 {
		image = podSandboxImage
	}

	// Step 2: Create the sandbox container.
	if r.GetRuntimeHandler() != "" && r.GetRuntimeHandler() != runtimeName {
		return nil, fmt.Errorf("RuntimeHandler %q not supported", r.GetRuntimeHandler())
	}

	// Step 3: Create Sandbox Checkpoint.
	if err = ds.checkpointManager.CreateCheckpoint(createResp.ID, constructPodSandboxCheckpoint(config)); err != nil {
		return nil, err
	}

	// Step 4: Start the sandbox container.
	// Assume kubelet's garbage collector would remove the sandbox later, if
	// startContainer failed.
	err = ds.client.StartContainer(createResp.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to start sandbox container for pod %q: %v", config.Metadata.Name, err)
	}

	// Step 5: Setup networking for the sandbox.
	// All pod networking is setup by a CNI plugin discovered at startup time.
	// This plugin assigns the pod ip, sets up routes inside the sandbox,
	// creates interfaces etc. In theory, its jurisdiction ends with pod
	// sandbox networking, but it might insert iptables rules or open ports
	// on the host as well, to satisfy parts of the pod spec that aren't
	// recognized by the CNI standard yet.
	cID := kubecontainer.BuildContainerID(runtimeName, createResp.ID)

	return resp, nil
}
````

### pod的生命周期和重启策略

pod在运行的过程中会被定义为各种状态，了解一些状态能帮助我们了解pod的调度策略。  
当 Pod 被创建之后，就会进入健康检查状态，当 Kubernetes 确定当前 Pod 已经能够接受外部的请求时，才会将流量打到
新的 Pod 上并继续对外提供服务，在这期间如果发生了错误就可能会触发重启机制。  
pod的重启侧罗包括  
- Always 只要失败，就会治重启
- OnFile 当容器终止运行，且退出吗不是0，就会重启
- Never  从来不会重启

重启的时间，是以2n来算。比如1,2,4,8.....最长延迟５分钟，并且在成功重启后的10分钟重置这个时间。