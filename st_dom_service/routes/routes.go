package routes

import (
	"st_dom_service/handlers"
	"st_dom_service/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(r *gin.Engine, stDomHandler *handlers.StDomHandler, sobaHandler *handlers.SobaHandler, aplikacijaHandler *handlers.AplikacijaHandler, prihvacenaAplikacijaHandler *handlers.PrihvacenaAplikacijaHandler, paymentHandler *handlers.PaymentHandler, healthHandler *handlers.HealthHandler, jwtSecret string) {
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

		// Inter-service communication endpoint (no auth required for service-to-service calls)
		interService := v1.Group("/internal")
		{
			interService.GET("/users/:userId/room-status", prihvacenaAplikacijaHandler.CheckUserRoomStatus) // Check if user has active room
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

			// User accepted applications routes
			prihvaceneAplikacije := user.Group("/prihvacene_aplikacije")
			{
				prihvaceneAplikacije.GET("/my", prihvacenaAplikacijaHandler.GetMyPrihvaceneAplikacije) // User gets their own accepted applications
				prihvaceneAplikacije.GET("/", prihvacenaAplikacijaHandler.GetAllPrihvaceneAplikacije)                    // Get all accepted applications (available to all authenticated users)
				prihvaceneAplikacije.GET("/user/:userId", prihvacenaAplikacijaHandler.GetPrihvaceneAplikacijeForUser)    // Get by user (available to all authenticated users)
				prihvaceneAplikacije.GET("/room/:sobaId", prihvacenaAplikacijaHandler.GetPrihvaceneAplikacijeForRoom)    // Get by room (available to all authenticated users)
				prihvaceneAplikacije.GET("/academic_year", prihvacenaAplikacijaHandler.GetPrihvaceneAplikacijeForAcademicYear) // Get by academic year (available to all authenticated users)
				prihvaceneAplikacije.POST("/checkout", prihvacenaAplikacijaHandler.CheckoutFromRoom)    // User voluntarily leaves room
			}

			// User payment routes
			payments := user.Group("/payments")
			{
				payments.GET("/my", paymentHandler.GetMyPayments)       // User gets their own payments
				payments.GET("/:id", paymentHandler.GetPayment)         // User gets their own, admin gets any
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

			// Admin accepted applications routes (Student ranking system)
			adminPrihvaceneAplikacije := admin.Group("/prihvacene_aplikacije")
			{
				adminPrihvaceneAplikacije.POST("/approve", prihvacenaAplikacijaHandler.ApproveAplikacija)                       // Approve application
				adminPrihvaceneAplikacije.POST("/evict", prihvacenaAplikacijaHandler.EvictStudent)                             // Evict student from room
				adminPrihvaceneAplikacije.GET("/:id", prihvacenaAplikacijaHandler.GetPrihvacenaAplikacija)                    // Get accepted application by ID
				adminPrihvaceneAplikacije.GET("/ranking/top", prihvacenaAplikacijaHandler.GetTopStudentsByProsek)             // Get top students overall
				adminPrihvaceneAplikacije.GET("/ranking/top/academic_year/:academicYear", prihvacenaAplikacijaHandler.GetTopStudentsByProsekForAcademicYear) // Get top students by year
				adminPrihvaceneAplikacije.GET("/ranking/top/room/:sobaId", prihvacenaAplikacijaHandler.GetTopStudentsByProsekForRoom) // Get top students by room
				adminPrihvaceneAplikacije.DELETE("/:id", prihvacenaAplikacijaHandler.DeletePrihvacenaAplikacija)               // Delete accepted application
			}

			// Admin payment routes
			adminPayments := admin.Group("/payments")
			{
				adminPayments.POST("/", paymentHandler.CreatePayment)                           // Create payment
				adminPayments.GET("/", paymentHandler.GetAllPayments)                          // Get all payments (with optional status filter)
				adminPayments.GET("/search", paymentHandler.SearchPaymentsByIndex)             // Search payments by student index pattern
				adminPayments.GET("/room/:sobaId", paymentHandler.GetPaymentsByRoom)           // Get payments by room
				adminPayments.GET("/user/:userId", paymentHandler.GetPaymentsByUser)           // Get payments by user
				adminPayments.GET("/aplikacija/:aplikacijaId", paymentHandler.GetPaymentsByAplikacija) // Get payments by application
				adminPayments.PUT("/:id", paymentHandler.UpdatePayment)                        // Update payment
				adminPayments.PATCH("/:id/mark-paid", paymentHandler.MarkPaymentAsPaid)        // Mark payment as paid
				adminPayments.PATCH("/:id/mark-unpaid", paymentHandler.MarkPaymentAsUnpaid)    // Mark payment as unpaid
				adminPayments.DELETE("/:id", paymentHandler.DeletePayment)                     // Delete payment
				adminPayments.POST("/update-overdue", paymentHandler.UpdateOverduePayments)    // Update overdue payments
			}
		}
	}
}
