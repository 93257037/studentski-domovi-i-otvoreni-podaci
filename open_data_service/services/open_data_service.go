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
	sobasCollection  *mongo.Collection
	stDomsCollection *mongo.Collection
}

// NewOpenDataService creates a new OpenDataService
func NewOpenDataService(sobasCollection, stDomsCollection *mongo.Collection) *OpenDataService {
	return &OpenDataService{
		sobasCollection:  sobasCollection,
		stDomsCollection: stDomsCollection,
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

