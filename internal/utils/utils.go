package utils

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GenerateHash(input string) string {
	hash := sha256.Sum256([]byte(input))
	return base64.URLEncoding.EncodeToString(hash[:])[:8]
}

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func ComparePasswords(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func GetTracker(ctx *fiber.Ctx) string {
	tracker := ctx.GetRespHeader("X-Request-ID")
	if tracker == "" {
		tracker = uuid.New().String()
		ctx.Set("X-Request-ID", tracker)
		return tracker
	}
	return tracker
}
