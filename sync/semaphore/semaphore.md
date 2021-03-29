## semaphore

### semaphore的作用

信号量是在并发编程中比较常见的一种同步机制，它会保证持有的计数器在0到初始化的权重之间，每次获取资源时都会将信号量中的计数器减去对应的数值，在释放时重新加回来，当遇到计数器大于信号量大小时就会进入休眠等待其他进程释放信号，我们常常会在控制访问资源的进程数量时用到。  

go中的`semaphore`，提供`sleep`和`wakeup`原语，使其能够在其它同步原语中的竞争情况下使用。当一个`goroutine`需要休眠时，将其进行集中存放，当需要`wakeup`时，再将其取出，重新放入调度器中。   

go中本身提供了`semaphore`的相关方法，不过只能在内部调用  

```go
// go/src/sync/runtime.go
func runtime_Semacquire(s *uint32)

func runtime_SemacquireMutex(s *uint32, lifo bool, skipframes int)

func runtime_Semrelease(s *uint32, handoff bool, skipframes int)
```

扩展包`golang.org/x/sync/semaphore`提供了一种带权重的信号量实现方式，我们可以按照不同的权重对资源的访问进行管理。   

### 如何使用


