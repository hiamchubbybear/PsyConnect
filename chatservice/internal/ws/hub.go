package ws

import (
	"log"
	"time"
)

type Hub struct {
	clients map[*Client]bool

	rooms map[string]map[*Client]bool

	register   chan *Client
	unregister chan *Client

	broadcast chan Message
}
	
type Message struct {
	ConversationID string    `json:"conversation_id"`
	UserID         string    `json:"user_id"`
	Text           string    `json:"text"`
	Timestamp      time.Time `json:"timestamp"`
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		rooms:      make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			if _, ok := h.rooms[client.conversationID]; !ok {
				h.rooms[client.conversationID] = make(map[*Client]bool)
			}
			h.rooms[client.conversationID][client] = true
			log.Printf("Client %s joined room %s", client.userID, client.conversationID)

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				if room, ok := h.rooms[client.conversationID]; ok {
					delete(room, client)
					if len(room) == 0 {
						delete(h.rooms, client.conversationID)
					}
				}
			}

		case msg := <-h.broadcast:
			if room, ok := h.rooms[msg.ConversationID]; ok {
				for client := range room {
					select {
					case client.send <- []byte(msg.Text):
					default:
						close(client.send)
						delete(room, client)
					}
				}
			}
		}
	}
}
