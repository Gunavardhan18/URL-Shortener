package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/guna/url-shortener/config"
	"github.com/guna/url-shortener/internal/utils"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
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
	Ping(c context.Context) error
	StoreInCache(c *fiber.Ctx, shortCode, longURL string, ttl time.Duration) error
	GetFromCache(c *fiber.Ctx, shortCode string) (string, error)
	DeleteFromCache(c *fiber.Ctx, shortCode string) error
}

func (rc *RedisClient) StoreInCache(c *fiber.Ctx, shortCode, longURL string, ttl time.Duration) error {
	ctx := c.Context()
	tracker := utils.GetTracker(c)
	err := rc.Client.Set(ctx, shortCode, longURL, ttl).Err()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tracker": tracker,
			"err":     err.Error(),
		}).Error(fmt.Sprintf("Failed to store %s in cache", shortCode))
	}
	return err
}

func (rc *RedisClient) GetFromCache(c *fiber.Ctx, shortCode string) (string, error) {
	ctx := c.Context()
	return rc.Client.Get(ctx, shortCode).Result()
}

func (rc *RedisClient) Ping(ctx context.Context) error {
	_, err := rc.Client.Ping(ctx).Result()
	return err
}

func (rc *RedisClient) DeleteFromCache(c *fiber.Ctx, shortCode string) error {
	ctx := c.Context()
	return rc.Client.Del(ctx, shortCode).Err()
}
