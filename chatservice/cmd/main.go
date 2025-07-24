package main

import (
	"chatservice/bootstrap"
	"chatservice/internal/db"
	"chatservice/internal/repository"
	"chatservice/internal/route"
)

func main() {
	
	db.InitDB()
	env := bootstrap.LoadEnv()
	repomanager := repository.NewRepositoryManager(env)

	route.RouterInit(env, repomanager)
}
