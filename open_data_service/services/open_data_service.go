package services

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"open_data_service/models"
	"sort"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// OpenDataService handles all open data operations
type OpenDataService struct {
	stDomsCollection               *mongo.Collection
	sobasCollection                *mongo.Collection
	aplikacijeCollection           *mongo.Collection
	prihvaceneAplikacijeCollection *mongo.Collection
}

// NewOpenDataService creates a new OpenDataService
func NewOpenDataService(
	stDomsCollection *mongo.Collection,
	sobasCollection *mongo.Collection,
	aplikacijeCollection *mongo.Collection,
	prihvaceneAplikacijeCollection *mongo.Collection,
) *OpenDataService {
	return &OpenDataService{
		stDomsCollection:               stDomsCollection,
		sobasCollection:                sobasCollection,
		aplikacijeCollection:           aplikacijeCollection,
		prihvaceneAplikacijeCollection: prihvaceneAplikacijeCollection,
	}
}

// ====================
// 1. Public Statistics Dashboard
// ====================

// GetPublicStatistics returns comprehensive public statistics about all dorms
func (s *OpenDataService) GetPublicStatistics() (*models.PublicStatistics, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stats := &models.PublicStatistics{
		LastUpdated: time.Now(),
	}

	// Get total dorms
	totalDorms, err := s.stDomsCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	stats.TotalDorms = int(totalDorms)

	// Get total rooms
	totalRooms, err := s.sobasCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	stats.TotalRooms = int(totalRooms)

	// Get all rooms to calculate capacity and amenities
	cursor, err := s.sobasCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rooms []models.Soba
	if err = cursor.All(ctx, &rooms); err != nil {
		return nil, err
	}

	// Calculate total capacity and amenities distribution
	amenitiesMap := make(map[string]int)
	roomTypeMap := make(map[int]int)
	totalCapacity := 0

	for _, room := range rooms {
		totalCapacity += room.Krevetnost
		roomTypeMap[room.Krevetnost]++

		for _, amenity := range room.Luksuzi {
			amenitiesMap[amenity]++
		}
	}

	stats.TotalCapacity = totalCapacity
	stats.AmenitiesDistribution = amenitiesMap
	stats.RoomTypeDistribution = roomTypeMap

	// Get total occupied spots from accepted applications
	totalOccupied, err := s.prihvaceneAplikacijeCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	stats.TotalOccupied = int(totalOccupied)

	// Calculate occupancy rate
	if totalCapacity > 0 {
		stats.OccupancyRate = float64(totalOccupied) / float64(totalCapacity) * 100
		stats.OccupancyRate = math.Round(stats.OccupancyRate*100) / 100
	}

	// Get application statistics
	appStats, err := s.getApplicationStatistics(ctx)
	if err != nil {
		return nil, err
	}
	stats.ApplicationStatistics = *appStats

	// Get per-dorm statistics
	dormStats, err := s.getDormStatistics(ctx)
	if err != nil {
		return nil, err
	}
	stats.DormStatistics = dormStats

	return stats, nil
}

// getApplicationStatistics calculates application-related statistics
func (s *OpenDataService) getApplicationStatistics(ctx context.Context) (*models.ApplicationStats, error) {
	stats := &models.ApplicationStats{}

	// Total applications
	total, err := s.aplikacijeCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	stats.TotalApplications = int(total)

	// Active applications
	active, err := s.aplikacijeCollection.CountDocuments(ctx, bson.M{"is_active": true})
	if err != nil {
		return nil, err
	}
	stats.ActiveApplications = int(active)

	// Accepted applications
	accepted, err := s.prihvaceneAplikacijeCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	stats.AcceptedApplications = int(accepted)

	// Acceptance rate
	if stats.TotalApplications > 0 {
		stats.AcceptanceRate = float64(stats.AcceptedApplications) / float64(stats.TotalApplications) * 100
		stats.AcceptanceRate = math.Round(stats.AcceptanceRate*100) / 100
	}

	// Average grade of accepted applications
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "avg_grade", Value: bson.D{{Key: "$avg", Value: "$prosek"}}},
		}}},
	}

	cursor, err := s.prihvaceneAplikacijeCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []struct {
		AvgGrade float64 `bson:"avg_grade"`
	}
	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	if len(result) > 0 {
		stats.AverageGradeOfAccepted = math.Round(result[0].AvgGrade*100) / 100
	}

	// Average grade of all applications
	cursor2, err := s.aplikacijeCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor2.Close(ctx)

	var result2 []struct {
		AvgGrade float64 `bson:"avg_grade"`
	}
	if err = cursor2.All(ctx, &result2); err != nil {
		return nil, err
	}
	if len(result2) > 0 {
		stats.AverageGradeOfApplications = math.Round(result2[0].AvgGrade*100) / 100
	}

	return stats, nil
}

