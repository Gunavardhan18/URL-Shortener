package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/guna/url-shortener/internal/services"
	"github.com/guna/url-shortener/internal/storage"
)

// Shorten URL Handler
func ShortenURL(db *storage.PostgresDB, cache *storage.RedisClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type Request struct {
			LongURL string `json:"long_url"`
		}
		ctx := context.Background()
		req := new(Request)
		if err := c.BodyParser(req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}

		shortCode, err := services.GenerateShortURL(ctx, req.LongURL, db, cache)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to shorten URL"})
		}

		return c.JSON(fiber.Map{"short_url": "http://localhost:8080/" + shortCode})
	}
}

// Redirect Handler
func RedirectURL(db *storage.PostgresDB, cache *storage.RedisClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.Background()
		shortCode := c.Params("shortCode")

		longURL, err := services.GetOriginalURL(ctx, shortCode, db, cache)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "URL not found"})
		}

		return c.Redirect(longURL, 301)
	}
}
