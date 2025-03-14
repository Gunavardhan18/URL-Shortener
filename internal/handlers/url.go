package handlers

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/guna/url-shortener/internal/models"
	"github.com/guna/url-shortener/internal/utils"
	"github.com/sirupsen/logrus"
)

// Shorten URL Handler
func (hdlr *handler) ShortenURL(c *fiber.Ctx) error {
	tracker := utils.GetTracker(c)
	userID := c.Locals("userID").(uint64)
	req := new(models.URLRequest)
	if err := c.BodyParser(req); err != nil {
		logrus.WithFields(logrus.Fields{
			"tracker": tracker,
			"err":     err.Error(),
		}).Error("Invalid request")
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.LongURL == "" {
		logrus.WithFields(logrus.Fields{
			"tracker": tracker,
		}).Error("Long URL is required")
		return c.Status(400).JSON(fiber.Map{"error": "Long URL is required"})
	}

	shortCode, err := hdlr.svc.GenerateShortURL(c, req.LongURL, userID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tracker": tracker,
			"err":     err.Error(),
		}).Error("Failed to shorten URL")
		return c.Status(500).JSON(fiber.Map{"error": "Failed to shorten URL"})
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	shortURL := fmt.Sprintf("%s/%s", baseURL, shortCode)

	return c.JSON(fiber.Map{"short_url": shortURL})
}

// Redirect Handler
func (hdlr *handler) RedirectURL(c *fiber.Ctx) error {
	tracker := utils.GetTracker(c)
	shortCode := c.Params("shortCode")

	longURL, err := hdlr.svc.GetOriginalURL(c, shortCode)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tracker":   tracker,
			"err":       err.Error(),
			"shortCode": shortCode,
		}).Error("URL not found")
		return c.Status(404).JSON(fiber.Map{"error": "URL not found"})
	}

	if !strings.HasPrefix(longURL, "http://") && !strings.HasPrefix(longURL, "https://") {
		longURL = "http://" + longURL // Default to HTTP if scheme is missing
	}

	logrus.WithFields(logrus.Fields{
		"tracker":   tracker,
		"shortCode": shortCode,
		"longURL":   longURL,
	}).Info("Redirecting to long URL")

	return c.Redirect(longURL, fiber.StatusFound)
}

// GetAllURLs Handler
func (hdlr *handler) GetAllURLs(c *fiber.Ctx) error {
	tracker := utils.GetTracker(c)
	userID := c.Locals("userID").(uint64)

	urls, err := hdlr.svc.GetAllURLs(c, userID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tracker": tracker,
			"err":     err.Error(),
		}).Error("Failed to fetch URLs")
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch URLs"})
	}

	logrus.WithFields(logrus.Fields{
		"tracker": tracker,
		"urls":    urls,
	}).Info("URLs fetched successfully")
	return c.JSON(fiber.Map{"urls": urls})
}

// DeleteURL Handler
func (hdlr *handler) DeleteURL(c *fiber.Ctx) error {
	shortCode := c.Params("shortCode")
	tracker := utils.GetTracker(c)

	err := hdlr.svc.DeleteURL(c, shortCode)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tracker": tracker,
			"err":     err.Error(),
		}).Error("Failed to delete URL")
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete URL"})
	}
	logrus.WithFields(logrus.Fields{
		"tracker":   tracker,
		"shortCode": shortCode,
	}).Info("URL deleted successfully")

	return c.JSON(fiber.Map{"message": "URL deleted successfully"})
}

// UpdateURL Handler
func (hdlr *handler) UpdateURL(c *fiber.Ctx) error {
	tracker := utils.GetTracker(c)
	shortCode := c.Params("shortCode")
	req := new(models.URLRequest)
	if err := c.BodyParser(req); err != nil {
		logrus.WithFields(logrus.Fields{
			"tracker": tracker,
			"err":     err.Error(),
		}).Error("UpdateURL: Invalid request")
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	err := hdlr.svc.UpdateURL(c, shortCode, req.LongURL)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tracker": tracker,
			"err":     err.Error(),
		}).Error("Failed to update URL")
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update URL"})
	}

	logrus.WithFields(logrus.Fields{
		"tracker":   tracker,
		"shortCode": shortCode,
	}).Info("URL updated successfully")

	return c.JSON(fiber.Map{"message": "URL updated successfully"})
}
