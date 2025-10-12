package services

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"open_data_service/models"
	"os"
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

		// Calculate average prosek of accepted applications for this dorm
		avgProsekPipeline := mongo.Pipeline{
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
			{{Key: "$group", Value: bson.D{
				{Key: "_id", Value: nil},
				{Key: "avg_prosek", Value: bson.D{{Key: "$avg", Value: "$prosek"}}},
			}}},
		}

		avgProsekCursor, err := s.prihvaceneAplikacijeCollection.Aggregate(ctx, avgProsekPipeline)
		if err != nil {
			return nil, err
		}

		var avgProsekResult []struct {
			AvgProsek float64 `bson:"avg_prosek"`
		}
		if err = avgProsekCursor.All(ctx, &avgProsekResult); err != nil {
			avgProsekCursor.Close(ctx)
			return nil, err
		}
		avgProsekCursor.Close(ctx)

		if len(avgProsekResult) > 0 && avgProsekResult[0].AvgProsek > 0 {
			stat.AverageProsek = math.Round(avgProsekResult[0].AvgProsek*100) / 100
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
	case "statistics", "dorm-statistics": // Support both old and new names
		return s.exportStatistics(ctx, format)
	case "application-analytics", "application-list": // Support both old and new names
		return s.exportApplicationAnalytics(ctx, format)
	case "accepted-applications":
		return s.exportAcceptedApplications(ctx, format)
	case "yearly-trends", "godisnja-kretanja":
		return s.exportYearlyTrends(ctx, format)
	case "dorm-trends":
		return s.exportDormTrends(ctx, format)
	case "amenities-report":
		return s.exportAmenitiesReport(ctx, format)
	case "occupancy-report":
		return s.exportOccupancyReport(ctx, format)
	case "room-types":
		return s.exportRoomTypes(ctx, format)
	default:
		return nil, fmt.Errorf("unknown dataset: %s", dataset)
	}
}

