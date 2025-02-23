package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/guna/url-shortener/internal/storage"
)

func HealthCheck(db *storage.PostgresDB, cache *storage.RedisClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.Background()

		// Check Database Connection
		if err := db.DB.Ping(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Database is down",
			})
		}

		// Check Redis Connection
		if _, err := cache.Client.Ping(ctx).Result(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Cache (Redis) is down",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "ok",
			"message": "Service, DB, and Cache are running",
		})
	}
}
