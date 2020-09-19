package main

import (
	"fmt"
	"sync"
)

type test struct {
	mapp map[string]string
	mx   sync.RWMutex
}

const rwmutexMaxReaders = 1 << 30

func main() {
	t := test{}
	t.mx.RLock()
	t.mx.Lock()
	t.mx.RLock()

	fmt.Println(rwmutexMaxReaders)
}
