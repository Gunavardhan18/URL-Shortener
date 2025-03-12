package app

import "github.com/guna/url-shortener/internal/middleware"

func (app *Application) setupRoutes() {
	app.SetupHealthRoutes()
	app.SetupUserRoutes()
	app.SetupURLRoutes()
}

func (app *Application) SetupHealthRoutes() {
	app.App.Get("/urlshortner/health", app.handler.HealthCheck)
}

func (app *Application) SetupUserRoutes() {
	app.App.Post("/urlshortner/auth/signup", app.handler.SignUp)
	app.App.Post("/urlshortner/auth/login", app.handler.Login)
	app.App.Post("/urlshortner/auth/delete", middleware.AuthMiddleware(), app.handler.DeleteAccount)
	app.App.Get("/urlshortner/user/profile", middleware.AuthMiddleware(), app.handler.GetProfile)
	app.App.Put("/urlshortner/user/profile", middleware.AuthMiddleware(), app.handler.UpdateProfile)
}

func (app *Application) SetupURLRoutes() {
	app.App.Post("urlshortner/shorten", app.handler.ShortenURL)
	app.App.Get("urlshortner/:shortCode", app.handler.RedirectURL)
	app.App.Put("urlshortner/:shortCode", app.handler.UpdateURL)
	app.App.Delete("/urlshortner/urls/:shortCode", app.handler.DeleteURL)
	app.App.Get("/urlshortner/urls", app.handler.GetAllURLs)
}
