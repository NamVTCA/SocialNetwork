package redis

import (
	"context"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"), // ví dụ: localhost:6379
		Password: "",                      // nếu có
		DB:       0,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic("⚠️ Cannot connect to Redis: " + err.Error())
	}
	return rdb
}
