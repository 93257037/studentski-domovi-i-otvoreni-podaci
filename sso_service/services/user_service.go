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

// UserService handles user-related operations
type UserService struct {
	collection      *mongo.Collection
	jwtSecret       string
	stDomServiceURL string
}

// NewUserService creates a new UserService
func NewUserService(collection *mongo.Collection, jwtSecret string, stDomServiceURL string) *UserService {
	return &UserService{
		collection:      collection,
		jwtSecret:       jwtSecret,
		stDomServiceURL: stDomServiceURL,
	}
}

// RegisterUser registers a new user
func (s *UserService) RegisterUser(req models.RegisterRequest) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if user already exists
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

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create new user
	user := models.NewUser(req, hashedPassword)

	// Insert user into database
	result, err := s.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return &user, nil
}

// LoginUser authenticates a user and returns a JWT token
func (s *UserService) LoginUser(req models.LoginRequest) (*models.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find user by email
	var user models.User
	err := s.collection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Username, user.Email, user.Role, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	// Clear password from response
	user.Password = ""

	return &models.LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// GetUserByID retrieves a user by ID
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

	// Clear password from response
	user.Password = ""
	return &user, nil
}

// DeleteUser deletes a user account
// First checks with st_dom_service if user has an active room assignment
func (s *UserService) DeleteUser(userID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if user exists
	var user models.User
	err := s.collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("user not found")
		}
		return err
	}

	// Call st_dom_service to check if user has an active room
	hasActiveRoom, err := s.checkUserHasActiveRoom(userID)
	if err != nil {
		return errors.New("failed to verify room status with dormitory service: " + err.Error())
	}

	if hasActiveRoom {
		return errors.New("cannot delete account: you must check out from your assigned room first")
	}

	// Proceed with deletion
	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": userID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}

// checkUserHasActiveRoom calls the st_dom_service to check if user has an active room
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

