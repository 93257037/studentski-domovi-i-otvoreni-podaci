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
	Ime           string             `bson:"ime" json:"ime" binding:"required"`
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
	Ime             string `json:"ime" binding:"required"`
	Address         string `json:"address" binding:"required"`
	TelephoneNumber string `json:"telephone_number" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
}

// UpdateStDomRequest represents the request body for updating a student dormitory
type UpdateStDomRequest struct {
	Ime             *string `json:"ime,omitempty"`
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
		Ime:             req.Ime,
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
	Prosek      int                `bson:"prosek" json:"prosek" binding:"required,min=6,max=10"`
	SobaID      primitive.ObjectID `bson:"soba_id" json:"soba_id" binding:"required"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id" binding:"required"`
	IsActive    bool               `bson:"is_active" json:"is_active"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// CreateAplikacijaRequest represents the request body for creating an application
type CreateAplikacijaRequest struct {
	BrojIndexa string             `json:"broj_indexa" binding:"required"`
	Prosek     int                `json:"prosek" binding:"required,min=6,max=10"`
	SobaID     primitive.ObjectID `json:"soba_id" binding:"required"`
}

// UpdateAplikacijaRequest represents the request body for updating an application
type UpdateAplikacijaRequest struct {
	BrojIndexa *string `json:"broj_indexa,omitempty"`
	Prosek     *int    `json:"prosek,omitempty"`
	IsActive   *bool   `json:"is_active,omitempty"`
}

// NewAplikacija creates a new application with default values
func NewAplikacija(req CreateAplikacijaRequest, userID primitive.ObjectID) Aplikacija {
	return Aplikacija{
		BrojIndexa: req.BrojIndexa,
		Prosek:     req.Prosek,
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

// PaymentStatus represents the payment status enum
type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusPaid    PaymentStatus = "paid"
	PaymentStatusOverdue PaymentStatus = "overdue"
)

// IsValid checks if the PaymentStatus value is valid
func (ps PaymentStatus) IsValid() bool {
	switch ps {
	case PaymentStatusPending, PaymentStatusPaid, PaymentStatusOverdue:
		return true
	}
	return false
}

// Payment represents a payment for a room rental
// A payment is created after an Aplikacija (application) is approved by admin
type Payment struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	AplikacijaID   primitive.ObjectID `bson:"aplikacija_id" json:"aplikacija_id" binding:"required"`
	Amount         float64            `bson:"amount" json:"amount" binding:"required,min=0"`
	PaymentPeriod  string             `bson:"payment_period" json:"payment_period" binding:"required"` // Format: "YYYY-MM"
	Status         PaymentStatus      `bson:"status" json:"status"`
	PaidAt         *time.Time         `bson:"paid_at,omitempty" json:"paid_at,omitempty"`
	DueDate        time.Time          `bson:"due_date" json:"due_date" binding:"required"`
	Notes          string             `bson:"notes,omitempty" json:"notes,omitempty"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}

// CreatePaymentRequest represents the request body for creating a payment
type CreatePaymentRequest struct {
	AplikacijaID  primitive.ObjectID `json:"aplikacija_id" binding:"required"`
	Amount        float64            `json:"amount" binding:"required,min=0"`
	PaymentPeriod string             `json:"payment_period" binding:"required"` // Format: "YYYY-MM"
	DueDate       time.Time          `json:"due_date" binding:"required"`
	Notes         string             `json:"notes,omitempty"`
}

// UpdatePaymentRequest represents the request body for updating a payment
type UpdatePaymentRequest struct {
	Amount        *float64       `json:"amount,omitempty"`
	PaymentPeriod *string        `json:"payment_period,omitempty"`
	Status        *PaymentStatus `json:"status,omitempty"`
	DueDate       *time.Time     `json:"due_date,omitempty"`
	Notes         *string        `json:"notes,omitempty"`
}

// MarkPaymentPaidRequest represents the request body for marking payment as paid
type MarkPaymentPaidRequest struct {
	PaidAt *time.Time `json:"paid_at,omitempty"` // If not provided, use current time
}

// NewPayment creates a new payment with default values
func NewPayment(req CreatePaymentRequest) Payment {
	return Payment{
		AplikacijaID:  req.AplikacijaID,
		Amount:        req.Amount,
		PaymentPeriod: req.PaymentPeriod,
		Status:        PaymentStatusPending,
		DueDate:       req.DueDate,
		Notes:         req.Notes,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// PrihvacenaAplikacija represents an accepted/approved application
// This is created when an admin approves a student's application
type PrihvacenaAplikacija struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	AplikacijaID   primitive.ObjectID `bson:"aplikacija_id" json:"aplikacija_id" binding:"required"`
	UserID         primitive.ObjectID `bson:"user_id" json:"user_id" binding:"required"`
	BrojIndexa     string             `bson:"broj_indexa" json:"broj_indexa" binding:"required"`
	Prosek         int                `bson:"prosek" json:"prosek" binding:"required,min=6,max=10"`
	SobaID         primitive.ObjectID `bson:"soba_id" json:"soba_id" binding:"required"`
	AcademicYear   string             `bson:"academic_year" json:"academic_year"` // Format: "2024/2025"
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}

// ApproveAplikacijaRequest represents the request body for approving an application
type ApproveAplikacijaRequest struct {
	AplikacijaID primitive.ObjectID `json:"aplikacija_id" binding:"required"`
	AcademicYear string             `json:"academic_year" binding:"required"` // Format: "2024/2025"
}

// EvictStudentRequest represents the request body for evicting a student
type EvictStudentRequest struct {
	UserID primitive.ObjectID `json:"user_id" binding:"required"`
	Reason string             `json:"reason" binding:"required"` // Reason for eviction
}

// NewPrihvacenaAplikacija creates a new accepted application from an Aplikacija
func NewPrihvacenaAplikacija(aplikacija *Aplikacija, academicYear string) PrihvacenaAplikacija {
	return PrihvacenaAplikacija{
		AplikacijaID: aplikacija.ID,
		UserID:       aplikacija.UserID,
		BrojIndexa:   aplikacija.BrojIndexa,
		Prosek:       aplikacija.Prosek,
		SobaID:       aplikacija.SobaID,
		AcademicYear: academicYear,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
