package main

import (
	"github.com/go-redis/redis"
)

func _redis() *redis.Client {
	t := conf()
	client := redis.NewClient(&redis.Options{
		Addr:     t.Get("redis.addr").(string),
		Password: t.Get("redis.pass").(string),
		DB:       0,
	})

	return client
}

func rGet(key string) string {
	val, err := _redis().Get(key).Result()
	if err != nil {
		return ""
	}

	return val
}

func rSet(key string, value string) bool {
	err := _redis().Set(key, value, 0).Err()
	if err != nil {
		return false
	}

	return true
}
