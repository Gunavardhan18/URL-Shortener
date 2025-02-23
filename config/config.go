package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Load environment variables
func Load() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
}

// Get port from .env or default
func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}
