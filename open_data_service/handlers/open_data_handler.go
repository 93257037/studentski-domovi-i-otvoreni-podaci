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
	service        *services.OpenDataService
	httpClient     *services.HTTPClientService
}

// NewOpenDataHandler creates a new OpenDataHandler
func NewOpenDataHandler(service *services.OpenDataService, httpClient *services.HTTPClientService) *OpenDataHandler {
	return &OpenDataHandler{
		service:    service,
		httpClient: httpClient,
	}
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

// SearchStDomsByIme godoc
// @Summary Search student dormitories by name
// @Description Search student dormitories using regex pattern matching on ime (name) field
// @Tags Open Data
// @Accept json
// @Produce json
// @Param ime query string true "Name search pattern (supports regex)"
// @Success 200 {object} map[string]interface{} "List of student dormitories"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/st-doms/search-by-ime [get]
func (h *OpenDataHandler) SearchStDomsByIme(c *gin.Context) {
	imePattern := c.Query("ime")
	if imePattern == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ime parameter is required",
		})
		return
	}

	stDoms, err := h.service.SearchStDomsByIme(imePattern)
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
// @Description Filter rooms using multiple criteria: luxury amenities, student dormitory, address, and bed capacity
// @Tags Open Data
// @Accept json
// @Produce json
// @Param luksuzi query string false "Comma-separated list of luxury amenities" example="klima,terasa"
// @Param st_dom_id query string false "Student dormitory ID"
// @Param address query string false "Address search pattern (regex, case-insensitive)"
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

	// Parse address pattern
	addressPattern := c.Query("address")

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

	rooms, err := h.service.AdvancedFilterRooms(luksuzi, stDomID, addressPattern, exact, min, max)
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

// GetTopFullStDoms godoc
// @Summary Get top 3 most full student dormitories
// @Description Returns the top 3 student dormitories with the most residents
// @Tags Statistics
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "List of top 3 most full dormitories"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/statistics/top-full-st-doms [get]
func (h *OpenDataHandler) GetTopFullStDoms(c *gin.Context) {
	stats, err := h.service.GetTopFullStDoms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get top full student dormitories",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  stats,
		"count": len(stats),
	})
}

// GetTopEmptyStDoms godoc
// @Summary Get top 3 most empty student dormitories
// @Description Returns the top 3 student dormitories with the fewest residents
// @Tags Statistics
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "List of top 3 most empty dormitories"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/statistics/top-empty-st-doms [get]
func (h *OpenDataHandler) GetTopEmptyStDoms(c *gin.Context) {
	stats, err := h.service.GetTopEmptyStDoms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get top empty student dormitories",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  stats,
		"count": len(stats),
	})
}

// GetStDomWithMostApplications godoc
// @Summary Get student dormitory with most applications
// @Description Returns the student dormitory with the highest number of active applications
// @Tags Statistics
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Student dormitory with most applications"
// @Failure 404 {object} map[string]interface{} "No applications found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/statistics/st-dom-most-applications [get]
func (h *OpenDataHandler) GetStDomWithMostApplications(c *gin.Context) {
	stats, err := h.service.GetStDomWithMostApplications()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get student dormitory with most applications",
		})
		return
	}

	if stats == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No applications found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": stats,
	})
}

// GetStDomWithHighestAverageProsek godoc
// @Summary Get student dormitory with highest average prosek
// @Description Returns the student dormitory with the highest average grade (prosek) of its residents
// @Tags Statistics
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Student dormitory with highest average prosek"
// @Failure 404 {object} map[string]interface{} "No residents found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/statistics/st-dom-highest-average-prosek [get]
func (h *OpenDataHandler) GetStDomWithHighestAverageProsek(c *gin.Context) {
	stats, err := h.service.GetStDomWithHighestAverageProsek()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get student dormitory with highest average prosek",
		})
		return
	}

	if stats == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No residents found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": stats,
	})
}

