package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Repair represents a scheduled repair for a room
type Repair struct {
	ID                      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	SobaID                  primitive.ObjectID `bson:"soba_id" json:"soba_id"`
	Description             string             `bson:"description" json:"description"`
	EstimatedCompletionDate time.Time          `bson:"estimated_completion_date" json:"estimated_completion_date"`
	CreatedBy               primitive.ObjectID `bson:"created_by" json:"created_by"` // Admin user ID
	Status                  string             `bson:"status" json:"status"`         // "scheduled", "in_progress", "completed", "cancelled"
	CreatedAt               time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt               time.Time          `bson:"updated_at" json:"updated_at"`
}

// CreateRepairRequest represents the request to create a repair
type CreateRepairRequest struct {
	SobaID                  string `json:"soba_id" binding:"required"`
	Description             string `json:"description" binding:"required"`
	EstimatedCompletionDate string `json:"estimated_completion_date" binding:"required"` // ISO 8601 format
}

// UpdateRepairRequest represents the request to update a repair
type UpdateRepairRequest struct {
	Description             string `json:"description"`
	EstimatedCompletionDate string `json:"estimated_completion_date"`
	Status                  string `json:"status"`
}

