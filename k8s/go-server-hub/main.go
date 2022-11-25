package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/hello", sayHello)

	log.Println("【默认项目】服务启动成功 监听端口 80")
	er := http.ListenAndServe("0.0.0.0:80", nil)
	if er != nil {
		log.Fatal("ListenAndServe: ", er)
	}
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	log.Println("request hello")
	data := map[string]interface{}{
		"status":  "ok",
		"message": "hello",
	}

	json.NewEncoder(w).Encode(&data)
}
