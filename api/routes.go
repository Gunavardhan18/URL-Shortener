package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guna/url-shortener/internal/handlers"
	"github.com/guna/url-shortener/internal/storage"
)

func SetupRoutes(app *fiber.App, db *storage.PostgresDB, cache *storage.RedisClient) {
	app.Get("/health", handlers.HealthCheck(db, cache))

	// URL Shortening APIs
	app.Post("/shorten", handlers.ShortenURL(db, cache))
	app.Get("/:shortCode", handlers.RedirectURL(db, cache))
}
