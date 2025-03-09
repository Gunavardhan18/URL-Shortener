package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func (hdlr *handler) HealthCheck(c *fiber.Ctx) error {
	ctx := context.Background()

	// Check Database Connection
	err := hdlr.svc.DBHealthCheck(ctx)
	if err != nil {
		logrus.Error("error: ", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Database is down",
		})
	}

	err = hdlr.svc.CacheHealthCheck(ctx)
	if err != nil {
		logrus.Error("error: ", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Cache (Redis) is down"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "ok",
		"message": "Service, DB, and Cache are running",
	})
}
