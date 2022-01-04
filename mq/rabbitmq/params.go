package rabbitmq

import (
	"Go-POINT/mq/rabbitmq/help"
	"log"
	"time"
)

type Handler func([]byte) HandleFLag

func DefaultHandler(key string, svc func([]byte) error) Handler {
	return func(data []byte) HandleFLag {
		defer func() {
			if err := recover(); err != nil {
				log.Println("jobPanic", "key", key, "err", err)
			}
		}()
		var start = time.Now()
		err := svc(data)
		elapsed := time.Since(start)

		if err != nil {
			if err == HandleRequeue {
				log.Println(key, "name: ", key, "data: ", string(data), "elapsed: ", elapsed, "err: ", "retry")
				return HandleRequeue
			} else if err == HandleDrop {
				log.Println(key, "name: ", key, "data: ", string(data), "elapsed: ", elapsed, "err:", "drop")
				return HandleDrop
			} else {
				log.Println(key, "name: ", key, "data: ", string(data), "elapsed: ", elapsed, "err: ", err)
			}
		} else {
			log.Println(key, "name: ", key, "data: ", string(data), "elapsed: ", elapsed)
		}
		return HandleSuccess
	}
}

type params struct {
	key         string // routing key
	queue       string // queue_name
	handler     Handler
	concurrency int     // 并发数
	prefetch    int     // 预获取
	retryQueue  []int64 // 重试队列
}

func (ps *params) Queue() string {
	if ps.queue != "" {
		return ps.queue
	}
	return ps.key
}

func (ps *params) Handle(data []byte) HandleFLag {
	return ps.handler(data)
}

type Param func(*params)

func WithRetry(startegy help.RetryStrategy, retry help.Retry) Param {
	return func(p *params) {
		if startegy == help.CUSTOMQUEUE {
			for _, delay := range retry.Queue {
				d, err := time.ParseDuration(delay)
				if err != nil {
					panic(err)
				}
				p.retryQueue = append(p.retryQueue, int64(d/time.Millisecond))
			}
			return
		}
		d, err := time.ParseDuration(retry.Delay)
		if err != nil {
			panic(err)
		}
		p.retryQueue = help.GetRetryQueue(int64(d/time.Millisecond), retry.Max, startegy)
	}
}

func WithConcurrency(concurrency int) Param {
	return func(p *params) {
		p.concurrency = concurrency
	}
}

func WithPrefetch(prefetch int) Param {
	return func(p *params) {
		p.prefetch = prefetch
	}
}

func WithQueue(queue string) Param {
	return func(p *params) {
		p.queue = queue
	}
}

func evaParam(param []Param) *params {
	ps := &params{}
	for _, p := range param {
		p(ps)
	}
	return ps
}
