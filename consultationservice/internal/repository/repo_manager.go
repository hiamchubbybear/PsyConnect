package repository

import (
	"consultationservice/bootstrap"
	"consultationservice/internal/db"
	"consultationservice/internal/grpc/handler"
	"consultationservice/internal/kafka"
	"log"
)

type RepositoryManager struct {
	SwipeRepo     *SwipeRepository
	ClientRepo    *ClientRepository
	TherapistRepo *TherapistRepository
	MatchingRepo  *MatchingRepository
	SessionRepo   *SessionRepository
	GrpcProfile   *handler.ProfileGrpc
	Kafka         *kafka.Producer
}

func NewRepositoryManager(env *bootstrap.Env) *RepositoryManager {
	grpcProfile, err := handler.NewProfileGrpc(env.GrpcAdd)
	if err != nil {
		log.Fatal("Failed to create grpc: ", err)
	}

	kafkaProducer, err := kafka.NewProducer(env)
	if err != nil {
		log.Fatal("Failed to initialize kafka: ", err)
	}

	clientRepo := NewClientRepository(db.GetClientCollection())
	if clientRepo == nil {
		log.Fatal("clientRepo is nil")
	}

	therapistRepo := NewTherapistRepository(db.GetTherapistCollection())
	if therapistRepo == nil {
		log.Fatal("therapistRepo is nil")
	}

	sessionRepo := NewSessionRepository(db.GetSessionCollection(), clientRepo, therapistRepo, nil)
	if sessionRepo == nil {
		log.Fatal("sessionRepo is nil")
	}
	swipesRepo := NewSwipeRepository(clientRepo, therapistRepo, sessionRepo, db.GetSwipedCollection())
	if swipesRepo == nil {
		log.Fatalf("swipesRepo is nil")
	}
	matchingRepo := NewMatchingRepository(clientRepo, grpcProfile, therapistRepo, sessionRepo)
	if matchingRepo == nil {
		log.Fatal("matchingRepo is nil")
	}

	sessionRepo.SetMatchingRepository(matchingRepo)

	return &RepositoryManager{
		ClientRepo:    clientRepo,
		TherapistRepo: therapistRepo,
		MatchingRepo:  matchingRepo,
		SessionRepo:   sessionRepo,
		GrpcProfile:   grpcProfile,
		Kafka:         kafkaProducer,
		SwipeRepo:     swipesRepo,
	}
}
