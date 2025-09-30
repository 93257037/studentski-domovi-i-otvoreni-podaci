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

// PaymentService handles payment-related operations
type PaymentService struct {
	collection *mongo.Collection
}

// NewPaymentService creates a new PaymentService
func NewPaymentService(collection *mongo.Collection) *PaymentService {
	return &PaymentService{
		collection: collection,
	}
}

// CreatePayment creates a new payment record (admin only)
func (s *PaymentService) CreatePayment(req models.CreatePaymentRequest, aplikacija *models.Aplikacija) (*models.Payment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if payment already exists for this period
	existingPayment, err := s.GetPaymentByAplikacijaAndPeriod(aplikacija.ID, req.PaymentPeriod)
	if err == nil && existingPayment != nil {
		return nil, errors.New("payment already exists for this period")
	}

	// Create new payment
	payment := models.NewPayment(req)

	// Insert into database
	result, err := s.collection.InsertOne(ctx, payment)
	if err != nil {
		return nil, err
	}

	payment.ID = result.InsertedID.(primitive.ObjectID)
	return &payment, nil
}

// GetPaymentByID retrieves a payment by ID
func (s *PaymentService) GetPaymentByID(id primitive.ObjectID) (*models.Payment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var payment models.Payment
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&payment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}

	return &payment, nil
}

// GetPaymentByAplikacijaAndPeriod checks if payment exists for a specific application and period
func (s *PaymentService) GetPaymentByAplikacijaAndPeriod(aplikacijaID primitive.ObjectID, period string) (*models.Payment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var payment models.Payment
	err := s.collection.FindOne(ctx, bson.M{
		"aplikacija_id":  aplikacijaID,
		"payment_period": period,
	}).Decode(&payment)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &payment, nil
}

// GetPaymentsByUserID retrieves all payments for a specific user
// Uses aggregation to join with aplikacije collection
func (s *PaymentService) GetPaymentsByUserID(userID primitive.ObjectID) ([]models.Payment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Aggregation pipeline to join with aplikacije and filter by user_id
	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "aplikacije"},
			{Key: "localField", Value: "aplikacija_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "aplikacija"},
		}}},
		{{Key: "$unwind", Value: "$aplikacija"}},
		{{Key: "$match", Value: bson.D{
			{Key: "aplikacija.user_id", Value: userID},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "aplikacija", Value: 0}, // Remove the joined aplikacija from result
		}}},
	}

	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var payments []models.Payment
	if err = cursor.All(ctx, &payments); err != nil {
		return nil, err
	}

	return payments, nil
}

// GetPaymentsBySobaID retrieves all payments for a specific room
// Uses aggregation to join with aplikacije collection
func (s *PaymentService) GetPaymentsBySobaID(sobaID primitive.ObjectID) ([]models.Payment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Aggregation pipeline to join with aplikacije and filter by soba_id
	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "aplikacije"},
			{Key: "localField", Value: "aplikacija_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "aplikacija"},
		}}},
		{{Key: "$unwind", Value: "$aplikacija"}},
		{{Key: "$match", Value: bson.D{
			{Key: "aplikacija.soba_id", Value: sobaID},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "aplikacija", Value: 0}, // Remove the joined aplikacija from result
		}}},
	}

	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var payments []models.Payment
	if err = cursor.All(ctx, &payments); err != nil {
		return nil, err
	}

	return payments, nil
}

