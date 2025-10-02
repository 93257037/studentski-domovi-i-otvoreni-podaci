package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sso_service/models"
	"sso_service/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserService - rukuje operacijama vezanim za korisnike
type UserService struct {
	collection      *mongo.Collection
	jwtSecret       string
	stDomServiceURL string
}

// kreira novi UserService sa kolekcijom baze, JWT tajnim kljucem i URL-om st_dom servisa
func NewUserService(collection *mongo.Collection, jwtSecret string, stDomServiceURL string) *UserService {
	return &UserService{
		collection:      collection,
		jwtSecret:       jwtSecret,
		stDomServiceURL: stDomServiceURL,
	}
}

// registruje novog korisnika - proverava da li vec postoji, hesuje lozinku i cuva u bazu
// vraca gresku ako korisnik sa istim email-om ili korisnickim imenom vec postoji
func (s *UserService) RegisterUser(req models.RegisterRequest) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var existingUser models.User
	err := s.collection.FindOne(ctx, bson.M{
		"$or": []bson.M{
			{"email": req.Email},
			{"username": req.Username},
		},
	}).Decode(&existingUser)

	if err == nil {
		return nil, errors.New("user with this email or username already exists")
	}

	if err != mongo.ErrNoDocuments {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := models.NewUser(req, hashedPassword)
	result, err := s.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return &user, nil
}

// prijavljuje korisnika - proverava email i lozinku, generiše JWT token
// vraca token i podatke o korisniku ako su podaci ispravni
func (s *UserService) LoginUser(req models.LoginRequest) (*models.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var user models.User
	err := s.collection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	token, err := utils.GenerateJWT(user.ID, user.Username, user.Email, user.Role, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return &models.LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// dobija korisnika po ID-u iz baze podataka
// vraca podatke o korisniku bez lozinke
func (s *UserService) GetUserByID(userID primitive.ObjectID) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := s.collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user.Password = ""
	return &user, nil
}

// brise korisnikov nalog - prvo poziva st_dom_service da proveri da li ima aktivnu sobu
// ne dozvoljava brisanje ako korisnik ima dodeljenu sobu
func (s *UserService) DeleteUser(userID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var user models.User
	err := s.collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("user not found")
		}
		return err
	}

	hasActiveRoom, err := s.checkUserHasActiveRoom(userID)
	if err != nil {
		return errors.New("failed to verify room status with dormitory service: " + err.Error())
	}

	if hasActiveRoom {
		return errors.New("cannot delete account: you must check out from your assigned room first")
	}

	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": userID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}

// poziva st_dom_service da proveri da li korisnik ima aktivnu sobu
// salje HTTP GET zahtev i parsira odgovor
func (s *UserService) checkUserHasActiveRoom(userID primitive.ObjectID) (bool, error) {
	url := fmt.Sprintf("%s/api/v1/internal/users/%s/room-status", s.stDomServiceURL, userID.Hex())
	
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, errors.New("failed to check room status")
	}

	var result struct {
		HasActiveRoom bool `json:"has_active_room"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	return result.HasActiveRoom, nil
}

