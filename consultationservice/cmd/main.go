package main

import (
	"consultationservice/bootstrap"
	"consultationservice/internal/db"
	handlers "consultationservice/internal/handler"
	"consultationservice/internal/kafka"
	"consultationservice/internal/repository"
	"consultationservice/internal/route"
)

func main() {
	db.InitDB()
	env := bootstrap.LoadEnv()
	repomanager := repository.NewRepositoryManager(env)

	sessionHandler := handlers.NewSessionHandler(env, repomanager)
	therapistHandler := handlers.NewTherapistHandler(env, repomanager)
	clientHandler := handlers.NewClientHandler(env, repomanager)
	matchHandler := handlers.NewMatchHandler(env, repomanager)

	kafka.NewConsumer(env)
	kafka.NewProducer(env)

	route.RouterInit(env, clientHandler, therapistHandler, matchHandler, sessionHandler)
}
