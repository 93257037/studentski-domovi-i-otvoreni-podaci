package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Luksuzi represents the luxury amenities enum
type Luksuzi string

const (
	LuksuziKlima          Luksuzi = "klima"
	LuksuziTerasa         Luksuzi = "terasa"
	LuksuziSopstvenoKupatilo Luksuzi = "sopstveno kupatilo"
	LuksuziStram          Luksuzi = "áram"
	LuksuziAblak          Luksuzi = "ablak"
	LuksuziNeisvrljanzid  Luksuzi = "neisvrljan zid"
)

// IsValid checks if the Luksuzi value is valid
func (l Luksuzi) IsValid() bool {
	switch l {
	case LuksuziKlima, LuksuziTerasa, LuksuziSopstvenoKupatilo, LuksuziStram, LuksuziAblak, LuksuziNeisvrljanzid:
		return true
	}
	return false
}

// StDom represents a student dormitory
type StDom struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Address       string             `bson:"address" json:"address" binding:"required"`
	TelephoneNumber string           `bson:"telephone_number" json:"telephone_number" binding:"required"`
	Email         string             `bson:"email" json:"email" binding:"required,email"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}

// Soba represents a room in a student dormitory
type Soba struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	StDomID       primitive.ObjectID   `bson:"st_dom_id" json:"st_dom_id" binding:"required"`
	Krevetnost    int                  `bson:"krevetnost" json:"krevetnost" binding:"required,min=1"`
	Luksuzi       []Luksuzi            `bson:"luksuzi" json:"luksuzi"`
	CreatedAt     time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time            `bson:"updated_at" json:"updated_at"`
}

// CreateStDomRequest represents the request body for creating a student dormitory
type CreateStDomRequest struct {
	Address         string `json:"address" binding:"required"`
	TelephoneNumber string `json:"telephone_number" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
}

// UpdateStDomRequest represents the request body for updating a student dormitory
type UpdateStDomRequest struct {
	Address         *string `json:"address,omitempty"`
	TelephoneNumber *string `json:"telephone_number,omitempty"`
	Email           *string `json:"email,omitempty"`
}

// CreateSobaRequest represents the request body for creating a room
type CreateSobaRequest struct {
	StDomID       primitive.ObjectID `json:"st_dom_id" binding:"required"`
	Krevetnost    int                `json:"krevetnost" binding:"required,min=1"`
	Luksuzi       []Luksuzi          `json:"luksuzi"`
}

// UpdateSobaRequest represents the request body for updating a room
type UpdateSobaRequest struct {
	Krevetnost    *int       `json:"krevetnost,omitempty"`
	Luksuzi       *[]Luksuzi `json:"luksuzi,omitempty"`
}

// NewStDom creates a new student dormitory with default values
func NewStDom(req CreateStDomRequest) StDom {
	return StDom{
		Address:         req.Address,
		TelephoneNumber: req.TelephoneNumber,
		Email:           req.Email,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

// NewSoba creates a new room with default values
func NewSoba(req CreateSobaRequest) Soba {
	return Soba{
		StDomID:       req.StDomID,
		Krevetnost:    req.Krevetnost,
		Luksuzi:       req.Luksuzi,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// Aplikacija represents a room application
type Aplikacija struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	BrojIndexa  string             `bson:"broj_indexa" json:"broj_indexa" binding:"required"`
	SobaID      primitive.ObjectID `bson:"soba_id" json:"soba_id" binding:"required"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id" binding:"required"`
	IsActive    bool               `bson:"is_active" json:"is_active"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// CreateAplikacijaRequest represents the request body for creating an application
type CreateAplikacijaRequest struct {
	BrojIndexa string             `json:"broj_indexa" binding:"required"`
	SobaID     primitive.ObjectID `json:"soba_id" binding:"required"`
}

// UpdateAplikacijaRequest represents the request body for updating an application
type UpdateAplikacijaRequest struct {
	BrojIndexa *string `json:"broj_indexa,omitempty"`
	IsActive   *bool   `json:"is_active,omitempty"`
}

// NewAplikacija creates a new application with default values
func NewAplikacija(req CreateAplikacijaRequest, userID primitive.ObjectID) Aplikacija {
	return Aplikacija{
		BrojIndexa: req.BrojIndexa,
		SobaID:     req.SobaID,
		UserID:     userID,
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// ValidateLuksuzi validates that all luxury amenities are valid
func ValidateLuksuzi(luksuzi []Luksuzi) bool {
	for _, l := range luksuzi {
		if !l.IsValid() {
			return false
		}
	}
	return true
}
