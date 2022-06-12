<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Redis 中的分布式锁如何使用](#redis-%E4%B8%AD%E7%9A%84%E5%88%86%E5%B8%83%E5%BC%8F%E9%94%81%E5%A6%82%E4%BD%95%E4%BD%BF%E7%94%A8)
  - [分布式锁的使用场景](#%E5%88%86%E5%B8%83%E5%BC%8F%E9%94%81%E7%9A%84%E4%BD%BF%E7%94%A8%E5%9C%BA%E6%99%AF)
  - [使用 Redis 来实现分布式锁](#%E4%BD%BF%E7%94%A8-redis-%E6%9D%A5%E5%AE%9E%E7%8E%B0%E5%88%86%E5%B8%83%E5%BC%8F%E9%94%81)
    - [使用 `set key value px milliseconds nx` 实现](#%E4%BD%BF%E7%94%A8-set-key-value-px-milliseconds-nx-%E5%AE%9E%E7%8E%B0)
    - [SETNX+Lua 实现](#setnxlua-%E5%AE%9E%E7%8E%B0)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Redis 中的分布式锁如何使用

### 分布式锁的使用场景

为了保证我们线上服务的并发性和安全性，目前我们的服务一般抛弃了单体应用，采用的都是扩展性很强的分布式架构。    

对于可变共享资源的访问，同一时刻，只能由一个线程或者进程去访问操作。这时候我们就需要做个标识，如果当前有线程或者进程在操作共享变量，我们就做个标记，标识当前资源正在被操作中， 其它的线程或者进程，就不能进行操作了。当前操作完成之后，删除标记，这样其他的线程或者进程，就能来申请共享变量的操作。通过上面的标记来保证同一时刻共享变量只能由一个线程或者进行持有。  

- 对于单体应用：多个线程之间访问可变共享变量，比较容易处理，可简单使用内存来存储标示即可；  

- 分布式应用：这种场景下比较麻烦，因为多个应用，部署的地址可能在不同的机房，一个在北京一个在上海。不能简单的存储标示在内存中了，这时候需要使用公共内存来记录该标示，栗如 Redis，MySQL 。。。   

### 使用 Redis 来实现分布式锁

这里来聊聊如何使用 Redis 实现分布式锁  

Redis 中分布式锁一般会用 `set key value px milliseconds nx` 或者 `SETNX+Lua`来实现。  

因为 `SETNX` 命令，需要配合 `EXPIRE` 设置过期时间,Redis 中单命令的执行是原子性的，组合命令就需要使用 Lua 才能保证原子性了。  

看下如何实现  

#### 使用 `set key value px milliseconds nx` 实现  

因为这个命令同时能够设置键值和过期时间，同时Redis中的单命令都是原子性的，所以加锁的时候使用这个命令即可  

```go
func (r *Redis) TryLock(ctx context.Context, key, value string, expire time.Duration) (isGetLock bool, err error) {
	// 使用 set nx
	res, err := r.Do(ctx, "set", key, value, "px", expire.Milliseconds(), "nx").Result()
	if err != nil {
		return false, err
	}
	if res == "OK" {
		return true, nil
	}
	return false, nil
}
```

#### SETNX+Lua 实现

如果使用 SETNX 命令，这个命令不能设置过期时间，需要配合 EXPIRE 命令来使用。  

因为是用到了两个命令，这时候两个命令的组合使用是不能保障原子性的，在一些并发比较大的时候，需要配合使用 Lua 脚本来保证命令的原子性。  

```go
func tryLockScript() string {
	script := `
		local key = KEYS[1]

		local value = ARGV[1] 
		local expireTime = ARGV[2] 
		local isSuccess = redis.call('SETNX', key, value)

		if isSuccess == 1 then
			redis.call('EXPIRE', key, expireTime)
			return "OK"
		end

		return "unLock"    `
	return script
}

func (r *Redis) TryLock(ctx context.Context, key, value string, expire time.Duration) (isGetLock bool, err error) {
	// 使用 Lua + SETNX
	res, err := r.Eval(ctx, tryLockScript(), []string{key}, value, expire.Seconds()).Result()
	if err != nil {
		return false, err
	}
	if res == "OK" {
		return true, nil
	}
	return false, nil
}
```

除了上面加锁两个命令的区别之外，在解锁的时候需要注意下不能误删除别的线程持有的锁   

为什么会出现这种情况呢，这里来分析下  

举个栗子  

1、线程1获取了锁，锁的过期时间为1s；  

2、线程1完成了业务操作，用时1.5s ，这时候线程1的锁已经被过期时间自动释放了，这把锁已经被别的线程获取了；    

3、但是线程1不知道，接着去释放锁，这时候就会将别的线程的锁，错误的释放掉。  

面对这种情况，其实也很好处理  

1、设置 value 具有唯一性；  

2、每次删除锁的时候，先去判断下 value 的值是否能对的上，不相同就表示，锁已经被别的线程获取了；    

看下代码实现  

```go
var UnLockErr = errors.New("未解锁成功")

func unLockScript() string {
	script := `
		local value = ARGV[1] 
		local key = KEYS[1]

		local keyValue = redis.call('GET', key)
		if tostring(keyValue) == tostring(value) then
			return redis.call('DEL', key)
		else
			return 0
		end
    `
	return script
}

func (r *Redis) Unlock(ctx context.Context, key, value string) error {
	res, err := r.Eval(ctx, unLockScript(), []string{key}, value).Result()
	if err != nil {
		return err
	}
	if res.(int64) == 0 {
		return UnLockErr
	}

	return nil
}
```

代码可参考[lock](https://github.com/boilingfrog/Go-POINT/blob/master/redis/lock/lock.go)

### 参考

【Redis核心技术与实战】https://time.geekbang.org/column/intro/100056701    
【Redis设计与实现】https://book.douban.com/subject/25900156/   
【EVAL简介】http://www.redis.cn/commands/eval.html   
【Redis学习笔记】https://github.com/boilingfrog/Go-POINT/tree/master/redis
【Redis Lua脚本调试器】http://www.redis.cn/topics/ldb.html    
【redis中Lua脚本的使用】https://boilingfrog.github.io/2022/06/06/Redis%E4%B8%AD%E7%9A%84%E5%8E%9F%E5%AD%90%E6%93%8D%E4%BD%9C(2)-redis%E4%B8%AD%E4%BD%BF%E7%94%A8Lua%E8%84%9A%E6%9C%AC%E4%BF%9D%E8%AF%81%E5%91%BD%E4%BB%A4%E5%8E%9F%E5%AD%90%E6%80%A7/  


