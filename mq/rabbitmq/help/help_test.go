package help

import (
	"fmt"
	"testing"
	"time"
)

func Test_fibonacci(t *testing.T) {

	queue := GetRetryQueue(10000, 10, FIBONACCI)

	for _, q := range queue {
		du := time.Duration(q) * time.Millisecond
		fmt.Println(du.String())
	}
}
