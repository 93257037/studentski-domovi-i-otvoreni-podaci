package services

import (
	"context"
	"errors"
	"st_dom_service/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PrihvacenaAplikacijaService handles accepted application operations
type PrihvacenaAplikacijaService struct {
	collection          *mongo.Collection
	aplikacijaService   *AplikacijaService
}

// NewPrihvacenaAplikacijaService creates a new PrihvacenaAplikacijaService
func NewPrihvacenaAplikacijaService(collection *mongo.Collection, aplikacijaService *AplikacijaService) *PrihvacenaAplikacijaService {
	return &PrihvacenaAplikacijaService{
		collection:        collection,
		aplikacijaService: aplikacijaService,
	}
}

// ApproveAplikacija approves an application and creates a PrihvacenaAplikacija entry
func (s *PrihvacenaAplikacijaService) ApproveAplikacija(req models.ApproveAplikacijaRequest) (*models.PrihvacenaAplikacija, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get the application
	aplikacija, err := s.aplikacijaService.GetAplikacijaByID(req.AplikacijaID)
	if err != nil {
		return nil, err
	}

	// Check if the application is active
	if !aplikacija.IsActive {
		return nil, errors.New("application is not active")
	}

	// Check if this application has already been accepted
	existingAccepted, err := s.GetPrihvacenaAplikacijaByAplikacijaID(req.AplikacijaID)
	if err == nil && existingAccepted != nil {
		return nil, errors.New("application has already been accepted")
	}

	// Check if this user already has an accepted application
	existingUserAccepted, err := s.GetPrihvaceneAplikacijeByUserID(aplikacija.UserID)
	if err == nil && len(existingUserAccepted) > 0 {
		return nil, errors.New("user already has an accepted application")
	}

	// Create the accepted application entry
	prihvacenaAplikacija := models.NewPrihvacenaAplikacija(aplikacija, req.AcademicYear)

	// Insert into database
	result, err := s.collection.InsertOne(ctx, prihvacenaAplikacija)
	if err != nil {
		return nil, err
	}

	prihvacenaAplikacija.ID = result.InsertedID.(primitive.ObjectID)

	// Mark the original application as inactive (accepted)
	isActive := false
	updateReq := models.UpdateAplikacijaRequest{
		IsActive: &isActive,
	}
	_, err = s.aplikacijaService.UpdateAplikacija(aplikacija.ID, updateReq, aplikacija.UserID)
	if err != nil {
		// Log the error but don't fail - the accepted application was created
		// In a production environment, you might want to handle this more carefully
		return &prihvacenaAplikacija, err
	}

	// Void/deny all other active applications from this user
	err = s.VoidAllOtherUserApplications(aplikacija.UserID, aplikacija.ID)
	if err != nil {
		// Log the error but don't fail - the main approval succeeded
		// In production, you might want to handle this differently
	}

	return &prihvacenaAplikacija, nil
}

// GetPrihvacenaAplikacijaByID retrieves an accepted application by ID
func (s *PrihvacenaAplikacijaService) GetPrihvacenaAplikacijaByID(id primitive.ObjectID) (*models.PrihvacenaAplikacija, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var prihvacenaAplikacija models.PrihvacenaAplikacija
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&prihvacenaAplikacija)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("accepted application not found")
		}
		return nil, err
	}

	return &prihvacenaAplikacija, nil
}

// GetPrihvacenaAplikacijaByAplikacijaID retrieves an accepted application by original application ID
func (s *PrihvacenaAplikacijaService) GetPrihvacenaAplikacijaByAplikacijaID(aplikacijaID primitive.ObjectID) (*models.PrihvacenaAplikacija, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var prihvacenaAplikacija models.PrihvacenaAplikacija
	err := s.collection.FindOne(ctx, bson.M{"aplikacija_id": aplikacijaID}).Decode(&prihvacenaAplikacija)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &prihvacenaAplikacija, nil
}

// GetAllPrihvaceneAplikacije retrieves all accepted applications
func (s *PrihvacenaAplikacijaService) GetAllPrihvaceneAplikacije() ([]models.PrihvacenaAplikacija, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var prihvaceneAplikacije []models.PrihvacenaAplikacija
	if err = cursor.All(ctx, &prihvaceneAplikacije); err != nil {
		return nil, err
	}

	return prihvaceneAplikacije, nil
}

// GetPrihvaceneAplikacijeByUserID retrieves accepted applications for a specific user
func (s *PrihvacenaAplikacijaService) GetPrihvaceneAplikacijeByUserID(userID primitive.ObjectID) ([]models.PrihvacenaAplikacija, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := s.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var prihvaceneAplikacije []models.PrihvacenaAplikacija
	if err = cursor.All(ctx, &prihvaceneAplikacije); err != nil {
		return nil, err
	}

	return prihvaceneAplikacije, nil
}

