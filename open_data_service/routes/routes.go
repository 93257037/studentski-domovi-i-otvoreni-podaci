package routes

import (
	"open_data_service/handlers"
	"open_data_service/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all routes for the open data service
func SetupRoutes(r *gin.Engine, handler *handlers.OpenDataHandler) {
	// Add CORS middleware
	r.Use(middleware.CORSMiddleware())

	// Health check endpoint
	r.GET("/health", handler.Health)

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Room filtering endpoints
		rooms := v1.Group("/rooms")
		{
			rooms.GET("", handler.GetAllRooms)                               // Get all rooms
			rooms.GET("/filter-by-luksuz", handler.FilterRoomsByLuksuz)      // Filter by luxury amenities
			rooms.GET("/filter-by-luksuz-and-stdom", handler.FilterRoomsByLuksuzAndStDom) // Filter by luxury and dorm
			rooms.GET("/filter-by-krevetnost", handler.FilterRoomsByKrevetnost)           // Filter by bed capacity
			rooms.GET("/advanced-filter", handler.AdvancedFilterRooms)                    // Advanced multi-criteria filter
		}

		// Student dormitory endpoints
		stDoms := v1.Group("/st-doms")
		{
			stDoms.GET("", handler.GetAllStDoms)                          // Get all student dormitories
			stDoms.GET("/search-by-address", handler.SearchStDomsByAddress) // Search by address (regex)
			stDoms.GET("/search-by-ime", handler.SearchStDomsByIme)         // Search by name (regex)
		}
	}
}

