<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [切片](#%E5%88%87%E7%89%87)
  - [什么是slice](#%E4%BB%80%E4%B9%88%E6%98%AFslice)
    - [slice的创建使用](#slice%E7%9A%84%E5%88%9B%E5%BB%BA%E4%BD%BF%E7%94%A8)
    - [slice使用的一点规范](#slice%E4%BD%BF%E7%94%A8%E7%9A%84%E4%B8%80%E7%82%B9%E8%A7%84%E8%8C%83)
  - [slice和数组的区别](#slice%E5%92%8C%E6%95%B0%E7%BB%84%E7%9A%84%E5%8C%BA%E5%88%AB)
  - [slice的append是如何发生的](#slice%E7%9A%84append%E6%98%AF%E5%A6%82%E4%BD%95%E5%8F%91%E7%94%9F%E7%9A%84)
  - [复制Slice和Map注意事项](#%E5%A4%8D%E5%88%B6slice%E5%92%8Cmap%E6%B3%A8%E6%84%8F%E4%BA%8B%E9%A1%B9)
    - [接收 Slice 和 Map 作为入参](#%E6%8E%A5%E6%94%B6-slice-%E5%92%8C-map-%E4%BD%9C%E4%B8%BA%E5%85%A5%E5%8F%82)
    - [返回 Slice 和 Map](#%E8%BF%94%E5%9B%9E-slice-%E5%92%8C-map)
  - [切片的截取](#%E5%88%87%E7%89%87%E7%9A%84%E6%88%AA%E5%8F%96)
    - [不发生扩容情况下修改新切片](#%E4%B8%8D%E5%8F%91%E7%94%9F%E6%89%A9%E5%AE%B9%E6%83%85%E5%86%B5%E4%B8%8B%E4%BF%AE%E6%94%B9%E6%96%B0%E5%88%87%E7%89%87)
    - [发生扩容情况下修改新的切片](#%E5%8F%91%E7%94%9F%E6%89%A9%E5%AE%B9%E6%83%85%E5%86%B5%E4%B8%8B%E4%BF%AE%E6%94%B9%E6%96%B0%E7%9A%84%E5%88%87%E7%89%87)
  - [总结](#%E6%80%BB%E7%BB%93)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 切片

### 什么是slice

Go中的切片，是我们经常用到的数据结构。有着比数组更灵活的用法，那么作者就去探究下什么是切片。

我们先来了解下切片的数据结构
````
type slice struct {
    array unsafe.Pointer // 指针
    len   int // 长度
    cap   int // 容量
}
````

切片一共三个属性：指针，指向底层的数组；长度，表示切片可用元素的个数，也就是说使用下标
对元素进行访问的时候，下标不能超过的长度；容量，底层数组的元素个数，容量》=长度。

![Aaron Swartz](https://github.com/zhan-liz/Go-POINT/blob/master/img/golang/slice_1.png?raw=true)

底层的数组是可以被多个切片同时指向的，因此对一个切片元素的操作可能会影响到其他的切片。

#### slice的创建使用

|     序号      | 方式 |             代码示例                                          |
| ------------ | ------------ | -------------------------------                      |
| 1            |    直接声明   |   var slice []int                                    |
| 2            |    new       |   slice := *new([]int)                               |
| 3            |    字面量     |   slice := []int{1,2,3,4,5}                          |
| 4            |    make      |   slice := make([]int, 5, 10)                        |
| 5            |从切片或数组截取|   slice := array[1:5] 或 slice := sourceSlice[1:5]   |

第一种创建出来的 slice 其实是一个 nil slice。它的长度和容量都为0。和nil比较的结果为true。

这里比较混淆的是empty slice，它的长度和容量也都为0，但是所有的空切片的数据指针都指向同一个地址 0xc42003bda0。空切片和 nil 比较的结果为false。

下面是它的内部结构：
![Aaron Swartz](https://github.com/zhan-liz/Go-POINT/blob/master/img/golang/slice_2.png?raw=true)

|     创建方式  |    nil切片               |             空切片                |
| ------------ | ------------------      | -------------------------------   |
| 方式一        |    var s1 []int         |  var s2 = []int{}                 |
| 方式二        |    var s4 = *new([]int) |  var s3 = make([]int, 0)          |
| 长度          |    0                    |     0                             |
| 容量          |    0                    |     0                             |
| 和nil比较      |true                     |   false                          |

nil 切片和空切片很相似，长度和容量都是0，官方建议尽量使用 nil 切片。

- 字面量

直接初始化表达式进行创建

````
 s1 := []int{0, 1, 2, 3,5,6}
````

- make 
````
slice := make([]int, 5, 10) // 长度为5，容量为10
````

#### slice使用的一点规范

- 根据 Uber Go代码风格指南

- nil 是一个有效的 slice

`nil` 是一个长度为 0 的 slice。意思是，

- 使用 `nil` 来替代长度为 0 的 slice 返回

  <table>
  <thead><tr><th>Bad</th><th>Good</th></tr></thead>
  <tbody>
  <tr><td>

  ```go
  if x == "" {
    return []int{}
  }
  ```

  </td><td>

  ```go
  if x == "" {
    return nil
  }
  ```

  </td></tr>
  </tbody></table>

- 检查一个空 slice，应该使用 `len(s) == 0`，而不是 `nil`。

  <table>
  <thead><tr><th>Bad</th><th>Good</th></tr></thead>
  <tbody>
  <tr><td>

  ```go
  func isEmpty(s []string) bool {
    return s == nil
  }
  ```

  </td><td>

  ```go
  func isEmpty(s []string) bool {
    return len(s) == 0
  }
  ```

  </td></tr>
  </tbody></table>

- The zero value (a slice declared with `var`) is usable immediately without
  `make()`.

- 零值（通过 `var` 声明的 slice）是立马可用的，并不需要 `make()` 。

  <table>
  <thead><tr><th>Bad</th><th>Good</th></tr></thead>
  <tbody>
  <tr><td>

  ```go
  nums := []int{}
  // or, nums := make([]int)

  if add1 {
    nums = append(nums, 1)
  }

  if add2 {
    nums = append(nums, 2)
  }
  ```

  </td><td>

  ```go
  var nums []int

  if add1 {
    nums = append(nums, 1)
  }

  if add2 {
    nums = append(nums, 2)
  }
  ```

  </td></tr>
  </tbody></table>
  
### slice和数组的区别

slice的底层是数组，slice是对数组的封装，它描述一个数组的片段。两者都可以通过下标访问单个元素。

数组是定长的，长度定义好，不能改变。在Go中数组是不常见的，因为长度是类型的一部分，限制了它的表达
能力，比如[3]int 和 [4]int 就是不同的类型。

切片可以动态的扩容，非常灵活。切片的类型和长度没有关系。

### slice的append是如何发生的

先看看append函数的原型：

````
func append(slice []Type, elems ...Type) []Type
````
append 函数的参数长度可变，因此可以追加多个值到 slice 中，还可以用 ... 传入 slice，直接追加一个切片。

````
slice = append(slice, elem1, elem2)
slice = append(slice, anotherSlice...)
````
append函数返回值是一个新的slice，Go编译器不允许调用了append函数后不使用返回值。

````
append(slice, elem1, elem2)
append(slice, anotherSlice...)
````
上面是不能编译通过的

使用 append 可以向 slice 追加元素，实际上是往底层数组添加元素。但是底层数组的长度是固定的，如果索引 len-1 所指向的元素已经是底层数组的最后一个元素，就没法再添加了。

这时，slice 会迁移到新的内存位置，新底层数组的长度也会增加，这样就可以放置新增的元素。同时，为了应对未来可能再次发生的 append 操作，新的底层数组的长度，也就是新 slice 的容量是留了一定的 buffer 的。否则，每次添加元素的时候，都会发生迁移，成本太高。

新slice预留buffer大小是有一定规律的。
````
// growslice handles slice growth during append.
// It is passed the slice element type, the old slice, and the desired new minimum capacity,
// and it returns a new slice with at least that capacity, with the old data
// copied into it.
// The new slice's length is set to the old slice's length,
// NOT to the new requested capacity.
// This is for codegen convenience. The old slice's length is used immediately
// to calculate where to write new values during an append.
// TODO: When the old backend is gone, reconsider this decision.
// The SSA backend might prefer the new length or to return only ptr/cap and save stack space.
func growslice(et *_type, old slice, cap int) slice {
	if raceenabled {
		callerpc := getcallerpc()
		racereadrangepc(old.array, uintptr(old.len*int(et.size)), callerpc, funcPC(growslice))
	}
	if msanenabled {
		msanread(old.array, uintptr(old.len*int(et.size)))
	}

	if et.size == 0 {
		if cap < old.cap {
			panic(errorString("growslice: cap out of range"))
		}
		// append should not create a slice with nil pointer but non-zero len.
		// We assume that append doesn't need to preserve old.array in this case.
		return slice{unsafe.Pointer(&zerobase), old.len, cap}
	}

	newcap := old.cap
	doublecap := newcap + newcap
	if cap > doublecap {
		newcap = cap
	} else {
		if old.len < 1024 {
			newcap = doublecap
		} else {
			// Check 0 < newcap to detect overflow
			// and prevent an infinite loop.
			for 0 < newcap && newcap < cap {
				newcap += newcap / 4
			}
			// Set newcap to the requested cap when
			// the newcap calculation overflowed.
			if newcap <= 0 {
				newcap = cap
			}
		}
	}

	var overflow bool
	var lenmem, newlenmem, capmem uintptr
	const ptrSize = unsafe.Sizeof((*byte)(nil))
	switch et.size {
	case 1:
		lenmem = uintptr(old.len)
		newlenmem = uintptr(cap)
		capmem = roundupsize(uintptr(newcap))
		overflow = uintptr(newcap) > _MaxMem
		newcap = int(capmem)
	case ptrSize:
		lenmem = uintptr(old.len) * ptrSize
		newlenmem = uintptr(cap) * ptrSize
		capmem = roundupsize(uintptr(newcap) * ptrSize)
		overflow = uintptr(newcap) > _MaxMem/ptrSize
		newcap = int(capmem / ptrSize)
	default:
		lenmem = uintptr(old.len) * et.size
		newlenmem = uintptr(cap) * et.size
		capmem = roundupsize(uintptr(newcap) * et.size)
		overflow = uintptr(newcap) > maxSliceCap(et.size)
		newcap = int(capmem / et.size)
	}

	// The check of overflow (uintptr(newcap) > maxSliceCap(et.size))
	// in addition to capmem > _MaxMem is needed to prevent an overflow
	// which can be used to trigger a segfault on 32bit architectures
	// with this example program:
	//
	// type T [1<<27 + 1]int64
	//
	// var d T
	// var s []T
	//
	// func main() {
	//   s = append(s, d, d, d, d)
	//   print(len(s), "\n")
	// }
	if cap < old.cap || overflow || capmem > _MaxMem {
		panic(errorString("growslice: cap out of range"))
	}

	var p unsafe.Pointer
	if et.kind&kindNoPointers != 0 {
		p = mallocgc(capmem, nil, false)
		memmove(p, old.array, lenmem)
		// The append() that calls growslice is going to overwrite from old.len to cap (which will be the new length).
		// Only clear the part that will not be overwritten.
		memclrNoHeapPointers(add(p, newlenmem), capmem-newlenmem)
	} else {
		// Note: can't use rawmem (which avoids zeroing of memory), because then GC can scan uninitialized memory.
		p = mallocgc(capmem, et, true)
		if !writeBarrier.enabled {
			memmove(p, old.array, lenmem)
		} else {
			for i := uintptr(0); i < lenmem; i += et.size {
				typedmemmove(et, add(p, i), add(old.array, i))
			}
		}
	}

	return slice{p, old.len, newcap}
}
````

其中这一段是重点的代码，我们可以看到
````
newcap := old.cap
	doublecap := newcap + newcap
	if cap > doublecap {
		newcap = cap
	} else {
		if old.len < 1024 {
			newcap = doublecap
		} else {
			// Check 0 < newcap to detect overflow
			// and prevent an infinite loop.
			for 0 < newcap && newcap < cap {
				newcap += newcap / 4
			}
			// Set newcap to the requested cap when
			// the newcap calculation overflowed.
			if newcap <= 0 {
				newcap = cap
			}
		}
	}
````

### 复制Slice和Map注意事项
slice 和 map 包含指向底层数据的指针，因此复制的时候需要当心。

slice 和 map 包含指向底层数据的指针，因此复制的时候需要当心。

#### 接收 Slice 和 Map 作为入参

需要留意的是，如果你保存了作为参数接收的 map 或 slice 的引用，可以通过引用修改它。

<table>
<thead><tr><th>Bad</th> <th>Good</th></tr></thead>
<tbody>
<tr>
<td>

```go
func (d *Driver) SetTrips(trips []Trip) {
  d.trips = trips
}

trips := ...
d1.SetTrips(trips)

// Did you mean to modify d1.trips?
trips[0] = ...
```

</td>
<td>

```go
func (d *Driver) SetTrips(trips []Trip) {
  d.trips = make([]Trip, len(trips))
  copy(d.trips, trips)
}

trips := ...
d1.SetTrips(trips)

// We can now modify trips[0] without affecting d1.trips.
trips[0] = ...
```

</td>
</tr>

</tbody>
</table>

#### 返回 Slice 和 Map

类似的，当心 map 或者 slice 暴露的内部状态是可以被修改的。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
type Stats struct {
  sync.Mutex

  counters map[string]int
}

// Snapshot 方法返回当前的状态
func (s *Stats) Snapshot() map[string]int {
  s.Lock()
  defer s.Unlock()

  return s.counters
}

// snapshot 不再被锁保护
snapshot := stats.Snapshot()
```

</td><td>

```go
type Stats struct {
  sync.Mutex

  counters map[string]int
}

func (s *Stats) Snapshot() map[string]int {
  s.Lock()
  defer s.Unlock()

  result := make(map[string]int, len(s.counters))
  for k, v := range s.counters {
    result[k] = v
  }
  return result
}

// 现在 Snapshot 是一个副本
snapshot := stats.Snapshot()
```

</td></tr>
</tbody></table>

### 切片的截取

基于已有 slice 创建新 slice 对象，被称为 reslice。新 slice 和老 slice 共用底层数组，新老 slice 对底层数组的更改都会影响到彼此。基于数
组创建的新 slice 对象也是同样的效果：对数组或 slice 元素作的更改都会影响到彼此。  

如果新截取的切片发生了扩容了，就会重新申请新的内存空间，这样就新截取的切片就不指向原来的地址空间了。  

截取的例子：

```go
 data := [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
 slice := data[2:4:6] // data[low, high, max]

max >= high >= low
```

或者

```go
 data := [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
 slice := data[2:4] // data[low, high]

```

这两个`low, high`对应的都是前闭后开。  

区别就是如果不设置`max`,那么新切片的长度cap就是原切片的长度`len(data)-low`  
设置了`max`那么新切片的长度cap就是原切片的长度`max-low`  

```go
package main

import "fmt"

func main() {
	s1 := []int{2, 3, 6, 2, 4, 5, 6, 7}
	fmt.Println(cap(s1), len(s1))

	s2 := s1[2:3:4]
	fmt.Println("带上max")
	fmt.Println(s2)
	fmt.Println(cap(s2), len(s2))

	fmt.Println("不带max")
	s3 := s1[2:3]
	fmt.Println(s3)
	fmt.Println(cap(s3), len(s3))
}
```

输出

```go
8 8
带上max
[6]
2 1
不带max
[6]
6 1
```

#### 不发生扩容情况下修改新切片

截取了新的切片，当不发生扩容的情况下，操作新的切片是会对老切片的数据产生影响，因为他们指向的是通一地址空间。  

```go
package main

import "fmt"

func main() {
	s1 := []int{2, 3, 6, 2, 4, 5, 6, 7}
	fmt.Println("原切片")
	fmt.Println("cap", cap(s1), "len", len(s1))

	s2 := s1[6:7]
	fmt.Println("新切片")
	fmt.Println("cap", cap(s2), "len", len(s2))
	fmt.Println(s2)

	s2 = append(s2, 100)
	fmt.Println("append之后的新切片")
	fmt.Println("cap", cap(s2), "len", len(s2))
	fmt.Println(cap(s2), len(s2))
	fmt.Println(s2)

	fmt.Println("老切片")
	fmt.Println(s1)
}
```

打印下结果

```go
原切片
cap 8 len 8
新切片
cap 2 len 1
[6]
append之后的新切片
cap 2 len 2
2 2
[6 100]
老切片
[2 3 6 2 4 5 6 100]
```

#### 发生扩容情况下修改新的切片

截取的新的切片发生扩容的情况下，新切片将指向一个新的数据空间。  

```go
package main

import "fmt"

func main() {
	s1 := []int{2, 3, 6, 2, 4, 5, 6, 7}
	fmt.Println("原切片")
	fmt.Println("cap", cap(s1), "len", len(s1))

	s2 := s1[6:7]
	fmt.Println("新切片")
	fmt.Println("cap", cap(s2), "len", len(s2))
	fmt.Println(s2)

	s2 = append(s2, 100)
	fmt.Println("append之后的新切片")
	fmt.Println("cap", cap(s2), "len", len(s2))
	fmt.Println(cap(s2), len(s2))
	fmt.Println(s2)

	fmt.Println("老切片")
	fmt.Println(s1)

	s2 = append(s2, 888)
	fmt.Println("append之后的新切片，发生扩容")
	fmt.Println("cap", cap(s2), "len", len(s2))
	fmt.Println(cap(s2), len(s2))
	fmt.Println(s2)

	fmt.Println("老切片")
	fmt.Println(s1)

}
```

打印下结果

```go
cap 8 len 8
新切片
cap 2 len 1
[6]
append之后的新切片
cap 2 len 2
2 2
[6 100]
老切片
[2 3 6 2 4 5 6 100]
append之后的新切片，发生扩容
cap 4 len 3
4 3
[6 100 888]
老切片
[2 3 6 2 4 5 6 100]
```

### 总结
- 切片是对底层数组的一个抽象，描述了它的一个片段。
- 切片实际上是一个结构体，它有三个字段：长度，容量，底层数据的地址。
- 多个切片可能共享同一个底层数组，这种情况下，对其中一个切片或者底层数组的更改，会影响到其他切片。
- append 函数会在切片容量不够的情况下，调用 growslice 函数获取所需要的内存，这称为扩容，扩容会改变元素原来的位置。
- 扩容策略并不是简单的扩为原切片容量的 2 倍或 1.25 倍，还有内存对齐的操作。扩容后的容量 >= 原容量的 2 倍或 1.25 倍。
- 当直接用切片作为函数参数时，可以改变切片的元素，不能改变切片本身；想要改变切片本身，可以将改变后的切片返回，函数调用者接收改变后的切片或者将切片指针作为函数参数。



### 参考
- 【快速理解Go数组和切片的内部实现原理】 https://i6448038.github.io/2018/08/11/array-and-slice-principle/
- 【GO代码风格指南 Uber Go】 https://github.com/uber-go/guide
- 【深度解密Go语言之Slice】 https://mp.weixin.qq.com/s/MTZ0C9zYsNrb8wyIm2D8BA    