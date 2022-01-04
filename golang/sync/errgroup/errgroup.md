<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [errgroup](#errgroup)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [如何使用](#%E5%A6%82%E4%BD%95%E4%BD%BF%E7%94%A8)
  - [实现原理](#%E5%AE%9E%E7%8E%B0%E5%8E%9F%E7%90%86)
  - [WithContext](#withcontext)
  - [Go](#go)
  - [Wait](#wait)
  - [错误的使用](#%E9%94%99%E8%AF%AF%E7%9A%84%E4%BD%BF%E7%94%A8)
  - [总结](#%E6%80%BB%E7%BB%93)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## errgroup

### 前言

来看下errgroup的实现  

### 如何使用

```go
func main() {
	var eg errgroup.Group

	eg.Go(func() error {
		return errors.New("test1")
	})

	eg.Go(func() error {
		return errors.New("test2")
	})

	if err := eg.Wait(); err != nil {
		fmt.Println(err)
	}
}
```

类比于`waitgroup`,`errgroup`增加了一个对`goroutine`错误收集的作用。  

不过需要注意的是：  

`errgroup`返回的第一个出错的`goroutine`抛出的`err`。  

`errgroup`中还可以加入`context`  

```go
func main() {
	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		// test1函数还可以在启动很多goroutine
		// 子节点都传入ctx，当test1报错，会把test1的子节点一一cancel
		return test1(ctx)
	})

	eg.Go(func() error {
		return test1(ctx)
	})

	if err := eg.Wait(); err != nil {
		fmt.Println(err)
	}
}

func test1(ctx context.Context) error {
	return errors.New("test2")
}
```

### 实现原理

代码很简单  

```go
type Group struct {
	// 一个取消的函数，主要来包装context.WithCancel的CancelFunc
	cancel func()

	// 还是借助于WaitGroup实现的
	wg sync.WaitGroup

	// 使用sync.Once实现只输出第一个err
	errOnce sync.Once

	// 记录下错误的信息
	err     error
}
```

还是在WaitGroup的基础上实现的  

### WithContext

```go
// 返回一个被context.WithCancel重新包装的ctx

func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &Group{cancel: cancel}, ctx
}
```

里面使用了`context`,通过`context.WithCancel`对传入的context进行了包装  

当`WithCancel`函数返回的`CancelFunc`被调用或者是父节点的`done channel`被关闭（父节点的 CancelFunc 被调用），此 context（子节点）的 `done channel` 也会被关闭。  

`errgroup`把返回的`CancelFunc`包进了自己的`cancel`中，来实现对使用`errgroup`的`ctx`启动的`goroutine`的取消操作。  

### Go

```go
// 启动取消阻塞的goroutine
// 记录第一个出错的goroutine的err信息
func (g *Group) Go(f func() error) {
	// 借助于waitgroup实现
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		// 执行出错
		if err := f(); err != nil {
			// 通过sync.Once记录下第一个出错的err信息
			g.errOnce.Do(func() {
				g.err = err
				// 如果包装了cancel，也就是context的CancelFunc，执行退出操作
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
	}()
}
```

1、借助于`waitgroup`实现对`goroutine`阻塞；  

2、通过`sync.Once`记录下，第一个出错的`goroutine`的错误信息；  

3、如果包装了`context`的`CancelFunc`，在出错的时候进行退出操作。  

### Wait

```go
// 阻塞所有的通过Go加入的goroutine，然后等待他们一个个执行完成
// 然后返回第一个出错的goroutine的错误信息
func (g *Group) Wait() error {
	// 借助于waitgroup实现
	g.wg.Wait()
	// 如果包装了cancel，也就是context的CancelFunc，执行退出操作
	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}
```

1、借助于`waitgroup`实现对`goroutine`阻塞；  

2、如果包装了`context`的`CancelFunc`，在出错的时候进行退出操作；  

3、抛出第一个出错的`goroutine`的错误信息。  

### 错误的使用

不过工作中发现一个`errgroup`错误使用的例子  

```go
func main() {
	eg := errgroup.Group{}
	var err error
	eg.Go(func() error {
		// 处理业务
		err = test1()
		return err
	})

	eg.Go(func() error {
		// 处理业务
		err = test1()
		return err
	})

	if err = eg.Wait(); err != nil {
		fmt.Println(err)
	}
}

func test1() error {
	return errors.New("test2")
}
```

很明显err被资源竞争了  

```go
$ go run -race main.go 
==================
WARNING: DATA RACE
Write at 0x00c0000801f0 by goroutine 8:
  main.main.func2()
      /Users/yj/Go/src/Go-POINT/sync/errgroup/main.go:23 +0x97
...
```

### 总结

`errgroup`相比比较简单，不过需要先弄明白`waitgroup`,`context`以及`sync.Once`,主要是借助这几个组件来实现的。  

`errgroup`可以带携带`context`,如果包装了`context`，会使用`context.WithCancel`进行超时，取消或者一些异常的情况  