// getDormStatistics calculates statistics for each dorm
func (s *OpenDataService) getDormStatistics(ctx context.Context) ([]models.DormStats, error) {
	// Get all dorms
	cursor, err := s.stDomsCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var dorms []models.StDom
	if err = cursor.All(ctx, &dorms); err != nil {
		return nil, err
	}

	var dormStats []models.DormStats

	for _, dorm := range dorms {
		stat := models.DormStats{
			DormID:   dorm.ID,
			DormName: dorm.Ime,
			Address:  dorm.Address,
			RoomTypes: make(map[int]int),
			Amenities: make(map[string]int),
		}

		// Get rooms for this dorm
		roomCursor, err := s.sobasCollection.Find(ctx, bson.M{"st_dom_id": dorm.ID})
		if err != nil {
			return nil, err
		}

		var rooms []models.Soba
		if err = roomCursor.All(ctx, &rooms); err != nil {
			roomCursor.Close(ctx)
			return nil, err
		}
		roomCursor.Close(ctx)

		stat.TotalRooms = len(rooms)

		// Calculate capacity, amenities, and room types
		for _, room := range rooms {
			stat.TotalCapacity += room.Krevetnost
			stat.RoomTypes[room.Krevetnost]++

			for _, amenity := range room.Luksuzi {
				stat.Amenities[amenity]++
			}
		}

		// Get occupied spots for this dorm
		pipeline := mongo.Pipeline{
			{{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "sobas"},
				{Key: "localField", Value: "soba_id"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "room"},
			}}},
			{{Key: "$unwind", Value: "$room"}},
			{{Key: "$match", Value: bson.D{
				{Key: "room.st_dom_id", Value: dorm.ID},
			}}},
			{{Key: "$count", Value: "total"}},
		}

		occupiedCursor, err := s.prihvaceneAplikacijeCollection.Aggregate(ctx, pipeline)
		if err != nil {
			return nil, err
		}

		var occupiedResult []struct {
			Total int `bson:"total"`
		}
		if err = occupiedCursor.All(ctx, &occupiedResult); err != nil {
			occupiedCursor.Close(ctx)
			return nil, err
		}
		occupiedCursor.Close(ctx)

		if len(occupiedResult) > 0 {
			stat.OccupiedSpots = occupiedResult[0].Total
		}

		stat.AvailableSpots = stat.TotalCapacity - stat.OccupiedSpots

		// Calculate occupancy rate
		if stat.TotalCapacity > 0 {
			stat.OccupancyRate = float64(stat.OccupiedSpots) / float64(stat.TotalCapacity) * 100
			stat.OccupancyRate = math.Round(stat.OccupancyRate*100) / 100
		}

		dormStats = append(dormStats, stat)
	}

	return dormStats, nil
}

// ====================
// 2. Room Availability Search
// ====================

