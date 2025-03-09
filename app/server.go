package app

import (
	"github.com/guna/url-shortener/internal/handlers"
	"github.com/guna/url-shortener/internal/services"
)

func (app *Application) setupServices() {
	app.service = services.NewService(app.database, app.cache)
}

func (app *Application) setupHandlers() {
	app.handler = handlers.NewHandler(app.service)
}
