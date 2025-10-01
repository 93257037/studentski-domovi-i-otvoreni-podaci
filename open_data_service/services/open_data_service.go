package services

import (
	"context"
	"open_data_service/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// OpenDataService handles business logic for open data operations
type OpenDataService struct {
	sobasCollection                *mongo.Collection
	stDomsCollection               *mongo.Collection
	aplikacijeCollection           *mongo.Collection
	prihvaceneAplikacijeCollection *mongo.Collection
}

// NewOpenDataService creates a new OpenDataService
func NewOpenDataService(sobasCollection, stDomsCollection, aplikacijeCollection, prihvaceneAplikacijeCollection *mongo.Collection) *OpenDataService {
	return &OpenDataService{
		sobasCollection:                sobasCollection,
		stDomsCollection:               stDomsCollection,
		aplikacijeCollection:           aplikacijeCollection,
		prihvaceneAplikacijeCollection: prihvaceneAplikacijeCollection,
	}
}

// FilterRoomsByLuksuz filters rooms by any combination of luxury amenities
// Returns rooms that have ALL specified luxury amenities
func (s *OpenDataService) FilterRoomsByLuksuz(luksuzi []models.Luksuzi) ([]models.Soba, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build filter - rooms must contain all specified luxury amenities
	filter := bson.M{}
	if len(luksuzi) > 0 {
		filter["luksuzi"] = bson.M{"$all": luksuzi}
	}

	cursor, err := s.sobasCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sobas []models.Soba
	if err = cursor.All(ctx, &sobas); err != nil {
		return nil, err
	}

	return sobas, nil
}

// FilterRoomsByLuksuzAndStDom filters rooms by luxury amenities and student dormitory
// Returns rooms that have ALL specified luxury amenities and belong to the specified st_dom
func (s *OpenDataService) FilterRoomsByLuksuzAndStDom(luksuzi []models.Luksuzi, stDomID primitive.ObjectID) ([]models.SobaWithStDom, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build filter
	filter := bson.M{
		"st_dom_id": stDomID,
	}
	if len(luksuzi) > 0 {
		filter["luksuzi"] = bson.M{"$all": luksuzi}
	}

	cursor, err := s.sobasCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sobas []models.Soba
	if err = cursor.All(ctx, &sobas); err != nil {
		return nil, err
	}

	// Fetch the st_dom information
	var stDom models.StDom
	err = s.stDomsCollection.FindOne(ctx, bson.M{"_id": stDomID}).Decode(&stDom)
	if err != nil {
		return nil, err
	}

	// Combine room data with st_dom information
	var result []models.SobaWithStDom
	for _, soba := range sobas {
		result = append(result, models.SobaWithStDom{
			ID:         soba.ID,
			StDomID:    soba.StDomID,
			Krevetnost: soba.Krevetnost,
			Luksuzi:    soba.Luksuzi,
			CreatedAt:  soba.CreatedAt,
			UpdatedAt:  soba.UpdatedAt,
			StDom:      &stDom,
		})
	}

	return result, nil
}

// FilterRoomsByKrevetnost filters rooms by bed capacity (krevetnost)
// Supports exact match, min, max, or range filtering
func (s *OpenDataService) FilterRoomsByKrevetnost(exact *int, min *int, max *int) ([]models.Soba, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build filter based on provided parameters
	filter := bson.M{}
	
	if exact != nil {
		// Exact match
		filter["krevetnost"] = *exact
	} else {
		// Range filtering
		krevetnostFilter := bson.M{}
		if min != nil {
			krevetnostFilter["$gte"] = *min
		}
		if max != nil {
			krevetnostFilter["$lte"] = *max
		}
		if len(krevetnostFilter) > 0 {
			filter["krevetnost"] = krevetnostFilter
		}
	}

	cursor, err := s.sobasCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sobas []models.Soba
	if err = cursor.All(ctx, &sobas); err != nil {
		return nil, err
	}

	return sobas, nil
}

// SearchStDomsByAddress searches student dormitories by address using regex pattern matching
// Returns dormitories whose address matches the provided pattern
func (s *OpenDataService) SearchStDomsByAddress(addressPattern string) ([]models.StDom, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Use regex for partial matching (case-insensitive)
	filter := bson.M{
		"address": bson.M{
			"$regex":   addressPattern,
			"$options": "i", // case-insensitive
		},
	}

	cursor, err := s.stDomsCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var stDoms []models.StDom
	if err = cursor.All(ctx, &stDoms); err != nil {
		return nil, err
	}

	return stDoms, nil
}

// SearchStDomsByIme searches student dormitories by name using regex pattern matching
// Returns dormitories whose ime (name) matches the provided pattern
func (s *OpenDataService) SearchStDomsByIme(imePattern string) ([]models.StDom, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Use regex for partial matching (case-insensitive)
	filter := bson.M{
		"ime": bson.M{
			"$regex":   imePattern,
			"$options": "i", // case-insensitive
		},
	}

	cursor, err := s.stDomsCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var stDoms []models.StDom
	if err = cursor.All(ctx, &stDoms); err != nil {
		return nil, err
	}

	return stDoms, nil
}

// AdvancedFilterRooms provides a comprehensive filtering endpoint combining multiple criteria
// Allows filtering by luksuzi, st_dom_id, krevetnost (exact, min, max), and address pattern all at once
func (s *OpenDataService) AdvancedFilterRooms(luksuzi []models.Luksuzi, stDomID *primitive.ObjectID, addressPattern string, exact *int, min *int, max *int) ([]models.SobaWithStDom, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// If address pattern is provided, first find matching dormitories
	var matchingStDomIDs []primitive.ObjectID
	if addressPattern != "" {
		stDomFilter := bson.M{
			"address": bson.M{
				"$regex":   addressPattern,
				"$options": "i", // case-insensitive
			},
		}
		
		stDomsCursor, err := s.stDomsCollection.Find(ctx, stDomFilter)
		if err != nil {
			return nil, err
		}
		defer stDomsCursor.Close(ctx)

		var matchingStDoms []models.StDom
		if err = stDomsCursor.All(ctx, &matchingStDoms); err != nil {
			return nil, err
		}

		// Extract IDs
		for _, stDom := range matchingStDoms {
			matchingStDomIDs = append(matchingStDomIDs, stDom.ID)
		}

		// If no dormitories match the address pattern, return empty result
		if len(matchingStDomIDs) == 0 {
			return []models.SobaWithStDom{}, nil
		}
	}

	// Build comprehensive filter for rooms
	filter := bson.M{}
	
	// Filter by luxury amenities
	if len(luksuzi) > 0 {
		filter["luksuzi"] = bson.M{"$all": luksuzi}
	}
	
	// Filter by student dormitory (either specific ID or address pattern matches)
	if stDomID != nil {
		// Specific dormitory ID provided
		filter["st_dom_id"] = *stDomID
	} else if len(matchingStDomIDs) > 0 {
		// Address pattern provided - filter by matching dormitory IDs
		filter["st_dom_id"] = bson.M{"$in": matchingStDomIDs}
	}
	
	// Filter by krevetnost
	if exact != nil {
		filter["krevetnost"] = *exact
	} else {
		krevetnostFilter := bson.M{}
		if min != nil {
			krevetnostFilter["$gte"] = *min
		}
		if max != nil {
			krevetnostFilter["$lte"] = *max
		}
		if len(krevetnostFilter) > 0 {
			filter["krevetnost"] = krevetnostFilter
		}
	}

	cursor, err := s.sobasCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sobas []models.Soba
	if err = cursor.All(ctx, &sobas); err != nil {
		return nil, err
	}

	// Group rooms by st_dom_id to minimize database queries
	stDomIDsMap := make(map[primitive.ObjectID]bool)
	for _, soba := range sobas {
		stDomIDsMap[soba.StDomID] = true
	}

	// Fetch all relevant st_doms in one query
	var stDomIDs []primitive.ObjectID
	for id := range stDomIDsMap {
		stDomIDs = append(stDomIDs, id)
	}

	stDomsMap := make(map[primitive.ObjectID]*models.StDom)
	if len(stDomIDs) > 0 {
		stDomsCursor, err := s.stDomsCollection.Find(ctx, bson.M{"_id": bson.M{"$in": stDomIDs}})
		if err != nil {
			return nil, err
		}
		defer stDomsCursor.Close(ctx)

		var stDoms []models.StDom
		if err = stDomsCursor.All(ctx, &stDoms); err != nil {
			return nil, err
		}

		for i := range stDoms {
			stDomsMap[stDoms[i].ID] = &stDoms[i]
		}
	}

	// Combine room data with st_dom information
	var result []models.SobaWithStDom
	for _, soba := range sobas {
		stDom := stDomsMap[soba.StDomID]
		result = append(result, models.SobaWithStDom{
			ID:         soba.ID,
			StDomID:    soba.StDomID,
			Krevetnost: soba.Krevetnost,
			Luksuzi:    soba.Luksuzi,
			CreatedAt:  soba.CreatedAt,
			UpdatedAt:  soba.UpdatedAt,
			StDom:      stDom,
		})
	}

	return result, nil
}

// GetAllRooms returns all rooms (for open data access)
func (s *OpenDataService) GetAllRooms() ([]models.Soba, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := s.sobasCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sobas []models.Soba
	if err = cursor.All(ctx, &sobas); err != nil {
		return nil, err
	}

	return sobas, nil
}

// GetAllStDoms returns all student dormitories (for open data access)
func (s *OpenDataService) GetAllStDoms() ([]models.StDom, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := s.stDomsCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var stDoms []models.StDom
	if err = cursor.All(ctx, &stDoms); err != nil {
		return nil, err
	}

	return stDoms, nil
}

// StDomOccupancyStats represents occupancy statistics for a student dormitory
type StDomOccupancyStats struct {
	StDom       models.StDom `json:"st_dom"`
	OccupiedCount int        `json:"occupied_count"`
	TotalCapacity int        `json:"total_capacity"`
	OccupancyRate float64    `json:"occupancy_rate"`
}

// GetTopFullStDoms returns the top 3 most full student dormitories based on the number of residents
func (s *OpenDataService) GetTopFullStDoms() ([]StDomOccupancyStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Aggregation pipeline to count residents per st_dom
	pipeline := mongo.Pipeline{
		// Stage 1: Lookup to join prihvacena_aplikacije with sobas to get st_dom_id
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "sobas"},
			{Key: "localField", Value: "soba_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "soba"},
		}}},
		// Stage 2: Unwind soba array
		bson.D{{Key: "$unwind", Value: "$soba"}},
		// Stage 3: Group by st_dom_id and count residents
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$soba.st_dom_id"},
			{Key: "occupied_count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		// Stage 4: Sort by occupied_count descending
		bson.D{{Key: "$sort", Value: bson.D{{Key: "occupied_count", Value: -1}}}},
		// Stage 5: Limit to top 3
		bson.D{{Key: "$limit", Value: 3}},
	}

	cursor, err := s.prihvaceneAplikacijeCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		ID            primitive.ObjectID `bson:"_id"`
		OccupiedCount int                `bson:"occupied_count"`
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}


	// Fetch st_dom details and calculate total capacity for each
	var stats []StDomOccupancyStats
	for _, result := range results {
		// Get st_dom details
		var stDom models.StDom
		err := s.stDomsCollection.FindOne(ctx, bson.M{"_id": result.ID}).Decode(&stDom)
		if err != nil {
			continue
		}

		// Calculate total capacity (sum of krevetnost for all rooms in this st_dom)
		capacityPipeline := mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.D{{Key: "st_dom_id", Value: result.ID}}}},
			bson.D{{Key: "$group", Value: bson.D{
				{Key: "_id", Value: nil},
				{Key: "total_capacity", Value: bson.D{{Key: "$sum", Value: "$krevetnost"}}},
			}}},
		}

		capacityCursor, err := s.sobasCollection.Aggregate(ctx, capacityPipeline)
		if err != nil {
			continue
		}

		var capacityResult []struct {
			TotalCapacity int `bson:"total_capacity"`
		}
		if err = capacityCursor.All(ctx, &capacityResult); err != nil || len(capacityResult) == 0 {
			capacityCursor.Close(ctx)
			continue
		}
		capacityCursor.Close(ctx)

		totalCapacity := capacityResult[0].TotalCapacity
		occupancyRate := 0.0
		if totalCapacity > 0 {
			occupancyRate = float64(result.OccupiedCount) / float64(totalCapacity) * 100
		}

		stats = append(stats, StDomOccupancyStats{
			StDom:         stDom,
			OccupiedCount: result.OccupiedCount,
			TotalCapacity: totalCapacity,
			OccupancyRate: occupancyRate,
		})
	}

	return stats, nil
}

