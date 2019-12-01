# unsafe


## 指针类型
首先我们先来了解下，GO里面的指针类型。

为什么需要指针类型呢？参考文献 go101.org 里举了这样一个例子：

````
func double(x int) {
	fmt.Println(x)
	x += x
	fmt.Println(x)
}

func main() {
	var a = 3
	double(a)
	fmt.Println(a)

}
````

double函数的作用是将3翻倍，但是实际上却没有做到，为什么呢？
因为go语言的函数操作都是值传递。double函数里面的x只是a的一个拷贝，
在函数内部对x的操作不能反馈到实参a。

其实在实际的编写代码的过程中我们会使用一个指针进行解决。

````

func double1(x *int) {
	*x += *x
	x = nil
}

func main() {
	var a = 3
	double1(&a)

	fmt.Println(a)

	p := &a

	double1(p)

	fmt.Println(*p)

}

````

其中有一个操作
````
x=nil
````

这个操作没有对我们的结果产生丝毫的影响。
其实也是很好理解的，因为我们知道go里面的函数中使用的都是值传递
x=nil，只是对&a的一个拷贝。

我们知道slice 和 map 包含指向底层数据的指针。
我们对它们的操作是会影响到，原参数的值。

````
func change(sl []int64) {
	sl[0] = 2
}

func main() {

	var sl = make([]int64, 2)
	change(sl)
	fmt.Println(sl)  // [2 0]
}
````
我们而已看到输出的值已经是[2 0]

这时候我们可以使用一个copy来操作
````
func change(sl []int64) {
	sl[0] = 2
}

func changeNo(sl []int64) {
	s2 := make([]int64, 2)
	copy(sl, s2)
	s2[0] = 2
}

func main() {

	var sl = make([]int64, 2)
	change(sl)
	fmt.Println(sl)

	changeNo(sl)
	fmt.Println(sl)
}
````

限制一：GO里面的指针不能进行数学的运算

````
错误
a := 5
p := &a

p++
p = &a + 3
````
限制二：不同类型的指针不能互相转换

````
错误的
func main(){
   a:=int(100)
   var f *float64

    f=&a
}
````
限制三：不同类型的指针不能使用==或!=比较。

限制四：不能类型的指针变量不能相互赋值。

## 什么是 unsafe














### 参考
- 【深度解密Go语言之Slice】 https://mp.weixin.qq.com/s/MTZ0C9zYsNrb8wyIm2D8BA    

















