// exportDorms exports dorm data (without database IDs, in Serbocroatian)
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
		// Create a custom structure without ID for JSON
		type DormExport struct {
			Ime             string    `json:"ime"`
			Adresa          string    `json:"adresa"`
			BrojTelefona    string    `json:"broj_telefona"`
			Email           string    `json:"email"`
			DatumKreiranja  time.Time `json:"datum_kreiranja"`
		}
		
		var exports []DormExport
		for _, dorm := range dorms {
			exports = append(exports, DormExport{
				Ime:            dorm.Ime,
				Adresa:         dorm.Address,
				BrojTelefona:   dorm.TelephoneNumber,
				Email:          dorm.Email,
				DatumKreiranja: dorm.CreatedAt,
			})
		}
		return exports, nil
	}

	// CSV format with Serbocroatian headers
	var csvData [][]string
	csvData = append(csvData, []string{"Ime", "Adresa", "Broj Telefona", "Email", "Datum Kreiranja"})

	for _, dorm := range dorms {
		csvData = append(csvData, []string{
			dorm.Ime,
			dorm.Address,
			dorm.TelephoneNumber,
			dorm.Email,
			dorm.CreatedAt.Format("2006-01-02 15:04:05"),
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
		// Create a filtered version without room_id and dorm_id
		var filteredRooms []map[string]interface{}
		for _, room := range rooms {
			filteredRooms = append(filteredRooms, map[string]interface{}{
				"dorm_name":        room.DormName,
				"dorm_address":     room.DormAddress,
				"capacity":         room.Capacity,
				"occupied":         room.Occupied,
				"available_spots":  room.AvailableSpots,
				"amenities":        room.Amenities,
				"is_available":     room.IsAvailable,
			})
		}
		return filteredRooms, nil
	}

	// CSV format with Serbocroatian headers (without IDs)
	var csvData [][]string
	csvData = append(csvData, []string{"Ime Doma", "Adresa", "Kapacitet", "Popunjeno", "Slobodno", "Pogodnosti", "Dostupna"})

	for _, room := range rooms {
		availableStr := "Ne"
		if room.IsAvailable {
			availableStr = "Da"
		}
		
		csvData = append(csvData, []string{
			room.DormName,
			room.DormAddress,
			fmt.Sprintf("%d", room.Capacity),
			fmt.Sprintf("%d", room.Occupied),
			fmt.Sprintf("%d", room.AvailableSpots),
			strings.Join(room.Amenities, ", "),
			availableStr,
		})
	}

	return csvData, nil
}

// exportStatistics exports statistical data (now called dorm-statistics)
func (s *OpenDataService) exportStatistics(ctx context.Context, format models.ExportFormat) (interface{}, error) {
	stats, err := s.GetPublicStatistics()
	if err != nil {
		return nil, err
	}

	if format == models.ExportFormatJSON {
		return stats, nil
	}

	// CSV format - export dorm statistics with Serbocroatian headers
	var csvData [][]string
	csvData = append(csvData, []string{"Ime Doma", "Adresa", "Ukupno Soba", "Ukupan Kapacitet", "Popunjeno", "Dostupno", "Stopa Popunjenosti (%)", "Prosečan Prosek"})

	for _, dormStat := range stats.DormStatistics {
		csvData = append(csvData, []string{
			dormStat.DormName,
			dormStat.Address,
			fmt.Sprintf("%d", dormStat.TotalRooms),
			fmt.Sprintf("%d", dormStat.TotalCapacity),
			fmt.Sprintf("%d", dormStat.OccupiedSpots),
			fmt.Sprintf("%d", dormStat.AvailableSpots),
			fmt.Sprintf("%.2f", dormStat.OccupancyRate),
			fmt.Sprintf("%.2f", dormStat.AverageProsek),
		})
	}

	return csvData, nil
}

// exportYearlyTrends exports yearly trends data
func (s *OpenDataService) exportYearlyTrends(ctx context.Context, format models.ExportFormat) (interface{}, error) {
	trends, err := s.GetApplicationTrends()
	if err != nil {
		return nil, err
	}

	if format == models.ExportFormatJSON {
		// Create filtered version without min/max grades
		var filteredTrends []map[string]interface{}
		for _, year := range trends.YearlyTrends {
			filteredTrends = append(filteredTrends, map[string]interface{}{
				"akademska_godina":    year.AcademicYear,
				"prijave":             year.TotalApplications,
				"prihvaceno":          year.AcceptedApplications,
				"stopa_prihvatanja":   year.AcceptanceRate,
				"prosecan_prosek":     year.AverageGrade,
			})
		}
		return map[string]interface{}{
			"godisnji_kretanja": filteredTrends,
			"generisano":        time.Now(),
		}, nil
	}

	// CSV format with Serbocroatian headers (without Min. Prosek and Maks. Prosek)
	var csvData [][]string
	csvData = append(csvData, []string{"Školska Godina", "Prijave", "Prihvaćeno", "Stopa Prihvatanja (%)", "Prosečan Prosek"})

	for _, year := range trends.YearlyTrends {
		csvData = append(csvData, []string{
			year.AcademicYear,
			fmt.Sprintf("%d", year.TotalApplications),
			fmt.Sprintf("%d", year.AcceptedApplications),
			fmt.Sprintf("%.2f", year.AcceptanceRate),
			fmt.Sprintf("%.2f", year.AverageGrade),
		})
	}

	return csvData, nil
}

// exportDormTrends exports dorm application trends data
func (s *OpenDataService) exportDormTrends(ctx context.Context, format models.ExportFormat) (interface{}, error) {
	trends, err := s.GetApplicationTrends()
	if err != nil {
		return nil, err
	}

	if format == models.ExportFormatJSON {
		// Create filtered version without IDs
		var filteredTrends []map[string]interface{}
		for _, dorm := range trends.DormTrends {
			filteredTrends = append(filteredTrends, map[string]interface{}{
				"ime_doma":          dorm.DormName,
				"ukupno_prijava":    dorm.TotalApplications,
				"prihvaceno":        dorm.AcceptedApplications,
				"stopa_prihvatanja": dorm.AcceptanceRate,
			})
		}
		return map[string]interface{}{
			"kretanja_po_domovima": filteredTrends,
			"generisano":           time.Now(),
		}, nil
	}

	// CSV format with Serbocroatian headers
	var csvData [][]string
	csvData = append(csvData, []string{"Ime Doma", "Ukupno Prijava", "Prihvaćeno", "Stopa Prihvatanja (%)"})

	for _, dorm := range trends.DormTrends {
		csvData = append(csvData, []string{
			dorm.DormName,
			fmt.Sprintf("%d", dorm.TotalApplications),
			fmt.Sprintf("%d", dorm.AcceptedApplications),
			fmt.Sprintf("%.2f", dorm.AcceptanceRate),
		})
	}

	return csvData, nil
}

// exportAcceptedApplications exports all accepted applications (without database IDs, in Serbocroatian)
func (s *OpenDataService) exportAcceptedApplications(ctx context.Context, format models.ExportFormat) (interface{}, error) {
	// Get all accepted applications with room and dorm information via aggregation
	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "sobas"},
			{Key: "localField", Value: "soba_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "room"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$room"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "st_doms"},
			{Key: "localField", Value: "room.st_dom_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "dorm"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$dorm"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
	}

	cursor, err := s.prihvaceneAplikacijeCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if format == models.ExportFormatJSON {
		// Create filtered version without database IDs
		var filteredApps []map[string]interface{}
		for _, result := range results {
			brojIndexa, _ := result["broj_indexa"].(string)
			prosek := result["prosek"]
			academicYear, _ := result["academic_year"].(string)
			createdAt, _ := result["created_at"].(primitive.DateTime)
			
			// Get dorm name
			dormName := "N/A"
			if dorm, ok := result["dorm"].(bson.M); ok {
				if ime, ok := dorm["ime"].(string); ok {
					dormName = ime
				}
			}
			
			// Convert prosek to float64
			var prosekFloat float64
			switch v := prosek.(type) {
			case float64:
				prosekFloat = v
			case float32:
				prosekFloat = float64(v)
			case int:
				prosekFloat = float64(v)
			case int32:
				prosekFloat = float64(v)
			case int64:
				prosekFloat = float64(v)
			}

			filteredApps = append(filteredApps, map[string]interface{}{
				"broj_indexa":    brojIndexa,
				"prosek":         prosekFloat,
				"ime_doma":       dormName,
				"academic_year":  academicYear,
				"kreirana":       createdAt.Time(),
			})
		}
		return filteredApps, nil
	}

	// CSV format with Serbocroatian headers (without IDs, without Datum Ažuriranja)
	var csvData [][]string
	csvData = append(csvData, []string{"Broj Indexa", "Prosek", "Ime Doma", "Akademska Godina", "Datum Kreiranja"})

	for _, result := range results {
		brojIndexa, _ := result["broj_indexa"].(string)
		prosek := result["prosek"]
		academicYear, _ := result["academic_year"].(string)
		createdAt, _ := result["created_at"].(primitive.DateTime)
		
		// Get dorm name
		dormName := "N/A"
		if dorm, ok := result["dorm"].(bson.M); ok {
			if ime, ok := dorm["ime"].(string); ok {
				dormName = ime
			}
		}
		
		// Convert prosek to float64 and format properly
		var prosekStr string
		switch v := prosek.(type) {
		case float64:
			prosekStr = fmt.Sprintf("%.2f", v)
		case float32:
			prosekStr = fmt.Sprintf("%.2f", v)
		case int:
			prosekStr = fmt.Sprintf("%.2f", float64(v))
		case int32:
			prosekStr = fmt.Sprintf("%.2f", float64(v))
		case int64:
			prosekStr = fmt.Sprintf("%.2f", float64(v))
		default:
			prosekStr = "0.00"
		}

		csvData = append(csvData, []string{
			brojIndexa,
			prosekStr,
			dormName,
			academicYear,
			createdAt.Time().Format("2006-01-02 15:04:05"),
		})
	}

	return csvData, nil
}

