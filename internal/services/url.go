package services

import (
	"context"

	"github.com/guna/url-shortener/internal/models"
	"github.com/guna/url-shortener/internal/utils"
)

type IURLService interface {
	GenerateShortURL(context.Context, string) (string, error)
	GetOriginalURL(context.Context, string) (string, error)
	GetAllURLs(ctx context.Context, userID uint64) ([]models.URL, error)
	DeleteURL(ctx context.Context, shortCode string) error
	UpdateURL(ctx context.Context, shortCode, longURL string) error
}

func (svc *Service) GenerateShortURL(ctx context.Context, longURL string) (string, error) {
	shortCode := utils.GenerateHash(longURL)
	err := svc.Storage.SaveURL(shortCode, longURL)
	if err != nil {
		return "", err
	}
	svc.Cache.StoreInCache(ctx, shortCode, longURL, -1)
	return shortCode, nil
}

func (svc *Service) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	url, err := svc.Cache.GetFromCache(ctx, shortCode)
	if err == nil {
		return url, nil
	}

	url, err = svc.Storage.GetURL(shortCode)
	if err != nil {
		return "", err
	}
	svc.Cache.StoreInCache(ctx, shortCode, url, -1)
	return url, nil
}

func (svc *Service) GetAllURLs(ctx context.Context, userID uint64) ([]models.URL, error) {
	urls, err := svc.Storage.GetAllURLs(userID)
	if err != nil {
		return nil, err
	}
	return urls, nil
}

func (svc *Service) DeleteURL(ctx context.Context, shortCode string) error {
	err := svc.Storage.DeleteURL(shortCode)
	if err != nil {
		return err
	}
	svc.Cache.DeleteFromCache(ctx, shortCode)
	return nil
}

func (svc *Service) UpdateURL(ctx context.Context, shortCode, longURL string) error {
	err := svc.Storage.UpdateURL(shortCode, longURL)
	if err != nil {
		return err
	}
	svc.Cache.StoreInCache(ctx, shortCode, longURL, -1)
	return nil
}
