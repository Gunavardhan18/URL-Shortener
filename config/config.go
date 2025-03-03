package config

import (
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

type ConnectionConfig struct {
	DB_URL    string
	RedisHost string
	Port      int
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
	portStr := viper.GetString("PORT")

	port, _ := strconv.ParseInt(portStr, 10, 16)
	return ConnectionConfig{
		DB_URL:    dbURL,
		RedisHost: redisAddress,
		Port:      int(port),
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
