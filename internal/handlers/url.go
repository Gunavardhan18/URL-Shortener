package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

// Shorten URL Handler
func (hdlr *handler) ShortenURL(c *fiber.Ctx) error {
	type Request struct {
		LongURL string `json:"long_url"`
	}
	ctx := context.Background()
	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	shortCode, err := hdlr.svc.GenerateShortURL(ctx, req.LongURL)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to shorten URL"})
	}

	return c.JSON(fiber.Map{"short_url": "http://localhost:8080/" + shortCode})
}

// Redirect Handler
func (hdlr *handler) RedirectURL(c *fiber.Ctx) error {
	ctx := context.Background()
	shortCode := c.Params("shortCode")

	longURL, err := hdlr.svc.GetOriginalURL(ctx, shortCode)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "URL not found"})
	}

	return c.Redirect(longURL, 301)
}

// GetAllURLs Handler
func (hdlr *handler) GetAllURLs(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := c.Locals("userID").(uint64)

	urls, err := hdlr.svc.GetAllURLs(ctx, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch URLs"})
	}

	return c.JSON(fiber.Map{"urls": urls})
}

// DeleteURL Handler
func (hdlr *handler) DeleteURL(c *fiber.Ctx) error {
	ctx := context.Background()
	shortCode := c.Params("shortCode")

	err := hdlr.svc.DeleteURL(ctx, shortCode)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete URL"})
	}

	return c.JSON(fiber.Map{"message": "URL deleted successfully"})
}

// UpdateURL Handler
func (hdlr *handler) UpdateURL(c *fiber.Ctx) error {
	type Request struct {
		LongURL string `json:"long_url"`
	}
	ctx := context.Background()
	shortCode := c.Params("shortCode")
	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	err := hdlr.svc.UpdateURL(ctx, shortCode, req.LongURL)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update URL"})
	}

	return c.JSON(fiber.Map{"message": "URL updated successfully"})
}
