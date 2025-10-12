package routes

import (
	"open_data_service/handlers"
	"open_data_service/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all routes for the open data service
func SetupRoutes(
	router *gin.Engine,
	openDataHandler *handlers.OpenDataHandler,
	healthHandler *handlers.HealthHandler,
) {
	// Apply CORS middleware
	router.Use(middleware.CORS())

	// Health check endpoint
	router.GET("/health", healthHandler.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Open Data routes group
		openData := v1.Group("/open-data")
		{
			// 1. Public Statistics Dashboard
			openData.GET("/statistics", openDataHandler.GetPublicStatistics)

			// 2. Room Availability Search
			openData.GET("/rooms/search", openDataHandler.SearchAvailableRooms)

			// 3. Dorm Comparison Tool
			openData.GET("/dorms/compare", openDataHandler.CompareDorms)
			openData.GET("/dorms/list", openDataHandler.GetDormList)

			// 4. Application Trends Analysis
			openData.GET("/trends/applications", openDataHandler.GetApplicationTrends)

			// 5. Real-time Occupancy Heatmap
			openData.GET("/occupancy/heatmap", openDataHandler.GetOccupancyHeatmap)

			// 6. Open Data Export (CSV/JSON)
			openData.GET("/export", openDataHandler.ExportData)

			// Helper endpoints
			openData.GET("/amenities", openDataHandler.GetAvailableAmenities)
		}
	}
}