// GetTopEmptyStDoms returns the top 3 most empty student dormitories based on available capacity
func (s *OpenDataService) GetTopEmptyStDoms() ([]StDomOccupancyStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get all st_doms
	allStDoms, err := s.GetAllStDoms()
	if err != nil {
		return nil, err
	}

	var stats []StDomOccupancyStats

	for _, stDom := range allStDoms {
		// Calculate total capacity for this st_dom
		capacityPipeline := mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.D{{Key: "st_dom_id", Value: stDom.ID}}}},
			bson.D{{Key: "$group", Value: bson.D{
				{Key: "_id", Value: nil},
				{Key: "total_capacity", Value: bson.D{{Key: "$sum", Value: "$krevetnost"}}},
			}}},
		}

		capacityCursor, err := s.sobasCollection.Aggregate(ctx, capacityPipeline)
		if err != nil {
			continue
		}

		var capacityResult []struct {
			TotalCapacity int `bson:"total_capacity"`
		}
		if err = capacityCursor.All(ctx, &capacityResult); err != nil || len(capacityResult) == 0 {
			capacityCursor.Close(ctx)
			continue
		}
		capacityCursor.Close(ctx)
		totalCapacity := capacityResult[0].TotalCapacity

		// Count occupied spots for this st_dom
		occupiedPipeline := mongo.Pipeline{
			bson.D{{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "sobas"},
				{Key: "localField", Value: "soba_id"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "soba"},
			}}},
			bson.D{{Key: "$unwind", Value: "$soba"}},
			bson.D{{Key: "$match", Value: bson.D{{Key: "soba.st_dom_id", Value: stDom.ID}}}},
			bson.D{{Key: "$count", Value: "occupied_count"}},
		}

		occupiedCursor, err := s.prihvaceneAplikacijeCollection.Aggregate(ctx, occupiedPipeline)
		if err != nil {
			continue
		}

		var occupiedResult []struct {
			OccupiedCount int `bson:"occupied_count"`
		}
		occupiedCount := 0
		if err = occupiedCursor.All(ctx, &occupiedResult); err == nil && len(occupiedResult) > 0 {
			occupiedCount = occupiedResult[0].OccupiedCount
		}
		occupiedCursor.Close(ctx)

		occupancyRate := 0.0
		if totalCapacity > 0 {
			occupancyRate = float64(occupiedCount) / float64(totalCapacity) * 100
		}

		stats = append(stats, StDomOccupancyStats{
			StDom:         stDom,
			OccupiedCount: occupiedCount,
			TotalCapacity: totalCapacity,
			OccupancyRate: occupancyRate,
		})
	}

	// Sort by occupied count (ascending) to get the emptiest
	// We'll sort in Go since we already have all the data
	for i := 0; i < len(stats)-1; i++ {
		for j := i + 1; j < len(stats); j++ {
			if stats[i].OccupiedCount > stats[j].OccupiedCount {
				stats[i], stats[j] = stats[j], stats[i]
			}
		}
	}

	// Return top 3
	if len(stats) > 3 {
		stats = stats[:3]
	}

	return stats, nil
}

