package repository

import (
	"chatservice/bootstrap"
	"chatservice/internal/db"
	"log"
)

type RepositoryManager struct {
	MessageRepo *ChatRepository
}

func NewRepositoryManager(env *bootstrap.Env) *RepositoryManager {
	messageRepo := NewChatRepository(db.GetChatCollection())
	if messageRepo == nil {
		log.Fatal("clientRepo is nil")
	}
	repoManager := RepositoryManager{
		MessageRepo: messageRepo,
	}
	return &repoManager
}