// SearchAvailableRooms searches for available rooms with filters
func (s *OpenDataService) SearchAvailableRooms(filters models.RoomSearchFilters) ([]models.RoomAvailability, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Build query
	query := bson.M{}

	if filters.DormID != "" {
		dormOID, err := primitive.ObjectIDFromHex(filters.DormID)
		if err != nil {
			return nil, fmt.Errorf("invalid dorm_id format")
		}
		query["st_dom_id"] = dormOID
	}

	if filters.MinCapacity > 0 {
		query["krevetnost"] = bson.M{"$gte": filters.MinCapacity}
	}

	if filters.MaxCapacity > 0 {
		if query["krevetnost"] != nil {
			query["krevetnost"].(bson.M)["$lte"] = filters.MaxCapacity
		} else {
			query["krevetnost"] = bson.M{"$lte": filters.MaxCapacity}
		}
	}

	if len(filters.Amenities) > 0 {
		query["luksuzi"] = bson.M{"$all": filters.Amenities}
	}

	// Set default limit and offset
	if filters.Limit <= 0 {
		filters.Limit = 50
	}
	if filters.Offset < 0 {
		filters.Offset = 0
	}

	// Get rooms matching query
	opts := options.Find().SetLimit(int64(filters.Limit)).SetSkip(int64(filters.Offset))
	cursor, err := s.sobasCollection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rooms []models.Soba
	if err = cursor.All(ctx, &rooms); err != nil {
		return nil, err
	}

	// Get all dorms for reference
	dormsMap, err := s.getAllDormsMap(ctx)
	if err != nil {
		return nil, err
	}

	// Calculate availability for each room
	var roomAvailabilities []models.RoomAvailability

	for _, room := range rooms {
		// Count occupied spots in this room
		occupied, err := s.prihvaceneAplikacijeCollection.CountDocuments(ctx, bson.M{"soba_id": room.ID})
		if err != nil {
			return nil, err
		}

		available := room.Krevetnost - int(occupied)
		isAvailable := available > 0

		// Filter by availability if requested
		if filters.OnlyAvailable && !isAvailable {
			continue
		}

		dorm := dormsMap[room.StDomID.Hex()]

		roomAvailability := models.RoomAvailability{
			RoomID:         room.ID,
			DormID:         room.StDomID,
			DormName:       dorm.Ime,
			DormAddress:    dorm.Address,
			Capacity:       room.Krevetnost,
			Occupied:       int(occupied),
			AvailableSpots: available,
			Amenities:      room.Luksuzi,
			IsAvailable:    isAvailable,
		}

		roomAvailabilities = append(roomAvailabilities, roomAvailability)
	}

	return roomAvailabilities, nil
}

// getAllDormsMap returns a map of dorm ID to dorm object
func (s *OpenDataService) getAllDormsMap(ctx context.Context) (map[string]models.StDom, error) {
	cursor, err := s.stDomsCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var dorms []models.StDom
	if err = cursor.All(ctx, &dorms); err != nil {
		return nil, err
	}

	dormsMap := make(map[string]models.StDom)
	for _, dorm := range dorms {
		dormsMap[dorm.ID.Hex()] = dorm
	}

	return dormsMap, nil
}

// ====================
// 3. Dorm Comparison Tool
// ====================

