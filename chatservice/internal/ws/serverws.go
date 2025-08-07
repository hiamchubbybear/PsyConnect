package ws

import (
	encoder "chatservice/internal/utils/conversation"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	receiver := query.Get("receiver")
	userID := query.Get("user_id")
	profileID := query.Get("profile_id")

	// userID := r.Header.Get("X-User-Id")
	// if userID == "" {
	// 	http.Error(w, "Unauthenticated", http.StatusUnauthorized)
	// 	return
	// }

	// profileID := r.Header.Get("X-Profile-Id")
	// if profileID == "" {
	// 	http.Error(w, "Unauthenticated", http.StatusUnauthorized)
	// 	return
	// }

	conversationID, err := encoder.New().EncodeConversationId(userID, receiver)
	if err != nil || conversationID == "" {
		http.Error(w, "Cannot generate conversation ID", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := &Client{
		ID:             userID,
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
