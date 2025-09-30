package handlers

import (
	"net/http"
	"open_data_service/models"
	"open_data_service/services"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OpenDataHandler handles HTTP requests for open data
type OpenDataHandler struct {
	service *services.OpenDataService
}

// NewOpenDataHandler creates a new OpenDataHandler
func NewOpenDataHandler(service *services.OpenDataService) *OpenDataHandler {
	return &OpenDataHandler{service: service}
}

// FilterRoomsByLuksuz godoc
// @Summary Filter rooms by luxury amenities
// @Description Filter rooms by any combination of luxury amenities (klima, terasa, sopstveno kupatilo, etc.)
// @Tags Open Data
// @Accept json
// @Produce json
// @Param luksuzi query string false "Comma-separated list of luxury amenities" example="klima,terasa"
// @Success 200 {object} map[string]interface{} "List of rooms"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/rooms/filter-by-luksuz [get]
func (h *OpenDataHandler) FilterRoomsByLuksuz(c *gin.Context) {
	// Parse luksuzi from query parameter
	luksuziStr := c.Query("luksuzi")
	var luksuzi []models.Luksuzi

	if luksuziStr != "" {
		luksuziList := strings.Split(luksuziStr, ",")
		for _, l := range luksuziList {
			l = strings.TrimSpace(l)
			luksuz := models.Luksuzi(l)
			if !luksuz.IsValid() {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid luxury amenity: " + l,
				})
				return
			}
			luksuzi = append(luksuzi, luksuz)
		}
	}

	rooms, err := h.service.FilterRoomsByLuksuz(luksuzi)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to filter rooms",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  rooms,
		"count": len(rooms),
	})
}

// FilterRoomsByLuksuzAndStDom godoc
// @Summary Filter rooms by luxury amenities and student dormitory
// @Description Filter rooms by luxury amenities and specific student dormitory
// @Tags Open Data
// @Accept json
// @Produce json
// @Param luksuzi query string false "Comma-separated list of luxury amenities" example="klima,terasa"
// @Param st_dom_id query string true "Student dormitory ID"
// @Success 200 {object} map[string]interface{} "List of rooms with dormitory info"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/rooms/filter-by-luksuz-and-stdom [get]
func (h *OpenDataHandler) FilterRoomsByLuksuzAndStDom(c *gin.Context) {
	// Parse st_dom_id
	stDomIDStr := c.Query("st_dom_id")
	if stDomIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "st_dom_id is required",
		})
		return
	}

	stDomID, err := primitive.ObjectIDFromHex(stDomIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid st_dom_id format",
		})
		return
	}

	// Parse luksuzi from query parameter
	luksuziStr := c.Query("luksuzi")
	var luksuzi []models.Luksuzi

	if luksuziStr != "" {
		luksuziList := strings.Split(luksuziStr, ",")
		for _, l := range luksuziList {
			l = strings.TrimSpace(l)
			luksuz := models.Luksuzi(l)
			if !luksuz.IsValid() {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid luxury amenity: " + l,
				})
				return
			}
			luksuzi = append(luksuzi, luksuz)
		}
	}

	rooms, err := h.service.FilterRoomsByLuksuzAndStDom(luksuzi, stDomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to filter rooms",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  rooms,
		"count": len(rooms),
	})
}

// FilterRoomsByKrevetnost godoc
// @Summary Filter rooms by bed capacity
// @Description Filter rooms by bed capacity (exact, min, max, or range)
// @Tags Open Data
// @Accept json
// @Produce json
// @Param exact query int false "Exact bed capacity"
// @Param min query int false "Minimum bed capacity"
// @Param max query int false "Maximum bed capacity"
// @Success 200 {object} map[string]interface{} "List of rooms"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/rooms/filter-by-krevetnost [get]
func (h *OpenDataHandler) FilterRoomsByKrevetnost(c *gin.Context) {
	var exact, min, max *int

	// Parse exact
	if exactStr := c.Query("exact"); exactStr != "" {
		val, err := strconv.Atoi(exactStr)
		if err != nil || val < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid exact value, must be a positive integer",
			})
			return
		}
		exact = &val
	}

	// Parse min
	if minStr := c.Query("min"); minStr != "" {
		val, err := strconv.Atoi(minStr)
		if err != nil || val < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid min value, must be a positive integer",
			})
			return
		}
		min = &val
	}

	// Parse max
	if maxStr := c.Query("max"); maxStr != "" {
		val, err := strconv.Atoi(maxStr)
		if err != nil || val < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid max value, must be a positive integer",
			})
			return
		}
		max = &val
	}

	// Validate that min <= max if both provided
	if min != nil && max != nil && *min > *max {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "min value cannot be greater than max value",
		})
		return
	}

	rooms, err := h.service.FilterRoomsByKrevetnost(exact, min, max)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to filter rooms",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  rooms,
		"count": len(rooms),
	})
}

