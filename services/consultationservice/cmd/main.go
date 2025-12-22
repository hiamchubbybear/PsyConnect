package main

import (
	"consultationservice/bootstrap"
	"consultationservice/internal/db"
	handlers "consultationservice/internal/handler"
	"consultationservice/internal/kafka"
	"consultationservice/internal/redis"
	"consultationservice/internal/repository"
	"consultationservice/internal/route"
	"consultationservice/pkg/logger"
	"log"
	"os"
	"strings"
)

func main() {
	env := bootstrap.LoadEnv()
	db.InitDB()

	// Initialize Kafka Logger
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
		log.Printf("⚠️ Warning: Redis initialization failed: %v", err)
	}
	repomanager := repository.NewRepositoryManager(env, redisClient)

	// Existing handlers
	sessionHandler := handlers.NewSessionHandler(env, repomanager)
	therapistHandler := handlers.NewTherapistHandler(env, repomanager)
	clientHandler := handlers.NewClientHandler(env, repomanager)
	matchHandler := handlers.NewMatchHandler(env, repomanager)
	swipeHandler := handlers.NewSwipeHandler(env, repomanager)

	// Newsfeed handlers
	postHandler := handlers.NewPostHandler(env, repomanager, redisClient)
	reactionHandler := handlers.NewReactionHandler(env, repomanager, redisClient)
	commentHandler := handlers.NewCommentHandler(env, repomanager, redisClient)
	socialHandler := handlers.NewSocialHandler(env, repomanager, redisClient)

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
	)
}
