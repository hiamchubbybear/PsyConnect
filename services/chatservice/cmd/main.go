package main

import (
	"fmt"
	"log"

	"chatservice/bootstrap"
	"chatservice/internal/db"
	"chatservice/internal/handler"
	"chatservice/internal/repository"
	"chatservice/internal/route"
	"chatservice/internal/ws"
	"chatservice/internal/ws/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	db.InitDB()

	// Load environment
	env := bootstrap.LoadEnv()
	repoManager := repository.NewRepositoryManager(env)

	// Initialize WebSocket hub
	hub := ws.NewHub(repoManager.MessageRepo)
	go hub.Run()

	// Setup Gin router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(cors.Default())

	// Initialize handlers
	chatHandler := handler.NewChatHandler(env, repoManager)
	conversationHandler := handler.NewConversationHandler(env, repoManager)

	// REST API endpoints
	router.POST("/chats", chatHandler.CreateChat)
	router.GET("/chats/:id", chatHandler.GetChatByID)
	router.GET("/chats/conversation/:conversationId", chatHandler.GetChatsByConversation)
	router.PUT("/chats/:id", chatHandler.UpdateChat)
	router.DELETE("/chats/:id", chatHandler.DeleteChat)

	router.POST("/conversations", chatHandler.CreateConversation)
	router.GET("/conversations/by-users", conversationHandler.GetConversationByUsers)

	// WebSocket endpoint (same port as REST API)
	router.GET("/ws", gin.WrapF(middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWS(hub, w, r)
	})))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Start unified server
	uri := fmt.Sprintf("%s:%s", env.Addr, env.Port)
	log.Printf("🚀 Chat Service (REST + WebSocket) running at %s", uri)
	log.Printf("   📡 REST API endpoints: http://%s/chats, /conversations", uri)
	log.Printf("   🔌 WebSocket endpoint: ws://%s/ws", uri)

	router.Run(uri)
}
