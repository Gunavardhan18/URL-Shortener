package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/guna/url-shortener/api"
	"github.com/guna/url-shortener/config"
	"github.com/guna/url-shortener/internal/storage"
)

func main() {
	// Load Config
	cfg := config.Load()

	// Connect to storage (PostgreSQL & Redis)
	db := storage.NewPostgresDB(cfg.DB_URL)
	cache := storage.NewRedisClient()

	// Initialize Fiber
	app := fiber.New()

	// Register routes
	api.SetupRoutes(app, db, cache)

	// Start Server
	port := config.GetPort()
	log.Fatal(app.Listen(":" + port))
	fmt.Println("Server running on port", port)
}
