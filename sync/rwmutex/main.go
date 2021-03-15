package main

import (
	"fmt"
	"sync"
)

const rwmutexMaxReaders = 1 << 30

//func main() {
//	t := test{
//		data: map[string]string{},
//		r:    sync.RWMutex{},
//	}
//	t.add()
//	t.read()
//}
//
//type test struct {
//	data map[string]string
//	r    sync.RWMutex
//}
//
//func (t test) read() {
//	t.r.RLock()
//	t.r.RLock()
//	t.r.Lock()
//	fmt.Println(t.data)
//	t.r.Unlock()
//	t.r.RUnlock()
//	t.r.RUnlock()
//}
//
//func (t test) add() {
//	t.r.Lock()
//	t.data["1"] = "test"
//	t.r.Unlock()
//}
type Store struct {
	a string
	b string

	sync.RWMutex
}

func main() {
	store := Store{}
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 1; i < 10000; i += 1 {
			fmt.Println("main write ", i)
			store.Write("111", "1111")
		}
	}()
	go func() {
		defer wg.Done()
		for i := 1; i < 10000; i += 1 {
			fmt.Println("main get ab", i)
			store.GetAB()
		}
	}()
	wg.Wait()
}

func (s *Store) GetA() string {
	fmt.Println("get a")
	s.RLock()
	fmt.Println("get a2")
	defer s.RUnlock()

	return s.a
}

func (s *Store) GetAB() (string, string) {
	fmt.Println("get ab")
	s.RLock()
	fmt.Println("get ab2")
	defer s.RUnlock()

	return s.GetA(), s.b
}

func (s *Store) Write(a, b string) {
	fmt.Println("write")
	s.Lock()
	defer s.Unlock()
	fmt.Println("write2")

	s.a = a
	s.b = b
}
