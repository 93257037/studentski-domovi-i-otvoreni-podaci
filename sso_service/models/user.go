package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User - predstavlja korisnika u sistemu sa svim potrebnim podacima
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username  string             `bson:"username" json:"username" binding:"required"`
	Email     string             `bson:"email" json:"email" binding:"required,email"`
	Password  string             `bson:"password" json:"-"`
	FirstName string             `bson:"first_name" json:"first_name" binding:"required"`
	LastName  string             `bson:"last_name" json:"last_name" binding:"required"`
	Role      string             `bson:"role" json:"role"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// RegisterRequest - zahtev za registraciju novog korisnika
type RegisterRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=20"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// LoginRequest - zahtev za prijavu korisnika
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse - odgovor za uspesnu prijavu sa tokenom i podacima o korisniku
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// kreira novog korisnika sa default vrednostima na osnovu zahteva za registraciju
// postavlja ulogu na "user" i trenutno vreme za kreiranje i azuriranje
func NewUser(req RegisterRequest, hashedPassword string) User {
	return User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

