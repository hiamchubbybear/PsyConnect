package repository

import (
	"log"

	"chatservice/bootstrap"
	"chatservice/internal/db"
)

type RepositoryManager struct {
	MessageRepo            *ChatRepository
	ConversationRepository *ConversationRepository
}

func NewRepositoryManager(env *bootstrap.Env) *RepositoryManager {
	messageRepo := NewChatRepository(db.GetChatCollection())
	conversationRepo := NewConversationRepository(db.GetConversationCollection())

	if messageRepo == nil {
		log.Fatal("messageRepo is nil")
	}
	if conversationRepo == nil {
		log.Fatal("conversationRepo is nil")
	}

	return &RepositoryManager{
		MessageRepo:            messageRepo,
		ConversationRepository: conversationRepo,
	}
}
