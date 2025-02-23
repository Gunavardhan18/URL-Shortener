package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient() *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &RedisClient{Client: rdb}
}

func (rc *RedisClient) StoreInCache(ctx context.Context, shortCode, longURL string, ttl time.Duration) {
	err := rc.Client.Set(ctx, shortCode, longURL, 0).Err()
	if err != nil {
		fmt.Println("Redis error:", err)
	}
}

func (rc *RedisClient) GetFromCache(ctx context.Context, shortCode string) (string, error) {
	return rc.Client.Get(ctx, shortCode).Result()
}
