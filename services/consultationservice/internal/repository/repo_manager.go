package repository

import (
	"consultationservice/bootstrap"
	"consultationservice/internal/db"
	"consultationservice/internal/grpc/handler"
	"consultationservice/internal/kafka"
	"consultationservice/internal/redis"
	"log"
)

type RepositoryManager struct {
	ClientRepo    *ClientRepository
	TherapistRepo *TherapistRepository
	MatchingRepo  *MatchRepository
	SessionRepo   *SessionRepository
	SwipeRepo     *SwipeRepository

	// Newsfeed repositories
	PostRepo     PostRepository
	ReactionRepo ReactionRepository
	CommentRepo  CommentRepository
	FollowRepo   FollowRepository
	BookmarkRepo BookmarkRepository

	GrpcProfile *handler.ProfileGrpc
	Kafka       *kafka.Producer
	Redis       *redis.RedisStore
}

func NewRepositoryManager(env *bootstrap.Env, redisClient redis.RedisStore) *RepositoryManager {
	log.Printf("NewRepositoryManager called with env: %+v", env)

	grpcAddr := env.GrpcAdd
	if grpcAddr == "" {
		grpcAddr = "profileservice:50051"
	}

	grpcProfile, err := handler.NewProfileGrpc(grpcAddr)
	if err != nil {
		log.Printf(" Warning: Failed to create grpc client: %v", err)
	}

	kafkaProducer, err := kafka.NewProducer(env)
	if err != nil {
		log.Printf(" Warning: Failed to initialize kafka: %v", err)
	}

	log.Printf("Initializing ClientRepository...")
	clientRepo := NewClientRepository(db.GetClientCollection(), redisClient)
	if clientRepo == nil {
		log.Printf(" Warning: clientRepo is nil")
	}
	log.Printf("ClientRepo created: %+v", clientRepo)

	log.Printf("Initializing TherapistRepository...")
	therapistRepo := NewTherapistRepository(db.GetTherapistCollection(), redisClient)
	if therapistRepo == nil {
		log.Printf(" Warning: therapistRepo is nil")
	}
	log.Printf("TherapistRepo created: %+v", therapistRepo)

	log.Printf("Initializing MatchRepository...")
	matchCollection := db.GetMatchedCollection()
	log.Printf("Got matchCollection: %+v", matchCollection)
	matchingRepo := NewMatchRepository(matchCollection, redisClient)
	if matchingRepo == nil {
		log.Printf(" Warning: matchingRepo is nil after NewMatchRepository")
	}
	log.Printf("MatchingRepo created successfully: %+v", matchingRepo)

	log.Printf("Initializing SessionRepository...")
	sessionRepo := NewSessionRepository(db.GetSessionCollection(), clientRepo, therapistRepo, matchingRepo, redisClient)
	if sessionRepo == nil {
		log.Printf(" Warning: sessionRepo is nil")
	}
	log.Printf("SessionRepo created: %+v", sessionRepo)

	log.Printf("Initializing SwipeRepository...")
	swipesRepo := NewSwipeRepository(clientRepo, therapistRepo, sessionRepo, db.GetSwipedCollection(), redisClient)
	if swipesRepo == nil {
		log.Printf(" Warning: swipesRepo is nil")
	}
	log.Printf("SwipesRepo created: %+v", swipesRepo)

	// Initialize newsfeed repositories
	log.Printf("Initializing Newsfeed repositories...")
	postRepo := NewPostRepo()
	reactionRepo := NewReactionRepo()
	commentRepo := NewCommentRepo()
	followRepo := NewFollowRepo()
	bookmarkRepo := NewBookmarkRepo()
	log.Printf("Newsfeed repos created successfully")

	repoManager := &RepositoryManager{
		ClientRepo:    clientRepo,
		TherapistRepo: therapistRepo,
		MatchingRepo:  matchingRepo,
		SessionRepo:   sessionRepo,
		GrpcProfile:   grpcProfile,
		Kafka:         kafkaProducer,
		SwipeRepo:     swipesRepo,
		Redis:         &redisClient,

		// Newsfeed repos
		PostRepo:     postRepo,
		ReactionRepo: reactionRepo,
		CommentRepo:  commentRepo,
		FollowRepo:   followRepo,
		BookmarkRepo: bookmarkRepo,
	}

	log.Printf("RepositoryManager created: %+v", repoManager)
	log.Printf("RepositoryManager.MatchingRepo: %+v", repoManager.MatchingRepo)

	return repoManager
}
