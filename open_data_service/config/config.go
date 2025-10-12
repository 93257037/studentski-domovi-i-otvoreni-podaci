package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config - holds all configuration values for the service
type Config struct {
	MongoDBURI      string
	DatabaseName    string
	Port            string
	GinMode         string
	StDomServiceURL string
}

// LoadConfig loads configuration from environment variables or config.env file
func LoadConfig() *Config {
	if err := godotenv.Load("config.env"); err != nil {
		log.Println("No config.env file found, using environment variables")
	}

	config := &Config{
		MongoDBURI:      getEnv("MONGODB_URI", "mongodb://localhost:27018"),
		DatabaseName:    getEnv("DATABASE_NAME", "st_dom_db"),
		Port:            getEnv("PORT", "8082"),
		GinMode:         getEnv("GIN_MODE", "debug"),
		StDomServiceURL: getEnv("ST_DOM_SERVICE_URL", "http://localhost:8081"),
	}

	return config
}

// getEnv gets an environment variable or returns a fallback value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
