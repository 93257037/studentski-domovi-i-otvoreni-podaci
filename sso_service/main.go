package main

import (
	"log"
	"sso_service/config"
	"sso_service/database"
	"sso_service/handlers"
	"sso_service/routes"
	"sso_service/services"

	"github.com/gin-gonic/gin"
)

// glavna funkcija - pokretanje SSO servisa
// ucitava konfiguraciju, povezuje se sa bazom, postavlja rute i pokretanje servera
func main() {
	cfg := config.LoadConfig()

	gin.SetMode(cfg.GinMode)

	db, err := database.NewMongoDB(cfg.MongoDBURI, cfg.DatabaseName)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Println("Error closing MongoDB connection:", err)
		}
	}()

	usersCollection := db.GetCollection("users")

	userService := services.NewUserService(usersCollection, cfg.JWTSecret, cfg.StDomServiceURL)

	authHandler := handlers.NewAuthHandler(userService)

	router := gin.Default()

	routes.SetupRoutes(router, authHandler, cfg.JWTSecret)

	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

