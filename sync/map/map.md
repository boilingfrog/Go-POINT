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

### load
