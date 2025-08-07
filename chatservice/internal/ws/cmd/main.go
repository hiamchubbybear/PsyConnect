package main

import (
	"chatservice/bootstrap"
	"chatservice/internal/db"
	"chatservice/internal/repository"
	"chatservice/internal/ws"
	"log"
	"net/http"
)

func main() {
	db.InitDB()
	env := bootstrap.LoadEnv()
	repomanager := repository.NewRepositoryManager(env)
	hub := ws.NewHub(repomanager.MessageRepo)
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWS(hub, w, r)
	})

	log.Println("Server started at :8085")
	log.Fatal(http.ListenAndServe(":8085", nil))
}
