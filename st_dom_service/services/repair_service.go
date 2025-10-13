package services

import (
	"context"
	"fmt"
	"st_dom_service/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RepairService struct {
	collection *mongo.Collection
}

func NewRepairService(db *mongo.Database) *RepairService {
	return &RepairService{
		collection: db.Collection("repairs"),
	}
}

// CreateRepair creates a new repair schedule
func (s *RepairService) CreateRepair(ctx context.Context, sobaID primitive.ObjectID, description string, estimatedCompletionDate time.Time, createdBy primitive.ObjectID) (*models.Repair, error) {
	repair := &models.Repair{
		ID:                      primitive.NewObjectID(),
		SobaID:                  sobaID,
		Description:             description,
		EstimatedCompletionDate: estimatedCompletionDate,
		CreatedBy:               createdBy,
		Status:                  "scheduled",
		CreatedAt:               time.Now(),
		UpdatedAt:               time.Now(),
	}

	_, err := s.collection.InsertOne(ctx, repair)
	if err != nil {
		return nil, err
	}

	return repair, nil
}

// GetAllRepairs retrieves all repairs
func (s *RepairService) GetAllRepairs(ctx context.Context) ([]models.Repair, error) {
	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var repairs []models.Repair
	if err = cursor.All(ctx, &repairs); err != nil {
		return nil, err
	}

	return repairs, nil
}

// GetRepairByID retrieves a repair by ID
func (s *RepairService) GetRepairByID(ctx context.Context, id primitive.ObjectID) (*models.Repair, error) {
	var repair models.Repair
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&repair)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("repair not found")
		}
		return nil, err
	}

	return &repair, nil
}

// GetRepairsByRoom retrieves all repairs for a specific room
func (s *RepairService) GetRepairsByRoom(ctx context.Context, sobaID primitive.ObjectID) ([]models.Repair, error) {
	cursor, err := s.collection.Find(ctx, bson.M{"soba_id": sobaID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var repairs []models.Repair
	if err = cursor.All(ctx, &repairs); err != nil {
		return nil, err
	}

	return repairs, nil
}

// GetRepairsByStatus retrieves all repairs with a specific status
func (s *RepairService) GetRepairsByStatus(ctx context.Context, status string) ([]models.Repair, error) {
	cursor, err := s.collection.Find(ctx, bson.M{"status": status})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var repairs []models.Repair
	if err = cursor.All(ctx, &repairs); err != nil {
		return nil, err
	}

	return repairs, nil
}

// UpdateRepair updates a repair
func (s *RepairService) UpdateRepair(ctx context.Context, id primitive.ObjectID, description string, estimatedCompletionDate *time.Time, status string) (*models.Repair, error) {
	update := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	if description != "" {
		update["$set"].(bson.M)["description"] = description
	}
	if estimatedCompletionDate != nil {
		update["$set"].(bson.M)["estimated_completion_date"] = *estimatedCompletionDate
	}
	if status != "" {
		update["$set"].(bson.M)["status"] = status
	}

	_, err := s.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return nil, err
	}

	return s.GetRepairByID(ctx, id)
}

// DeleteRepair deletes a repair
func (s *RepairService) DeleteRepair(ctx context.Context, id primitive.ObjectID) error {
	_, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

