package config

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

// 初始化Redis
func initRedis() *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr: Viper.RedisConfig.Addr,
		DB:   Viper.RedisConfig.DB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil
	}
	return rdb
}
