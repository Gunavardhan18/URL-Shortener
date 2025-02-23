package services

import (
	"context"

	"github.com/guna/url-shortener/internal/storage"
	"github.com/guna/url-shortener/internal/utils"
)

func GenerateShortURL(ctx context.Context, longURL string, db *storage.PostgresDB, cache *storage.RedisClient) (string, error) {
	shortCode := utils.GenerateHash(longURL)
	err := db.SaveURL(shortCode, longURL)
	if err != nil {
		return "", err
	}
	cache.StoreInCache(ctx, shortCode, longURL, -1)
	return shortCode, nil
}

func GetOriginalURL(ctx context.Context, shortCode string, db *storage.PostgresDB, cache *storage.RedisClient) (string, error) {
	url, err := cache.GetFromCache(ctx, shortCode)
	if err == nil {
		return url, nil
	}

	url, err = db.GetURL(shortCode)
	if err != nil {
		return "", err
	}
	cache.StoreInCache(ctx, shortCode, url, -1)
	return url, nil
}