// CompareDorms compares multiple dorms side-by-side
func (s *OpenDataService) CompareDorms(dormIDs []string) (*models.DormComparison, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var dormOIDs []primitive.ObjectID
	for _, idStr := range dormIDs {
		oid, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid dorm_id: %s", idStr)
		}
		dormOIDs = append(dormOIDs, oid)
	}

	// Get dorms
	cursor, err := s.stDomsCollection.Find(ctx, bson.M{"_id": bson.M{"$in": dormOIDs}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var dorms []models.StDom
	if err = cursor.All(ctx, &dorms); err != nil {
		return nil, err
	}

	var comparisonDetails []models.DormComparisonDetail

	for _, dorm := range dorms {
		detail := models.DormComparisonDetail{
			DormID:   dorm.ID,
			DormName: dorm.Ime,
			Address:  dorm.Address,
			ContactInfo: models.ContactInfo{
				Phone: dorm.TelephoneNumber,
				Email: dorm.Email,
			},
			RoomDistribution: make(map[int]int),
			AmenitiesOffered: make(map[string]int),
		}

		// Get rooms for this dorm
		roomCursor, err := s.sobasCollection.Find(ctx, bson.M{"st_dom_id": dorm.ID})
		if err != nil {
			return nil, err
		}

		var rooms []models.Soba
		if err = roomCursor.All(ctx, &rooms); err != nil {
			roomCursor.Close(ctx)
			return nil, err
		}
		roomCursor.Close(ctx)

		// Calculate capacity info
		totalCapacity := 0
		for _, room := range rooms {
			totalCapacity += room.Krevetnost
			detail.RoomDistribution[room.Krevetnost]++

			for _, amenity := range room.Luksuzi {
				detail.AmenitiesOffered[amenity]++
			}
		}

		detail.Capacity.TotalRooms = len(rooms)
		detail.Capacity.TotalCapacity = totalCapacity

		// Get occupied spots
		pipeline := mongo.Pipeline{
			{{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "sobas"},
				{Key: "localField", Value: "soba_id"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "room"},
			}}},
			{{Key: "$unwind", Value: "$room"}},
			{{Key: "$match", Value: bson.D{
				{Key: "room.st_dom_id", Value: dorm.ID},
			}}},
			{{Key: "$count", Value: "total"}},
		}

		occupiedCursor, err := s.prihvaceneAplikacijeCollection.Aggregate(ctx, pipeline)
		if err != nil {
			return nil, err
		}

		var occupiedResult []struct {
			Total int `bson:"total"`
		}
		if err = occupiedCursor.All(ctx, &occupiedResult); err != nil {
			occupiedCursor.Close(ctx)
			return nil, err
		}
		occupiedCursor.Close(ctx)

		if len(occupiedResult) > 0 {
			detail.Capacity.OccupiedSpots = occupiedResult[0].Total
		}

		detail.Capacity.AvailableSpots = detail.Capacity.TotalCapacity - detail.Capacity.OccupiedSpots

		if detail.Capacity.TotalCapacity > 0 {
			detail.Capacity.OccupancyRate = float64(detail.Capacity.OccupiedSpots) / float64(detail.Capacity.TotalCapacity) * 100
			detail.Capacity.OccupancyRate = math.Round(detail.Capacity.OccupancyRate*100) / 100
		}

		// Get application metrics for this dorm
		appMetrics, err := s.getDormApplicationMetrics(ctx, dorm.ID)
		if err != nil {
			return nil, err
		}
		detail.ApplicationMetrics = *appMetrics

		comparisonDetails = append(comparisonDetails, detail)
	}

	return &models.DormComparison{
		Dorms:          comparisonDetails,
		ComparisonDate: time.Now(),
	}, nil
}

// getDormApplicationMetrics calculates application metrics for a specific dorm
func (s *OpenDataService) getDormApplicationMetrics(ctx context.Context, dormID primitive.ObjectID) (*models.ApplicationMetrics, error) {
	metrics := &models.ApplicationMetrics{}

	// Get all applications for rooms in this dorm
	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "sobas"},
			{Key: "localField", Value: "soba_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "room"},
		}}},
		{{Key: "$unwind", Value: "$room"}},
		{{Key: "$match", Value: bson.D{
			{Key: "room.st_dom_id", Value: dormID},
		}}},
	}

	cursor, err := s.aplikacijeCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var applications []models.Aplikacija
	if err = cursor.All(ctx, &applications); err != nil {
		return nil, err
	}

	metrics.TotalApplications = len(applications)

	// Calculate average grade
	if len(applications) > 0 {
		totalGrade := 0
		for _, app := range applications {
			totalGrade += app.Prosek
		}
		metrics.AverageGrade = float64(totalGrade) / float64(len(applications))
		metrics.AverageGrade = math.Round(metrics.AverageGrade*100) / 100
	}

	// Get accepted applications
	acceptedPipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "sobas"},
			{Key: "localField", Value: "soba_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "room"},
		}}},
		{{Key: "$unwind", Value: "$room"}},
		{{Key: "$match", Value: bson.D{
			{Key: "room.st_dom_id", Value: dormID},
		}}},
		{{Key: "$count", Value: "total"}},
	}

	acceptedCursor, err := s.prihvaceneAplikacijeCollection.Aggregate(ctx, acceptedPipeline)
	if err != nil {
		return nil, err
	}
	defer acceptedCursor.Close(ctx)

	var acceptedResult []struct {
		Total int `bson:"total"`
	}
	if err = acceptedCursor.All(ctx, &acceptedResult); err != nil {
		return nil, err
	}

	if len(acceptedResult) > 0 {
		metrics.AcceptedApplications = acceptedResult[0].Total
	}

	// Calculate acceptance rate
	if metrics.TotalApplications > 0 {
		metrics.AcceptanceRate = float64(metrics.AcceptedApplications) / float64(metrics.TotalApplications) * 100
		metrics.AcceptanceRate = math.Round(metrics.AcceptanceRate*100) / 100
	}

	return metrics, nil
}

