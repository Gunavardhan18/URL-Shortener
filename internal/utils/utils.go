package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"math/rand"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GenerateSalt() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	salt := make([]byte, 6)
	for i := range salt {
		salt[i] = charset[rand.Intn(len(charset))]
	}
	return string(salt)
}

// GenerateShortCode creates a short URL code using hashing + salt
func GenerateShortCode(url string) string {
	salt := GenerateSalt()                                // Generate a random salt
	hash := sha256.Sum256([]byte(url + salt))             // Hash (URL + Salt)
	return base64.URLEncoding.EncodeToString(hash[:])[:8] // Take first 8 chars
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
