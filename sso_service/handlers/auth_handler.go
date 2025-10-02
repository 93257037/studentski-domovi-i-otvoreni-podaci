package handlers

import (
	"net/http"
	"sso_service/models"
	"sso_service/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuthHandler - rukuje zahtevima za autentifikaciju
type AuthHandler struct {
	userService *services.UserService
}

// kreira novi AuthHandler sa prosledjenim user servisom
func NewAuthHandler(userService *services.UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

// registracija korisnika - prima podatke, validira ih i kreira novi nalog
// vraca gresku ako korisnik vec postoji ili su podaci neispravni
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.RegisterUser(req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	user.Password = ""

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}

// prijava korisnika - proverava email i lozinku, vraca JWT token ako su ispravni
// vraca gresku ako su podaci neispravni
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.userService.LoginUser(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   response.Token,
		"user":    response.User,
	})
}

// dobija profil trenutno ulogovanog korisnika na osnovu JWT tokena
// vraca podatke o korisniku bez lozinke
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	user, err := h.userService.GetUserByID(userID.(primitive.ObjectID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// brise nalog korisnika - prvo proverava da li korisnik ima aktivnu sobu
// ne dozvoljava brisanje ako korisnik ima dodeljenu sobu u domu
func (h *AuthHandler) DeleteAccount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	err := h.userService.DeleteUser(userID.(primitive.ObjectID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account deleted successfully",
	})
}

// provera zdravlja servisa - vraca status da li je servis aktivan
func (h *AuthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "sso_service",
	})
}

