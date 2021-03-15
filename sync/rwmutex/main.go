package main

import (
	"fmt"
	"sync"
)

type test struct {
	data map[string]string
	r    sync.RWMutex
}

const rwmutexMaxReaders = 1 << 30

func main() {
	t := test{
		data: map[string]string{},
		r:    sync.RWMutex{},
	}
	t.add()
	t.read()
}

func (t test) read() {
	t.r.RLock()
	t.r.RLock()

	fmt.Println(t.data)
	t.r.RUnlock()
	t.r.RUnlock()

}

func (t test) add() {
	t.r.Lock()
	t.data["1"] = "test"
	t.r.Unlock()
}