// GetPrihvaceneAplikacije godoc
// @Summary Get all accepted applications from st_dom_service
// @Description Retrieves all accepted applications by calling st_dom_service inter-service communication
// @Tags Inter-Service Communication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "List of accepted applications"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/inter-service/prihvacene-aplikacije [get]
func (h *OpenDataHandler) GetPrihvaceneAplikacije(c *gin.Context) {
	// Extract Authorization header from the incoming request
	authHeader := c.GetHeader("Authorization")
	
	response, err := h.httpClient.GetPrihvaceneAplikacije(authHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve accepted applications from st_dom_service: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data retrieved from st_dom_service via inter-service communication",
		"data":    response.PrihvaceneAplikacije,
		"count":   response.Count,
		"source":  "st_dom_service",
	})
}

// GetPrihvaceneAplikacijeForUser godoc
// @Summary Get accepted applications for a specific user from st_dom_service
// @Description Retrieves accepted applications for a specific user by calling st_dom_service
// @Tags Inter-Service Communication
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "List of accepted applications for user"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/inter-service/prihvacene-aplikacije/user/{userId} [get]
func (h *OpenDataHandler) GetPrihvaceneAplikacijeForUser(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
		})
		return
	}

	// Extract Authorization header from the incoming request
	authHeader := c.GetHeader("Authorization")

	response, err := h.httpClient.GetPrihvaceneAplikacijeForUser(userID, authHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve accepted applications for user from st_dom_service: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data retrieved from st_dom_service via inter-service communication",
		"data":    response.PrihvaceneAplikacije,
		"count":   response.Count,
		"user_id": userID,
		"source":  "st_dom_service",
	})
}

// GetPrihvaceneAplikacijeForRoom godoc
// @Summary Get accepted applications for a specific room from st_dom_service
// @Description Retrieves accepted applications for a specific room by calling st_dom_service
// @Tags Inter-Service Communication
// @Accept json
// @Produce json
// @Param roomId path string true "Room ID"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "List of accepted applications for room"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/inter-service/prihvacene-aplikacije/room/{roomId} [get]
func (h *OpenDataHandler) GetPrihvaceneAplikacijeForRoom(c *gin.Context) {
	roomID := c.Param("roomId")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Room ID is required",
		})
		return
	}

	// Extract Authorization header from the incoming request
	authHeader := c.GetHeader("Authorization")

	response, err := h.httpClient.GetPrihvaceneAplikacijeForRoom(roomID, authHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve accepted applications for room from st_dom_service: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data retrieved from st_dom_service via inter-service communication",
		"data":    response.PrihvaceneAplikacije,
		"count":   response.Count,
		"room_id": roomID,
		"source":  "st_dom_service",
	})
}

// GetPrihvaceneAplikacijeForAcademicYear godoc
// @Summary Get accepted applications for a specific academic year from st_dom_service
// @Description Retrieves accepted applications for a specific academic year by calling st_dom_service
// @Tags Inter-Service Communication
// @Accept json
// @Produce json
// @Param academicYear path string true "Academic Year (e.g., 2024/2025)"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "List of accepted applications for academic year"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/inter-service/prihvacene-aplikacije/academic-year/{academicYear} [get]
func (h *OpenDataHandler) GetPrihvaceneAplikacijeForAcademicYear(c *gin.Context) {
	academicYear := c.Param("academicYear")
	if academicYear == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Academic year is required",
		})
		return
	}

	// Extract Authorization header from the incoming request
	authHeader := c.GetHeader("Authorization")

	response, err := h.httpClient.GetPrihvaceneAplikacijeForAcademicYear(academicYear, authHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve accepted applications for academic year from st_dom_service: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Data retrieved from st_dom_service via inter-service communication",
		"data":          response.PrihvaceneAplikacije,
		"count":         response.Count,
		"academic_year": academicYear,
		"source":        "st_dom_service",
	})
}

// CheckStDomServiceHealth godoc
// @Summary Check health of st_dom_service
// @Description Checks if st_dom_service is available and responding
// @Tags Inter-Service Communication
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Service health status"
// @Failure 500 {object} map[string]interface{} "Service unavailable"
// @Router /api/v1/inter-service/health [get]
func (h *OpenDataHandler) CheckStDomServiceHealth(c *gin.Context) {
	err := h.httpClient.HealthCheck()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "unhealthy",
			"service": "st_dom_service",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "st_dom_service",
		"message": "st_dom_service is available and responding",
	})
}

