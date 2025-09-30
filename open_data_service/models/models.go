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
	Ime           string             `bson:"ime" json:"ime"`
	Address       string             `bson:"address" json:"address"`
	TelephoneNumber string           `bson:"telephone_number" json:"telephone_number"`
	Email         string             `bson:"email" json:"email"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}

// Soba represents a room in a student dormitory
type Soba struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	StDomID       primitive.ObjectID   `bson:"st_dom_id" json:"st_dom_id"`
	Krevetnost    int                  `bson:"krevetnost" json:"krevetnost"`
	Luksuzi       []Luksuzi            `bson:"luksuzi" json:"luksuzi"`
	CreatedAt     time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time            `bson:"updated_at" json:"updated_at"`
}

// SobaWithStDom represents a room with its student dormitory information
type SobaWithStDom struct {
	ID            primitive.ObjectID   `json:"id,omitempty"`
	StDomID       primitive.ObjectID   `json:"st_dom_id"`
	Krevetnost    int                  `json:"krevetnost"`
	Luksuzi       []Luksuzi            `json:"luksuzi"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
	StDom         *StDom               `json:"st_dom,omitempty"`
}

