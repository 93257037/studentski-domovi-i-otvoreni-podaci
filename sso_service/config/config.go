package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config - drzi sve konfiguracione vrednosti za servis
type Config struct {
	MongoDBURI        string
	DatabaseName      string
	JWTSecret         string
	Port              string
	GinMode           string
	StDomServiceURL   string
}

// ucitava konfiguraciju iz environment varijabli ili config.env fajla
// postavlja default vrednosti ako varijable nisu definisane
func LoadConfig() *Config {
	if err := godotenv.Load("config.env"); err != nil {
		log.Println("No config.env file found, using environment variables")
	}

	config := &Config{
		MongoDBURI:      getEnv("MONGODB_URI", "mongodb://localhost:27018"),
		DatabaseName:    getEnv("DATABASE_NAME", "sso_db"),
		JWTSecret:       getEnv("JWT_SECRET", "default_jwt_secret_change_in_production"),
		Port:            getEnv("PORT", "8080"),
		GinMode:         getEnv("GIN_MODE", "debug"),
		StDomServiceURL: getEnv("ST_DOM_SERVICE_URL", "http://localhost:8081"),
	}

	return config
}

// dobija environment varijablu ili vraca default vrednost ako ne postoji
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

