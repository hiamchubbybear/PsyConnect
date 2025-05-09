package main

import (
	"consultationservice/bootstrap"
	"consultationservice/internal/db"
	"consultationservice/internal/handler"
	"consultationservice/internal/route"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var mongoClient *mongo.Client

func main() {
	db.InitDB()
	env := bootstrap.LoadEnv()
	handlers.InitTherapistHandler()
	handlers.InitClientHandler()
	handlers.InitMatchHandler()
	route.RouterInit(env)
}
