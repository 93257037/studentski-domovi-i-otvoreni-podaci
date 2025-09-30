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

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Set Gin mode
	gin.SetMode(cfg.GinMode)

	// Connect to MongoDB
	db, err := database.NewMongoDB(cfg.MongoDBURI, cfg.DatabaseName)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Println("Error closing MongoDB connection:", err)
		}
	}()

	// Get users collection
	usersCollection := db.GetCollection("users")

	// Initialize services
	userService := services.NewUserService(usersCollection, cfg.JWTSecret, cfg.StDomServiceURL)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userService)

	// Initialize Gin router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router, authHandler, cfg.JWTSecret)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

