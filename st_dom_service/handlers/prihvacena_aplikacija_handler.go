package handlers

import (
	"fmt"
	"net/http"
	"st_dom_service/models"
	"st_dom_service/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PrihvacenaAplikacijaHandler - rukuje zahtevima vezanim za prihvacene aplikacije
type PrihvacenaAplikacijaHandler struct {
	prihvacenaAplikacijaService *services.PrihvacenaAplikacijaService
}

// kreira novi PrihvacenaAplikacijaHandler sa potrebnim servisom
func NewPrihvacenaAplikacijaHandler(prihvacenaAplikacijaService *services.PrihvacenaAplikacijaService) *PrihvacenaAplikacijaHandler {
	return &PrihvacenaAplikacijaHandler{
		prihvacenaAplikacijaService: prihvacenaAplikacijaService,
	}
}

// odobrava aplikaciju za sobu - samo administratori mogu odobriti aplikacije
// kreira prihvacenu aplikaciju i generiše racun za placanje
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

// dobija prihvacenu aplikaciju po ID-u
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

// dobija sve prihvacene aplikacije
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

// dobija sve prihvacene aplikacije za odredjenog korisnika
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

// dobija sve prihvacene aplikacije za odredjenu sobu
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

// dobija sve prihvacene aplikacije za odredjenu skolsku godinu
func (h *PrihvacenaAplikacijaHandler) GetPrihvaceneAplikacijeForAcademicYear(c *gin.Context) {
	academicYear := c.Query("academic_year")
	fmt.Printf("DEBUG: st_dom_service received academicYear query: '%s'\n", academicYear)
	fmt.Printf("DEBUG: Full request path: %s\n", c.Request.URL.Path)
	fmt.Printf("DEBUG: Full request URL: %s\n", c.Request.URL.String())
	
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

// dobija najbolje studente rangirane po proseku
// prima limit kao parametar (default 10)
func (h *PrihvacenaAplikacijaHandler) GetTopStudentsByProsek(c *gin.Context) {
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

// dobija najbolje studente za odredjenu skolsku godinu rangirane po proseku
func (h *PrihvacenaAplikacijaHandler) GetTopStudentsByProsekForAcademicYear(c *gin.Context) {
	academicYear := c.Param("academicYear")
	if academicYear == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Academic year is required"})
		return
	}

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

// dobija najbolje studente za odredjenu sobu rangirane po proseku
func (h *PrihvacenaAplikacijaHandler) GetTopStudentsByProsekForRoom(c *gin.Context) {
	sobaIDParam := c.Param("sobaId")
	sobaID, err := primitive.ObjectIDFromHex(sobaIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID format"})
		return
	}

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

// dobija prihvacene aplikacije trenutno ulogovanog korisnika
func (h *PrihvacenaAplikacijaHandler) GetMyPrihvaceneAplikacije(c *gin.Context) {
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

// brise prihvacenu aplikaciju - samo administratori
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

// izbacuje studenta iz sobe - samo administratori mogu izbaciti studenta
// prima razlog izbacivanja i uklanja studenta iz sobe
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

// student dobrovoljno napusta sobu - korisnik moze sam da se odjavi iz sobe
func (h *PrihvacenaAplikacijaHandler) CheckoutFromRoom(c *gin.Context) {
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

// proverava da li korisnik ima aktivnu sobu - za komunikaciju izmedju servisa
// koristi se od strane SSO servisa da proveri da li korisnik moze da obrise nalog
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