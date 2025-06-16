package main

import (
	"errors"
	"fmt"
	"strconv"
)

func main() {
	slice1 := []int{1, 2, 3, 4, 5}
	slice2 := []int{5, 4, 3}
	fmt.Printf("%p\n", slice1)
	copy(slice1, slice2)

	slice2[0] = 9

	slice1[3] = 9

	fmt.Printf("%p\n", slice1)
}

func DB(id string) (int, error) {
	if len(id) != 24 {
		return 0, errors.New("出错了")
	}
	key, err := strconv.ParseInt(id[len(id)-3:], 16, 0)
	if err != nil {
		return 0, errors.New("出错了")
	}
	i := int(key) % 28
	return i + 1, nil
}
