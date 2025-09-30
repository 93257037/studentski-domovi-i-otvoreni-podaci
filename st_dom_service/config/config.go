package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration values
type Config struct {
	MongoDBURI   string
	DatabaseName string
	JWTSecret    string
	Port         string
	GinMode      string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Try to load from config.env file
	if err := godotenv.Load("config.env"); err != nil {
		log.Println("No config.env file found, using environment variables")
	}

	config := &Config{
		MongoDBURI:   getEnv("MONGODB_URI", "mongodb://localhost:27018"),
		DatabaseName: getEnv("DATABASE_NAME", "st_dom_db"),
		JWTSecret:    getEnv("JWT_SECRET", "default_jwt_secret_change_in_production"),
		Port:         getEnv("PORT", "8081"),
		GinMode:      getEnv("GIN_MODE", "debug"),
	}

	return config
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
