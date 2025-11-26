package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"chatservice/bootstrap"
	"chatservice/internal/db"
	"chatservice/internal/handler"
	"chatservice/internal/middleware"
	"chatservice/internal/repository"
	"chatservice/internal/ws"
	wsmiddleware "chatservice/internal/ws/middleware"
	"chatservice/pkg/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	db.InitDB()

	// Load environment
	env := bootstrap.LoadEnv()

	// Initialize Kafka Logger
	kafkaBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	if len(kafkaBrokers) == 0 || kafkaBrokers[0] == "" {
		kafkaBrokers = []string{"localhost:9092"}
	}

	kafkaLogger := logger.NewKafkaLogger(logger.Config{
		Brokers:     kafkaBrokers,
		Topic:       "logging-service",
		ServiceName: "chat-service",
		Environment: os.Getenv("ENVIRONMENT"),
		Version:     "1.0.0",
	})
	defer kafkaLogger.Close()

	kafkaLogger.Info("Chat service starting", map[string]interface{}{
		"port":        env.Port,
		"environment": os.Getenv("ENVIRONMENT"),
	})

	repoManager := repository.NewRepositoryManager(env)

	// Initialize WebSocket hub
	hub := ws.NewHub(repoManager.MessageRepo)
	go hub.Run()

	kafkaLogger.Info("WebSocket hub initialized", nil)

	// Setup Gin router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(cors.Default())

	// Add logging middleware
	router.Use(middleware.LoggingMiddleware(kafkaLogger))

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
	router.GET("/ws", gin.WrapF(wsmiddleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWS(hub, w, r)
	})))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Start unified server
	uri := fmt.Sprintf("%s:%s", env.Addr, env.Port)

	kafkaLogger.Info("Chat service started successfully", map[string]interface{}{
		"uri": uri,
	})

	log.Printf("🚀 Chat Service (REST + WebSocket) running at %s", uri)
	log.Printf("   📡 REST API endpoints: http://%s/chats, /conversations", uri)
	log.Printf("   🔌 WebSocket endpoint: ws://%s/ws", uri)

	router.Run(uri)
}