// ====================
// 4. Application Trends Analysis
// ====================

// GetApplicationTrends returns historical trends of applications by academic year
func (s *OpenDataService) GetApplicationTrends() (*models.ApplicationTrends, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	trends := &models.ApplicationTrends{
		GeneratedAt: time.Now(),
	}

	// Get yearly trends from accepted applications (grouped by academic year)
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$academic_year"},
			{Key: "total", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "avg_grade", Value: bson.D{{Key: "$avg", Value: "$prosek"}}},
			{Key: "min_grade", Value: bson.D{{Key: "$min", Value: "$prosek"}}},
			{Key: "max_grade", Value: bson.D{{Key: "$max", Value: "$prosek"}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}},
	}

	cursor, err := s.prihvaceneAplikacijeCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var yearlyResults []struct {
		AcademicYear string  `bson:"_id"`
		Total        int     `bson:"total"`
		AvgGrade     float64 `bson:"avg_grade"`
		MinGrade     int     `bson:"min_grade"`
		MaxGrade     int     `bson:"max_grade"`
	}

	if err = cursor.All(ctx, &yearlyResults); err != nil {
		return nil, err
	}

	for _, result := range yearlyResults {
		// Count total applications for this year (from aplikacije collection)
		// We need to match applications with accepted ones to find the academic year
		totalApps, err := s.countApplicationsByYear(ctx, result.AcademicYear)
		if err != nil {
			return nil, err
		}

		yearlyTrend := models.YearlyTrend{
			AcademicYear:         result.AcademicYear,
			TotalApplications:    totalApps,
			AcceptedApplications: result.Total,
			AverageGrade:         math.Round(result.AvgGrade*100) / 100,
			MinGrade:             result.MinGrade,
			MaxGrade:             result.MaxGrade,
		}

		if totalApps > 0 {
			yearlyTrend.AcceptanceRate = float64(result.Total) / float64(totalApps) * 100
			yearlyTrend.AcceptanceRate = math.Round(yearlyTrend.AcceptanceRate*100) / 100
		}

		trends.YearlyTrends = append(trends.YearlyTrends, yearlyTrend)
	}

	// Get per-dorm trends
	dormTrends, err := s.getDormApplicationTrends(ctx)
	if err != nil {
		return nil, err
	}
	trends.DormTrends = dormTrends

	// Calculate overall metrics
	trends.OverallMetrics = s.calculateTrendMetrics(trends.YearlyTrends)

	return trends, nil
}

// countApplicationsByYear counts total applications for a given academic year
func (s *OpenDataService) countApplicationsByYear(ctx context.Context, academicYear string) (int, error) {
	// Since applications don't have academic_year, we estimate based on creation date
	// This is a simplified approach - in production, you might want to add academic_year to applications
	
	// For now, count all applications as an approximation
	// You could improve this by adding academic_year field to Aplikacija model
	total, err := s.aplikacijeCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}

	return int(total), nil
}

