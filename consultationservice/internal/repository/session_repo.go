package repository

import (
	"consultationservice/internal/dto"
	"context"
	"errors"
	"log"

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
	update := bson.D{{"$push", bson.D{{"current_session", sessionID}}}}

	_, _ = r.clientRepo.MongoDBCollection.UpdateOne(
		context.Background(),
		bson.D{{"profile_id", clientId}},
		update,
	)

	_, _ = r.therapistRepo.MongoDBCollection.UpdateOne(
		context.Background(),
		bson.D{{"profile_id", therapistId}},
		update,
	)

	return res, nil
}
func (r *SessionRepository) DeleteCurrentSession(session dto.DeleteSessionRequest) (bool, error) {
	if session.SessionID == "" {
		return true, errors.New("Session id couldn't be empty")
	}
	return false, nil
}
func (s *SessionRepository) SetMatchingRepository(matchingRepo *MatchingRepository) {
	s.matchingRepo = matchingRepo
}
