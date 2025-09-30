package routes

import (
	"st_dom_service/handlers"
	"st_dom_service/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(r *gin.Engine, stDomHandler *handlers.StDomHandler, sobaHandler *handlers.SobaHandler, aplikacijaHandler *handlers.AplikacijaHandler, healthHandler *handlers.HealthHandler, jwtSecret string) {
	// Add CORS middleware
	r.Use(middleware.CORSMiddleware())

	// Health check endpoint
	r.GET("/health", healthHandler.Health)

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Public student dormitory routes (read-only)
		stDoms := v1.Group("/st_doms")
		{
			stDoms.GET("/", stDomHandler.GetAllStDoms)
			stDoms.GET("/:id", stDomHandler.GetStDom)
			stDoms.GET("/:id/rooms", stDomHandler.GetStDomRooms)
		}

		// Public room routes (read-only)
		sobas := v1.Group("/sobas")
		{
			sobas.GET("/", sobaHandler.GetAllSobas)
			sobas.GET("/:id", sobaHandler.GetSoba)
		}

		// User routes (authentication required)
		user := v1.Group("/")
		user.Use(middleware.AuthMiddleware(jwtSecret))
		{
			// User application routes
			aplikacije := user.Group("/aplikacije")
			{
				aplikacije.POST("/", aplikacijaHandler.CreateAplikacija)       // User only
				aplikacije.GET("/my", aplikacijaHandler.GetMyAplikacije)       // User gets their own
				aplikacije.GET("/:id", aplikacijaHandler.GetAplikacija)        // User gets their own, admin gets any
				aplikacije.PUT("/:id", aplikacijaHandler.UpdateAplikacija)     // User updates their own
				aplikacije.DELETE("/:id", aplikacijaHandler.DeleteAplikacija)  // User deletes their own
			}
		}

		// Admin-only routes (authentication + admin role required)
		admin := v1.Group("/")
		admin.Use(middleware.AuthMiddleware(jwtSecret))
		admin.Use(middleware.RoleMiddleware("admin"))
		{
			// Admin student dormitory routes
			adminStDoms := admin.Group("/st_doms")
			{
				adminStDoms.POST("/", stDomHandler.CreateStDom)
				adminStDoms.PUT("/:id", stDomHandler.UpdateStDom)
				adminStDoms.DELETE("/:id", stDomHandler.DeleteStDom)
			}

			// Admin room routes
			adminSobas := admin.Group("/sobas")
			{
				adminSobas.POST("/", sobaHandler.CreateSoba)
				adminSobas.PUT("/:id", sobaHandler.UpdateSoba)
				adminSobas.DELETE("/:id", sobaHandler.DeleteSoba)
			}

			// Admin application routes
			adminAplikacije := admin.Group("/aplikacije")
			{
				adminAplikacije.GET("/", aplikacijaHandler.GetAllAplikacije)           // Admin gets all
				adminAplikacije.GET("/room/:sobaId", aplikacijaHandler.GetAplikacijeForRoom) // Admin gets by room
			}
		}
	}
}
