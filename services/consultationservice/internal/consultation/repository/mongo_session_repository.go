package repository

import (
	"consultationservice/internal/consultation/domain"
	domain2 "consultationservice/internal/payment/domain"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoSessionRepository struct {
	collection *mongo.Collection
}


func NewMongoSessionRepository(collection *mongo.Collection) *MongoSessionRepository {
	return &MongoSessionRepository{
		collection: collection,
	}
}

func (r *MongoSessionRepository) UpdateStatus(
	ctx context.Context,
	sessionID string,
	status domain.SessionStatus,
) error {

	filter := bson.M{"_id": sessionID}
	now := time.Now().UTC()
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": now,
		},
		"$push": bson.M{
			"logs": domain.SessionLog{
				Status:    status,
				Timestamp: now,
				Message:   "Status updated to " + string(status),
			},
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("session not found")
	}

	return nil
}

func (r *MongoSessionRepository) UpdatePayment(
	ctx context.Context,
	sessionID string,
	paymentStatus domain2.PaymentStatus,
	paymentID *string,
) error {

	filter := bson.M{"_id": sessionID}
	now := time.Now().UTC()

	updateFields := bson.M{
		"payment_status": paymentStatus,
		"updated_at":     now,
	}

	if paymentID != nil {
		updateFields["payment_id"] = paymentID
	}

	update := bson.M{
		"$set": updateFields,
		"$push": bson.M{
			"logs": domain.SessionLog{
				Status:    "payment_update",
				Timestamp: now,
				Message:   "Payment status updated to " + string(paymentStatus),
			},
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("session not found")
	}

	return nil
}
func (r *MongoSessionRepository) AttachCallSession(
	ctx context.Context,
	sessionID string,
	callSessionID string,
) error {

	filter := bson.M{
		"_id":             sessionID,
		"call_session_id": bson.M{"$exists": false},
	}

	update := bson.M{
		"$set": bson.M{
			"call_session_id": callSessionID,
			"updated_at":      time.Now().UTC(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("session not found or call already started")
	}

	return nil
}

func (r *MongoSessionRepository) Cancel(
	ctx context.Context,
	sessionID string,
	cancelMeta domain.CancelData,
) error {

	filter := bson.M{
		"_id": sessionID,
		"status": bson.M{
			"$ne": domain.SessionStatusCompleted,
		},
	}

	now := time.Now().UTC()
	update := bson.M{
		"$set": bson.M{
			"status":           domain.SessionStatusCancelled,
			"cancel_meta_data": cancelMeta,
			"updated_at":       now,
		},
		"$push": bson.M{
			"logs": domain.SessionLog{
				Status:    domain.SessionStatusCancelled,
				Timestamp: now,
				Message:   "Session cancelled",
			},
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("session not found or cannot be cancelled")
	}

	return nil
}

func (r *MongoSessionRepository) Create(ctx context.Context, session *domain.Session) error {
	_, err := r.collection.InsertOne(ctx, session)
	return err
}

func (r *MongoSessionRepository) GetByID(ctx context.Context, id string) (*domain.Session, error) {
	var session domain.Session

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {

		err = r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&session)
	} else {
		err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&session)
	}

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("session not found")
		}
		return nil, err
	}

	return &session, nil
}

func (r *MongoSessionRepository) GetByTherapist(ctx context.Context, therapistID string) ([]*domain.Session, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"therapist_id": therapistID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sessions []*domain.Session
	if err := cursor.All(ctx, &sessions); err != nil {
		return nil, err
	}

	return sessions, nil
}

func (r *MongoSessionRepository) GetByClient(ctx context.Context, clientID string) ([]*domain.Session, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"client_id": clientID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sessions []*domain.Session
	if err := cursor.All(ctx, &sessions); err != nil {
		return nil, err
	}

	return sessions, nil
}

func (r *MongoSessionRepository) GetAll(ctx context.Context) ([]*domain.Session, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var sessions []*domain.Session
	if err := cursor.All(ctx, &sessions); err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *MongoSessionRepository) Update(ctx context.Context, session *domain.Session) error {
	session.UpdatedAt = time.Now().UTC()
	filter := bson.M{"_id": session.SessionID}
	update := bson.M{"$set": session}
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("session not found")
	}
	return nil
}

func (r *MongoSessionRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {

		result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
		if err != nil {
			return err
		}
		if result.DeletedCount == 0 {
			return errors.New("session not found")
		}
		return nil
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("session not found")
	}

	return nil
}

func (r *MongoSessionRepository) CheckTimeOverlap(ctx context.Context, sessions []*domain.Session, startTime, endTime string, day string) (bool, error) {

	layout := "15:04"
	newStart, err := time.Parse(layout, startTime)
	if err != nil {
		return false, err
	}
	newEnd, err := time.Parse(layout, endTime)
	if err != nil {
		return false, err
	}

	for _, session := range sessions {

		if !session.IsActive() {
			continue
		}

		existingStart := session.StartTime
		existingEnd := session.EndTime

		if newStart.Before(existingEnd) && newEnd.After(existingStart) {
			return true, nil
		}
	}

	return false, nil
}
