package handlers

import (
	"net/http"
	"st_dom_service/models"
	"st_dom_service/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// StDomHandler - rukuje zahtevima vezanim za studentske domove
type StDomHandler struct {
	stDomService *services.StDomService
	sobaService  *services.SobaService
}

// kreira novi StDomHandler sa potrebnim servisima
func NewStDomHandler(stDomService *services.StDomService, sobaService *services.SobaService) *StDomHandler {
	return &StDomHandler{
		stDomService: stDomService,
		sobaService:  sobaService,
	}
}

// kreira novi studentski dom - prima podatke, validira ih i cuva u bazu
// proverava da li vec postoji dom sa istom adresom
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

// dobija studentski dom po ID-u iz baze podataka
// vraca gresku ako ID nije valjan ili dom ne postoji
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

// dobija sve studentske domove iz baze podataka
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

// azurira podatke o studentskom domu
// prima ID i nove podatke, validira ih i cuva promene
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

// brise studentski dom i sve povezane sobe
// prvo brise sve sobe koje pripadaju domu, zatim brise sam dom
func (h *StDomHandler) DeleteStDom(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	err = h.sobaService.DeleteSobasByStDomID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete associated rooms"})
		return
	}

	err = h.stDomService.DeleteStDom(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Student dormitory and associated rooms deleted successfully",
	})
}

// dobija sve dostupne sobe za odredjeni studentski dom
// vraca samo sobe koje nisu potpuno popunjene
func (h *StDomHandler) GetStDomRooms(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	_, err = h.stDomService.GetStDomByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student dormitory not found"})
		return
	}

	sobas, err := h.sobaService.GetAvailableSobasByStDomID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sobas": sobas,
	})
}