// SearchStDomsByAddress godoc
// @Summary Search student dormitories by address
// @Description Search student dormitories using regex pattern matching on address field
// @Tags Open Data
// @Accept json
// @Produce json
// @Param address query string true "Address search pattern (supports regex)"
// @Success 200 {object} map[string]interface{} "List of student dormitories"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/st-doms/search-by-address [get]
func (h *OpenDataHandler) SearchStDomsByAddress(c *gin.Context) {
	addressPattern := c.Query("address")
	if addressPattern == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "address parameter is required",
		})
		return
	}

	stDoms, err := h.service.SearchStDomsByAddress(addressPattern)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to search student dormitories",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  stDoms,
		"count": len(stDoms),
	})
}

// AdvancedFilterRooms godoc
// @Summary Advanced room filtering
// @Description Filter rooms using multiple criteria: luxury amenities, student dormitory, and bed capacity
// @Tags Open Data
// @Accept json
// @Produce json
// @Param luksuzi query string false "Comma-separated list of luxury amenities" example="klima,terasa"
// @Param st_dom_id query string false "Student dormitory ID"
// @Param exact query int false "Exact bed capacity"
// @Param min query int false "Minimum bed capacity"
// @Param max query int false "Maximum bed capacity"
// @Success 200 {object} map[string]interface{} "List of rooms with dormitory info"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/rooms/advanced-filter [get]
func (h *OpenDataHandler) AdvancedFilterRooms(c *gin.Context) {
	// Parse luksuzi
	luksuziStr := c.Query("luksuzi")
	var luksuzi []models.Luksuzi

	if luksuziStr != "" {
		luksuziList := strings.Split(luksuziStr, ",")
		for _, l := range luksuziList {
			l = strings.TrimSpace(l)
			luksuz := models.Luksuzi(l)
			if !luksuz.IsValid() {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid luxury amenity: " + l,
				})
				return
			}
			luksuzi = append(luksuzi, luksuz)
		}
	}

	// Parse st_dom_id
	var stDomID *primitive.ObjectID
	if stDomIDStr := c.Query("st_dom_id"); stDomIDStr != "" {
		id, err := primitive.ObjectIDFromHex(stDomIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid st_dom_id format",
			})
			return
		}
		stDomID = &id
	}

	// Parse krevetnost parameters
	var exact, min, max *int

	if exactStr := c.Query("exact"); exactStr != "" {
		val, err := strconv.Atoi(exactStr)
		if err != nil || val < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid exact value, must be a positive integer",
			})
			return
		}
		exact = &val
	}

	if minStr := c.Query("min"); minStr != "" {
		val, err := strconv.Atoi(minStr)
		if err != nil || val < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid min value, must be a positive integer",
			})
			return
		}
		min = &val
	}

	if maxStr := c.Query("max"); maxStr != "" {
		val, err := strconv.Atoi(maxStr)
		if err != nil || val < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid max value, must be a positive integer",
			})
			return
		}
		max = &val
	}

	// Validate min <= max if both provided
	if min != nil && max != nil && *min > *max {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "min value cannot be greater than max value",
		})
		return
	}

	rooms, err := h.service.AdvancedFilterRooms(luksuzi, stDomID, exact, min, max)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to filter rooms",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  rooms,
		"count": len(rooms),
	})
}

// GetAllRooms godoc
// @Summary Get all rooms
// @Description Get all available rooms in the system
// @Tags Open Data
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "List of all rooms"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/rooms [get]
func (h *OpenDataHandler) GetAllRooms(c *gin.Context) {
	rooms, err := h.service.GetAllRooms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve rooms",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  rooms,
		"count": len(rooms),
	})
}

// GetAllStDoms godoc
// @Summary Get all student dormitories
// @Description Get all student dormitories in the system
// @Tags Open Data
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "List of all student dormitories"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/st-doms [get]
func (h *OpenDataHandler) GetAllStDoms(c *gin.Context) {
	stDoms, err := h.service.GetAllStDoms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve student dormitories",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  stDoms,
		"count": len(stDoms),
	})
}

// Health godoc
// @Summary Health check
// @Description Check if the service is running
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Service status"
// @Router /health [get]
func (h *OpenDataHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"service": "open_data_service",
	})
}

