package route

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"chatservice/bootstrap"
	handlers "chatservice/internal/handler"
	"chatservice/internal/repository"
)

func RouterInit(
	env *bootstrap.Env,
	repoManager *repository.RepositoryManager,
) {

	urI := fmt.Sprintf("%v:%v", env.Addr, env.Port)
	router := gin.Default()

	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()
	chatHandler := handlers.NewChatHandler(env, repoManager)

	router.POST("/chats", chatHandler.CreateChat)
	router.GET("/chats/:id", chatHandler.GetChatByID)
	router.GET("/chats/conversation/:conversationId", chatHandler.GetChatsByConversation)
	router.PUT("/chats/:id", chatHandler.UpdateChat)
	router.DELETE("/chats/:id", chatHandler.DeleteChat)

	router.Run(urI)
}
