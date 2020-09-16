## 设计模式

###装饰模式

装饰者模式(Decorator Pattern)：  

动态地给一个对象增加一些额外的职责，增加对象功能来说，装饰模式比生成子类实现更为灵活。装饰模式是一种对象结构型模式。  

在装饰者模式中，为了让系统具有更好的灵活性和可扩展性，我们通常会定义一个抽象装饰类，而将具体的装饰类作为它的子类。


### 单例模式

单例模式(Singleton Pattern)：确保某一个类只有一个实例，而且自行实例化并向整个系统提供这个实例，这个类称为单例类，它提供全局访问的方法。单例模式是一种对象创建型模式。  

单例模式有三个要点：  

1、构造方法私有化；  
2、实例化的变量引用私有化；  
3、获取实例的方法共有；  

懒汉模式(需要的时候才去加载)

当程序需要这个实例的时候才去创建对象，就如同一个人懒的饿到不行了才去吃东西。  
如果一个对象使用频率不高，占用内存还特别大，明显就不合适用饿汉式了，这时就需要一种懒加载的思想。  

```go
package singleton
import "sync"
type singleton1 struct {  // 首字母小写，确保外部无法直接实例化
}
var once sync.Once
var instance *singleton1
func GetInstance() *singleton1 {
	once.Do(func() {       // 确保只执行一次，且执行时加锁
		instance = &singleton1{}
	})
	
	return instance
}
```


饿汉模式（不管时候需要都提前加载好）

不管程序是否需要这个对象的实例，总是在类加载的时候就先创建好实例，理解起来就像不管一个人想不想吃东西都把吃的先买好，如同饿怕了一样。  

```go
package singleton
import "sync"
type singleton2 struct {  // 首字母小写，确保外部无法直接实例化
	id string
}
var once sync.Once
var instance *singleton2
func init() {     // 加载时执行
	once.Do(func() {  // 确保只执行一次，且执行时加锁
		instance = &singleton2{}
	})
}
func GetInstance() *singleton2 {
	return instance
}

```

### 适配器模式

该模式适用于，被调用者的接口已经定型的情况下（如：已经在运行的服务），而调用者定义的接口又不兼容被调用者提供的接口，这时可以利用一个适配器类提供接口转换功能。  

