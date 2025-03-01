package database

import "github.com/go-redis/redis/v8"

type WriteRedis interface {
	RedisSet(key string, value interface{}) error
	RedisDelete(key string) (int64, error)

	PushToQueue(value interface{}) error
}

type ReadRedis interface {
	RedisGet(key string) ([]byte, error)
	PopFromQueue() (interface{}, error)
}

type CacheAccess interface {
	RedisDb() *redis.Client
}
type CacheManager interface {
	WriteRedis
	ReadRedis
	CacheAccess
}
