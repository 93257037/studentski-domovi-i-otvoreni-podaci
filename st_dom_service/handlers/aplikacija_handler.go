package handlers

import (
	"net/http"
	"st_dom_service/models"
	"st_dom_service/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AplikacijaHandler - rukuje zahtevima vezanim za aplikacije za sobe
type AplikacijaHandler struct {
	aplikacijaService *services.AplikacijaService
	sobaService       *services.SobaService
}

// kreira novi AplikacijaHandler sa potrebnim servisima
func NewAplikacijaHandler(aplikacijaService *services.AplikacijaService, sobaService *services.SobaService) *AplikacijaHandler {
	return &AplikacijaHandler{
		aplikacijaService: aplikacijaService,
		sobaService:       sobaService,
	}
}

// kreira novu aplikaciju za sobu - samo korisnici mogu da kreiraju aplikacije
// izvlaci podatke iz JWT tokena i proverava da li je soba dostupna
func (h *AplikacijaHandler) CreateAplikacija(c *gin.Context) {
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

	if userRole != "user" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only users can create applications"})
		return
	}

	var req models.CreateAplikacijaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

// dobija aplikaciju po ID-u - korisnici mogu videti samo svoje aplikacije
// administratori mogu videti sve aplikacije
func (h *AplikacijaHandler) GetAplikacija(c *gin.Context) {
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

	aplikacija, err := h.aplikacijaService.GetAplikacijaByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if userRole == "user" && aplikacija.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: application does not belong to user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"aplikacija": aplikacija,
	})
}

// dobija sve aplikacije trenutno ulogovanog korisnika
func (h *AplikacijaHandler) GetMyAplikacije(c *gin.Context) {
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

	aplikacije, err := h.aplikacijaService.GetAplikacijeByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"aplikacije": aplikacije,
	})
}

// dobija sve aplikacije - samo za administratore
func (h *AplikacijaHandler) GetAllAplikacije(c *gin.Context) {
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

// dobija sve aplikacije za odredjenu sobu - samo za administratore
func (h *AplikacijaHandler) GetAplikacijeForRoom(c *gin.Context) {
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

// azurira aplikaciju - korisnik moze azurirati samo svoju aplikaciju
func (h *AplikacijaHandler) UpdateAplikacija(c *gin.Context) {
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

// brise aplikaciju - korisnik moze brisati svoju, admin moze brisati bilo koju
func (h *AplikacijaHandler) DeleteAplikacija(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	userRole, roleExists := c.Get("role")
	if !roleExists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found in token"})
		return
	}

	if userRole == "admin" {
		err = h.aplikacijaService.DeleteAplikacijaByID(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Application deleted successfully",
		})
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

	err = h.aplikacijaService.DeleteAplikacija(id, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Application deleted successfully",
	})
}