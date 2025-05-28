package cache

import (
	"context"
	"entry_task/config"
	"github.com/redis/go-redis/v9"
)

var redisDb *redis.Client

func init() {
	redisDb = redis.NewClient(&redis.Options{
		Addr:     config.ProjectConfig.Redis.Addr,
		Password: config.ProjectConfig.Redis.Password,
		DB:       config.ProjectConfig.Redis.DB,
	})

	_, err := redisDb.Ping(context.Background()).Result()
	if err != nil {
		panic("fail to connect redis, err:" + err.Error())
	}
}

func GetRedisClient() *redis.Client {
	return redisDb
}
