package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	server()
}

func handler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(3 * time.Second)
	fmt.Println("测试超时")

	w.Write([]byte("hello world"))
}

func serverHandle() {
	srv := http.Server{
		Addr:         ":8081",
		WriteTimeout: 1 * time.Second,
		Handler:      http.TimeoutHandler(http.HandlerFunc(handler), 5*time.Second, "Timeout!\n"),
	}
	if err := srv.ListenAndServe(); err != nil {
		os.Exit(1)
	}
}

func server() {
	// 开启pprof
	go func() {
		srv := http.Server{
			Addr:         ":8081",
			WriteTimeout: 1 * time.Second,
			Handler:      http.TimeoutHandler(http.HandlerFunc(handler), 1*time.Second, "Timeout!\n"),
		}
		if err := srv.ListenAndServe(); err != nil {
			os.Exit(1)
		}
	}()

	// 路由，访问，触发内存泄露的代码判断
	http.HandleFunc("/test", handler)

	// 阻塞
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
