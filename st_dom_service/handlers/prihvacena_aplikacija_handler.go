package handlers

import (
	"net/http"
	"st_dom_service/models"
	"st_dom_service/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PrihvacenaAplikacijaHandler handles accepted application-related requests
type PrihvacenaAplikacijaHandler struct {
	prihvacenaAplikacijaService *services.PrihvacenaAplikacijaService
}

// NewPrihvacenaAplikacijaHandler creates a new PrihvacenaAplikacijaHandler
func NewPrihvacenaAplikacijaHandler(prihvacenaAplikacijaService *services.PrihvacenaAplikacijaService) *PrihvacenaAplikacijaHandler {
	return &PrihvacenaAplikacijaHandler{
		prihvacenaAplikacijaService: prihvacenaAplikacijaService,
	}
}

// ApproveAplikacija handles approving an application (admin only)
func (h *PrihvacenaAplikacijaHandler) ApproveAplikacija(c *gin.Context) {
	var req models.ApproveAplikacijaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prihvacenaAplikacija, err := h.prihvacenaAplikacijaService.ApproveAplikacija(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":              "Application approved successfully",
		"prihvacena_aplikacija": prihvacenaAplikacija,
	})
}

// GetPrihvacenaAplikacija handles retrieving an accepted application by ID
func (h *PrihvacenaAplikacijaHandler) GetPrihvacenaAplikacija(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	prihvacenaAplikacija, err := h.prihvacenaAplikacijaService.GetPrihvacenaAplikacijaByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"prihvacena_aplikacija": prihvacenaAplikacija,
	})
}

// GetAllPrihvaceneAplikacije handles retrieving all accepted applications
func (h *PrihvacenaAplikacijaHandler) GetAllPrihvaceneAplikacije(c *gin.Context) {
	prihvaceneAplikacije, err := h.prihvacenaAplikacijaService.GetAllPrihvaceneAplikacije()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"prihvacene_aplikacije": prihvaceneAplikacije,
		"count":                  len(prihvaceneAplikacije),
	})
}

// GetPrihvaceneAplikacijeForUser handles retrieving accepted applications for a user
func (h *PrihvacenaAplikacijaHandler) GetPrihvaceneAplikacijeForUser(c *gin.Context) {
	userIDParam := c.Param("userId")
	userID, err := primitive.ObjectIDFromHex(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	prihvaceneAplikacije, err := h.prihvacenaAplikacijaService.GetPrihvaceneAplikacijeByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"prihvacene_aplikacije": prihvaceneAplikacije,
		"count":                  len(prihvaceneAplikacije),
	})
}

// GetPrihvaceneAplikacijeForRoom handles retrieving accepted applications for a room
func (h *PrihvacenaAplikacijaHandler) GetPrihvaceneAplikacijeForRoom(c *gin.Context) {
	sobaIDParam := c.Param("sobaId")
	sobaID, err := primitive.ObjectIDFromHex(sobaIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID format"})
		return
	}

	prihvaceneAplikacije, err := h.prihvacenaAplikacijaService.GetPrihvaceneAplikacijeBySobaID(sobaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"prihvacene_aplikacije": prihvaceneAplikacije,
		"count":                  len(prihvaceneAplikacije),
	})
}

// GetPrihvaceneAplikacijeForAcademicYear handles retrieving accepted applications for an academic year
func (h *PrihvacenaAplikacijaHandler) GetPrihvaceneAplikacijeForAcademicYear(c *gin.Context) {
	academicYear := c.Param("academicYear")
	if academicYear == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Academic year is required"})
		return
	}

	prihvaceneAplikacije, err := h.prihvacenaAplikacijaService.GetPrihvaceneAplikacijeByAcademicYear(academicYear)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"prihvacene_aplikacije": prihvaceneAplikacije,
		"count":                  len(prihvaceneAplikacije),
	})
}

// GetTopStudentsByProsek handles retrieving top N students ranked by prosek
func (h *PrihvacenaAplikacijaHandler) GetTopStudentsByProsek(c *gin.Context) {
	// Get limit from query parameter, default to 10
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	topStudents, err := h.prihvacenaAplikacijaService.GetTopStudentsByProsek(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"top_students": topStudents,
		"count":        len(topStudents),
		"limit":        limit,
	})
}

// GetTopStudentsByProsekForAcademicYear handles retrieving top N students for a specific academic year
func (h *PrihvacenaAplikacijaHandler) GetTopStudentsByProsekForAcademicYear(c *gin.Context) {
	academicYear := c.Param("academicYear")
	if academicYear == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Academic year is required"})
		return
	}

	// Get limit from query parameter, default to 10
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	topStudents, err := h.prihvacenaAplikacijaService.GetTopStudentsByProsekForAcademicYear(academicYear, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"top_students":   topStudents,
		"count":          len(topStudents),
		"limit":          limit,
		"academic_year": academicYear,
	})
}

// GetTopStudentsByProsekForRoom handles retrieving top N students for a specific room
func (h *PrihvacenaAplikacijaHandler) GetTopStudentsByProsekForRoom(c *gin.Context) {
	sobaIDParam := c.Param("sobaId")
	sobaID, err := primitive.ObjectIDFromHex(sobaIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID format"})
		return
	}

	// Get limit from query parameter, default to 10
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	topStudents, err := h.prihvacenaAplikacijaService.GetTopStudentsByProsekForRoom(sobaID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"top_students": topStudents,
		"count":        len(topStudents),
		"limit":        limit,
		"soba_id":      sobaID,
	})
}

// GetMyPrihvaceneAplikacije handles retrieving accepted applications for the current user
func (h *PrihvacenaAplikacijaHandler) GetMyPrihvaceneAplikacije(c *gin.Context) {
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

	prihvaceneAplikacije, err := h.prihvacenaAplikacijaService.GetPrihvaceneAplikacijeByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"prihvacene_aplikacije": prihvaceneAplikacije,
		"count":                  len(prihvaceneAplikacije),
	})
}

// DeletePrihvacenaAplikacija handles deleting an accepted application (admin only)
func (h *PrihvacenaAplikacijaHandler) DeletePrihvacenaAplikacija(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	err = h.prihvacenaAplikacijaService.DeletePrihvacenaAplikacija(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Accepted application deleted successfully",
	})
}

// EvictStudent handles evicting a student from their room (admin only)
func (h *PrihvacenaAplikacijaHandler) EvictStudent(c *gin.Context) {
	var req models.EvictStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.prihvacenaAplikacijaService.EvictStudent(req.UserID, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Student evicted successfully",
		"user_id": req.UserID,
		"reason":  req.Reason,
	})
}

// CheckoutFromRoom handles a student voluntarily leaving their room
func (h *PrihvacenaAplikacijaHandler) CheckoutFromRoom(c *gin.Context) {
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

	err := h.prihvacenaAplikacijaService.CheckoutStudent(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully checked out from room",
	})
}

// CheckUserRoomStatus checks if a user has an active room assignment (for inter-service communication)
func (h *PrihvacenaAplikacijaHandler) CheckUserRoomStatus(c *gin.Context) {
	userIDParam := c.Param("userId")
	userID, err := primitive.ObjectIDFromHex(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	hasRoom, err := h.prihvacenaAplikacijaService.CheckUserHasActiveRoom(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"has_active_room": hasRoom,
	})
}