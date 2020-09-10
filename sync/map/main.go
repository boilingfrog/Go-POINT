package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

/*

场景：在一个高并发的web服务器中，要限制IP的频繁访问。现模拟100个IP同时并发访问服务器，每个IP要重复访问1000次。每个IP三分钟之内只能访问一次。
修改以下代码完成该过程，要求能成功输出 success:100

*/

type Ban struct {
	visitIPs map[string]time.Time
}

func NewBan() *Ban {
	return &Ban{visitIPs: make(map[string]time.Time)}
}
func (o *Ban) visit(ip string) bool {
	if _, ok := o.visitIPs[ip]; ok {
		return true
	}
	o.visitIPs[ip] = time.Now()
	return false
}
func main() {
	var success int64 = 0
	wg := sync.WaitGroup{}
	ban := NewBan()
	wg.Add(1000)

	for i := 0; i < 100; i++ {
		go func() {
			ip := fmt.Sprintf("192.168.1.%d", 1)
			if !ban.visit(ip) {
				atomic.AddInt64(&success, 1)
			}
			wg.Done()
		}()
	}
	fmt.Println("success:", success)
}
