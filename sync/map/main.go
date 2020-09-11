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
	M sync.Map
}

func (o *Ban) visit(ip string) bool {
	fmt.Printf("%s 进来了\n", ip)
	// Load 方法返回两个值，一个是如果能拿到的 key 的 value
	// 还有一个是否能够拿到这个值的 bool 结果
	v, ok := o.M.Load(ip) // sync.Map.Load 去查询对应 key 的值
	if !ok {
		// 如果没有，说明可以访问
		fmt.Printf("名单里没有 %s，可以访问\n", ip)
		// 将用户名存入到 Ban List 中
		o.M.Store(ip, time.Now())
		return false
	}

	// 如果有，则判断用户的时间距离现在是否已经超过了 180 秒，也就是3分钟
	if time.Now().Second()-v.(time.Time).Second() > 180 {
		// 超过则可以继续访问
		fmt.Printf("时间为：%d-%d\n", v.(time.Time).Second(), time.Now().Second())
		// 同时重新存入时间
		o.M.Store(ip, time.Now())
		return false
	}
	// 否则不能访问
	fmt.Printf("名单里有 %s，拒绝访问\n", ip)
	return true
}
func main() {
	var success int64 = 0
	wg := sync.WaitGroup{}
	ban := new(Ban)
	for j := 0; j < 1000; j++ {
		wg.Add(100)
		for i := 0; i < 100; i++ {
			go func(m int) {
				defer wg.Done()
				ip := fmt.Sprintf("192.168.1.%d", m)
				if !ban.visit(ip) {
					fmt.Println(m)
					atomic.AddInt64(&success, 1)
				}
			}(i)
		}
		wg.Wait()
	}
	fmt.Println("success:", success)
}
