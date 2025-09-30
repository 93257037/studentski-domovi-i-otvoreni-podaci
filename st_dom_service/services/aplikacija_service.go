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

// AplikacijaService handles application-related operations
type AplikacijaService struct {
	collection *mongo.Collection
}

// NewAplikacijaService creates a new AplikacijaService
func NewAplikacijaService(collection *mongo.Collection) *AplikacijaService {
	return &AplikacijaService{
		collection: collection,
	}
}

// CreateAplikacija creates a new application
func (s *AplikacijaService) CreateAplikacija(req models.CreateAplikacijaRequest, userID primitive.ObjectID) (*models.Aplikacija, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if user already has an active application for this room
	existingApp, err := s.GetActiveAplikacijaByUserAndRoom(userID, req.SobaID)
	if err == nil && existingApp != nil {
		return nil, errors.New("user already has an active application for this room")
	}

	// Create new application
	aplikacija := models.NewAplikacija(req, userID)

	// Insert into database
	result, err := s.collection.InsertOne(ctx, aplikacija)
	if err != nil {
		return nil, err
	}

	aplikacija.ID = result.InsertedID.(primitive.ObjectID)
	return &aplikacija, nil
}

// GetAplikacijaByID retrieves an application by ID
func (s *AplikacijaService) GetAplikacijaByID(id primitive.ObjectID) (*models.Aplikacija, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var aplikacija models.Aplikacija
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&aplikacija)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("application not found")
		}
		return nil, err
	}

	return &aplikacija, nil
}

// GetAplikacijeByUserID retrieves all applications for a specific user
func (s *AplikacijaService) GetAplikacijeByUserID(userID primitive.ObjectID) ([]models.Aplikacija, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := s.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var aplikacije []models.Aplikacija
	if err = cursor.All(ctx, &aplikacije); err != nil {
		return nil, err
	}

	return aplikacije, nil
}

// GetActiveAplikacijaByUserAndRoom checks if user has an active application for a specific room
func (s *AplikacijaService) GetActiveAplikacijaByUserAndRoom(userID, sobaID primitive.ObjectID) (*models.Aplikacija, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var aplikacija models.Aplikacija
	err := s.collection.FindOne(ctx, bson.M{
		"user_id":   userID,
		"soba_id":   sobaID,
		"is_active": true,
	}).Decode(&aplikacija)
	
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No active application found
		}
		return nil, err
	}

	return &aplikacija, nil
}

// GetAplikacijeBySobaID retrieves all applications for a specific room
func (s *AplikacijaService) GetAplikacijeBySobaID(sobaID primitive.ObjectID) ([]models.Aplikacija, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := s.collection.Find(ctx, bson.M{"soba_id": sobaID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var aplikacije []models.Aplikacija
	if err = cursor.All(ctx, &aplikacije); err != nil {
		return nil, err
	}

	return aplikacije, nil
}

// UpdateAplikacija updates an application
func (s *AplikacijaService) UpdateAplikacija(id primitive.ObjectID, req models.UpdateAplikacijaRequest, userID primitive.ObjectID) (*models.Aplikacija, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// First, verify that the application belongs to the user
	currentApp, err := s.GetAplikacijaByID(id)
	if err != nil {
		return nil, err
	}

	if currentApp.UserID != userID {
		return nil, errors.New("unauthorized: application does not belong to user")
	}

	// Build update document
	update := bson.M{"$set": bson.M{"updated_at": time.Now()}}
	
	if req.BrojIndexa != nil {
		update["$set"].(bson.M)["broj_indexa"] = *req.BrojIndexa
	}
	if req.IsActive != nil {
		update["$set"].(bson.M)["is_active"] = *req.IsActive
	}

	// Update the document
	result, err := s.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, errors.New("application not found")
	}

	// Return updated document
	return s.GetAplikacijaByID(id)
}

// DeleteAplikacija deletes an application (only by the owner)
func (s *AplikacijaService) DeleteAplikacija(id primitive.ObjectID, userID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// First, verify that the application belongs to the user
	currentApp, err := s.GetAplikacijaByID(id)
	if err != nil {
		return err
	}

	if currentApp.UserID != userID {
		return errors.New("unauthorized: application does not belong to user")
	}

	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("application not found")
	}

	return nil
}

// GetAllAplikacije retrieves all applications (admin only)
func (s *AplikacijaService) GetAllAplikacije() ([]models.Aplikacija, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var aplikacije []models.Aplikacija
	if err = cursor.All(ctx, &aplikacije); err != nil {
		return nil, err
	}

	return aplikacije, nil
}