// getDormApplicationTrends calculates application trends per dorm
func (s *OpenDataService) getDormApplicationTrends(ctx context.Context) ([]models.DormApplicationTrend, error) {
	// Get all dorms
	dormCursor, err := s.stDomsCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer dormCursor.Close(ctx)

	var dorms []models.StDom
	if err = dormCursor.All(ctx, &dorms); err != nil {
		return nil, err
	}

	var dormTrends []models.DormApplicationTrend

	for _, dorm := range dorms {
		// Get applications for this dorm
		pipeline := mongo.Pipeline{
			{{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "sobas"},
				{Key: "localField", Value: "soba_id"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "room"},
			}}},
			{{Key: "$unwind", Value: "$room"}},
			{{Key: "$match", Value: bson.D{
				{Key: "room.st_dom_id", Value: dorm.ID},
			}}},
			{{Key: "$count", Value: "total"}},
		}

		cursor, err := s.aplikacijeCollection.Aggregate(ctx, pipeline)
		if err != nil {
			return nil, err
		}

		var appResult []struct {
			Total int `bson:"total"`
		}
		if err = cursor.All(ctx, &appResult); err != nil {
			cursor.Close(ctx)
			return nil, err
		}
		cursor.Close(ctx)

		totalApps := 0
		if len(appResult) > 0 {
			totalApps = appResult[0].Total
		}

		// Get accepted applications
		acceptedCursor, err := s.prihvaceneAplikacijeCollection.Aggregate(ctx, pipeline)
		if err != nil {
			return nil, err
		}

		var acceptedResult []struct {
			Total int `bson:"total"`
		}
		if err = acceptedCursor.All(ctx, &acceptedResult); err != nil {
			acceptedCursor.Close(ctx)
			return nil, err
		}
		acceptedCursor.Close(ctx)

		acceptedApps := 0
		if len(acceptedResult) > 0 {
			acceptedApps = acceptedResult[0].Total
		}

		trend := models.DormApplicationTrend{
			DormID:               dorm.ID,
			DormName:             dorm.Ime,
			TotalApplications:    totalApps,
			AcceptedApplications: acceptedApps,
		}

		if totalApps > 0 {
			trend.AcceptanceRate = float64(acceptedApps) / float64(totalApps) * 100
			trend.AcceptanceRate = math.Round(trend.AcceptanceRate*100) / 100
		}

		dormTrends = append(dormTrends, trend)
	}

	return dormTrends, nil
}

// calculateTrendMetrics calculates overall trend metrics
func (s *OpenDataService) calculateTrendMetrics(yearlyTrends []models.YearlyTrend) models.TrendMetrics {
	metrics := models.TrendMetrics{
		TotalYears:      len(yearlyTrends),
		TrendDirection:  "stable",
	}

	if len(yearlyTrends) == 0 {
		return metrics
	}

	// Calculate average applications per year
	totalApps := 0
	for _, trend := range yearlyTrends {
		totalApps += trend.TotalApplications
	}
	metrics.AverageApplicationsPerYear = totalApps / len(yearlyTrends)

	// Determine trend direction (simple comparison of first and last year)
	if len(yearlyTrends) > 1 {
		first := float64(yearlyTrends[0].TotalApplications)
		last := float64(yearlyTrends[len(yearlyTrends)-1].TotalApplications)

		if last > first*1.1 { // 10% increase
			metrics.TrendDirection = "increasing"
		} else if last < first*0.9 { // 10% decrease
			metrics.TrendDirection = "decreasing"
		}
	}

	return metrics
}

// ====================
// 5. Real-time Occupancy Heatmap
// ====================

