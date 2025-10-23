package main

import (
	"log"
	"net/http"

	"chatservice/bootstrap"
	"chatservice/internal/db"
	"chatservice/internal/repository"
	"chatservice/internal/ws"
	"chatservice/internal/ws/middleware"
)

func main() {
	db.InitDB()
	env := bootstrap.LoadEnv()
	repoManager := repository.NewRepositoryManager(env)

	hub := ws.NewHub(repoManager.MessageRepo)
	go hub.Run()

	http.HandleFunc("/ws", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWS(hub, w, r)
	}))

	log.Println(" WebSocket server running at :8085")
	log.Fatal(http.ListenAndServe(":8085", nil))
}
