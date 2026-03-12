package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"chatservice/internal/ws/middleware"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	receiver := r.URL.Query().Get("receiver")
	conversationID := r.URL.Query().Get("conversationId")
	profileID := r.Context().Value(middleware.ProfileIDKey).(string)
	log.Printf("[ServeWS] New connection attempt: profileID=%s, receiver=%s, conversationID=%s", profileID, receiver, conversationID)
	if receiver == "" || conversationID == "" {
		http.Error(w, "Cannot generate conversation ID", http.StatusBadRequest)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := &Client{
		ID:             profileID,
		ProfileID:      profileID,
		ConversationID: conversationID,
		Conn:           conn,
		Send:           make(chan []byte, 256),
		Hub:            hub,
	}

	hub.Register <- client
	go client.WritePump()
	go client.ReadPump()
}
