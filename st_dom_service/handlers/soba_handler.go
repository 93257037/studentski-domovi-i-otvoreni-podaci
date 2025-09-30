package handlers

import (
	"net/http"
	"st_dom_service/models"
	"st_dom_service/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SobaHandler handles room-related requests
type SobaHandler struct {
	sobaService  *services.SobaService
	stDomService *services.StDomService
}

// NewSobaHandler creates a new SobaHandler
func NewSobaHandler(sobaService *services.SobaService, stDomService *services.StDomService) *SobaHandler {
	return &SobaHandler{
		sobaService:  sobaService,
		stDomService: stDomService,
	}
}

// CreateSoba handles creating a new room
func (h *SobaHandler) CreateSoba(c *gin.Context) {
	var req models.CreateSobaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the dormitory exists
	_, err := h.stDomService.GetStDomByID(req.StDomID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Student dormitory not found"})
		return
	}

	soba, err := h.sobaService.CreateSoba(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Room created successfully",
		"soba":    soba,
	})
}

// GetSoba handles retrieving a room by ID
func (h *SobaHandler) GetSoba(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	soba, err := h.sobaService.GetSobaByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"soba": soba,
	})
}

// GetAllSobas handles retrieving all rooms
func (h *SobaHandler) GetAllSobas(c *gin.Context) {
	sobas, err := h.sobaService.GetAllSobas()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sobas": sobas,
	})
}

// UpdateSoba handles updating a room
func (h *SobaHandler) UpdateSoba(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req models.UpdateSobaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	soba, err := h.sobaService.UpdateSoba(id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Room updated successfully",
		"soba":    soba,
	})
}

// DeleteSoba handles deleting a room
func (h *SobaHandler) DeleteSoba(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	err = h.sobaService.DeleteSoba(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Room deleted successfully",
	})
}
