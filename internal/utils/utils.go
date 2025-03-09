package utils

import (
	"crypto/sha256"
	"encoding/base64"
)

func GenerateHash(input string) string {
	hash := sha256.Sum256([]byte(input))
	return base64.URLEncoding.EncodeToString(hash[:])[:8]
}

func HashPassword(password string) (string, error) {
	hash := sha256.Sum256([]byte(password))
	return base64.URLEncoding.EncodeToString(hash[:]), nil
}

func ComparePasswords(hashedPassword, password string) bool {
	return hashedPassword == password
}

func GenerateJWT(userID uint64) (string, error) {
	return "", nil
}
