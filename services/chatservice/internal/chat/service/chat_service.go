package service

import (
	"chatservice/internal/chat/model/model"
	"chatservice/internal/chat/repository/repository"
	"chatservice/internal/ws"
	"encoding/json"
	"log"
	"time"
)

type ChatService struct {
	chatRepo *repository.ChatRepository
}

func NewChatService(chatRepo *repository.ChatRepository) *ChatService {
	return &ChatService{
		chatRepo: chatRepo,
	}
}

func (s *ChatService) HandleChatMessage(hub *ws.Hub, message ws.Message) {
	log.Printf("Handling chat message from: %s", message.SenderID)

	// Unmarshal Data from json.RawMessage
	var data map[string]interface{}
	if err := json.Unmarshal(message.Data, &data); err != nil {
		log.Println("Invalid chat message: failed to parse data")
		return
	}

	text, ok := data["text"].(string)
	if !ok {
		log.Println("Invalid chat message: missing text")
		return
	}

	newChat := &model.Chat{
		SenderID:       message.SenderID,
		ConversationID: message.ConversationID,
		Text:           text,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	res, err := s.chatRepo.CreateChat(newChat)
	if err != nil || res == nil {
		log.Printf("Failed to save chat: %v", err)
		return
	}

	log.Printf("Chat saved: sender=%s text=%s", newChat.SenderID, newChat.Text)

	// Full payload for generic message handling
	payload := map[string]interface{}{
		"id":             newChat.ID,
		"senderId":       newChat.SenderID,
		"content":        newChat.Text,
		"timestamp":      newChat.CreatedAt.Format(time.RFC3339),
		"conversationId": newChat.ConversationID,
	}

	// Marshal payload to json.RawMessage
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal payload: %v", err)
		return
	}

	broadcastMsg := ws.Message{
		Type:           ws.MessageTypeChat,
		ConversationID: message.ConversationID,
		SenderID:       message.SenderID,
		Data:           payloadBytes,
	}

	hub.BroadcastToRoom(message.ConversationID, broadcastMsg)
}

func (s *ChatService) GetChatHistory(conversationID string, limit int, before time.Time) ([]*model.Chat, error) {
	return s.chatRepo.FindChatsByConversation(conversationID, limit, before)
}

func (s *ChatService) DeleteChat(chatID string) (int64, error) {
	return s.chatRepo.DeleteChatByID(chatID)
}
