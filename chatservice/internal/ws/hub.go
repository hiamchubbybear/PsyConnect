package ws

import (
	"chatservice/internal/model"
	"chatservice/internal/repository"
	"fmt"
	"time"
)

type Message struct {
	ConversationID string
	SenderID       string
	Content        []byte
}

type Hub struct {
	Clients    map[string][]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
	ChatRepo   *repository.ChatRepository
}

func NewHub(repo *repository.ChatRepository) *Hub {
	return &Hub{
		Clients:    make(map[string][]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message),
		ChatRepo:   repo,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.ConversationID] = append(h.Clients[client.ConversationID], client)

		case client := <-h.Unregister:
			clients := h.Clients[client.ConversationID]
			for i, c := range clients {
				if c == client {
					h.Clients[client.ConversationID] = append(clients[:i], clients[i+1:]...)
					break
				}
			}

		case message := <-h.Broadcast:
			h.handleMessage(message)
		}
	}
}
func (h *Hub) handleMessage(message Message) {
	var newChat = &model.Chat{
		SenderID:       message.SenderID,
		ConversationID: message.ConversationID,
		Text:           string(message.Content),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	res, err := h.ChatRepo.CreateChat(newChat)
	if err != nil || res == nil {
		fmt.Println("Save error:", err)
	}
	fmt.Println(newChat)
	for _, client := range h.Clients[message.ConversationID] {
		if client.ID != message.SenderID {
			client.Send <- message.Content
		}
	}
}
