package handlers

import (
	"net/http"
	"st_dom_service/models"
	"st_dom_service/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AplikacijaHandler handles application-related requests
type AplikacijaHandler struct {
	aplikacijaService *services.AplikacijaService
	sobaService       *services.SobaService
}

// NewAplikacijaHandler creates a new AplikacijaHandler
func NewAplikacijaHandler(aplikacijaService *services.AplikacijaService, sobaService *services.SobaService) *AplikacijaHandler {
	return &AplikacijaHandler{
		aplikacijaService: aplikacijaService,
		sobaService:       sobaService,
	}
}

// CreateAplikacija handles creating a new application (user only)
func (h *AplikacijaHandler) CreateAplikacija(c *gin.Context) {
	// Extract user ID from JWT token
	userIDClaim, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	userID, ok := userIDClaim.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Extract user role from JWT token
	userRole, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found in token"})
		return
	}

	if userRole != "user" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only users can create applications"})
		return
	}

	var req models.CreateAplikacijaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the room exists
	_, err := h.sobaService.GetSobaByID(req.SobaID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room not found"})
		return
	}

	aplikacija, err := h.aplikacijaService.CreateAplikacija(req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Application created successfully",
		"aplikacija": aplikacija,
	})
}

// GetAplikacija handles retrieving an application by ID
func (h *AplikacijaHandler) GetAplikacija(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Extract user ID and role from JWT token
	userIDClaim, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	userID, ok := userIDClaim.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
		return
	}

	userRole, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found in token"})
		return
	}

	aplikacija, err := h.aplikacijaService.GetAplikacijaByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Users can only see their own applications, admins can see all
	if userRole == "user" && aplikacija.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: application does not belong to user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"aplikacija": aplikacija,
	})
}

// GetMyAplikacije handles retrieving all applications for the current user
func (h *AplikacijaHandler) GetMyAplikacije(c *gin.Context) {
	// Extract user ID from JWT token
	userIDClaim, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	userIDStr, ok := userIDClaim.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	aplikacije, err := h.aplikacijaService.GetAplikacijeByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"aplikacije": aplikacije,
	})
}

// GetAllAplikacije handles retrieving all applications (admin only)
func (h *AplikacijaHandler) GetAllAplikacije(c *gin.Context) {
	// Extract user role from JWT token
	userRole, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found in token"})
		return
	}

	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	aplikacije, err := h.aplikacijaService.GetAllAplikacije()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"aplikacije": aplikacije,
	})
}

// GetAplikacijeForRoom handles retrieving all applications for a specific room (admin only)
func (h *AplikacijaHandler) GetAplikacijeForRoom(c *gin.Context) {
	// Extract user role from JWT token
	userRole, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found in token"})
		return
	}

	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	idParam := c.Param("sobaId")
	sobaID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID format"})
		return
	}

	// Check if the room exists
	_, err = h.sobaService.GetSobaByID(sobaID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room not found"})
		return
	}

	aplikacije, err := h.aplikacijaService.GetAplikacijeBySobaID(sobaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"aplikacije": aplikacije,
	})
}

// UpdateAplikacija handles updating an application (user can update their own)
func (h *AplikacijaHandler) UpdateAplikacija(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Extract user ID from JWT token
	userIDClaim, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	userIDStr, ok := userIDClaim.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	var req models.UpdateAplikacijaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	aplikacija, err := h.aplikacijaService.UpdateAplikacija(id, req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Application updated successfully",
		"aplikacija": aplikacija,
	})
}

// DeleteAplikacija handles deleting an application (user can delete their own)
func (h *AplikacijaHandler) DeleteAplikacija(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Extract user ID from JWT token
	userIDClaim, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	userIDStr, ok := userIDClaim.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.aplikacijaService.DeleteAplikacija(id, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Application deleted successfully",
	})
}