// StDomApplicationStats represents application statistics for a student dormitory
type StDomApplicationStats struct {
	StDom            models.StDom `json:"st_dom"`
	ApplicationCount int          `json:"application_count"`
}

// GetStDomWithMostApplications returns the student dormitory with the most applications
func (s *OpenDataService) GetStDomWithMostApplications() (*StDomApplicationStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Aggregation pipeline to count active applications per st_dom
	pipeline := mongo.Pipeline{
		// Stage 1: Match only active applications
		bson.D{{Key: "$match", Value: bson.D{{Key: "is_active", Value: true}}}},
		// Stage 2: Lookup to join aplikacije with sobas to get st_dom_id
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "sobas"},
			{Key: "localField", Value: "soba_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "soba"},
		}}},
		// Stage 3: Unwind soba array
		bson.D{{Key: "$unwind", Value: "$soba"}},
		// Stage 4: Group by st_dom_id and count applications
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$soba.st_dom_id"},
			{Key: "application_count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		// Stage 5: Sort by application_count descending
		bson.D{{Key: "$sort", Value: bson.D{{Key: "application_count", Value: -1}}}},
		// Stage 6: Limit to top 1
		bson.D{{Key: "$limit", Value: 1}},
	}

	cursor, err := s.aplikacijeCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		ID               primitive.ObjectID `bson:"_id"`
		ApplicationCount int                `bson:"application_count"`
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}

	// Get st_dom details
	var stDom models.StDom
	err = s.stDomsCollection.FindOne(ctx, bson.M{"_id": results[0].ID}).Decode(&stDom)
	if err != nil {
		return nil, err
	}

	return &StDomApplicationStats{
		StDom:            stDom,
		ApplicationCount: results[0].ApplicationCount,
	}, nil
}

