package handlers

import (
	"fmt"
	"net/http"
	"open_data_service/models"
	"open_data_service/services"
	"strings"

	"github.com/gin-gonic/gin"
)

// OpenDataHandler handles all open data API requests
type OpenDataHandler struct {
	openDataService *services.OpenDataService
}

// NewOpenDataHandler creates a new OpenDataHandler
func NewOpenDataHandler(openDataService *services.OpenDataService) *OpenDataHandler {
	return &OpenDataHandler{
		openDataService: openDataService,
	}
}

// ====================
// 1. Public Statistics Dashboard
// ====================

// GetPublicStatistics returns comprehensive public statistics about all dorms
// GET /api/v1/open-data/statistics
func (h *OpenDataHandler) GetPublicStatistics(c *gin.Context) {
	stats, err := h.openDataService.GetPublicStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statistics": stats,
	})
}

// ====================
// 2. Room Availability Search
// ====================

// SearchAvailableRooms searches for available rooms with filters
// GET /api/v1/open-data/rooms/search
// Query params: dorm_id, min_capacity, max_capacity, amenities, only_available, limit, offset
func (h *OpenDataHandler) SearchAvailableRooms(c *gin.Context) {
	var filters models.RoomSearchFilters

	// Parse query parameters
	filters.DormID = c.Query("dorm_id")
	
	if minCap := c.Query("min_capacity"); minCap != "" {
		var minCapInt int
		if _, err := fmt.Sscanf(minCap, "%d", &minCapInt); err == nil {
			filters.MinCapacity = minCapInt
		}
	}

	if maxCap := c.Query("max_capacity"); maxCap != "" {
		var maxCapInt int
		if _, err := fmt.Sscanf(maxCap, "%d", &maxCapInt); err == nil {
			filters.MaxCapacity = maxCapInt
		}
	}

	// Parse amenities (comma-separated)
	if amenitiesStr := c.Query("amenities"); amenitiesStr != "" {
		filters.Amenities = strings.Split(amenitiesStr, ",")
	}

	filters.OnlyAvailable = c.Query("only_available") == "true"

	if limit := c.Query("limit"); limit != "" {
		var limitInt int
		if _, err := fmt.Sscanf(limit, "%d", &limitInt); err == nil {
			filters.Limit = limitInt
		}
	}

	if offset := c.Query("offset"); offset != "" {
		var offsetInt int
		if _, err := fmt.Sscanf(offset, "%d", &offsetInt); err == nil {
			filters.Offset = offsetInt
		}
	}

	rooms, err := h.openDataService.SearchAvailableRooms(filters)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rooms": rooms,
		"count": len(rooms),
	})
}

// ====================
// 3. Dorm Comparison Tool
// ====================

// CompareDorms compares multiple dorms side-by-side
// GET /api/v1/open-data/dorms/compare
// Query params: dorm_ids (comma-separated list of dorm IDs)
func (h *OpenDataHandler) CompareDorms(c *gin.Context) {
	dormIDsStr := c.Query("dorm_ids")
	if dormIDsStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dorm_ids parameter is required"})
		return
	}

	dormIDs := strings.Split(dormIDsStr, ",")
	if len(dormIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "at least one dorm_id must be provided"})
		return
	}

	if len(dormIDs) > 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "maximum 10 dorms can be compared at once"})
		return
	}

	comparison, err := h.openDataService.CompareDorms(dormIDs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comparison": comparison,
	})
}

// ====================
// 4. Application Trends Analysis
// ====================

// GetApplicationTrends returns historical trends of applications
// GET /api/v1/open-data/trends/applications
func (h *OpenDataHandler) GetApplicationTrends(c *gin.Context) {
	trends, err := h.openDataService.GetApplicationTrends()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"trends": trends,
	})
}

// ====================
// 5. Real-time Occupancy Heatmap
// ====================

