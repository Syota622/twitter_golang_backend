package config

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/go-redis/redis/v8"
)

// ConnectRedis Redisに接続
func ConnectRedis(envConfig EnvConfig) *redis.Client {
	redisDBInt, _ := strconv.Atoi(envConfig.RedisDB)
	rdb := redis.NewClient(&redis.Options{
		Addr:     envConfig.RedisAddr,
		Password: envConfig.RedisPassword,
		DB:       redisDBInt,
	})
	ctx := context.Background()
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis接続エラー: %v", err)
	}
	fmt.Println(pong, "Redisに接続しました！")
	return rdb
}
