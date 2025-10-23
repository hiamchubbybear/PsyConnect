package route

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"chatservice/bootstrap"
	"chatservice/internal/handler"
	"chatservice/internal/repository"
)

func RouterInit(env *bootstrap.Env, repoManager *repository.RepositoryManager) {
	uri := fmt.Sprintf("%v:%v", env.Addr, env.Port)

	router := gin.New()
	router.Use(gin.Recovery())

	chatHandler := handler.NewChatHandler(env, repoManager)
	conversationHandler := handler.NewConversationHandler(env, repoManager)

	router.POST("/chats", chatHandler.CreateChat)
	router.GET("/chats/:id", chatHandler.GetChatByID)
	router.GET("/chats/conversation/:conversationId", chatHandler.GetChatsByConversation)
	router.PUT("/chats/:id", chatHandler.UpdateChat)
	router.DELETE("/chats/:id", chatHandler.DeleteChat)

	router.POST("/conversations", chatHandler.CreateConversation)
	router.GET("/conversations/by-users", conversationHandler.GetConversationByUsers)
	router.Run(uri)
}
