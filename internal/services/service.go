package services

import (
	"context"

	"github.com/guna/url-shortener/internal/storage"
)

type Service struct {
	Storage storage.IStorage
	Cache   storage.ICacheStorage
}

func NewService(DB *storage.PostgresDB, Cache *storage.RedisClient) Iservice {
	return &Service{
		Storage: DB,
		Cache:   Cache,
	}
}

type Iservice interface {
	IURLService
	IUserService
	IHealthService
}

type IHealthService interface {
	DBHealthCheck(ctx context.Context) error
	CacheHealthCheck(ctx context.Context) error
}

func (svc *Service) CacheHealthCheck(ctx context.Context) error {
	return svc.Cache.Ping(ctx)
}

func (svc *Service) DBHealthCheck(ctx context.Context) error {
	return svc.Storage.Ping()
}
