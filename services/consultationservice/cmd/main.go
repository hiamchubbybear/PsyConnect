package main

import (
	"log"
	"os"
	"strings"

	"consultationservice/bootstrap"
	"consultationservice/internal/db"
	"consultationservice/internal/kafka"
	"consultationservice/internal/redis"
	"consultationservice/internal/repository"
	"consultationservice/internal/route"
	"consultationservice/pkg/logger"

	// DDD imports
	consultationRepo "consultationservice/internal/consultation/repository"
	httpHandler "consultationservice/internal/consultation/transport/http"
	"consultationservice/internal/consultation/usecase"

	// DDD imports
	clientRepo "consultationservice/internal/client/repository"
	clientHTTP "consultationservice/internal/client/transport/http"
	clientUseCase "consultationservice/internal/client/usecase"

	therapistRepo "consultationservice/internal/therapist/repository"
	therapistHTTP "consultationservice/internal/therapist/transport/http"
	therapistUseCase "consultationservice/internal/therapist/usecase"

	matchRepo "consultationservice/internal/matching/repository"
	matchHTTP "consultationservice/internal/matching/transport/http"
	matchUseCase "consultationservice/internal/matching/usecase"

	swipeRepo "consultationservice/internal/swipe/repository"
	swipeHTTP "consultationservice/internal/swipe/transport/http"
	swipeUseCase "consultationservice/internal/swipe/usecase"

	reactionRepo "consultationservice/internal/newsfeed/reaction/repository"
	reactionHTTP "consultationservice/internal/newsfeed/reaction/transport/http"
	reactionUseCase "consultationservice/internal/newsfeed/reaction/usecase"

	commentRepo "consultationservice/internal/newsfeed/comment/repository"
	commentHTTP "consultationservice/internal/newsfeed/comment/transport/http"
	commentUseCase "consultationservice/internal/newsfeed/comment/usecase"

	socialRepo "consultationservice/internal/newsfeed/social/repository"
	socialHTTP "consultationservice/internal/newsfeed/social/transport/http"
	socialUseCase "consultationservice/internal/newsfeed/social/usecase"

	postRepo "consultationservice/internal/newsfeed/post/repository"
	postHTTP "consultationservice/internal/newsfeed/post/transport/http"
	postUseCase "consultationservice/internal/newsfeed/post/usecase"

	groupRepo "consultationservice/internal/newsfeed/group/repository"
	groupHTTP "consultationservice/internal/newsfeed/group/transport/http"
	groupUseCase "consultationservice/internal/newsfeed/group/usecase"
)

