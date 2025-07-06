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
	MatchingRepo  *MatchRepository
	SessionRepo   *SessionRepository
	GrpcProfile   *handler.ProfileGrpc
	Kafka         *kafka.Producer
}

func NewRepositoryManager(env *bootstrap.Env) *RepositoryManager {
	log.Printf("NewRepositoryManager called with env: %+v", env)

	grpcProfile, err := handler.NewProfileGrpc(env.GrpcAdd)
	if err != nil {
		log.Fatal("Failed to create grpc: ", err)
	}

	kafkaProducer, err := kafka.NewProducer(env)
	if err != nil {
		log.Fatal("Failed to initialize kafka: ", err)
	}

	log.Printf("Initializing ClientRepository...")
	clientRepo := NewClientRepository(db.GetClientCollection())
	if clientRepo == nil {
		log.Fatal("clientRepo is nil")
	}
	log.Printf("ClientRepo created: %+v", clientRepo)

	log.Printf("Initializing TherapistRepository...")
	therapistRepo := NewTherapistRepository(db.GetTherapistCollection())
	if therapistRepo == nil {
		log.Fatal("therapistRepo is nil")
	}
	log.Printf("TherapistRepo created: %+v", therapistRepo)

	log.Printf("Initializing MatchRepository...")
	matchCollection := db.GetMatchedCollection()
	log.Printf("Got matchCollection: %+v", matchCollection)
	matchingRepo := NewMatchRepository(matchCollection)
	if matchingRepo == nil {
		log.Fatal("matchingRepo is nil after NewMatchRepository")
	}
	log.Printf("MatchingRepo created successfully: %+v", matchingRepo)

	log.Printf("Initializing SessionRepository...")
	sessionRepo := NewSessionRepository(db.GetSessionCollection(), clientRepo, therapistRepo, matchingRepo)
	if sessionRepo == nil {
		log.Fatal("sessionRepo is nil")
	}
	log.Printf("SessionRepo created: %+v", sessionRepo)

	log.Printf("Initializing SwipeRepository...")
	swipesRepo := NewSwipeRepository(clientRepo, therapistRepo, sessionRepo, db.GetSwipedCollection())
	if swipesRepo == nil {
		log.Fatalf("swipesRepo is nil")
	}
	log.Printf("SwipesRepo created: %+v", swipesRepo)

	repoManager := &RepositoryManager{
		ClientRepo:    clientRepo,
		TherapistRepo: therapistRepo,
		MatchingRepo:  matchingRepo,
		SessionRepo:   sessionRepo,
		GrpcProfile:   grpcProfile,
		Kafka:         kafkaProducer,
		SwipeRepo:     swipesRepo,
	}

	log.Printf("RepositoryManager created: %+v", repoManager)
	log.Printf("RepositoryManager.MatchingRepo: %+v", repoManager.MatchingRepo)

	return repoManager
}
