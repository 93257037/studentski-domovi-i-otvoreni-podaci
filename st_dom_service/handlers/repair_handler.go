package handlers

import (
	"net/http"
	"st_dom_service/models"
	"st_dom_service/services"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RepairHandler struct {
	repairService *services.RepairService
}

func NewRepairHandler(repairService *services.RepairService) *RepairHandler {
	return &RepairHandler{
		repairService: repairService,
	}
}

// CreateRepair creates a new repair schedule
// POST /api/v1/repairs (admin only)
func (h *RepairHandler) CreateRepair(c *gin.Context) {
	var req models.CreateRepairRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse room ID
	sobaID, err := primitive.ObjectIDFromHex(req.SobaID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	// Parse estimated completion date
	estimatedCompletionDate, err := time.Parse(time.RFC3339, req.EstimatedCompletionDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use ISO 8601 format (YYYY-MM-DDTHH:MM:SSZ)"})
		return
	}

	// Get user ID from context (set by auth middleware)
	userIDClaim, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	createdBy, ok := userIDClaim.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
		return
	}

	repair, err := h.repairService.CreateRepair(c.Request.Context(), sobaID, req.Description, estimatedCompletionDate, createdBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Repair scheduled successfully",
		"repair":  repair,
	})
}

// GetAllRepairs retrieves all repairs
// GET /api/v1/repairs (admin only)
func (h *RepairHandler) GetAllRepairs(c *gin.Context) {
	repairs, err := h.repairService.GetAllRepairs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"repairs": repairs,
		"count":   len(repairs),
	})
}

// GetRepair retrieves a repair by ID
// GET /api/v1/repairs/:id (admin only)
func (h *RepairHandler) GetRepair(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repair ID"})
		return
	}

	repair, err := h.repairService.GetRepairByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"repair": repair})
}

// GetRepairsByRoom retrieves all repairs for a specific room
// GET /api/v1/repairs/room/:roomId (admin only)
func (h *RepairHandler) GetRepairsByRoom(c *gin.Context) {
	roomIDParam := c.Param("roomId")
	roomID, err := primitive.ObjectIDFromHex(roomIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	repairs, err := h.repairService.GetRepairsByRoom(c.Request.Context(), roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"repairs": repairs,
		"count":   len(repairs),
	})
}

// GetRepairsByStatus retrieves all repairs with a specific status
// GET /api/v1/repairs/status/:status (admin only)
func (h *RepairHandler) GetRepairsByStatus(c *gin.Context) {
	status := c.Param("status")

	repairs, err := h.repairService.GetRepairsByStatus(c.Request.Context(), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"repairs": repairs,
		"count":   len(repairs),
		"status":  status,
	})
}

// UpdateRepair updates a repair
// PUT /api/v1/repairs/:id (admin only)
func (h *RepairHandler) UpdateRepair(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repair ID"})
		return
	}

	var req models.UpdateRepairRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var estimatedCompletionDate *time.Time
	if req.EstimatedCompletionDate != "" {
		parsed, err := time.Parse(time.RFC3339, req.EstimatedCompletionDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}
		estimatedCompletionDate = &parsed
	}

	repair, err := h.repairService.UpdateRepair(c.Request.Context(), id, req.Description, estimatedCompletionDate, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Repair updated successfully",
		"repair":  repair,
	})
}

// DeleteRepair deletes a repair
// DELETE /api/v1/repairs/:id (admin only)
func (h *RepairHandler) DeleteRepair(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repair ID"})
		return
	}

	err = h.repairService.DeleteRepair(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Repair deleted successfully"})
}

