<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [sort 包源码解读](#sort-%E5%8C%85%E6%BA%90%E7%A0%81%E8%A7%A3%E8%AF%BB)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [如何使用](#%E5%A6%82%E4%BD%95%E4%BD%BF%E7%94%A8)
    - [基本数据类型切片的排序](#%E5%9F%BA%E6%9C%AC%E6%95%B0%E6%8D%AE%E7%B1%BB%E5%9E%8B%E5%88%87%E7%89%87%E7%9A%84%E6%8E%92%E5%BA%8F)
    - [自定义 Less 排序比较器](#%E8%87%AA%E5%AE%9A%E4%B9%89-less-%E6%8E%92%E5%BA%8F%E6%AF%94%E8%BE%83%E5%99%A8)
    - [自定义数据结构的排序](#%E8%87%AA%E5%AE%9A%E4%B9%89%E6%95%B0%E6%8D%AE%E7%BB%93%E6%9E%84%E7%9A%84%E6%8E%92%E5%BA%8F)
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
	sort.Ints(s)
	// 正序
	fmt.Println(s)
	// 倒序
	sort.Sort(sort.Reverse(sort.IntSlice(s)))
	fmt.Println(s)
	// 稳定排序
	sort.Stable(sort.IntSlice(s))
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
[1 2 3 4 5 6]
[6 5 4 3 2 1]
[1 2 3 4 5 6]
[a c d f r s]
[0.11 1.33 4.22 4.78 6.77 8.99]
```

sort 本身不是稳定排序，需要稳定排序使用`sort.Stable`，同时排序默认是升序，降序可使用`sort.Reverse`  

#### 自定义 Less 排序比较器

使用 `sort.Slice`，`sort.Slice`中提供了 less 函数，我们，可以自定义这个函数，然后通过`sort.Slice`进行排序，`sort.Slice`不是稳定排序，稳定排序可使用`sort.SliceStable`

```go
func TestSortSlice(t *testing.T) {
	s := []int{5, 2, 6, 3, 1, 4}
	sort.Slice(s, func(i, j int) bool {
		return s[j] > s[i]
	})
	// 正序
	fmt.Println(s)
	// 倒序
	sort.Slice(s, func(i, j int) bool {
		return s[j] < s[i]
	})
	fmt.Println(s)

	// 稳定排序
	sort.SliceStable(s, func(i, j int) bool {
		return s[j] < s[i]
	})
	fmt.Println(s)
}
```

看下输出  

```
[1 2 3 4 5 6]
[6 5 4 3 2 1]
[6 5 4 3 2 1]
```

#### 自定义数据结构的排序







### 参考

【Golang sort 排序】https://blog.csdn.net/K346K346/article/details/118314382    