// StDomAverageProsekStats represents average prosek statistics for a student dormitory
type StDomAverageProsekStats struct {
	StDom         models.StDom `json:"st_dom"`
	AverageProsek float64      `json:"average_prosek"`
	ResidentCount int          `json:"resident_count"`
}

// GetStDomWithHighestAverageProsek returns the student dormitory with the highest average prosek
func (s *OpenDataService) GetStDomWithHighestAverageProsek() (*StDomAverageProsekStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Aggregation pipeline to calculate average prosek per st_dom
	pipeline := mongo.Pipeline{
		// Stage 1: Lookup to join prihvacena_aplikacije with sobas to get st_dom_id
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "sobas"},
			{Key: "localField", Value: "soba_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "soba"},
		}}},
		// Stage 2: Unwind soba array
		bson.D{{Key: "$unwind", Value: "$soba"}},
		// Stage 3: Group by st_dom_id and calculate average prosek
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$soba.st_dom_id"},
			{Key: "average_prosek", Value: bson.D{{Key: "$avg", Value: "$prosek"}}},
			{Key: "resident_count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		// Stage 4: Sort by average_prosek descending
		bson.D{{Key: "$sort", Value: bson.D{{Key: "average_prosek", Value: -1}}}},
		// Stage 5: Limit to top 1
		bson.D{{Key: "$limit", Value: 1}},
	}

	cursor, err := s.prihvaceneAplikacijeCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		ID            primitive.ObjectID `bson:"_id"`
		AverageProsek float64            `bson:"average_prosek"`
		ResidentCount int                `bson:"resident_count"`
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}

	// Get st_dom details
	var stDom models.StDom
	err = s.stDomsCollection.FindOne(ctx, bson.M{"_id": results[0].ID}).Decode(&stDom)
	if err != nil {
		return nil, err
	}

	return &StDomAverageProsekStats{
		StDom:         stDom,
		AverageProsek: results[0].AverageProsek,
		ResidentCount: results[0].ResidentCount,
	}, nil
}

