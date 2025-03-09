package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/guna/url-shortener/config"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(cfg config.ConnectionConfig) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
		os.Exit(1)
	}
	log.Println("Connected to Redis successfully!")
	return &RedisClient{Client: rdb}
}

type ICacheStorage interface {
	Ping(ctx context.Context) error
	StoreInCache(ctx context.Context, shortCode, longURL string, ttl time.Duration) error
	GetFromCache(ctx context.Context, shortCode string) (string, error)
	DeleteFromCache(ctx context.Context, shortCode string) error
}

func (rc *RedisClient) StoreInCache(ctx context.Context, shortCode, longURL string, ttl time.Duration) error {
	err := rc.Client.Set(ctx, shortCode, longURL, 0).Err()
	if err != nil {
		fmt.Println("Redis error:", err)
	}
	return err
}

func (rc *RedisClient) GetFromCache(ctx context.Context, shortCode string) (string, error) {
	return rc.Client.Get(ctx, shortCode).Result()
}

func (rc *RedisClient) Ping(ctx context.Context) error {
	_, err := rc.Client.Ping(ctx).Result()
	return err
}

func (rc *RedisClient) DeleteFromCache(ctx context.Context, shortCode string) error {
	return rc.Client.Del(ctx, shortCode).Err()
}
