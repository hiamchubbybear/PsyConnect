package repository

import (
	"consultationservice/internal/dto"
	grpc "consultationservice/internal/grpc/handler"
	"context"
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

func NewSessionRepository(collection *mongo.Collection, clientRepo *ClientRepository, grpcProfile *grpc.ProfileGrpc, therapistRepo *TherapistRepository, matchingRepo *MatchingRepository) *SessionRepository {
	return &SessionRepository{
		clientRepo:        clientRepo,
		therapistRepo:     therapistRepo,
		MongoDBCollection: clientRepo.MongoDBCollection,
		matchingRepo:      matchingRepo,
	}
}
func (r *SessionRepository) CreateNewSession(session dto.SessionRequest) (interface{}, error) {
	var (
		clientId    = session.ClientID
		therapistId = session.TherapistID
	)
	_, terr := r.therapistRepo.FindTherapistMatchingProfile(clientId)
	if terr != nil {
		log.Printf("Failed to find client profile")
		return nil, terr
	}
	_, cerr := r.clientRepo.FindClientMatchingProfile(therapistId)
	if cerr != nil {
		log.Printf("Failed to find therapist profile")
		return nil, cerr
	}
	validTime, err := r.matchingRepo.CheckValidSessionTimeAndDays(session.SessionTime, clientId, therapistId)
	if !validTime && err != nil {
		return nil, err
	}
	res, err := r.MongoDBCollection.InsertOne(context.Background(), session)
	if err != nil {
		return nil, err
	}
	sessionID := res.InsertedID
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
