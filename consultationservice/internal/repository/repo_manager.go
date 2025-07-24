package repository

import (
	"chatservice/bootstrap"
	"chatservice/internal/db"
	"chatservice/internal/repository"
	"log"
)

type RepositoryManager struct {
	chatRepo *repository.ChatRepository
}

func NewRepositoryManager(env *bootstrap.Env) *RepositoryManager {
	log.Printf("NewRepositoryManager called with env: %+v", env)

	chatRepo := repository.NewChatRepository(db.GetChatCollection())
	if chatRepo == nil {
		log.Fatal("chatRepo is nil")
	}
	log.Printf("chatRepo created: %+v", chatRepo)

	log.Printf("Initializing TherapistRepository...")

	repoManager := &RepositoryManager{
		chatRepo: chatRepo,
	}

	return repoManager
}
