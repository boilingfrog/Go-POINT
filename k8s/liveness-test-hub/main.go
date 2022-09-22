package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var started = time.Now()

func main() {
	http.HandleFunc("/hello", sayHello)
	http.HandleFunc("/healthz", healthz)

	log.Println("【默认项目】服务启动成功 监听端口 8001")
	er := http.ListenAndServe("0.0.0.0:8001", nil)
	if er != nil {
		log.Fatal("ListenAndServe: ", er)
	}
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	data := map[string]interface{}{
		"status":  "ok",
		"message": "hello",
	}

	json.NewEncoder(w).Encode(&data)
}

func healthz(w http.ResponseWriter, r *http.Request) {

	duration := time.Now().Sub(started)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if duration.Seconds() > 10 {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("error: %v", duration.Seconds())))
	} else {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}
}
