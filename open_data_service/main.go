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

// main starts the open_data_service
// loads configuration, connects to database, creates services and handlers, sets up routes
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

	// Get collections
	stDomsCollection := db.GetCollection("st_doms")
	sobasCollection := db.GetCollection("sobas")
	aplikacijeCollection := db.GetCollection("aplikacije")
	prihvaceneAplikacijeCollection := db.GetCollection("prihvacene_aplikacije")

	// Create services
	openDataService := services.NewOpenDataService(
		stDomsCollection,
		sobasCollection,
		aplikacijeCollection,
		prihvaceneAplikacijeCollection,
	)

	// Create handlers
	openDataHandler := handlers.NewOpenDataHandler(openDataService)
	healthHandler := handlers.NewHealthHandler()

	// Create router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router, openDataHandler, healthHandler)

	log.Printf("Open Data Service starting on port %s", cfg.Port)
	log.Println("Available endpoints:")
	log.Println("  GET /health")
	log.Println("  GET /api/v1/open-data/statistics")
	log.Println("  GET /api/v1/open-data/rooms/search")
	log.Println("  GET /api/v1/open-data/rooms/:roomId/applications")
	log.Println("  GET /api/v1/open-data/dorms/compare")
	log.Println("  GET /api/v1/open-data/dorms/list")
	log.Println("  GET /api/v1/open-data/trends/applications")
	log.Println("  GET /api/v1/open-data/occupancy/heatmap")
	log.Println("  GET /api/v1/open-data/export")
	log.Println("  GET /api/v1/open-data/amenities")

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
