package main

import (
	"chatservice/internal/ws"
	"log"
	"net/http"
)

func main() {
	hub := ws.NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWS(hub, w, r)
	})

	log.Println("Server started at :8085")
	log.Fatal(http.ListenAndServe(":8085", nil))
}
