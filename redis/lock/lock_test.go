package lock

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/go-redis/redis/v8"
)

func NewRedis() *Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	return &Redis{rdb}
}

func TestTryLock(t *testing.T) {
	redis := NewRedis()
	createdDate, _ := strconv.Atoi(time.Now().Format("20060102150405"))

	var GetRandomKey = func(key string) string {
		return fmt.Sprintf("test:%s:%d", key, createdDate)
	}

	tests := []struct {
		name            string
		key             string
		Value           string
		isTryLock       bool // 是否加锁
		isLockSuccess   bool // 是否成功加锁
		isUnLock        bool // 是否解锁
		isUnLockSuccess bool // 是否成功解锁
	}{
		{
			name:            "测试 加锁-解锁-1",
			key:             GetRandomKey("test-DefaultLock"),
			isTryLock:       true,
			isUnLock:        true,
			Value:           "123",
			isLockSuccess:   true,
			isUnLockSuccess: true,
		},
		{
			name:            "测试 加锁-解锁-2",
			key:             GetRandomKey("test-DefaultLock"),
			isTryLock:       true,
			isUnLock:        true,
			Value:           "123",
			isLockSuccess:   true,
			isUnLockSuccess: true,
		},
		{
			name:            "测试 加锁-解锁-3",
			key:             GetRandomKey("test-DefaultLock"),
			isTryLock:       true,
			isUnLock:        true,
			Value:           "123",
			isLockSuccess:   true,
			isUnLockSuccess: true,
		},

		{
			name:            "测试 线程1-抢占加锁",
			key:             GetRandomKey("test-DefaultLock"),
			isTryLock:       true,
			isUnLock:        false,
			Value:           "123",
			isLockSuccess:   true,
			isUnLockSuccess: false,
		},
		{
			name:            "测试 线程2-抢占加锁",
			key:             GetRandomKey("test-DefaultLock"),
			isTryLock:       true,
			isUnLock:        false,
			Value:           "456",
			isLockSuccess:   false,
			isUnLockSuccess: false,
		},
		{
			name:            "测试 线程3-抢占加锁",
			key:             GetRandomKey("test-DefaultLock"),
			isTryLock:       true,
			isUnLock:        false,
			Value:           "789",
			isLockSuccess:   false,
			isUnLockSuccess: false,
		},
		{
			name:            "测试 线程1-解锁-失败",
			key:             GetRandomKey("test-DefaultLock"),
			isTryLock:       false,
			isUnLock:        true,
			Value:           "1234567",
			isLockSuccess:   false,
			isUnLockSuccess: false,
		},
		{
			name:            "测试 线程1-解锁-成功",
			key:             GetRandomKey("test-DefaultLock"),
			isTryLock:       false,
			isUnLock:        true,
			Value:           "123",
			isLockSuccess:   false,
			isUnLockSuccess: true,
		},
		{
			name:          "测试 线程4-抢占加锁",
			key:           GetRandomKey("test-DefaultLock"),
			isTryLock:     true,
			isUnLock:      false,
			Value:         "789",
			isLockSuccess: true,
		},
	}
	for _, item := range tests {
		t.Run(item.name, func(t *testing.T) {
			ctx := context.Background()
			// 测试需要加锁
			if item.isTryLock {
				isGetLock, err := redis.TryLock(ctx, item.key, item.Value, time.Second*10)
				t.Log(err)
				assert.Equal(t, item.isLockSuccess, isGetLock)
			}

			// 测试需要解锁
			if item.isUnLock {
				isUnlock, err := redis.Unlock(ctx, item.key, item.Value)
				if !item.isUnLockSuccess {
					assert.Equal(t, false, isUnlock)
				} else {
					assert.Equal(t, err, nil)
				}
			}
		})
	}
}
