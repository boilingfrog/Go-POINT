package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	contextTimeOut()
}

func timeOut() {
	c := &http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := c.Get("http://127.0.0.1:8081/test")
	fmt.Println(resp)
	fmt.Println(err)
}

func contextTimeOut() {
	ctx, cancel := context.WithCancel(context.TODO())
	timer := time.AfterFunc(5*time.Second, func() {
		cancel()
	})

	fmt.Println(timer)
	req, err := http.NewRequest("GET", "http://127.0.0.1:8081/test", nil)
	if err != nil {
		log.Fatal(err)
	}
	req = req.WithContext(ctx)
}
