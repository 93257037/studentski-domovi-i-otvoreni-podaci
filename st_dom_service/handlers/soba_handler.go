package handlers

import (
	"net/http"
	"st_dom_service/models"
	"st_dom_service/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SobaHandler - rukuje zahtevima vezanim za sobe
type SobaHandler struct {
	sobaService  *services.SobaService
	stDomService *services.StDomService
}

// kreira novi SobaHandler sa potrebnim servisima
func NewSobaHandler(sobaService *services.SobaService, stDomService *services.StDomService) *SobaHandler {
	return &SobaHandler{
		sobaService:  sobaService,
		stDomService: stDomService,
	}
}

// kreira novu sobu u studentskom domu
// proverava da li dom postoji pre kreiranja sobe
func (h *SobaHandler) CreateSoba(c *gin.Context) {
	var req models.CreateSobaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

// dobija sobu po ID-u iz baze podataka
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

// dobija sve sobe iz baze podataka
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

// azurira podatke o sobi
// prima ID sobe i nove podatke za azuriranje
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

// brise sobu iz baze podataka
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
