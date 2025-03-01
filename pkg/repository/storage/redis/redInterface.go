package redis

import "github.com/go-redis/redis/v8"

type Redis struct {
	Red *redis.Client
}

func NewRedisConnection(rdb *redis.Client) *Redis {
	return &Redis{Red: rdb}
}

func (rdb *Redis) RedisDb() *redis.Client {
	return rdb.Red
}
