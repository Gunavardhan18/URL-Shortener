package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type ConnectionConfig struct {
	DB_URL    string
	RedisHost string
	Port      string
}

// Load environment variables
func Load() ConnectionConfig {
	viper.SetConfigName("env")         // File name (without extension)
	viper.SetConfigType("yaml")        // File type
	viper.AddConfigPath("../config/.") // Look in the current directory

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	dbURL := viper.GetString("DB_URL")
	redisAddress := viper.GetString("REDIS_ADDR")
	port := viper.GetString("PORT")
	return ConnectionConfig{
		DB_URL:    dbURL,
		RedisHost: redisAddress,
		Port:      port,
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
