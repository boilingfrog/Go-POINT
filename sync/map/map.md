<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [sync.map](#syncmap)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [深入了解下](#%E6%B7%B1%E5%85%A5%E4%BA%86%E8%A7%A3%E4%B8%8B)
  - [查看下具体的实现](#%E6%9F%A5%E7%9C%8B%E4%B8%8B%E5%85%B7%E4%BD%93%E7%9A%84%E5%AE%9E%E7%8E%B0)
    - [Load](#load)
    - [Store](#store)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## sync.map

### 前言

Go中的map不是并发安全的，在Go1.9之后，引入了`sync.Map`,并发安全的map。  

### 深入了解下

对于map，我们常用的做法就是加锁。  

对于`sync.Map`这些操作则是不需要的，来看下`sync.Map`的特点:  

- 1、以空间换效率，通过read和dirty两个map来提高读取效率
- 2、优先从read map中读取(无锁)，否则再从dirty map中读取(加锁)
- 3、动态调整，当misses次数过多时，将dirty map提升为read map
- 4、延迟删除，删除只是为value打一个标记，在dirty map提升时才执行真正的删除

简单的使用栗子：

```go
func syncMapDemo() {

	var smp sync.Map

	// 数据写入
	smp.Store("name", "小红")
	smp.Store("age", 18)

	// 数据读取
	name, _ := smp.Load("name")
	fmt.Println(name)

	age, _ := smp.Load("age")
	fmt.Println(age)

	// 遍历
	smp.Range(func(key, value interface{}) bool {
		fmt.Println(key, value)
		return true
	})

	// 删除
	smp.Delete("age")
	age, ok := smp.Load("age")
	fmt.Println("删除后的查询")
	fmt.Println(age, ok)

	// 读取或写入,存在就读取，不存在就写入
	smp.LoadOrStore("age", 100)
	age, _ = smp.Load("age")
	fmt.Println("不存在")
	fmt.Println(age)

	smp.LoadOrStore("age", 99)
	age, _ = smp.Load("age")
	fmt.Println("存在")
	fmt.Println(age)
}
```

### 查看下具体的实现

```go
// sync/map.go
type Map struct {
   // 当写read map 或读写dirty map时 需要上锁
   mu Mutex

   // read map的 k v(entry) 是不变的，删除只是打标记，插入新key会加锁写到dirty中
   // 因此对read map的读取无需加锁
   read atomic.Value // 保存readOnly结构体

   // dirty map 对dirty map的操作需要持有mu锁
   dirty map[interface{}]*entry

   // 当Load操作在read map中未找到，尝试从dirty中进行加载时(不管是否存在)，misses+1
   // 当misses达到diry map len时，dirty被提升为read 并且重新分配dirty
   misses int
}

// read map数据结构
type readOnly struct {
   m       map[interface{}]*entry
   // 为true时代表dirty map中含有m中没有的元素
   amended bool
}

type entry struct {
   // 指向实际的interface{}
   // p有三种状态:
   // p == nil: 键值已经被删除，此时，m.dirty==nil 或 m.dirty[k]指向该entry
   // p == expunged: 键值已经被删除， 此时, m.dirty!=nil 且 m.dirty不存在该键值
   // 其它情况代表实际interface{}地址 如果m.dirty!=nil 则 m.read[key] 和 m.dirty[key] 指向同一个entry
   // 当删除key时，并不实际删除，先CAS entry.p为nil 等到每次dirty map创建时(dirty提升后的第一次新建Key)，会将entry.p由nil CAS为expunged
   p unsafe.Pointer // *interface{}
}
```

`read map` 和 `dirty map` 的存储方式是不一致的。  

前者使用 `atomic.Value`，后者只是单纯的使用 map。原因是 `read map` 使用 `lock free` 操作，必须保证 `load/store` 的原子性；而 `dirty map` 的 `load+store` 操作是由 lock（就是 mu）来保护的。

1、read和dirty通过entry包装value，这样使得value的变化和map的变化隔离，前者可以用atomic无锁完成  
2、Map的read字段结构体定义为readOnly，这只是针对map[interface{}]*entry而言的，entry内的内容以及amended字段都是可以变的  
3、大部分情况下，对已有key的删除(entry.p置为nil)和更新可以直接通过修改entry.p来完成  


#### Load

```go
func (m *Map) Load(key interface{}) (value interface{}, ok bool) {
    // 首先在通过atomic的原子操作读取内容
	read, _ := m.read.Load().(readOnly)
	e, ok := read.m[key]
    // 如果没在 read 中找到，并且 amended 为 true，即 dirty 中存在 read 中没有的 key
	if !ok && read.amended {
       // read调用了atomic的原子性，所以不用加锁，但是dirty map[interface{}]*entry就需要了，用的互斥锁
		m.mu.Lock()
        // double check，避免在加锁的时候dirty map提升为read map
		read, _ = m.read.Load().(readOnly)
		e, ok = read.m[key]
        // 还是没有找到
		if !ok && read.amended {
            // 从 dirty 中找
			e, ok = m.dirty[key]
            // 不管dirty中有没有找到 都增加misses计数 该函数可能将dirty map提升为readmap
			m.missLocked()
		}
		m.mu.Unlock()
	}
	if !ok {
		return nil, false
	}
	return e.load()
}

// 从entry中atomic load实际interface{}
func (e *entry) load() (value interface{}, ok bool) {
  p := atomic.LoadPointer(&e.p)
  if p == nil || p == expunged {
    return nil, false
  }
  return *(*interface{})(p), true
}

```

梳理下处理的逻辑：  

1、首先是 fast path，直接在 read 中找，如果找到了直接调用 entry 的 load 方法，取出其中的值。  
2、如果 read 中没有这个 key，且 amended 为 fase，说明 dirty 为空，那直接返回 空和 false。   
3、如果 read 中没有这个 key，且 amended 为 true，说明 dirty 中可能存在我们要找的 key。当然要先上锁，再尝试去 dirty 中查找。在这之前，仍然有一个 double check 的操作。若还是没有在 read 中找到，那么就从 dirty 中找。不管 dirty 中有没有找到，都要"记一笔"，因为在 dirty 被提升为 read 之前，都会进入这条路径  

```go
// 增加misses计数，并在必要的时候提升dirty map
// 如果 misses 值小于 m.dirty 的长度，就直接返回。否则，将 m.dirty 晋升为 read，并清空 dirty，清空 misses 计数值。这样，之前
// 一段时间新加入的 key 都会进入到 read 中，从而能够提升 read 的命中率。
func (m *Map) missLocked() {
  m.misses++
  if m.misses < len(m.dirty) {
    return
  }
  // 提升过程很简单，直接将m.dirty赋给m.read.m
  // 提升完成之后 amended == false m.dirty == nil
  // m.dirty并不立即创建被拷贝元素，而是延迟创建
  m.read.Store(readOnly{m: m.dirty})
  m.dirty = nil
  m.misses = 0
}
```

对于`missLocked`会直接将 misses 的值加 1，表示一次未命中，如果 misses 值小于 m.dirty 的长度，就直接返回。否则，将 m.dirty 晋升为 read，并清空 dirty，清空 misses 计数值。这样，之前一段时间新加入的 key 都会进入到 read 中，从而能够提升 read 的命中率。  
             
#### Store

```go
// Store sets the value for a key.
func (m *Map) Store(key, value interface{}) {
    // 如果read map中存在该key  则尝试直接更改(由于修改的是entry内部的pointer，因此dirty map也可见)
	read, _ := m.read.Load().(readOnly)
	if e, ok := read.m[key]; ok && e.tryStore(&value) {
		return
	}

	m.mu.Lock()
	read, _ = m.read.Load().(readOnly)
	if e, ok := read.m[key]; ok {
		if e.unexpungeLocked() {
			// 如果 read map 中存在该 key，但 p == expunged，则说明 m.dirty != nil 并且 m.dirty 中不存在该 key 值 此时:
			//    a. 将 p 的状态由 expunged  更改为 nil
			//    b. dirty map 插入 key
			m.dirty[key] = e
		}
        // 更新 entry.p = value (read map 和 dirty map 指向同一个 entry)
		e.storeLocked(&value)
	} else if e, ok := m.dirty[key]; ok {
        // 如果 read map 中不存在该 key，但 dirty map 中存在该 key，直接写入更新 entry(read map 中仍然没有这个 key)
		e.storeLocked(&value)
	} else {
		// 如果read map和dirty map中都不存在该key，则:
		//    a. 如果dirty map为空，则需要创建dirty map，并从read map中拷贝未删除的元素
		//    b. 更新amended字段，标识dirty map中存在read map中没有的key
		//    c. 将k v写入dirty map中，read.m不变
		if !read.amended {

			m.dirtyLocked()
			m.read.Store(readOnly{m: read.m, amended: true})
		}
		m.dirty[key] = newEntry(value)
	}
	m.mu.Unlock()
}

// 尝试直接更新entry 如果p == expunged 返回false
func (e *entry) tryStore(i *interface{}) bool {
  p := atomic.LoadPointer(&e.p)
  if p == expunged {
    return false
  }
  for {
    if atomic.CompareAndSwapPointer(&e.p, p, unsafe.Pointer(i)) {
      return true
    }
    p = atomic.LoadPointer(&e.p)
    if p == expunged {
      return false
    }
  }
}

func (e *entry) unexpungeLocked() (wasExpunged bool) {
  return atomic.CompareAndSwapPointer(&e.p, expunged, nil)
}

// 如果 dirty map为nil，则从read map中拷贝元素到dirty map
func (m *Map) dirtyLocked() {
  if m.dirty != nil {
    return
  }

  read, _ := m.read.Load().(readOnly)
  m.dirty = make(map[interface{}]*entry, len(read.m))
  for k, e := range read.m {
    // a. 将所有为 nil的 p 置为 expunged
    // b. 只拷贝不为expunged 的 p
    if !e.tryExpungeLocked() {
      m.dirty[k] = e
    }
  }
}

func (e *entry) tryExpungeLocked() (isExpunged bool) {
  p := atomic.LoadPointer(&e.p)
  for p == nil {
    if atomic.CompareAndSwapPointer(&e.p, nil, expunged) {
      return true
    }
    p = atomic.LoadPointer(&e.p)
  }
  return p == expunged
}
```
梳理下流程：  

1、首先还是去read map中查询，存在并且p!=expunged,直接修改。（由于修改的是 entry 内部的 pointer，因此 dirty map 也可见）
2、如果read map中存在该key，但p == expunged。加锁更新p的状态，然后直接更新该entry (此时m.dirty==nil或m.dirty[key]==e)  
3、如果read map中不存在该Key，但dirty map中存在该key，直接写入更新entry(read map中仍然没有)   
4、如果read map和dirty map都不存在该key     
- a. 如果dirty map为空，则需要创建dirty map，并从read map中拷贝未删除的元素  
- b. 更新amended字段，标识dirty map中存在read map中没有的key  
- c. 将k v写入dirty map中，read.m不变  



### 参考
【Go sync.Map 实现】https://wudaijun.com/2018/02/go-sync-map-implement/