package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"notificationservice/pkg/api"
	"notificationservice/pkg/config"
	"notificationservice/pkg/database"
	"notificationservice/pkg/email"
	"notificationservice/pkg/firebase"
	"notificationservice/pkg/handlers"
	"notificationservice/pkg/kafka"
)

func main() {
	log.Println("🚀 Starting PsyConnect Notification Service (Go)")

	
	cfg := config.Load()
	log.Printf("✅ Configuration loaded")
	log.Printf("📧 SMTP: %s:%d", cfg.SMTPHost, cfg.SMTPPort)
	log.Printf("📨 Kafka Brokers: %v", cfg.KafkaBrokers)
	log.Printf("🌍 Environment: %s", cfg.Environment)

	
	db, err := database.NewDatabase(cfg.DatabaseDSN)
	if err != nil {
		log.Printf("⚠️  Database connection failed: %v (continuing without database)", err)
		db = nil
	} else {
		log.Println("✅ Database connected")
		defer db.Close()
	}

	
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

	
	emailService := email.NewEmailService(cfg)
	log.Println("✅ Email service initialized")

	
	var notifService *handlers.NotificationService
	if db != nil {
		notifService = handlers.NewNotificationService(db, fcmService)
		log.Println("✅ Notification service initialized")
	}

	
	consumer := kafka.NewConsumer(cfg, emailService, notifService)
	log.Println("✅ Kafka consumer initialized")

	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	
	go func() {
		if err := consumer.Start(ctx); err != nil {
			log.Printf("❌ Kafka consumer error: %v", err)
		}
	}()

	
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

	
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	log.Println("✅ Notification service is running...")
	log.Printf("🏥 Health check: http://localhost:%s/health", cfg.ServerPort)

	<-sigChan
	log.Println("🛑 Shutting down gracefully...")

	
	cancel()

	
	time.Sleep(2 * time.Second)

	log.Println("👋 Notification service stopped")
}
