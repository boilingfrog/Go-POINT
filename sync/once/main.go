package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {

	once := Once{}
	once.doSlow(func() {
		fmt.Println(12)
	})
}

type Once struct {
	done int32
	m    sync.Mutex
}

func (o *Once) doSlow(f func()) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		defer atomic.StoreInt32(&o.done, 1)
		f()
	}
}
