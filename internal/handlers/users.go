package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guna/url-shortener/internal/models"
	"github.com/guna/url-shortener/internal/utils"
	"github.com/sirupsen/logrus"
)

// SignUp handles user registration
func (h *handler) SignUp(c *fiber.Ctx) error {
	var user models.User
	ctx := c.Context()
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	if err := h.svc.CreateUser(ctx, &user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	logrus.Info("User created successfully")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"user":    user,
	})
}

// Login handles user authentication
func (h *handler) Login(c *fiber.Ctx) error {
	loginRequest := models.LoginRequest{}

	ctx := c.Context()

	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user, err := h.svc.AuthenticateUser(ctx, loginRequest.Email, loginRequest.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}
	c.Set("Authorization", "Bearer "+token)

	logrus.Info("Login successful")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
	})
}

// GetProfile retrieves the user's profile
func (h *handler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint64)
	ctx := c.Context()

	user, err := h.svc.GetUserByID(ctx, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	logrus.Info("User profile retrieved successfully")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": user,
	})
}

// UpdateProfile updates the user's profile
func (h *handler) UpdateProfile(c *fiber.Ctx) error {
	ctx := c.Context()
	var updateData models.UpdateUserRequest

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if updateData.OldPassword == updateData.NewPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "New password cannot be same as old password",
		})
	}

	if err := h.svc.UpdateUser(ctx, &updateData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update profile",
		})
	}

	logrus.Info("Profile updated successfully")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Profile updated successfully",
	})
}

// Logout logs out the user
func (h *handler) Logout(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := c.Locals("userID").(uint64)
	if err := h.svc.LogoutUser(ctx, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to logout user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User logged out successfully",
	})
}
