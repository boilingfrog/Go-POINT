package main

import (
	"sync"
)

type test struct {
	mapp map[string]string
	mx   sync.RWMutex
}

const rwmutexMaxReaders = 1 << 30

func main() {
	t := test{}
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		t.mx.RLock()
	}()

	go func() {
		defer wg.Done()
		t.mx.Lock()
	}()

	wg.Wait()

}
