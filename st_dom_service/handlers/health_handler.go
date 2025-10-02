package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler - rukuje zahtevima za proveru zdravlja servisa
type HealthHandler struct{}

// kreira novi HealthHandler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// provera zdravlja servisa - vraca status da li je servis aktivan
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "st_dom_service",
	})
}