// GetOccupancyHeatmap returns real-time occupancy data for visualization
func (s *OpenDataService) GetOccupancyHeatmap() (*models.OccupancyHeatmap, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	heatmap := &models.OccupancyHeatmap{
		GeneratedAt: time.Now(),
	}

	// Get all dorms
	cursor, err := s.stDomsCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var dorms []models.StDom
	if err = cursor.All(ctx, &dorms); err != nil {
		return nil, err
	}

	var points []models.DormOccupancyPoint
	var occupancyRates []float64

	for _, dorm := range dorms {
		point := models.DormOccupancyPoint{
			DormID:   dorm.ID,
			DormName: dorm.Ime,
			Address:  dorm.Address,
		}

		// Get total capacity
		pipeline := mongo.Pipeline{
			{{Key: "$match", Value: bson.D{{Key: "st_dom_id", Value: dorm.ID}}}},
			{{Key: "$group", Value: bson.D{
				{Key: "_id", Value: nil},
				{Key: "total_capacity", Value: bson.D{{Key: "$sum", Value: "$krevetnost"}}},
			}}},
		}

		capacityCursor, err := s.sobasCollection.Aggregate(ctx, pipeline)
		if err != nil {
			return nil, err
		}

		var capacityResult []struct {
			TotalCapacity int `bson:"total_capacity"`
		}
		if err = capacityCursor.All(ctx, &capacityResult); err != nil {
			capacityCursor.Close(ctx)
			return nil, err
		}
		capacityCursor.Close(ctx)

		if len(capacityResult) > 0 {
			point.TotalCapacity = capacityResult[0].TotalCapacity
		}

		// Get occupied spots
		occupiedPipeline := mongo.Pipeline{
			{{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "sobas"},
				{Key: "localField", Value: "soba_id"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "room"},
			}}},
			{{Key: "$unwind", Value: "$room"}},
			{{Key: "$match", Value: bson.D{
				{Key: "room.st_dom_id", Value: dorm.ID},
			}}},
			{{Key: "$count", Value: "total"}},
		}

		occupiedCursor, err := s.prihvaceneAplikacijeCollection.Aggregate(ctx, occupiedPipeline)
		if err != nil {
			return nil, err
		}

		var occupiedResult []struct {
			Total int `bson:"total"`
		}
		if err = occupiedCursor.All(ctx, &occupiedResult); err != nil {
			occupiedCursor.Close(ctx)
			return nil, err
		}
		occupiedCursor.Close(ctx)

		if len(occupiedResult) > 0 {
			point.OccupiedSpots = occupiedResult[0].Total
		}

		point.AvailableSpots = point.TotalCapacity - point.OccupiedSpots

		// Calculate occupancy rate
		if point.TotalCapacity > 0 {
			point.OccupancyRate = float64(point.OccupiedSpots) / float64(point.TotalCapacity) * 100
			point.OccupancyRate = math.Round(point.OccupancyRate*100) / 100
			occupancyRates = append(occupancyRates, point.OccupancyRate)
		}

		// Determine status
		if point.OccupancyRate >= 80 {
			point.Status = "high"
		} else if point.OccupancyRate >= 50 {
			point.Status = "medium"
		} else {
			point.Status = "low"
		}

		points = append(points, point)
	}

	heatmap.Dorms = points

	// Calculate summary
	heatmap.Summary = s.calculateOccupancySummary(points, occupancyRates)

	return heatmap, nil
}

// calculateOccupancySummary calculates summary statistics for occupancy
func (s *OpenDataService) calculateOccupancySummary(points []models.DormOccupancyPoint, rates []float64) models.OccupancySummary {
	summary := models.OccupancySummary{}

	if len(rates) == 0 {
		return summary
	}

	// Calculate average
	total := 0.0
	for _, rate := range rates {
		total += rate
	}
	summary.AverageOccupancy = math.Round((total/float64(len(rates)))*100) / 100

	// Find highest and lowest
	sort.Float64s(rates)
	summary.LowestOccupancy = rates[0]
	summary.HighestOccupancy = rates[len(rates)-1]

	// Count full and empty dorms
	for _, point := range points {
		if point.OccupancyRate >= 100 {
			summary.FullDorms++
		}
		if point.OccupancyRate == 0 {
			summary.EmptyDorms++
		}
	}

	return summary
}

