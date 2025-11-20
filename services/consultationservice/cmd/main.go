package main

import (
	"consultationservice/bootstrap"
	"consultationservice/internal/db"
	handlers "consultationservice/internal/handler"
	"consultationservice/internal/kafka"
	"consultationservice/internal/redis"
	"consultationservice/internal/repository"
	"consultationservice/internal/route"
	"log"
)

func main() {
	db.InitDB()
	env := bootstrap.LoadEnv()
	redisClient, err := redis.NewRedisStore(env)
	if err != nil {
		log.Println(err)
		panic(err)
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

	route.RouterInit(
		env,
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
