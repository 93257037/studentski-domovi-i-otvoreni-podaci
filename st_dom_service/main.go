package main

import (
	"log"
	"st_dom_service/config"
	"st_dom_service/database"
	"st_dom_service/handlers"
	"st_dom_service/routes"
	"st_dom_service/services"

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

	// Get collections
	stDomsCollection := db.GetCollection("st_doms")
	sobasCollection := db.GetCollection("sobas")
	aplikacijeCollection := db.GetCollection("aplikacije")
	prihvaceneAplikacijeCollection := db.GetCollection("prihvacene_aplikacije")
	paymentsCollection := db.GetCollection("payments")

	// Initialize services
	stDomService := services.NewStDomService(stDomsCollection)
	sobaService := services.NewSobaService(sobasCollection, prihvaceneAplikacijeCollection)
	aplikacijaService := services.NewAplikacijaService(aplikacijeCollection)
	paymentService := services.NewPaymentService(paymentsCollection)
	prihvacenaAplikacijaService := services.NewPrihvacenaAplikacijaService(prihvaceneAplikacijeCollection, aplikacijaService, paymentService)

	// Initialize handlers
	stDomHandler := handlers.NewStDomHandler(stDomService, sobaService)
	sobaHandler := handlers.NewSobaHandler(sobaService, stDomService)
	aplikacijaHandler := handlers.NewAplikacijaHandler(aplikacijaService, sobaService)
	prihvacenaAplikacijaHandler := handlers.NewPrihvacenaAplikacijaHandler(prihvacenaAplikacijaService)
	paymentHandler := handlers.NewPaymentHandler(paymentService, aplikacijaService, sobaService)
	healthHandler := handlers.NewHealthHandler()

	// Initialize Gin router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router, stDomHandler, sobaHandler, aplikacijaHandler, prihvacenaAplikacijaHandler, paymentHandler, healthHandler, cfg.JWTSecret)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
