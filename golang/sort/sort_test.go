package sort

import (
	"fmt"
	"sort"
	"testing"
)

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