func main() {
	env := bootstrap.LoadEnv()
	db.InitDB()

	kafkaBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	if len(kafkaBrokers) == 0 || kafkaBrokers[0] == "" {
		kafkaBrokers = []string{"localhost:9092"}
	}

	kafkaLogger := logger.NewKafkaLogger(logger.Config{
		Brokers:     kafkaBrokers,
		Topic:       "logging-service",
		ServiceName: "consultation-service",
		Environment: os.Getenv("ENVIRONMENT"),
		Version:     "1.0.0",
	})
	defer kafkaLogger.Close()

	kafkaLogger.Info("Consultation service starting", map[string]interface{}{
		"port":        os.Getenv("PORT"),
		"environment": os.Getenv("ENVIRONMENT"),
	})

	redisClient, err := redis.NewRedisStore(env)
	if err != nil {
		kafkaLogger.Warn("Failed to initialize Redis (continuing without cache)", map[string]interface{}{
			"error": err.Error(),
		})
		log.Printf("Warning: Redis initialization failed: %v", err)
	}
	// ===== Initialize Legacy Repositories (for adapters) =====
	// These will be replaced with pure DDD implementations
	legacyPostRepo := repository.NewPostRepo()
	legacyReactionRepo := repository.NewReactionRepo()
	legacyBookmarkRepo := repository.NewBookmarkRepo()

	// Initialize Kafka Producer
	kafkaProducer, err := kafka.NewProducer(env)
	if err != nil {
		log.Printf("⚠️  Warning: Failed to initialize Kafka producer: %v", err)
	}

	// ===== DDD Session Components =====
	// Initialize DDD Session Repository
	sessionRepo := consultationRepo.NewMongoSessionRepository(db.GetSessionCollection())

	// Initialize Session UseCases
	createSessionUC := usecase.NewCreateSessionUseCase(sessionRepo)
	getSessionUC := usecase.NewGetSessionUseCase(sessionRepo)
	deleteSessionUC := usecase.NewDeleteSessionUseCase(sessionRepo)
	startCallUC := usecase.NewStartCallUseCase(sessionRepo, kafkaProducer)

	// Initialize DDD HTTP Handler
	sessionHandler := httpHandler.NewHandler(
		createSessionUC,
		getSessionUC,
		deleteSessionUC,
		startCallUC,
	)
	// ===== End DDD Session Components =====

	// ===== DDD Client Components =====
	// Initialize DDD Client Repository
	clientRepository := clientRepo.NewMongoClientRepository(db.GetClientCollection(), redisClient)

	// Initialize Client UseCases
	createClientUC := clientUseCase.NewCreateClientUseCase(clientRepository)
	getClientUC := clientUseCase.NewGetClientUseCase(clientRepository)
	updateClientUC := clientUseCase.NewUpdateClientUseCase(clientRepository)
	deleteClientUC := clientUseCase.NewDeleteClientUseCase(clientRepository)

	// Initialize DDD Client HTTP Handler
	clientHandler := clientHTTP.NewHandler(
		createClientUC,
		getClientUC,
		updateClientUC,
		deleteClientUC,
	)
	// ===== End DDD Client Components =====

	// ===== DDD Therapist Components =====
	// Initialize DDD Therapist Repository
	therapistRepository := therapistRepo.NewMongoTherapistRepository(db.GetTherapistCollection(), redisClient)

	// Initialize Therapist UseCases
	createTherapistUC := therapistUseCase.NewCreateTherapistUseCase(therapistRepository)
	getTherapistUC := therapistUseCase.NewGetTherapistUseCase(therapistRepository)
	updateTherapistUC := therapistUseCase.NewUpdateTherapistUseCase(therapistRepository)
	deleteTherapistUC := therapistUseCase.NewDeleteTherapistUseCase(therapistRepository)

	// Initialize DDD Therapist HTTP Handler
	therapistHandler := therapistHTTP.NewHandler(
		createTherapistUC,
		getTherapistUC,
		updateTherapistUC,
		deleteTherapistUC,
	)
	// ===== End DDD Therapist Components =====

	// ===== DDD Matching Components =====
	matchRepository := matchRepo.NewMongoMatchRepository(db.GetMatchCollection(), redisClient)
	createMatchUC := matchUseCase.NewCreateMatchUseCase(matchRepository, clientRepository, therapistRepository)
	getClientMatchesUC := matchUseCase.NewGetClientMatchesUseCase(matchRepository)

	matchHandler := matchHTTP.NewHandler(createMatchUC, getClientMatchesUC)
	// ===== End DDD Matching Components =====

	// ===== DDD Swipe Components =====
	swipeRepository := swipeRepo.NewMongoSwipeRepository(db.GetSwipedCollection(), matchRepository, redisClient)
	insertSwipeUC := swipeUseCase.NewInsertSwipeUseCase(swipeRepository)
	swipeAndMatchUC := swipeUseCase.NewSwipeAndMatchUseCase(swipeRepository)
	recommendUC := swipeUseCase.NewRecommendUseCase(clientRepository, therapistRepository, swipeRepository)

	swipeHandler := swipeHTTP.NewHandler(insertSwipeUC, swipeAndMatchUC, recommendUC)
	// ===== End DDD Swipe Components =====

	// ===== DDD Reaction Components =====
	reactionRepository := reactionRepo.NewMongoReactionRepository(db.GetReactionCollection(), redisClient)

	// Create adapter for legacy PostRepository
	postRepoAdapter := reactionUseCase.NewPostRepositoryAdapter(legacyPostRepo)

	toggleReactionUC := reactionUseCase.NewToggleReactionUseCase(reactionRepository, postRepoAdapter)
	getReactionsUC := reactionUseCase.NewGetReactionsUseCase(reactionRepository)

	reactionHandler := reactionHTTP.NewHandler(toggleReactionUC, getReactionsUC, kafkaProducer)
	// ===== End DDD Reaction Components =====

	// ===== DDD Comment Components =====
	commentRepository := commentRepo.NewMongoCommentRepository(db.GetCommentCollection(), redisClient)
	postRepoAdapterComment := commentUseCase.NewPostRepositoryAdapter(legacyPostRepo)

	createCommentUC := commentUseCase.NewCreateCommentUseCase(commentRepository, postRepoAdapterComment)
	getCommentsUC := commentUseCase.NewGetCommentsUseCase(commentRepository)
	updateCommentUC := commentUseCase.NewUpdateCommentUseCase(commentRepository)
	deleteCommentUC := commentUseCase.NewDeleteCommentUseCase(commentRepository, postRepoAdapterComment)
	getRepliesUC := commentUseCase.NewGetRepliesUseCase(commentRepository)

	commentHandler := commentHTTP.NewHandler(createCommentUC, getCommentsUC, updateCommentUC, deleteCommentUC, getRepliesUC, kafkaProducer)
	// ===== End DDD Comment Components =====

	// ===== DDD Social Components =====
	followRepository := socialRepo.NewMongoFollowRepository(db.GetFollowCollection())
	bookmarkRepository := socialRepo.NewMongoBookmarkRepository(db.GetBookmarkCollection())

	followUserUC := socialUseCase.NewFollowUserUseCase(followRepository)
	unfollowUserUC := socialUseCase.NewUnfollowUserUseCase(followRepository)
	bookmarkPostUC := socialUseCase.NewBookmarkPostUseCase(bookmarkRepository)
	unbookmarkPostUC := socialUseCase.NewUnbookmarkPostUseCase(bookmarkRepository)

	socialHandler := socialHTTP.NewHandler(followUserUC, unfollowUserUC, bookmarkPostUC, unbookmarkPostUC)
	// ===== End DDD Social Components =====

	// ===== DDD Post Components =====
	legacyPostRepoAdapter := postRepo.NewLegacyPostRepositoryAdapter(legacyPostRepo)

	createPostUC := postUseCase.NewCreatePostUseCase(legacyPostRepoAdapter)
	getPostByIDUC := postUseCase.NewGetPostByIDUseCase(legacyPostRepoAdapter)
	updatePostUC := postUseCase.NewUpdatePostUseCase(legacyPostRepoAdapter)
	deletePostUC := postUseCase.NewDeletePostUseCase(legacyPostRepoAdapter)
	getFeedUC := postUseCase.NewGetFeedUseCase(legacyPostRepoAdapter)
	getTrendingUC := postUseCase.NewGetTrendingPostsUseCase(legacyPostRepoAdapter)
	searchPostsUC := postUseCase.NewSearchPostsUseCase(legacyPostRepoAdapter)
	getUserPostsUC := postUseCase.NewGetUserPostsUseCase(legacyPostRepoAdapter)
	getByTagUC := postUseCase.NewGetPostsByTagUseCase(legacyPostRepoAdapter)
	getByCategoryUC := postUseCase.NewGetPostsByCategoryUseCase(legacyPostRepoAdapter)
	incrementViewUC := postUseCase.NewIncrementViewCountUseCase(legacyPostRepoAdapter)
	getPopularTagsUC := postUseCase.NewGetPopularTagsUseCase(legacyPostRepoAdapter)

	postHandler := postHTTP.NewHandler(
		createPostUC, getPostByIDUC, updatePostUC, deletePostUC,
		getFeedUC, getTrendingUC, searchPostsUC, getUserPostsUC,
		getByTagUC, getByCategoryUC, incrementViewUC, getPopularTagsUC,
		redisClient, legacyReactionRepo, legacyBookmarkRepo,
	)
	// ===== End DDD Post Components =====

	// ===== DDD Group Components =====
	groupRepository := groupRepo.NewGroupRepo()

	createGroupUC := groupUseCase.NewCreateGroupUseCase(groupRepository)
	getGroupsUC := groupUseCase.NewGetGroupsUseCase(groupRepository)
	getGroupByIDUC := groupUseCase.NewGetGroupByIDUseCase(groupRepository)
	joinGroupUC := groupUseCase.NewJoinGroupUseCase(groupRepository)

	groupHandler := groupHTTP.NewHandler(
		createGroupUC,
		getGroupsUC,
		getGroupByIDUC,
		joinGroupUC,
	)
	// ===== End DDD Group Components =====

	// therapistHandler := handlers.NewTherapistHandler(env, repomanager) // Legacy removed
	// clientHandler now uses DDD version above
	// matchHandler used to be here
	// swipeHandler used to be here

	// ===== Legacy Handlers (to be removed) =====
	// postHandler used to be here (now DDD)
	// reactionHandler used to be here (now DDD)
	// commentHandler used to be here (now DDD)
	// socialHandler used to be here (now DDD)

	kafka.NewConsumer(env)
	kafka.NewProducer(env)

	kafkaLogger.Info("Consultation service initialized successfully", nil)

	route.RouterInit(
		env,
		kafkaLogger,
		clientHandler,
		therapistHandler,
		matchHandler,
		sessionHandler,
		swipeHandler,
		postHandler,
		reactionHandler,
		commentHandler,
		socialHandler,
		groupHandler,
	)
}
