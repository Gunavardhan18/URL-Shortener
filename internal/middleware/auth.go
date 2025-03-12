package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/guna/url-shortener/internal/utils"
	"github.com/sirupsen/logrus"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		tracker := utils.GetTracker(c)
		if authHeader == "" {
			logrus.WithFields(logrus.Fields{
				"tracker": tracker,
			}).Error("Missing token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing token",
			})
		}

		if !strings.Contains(authHeader, "Bearer ") {
			logrus.WithFields(logrus.Fields{
				"tracker": tracker,
				"token":   authHeader,
			}).Error("Invalid token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		tokenString := authHeader[len("Bearer "):]
		userID, err := utils.ValidateJWT(tokenString)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"tracker": tracker,
				"err":     err.Error(),
			}).Error("Invalid token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		c.Locals("userID", userID)
		return c.Next()
	}
}
