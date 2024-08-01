package redis

import (
	"github.com/go-redis/redis/v8"
)

func RedisGet(rdb *redis.Client, key string) ([]byte, error) {
	serialized, err := rdb.Get(Ctx, key).Bytes()
	return serialized, err
}
