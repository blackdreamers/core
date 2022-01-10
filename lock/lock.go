package lock

import (
	"context"
	"sync"
	"time"

	"github.com/blackdreamers/core/cache/redis"
	"github.com/blackdreamers/core/consts"
	log "github.com/blackdreamers/core/logger"
	"github.com/blackdreamers/core/retry"
)

type Lock struct {
	ctx        context.Context
	key        string
	expireTime time.Duration
	mu         sync.Mutex
	retry      *retry.Retry
}

// NewLock 创建锁
func NewLock(key string) *Lock {
	return &Lock{
		ctx:        context.Background(),
		key:        key,
		expireTime: 5 * time.Minute,
		mu:         sync.Mutex{},
		retry:      retry.NewRetry("lock", retry.Delay(100*time.Millisecond), retry.Attempts(uint(5))),
	}
}

// Key 获取当前锁的redis key
func (l *Lock) Key() string {
	return l.key
}

// SetExpire 设置锁过期时间，默认5分钟
func (l *Lock) SetExpire(expireTime time.Duration) {
	l.expireTime = expireTime
}

// Lock 获取锁
func (l *Lock) Lock() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	result, err := redis.SetNX(l.ctx, l.key, "true", l.expireTime)
	if err != nil {
		log.Fields(consts.ErrKey, err, "key", l.key).Log(log.ErrorLevel, "get lock err")
	}

	return result
}

// UnLock 释放锁
func (l *Lock) UnLock() {
	err := l.retry.Do(
		func() error {
			_, err := redis.Del(l.ctx, l.key)
			if err != nil {
				return err
			}
			return nil
		},
		retry.OnRetry(func(n uint, err error) {
			log.Fields(
				consts.ErrKey, err,
				"key", l.key,
				"num", n,
			).Log(log.WarnLevel, "retry")
		}),
	)
	if err != nil {
		log.Fields(consts.ErrKey, err, "key", l.key).Log(log.ErrorLevel, "unlock err")
	}
}

// Expire 重置锁时间
func (l *Lock) Expire(expireTime time.Duration) error {
	count, err := redis.Exists(l.ctx, l.key)
	if err != nil {
		return err
	}

	if expireTime != 0 {
		l.expireTime = expireTime
	}

	if count > 0 {
		_, err = redis.Expire(l.ctx, l.key, l.expireTime)
		if err != nil {
			return err
		}
	} else {
		l.Lock()
	}

	return nil
}