// GetPaymentsByAplikacijaID retrieves all payments for a specific application
func (s *PaymentService) GetPaymentsByAplikacijaID(aplikacijaID primitive.ObjectID) ([]models.Payment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := s.collection.Find(ctx, bson.M{"aplikacija_id": aplikacijaID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var payments []models.Payment
	if err = cursor.All(ctx, &payments); err != nil {
		return nil, err
	}

	return payments, nil
}

// GetAllPayments retrieves all payments (admin only)
func (s *PaymentService) GetAllPayments() ([]models.Payment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var payments []models.Payment
	if err = cursor.All(ctx, &payments); err != nil {
		return nil, err
	}

	return payments, nil
}

// GetPaymentsByStatus retrieves all payments with a specific status (admin only)
func (s *PaymentService) GetPaymentsByStatus(status models.PaymentStatus) ([]models.Payment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := s.collection.Find(ctx, bson.M{"status": status})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var payments []models.Payment
	if err = cursor.All(ctx, &payments); err != nil {
		return nil, err
	}

	return payments, nil
}

// SearchPaymentsByIndex searches payments by student index pattern (admin only)
// Uses aggregation to join with aplikacije and filter by broj_indexa pattern
// If status is provided, also filters by payment status
func (s *PaymentService) SearchPaymentsByIndex(indexPattern string, status *models.PaymentStatus) ([]models.Payment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build aggregation pipeline
	pipeline := mongo.Pipeline{
		// Join with aplikacije collection
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "aplikacije"},
			{Key: "localField", Value: "aplikacija_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "aplikacija"},
		}}},
		{{Key: "$unwind", Value: "$aplikacija"}},
	}

	// Build match conditions
	matchConditions := bson.D{}
	
	// Add index pattern filter (case-insensitive regex)
	if indexPattern != "" {
		matchConditions = append(matchConditions, bson.E{
			Key: "aplikacija.broj_indexa",
			Value: bson.D{{Key: "$regex", Value: "^" + indexPattern}, {Key: "$options", Value: "i"}},
		})
	}

	// Add status filter if provided
	if status != nil {
		matchConditions = append(matchConditions, bson.E{Key: "status", Value: *status})
	}

	// Add match stage if there are conditions
	if len(matchConditions) > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: matchConditions}})
	}

	// Remove the joined aplikacija from result to keep payment structure clean
	pipeline = append(pipeline, bson.D{{Key: "$project", Value: bson.D{
		{Key: "aplikacija", Value: 0},
	}}})

	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var payments []models.Payment
	if err = cursor.All(ctx, &payments); err != nil {
		return nil, err
	}

	return payments, nil
}

// UpdatePayment updates a payment (admin only)
func (s *PaymentService) UpdatePayment(id primitive.ObjectID, req models.UpdatePaymentRequest) (*models.Payment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build update document
	update := bson.M{"$set": bson.M{"updated_at": time.Now()}}

	if req.Amount != nil {
		update["$set"].(bson.M)["amount"] = *req.Amount
	}
	if req.PaymentPeriod != nil {
		update["$set"].(bson.M)["payment_period"] = *req.PaymentPeriod
	}
	if req.Status != nil {
		if !req.Status.IsValid() {
			return nil, errors.New("invalid payment status")
		}
		update["$set"].(bson.M)["status"] = *req.Status
	}
	if req.DueDate != nil {
		update["$set"].(bson.M)["due_date"] = *req.DueDate
	}
	if req.Notes != nil {
		update["$set"].(bson.M)["notes"] = *req.Notes
	}

	// Update the document
	result, err := s.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, errors.New("payment not found")
	}

	// Return updated document
	return s.GetPaymentByID(id)
}

// MarkPaymentAsPaid marks a payment as paid (admin only)
func (s *PaymentService) MarkPaymentAsPaid(id primitive.ObjectID, paidAt *time.Time) (*models.Payment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// If paidAt is not provided, use current time
	if paidAt == nil {
		now := time.Now()
		paidAt = &now
	}

	update := bson.M{
		"$set": bson.M{
			"status":     models.PaymentStatusPaid,
			"paid_at":    paidAt,
			"updated_at": time.Now(),
		},
	}

	result, err := s.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, errors.New("payment not found")
	}

	return s.GetPaymentByID(id)
}

// MarkPaymentAsUnpaid marks a payment as unpaid/pending (admin only)
func (s *PaymentService) MarkPaymentAsUnpaid(id primitive.ObjectID) (*models.Payment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"status":     models.PaymentStatusPending,
			"updated_at": time.Now(),
		},
		"$unset": bson.M{
			"paid_at": "",
		},
	}

	result, err := s.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, errors.New("payment not found")
	}

	return s.GetPaymentByID(id)
}

// DeletePayment deletes a payment (admin only)
func (s *PaymentService) DeletePayment(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("payment not found")
	}

	return nil
}

// UpdateOverduePayments updates payment status to overdue for payments past due date
func (s *PaymentService) UpdateOverduePayments() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find all pending payments with due_date in the past
	filter := bson.M{
		"status": models.PaymentStatusPending,
		"due_date": bson.M{
			"$lt": time.Now(),
		},
	}

	update := bson.M{
		"$set": bson.M{
			"status":     models.PaymentStatusOverdue,
			"updated_at": time.Now(),
		},
	}

	result, err := s.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return 0, err
	}

	return result.ModifiedCount, nil
}
