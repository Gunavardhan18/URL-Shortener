package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guna/url-shortener/internal/models"
	"github.com/guna/url-shortener/internal/utils"
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
		"user":    user,
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": user,
	})
}

// UpdateProfile updates the user's profile
func (h *handler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint64)
	ctx := c.Context()
	var updateData models.User

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	updateData.ID = userID
	if err := h.svc.UpdateUser(ctx, &updateData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update profile",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Profile updated successfully",
	})
}

// Logout logs out the user
func (h *handler) Logout(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint64)
	ctx := c.Context()

	if err := h.svc.LogoutUser(ctx, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete account",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Account deleted successfully",
	})
}
