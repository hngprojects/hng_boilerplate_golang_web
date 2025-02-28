package redis

func (rdb *Redis) RedisDelete(key string) (int64, error) {
	deleted, err := rdb.Red.Del(Ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
