package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guna/url-shortener/internal/services"
)

type handler struct {
	svc services.Iservice
}

func NewHandler(svc services.Iservice) *handler {
	return &handler{
		svc: svc,
	}
}

type IURLHandler interface {
	HealthCheck(c *fiber.Ctx) error
	SignUp(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	DeleteAccount(c *fiber.Ctx) error
	GetProfile(c *fiber.Ctx) error
	UpdateProfile(c *fiber.Ctx) error
	ShortenURL(c *fiber.Ctx) error
	RedirectURL(c *fiber.Ctx) error
	UpdateURL(c *fiber.Ctx) error
	DeleteURL(c *fiber.Ctx) error
	GetAllURLs(c *fiber.Ctx) error
}
