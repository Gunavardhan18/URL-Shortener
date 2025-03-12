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
	tracker := utils.GetTracker(c)
	if err := c.BodyParser(&user); err != nil {
		logrus.WithFields(logrus.Fields{
			"tracker": tracker,
		}).Error("Invalid request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	requiredFields := map[string]string{
		"Password": user.Password,
		"Email":    user.Email,
		"Name":     user.Name,
	}

	for field, value := range requiredFields {
		if value == "" {
			logrus.WithFields(logrus.Fields{
				"tracker": tracker,
			}).Errorf(field + " is required")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": field + " is required",
				"msg":   "User registration failed",
			})
		}
	}

	if err := h.svc.CreateUser(c, &user); err != nil {
		logrus.WithFields(logrus.Fields{
			"tracker": tracker,
			"err":     err.Error(),
		}).Error("Failed to create user")
		if err.Error() == models.ErrUserNameExists || err.Error() == models.ErrUserWithEmailExists {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
				"msg":   "User registration failed",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
			"msg":   "User registration failed",
		})
	}

	logrus.WithFields(logrus.Fields{
		"tracker": tracker,
	}).Info("User created successfully")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
	})
}

// Login handles user authentication
func (h *handler) Login(c *fiber.Ctx) error {
	loginRequest := models.LoginRequest{}

	tracker := utils.GetTracker(c)

	if err := c.BodyParser(&loginRequest); err != nil {
		logrus.WithFields(logrus.Fields{
			"tracker": tracker,
		}).Error("Invalid request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"msg":   "Login Failed",
		})
	}

	user, err := h.svc.AuthenticateUser(c, loginRequest.Email, loginRequest.Password)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tracker": tracker,
			"err":     err.Error(),
		}).Error("Login failed")
		if err.Error() == models.ErrInvalidCredentials || err.Error() == models.ErrPasswordDoestMatch {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
				"msg":   "Login Failed",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
			"msg":   "Login Failed",
		})
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tracker": tracker,
			"err":     err.Error(),
		}).Error("Failed to generate token")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
			"msg":   "Login Failed",
		})
	}
	c.Set("Authorization", "Bearer "+token)

	logrus.WithFields(logrus.Fields{
		"tracker": tracker,
	}).Info("Login successful")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful, hello " + user.Name,
	})
}

// GetProfile retrieves the user's profile
func (h *handler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint64)
	tracker := utils.GetTracker(c)
	user, err := h.svc.GetUserByID(c, userID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tracker": tracker,
			"err":     err.Error(),
		}).Error("User not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	profileRespone := models.ProfileResponse{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Status: user.Status,
	}
	logrus.WithFields(logrus.Fields{
		"tracker": tracker,
	}).Info("Profile retrieved successfully")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": profileRespone,
	})
}

// UpdateProfile updates the user's profile
func (h *handler) UpdateProfile(c *fiber.Ctx) error {
	var updateData models.UpdateUserRequest
	userID := c.Locals("userID").(uint64)
	tracker := utils.GetTracker(c)
	if err := c.BodyParser(&updateData); err != nil {
		logrus.WithFields(logrus.Fields{
			"tracker": tracker,
		}).Error("Invalid request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	updateData.UserID = userID

	requiredFields := map[string]string{
		"old_password": updateData.OldPassword,
		"new_password": updateData.NewPassword,
		"Email":        updateData.Email,
	}

	for field, value := range requiredFields {
		if value == "" {
			logrus.WithFields(logrus.Fields{
				"tracker": tracker,
			}).Errorf(field + " is required")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": field + " is required",
				"msg":   "User registration failed",
			})
		}
	}

	if updateData.UpdateName && updateData.Name == "" {
		logrus.WithFields(logrus.Fields{
			"tracker": tracker,
		}).Error("Name is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name is required",
		})
	}

	if updateData.OldPassword == updateData.NewPassword {
		logrus.WithFields(logrus.Fields{
			"tracker": tracker,
		}).Error("New password cannot be same as old password")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "New password cannot be same as old password",
		})
	}

	if err := h.svc.UpdateUser(c, &updateData); err != nil {
		logrus.WithFields(logrus.Fields{
			"tracker": tracker,
			"err":     err.Error(),
		}).Error("Failed to update profile")
		if err.Error() == models.ErrInvalidCredentials || err.Error() == models.ErrUserNameExists {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update profile",
		})
	}

	logrus.WithFields(logrus.Fields{
		"tracker": tracker,
	}).Info("Profile updated successfully")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Profile updated successfully",
	})
}

// Logout logs out the user
func (h *handler) DeleteAccount(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint64)
	if err := h.svc.LogoutUser(c, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to logout user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User logged out successfully",
	})
}
