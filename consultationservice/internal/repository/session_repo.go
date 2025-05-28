package repository

import (
	"consultationservice/internal/dto"
	"consultationservice/internal/model"
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SessionRepository struct {
	MongoDBCollection *mongo.Collection
	clientRepo        *ClientRepository
	therapistRepo     *TherapistRepository
	matchingRepo      *MatchingRepository
}

func NewSessionRepository(collection *mongo.Collection, clientRepo *ClientRepository, therapistRepo *TherapistRepository, matchingRepo *MatchingRepository) *SessionRepository {
	return &SessionRepository{
		clientRepo:        clientRepo,
		therapistRepo:     therapistRepo,
		MongoDBCollection: collection,
		matchingRepo:      matchingRepo,
	}
}
func (r *SessionRepository) CreateNewSession(session dto.SessionRequest) (interface{}, error) {
	var (
		clientId    = session.ClientID
		therapistId = session.TherapistID
	)

	_, terr := r.therapistRepo.FindTherapistMatchingProfile(therapistId)
	if terr != nil {
		log.Printf("Failed to find therapist profile : %v", terr)
		return nil, terr
	}
	log.Printf("Find client have client id %v", clientId)
	client, cerr := r.clientRepo.FindClientMatchingProfile("c001")
	if cerr != nil {
		log.Printf("Failed to find client profile: %v", cerr)
		return nil, cerr
	}
	if client == nil || client.ProfileId == "" {
		log.Println("Client not found or invalid")
		return nil, errors.New("cannot find client")
	}
	validTime, err := r.matchingRepo.CheckValidSessionTimeAndDays(session.SessionTime, clientId, therapistId)
	if err != nil {
		return nil, err
	}
	if !validTime {
		return nil, errors.New("Invalid time")
	}
	log.Printf("Session 1 ")
	res, err := r.MongoDBCollection.InsertOne(context.Background(), session)
	if err != nil {
		return nil, err
	}

	sessionID := res.InsertedID
	log.Printf("Session %v", sessionID)
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "current_session", Value: sessionID}}}}

	_, _ = r.clientRepo.MongoDBCollection.UpdateOne(
		context.Background(),
		bson.D{{Key: "profile_id", Value: clientId}},
		update,
	)

	_, _ = r.therapistRepo.MongoDBCollection.UpdateOne(
		context.Background(),
		bson.D{{Key: "profile_id", Value: therapistId}},
		update,
	)

	return res, nil
}
func (r *SessionRepository) DeleteCurrentSession(session dto.DeleteSessionRequest) (bool, error) {
	ctx := context.Background()

	if session.SessionID == "" {
		return false, errors.New("session id couldn't be empty")
	}
	objectID, err := primitive.ObjectIDFromHex(session.SessionID)
	if err != nil {
		return false, errors.New("invalid session id format")
	}

	res, err := r.MongoDBCollection.DeleteOne(ctx, bson.D{{Key: "_id", Value: objectID}})
	if err != nil {
		log.Printf("Failed to delete session: %v", err)
		return false, err
	}

	if res.DeletedCount == 0 {
		log.Println("Session not found")
		return false, errors.New("session not found")
	}

	therapistUpdate := bson.D{{Key: "$pull", Value: bson.D{{Key: "current_session", Value: objectID}}}}
	_, err = r.therapistRepo.MongoDBCollection.UpdateMany(ctx, bson.D{}, therapistUpdate)
	if err != nil {
		log.Printf("Failed to update therapist: %v", err)
		return false, err
	}

	clientUpdate := bson.D{{Key: "$pull", Value: bson.D{{Key: "current_session", Value: objectID}}}}
	_, err = r.clientRepo.MongoDBCollection.UpdateMany(ctx, bson.D{}, clientUpdate)
	if err != nil {
		log.Printf("Failed to update client: %v", err)
		return false, err
	}
	return true, nil
}
func (r *SessionRepository) GetAllSessions() ([]model.Session, error) {
	var sessions []model.Session
	cursor, err := r.MongoDBCollection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(context.Background(), &sessions); err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *SessionRepository) GetSessionByID(id string) (*model.Session, error) {
	var session model.Session
	err := r.MongoDBCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *SessionRepository) SetMatchingRepository(matchingRepo *MatchingRepository) {
	s.matchingRepo = matchingRepo
}
