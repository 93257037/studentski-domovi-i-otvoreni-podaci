package routes

import (
	"sso_service/handlers"
	"sso_service/middleware"

	"github.com/gin-gonic/gin"
)

// postavlja sve rute za aplikaciju - javne i zasticene
// javne rute su za registraciju i prijavu, zasticene zahtevaju JWT token
func SetupRoutes(r *gin.Engine, authHandler *handlers.AuthHandler, jwtSecret string) {
	r.Use(middleware.CORSMiddleware())

	r.GET("/health", authHandler.Health)

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware(jwtSecret))
		{
			protected.GET("/profile", authHandler.GetProfile)
			protected.DELETE("/account", authHandler.DeleteAccount)
		}
	}
}