// GetPrihvaceneAplikacijeBySobaID retrieves accepted applications for a specific room
func (s *PrihvacenaAplikacijaService) GetPrihvaceneAplikacijeBySobaID(sobaID primitive.ObjectID) ([]models.PrihvacenaAplikacija, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := s.collection.Find(ctx, bson.M{"soba_id": sobaID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var prihvaceneAplikacije []models.PrihvacenaAplikacija
	if err = cursor.All(ctx, &prihvaceneAplikacije); err != nil {
		return nil, err
	}

	return prihvaceneAplikacije, nil
}

// GetPrihvaceneAplikacijeByAcademicYear retrieves accepted applications for a specific academic year
func (s *PrihvacenaAplikacijaService) GetPrihvaceneAplikacijeByAcademicYear(academicYear string) ([]models.PrihvacenaAplikacija, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := s.collection.Find(ctx, bson.M{"academic_year": academicYear})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var prihvaceneAplikacije []models.PrihvacenaAplikacija
	if err = cursor.All(ctx, &prihvaceneAplikacije); err != nil {
		return nil, err
	}

	return prihvaceneAplikacije, nil
}

// GetTopStudentsByProsek retrieves top N students ranked by their average grade (prosek)
// Returns students in descending order (highest prosek first)
func (s *PrihvacenaAplikacijaService) GetTopStudentsByProsek(limit int) ([]models.PrihvacenaAplikacija, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set default limit if not specified
	if limit <= 0 {
		limit = 10
	}

	// Create options to sort by prosek (descending) and limit results
	findOptions := options.Find().
		SetSort(bson.D{{Key: "prosek", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := s.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var topStudents []models.PrihvacenaAplikacija
	if err = cursor.All(ctx, &topStudents); err != nil {
		return nil, err
	}

	return topStudents, nil
}

// GetTopStudentsByProsekForAcademicYear retrieves top N students for a specific academic year
func (s *PrihvacenaAplikacijaService) GetTopStudentsByProsekForAcademicYear(academicYear string, limit int) ([]models.PrihvacenaAplikacija, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set default limit if not specified
	if limit <= 0 {
		limit = 10
	}

	// Create options to sort by prosek (descending) and limit results
	findOptions := options.Find().
		SetSort(bson.D{{Key: "prosek", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := s.collection.Find(ctx, bson.M{"academic_year": academicYear}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var topStudents []models.PrihvacenaAplikacija
	if err = cursor.All(ctx, &topStudents); err != nil {
		return nil, err
	}

	return topStudents, nil
}

// GetTopStudentsByProsekForRoom retrieves top N students for a specific room
func (s *PrihvacenaAplikacijaService) GetTopStudentsByProsekForRoom(sobaID primitive.ObjectID, limit int) ([]models.PrihvacenaAplikacija, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set default limit if not specified
	if limit <= 0 {
		limit = 10
	}

	// Create options to sort by prosek (descending) and limit results
	findOptions := options.Find().
		SetSort(bson.D{{Key: "prosek", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := s.collection.Find(ctx, bson.M{"soba_id": sobaID}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var topStudents []models.PrihvacenaAplikacija
	if err = cursor.All(ctx, &topStudents); err != nil {
		return nil, err
	}

	return topStudents, nil
}

// VoidAllOtherUserApplications marks all other active applications from a user as inactive
// This is called when one of their applications gets accepted
func (s *PrihvacenaAplikacijaService) VoidAllOtherUserApplications(userID primitive.ObjectID, approvedAplikacijaID primitive.ObjectID) error {
	// Get all active applications for this user
	allUserApplications, err := s.aplikacijaService.GetAplikacijeByUserID(userID)
	if err != nil {
		return err
	}

	// Update each application (except the one being approved) to set is_active = false
	isActive := false
	updateReq := models.UpdateAplikacijaRequest{
		IsActive: &isActive,
	}

	for _, app := range allUserApplications {
		// Skip the application being approved and already inactive ones
		if app.ID == approvedAplikacijaID || !app.IsActive {
			continue
		}

		// Mark this application as inactive (voided)
		_, err = s.aplikacijaService.UpdateAplikacija(app.ID, updateReq, userID)
		if err != nil {
			// Continue even if one fails - we want to try to void all
			continue
		}
	}

	return nil
}

// DeletePrihvacenaAplikacija deletes an accepted application (admin only)
func (s *PrihvacenaAplikacijaService) DeletePrihvacenaAplikacija(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("accepted application not found")
	}

	return nil
}

// EvictStudent evicts a student from their room (admin only)
// This removes the accepted application, freeing up the room spot
func (s *PrihvacenaAplikacijaService) EvictStudent(userID primitive.ObjectID, reason string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find the user's active accepted application
	var prihvacenaAplikacija models.PrihvacenaAplikacija
	err := s.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&prihvacenaAplikacija)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("user does not have an active room assignment")
		}
		return err
	}

	// Delete the accepted application
	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": prihvacenaAplikacija.ID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("failed to evict student")
	}

	return nil
}

// CheckoutStudent allows a student to voluntarily leave their room
func (s *PrihvacenaAplikacijaService) CheckoutStudent(userID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find the user's active accepted application
	var prihvacenaAplikacija models.PrihvacenaAplikacija
	err := s.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&prihvacenaAplikacija)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("you do not have an active room assignment")
		}
		return err
	}

	// Delete the accepted application
	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": prihvacenaAplikacija.ID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("failed to checkout from room")
	}

	return nil
}