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

// StDomService handles student dormitory-related operations
type StDomService struct {
	collection *mongo.Collection
}

// NewStDomService creates a new StDomService
func NewStDomService(collection *mongo.Collection) *StDomService {
	return &StDomService{
		collection: collection,
	}
}

// CreateStDom creates a new student dormitory
func (s *StDomService) CreateStDom(req models.CreateStDomRequest) (*models.StDom, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if dormitory with same address already exists
	var existingStDom models.StDom
	err := s.collection.FindOne(ctx, bson.M{"address": req.Address}).Decode(&existingStDom)
	if err == nil {
		return nil, errors.New("student dormitory with this address already exists")
	}
	if err != mongo.ErrNoDocuments {
		return nil, err
	}

	// Create new dormitory
	stDom := models.NewStDom(req)

	// Insert into database
	result, err := s.collection.InsertOne(ctx, stDom)
	if err != nil {
		return nil, err
	}

	stDom.ID = result.InsertedID.(primitive.ObjectID)
	return &stDom, nil
}

// GetStDomByID retrieves a student dormitory by ID
func (s *StDomService) GetStDomByID(id primitive.ObjectID) (*models.StDom, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var stDom models.StDom
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&stDom)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("student dormitory not found")
		}
		return nil, err
	}

	return &stDom, nil
}

// GetAllStDoms retrieves all student dormitories
func (s *StDomService) GetAllStDoms() ([]models.StDom, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := s.collection.Find(ctx, bson.M{})
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

// UpdateStDom updates a student dormitory
func (s *StDomService) UpdateStDom(id primitive.ObjectID, req models.UpdateStDomRequest) (*models.StDom, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build update document
	update := bson.M{"$set": bson.M{"updated_at": time.Now()}}
	if req.Address != nil {
		update["$set"].(bson.M)["address"] = *req.Address
	}
	if req.TelephoneNumber != nil {
		update["$set"].(bson.M)["telephone_number"] = *req.TelephoneNumber
	}
	if req.Email != nil {
		update["$set"].(bson.M)["email"] = *req.Email
	}

	// Check if address is being updated and already exists
	if req.Address != nil {
		var existingStDom models.StDom
		err := s.collection.FindOne(ctx, bson.M{
			"address": *req.Address,
			"_id":     bson.M{"$ne": id},
		}).Decode(&existingStDom)
		if err == nil {
			return nil, errors.New("student dormitory with this address already exists")
		}
		if err != mongo.ErrNoDocuments {
			return nil, err
		}
	}

	// Update the document
	result, err := s.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, errors.New("student dormitory not found")
	}

	// Return updated document
	return s.GetStDomByID(id)
}

// DeleteStDom deletes a student dormitory
func (s *StDomService) DeleteStDom(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("student dormitory not found")
	}

	return nil
}
