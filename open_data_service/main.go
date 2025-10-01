package main

import (
	"log"
	"open_data_service/config"
	"open_data_service/database"
	"open_data_service/handlers"
	"open_data_service/routes"
	"open_data_service/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Set Gin mode
	gin.SetMode(cfg.GinMode)

	// Connect to MongoDB (same database as st_dom_service)
	db, err := database.NewMongoDB(cfg.MongoDBURI, cfg.DatabaseName)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Println("Error closing MongoDB connection:", err)
		}
	}()

	// Get collections
	sobasCollection := db.GetCollection("sobas")
	stDomsCollection := db.GetCollection("st_doms")
	aplikacijeCollection := db.GetCollection("aplikacije")
	prihvaceneAplikacijeCollection := db.GetCollection("prihvacene_aplikacije")

	// Initialize services
	openDataService := services.NewOpenDataService(sobasCollection, stDomsCollection, aplikacijeCollection, prihvaceneAplikacijeCollection)
	
	// Initialize HTTP client service for inter-service communication
	httpClientService := services.NewHTTPClientService(cfg.StDomServiceURL)

	// Initialize handlers
	openDataHandler := handlers.NewOpenDataHandler(openDataService, httpClientService)

	// Initialize Gin router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router, openDataHandler)

	// Start server
	log.Printf("Open Data Service starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

