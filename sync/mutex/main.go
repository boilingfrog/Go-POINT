package main

import (
	"fmt"
	"sync"
)

var (
	mu      sync.Mutex
	balance int
)

const (
	// mutex is locked
	// 是否加锁的标识
	mutexLocked = 1 << iota
	mutexWoken
	mutexStarving
	mutexWaiterShift = iota
)

func main() {
	a, b := 60, 13

	fmt.Println(a & b)

}

func Deposit(amount int) {
	mu.Lock()
	balance = balance + amount
	mu.Unlock()
}

func Balance() int {
	mu.Lock()
	b := balance
	mu.Unlock()
	return b
}
