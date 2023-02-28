package data

import "time"

func NewRedis(data *Data) *RedisRepo {
	return &RedisRepo{
		data: data,
	}
}

type RedisRepo struct {
	data *Data
}

func (r *RedisRepo) Get(key string) (string, error) {
	return r.data.redis.Get(key).Result()
}

func (r *RedisRepo) Set(key string, val string) (string, error) {
	return r.data.redis.Set(key, val, time.Duration(time.Hour)).Result()
}
