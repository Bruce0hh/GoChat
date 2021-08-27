package config

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

// 初始化Redis
func initRedis() *redis.Client {
	redisConfig := Viper.RedisConfig
	rdb := redis.NewClient(&redis.Options{
		Addr: redisConfig.Addr,
		DB:   redisConfig.DB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil
	}
	return rdb
}
