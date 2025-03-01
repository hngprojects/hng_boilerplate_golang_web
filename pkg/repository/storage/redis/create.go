package redis

import (
	"encoding/json"
	"fmt"
	"time"
)

func (rdb *Redis) RedisSet(key string, value interface{}) error {
	serialized, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rdb.Red.Set(Ctx, key, serialized, 24*time.Hour).Err()
}

func (rdb *Redis) PushToQueue(value interface{}) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("could not marshal struct: %v", err)
	}

	err = rdb.Red.LPush(Ctx, KeyName, jsonValue).Err()
	if err != nil {
		fmt.Println("could not push to Redis queue: ", err)
	}

	return nil
}
