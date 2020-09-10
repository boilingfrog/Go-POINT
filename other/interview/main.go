package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

//func main() {
//	ch := make(chan struct{})
//	go func() {
//		fmt.Println("start working")
//		time.Sleep(time.Second * 1)
//		ch <- struct{}{}
//	}()
//
//	<-ch
//
//	fmt.Println("finished")
//}

func gen(ctx context.Context) <-chan int {
	ch := make(chan int)
	go func() {
		var n int
		for {
			select {
			case <-ctx.Done():
				return
			case ch <- n:
				n++
				time.Sleep(time.Second)
			}
		}
	}()
	return ch
}

func main() {
	mux := http.NewServeMux()

	rh := http.RedirectHandler("http://www.baidu.com", 307)
	mux.Handle("/foo", rh)

	var dmap = make(map[int]string, 1)
	_ = dmap

	log.Println("Listening...")
	http.ListenAndServe(":3000", mux)
}
