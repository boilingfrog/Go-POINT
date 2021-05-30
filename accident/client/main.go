package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	transportTimeout()

	transportTimeout()
}

func transportTimeout() {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 3 * time.Second,
		}).DialContext,
	}

	c := http.Client{Transport: transport}

	resp, err := c.Get("http://127.0.0.1:8081/test")
	fmt.Println(resp)
	fmt.Println(err)
}

func httpClientTimeout() {
	c := &http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := c.Get("http://127.0.0.1:8081/test")
	fmt.Println(resp)
	fmt.Println(err)
}

func contextTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequest("GET", "http://127.0.0.1:8081/test", nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	fmt.Println(resp)
	fmt.Println(err)
}
