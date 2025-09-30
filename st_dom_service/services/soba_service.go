package services

import (
	"context"
	"errors"
	"st_dom_service/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// SobaService handles room-related operations
type SobaService struct {
	collection *mongo.Collection
}

// NewSobaService creates a new SobaService
func NewSobaService(collection *mongo.Collection) *SobaService {
	return &SobaService{
		collection: collection,
	}
}

// CreateSoba creates a new room
func (s *SobaService) CreateSoba(req models.CreateSobaRequest) (*models.Soba, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Validate luxury amenities
	if !models.ValidateLuksuzi(req.Luksuzi) {
		return nil, errors.New("invalid luxury amenities")
	}

	// Create new room
	soba := models.NewSoba(req)

	// Insert into database
	result, err := s.collection.InsertOne(ctx, soba)
	if err != nil {
		return nil, err
	}

	soba.ID = result.InsertedID.(primitive.ObjectID)
	return &soba, nil
}

// GetSobaByID retrieves a room by ID
func (s *SobaService) GetSobaByID(id primitive.ObjectID) (*models.Soba, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var soba models.Soba
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&soba)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("room not found")
		}
		return nil, err
	}

	return &soba, nil
}

// GetAllSobas retrieves all rooms
func (s *SobaService) GetAllSobas() ([]models.Soba, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := s.collection.Find(ctx, bson.M{})
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

// GetSobasByStDomID retrieves all rooms for a specific dormitory
func (s *SobaService) GetSobasByStDomID(stDomID primitive.ObjectID) ([]models.Soba, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := s.collection.Find(ctx, bson.M{"st_dom_id": stDomID})
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

// UpdateSoba updates a room
func (s *SobaService) UpdateSoba(id primitive.ObjectID, req models.UpdateSobaRequest) (*models.Soba, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build update document
	update := bson.M{"$set": bson.M{"updated_at": time.Now()}}
	
	if req.Krevetnost != nil {
		update["$set"].(bson.M)["krevetnost"] = *req.Krevetnost
	}
	if req.Luksuzi != nil {
		// Validate luxury amenities
		if !models.ValidateLuksuzi(*req.Luksuzi) {
			return nil, errors.New("invalid luxury amenities")
		}
		update["$set"].(bson.M)["luksuzi"] = *req.Luksuzi
	}


	// Update the document
	result, err := s.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, errors.New("room not found")
	}

	// Return updated document
	return s.GetSobaByID(id)
}

// DeleteSoba deletes a room
func (s *SobaService) DeleteSoba(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("room not found")
	}

	return nil
}

// DeleteSobasByStDomID deletes all rooms for a specific dormitory
func (s *SobaService) DeleteSobasByStDomID(stDomID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := s.collection.DeleteMany(ctx, bson.M{"st_dom_id": stDomID})
	return err
}
