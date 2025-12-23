package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/psyconnect/notificationservice/pkg/api"
	"github.com/psyconnect/notificationservice/pkg/config"
	"github.com/psyconnect/notificationservice/pkg/database"
	"github.com/psyconnect/notificationservice/pkg/email"
	"github.com/psyconnect/notificationservice/pkg/firebase"
	"github.com/psyconnect/notificationservice/pkg/handlers"
	"github.com/psyconnect/notificationservice/pkg/kafka"
)

func main() {
	log.Println("🚀 Starting PsyConnect Notification Service (Go)")

	// Load configuration
	cfg := config.Load()
	log.Printf("✅ Configuration loaded")
	log.Printf("📧 SMTP: %s:%d", cfg.SMTPHost, cfg.SMTPPort)
	log.Printf("📨 Kafka Brokers: %v", cfg.KafkaBrokers)
	log.Printf("🌍 Environment: %s", cfg.Environment)

	// Initialize database
	db, err := database.NewDatabase(cfg.DatabaseDSN)
	if err != nil {
		log.Printf("⚠️  Database connection failed: %v (continuing without database)", err)
		db = nil
	} else {
		log.Println("✅ Database connected")
		defer db.Close()
	}

	// Initialize Firebase/FCM
	var fcmService *firebase.FCMService
	if cfg.FirebaseCredentials != "" {
		fcmService, err = firebase.NewFCMService(cfg.FirebaseCredentials)
		if err != nil {
			log.Printf("⚠️  Firebase initialization failed: %v (continuing without FCM)", err)
			fcmService = nil
		} else {
			log.Println("✅ Firebase/FCM initialized")
		}
	} else {
		log.Println("⚠️  No Firebase credentials provided (push notifications disabled)")
	}

	// Initialize email service
	emailService := email.NewEmailService(cfg)
	log.Println("✅ Email service initialized")

	// Initialize notification service
	var notifService *handlers.NotificationService
	if db != nil {
		notifService = handlers.NewNotificationService(db, fcmService)
		log.Println("✅ Notification service initialized")
	}

	// Initialize Kafka consumer
	consumer := kafka.NewConsumer(cfg, emailService, notifService)
	log.Println("✅ Kafka consumer initialized")

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start Kafka consumer in goroutine
	go func() {
		if err := consumer.Start(ctx); err != nil {
			log.Printf("❌ Kafka consumer error: %v", err)
		}
	}()

	// Start HTTP API server
	if notifService != nil {
		apiHandler := api.NewAPI(emailService, notifService)
		router := api.SetupRouter(apiHandler)

		go func() {
			log.Printf("🌐 HTTP API server listening on port %s", cfg.ServerPort)
			if err := router.Run(":" + cfg.ServerPort); err != nil {
				log.Printf("❌ HTTP server error: %v", err)
			}
		}()
	} else {
		log.Println("⚠️  HTTP API disabled (database not available)")
	}

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	log.Println("✅ Notification service is running...")
	log.Printf("🏥 Health check: http://localhost:%s/health", cfg.ServerPort)

	<-sigChan
	log.Println("🛑 Shutting down gracefully...")

	// Cancel context to stop consumer
	cancel()

	// Give services time to cleanup
	time.Sleep(2 * time.Second)

	log.Println("👋 Notification service stopped")
}
