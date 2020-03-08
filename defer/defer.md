- [defer](#defer)
  * [前言](#前言)
  * [defer的定义](#defer)

## defer

### 前言

defer作为go里面一个延迟调用的机制，它的存在能够大大的帮助我们优化我们的代码结构。但是我们
要弄明白defer的使用机制，不然我们的程序会发生很多莫名的问题。

### defer的定义 

defer用于延迟指定的函数，只能出现在函数的内部，由defer关键字以及指针对某个函数的调用表达式组成。这里被调用的函数
称为延迟函数。

### defer执行的规则

- 当外围函数中的语句执行完毕之时，只有当其中所有的延迟函数都执行完毕，外围函数才会真正的执结束执行。
- 当执行外围函数的return语句时，只有其中所有的延迟函数都执行完毕后，该外围函数才会真正的返回。
- 当外围函数中的代码引起运行恐慌时，只有当其中所有的延迟函数u调用到都执行完毕后，该运行恐慌才会真正被扩散至调用函数。

### 为什么需要defer

程序员在编程的时候，经常需要打开一些资源，比如数据库连接、文件、锁等，这些资源需要在用完之后释放掉，否则会造成内存泄漏。  

但是程序员都是人，是人就会犯错。因此经常有程序员忘记关闭这些资源。Golang直接在语言层面提供defer关键字，在打开资源语句的
下一行，就可以直接用defer语句来注册函数结束后执行关闭资源的操作。因为这样一颗“小小”的语法糖，程序员忘写关闭资源语句的情
况就大大地减少了。  

展示下我的错误的代码：  

````go
	wg.Add(len(hostMap))
	for _, item := range hostMap {
		go func(host string, port int64) {
			wg.Done()
			cli, err := h.NewSyntecAuthCli(host, port)
			if err != nil {
				wg.Done()
				errFlag = true
				h.Log.Warn(err)
				res = apierror.HandleError(err)
				return
			}

			result, err := cli.GetMachineListAll(h.Ctx, &params)
			if err != nil {
				wg.Done()
				errFlag = true
				h.Log.Warn(err)
				res = apierror.HandleError(err)
				return
			}
			wg.Done()
			output.Result = append(output.Result, result.Result...)
		}(item.Host, item.Port)
	}
	wg.Wait()
````
如果不使用defer那么就要在每个err里面加上 `wg.Done()`，当然这也是很容易忘记的。有defer
的存在就不一样了。我们只用在函数里面加上defer就好了。

````go
	wg.Add(len(hostMap))
	for _, item := range hostMap {
		go func(host string, port int64) {
			defer wg.Done()
			cli, err := h.NewSyntecAuthCli(host, port)
			if err != nil {
				errFlag = true
				h.Log.Warn(err)
				res = apierror.HandleError(err)
				return
			}

			result, err := cli.GetMachineListAll(h.Ctx, &params)
			if err != nil {
				errFlag = true
				h.Log.Warn(err)
				res = apierror.HandleError(err)
				return
			}
			output.Result = append(output.Result, result.Result...)
		}(item.Host, item.Port)
	}
	wg.Wait()
````
但是，defer并不是非常完美的，defer会有小小地延迟，对时间要求特别特别特别高的程序，可以避免使用它，其他一般忽略它带来的延迟。

当然defer也是不能滥用的，比如下面的  
````go
	i := 0
	rw.Lock()
	i = 2
	defer rw.Unlock()

	rw.Lock()
	i = 6
	defer rw.Unlock()

	fmt.Println(i)
````
defer是在函数退出的时执行的，所以第二个锁，去获取锁的时候，第一个锁还没有释放，所以就报错了。  
当然这是滥用造成的，我们应该去掉defer



### 参考
【go语言并发编程实战】   
【Golang之轻松化解defer的温柔陷阱】https://www.cnblogs.com/qcrao-2018/p/10367346.html#%E4%BB%80%E4%B9%88%E6%98%AFdefer  e