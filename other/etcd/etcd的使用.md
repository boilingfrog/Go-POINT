<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [etcd的使用](#etcd%E7%9A%84%E4%BD%BF%E7%94%A8)
  - [什么是etcd](#%E4%BB%80%E4%B9%88%E6%98%AFetcd)
  - [etcd的特点](#etcd%E7%9A%84%E7%89%B9%E7%82%B9)
  - [etcd的应用场景](#etcd%E7%9A%84%E5%BA%94%E7%94%A8%E5%9C%BA%E6%99%AF)
    - [服务注册与发现](#%E6%9C%8D%E5%8A%A1%E6%B3%A8%E5%86%8C%E4%B8%8E%E5%8F%91%E7%8E%B0)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## etcd的使用

### 什么是etcd

ETCD是一个分布式、可靠的`key-value`存储的分布式系统，用于存储分布式系统中的关键数据；当然，它不仅仅用于存储，还提供配置共享及服务发现；基于Go语言实现  。

### etcd的特点

- 完全复制：集群中的每个节点都可以使用完整的存档

- 高可用性：Etcd可用于避免硬件的单点故障或网络问题

- 一致性：每次读取都会返回跨多主机的最新写入

- 简单：包括一个定义良好、面向用户的API（gRPC）

- 安全：实现了带有可选的客户端证书身份验证的自动化TLS

- 可靠：使用Raft算法实现了强一致、高可用的服务存储目录

### etcd的应用场景



### etcd部署

在使用之前先构建一个etcd

#### 服务注册与发现

服务发现还能注册

服务注册发现解决的是分布式系统中最常见的问题之一，即在同一个分布式系统中，找到我们需要的目标服务，建立连接，然后完成整个链路的调度。   

本质上来说，服务发现就是想要了解集群中是否有进程在监听 udp 或 tcp 端口，并且通过名字就可以查找和连接。要解决服务发现的问题，需要有下面三大支柱，缺一不可。  

1、**一个强一致性、高可用的服务存储目录**。基于Raft算法的etcd天生就是这样一个强一致性高可用的服务存储目录。  

2、**一种注册服务和监控服务健康状态的机制**。用户可以在etcd中注册服务，并且对注册的服务设置`key TTL`，定时保持服务的心跳以达到监控健康状态的效果。   

3、**一种查找和连接服务的机制**。通过在 etcd 指定的主题下注册的服务也能在对应的主题下查找到。为了确保连接，我们可以在每个服务机器上都部署一个Proxy模式的etcd，这样就可以确保能访问etcd集群的服务都能互相连接。    

<img src="/img/etcd_1.png" alt="etcd" align=center/>

一个用户的api请求，可能调用多个微服务资源，这些服务我们可以使用etcd进行服务注册和服务发现，当每个服务启动的时候就注册到etcd中，当我们需要使用的时候，直接在etcd中寻找，调用即可。  

当然，每个服务的实例不止一个，比如我们的用户服务，我们可能启动了多个实例，这些实例在服务启动过程中全部注册到了etcd中，但是某个实例可能出现故障重启，这时候就etcd在进行转发的时候，就会屏蔽到故障的实例节点，只向正常运行的实例，进行请求转发。   

<img src="/img/etcd-register.png" alt="etcd" align=center/>

来看个服务注册发现的demo

这里放一段比较核心的代码，这里摘录了我们线上正在使用的etcd实现grpc服务注册和发现的，具体的实现可参考,[etcd实现grpc的服务注册和服务发现](https://github.com/boilingfrog/daily-test/tree/master/etcd/discovery)  

对于etcd中的连接，我们每个都维护一个租约，通过KeepAlive自动续保。如果租约过期则所有附加在租约上的key将过期并被删除，即所对应的服务被拿掉。  

```go
// Register for grpc server
type Register struct {
	EtcdAddrs   []string
	DialTimeout int

	closeCh     chan struct{}
	leasesID    clientv3.LeaseID
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse

	srvInfo Server
	srvTTL  int64
	cli     *clientv3.Client
	logger  *zap.Logger
}

// NewRegister create a register base on etcd
func NewRegister(etcdAddrs []string, logger *zap.Logger) *Register {
	return &Register{
		EtcdAddrs:   etcdAddrs,
		DialTimeout: 3,
		logger:      logger,
	}
}

// Register a service
func (r *Register) Register(srvInfo Server, ttl int64) (chan<- struct{}, error) {
	var err error

	if strings.Split(srvInfo.Addr, ":")[0] == "" {
		return nil, errors.New("invalid ip")
	}

	if r.cli, err = clientv3.New(clientv3.Config{
		Endpoints:   r.EtcdAddrs,
		DialTimeout: time.Duration(r.DialTimeout) * time.Second,
	}); err != nil {
		return nil, err
	}

	r.srvInfo = srvInfo
	r.srvTTL = ttl

	if err = r.register(); err != nil {
		return nil, err
	}

	r.closeCh = make(chan struct{})

	go r.keepAlive()

	return r.closeCh, nil
}

// Stop stop register
func (r *Register) Stop() {
	r.closeCh <- struct{}{}
}

// register 注册节点
func (r *Register) register() error {
	leaseCtx, cancel := context.WithTimeout(context.Background(), time.Duration(r.DialTimeout)*time.Second)
	defer cancel()

	// 分配一个租约
	leaseResp, err := r.cli.Grant(leaseCtx, r.srvTTL)
	if err != nil {
		return err
	}
	r.leasesID = leaseResp.ID 
	// 自动定时的续约某个租约。
	if r.keepAliveCh, err = r.cli.KeepAlive(context.Background(), leaseResp.ID); err != nil {
		return err
	}

	data, err := json.Marshal(r.srvInfo)
	if err != nil {
		return err
	}
	_, err = r.cli.Put(context.Background(), BuildRegPath(r.srvInfo), string(data), clientv3.WithLease(r.leasesID))
	return err
}

// unregister 删除节点
func (r *Register) unregister() error {
	_, err := r.cli.Delete(context.Background(), BuildRegPath(r.srvInfo))
	return err
}

// keepAlive
func (r *Register) keepAlive() {
	ticker := time.NewTicker(time.Duration(r.srvTTL) * time.Second)
	for {
		select {
		case <-r.closeCh:
			if err := r.unregister(); err != nil {
				r.logger.Error("unregister failed", zap.Error(err))
			}
			if _, err := r.cli.Revoke(context.Background(), r.leasesID); err != nil {
				r.logger.Error("revoke failed", zap.Error(err))
			}
			return
		case res := <-r.keepAliveCh:
			if res == nil {
				if err := r.register(); err != nil {
					r.logger.Error("register failed", zap.Error(err))
				}
			}
		case <-ticker.C:
			if r.keepAliveCh == nil {
				if err := r.register(); err != nil {
					r.logger.Error("register failed", zap.Error(err))
				}
			}
		}
	}
}
```


### 参考

【一文入门ETCD】https://juejin.cn/post/6844904031186321416   
【etcd：从应用场景到实现原理的全方位解读】https://www.infoq.cn/article/etcd-interpretation-application-scenario-implement-principle   
【Etcd 架构与实现解析】http://jolestar.com/etcd-architecture/   
【linux单节点和集群的etcd】https://www.jianshu.com/p/07ca88b6ff67  