package help

import (
	"math"
)

type RetryStrategy int

const (
	FIBONACCI RetryStrategy = iota
	EXPONENTIAL
	LINEAR
	CUSTOMQUEUE
)

type Retry struct {
	Delay string
	Max   int
	Queue []string
}

func fibonacci() func() int {
	var x, y int = 0, 1
	return func() int {
		x, y = y, x+y
		return y
	}
}

func fib(times int32) (num int64) {
	f := fibonacci()
	num = int64(f())
	for i := 0; i < int(times); i++ {
		num = int64(f())
	}
	return
}

var retryStrategies = map[RetryStrategy]func(int32) int64{
	FIBONACCI:   fib,
	EXPONENTIAL: func(t int32) int64 { return int64(math.Pow(2, float64(t))) },
	LINEAR:      func(t int32) int64 { return int64(t) + 1 },
}

func GetRetryDelay(delayQueue []int64, delay, count int32, strategy RetryStrategy) int64 {
	if strategy == CUSTOMQUEUE && len(delayQueue) > 0 {
		return delayQueue[count]
	}

	fn, ok := retryStrategies[strategy]
	if !ok {
		fn = retryStrategies[FIBONACCI]
	}
	return int64(delay) * fn(count)
}

func GetRetryQueue(delay int64, max int, strategy RetryStrategy) []int64 {
	fn, ok := retryStrategies[strategy]
	if !ok {
		fn = retryStrategies[FIBONACCI]
	}
	queue := []int64{}

	for i := 0; i < max; i++ {
		queue = append(queue, delay*fn(int32(i)))
	}
	return queue
}
