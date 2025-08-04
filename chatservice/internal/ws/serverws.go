package ws

import (
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
	userID := r.Header.Get("X-User-Id")
	if userID == "" {
		http.Error(w, "Unauthenticated ", http.StatusUnauthorized)
		return
	}
	profileId := r.Header.Get("X-Profile-Id")
	if profileId == "" {
		http.Error(w, "Unauthenticated ", http.StatusUnauthorized)
		return
	}
	conversationID := query.Get("conversation_id")
	if conversationID == "" {
		http.Error(w, "Invalid conversation id  ", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		hub:            hub,
		conn:           conn,
		send:           make(chan []byte, 256),
		userID:         userID,
		conversationID: conversationID,
		profileId:      profileId,
	}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
