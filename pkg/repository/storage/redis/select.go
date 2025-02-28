package redis

import (
	"encoding/json"
	"fmt"
)

func (rdb *Redis) RedisGet(key string) ([]byte, error) {
	serialized, err := rdb.Red.Get(Ctx, key).Bytes()
	return serialized, err
}

func (rdb *Redis) PopFromQueue() (interface{}, error) {
	var response interface{}

	jsonValue, err := rdb.Red.RPop(Ctx, KeyName).Result()
	if err != nil {
		return response, fmt.Errorf("could not pop from Redis queue: %v", err)
	}

	err = json.Unmarshal([]byte(jsonValue), &response)
	if err != nil {
		return response, fmt.Errorf("could not unmarshal JSON: %v", err)
	}

	return response, nil
}
