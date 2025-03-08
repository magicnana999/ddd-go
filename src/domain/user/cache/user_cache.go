package cache

import (
	"github.com/go-redis/redis/v8"
	"github.com/magicnana999/ddd-go/infrastructure"
	"sync"
)

var (
	UserCacheInstance *UserCache
	uciOnce           sync.Once
)

type UserCache struct {
	rdb *redis.Client
}

func InitUserCache() *UserCache {
	uciOnce.Do(func() {
		UserCacheInstance = &UserCache{
			rdb: infrastructure.InitRedis(),
		}
	})
	return UserCacheInstance
}
