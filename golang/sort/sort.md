<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [sort 包源码解读](#sort-%E5%8C%85%E6%BA%90%E7%A0%81%E8%A7%A3%E8%AF%BB)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [如何使用](#%E5%A6%82%E4%BD%95%E4%BD%BF%E7%94%A8)
    - [基本数据类型切片的排序](#%E5%9F%BA%E6%9C%AC%E6%95%B0%E6%8D%AE%E7%B1%BB%E5%9E%8B%E5%88%87%E7%89%87%E7%9A%84%E6%8E%92%E5%BA%8F)
    - [自定义 Less 排序比较器](#%E8%87%AA%E5%AE%9A%E4%B9%89-less-%E6%8E%92%E5%BA%8F%E6%AF%94%E8%BE%83%E5%99%A8)
    - [自定义数据结构的排序](#%E8%87%AA%E5%AE%9A%E4%B9%89%E6%95%B0%E6%8D%AE%E7%BB%93%E6%9E%84%E7%9A%84%E6%8E%92%E5%BA%8F)
  - [分析下源码](#%E5%88%86%E6%9E%90%E4%B8%8B%E6%BA%90%E7%A0%81)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## sort 包源码解读

### 前言

我们的代码业务中很多地方需要我们自己进行排序操作，go 标准库中是提供了 sort 包是实现排序功能的，这里来看下生产级别的排序功能是如何实现的。   

### 如何使用

先来看下 sort 提供的主要功能  

- 对基本数据类型切片的排序支持  

- 自定义 Less 排序比较器  

- 自定义数据结构的排序  

- 判断基本数据类型切片是否已经排好序  

- 基本数据元素查找

#### 基本数据类型切片的排序

sort 包中已经实现了对 `[]int, []float, []string` 这几种类型的排序    

```go
func TestSort(t *testing.T) {
	s := []int{5, 2, 6, 3, 1, 4}
	fmt.Println("是否排好序了", sort.IntsAreSorted(s))
	sort.Ints(s)
	// 正序
	fmt.Println(s)
	// 倒序
	sort.Sort(sort.Reverse(sort.IntSlice(s)))
	fmt.Println(s)
	// 稳定排序
	sort.Stable(sort.IntSlice(s))
	fmt.Println("是否排好序了", sort.IntsAreSorted(s))
	fmt.Println("查找是否存在", sort.SearchInts(s, 5))
	fmt.Println(s)

	str := []string{"s", "f", "d", "c", "r", "a"}
	sort.Strings(str)
	fmt.Println(str)

	flo := []float64{1.33, 4.78, 0.11, 6.77, 8.99, 4.22}
	sort.Float64s(flo)
	fmt.Println(flo)
}
```

看下输出  

```
是否排好序了 false
[1 2 3 4 5 6]
[6 5 4 3 2 1]
是否排好序了 true
查找是否存在 4
[1 2 3 4 5 6]
[a c d f r s]
[0.11 1.33 4.22 4.78 6.77 8.99]
```

sort 本身不是稳定排序，需要稳定排序使用`sort.Stable`，同时排序默认是升序，降序可使用`sort.Reverse`  

#### 自定义 Less 排序比较器

如果我们需要进行的排序的内容是一些复杂的结构，例如下面的栗子，是个结构体，根据结构体中的某一个属性进行排序，这时候可以通过自定义 Less 比较器实现

使用 `sort.Slice`，`sort.Slice`中提供了 less 函数，我们，可以自定义这个函数，然后通过`sort.Slice`进行排序，`sort.Slice`不是稳定排序，稳定排序可使用`sort.SliceStable`

```go
type Person struct {
	Name string
	Age  int
}

func TestSortSlice(t *testing.T) {
	people := []Person{
		{"Bob", 31},
		{"John", 42},
		{"Michael", 17},
		{"Jenny", 26},
	}

	sort.Slice(people, func(i, j int) bool {
		return people[i].Age < people[j].Age
	})
	// Age正序
	fmt.Println(people)
	// Age倒序
	sort.Slice(people, func(i, j int) bool {
		return people[i].Age > people[j].Age
	})
	fmt.Println(people)

	// 稳定排序
	sort.SliceStable(people, func(i, j int) bool {
		return people[i].Age > people[j].Age
	})
	fmt.Println(people)
}
```

看下输出  

```
[{Michael 17} {Jenny 26} {Bob 31} {John 42}]
[{John 42} {Bob 31} {Jenny 26} {Michael 17}]
[{John 42} {Bob 31} {Jenny 26} {Michael 17}]
```

#### 自定义数据结构的排序

对自定义结构的排序，除了可以自定义 Less 排序比较器之外，sort 包中也提供了`sort.Interface`接口，我们只要实现了`sort.Interface`中提供的三个方法，即可通过 sort 包内的函数完成排序，查找等操作   

```go
// An implementation of Interface can be sorted by the routines in this package.
// The methods refer to elements of the underlying collection by integer index.
type Interface interface {
	// Len is the number of elements in the collection.
	Len() int

	// Less reports whether the element with index i
	// must sort before the element with index j.
	//
	// If both Less(i, j) and Less(j, i) are false,
	// then the elements at index i and j are considered equal.
	// Sort may place equal elements in any order in the final result,
	// while Stable preserves the original input order of equal elements.
	//
	// Less must describe a transitive ordering:
	//  - if both Less(i, j) and Less(j, k) are true, then Less(i, k) must be true as well.
	//  - if both Less(i, j) and Less(j, k) are false, then Less(i, k) must be false as well.
	//
	// Note that floating-point comparison (the < operator on float32 or float64 values)
	// is not a transitive ordering when not-a-number (NaN) values are involved.
	// See Float64Slice.Less for a correct implementation for floating-point values.
	Less(i, j int) bool

	// Swap swaps the elements with indexes i and j.
	Swap(i, j int)
}
```

来看下如何使用  

```go
type ByAge []Person

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }

func TestSortStruct(t *testing.T) {
	people := []Person{
		{"Bob", 31},
		{"John", 42},
		{"Michael", 17},
		{"Jenny", 26},
	}

	sort.Sort(ByAge(people))
	fmt.Println(people)
}
```

输出  

```
[{Michael 17} {Jenny 26} {Bob 31} {John 42}]
```

当然 sort 包中已经实现的`[]int, []float, []string` 这几种类型的排序也是实现了`sort.Interface`接口  

对于上面的三种排序，第一种和第二种基本上就能满足我们的额需求了，不过第三种灵活性更强。  

### 分析下源码

sort 中的排序算法用到了，quickSort(快排),heapSort(堆排序),insertionSort(插入排序),shellSort(希尔排序)  

先来分析下这几种排序算法的使用   

````go
func quickSort(data Interface, a, b, maxDepth int) {
	// 切片长度大于12的时候使用快排
	for b-a > 12 { // Use ShellSort for slices <= 12 elements
		// maxDepth 返回快速排序应该切换的阈值
		// 进行堆排序
		// 当 maxDepth为0的时候进行堆排序
		if maxDepth == 0 {
			heapSort(data, a, b)
			return
		}
		maxDepth--
		// doPivot 是快排核心算法，它取一点为轴，把不大于轴的元素放左边，大于轴的元素放右边，返回小于轴部分数据的最后一个下标，以及大于轴部分数据的第一个下标
		// 下标位置 a...mlo,pivot,mhi...b
		// data[a...mlo] <= data[pivot]
		// data[mhi...b] > data[pivot]
		// 和中位数一样的数据就不用在进行交换了，维护这个范围值能减少数据的次数  
		mlo, mhi := doPivot(data, a, b)
		// 避免递归过深
		// 循环是比递归节省时间的，如果有大规模的子节点，让小的先递归，达到了 maxDepth 也就是可以触发堆排序的条件了，然后使用堆排序进行排序
		if mlo-a < b-mhi {
			quickSort(data, a, mlo, maxDepth)
			a = mhi // i.e., quickSort(data, mhi, b)
		} else {
			quickSort(data, mhi, b, maxDepth)
			b = mlo // i.e., quickSort(data, a, mlo)
		}
	}
	// 如果切片的长度大于1小于等于12的时候，使用 shell 排序  
	if b-a > 1 {
		// Do ShellSort pass with gap 6
		// It could be written in this simplified form cause b-a <= 12
		// 这里先做一轮shell 排序
		for i := a + 6; i < b; i++ {
			if data.Less(i, i-6) {
				data.Swap(i, i-6)
			}
		}
		// 进行插入排序
		insertionSort(data, a, b)
	}
}

// maxDepth 返回快速排序应该切换的阈值
// 进行堆排序
func maxDepth(n int) int {
	var depth int
	for i := n; i > 0; i >>= 1 {
		depth++
	}
	return depth * 2
}

// doPivot 是快排核心算法，它取一点为轴，把不大于轴的元素放左边，大于轴的元素放右边，返回小于轴部分数据的最后一个下标，以及大于轴部分数据的第一个下标
// 下标位置 lo...midlo,pivot,midhi...hi
// data[lo...midlo] <= data[pivot]
// data[midhi...hi] > data[pivot]
func doPivot(data Interface, lo, hi int) (midlo, midhi int) {
	m := int(uint(lo+hi) >> 1) // Written like this to avoid integer overflow.
	// 这里用到了 Tukey's ninther 算法，文章链接 https://www.johndcook.com/blog/2009/06/23/tukey-median-ninther/
	// 通过该算法求出中位数
	if hi-lo > 40 {
		// Tukey's ``Ninther,'' median of three medians of three.
		s := (hi - lo) / 8
		medianOfThree(data, lo, lo+s, lo+2*s)
		medianOfThree(data, m, m-s, m+s)
		medianOfThree(data, hi-1, hi-1-s, hi-1-2*s)
	}

	// 求出中位数 data[m] <= data[lo] <= data[hi-1]
	medianOfThree(data, lo, m, hi-1)

	// Invariants are:
	//	data[lo] = pivot (set up by ChoosePivot)
	//	data[lo < i < a] < pivot
	//	data[a <= i < b] <= pivot
	//	data[b <= i < c] unexamined
	//	data[c <= i < hi-1] > pivot
	//	data[hi-1] >= pivot
	// 中位数
	pivot := lo
	a, c := lo+1, hi-1

	// 将左边和中位数进行比较， 一直到不满足条件为止
	for ; a < c && data.Less(a, pivot); a++ {
	}
	b := a
	for {
		// 对比保证  data[b] <= pivot
		// 和上面的有点重合，不过在处理上面的for 循环，会发生数据交换的情况，可能是为了更加严谨吧
		for ; b < c && !data.Less(pivot, b); b++ {
		}
		// 对比保证  data[c-1] > pivot
		for ; b < c && data.Less(pivot, c-1); c-- { // data[c-1] > pivot
		}
		// 左边和右边重合或者已经再右边的右侧
		if b >= c {
			break
		}
		// data[b] > pivot; data[c-1] <= pivot
		// 左侧的数据大于右侧，交换，然后接着排序
		data.Swap(b, c-1)
		b++
		c--
	}
	// 如果 hi-c<3 则存在重复项（按中位数为 9 的属性）。
	// 让我们稍微保守一点，将边框设置为 5。
	protect := hi-c < 5
	if !protect && hi-c < (hi-lo)/4 {
		// 用一些特殊的点和中间数进行比较
		dups := 0
		// 处理使 data[hi-1] = pivot
		if !data.Less(pivot, hi-1) {
			data.Swap(c, hi-1)
			c++
			dups++
		}
		// 处理使 data[b-1] = pivot
		if !data.Less(b-1, pivot) {
			b--
			dups++
		}
		// m-lo = (hi-lo)/2 > 6
		// b-lo > (hi-lo)*3/4-1 > 8
		// ==> m < b ==> data[m] <= pivot
		if !data.Less(m, pivot) { // data[m] = pivot
			data.Swap(m, b-1)
			b--
			dups++
		}
		// 如果上面的 if 进入了两次， 就证明现在是偏态分布（也就是左右不平衡的）
		protect = dups > 1
	}
	// 不平衡，接着进行处理
	if protect {
		// Protect against a lot of duplicates
		// Add invariant:
		//	data[a <= i < b] unexamined
		//	data[b <= i < c] = pivot
		for {
			// 处理使 data[b] == pivot
			for ; a < b && !data.Less(b-1, pivot); b-- {
			}
			// 处理使 data[a] < pivot
			for ; a < b && data.Less(a, pivot); a++ {
			}
			if a >= b {
				break
			}
			// data[a] == pivot; data[b-1] < pivot
			data.Swap(a, b-1)
			a++
			b--
		}
	}
	// 交换中位数到中间
	data.Swap(pivot, b-1)
	return b - 1, c
}
````

对于这几种排序算法的使用，sort 包中是混合使用的  

1、如果切片长度大于12的时候使用快排，使用快排的时候，如果满足了使用堆排序的条件没这个排序对于后面的数据的处理，又会转换成堆排序；   

2、切片长度小于12了，就使用 shell 排序，shell 排序只处理一轮数据，后面数据的排序使用插入排序；  

堆排序和插入排序就是正常的排序处理了    

```go
// insertionSort sorts data[a:b] using insertion sort.
// 插入排序
func insertionSort(data Interface, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && data.Less(j, j-1); j-- {
			data.Swap(j, j-1)
		}
	}
}

// 堆排序
func heapSort(data Interface, a, b int) {
	first := a
	lo := 0
	hi := b - a

	// Build heap with greatest element at top.
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDown(data, i, hi, first)
	}

	// Pop elements, largest first, into end of data.
	for i := hi - 1; i >= 0; i-- {
		data.Swap(first, first+i)
		siftDown(data, lo, i, first)
	}
}
```





### 参考

【Golang sort 排序】https://blog.csdn.net/K346K346/article/details/118314382    
【文中示例代码】https://github.com/boilingfrog/Go-POINT/blob/master/golang/sort/sort_test.go  
【】https://www.johndcook.com/blog/2009/06/23/tukey-median-ninther/  

