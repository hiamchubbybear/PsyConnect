package route

import (
	"chatservice/internal/chat/repository/repository"
	"chatservice/internal/handler"
	"fmt"

	"github.com/gin-gonic/gin"

	"chatservice/bootstrap"
	"chatservice/internal/chat/service"
	"chatservice/internal/kafka"
	"log"
)

func RouterInit(env *bootstrap.Env, repoManager *repository.RepositoryManager) {
	uri := fmt.Sprintf("%v:%v", env.Addr, env.Port)

	router := gin.New()
	router.Use(gin.Recovery())

	
	router.GET("/actuator/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	
	producer, err := kafka.NewProducer(env)
	if err != nil {
		log.Fatalf("Failed to initialize Kafka producer: %v", err)
	}
	
	

	chatService := service.NewChatService(repoManager.MessageRepo, producer)
	chatHandler := handler.NewChatHandler(env, repoManager, chatService)
	conversationHandler := handler.NewConversationHandler(env, repoManager)

	router.POST("/chats/call/start", chatHandler.StartCall)
	router.POST("/chats", chatHandler.CreateChat)
	router.GET("/chats/:id", chatHandler.GetChatByID)
	router.GET("/chats/conversation/:conversationId", chatHandler.GetChatsByConversation)
	router.PUT("/chats/:id", chatHandler.UpdateChat)
	router.DELETE("/chats/:id", chatHandler.DeleteChat)

	router.POST("/conversations", chatHandler.CreateConversation)
	router.GET("/conversations/by-users", conversationHandler.GetConversationByUsers)
	router.GET("/conversations/me", conversationHandler.GetRecentConversations)

	router.Run(uri)
}