// ====================
// 6. Open Data Export (CSV/JSON)
// ====================

// ExportData exports data in CSV or JSON format
func (s *OpenDataService) ExportData(dataset string, format models.ExportFormat) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	switch dataset {
	case "dorms":
		return s.exportDorms(ctx, format)
	case "rooms":
		return s.exportRooms(ctx, format)
	case "statistics":
		return s.exportStatistics(ctx, format)
	default:
		return nil, fmt.Errorf("unknown dataset: %s", dataset)
	}
}

// exportDorms exports dorm data
func (s *OpenDataService) exportDorms(ctx context.Context, format models.ExportFormat) (interface{}, error) {
	cursor, err := s.stDomsCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var dorms []models.StDom
	if err = cursor.All(ctx, &dorms); err != nil {
		return nil, err
	}

	if format == models.ExportFormatJSON {
		return dorms, nil
	}

	// CSV format
	var csvData [][]string
	csvData = append(csvData, []string{"ID", "Name", "Address", "Phone", "Email", "Created At"})

	for _, dorm := range dorms {
		csvData = append(csvData, []string{
			dorm.ID.Hex(),
			dorm.Ime,
			dorm.Address,
			dorm.TelephoneNumber,
			dorm.Email,
			dorm.CreatedAt.Format(time.RFC3339),
		})
	}

	return csvData, nil
}

// exportRooms exports room data with availability
func (s *OpenDataService) exportRooms(ctx context.Context, format models.ExportFormat) (interface{}, error) {
	// Get rooms with availability data
	rooms, err := s.SearchAvailableRooms(models.RoomSearchFilters{Limit: 1000})
	if err != nil {
		return nil, err
	}

	if format == models.ExportFormatJSON {
		return rooms, nil
	}

	// CSV format
	var csvData [][]string
	csvData = append(csvData, []string{"Room ID", "Dorm ID", "Dorm Name", "Capacity", "Occupied", "Available", "Amenities", "Is Available"})

	for _, room := range rooms {
		csvData = append(csvData, []string{
			room.RoomID.Hex(),
			room.DormID.Hex(),
			room.DormName,
			fmt.Sprintf("%d", room.Capacity),
			fmt.Sprintf("%d", room.Occupied),
			fmt.Sprintf("%d", room.AvailableSpots),
			strings.Join(room.Amenities, ";"),
			fmt.Sprintf("%t", room.IsAvailable),
		})
	}

	return csvData, nil
}

// exportStatistics exports statistical data
func (s *OpenDataService) exportStatistics(ctx context.Context, format models.ExportFormat) (interface{}, error) {
	stats, err := s.GetPublicStatistics()
	if err != nil {
		return nil, err
	}

	if format == models.ExportFormatJSON {
		return stats, nil
	}

	// CSV format - export dorm statistics
	var csvData [][]string
	csvData = append(csvData, []string{"Dorm Name", "Address", "Total Rooms", "Total Capacity", "Occupied", "Available", "Occupancy Rate"})

	for _, dormStat := range stats.DormStatistics {
		csvData = append(csvData, []string{
			dormStat.DormName,
			dormStat.Address,
			fmt.Sprintf("%d", dormStat.TotalRooms),
			fmt.Sprintf("%d", dormStat.TotalCapacity),
			fmt.Sprintf("%d", dormStat.OccupiedSpots),
			fmt.Sprintf("%d", dormStat.AvailableSpots),
			fmt.Sprintf("%.2f", dormStat.OccupancyRate),
		})
	}

	return csvData, nil
}

// FormatCSV formats CSV data to string
func FormatCSV(data [][]string) (string, error) {
	var builder strings.Builder
	writer := csv.NewWriter(&builder)

	if err := writer.WriteAll(data); err != nil {
		return "", err
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}

	return builder.String(), nil
}

// FormatJSON formats data to JSON string
func FormatJSON(data interface{}) (string, error) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
