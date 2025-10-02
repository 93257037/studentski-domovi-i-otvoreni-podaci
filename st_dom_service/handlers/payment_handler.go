package handlers

import (
	"net/http"
	"st_dom_service/models"
	"st_dom_service/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PaymentHandler - rukuje zahtevima vezanim za placanja
type PaymentHandler struct {
	paymentService    *services.PaymentService
	aplikacijaService *services.AplikacijaService
	sobaService       *services.SobaService
}

// kreira novi PaymentHandler sa potrebnim servisima
func NewPaymentHandler(paymentService *services.PaymentService, aplikacijaService *services.AplikacijaService, sobaService *services.SobaService) *PaymentHandler {
	return &PaymentHandler{
		paymentService:    paymentService,
		aplikacijaService: aplikacijaService,
		sobaService:       sobaService,
	}
}

// kreira novo placanje - samo administratori mogu kreirati placanja
// proverava da li aplikacija postoji i da li je aktivna
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var req models.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	aplikacija, err := h.aplikacijaService.GetAplikacijaByID(req.AplikacijaID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Application not found"})
		return
	}

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

// dobija placanje po ID-u - korisnici mogu videti samo svoja placanja
// administratori mogu videti sva placanja
func (h *PaymentHandler) GetPayment(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

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

	if userRole == "user" {
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

// dobija sva placanja trenutno ulogovanog korisnika
func (h *PaymentHandler) GetMyPayments(c *gin.Context) {
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

// dobija sva placanja - samo za administratore
// moze da filtrira po statusu placanja
func (h *PaymentHandler) GetAllPayments(c *gin.Context) {
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

// pretrazuje placanja po indeksu studenta - samo za administratore
// prima pattern indeksa i opciono status placanja
func (h *PaymentHandler) SearchPaymentsByIndex(c *gin.Context) {
	indexPattern := c.Query("index")
	if indexPattern == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Index pattern is required"})
		return
	}

	var statusPtr *models.PaymentStatus
	statusParam := c.Query("status")
	if statusParam != "" {
		status := models.PaymentStatus(statusParam)
		if !status.IsValid() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment status. Valid values: pending, paid, overdue"})
			return
		}
		statusPtr = &status
	}

	payments, err := h.paymentService.SearchPaymentsByIndex(indexPattern, statusPtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payments": payments,
		"count":    len(payments),
		"index_pattern": indexPattern,
		"status": statusParam,
	})
}

// dobija sva placanja za odredjenu sobu - samo za administratore
func (h *PaymentHandler) GetPaymentsByRoom(c *gin.Context) {
	idParam := c.Param("sobaId")
	sobaID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID format"})
		return
	}

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

// dobija sva placanja za odredjenog korisnika - samo za administratore
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

// dobija sva placanja za odredjenu aplikaciju - samo za administratore
func (h *PaymentHandler) GetPaymentsByAplikacija(c *gin.Context) {
	idParam := c.Param("aplikacijaId")
	aplikacijaID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application ID format"})
		return
	}

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

// azurira placanje - samo administratori mogu azurirati placanja
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

// oznacava placanje kao placeno - samo administratori
// moze da primi datum placanja ili koristi trenutno vreme
func (h *PaymentHandler) MarkPaymentAsPaid(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req models.MarkPaymentPaidRequest
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

// oznacava placanje kao neplaceno - samo administratori
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

// brise placanje - samo administratori mogu brisati placanja
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

// azurira zakasnela placanja - samo administratori
// prolazi kroz sva placanja i oznacava zakasnela
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
