package main

import (
	"fmt"
	"log"

	"github.com/guna/url-shortener/app"
	"github.com/guna/url-shortener/config"
	"github.com/guna/url-shortener/internal/storage"
)

func main() {
	// Load Config
	cfg := config.Load()

	// Connect to storage (PostgreSQL & Redis)
	db := storage.NewPostgresDB(cfg.DB_URL)
	cache := storage.NewRedisClient(cfg)

	// Initialize Fiber
	app := app.NewApplication(db, cache)

	// setup required components
	app.SetupComponents()
	log.Fatal(app.App.Listen(":" + cfg.Port))
	fmt.Println("Server running on port", cfg.Port)
}