// GetOccupancyHeatmap returns real-time occupancy data for visualization
// GET /api/v1/open-data/occupancy/heatmap
func (h *OpenDataHandler) GetOccupancyHeatmap(c *gin.Context) {
	heatmap, err := h.openDataService.GetOccupancyHeatmap()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"heatmap": heatmap,
	})
}

// ====================
// 6. Open Data Export (CSV/JSON)
// ====================

// ExportData exports data in CSV or JSON format
// GET /api/v1/open-data/export
// Query params: dataset (dorms, rooms, dorm-statistics, application-list, accepted-applications, dorm-trends, amenities-report, occupancy-report, room-types), format (csv, json)
func (h *OpenDataHandler) ExportData(c *gin.Context) {
	dataset := c.Query("dataset")
	if dataset == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dataset parameter is required (dorms, rooms, dorm-statistics, application-list, accepted-applications, dorm-trends, amenities-report, occupancy-report, or room-types)"})
		return
	}

	formatStr := c.Query("format")
	if formatStr == "" {
		formatStr = "json"
	}

	var format models.ExportFormat
	if formatStr == "csv" {
		format = models.ExportFormatCSV
	} else if formatStr == "json" {
		format = models.ExportFormatJSON
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format must be 'csv' or 'json'"})
		return
	}

	data, err := h.openDataService.ExportData(dataset, format)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if format == models.ExportFormatCSV {
		csvData, ok := data.([][]string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to convert data to CSV format"})
			return
		}

		csvString, err := services.FormatCSV(csvData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", "attachment; filename="+dataset+".csv")
		c.String(http.StatusOK, csvString)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	}
}

// ====================
// Additional Helper Endpoints
// ====================

// GetRoomApplications proxies the request to st_dom_service to get accepted applications for a room
// GET /api/v1/open-data/rooms/:roomId/applications
func (h *OpenDataHandler) GetRoomApplications(c *gin.Context) {
	roomID := c.Param("roomId")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room ID is required"})
		return
	}

	// Get the authorization token from the request header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization token required"})
		return
	}

	// Call st_dom_service endpoint
	applications, err := h.openDataService.GetRoomApplicationsFromStDomService(roomID, authHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": applications,
	})
}

// GetApplicationsByAcademicYear proxies the request to st_dom_service to get accepted applications by academic year
// GET /api/v1/open-data/applications/academic-year
// This is now a public endpoint (no authentication required)
func (h *OpenDataHandler) GetApplicationsByAcademicYear(c *gin.Context) {
	academicYear := c.Query("academic_year")
	if academicYear == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "academic_year query parameter is required"})
		return
	}

	// Call st_dom_service endpoint (no auth required for public data)
	applications, err := h.openDataService.GetApplicationsByAcademicYearFromStDomService(academicYear)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": applications,
	})
}

// GetDormList returns a simple list of all dorms (for dropdown/selection)
// GET /api/v1/open-data/dorms/list
func (h *OpenDataHandler) GetDormList(c *gin.Context) {
	stats, err := h.openDataService.GetPublicStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type DormListItem struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Address string `json:"address"`
	}

	var dormList []DormListItem
	for _, dorm := range stats.DormStatistics {
		dormList = append(dormList, DormListItem{
			ID:      dorm.DormID.Hex(),
			Name:    dorm.DormName,
			Address: dorm.Address,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"dorms": dormList,
		"count": len(dormList),
	})
}

// GetAvailableAmenities returns a list of all unique amenities available
// GET /api/v1/open-data/amenities
func (h *OpenDataHandler) GetAvailableAmenities(c *gin.Context) {
	stats, err := h.openDataService.GetPublicStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	amenities := make([]string, 0, len(stats.AmenitiesDistribution))
	for amenity := range stats.AmenitiesDistribution {
		amenities = append(amenities, amenity)
	}

	c.JSON(http.StatusOK, gin.H{
		"amenities": amenities,
		"count":     len(amenities),
	})
}
