<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [etcd中watch的源码解析](#etcd%E4%B8%ADwatch%E7%9A%84%E6%BA%90%E7%A0%81%E8%A7%A3%E6%9E%90)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [client端的代码](#client%E7%AB%AF%E7%9A%84%E4%BB%A3%E7%A0%81)
  - [server端的代码实现](#server%E7%AB%AF%E7%9A%84%E4%BB%A3%E7%A0%81%E5%AE%9E%E7%8E%B0)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## etcd中watch的源码解析

### 前言

etcd是一个cs网络架构，源码分析应该涉及到client端，server端。client主要是提供操作来请求对key监听，并且接收key变更时的通知。server要能做到接收key监听请求，并且启动定时器等方法来对key进行监听，有变更时通知client。  

### client端的代码  

```go
// client/v3/watch.go

type Watcher interface {
	// 在键或前缀上监听。将监听的事件
	// 通过定义的返回的channel进行返回。如果修订等待通过
	// 监听被压缩，然后监听将被服务器取消，
	// 客户端将发布压缩的错误观察响应，并且通道将关闭。
	// 如果请求的修订为 0 或未指定，则返回的通道将
	// 返回服务器收到监视请求后发生的监视事件。
	// 如果上下文“ctx”被取消或超时，返回的“WatchChan”关闭，
	// 并且来自此关闭通道的“WatchResponse”具有零事件且为零“Err()”。
	// 一旦不再使用观察者，上下文“ctx”必须被取消，
	// 释放相关资源。
	//
	// 如果上下文是“context.Background/TODO”，则返回“WatchChan”
	// 不会被关闭和阻塞直到事件被触发，除非服务器
	// 返回一个不可恢复的错误（例如 ErrCompacted）。
	// 例如，当上下文通过“WithRequireLeader”和
	// 连接的服务器没有领导者（例如，由于网络分区），
	// 将返回错误“etcdserver: no leader”（ErrNoLeader），
	// 然后 "WatchChan" 以非零 "Err()" 关闭。
	// 为了防止观察流卡在分区节点中，
	// 确保使用“WithRequireLeader”包装上下文。
	//
	// 否则，只要上下文没有被取消或超时，
	// watch 将永远重试其他可恢复的错误，直到重新连接。
	//
	// TODO：在最后一个“WatchResponse”消息中显式设置上下文错误并关闭通道？
	// 目前，客户端上下文被永远不会关闭的“valCtx”覆盖。
	// TODO(v3.4): 配置watch重试策略，限制最大重试次数
	//（参见 https://github.com/etcd-io/etcd/issues/8980）
	Watch(ctx context.Context, key string, opts ...OpOption) WatchChan

	// RequestProgress requests a progress notify response be sent in all watch channels.
	RequestProgress(ctx context.Context) error

	// Close closes the watcher and cancels all watch requests.
	Close() error
}

// watcher implements the Watcher interface
type watcher struct {
	remote   pb.WatchClient
	callOpts []grpc.CallOption

	// mu protects the grpc streams map
	mu sync.Mutex

	// streams 保存所有由 ctx 值键控的活动 grpc 流。
	streams map[string]*watchGrpcStream
	lg      *zap.Logger
}

// watchGrpcStream 跟踪附加到单个 grpc 流的所有watch资源。
type watchGrpcStream struct {
	owner    *watcher
	remote   pb.WatchClient
	callOpts []grpc.CallOption

	// ctx 控制内部的remote.Watch requests
	ctx context.Context
	// ctxKey 用来找流的上下文信息
	ctxKey string
	cancel context.CancelFunc

	// substreams 持有此 grpc 流上的所有活动的watchers
	substreams map[int64]*watcherStream
	// 恢复保存此 grpc 流上的所有正在恢复的观察者
	resuming []*watcherStream

	// reqc 从 Watch() 向主协程发送观察请求
	reqc chan watchStreamRequest
	// respc 从 watch 客户端接收数据
	respc chan *pb.WatchResponse
	// donec 通知广播进行退出
	donec chan struct{}
	// errc transmits errors from grpc Recv to the watch stream reconnect logic
	errc chan error
	// Closec 获取关闭观察者的观察者流
	closingc chan *watcherStream
	// 当所有子流 goroutine 都退出时，wg 完成
	wg sync.WaitGroup

	// resumec 关闭以表示所有子流都应开始恢复
	resumec chan struct{}
	// closeErr 是关闭监视流的错误
	closeErr error

	lg *zap.Logger
}

// watcherStream 代表注册的观察者
// watch()时，构造watchgrpcstream时构造的watcherStream，用于封装一个watch rpc请求，包含订阅监听key，通知key变更通道，一些重要标志。
type watcherStream struct {
	// initReq 是发起这个请求的请求
	initReq watchRequest

	// outc 向订阅者发布watch响应
	outc chan WatchResponse
	// recvc buffers watch responses before publishing
	recvc chan *WatchResponse
	// 当 watcherStream goroutine 停止时 donec 关闭
	donec chan struct{}
	// 当应该安排流关闭时，closures 设置为 true。
	closing bool
	// id 是在 grpc 流上注册的 watch id
	id int64

	// buf 保存从 etcd 收到但尚未被客户端消费的所有事件
	buf []*WatchResponse
}

// 1、key是否满足watch的条件
// 2、过滤监听事件
// 3、构造watch请求
// 4、查找或分配新的grpc watch stream
// 5、发送watch请求到reqc通道
// 6、返回WatchResponse 接收chan给客户端
func (w *watcher) Watch(ctx context.Context, key string, opts ...OpOption) WatchChan {
	ow := opWatch(key, opts...)

	var filters []pb.WatchCreateRequest_FilterType
	if ow.filterPut {
		filters = append(filters, pb.WatchCreateRequest_NOPUT)
	}
	if ow.filterDelete {
		filters = append(filters, pb.WatchCreateRequest_NODELETE)
	}

	wr := &watchRequest{
		ctx:            ctx,
		createdNotify:  ow.createdNotify,
		key:            string(ow.key),
		end:            string(ow.end),
		rev:            ow.rev,
		progressNotify: ow.progressNotify,
		fragment:       ow.fragment,
		filters:        filters,
		prevKV:         ow.prevKV,
		retc:           make(chan chan WatchResponse, 1),
	}

	ok := false
	ctxKey := streamKeyFromCtx(ctx)

	var closeCh chan WatchResponse
	for {
		// 查找或分配适当的 grpc 监视流
		w.mu.Lock()
		if w.streams == nil {
			// closed
			w.mu.Unlock()
			ch := make(chan WatchResponse)
			close(ch)
			return ch
		}

		// streams是一个map,保存所有由 ctx 值键控的活动 grpc 流
		// 如果该请求对应的流为空,则新建
		wgs := w.streams[ctxKey]
		if wgs == nil {
            // newWatcherGrpcStream new一个watch grpc stream来传输watch请求
            // 创建goroutine来处理监听key的watch各种事件
			wgs = w.newWatcherGrpcStream(ctx)
			w.streams[ctxKey] = wgs
		}
		donec := wgs.donec
		reqc := wgs.reqc
		w.mu.Unlock()

		// couldn't create channel; return closed channel
		if closeCh == nil {
			closeCh = make(chan WatchResponse, 1)
		}

		// 等待接收值
		select {
		// reqc 从 Watch() 向主协程发送观察请求
		case reqc <- wr:
			ok = true
		case <-wr.ctx.Done():
			ok = false
		case <-donec:
			ok = false
			if wgs.closeErr != nil {
				closeCh <- WatchResponse{Canceled: true, closeErr: wgs.closeErr}
				break
			}
			// 重试，可能已经从没有 ctxs 中删除了流
			continue
		}

		// receive channel
		if ok {
			select {
			case ret := <-wr.retc:
				return ret
			case <-ctx.Done():
			case <-donec:
				if wgs.closeErr != nil {
					closeCh <- WatchResponse{Canceled: true, closeErr: wgs.closeErr}
					break
				}
				// 重试，可能已经从没有 ctxs 中删除了流
				continue
			}
		}
		break
	}

	close(closeCh)
	return closeCh
}

// newWatcherGrpcStream new一个watch grpc stream来传输watch请求
func (w *watcher) newWatcherGrpcStream(inctx context.Context) *watchGrpcStream {
    ctx, cancel := context.WithCancel(&valCtx{inctx})

    //构造watchGrpcStream
    wgs := &watchGrpcStream{
        owner:      w,
        remote:     w.remote,
        callOpts:   w.callOpts,
        ctx:        ctx,
        ctxKey:     streamKeyFromCtx(inctx),
        cancel:     cancel,
        substreams: make(map[int64]*watcherStream),
        respc:      make(chan *pb.WatchResponse),
        reqc:       make(chan watchStreamRequest),
        donec:      make(chan struct{}),
        errc:       make(chan error, 1),
        closingc:   make(chan *watcherStream),
        resumec:    make(chan struct{}),
    }

    // 创建goroutine来处理监听key的watch各种事件
    go wgs.run()
    return wgs
}

// 通过etcd grpc服务器启动一个watch stream
// run 管理watch 的事件chan
func (w *watchGrpcStream) run() {
	var wc pb.Watch_WatchClient
	var closeErr error
	...
	// 创建一个grpc client连接etcd grpc server。
	if wc, closeErr = w.newWatchClient(); closeErr != nil {
		return
	}

	cancelSet := make(map[int64]struct{})

	// select检测各个chan的事件（reqc、respc、errc、closingc）
	var cur *pb.WatchResponse
	for {
		select {
			...
		}
	}
}

// 1、将所有订阅的stream标记为恢复
// 2、连接到grpc stream，并且接受watch取消
// 3、关闭出错的client stream，并且创建goroutine，用于转发从run()得到的响应给订阅者
// 4、创建goroutine接收来自新grpc流的数据
func (w *watchGrpcStream) newWatchClient() (pb.Watch_WatchClient, error) {
	// 将所有订阅的stream标记为恢复
	close(w.resumec)
	w.resumec = make(chan struct{})
	w.joinSubstreams()
	for _, ws := range w.substreams {
		ws.id = -1
		w.resuming = append(w.resuming, ws)
	}
	// 去掉无用，即为nil的stream
	var resuming []*watcherStream
	for _, ws := range w.resuming {
		if ws != nil {
			resuming = append(resuming, ws)
		}
	}
	w.resuming = resuming
	w.substreams = make(map[int64]*watcherStream)

	// 连接到grpc stream，并且接受watch取消
	stopc := make(chan struct{})
	donec := w.waitCancelSubstreams(stopc)
	wc, err := w.openWatchClient()
	close(stopc)
	<-donec

	// 对于client出错的stream，可以关闭，并且创建一个goroutine，用于转发从run()得到的响应给订阅者
	for _, ws := range w.resuming {
		if ws.closing {
			continue
		}
		ws.donec = make(chan struct{})
		w.wg.Add(1)
		go w.serveSubstream(ws, w.resumec)
	}

	if err != nil {
		return nil, v3rpc.Error(err)
	}

	// 创建goroutine接收来自新grpc流的数据
	go w.serveWatchClient(wc)
	return wc, nil
}
```

总结：  

1、`etcd v3 API`采用了gRPC ，而 gRPC 又利用了`HTTP/2 TCP` 链接多路复用（` multiple stream per tcp connection ）`，这样同一个Client的不同watch可以共享同一个TCP连接。  

2、watch支持指定单个 key，也可以指定一个 key 的前缀；  

3、Watch观察将要发生或者已经发生的事件，输入和输出都是流，输入流用于创建和取消观察，输出流发送事件；  

4、WatcherGrpcStream会启动一个协程专门用于通过 gRPC client stream 接收Server端的 watch response，然后将watch response send 到WatcherGrpcStream的watch response channel。  
 
5、 WatcherGrpcStream 也有一个专门的 协程专门用于从watch response channel 读数据，拿到watch response之后，会根据response里面的watchId 从WatcherGrpcStream的map[watchID] WatcherStream 中拿到对应的WatcherStream，并send到WatcherStream里面的WatchReponse channel。  

6、这里的watchId其实是Server端返回给client端的，当client Send Watch request给Server端时候，response会带上watchId, 这个watchId是与watch key是一一对应关系，然后client会建立WatchId与WatcherStream的映射关系。  

7、WatcherStream是具体的 watch response的处理结构，对于每个watch key，WatcherGrpcStream 也会启动一个专门的协程处理WatcherStream里面的watch response channel。  

### server端的代码实现  

来看下总体的架构  

1、etcd服务端创建newWatchableStore开启group监听；  

2、调用mvcc中syncWatchers将所有未通知的事件通知给所有的监听者；  

3、对watcher通道阻塞时存入victim中数据，开启syncVictimsLoop；  

4、watchServer响应客户端请求，发起watchStream及watcher实例新建，并将其添加至unsynced或synced中；  

5、client端通过grpc proxy向watcherServer发送watcher请求；  

6、grpc proxy提供对同一个key的多次watch合并减少etcd server中重复watcher创建，以提高etcd server稳定性。  


<img src="/img/etcd-server.png" alt="etcd" align=center/>

主要的文件  

```
/etcdserver/api/v3rpc/watch.go    watch 服务端实现

/mvcc/watcher.go  主要封装 watchStream 的实现

/mvcc/watchable_store.go  watch 版本的 KV 存储实现

/mvcc/watchable_store_txn.go  主要实现事务提交后 End() 函数的处理

/mvcc/watcher_group.go
```

```go
// 文件：/etcdserver/api/v3rpc/watch.go
type watchServer struct {
	...
	watchable mvcc.WatchableKV // 键值存储
	....
}

type serverWatchStream struct {
	...
	watchable mvcc.WatchableKV //kv 存储
	...
	// 与客户端进行连接的 Stream
	gRPCStream  pb.Watch_WatchServer
	// key 变动的消息管道
	watchStream mvcc.WatchStream
	// 响应客户端请求的消息管道
	ctrlStream  chan *pb.WatchResponse
	...
	// 该类型的 watch，服务端会定时发送类似心跳消息
	progress map[mvcc.WatchID]bool
	// 该类型表明，对于/a/b 这样的监听范围, 如果 b 变化了， 前缀/a也需要通知
	prevKV map[mvcc.WatchID]bool
	// 该类型表明，传输数据量大于阈值，需要拆分发送
	fragment map[mvcc.WatchID]bool
} 
```

```go
// 文件：/mvcc/watcher.go 
// 响应结构体
type WatchResponse struct {
	WatchID WatchID
	// 当前 watchResponse 实例创建时对应的 revision 值
	Revision int64
	// 压缩操作对应的 revison
	CompactRevision int64
}

type watchStream struct {
	// 用来记录关联的 watchableStore
	watchable watchable
	// event 事件写入通道
	ch        chan WatchResponse
	...
	cancels  map[WatchID]cancelFunc
	// 用来记录唯一标识与 watcher 的实例的关系
	watchers map[WatchID]*watcher
}
```

```go
// 文件：/mvcc/watcher_group.go 
type eventBatch struct {
	// 其中记录的 Event 实例是按照 revision 排序的
	evs []mvccpb.Event
	// 记录当前 eventBatch 中，有多少个来自不同的 main revison,
	revs int
	// 当前 eventBatch 记录的 Event 个数达到上限之后，后续 Event 实例无法加入该 eventBatch 中
	// 该字段记录了无法加入该 eventBatch 实例的第一个main revision
	moreRev int64
}

type watcherBatch map[*watcher]*eventBatch
type watcherSet map[*watcher]struct{}
type watcherSetByKey map[string]watcherSet

type watcherGroup struct {
	// 记录监听单个 Key 的 watch 实例
	keyWatchers watcherSetByKey
	// 记录进行范围监听的 watcher 实例。
	ranges adt.IntervalTree
	// 记录当前 wathcer 的全部实例
	watchers watcherSet
}
```

```go
type WatchStream interface {
	// Watch 创建了一个观察者. 观察者监听发生在给定的键或范围[key, end]上的事件的变化。
	//
	// 整个事件历史可以被观察，除非压缩。
	// 如果"startRev" <=0, watch观察当前之后的事件。

	// 将返回watcher的id，它显示为WatchID
	// 通过流通道发送给创建的监视器的事件。
	// watch ID在不等于AutoWatchID的时候被使用，否则将会返回一个自增的id
	Watch(id WatchID, key, end []byte, startRev int64, fcs ...FilterFunc) (WatchID, error)

	// Chan返回一个Chan。所有的观察响应将被发送到返回的chan。
	Chan() <-chan WatchResponse

	// RequestProgress请求给定ID的观察者的进度。响应只在观察者当前同步时发送。
	// 响应将通过附加的WatchRespone Chan发送,使用这个流来确保正确的排序。
	// 相应不包含事件。响应中的修订是进度的观察者，因为观察者当前已同步。
	RequestProgress(id WatchID)

	// Cancel 通过给出它的 ID 来取消一个观察者。如果 watcher 不存在，则会报错
	Cancel(id WatchID) error

	// Close closes Chan and release all related resources.
	Close()

	// Rev 返回流监视的 KV 的当前版本。
	Rev() int64
}
```

看下主要函数的实现  

整体的设计思路  

1、etcd服务端创建newWatchableStore开启group监听；  

2、

```go
// 文件 /mvcc/watchable_store.go
type watchableStore struct {
	*store
	mu sync.RWMutex
	// 当ch被阻塞时，对应 watcherBatch 实例会暂时记录到这个字段
	victims []watcherBatch
	// 当有新的 watcherBatch 实例添加到 victims 字段时，会向该通道发送消息
	victimc chan struct{}
	// 未同步的 watcher
	unsynced watcherGroup
	// 已完成同步的 watcher
	synced watcherGroup
	stopc chan struct{}
	wg sync.WaitGroup
}

type watcher struct {
	// 监听起始值
	key []byte
	// 监听终止值，  key 和 end 共同组成一个键值范围
	end []byte
	// 是否被阻塞
	victim bool
	// 是否压缩
	compacted bool
	...
	// 最小的 revision main
	minRev int64
	id     WatchID
	...
	ch chan<- WatchResponse
}

// server/mvcc/watchable_store.go
func newWatchableStore(lg *zap.Logger, b backend.Backend, le lease.Lessor, cfg StoreConfig) *watchableStore {
	if lg == nil {
		lg = zap.NewNop()
	}
	s := &watchableStore{
		store:    NewStore(lg, b, le, cfg),
		victimc:  make(chan struct{}, 1),
		unsynced: newWatcherGroup(),
		synced:   newWatcherGroup(),
		stopc:    make(chan struct{}),
	}
	s.store.ReadView = &readView{s}
	s.store.WriteView = &writeView{s}
	if s.le != nil {
		// use this store as the deleter so revokes trigger watch events
		s.le.SetRangeDeleter(func() lease.TxnDelete { return s.Write(traceutil.TODO()) })
	}
	s.wg.Add(2)
	// 开2个协程
    // syncWatchersLoop 每 100 毫秒同步一次未同步映射中的观察者。
	go s.syncWatchersLoop()
    // syncVictimsLoop 同步预先发送未成功的watchers
	go s.syncVictimsLoop()
	return s
}
```

总结

1、初始化一个watchableStore；  

2、启动了两个协程  

- syncWatchersLoop:每 100 毫秒同步一次未同步映射中的观察者；  

- 







```go
// server/mvcc/watcher.go
type watchStream struct {
	// 用来记录关联的 watchableStore
	watchable watchable
	// event 事件写入通道
	ch        chan WatchResponse
	...
	cancels  map[WatchID]cancelFunc
	// 用来记录唯一标识与 watcher 的实例的关系
	watchers map[WatchID]*watcher
}

// Watch 在流中创建一个新的 watcher 并返回它的 WatchID。
func (ws *watchStream) Watch(id WatchID, key, end []byte, startRev int64, fcs ...FilterFunc) (WatchID, error) {
	// prevent wrong range where key >= end lexicographically
	// watch request with 'WithFromKey' has empty-byte range end
	if len(end) != 0 && bytes.Compare(key, end) != -1 {
		return -1, ErrEmptyWatcherRange
	}

	ws.mu.Lock()
	defer ws.mu.Unlock()
	if ws.closed {
		return -1, ErrEmptyWatcherRange
	}

	// watch ID在不等于AutoWatchID的时候被使用，否则将会返回一个自增的id
	if id == AutoWatchID {
		for ws.watchers[ws.nextID] != nil {
			ws.nextID++
		}
		id = ws.nextID
		ws.nextID++
	} else if _, ok := ws.watchers[id]; ok {
		return -1, ErrWatcherDuplicateID
	}

	w, c := ws.watchable.watch(key, end, startRev, id, ws.ch, fcs...)

	ws.cancels[id] = c
	ws.watchers[id] = w
	return id, nil
}

func (s *watchableStore) watch(key, end []byte, startRev int64, id WatchID, ch chan<- WatchResponse, fcs ...FilterFunc) (*watcher, cancelFunc) {
	wa := &watcher{
		key:    key,
		end:    end,
		minRev: startRev,
		id:     id,
		ch:     ch,
		fcs:    fcs,
	}
    // 先上一把大的互斥锁
    // 多个watch操作，通过这个互斥锁，保证数据的顺序
	s.mu.Lock()
    // 里面上一把小的读锁
    // 读操作优先，保护读操作
	s.revMu.RLock()
	// 比较 startRev 和 currentRev，决定添加的 watcher 实例是否已经同步
	synced := startRev > s.store.currentRev || startRev == 0
	if synced {
		wa.minRev = s.store.currentRev + 1
		if startRev > wa.minRev {
			wa.minRev = startRev
		}
        // 添加到已同步的 watcher中
		s.synced.add(wa)
	} else {
		slowWatcherGauge.Inc()
        // 添加到未同步的 watcher中
		s.unsynced.add(wa)
	}
	s.revMu.RUnlock()
	s.mu.Unlock()

	watcherGauge.Inc()

	return wa, func() { s.cancelWatcher(wa) }
}
```
