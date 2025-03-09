package app

func (app *Application) setupRoutes() {
	app.SetupHealthRoutes()
	app.SetupUserRoutes()
	app.SetupURLRoutes()
}

func (app *Application) SetupHealthRoutes() {
	app.App.Get("/health", app.handler.HealthCheck)
}

func (app *Application) SetupUserRoutes() {
	app.App.Post("/api/auth/signup", app.handler.SignUp)
	app.App.Post("/api/auth/login", app.handler.Login)
	app.App.Delete("/api/auth/logout", app.handler.Logout)
	app.App.Get("/api/user/profile", app.handler.GetProfile)
}

func (app *Application) SetupURLRoutes() {
	app.App.Post("/shorten", app.handler.ShortenURL)
	app.App.Get("/:shortCode", app.handler.RedirectURL)
	app.App.Put("/:shortCode", app.handler.UpdateURL)
	app.App.Delete("/api/urls/:shortCode", app.handler.DeleteURL)
	app.App.Get("/api/urls", app.handler.GetAllURLs)
}
