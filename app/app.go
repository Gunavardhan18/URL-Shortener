package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guna/url-shortener/internal/handlers"
	"github.com/guna/url-shortener/internal/services"
	"github.com/guna/url-shortener/internal/storage"
)

type Application struct {
	App      *fiber.App
	database *storage.PostgresDB
	cache    *storage.RedisClient
	service  services.Iservice
	handler  handlers.IURLHandler
}

func NewApplication(db *storage.PostgresDB, cache *storage.RedisClient) *Application {
	app := fiber.New(
		fiber.Config{
			AppName:       "URL-shortener",
			CaseSensitive: true,
		},
	)
	return &Application{App: app, database: db, cache: cache}
}

func (app *Application) SetupComponents() {
	app.setupServices()
	app.setupHandlers()
	app.setupRoutes()
}
