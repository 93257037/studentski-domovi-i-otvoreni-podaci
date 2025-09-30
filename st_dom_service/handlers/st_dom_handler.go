package handlers

import (
	"net/http"
	"st_dom_service/models"
	"st_dom_service/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// StDomHandler handles student dormitory-related requests
type StDomHandler struct {
	stDomService *services.StDomService
	sobaService  *services.SobaService
}

// NewStDomHandler creates a new StDomHandler
func NewStDomHandler(stDomService *services.StDomService, sobaService *services.SobaService) *StDomHandler {
	return &StDomHandler{
		stDomService: stDomService,
		sobaService:  sobaService,
	}
}

// CreateStDom handles creating a new student dormitory
func (h *StDomHandler) CreateStDom(c *gin.Context) {
	var req models.CreateStDomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stDom, err := h.stDomService.CreateStDom(req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Student dormitory created successfully",
		"st_dom":  stDom,
	})
}

// GetStDom handles retrieving a student dormitory by ID
func (h *StDomHandler) GetStDom(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	stDom, err := h.stDomService.GetStDomByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"st_dom": stDom,
	})
}

// GetAllStDoms handles retrieving all student dormitories
func (h *StDomHandler) GetAllStDoms(c *gin.Context) {
	stDoms, err := h.stDomService.GetAllStDoms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"st_doms": stDoms,
	})
}

// UpdateStDom handles updating a student dormitory
func (h *StDomHandler) UpdateStDom(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req models.UpdateStDomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stDom, err := h.stDomService.UpdateStDom(id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Student dormitory updated successfully",
		"st_dom":  stDom,
	})
}

// DeleteStDom handles deleting a student dormitory
func (h *StDomHandler) DeleteStDom(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// First delete all rooms associated with this dormitory
	err = h.sobaService.DeleteSobasByStDomID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete associated rooms"})
		return
	}

	// Then delete the dormitory
	err = h.stDomService.DeleteStDom(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Student dormitory and associated rooms deleted successfully",
	})
}

// GetStDomRooms handles retrieving all rooms for a specific dormitory
func (h *StDomHandler) GetStDomRooms(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// First check if dormitory exists
	_, err = h.stDomService.GetStDomByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student dormitory not found"})
		return
	}

	// Get rooms
	sobas, err := h.sobaService.GetSobasByStDomID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sobas": sobas,
	})
}
