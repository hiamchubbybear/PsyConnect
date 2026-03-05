package main

import (
	"chatservice/internal/chat/repository/repository"
	"chatservice/internal/chat/service"

	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"chatservice/bootstrap"
	"chatservice/internal/db"
	"chatservice/internal/handler"
	"chatservice/internal/kafka"
	"chatservice/internal/middleware"
	"chatservice/internal/signaling"
	"chatservice/internal/ws"
	wsmiddleware "chatservice/internal/ws/middleware"
	"chatservice/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {

	env := bootstrap.LoadEnv()
	db.InitDB(env)
	kafkaBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	if len(kafkaBrokers) == 0 || kafkaBrokers[0] == "" {
		kafkaBrokers = []string{"kafka:9094"}
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

	hub := ws.NewHub()
	go hub.Run()

	
	signalingService := signaling.NewSignalingService(hub)
	log.Printf("📞 WebRTC signaling service initialized")

	
	go func() {
		ticker := time.NewTicker(30 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			signalingService.CleanupOldSessions()
		}
	}()

	
	kafkaProducer, err := kafka.NewProducer(env)
	if err != nil {
		log.Fatalf("Failed to initialize Kafka producer: %v", err)
	}
	defer kafkaProducer.Close()

	chatService := service.NewChatService(repoManager.MessageRepo, kafkaProducer)
	hub.RegisterHandler("chat", func(hub *ws.Hub, msg ws.Message) {
		chatService.HandleChatMessage(hub, msg)
	})

	kafkaLogger.Info("WebSocket hub initialized with handlers", nil)

	router := gin.New()
	router.Use(gin.Recovery())

	router.Use(middleware.LoggingMiddleware(kafkaLogger))

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

	router.GET("/ws", gin.WrapF(wsmiddleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWS(hub, w, r)
	})))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	uri := fmt.Sprintf("%s:%s", env.Addr, env.Port)

	kafkaLogger.Info("Chat service started successfully", map[string]interface{}{
		"uri": uri,
	})

	log.Printf("Chat Service (REST + WebSocket) running at %s", uri)
	log.Printf("REST API endpoints: http://%s/chats, /conversations", uri)
	log.Printf("WebSocket endpoint: ws://%s/ws", uri)
	router.Run(uri)
}