// exportApplicationAnalytics exports detailed application list
// Note: "Aktivna" field means whether the application is still pending (true) or has been processed/closed (false)
func (s *OpenDataService) exportApplicationAnalytics(ctx context.Context, format models.ExportFormat) (interface{}, error) {
	// Get all applications with room and dorm information
	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "sobas"},
			{Key: "localField", Value: "soba_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "room"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$room"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "st_doms"},
			{Key: "localField", Value: "room.st_dom_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "dorm"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$dorm"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
	}

	cursor, err := s.aplikacijeCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	type ApplicationAnalytic struct {
		BrojIndexa   string    `bson:"broj_indexa" json:"broj_indexa"`
		DormName     string    `json:"ime_doma"`
		RoomCapacity int       `json:"kapacitet_sobe"`
		Grade        int       `bson:"prosek" json:"prosek"`
		IsActive     bool      `bson:"is_active" json:"aktivna"`
		CreatedAt    time.Time `bson:"created_at" json:"datum_kreiranja"`
	}

	var results []struct {
		BrojIndexa   string                 `bson:"broj_indexa"`
		Grade        int                    `bson:"prosek"`
		IsActive     bool                   `bson:"is_active"`
		CreatedAt    time.Time              `bson:"created_at"`
		Room         map[string]interface{} `bson:"room"`
		Dorm         map[string]interface{} `bson:"dorm"`
	}
	
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	var analytics []ApplicationAnalytic
	for _, result := range results {
		analytic := ApplicationAnalytic{
			BrojIndexa: result.BrojIndexa,
			Grade:      result.Grade,
			IsActive:   result.IsActive,
			CreatedAt:  result.CreatedAt,
		}

		// Extract dorm name from nested document
		if result.Dorm != nil {
			if dormName, ok := result.Dorm["ime"].(string); ok {
				analytic.DormName = dormName
			}
		}

		// Extract room capacity from nested document
		if result.Room != nil {
			if krevetnost, ok := result.Room["krevetnost"].(int32); ok {
				analytic.RoomCapacity = int(krevetnost)
			} else if krevetnost, ok := result.Room["krevetnost"].(int64); ok {
				analytic.RoomCapacity = int(krevetnost)
			} else if krevetnost, ok := result.Room["krevetnost"].(int); ok {
				analytic.RoomCapacity = krevetnost
			}
		}

		analytics = append(analytics, analytic)
	}

	if format == models.ExportFormatJSON {
		return analytics, nil
	}

	// CSV format
	var csvData [][]string
	csvData = append(csvData, []string{"Broj Indexa", "Ime Doma", "Kapacitet Sobe", "Prosek", "Aktivna (da/ne)", "Datum Kreiranja"})

	for _, app := range analytics {
		isActiveStr := "ne"
		if app.IsActive {
			isActiveStr = "da"
		}
		
		csvData = append(csvData, []string{
			app.BrojIndexa,
			app.DormName,
			fmt.Sprintf("%d", app.RoomCapacity),
			fmt.Sprintf("%d", app.Grade),
			isActiveStr,
			app.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return csvData, nil
}

// exportAmenitiesReport exports amenities distribution report
func (s *OpenDataService) exportAmenitiesReport(ctx context.Context, format models.ExportFormat) (interface{}, error) {
	stats, err := s.GetPublicStatistics()
	if err != nil {
		return nil, err
	}

	type AmenityReport struct {
		Amenity      string `json:"amenity"`
		TotalRooms   int    `json:"total_rooms"`
		Percentage   float64 `json:"percentage"`
	}

	var amenityReports []AmenityReport
	totalRooms := stats.TotalRooms

	for amenity, count := range stats.AmenitiesDistribution {
		percentage := 0.0
		if totalRooms > 0 {
			percentage = float64(count) / float64(totalRooms) * 100
			percentage = math.Round(percentage*100) / 100
		}

		amenityReports = append(amenityReports, AmenityReport{
			Amenity:    amenity,
			TotalRooms: count,
			Percentage: percentage,
		})
	}

	// Sort by total rooms descending
	sort.Slice(amenityReports, func(i, j int) bool {
		return amenityReports[i].TotalRooms > amenityReports[j].TotalRooms
	})

	if format == models.ExportFormatJSON {
		return map[string]interface{}{
			"amenities":   amenityReports,
			"total_rooms": totalRooms,
			"generated_at": time.Now(),
		}, nil
	}

	// CSV format
	var csvData [][]string
	csvData = append(csvData, []string{"Amenity", "Rooms with Amenity", "Percentage of Total Rooms"})

	for _, report := range amenityReports {
		csvData = append(csvData, []string{
			report.Amenity,
			fmt.Sprintf("%d", report.TotalRooms),
			fmt.Sprintf("%.2f%%", report.Percentage),
		})
	}

	return csvData, nil
}

// exportOccupancyReport exports detailed occupancy analysis
func (s *OpenDataService) exportOccupancyReport(ctx context.Context, format models.ExportFormat) (interface{}, error) {
	heatmap, err := s.GetOccupancyHeatmap()
	if err != nil {
		return nil, err
	}

	type OccupancyDetail struct {
		DormID          string  `json:"dorm_id"`
		DormName        string  `json:"dorm_name"`
		Address         string  `json:"address"`
		TotalCapacity   int     `json:"total_capacity"`
		OccupiedSpots   int     `json:"occupied_spots"`
		AvailableSpots  int     `json:"available_spots"`
		OccupancyRate   float64 `json:"occupancy_rate"`
		Status          string  `json:"status"`
	}

	var occupancyDetails []OccupancyDetail

	for _, dorm := range heatmap.Dorms {
		occupancyDetails = append(occupancyDetails, OccupancyDetail{
			DormID:         dorm.DormID.Hex(),
			DormName:       dorm.DormName,
			Address:        dorm.Address,
			TotalCapacity:  dorm.TotalCapacity,
			OccupiedSpots:  dorm.OccupiedSpots,
			AvailableSpots: dorm.AvailableSpots,
			OccupancyRate:  dorm.OccupancyRate,
			Status:         dorm.Status,
		})
	}

	if format == models.ExportFormatJSON {
		return map[string]interface{}{
			"dorms":      occupancyDetails,
			"summary":    heatmap.Summary,
			"generated_at": time.Now(),
		}, nil
	}

	// CSV format
	var csvData [][]string
	csvData = append(csvData, []string{"Dorm ID", "Dorm Name", "Address", "Total Capacity", "Occupied", "Available", "Occupancy Rate", "Status"})

	for _, detail := range occupancyDetails {
		csvData = append(csvData, []string{
			detail.DormID,
			detail.DormName,
			detail.Address,
			fmt.Sprintf("%d", detail.TotalCapacity),
			fmt.Sprintf("%d", detail.OccupiedSpots),
			fmt.Sprintf("%d", detail.AvailableSpots),
			fmt.Sprintf("%.2f%%", detail.OccupancyRate),
			detail.Status,
		})
	}

	return csvData, nil
}

// exportRoomTypes exports room types and availability breakdown
func (s *OpenDataService) exportRoomTypes(ctx context.Context, format models.ExportFormat) (interface{}, error) {
	stats, err := s.GetPublicStatistics()
	if err != nil {
		return nil, err
	}

	type RoomTypeDetail struct {
		Capacity          int     `json:"kapacitet"`
		TotalRooms        int     `json:"ukupno_soba"`
		TotalCapacity     int     `json:"ukupan_kapacitet"`
		OccupiedSpots     int     `json:"popunjeno_mesta"`
		AvailableSpots    int     `json:"dostupno_mesta"`
		OccupancyRate     float64 `json:"stopa_popunjenosti"`
	}

	var roomTypeDetails []RoomTypeDetail

	// Get room type statistics with occupancy
	for capacity, roomCount := range stats.RoomTypeDistribution {
		totalCap := capacity * roomCount

		// Count occupied spots for this room type
		pipeline := mongo.Pipeline{
			{{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "sobas"},
				{Key: "localField", Value: "soba_id"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "room"},
			}}},
			{{Key: "$unwind", Value: "$room"}},
			{{Key: "$match", Value: bson.D{
				{Key: "room.krevetnost", Value: capacity},
			}}},
			{{Key: "$count", Value: "total"}},
		}

		cursor, err := s.prihvaceneAplikacijeCollection.Aggregate(ctx, pipeline)
		if err != nil {
			return nil, err
		}

		var result []struct {
			Total int `bson:"total"`
		}
		if err = cursor.All(ctx, &result); err != nil {
			cursor.Close(ctx)
			return nil, err
		}
		cursor.Close(ctx)

		occupied := 0
		if len(result) > 0 {
			occupied = result[0].Total
		}

		available := totalCap - occupied
		occupancyRate := 0.0
		if totalCap > 0 {
			occupancyRate = float64(occupied) / float64(totalCap) * 100
			occupancyRate = math.Round(occupancyRate*100) / 100
		}

		roomTypeDetails = append(roomTypeDetails, RoomTypeDetail{
			Capacity:       capacity,
			TotalRooms:     roomCount,
			TotalCapacity:  totalCap,
			OccupiedSpots:  occupied,
			AvailableSpots: available,
			OccupancyRate:  occupancyRate,
		})
	}

	// Sort by capacity
	sort.Slice(roomTypeDetails, func(i, j int) bool {
		return roomTypeDetails[i].Capacity < roomTypeDetails[j].Capacity
	})

	if format == models.ExportFormatJSON {
		return map[string]interface{}{
			"tipovi_soba":  roomTypeDetails,
			"generisano": time.Now(),
		}, nil
	}

	// CSV format with Serbocroatian headers
	var csvData [][]string
	csvData = append(csvData, []string{"Kapacitet Sobe", "Ukupno Soba", "Ukupan Kapacitet", "Popunjeno Mesta", "Dostupno Mesta", "Stopa Popunjenosti (%)"})

	for _, detail := range roomTypeDetails {
		csvData = append(csvData, []string{
			fmt.Sprintf("Soba za %d osoba", detail.Capacity),
			fmt.Sprintf("%d", detail.TotalRooms),
			fmt.Sprintf("%d", detail.TotalCapacity),
			fmt.Sprintf("%d", detail.OccupiedSpots),
			fmt.Sprintf("%d", detail.AvailableSpots),
			fmt.Sprintf("%.2f", detail.OccupancyRate),
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

// ====================
// Helper Functions
// ====================

// GetRoomApplicationsFromStDomService fetches accepted applications for a room from st_dom_service
func (s *OpenDataService) GetRoomApplicationsFromStDomService(roomID string, authHeader string) (interface{}, error) {
	// Get ST_DOM_SERVICE_URL from environment variable
	stDomServiceURL := os.Getenv("ST_DOM_SERVICE_URL")
	if stDomServiceURL == "" {
		stDomServiceURL = "http://st_dom_service:8081" // default for Docker
	}

	// Construct the endpoint URL
	url := fmt.Sprintf("%s/api/v1/prihvacene_aplikacije/room/%s", stDomServiceURL, roomID)

	// Create HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set the authorization header
	req.Header.Set("Authorization", authHeader)

	// Make the request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call st_dom_service: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("st_dom_service returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse JSON response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return result, nil
}

// GetApplicationsByAcademicYearFromStDomService fetches accepted applications by academic year from st_dom_service
// Now accessible without authentication for public open data (proxies to st_dom_service)
func (s *OpenDataService) GetApplicationsByAcademicYearFromStDomService(academicYear string) (interface{}, error) {
	// Get ST_DOM_SERVICE_URL from environment variable
	stDomServiceURL := os.Getenv("ST_DOM_SERVICE_URL")
	if stDomServiceURL == "" {
		stDomServiceURL = "http://st_dom_service:8081" // default for Docker
	}

	// Construct the endpoint URL with query parameter
	url := fmt.Sprintf("%s/api/v1/prihvacene_aplikacije/academic_year?academic_year=%s", stDomServiceURL, academicYear)

	// Create HTTP request (no auth required - public endpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Make the request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call st_dom_service: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("st_dom_service returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse JSON response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return result, nil
}

// FormatJSON formats data to JSON string
func FormatJSON(data interface{}) (string, error) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
