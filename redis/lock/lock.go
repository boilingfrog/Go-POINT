package lock

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	*redis.Client
}

func tryLockScript() string {
	script := `
		local key = KEYS[1]

		local value = ARGV[1] 
		local expireTime = ARGV[2] 
		local isSuccess = redis.call('SETNX', key, value)

		if isSuccess == 1 then
			redis.call('EXPIRE', key, expireTime)
			return "OK"
		end

		return "unLock"    `
	return script
}

func unLockScript() string {
	script := `
		local value = ARGV[1] 
		local key = KEYS[1]

		local keyValue = redis.call('GET', key)
		if tostring(keyValue) == tostring(value) then
			return redis.call('DEL', key)
		else
			return 0
		end
    `
	return script
}

var UnLockErr = errors.New("未解锁成功")

// 使用 set nx
// res, err := r.Do(ctx, "set", key, value, "px", expire.Milliseconds(), "nx").Result()
func (r *Redis) TryLock(ctx context.Context, key, value string, expire time.Duration) (isGetLock bool, err error) {
	// 使用 Lua + SETNX
	res, err := r.Eval(ctx, tryLockScript(), []string{key}, value, expire.Seconds()).Result()
	if err != nil {
		return false, err
	}
	if res == "OK" {
		return true, nil
	}
	return false, nil
}

func (r *Redis) Unlock(ctx context.Context, key, value string) (bool, error) {
	res, err := r.Eval(ctx, unLockScript(), []string{key}, value).Result()
	if err != nil {
		return false, err
	}

	return res.(int64) != 0, nil
}
