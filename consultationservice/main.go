package main

import (
	"consultationservice/bootstrap"
	"consultationservice/db"
	handlers "consultationservice/handler"
	"consultationservice/route"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var mongoClient *mongo.Client

func main() {
	db.InitDB()
	env := bootstrap.LoadEnv()
	handlers.InitTherapistHandler()
	//grcpHandler.InitGrpcClient()
	routes := route.RouterInit()
	urI := fmt.Sprintf("%v:%v", env.Addr, env.Port)
	routes.Run(urI)

}
