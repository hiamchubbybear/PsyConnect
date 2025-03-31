package main

import (
	"consultationservice/bootstrap"
	"consultationservice/db"
	handlers "consultationservice/handler"
	"consultationservice/route"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var mongoClient *mongo.Client

func main() {
	db.InitDB()
	env := bootstrap.LoadEnv()
	handlers.InitTherapistHandler()
	route.RouterInit(env)

}
