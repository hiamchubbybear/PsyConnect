package repository

import (
	"chatservice/bootstrap"
	"chatservice/internal/db"
	"log"
)

type RepositoryManager struct {
	chatRepo *ChatRepository
}

func NewRepositoryManager(env *bootstrap.Env) *RepositoryManager {
	messageRepo := NewChatRepository(db.GetChatCollection())
	if messageRepo == nil {
		log.Fatal("clientRepo is nil")
	}
	repoManager := RepositoryManager{
		chatRepo: messageRepo,
	}
	return &repoManager
}
