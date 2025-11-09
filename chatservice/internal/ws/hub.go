package ws

import (
	"encoding/json"
	"fmt"
	"time"

	"chatservice/internal/model"
	"chatservice/internal/repository"
)

type Message struct {
	ConversationID string `json:"conversationId" bson:"conversation_id"`
	SenderID       string `json:"sendId"`
	Content        []byte `json:"text" `
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
			h.handleMessageV1(message)
		}
	}
}

// Deprecated : Replace handleMessageV1 insteads
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

func (h *Hub) handleMessageV1(message Message) {
	fmt.Println("Handler save messsage")
	newChat := &model.Chat{
		SenderID:       message.SenderID,
		ConversationID: message.ConversationID,
		Text:           string(message.Content),
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	res, err := h.ChatRepo.CreateChat(newChat)
	if err != nil || res == nil {
		fmt.Println("Save error:", err)
		return
	}

	fmt.Printf("[Chat Saved] sender=%s text=%s\n", newChat.SenderID, newChat.Text)

	for _, client := range h.Clients[message.ConversationID] {

		payload := map[string]interface{}{
			"id":             newChat.ID,
			"userId":         newChat.SenderID,
			"senderId":       newChat.SenderID,
			"userName":       "",
			"userAvatar":     "",
			"content":        newChat.Text,
			"timestamp":      newChat.CreatedAt.Format(time.RFC3339),
			"isMine":         nil,
			"conversationId": newChat.ConversationID,
		}

		data, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("Error encoding message:", err)
			continue
		}

		client.Send <- data
	}
}
