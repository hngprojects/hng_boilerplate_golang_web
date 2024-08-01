package redis

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

func RedisSet(rdb *redis.Client, key string, value interface{}) error {
	serialized, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rdb.Set(Ctx, key, serialized, 24*time.Hour).Err()
}
