package services

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/guna/url-shortener/internal/models"
	"github.com/guna/url-shortener/internal/utils"
	"github.com/sirupsen/logrus"
)

type IURLService interface {
	GenerateShortURL(*fiber.Ctx, string, uint64) (string, error)
	GetOriginalURL(*fiber.Ctx, string) (string, error)
	GetAllURLs(ctx *fiber.Ctx, userID uint64) (*models.GetAllURLsResponse, error)
	DeleteURL(ctx *fiber.Ctx, shortCode string) error
	UpdateURL(ctx *fiber.Ctx, shortCode, longURL string) error
}

func (svc *Service) GenerateShortURL(ctx *fiber.Ctx, longURL string, userId uint64) (string, error) {
	shortCode := utils.GenerateShortCode(longURL)
	url := models.URL{
		ShortURL:  shortCode,
		LongURL:   longURL,
		UserID:    userId,
		Status:    models.URLActive,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7), // 1 week
	}
	err := svc.Storage.SaveURL(url)
	if err != nil {
		return "", err
	}
	svc.Cache.StoreInCache(ctx, shortCode, longURL, 60*time.Minute)
	return shortCode, nil
}

func (svc *Service) GetOriginalURL(ctx *fiber.Ctx, shortCode string) (string, error) {
	tracker := utils.GetTracker(ctx)
	url, err := svc.Cache.GetFromCache(ctx, shortCode)
	if err == nil {
		err := svc.Storage.RegisterURLAnalytics(ctx, shortCode)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"tracker": tracker,
				"err":     err.Error(),
			}).Error("Failed to register URL analytics")
		}
		return url, nil
	}

	url, err = svc.Storage.GetURL(shortCode)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tracker": tracker,
			"err":     err.Error(),
		}).Error("Failed to get original URL")
		return "", err
	}
	svc.Cache.StoreInCache(ctx, shortCode, url, 60*time.Minute)
	err = svc.Storage.RegisterURLAnalytics(ctx, shortCode)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tracker": tracker,
			"err":     err.Error(),
		}).Error("Failed to register URL analytics")
	}
	return url, nil
}

func (svc *Service) GetAllURLs(ctx *fiber.Ctx, userID uint64) (*models.GetAllURLsResponse, error) {
	response := &models.GetAllURLsResponse{}
	tracker := utils.GetTracker(ctx)
	urls, err := svc.Storage.GetAllURLs(userID)
	if err != nil {
		return nil, err
	}
	for _, url := range urls {
		url.Clicks, err = svc.Storage.GetClicks(userID, url.ShortURL)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"tracker": tracker,
				"err":     err.Error(),
			}).Error("Failed to get URL clicks")
		}
	}
	response.URLs = urls
	return response, nil
}

func (svc *Service) DeleteURL(ctx *fiber.Ctx, shortCode string) error {
	err := svc.Storage.DeleteURL(shortCode)
	if err != nil {
		return err
	}
	svc.Cache.DeleteFromCache(ctx, shortCode)
	return nil
}

func (svc *Service) UpdateURL(ctx *fiber.Ctx, shortCode, longURL string) error {
	err := svc.Storage.UpdateURL(shortCode, longURL)
	if err != nil {
		return err
	}
	svc.Cache.StoreInCache(ctx, shortCode, longURL, 60*time.Minute)
	return nil
}
