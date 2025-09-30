package handlers

import (
	"net/http"
	"st_dom_service/models"
	"st_dom_service/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PaymentHandler handles payment-related requests
type PaymentHandler struct {
	paymentService    *services.PaymentService
	aplikacijaService *services.AplikacijaService
	sobaService       *services.SobaService
}

// NewPaymentHandler creates a new PaymentHandler
func NewPaymentHandler(paymentService *services.PaymentService, aplikacijaService *services.AplikacijaService, sobaService *services.SobaService) *PaymentHandler {
	return &PaymentHandler{
		paymentService:    paymentService,
		aplikacijaService: aplikacijaService,
		sobaService:       sobaService,
	}
}

// CreatePayment handles creating a new payment (admin only)
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var req models.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify that the application exists
	aplikacija, err := h.aplikacijaService.GetAplikacijaByID(req.AplikacijaID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Application not found"})
		return
	}

	// Check if application is active
	if !aplikacija.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Application is not active"})
		return
	}

	payment, err := h.paymentService.CreatePayment(req, aplikacija)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Payment created successfully",
		"payment": payment,
	})
}

// GetPayment handles retrieving a payment by ID
func (h *PaymentHandler) GetPayment(c *gin.Context) {
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

	payment, err := h.paymentService.GetPaymentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Users can only see their own payments, admins can see all
	if userRole == "user" {
		// Get the associated application to check ownership
		aplikacija, err := h.aplikacijaService.GetAplikacijaByID(payment.AplikacijaID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify payment ownership"})
			return
		}
		if aplikacija.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: payment does not belong to user"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"payment": payment,
	})
}

// GetMyPayments handles retrieving all payments for the current user
func (h *PaymentHandler) GetMyPayments(c *gin.Context) {
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

	payments, err := h.paymentService.GetPaymentsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payments": payments,
	})
}

// GetAllPayments handles retrieving all payments (admin only)
func (h *PaymentHandler) GetAllPayments(c *gin.Context) {
	// Check for status filter
	statusParam := c.Query("status")
	
	var payments []models.Payment
	var err error

	if statusParam != "" {
		status := models.PaymentStatus(statusParam)
		if !status.IsValid() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment status"})
			return
		}
		payments, err = h.paymentService.GetPaymentsByStatus(status)
	} else {
		payments, err = h.paymentService.GetAllPayments()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payments": payments,
	})
}

// GetPaymentsByRoom handles retrieving all payments for a specific room (admin only)
func (h *PaymentHandler) GetPaymentsByRoom(c *gin.Context) {
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

	payments, err := h.paymentService.GetPaymentsBySobaID(sobaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payments": payments,
	})
}

// GetPaymentsByUser handles retrieving all payments for a specific user (admin only)
func (h *PaymentHandler) GetPaymentsByUser(c *gin.Context) {
	idParam := c.Param("userId")
	userID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	payments, err := h.paymentService.GetPaymentsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payments": payments,
	})
}

// GetPaymentsByAplikacija handles retrieving all payments for a specific application (admin only)
func (h *PaymentHandler) GetPaymentsByAplikacija(c *gin.Context) {
	idParam := c.Param("aplikacijaId")
	aplikacijaID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application ID format"})
		return
	}

	// Check if the application exists
	_, err = h.aplikacijaService.GetAplikacijaByID(aplikacijaID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Application not found"})
		return
	}

	payments, err := h.paymentService.GetPaymentsByAplikacijaID(aplikacijaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payments": payments,
	})
}

// UpdatePayment handles updating a payment (admin only)
func (h *PaymentHandler) UpdatePayment(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req models.UpdatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.paymentService.UpdatePayment(id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Payment updated successfully",
		"payment": payment,
	})
}

// MarkPaymentAsPaid handles marking a payment as paid (admin only)
func (h *PaymentHandler) MarkPaymentAsPaid(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req models.MarkPaymentPaidRequest
	// Don't return error if body is empty, just use default (current time)
	_ = c.ShouldBindJSON(&req)

	payment, err := h.paymentService.MarkPaymentAsPaid(id, req.PaidAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Payment marked as paid",
		"payment": payment,
	})
}

// MarkPaymentAsUnpaid handles marking a payment as unpaid/pending (admin only)
func (h *PaymentHandler) MarkPaymentAsUnpaid(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	payment, err := h.paymentService.MarkPaymentAsUnpaid(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Payment marked as unpaid",
		"payment": payment,
	})
}

// DeletePayment handles deleting a payment (admin only)
func (h *PaymentHandler) DeletePayment(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	err = h.paymentService.DeletePayment(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Payment deleted successfully",
	})
}

// UpdateOverduePayments handles updating overdue payments (admin only)
func (h *PaymentHandler) UpdateOverduePayments(c *gin.Context) {
	count, err := h.paymentService.UpdateOverduePayments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Overdue payments updated",
		"updated_count": count,
	})
}
