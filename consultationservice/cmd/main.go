package main

import (
	"consultationservice/bootstrap"
	"consultationservice/internal/db"
	"consultationservice/internal/handler"
	"consultationservice/internal/kafka"
	"consultationservice/internal/route"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var mongoClient *mongo.Client

func main() {
	db.InitDB()
	env := bootstrap.LoadEnv()
	handlers.InitTherapistHandler(env)
	handlers.InitClientHandler(env)
	handlers.InitMatchHandler(env)
	kafka.NewConsumer(env)
	kafka.NewProducer(env)
	route.RouterInit(env)
}
