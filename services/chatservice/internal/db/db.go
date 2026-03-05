package db

import (
	"chatservice/bootstrap"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient *mongo.Client
	dbName      string
	once        sync.Once
)

func InitDB(env *bootstrap.Env) *mongo.Client {
	once.Do(func() {
		mongoURI := env.DatabaseMongoUri
		if mongoURI == "" {
			dbUser := env.DatabaseUser
			dbPass := env.DatabasePassword
			dbHost := env.DatabaseHost
			dbPort := env.DatabasePort

			log.Println("DB_USER =", dbUser)
			log.Println("DB_HOST =", dbHost)
			log.Println("DB_PORT =", dbPort)
			if dbUser == "" || dbPass == "" {
				log.Println("Warning: DB_USER or DB_PASS not set, connecting without auth")
				mongoURI = fmt.Sprintf("mongodb://%s:%s/", dbHost, dbPort)
			} else {
				mongoURI = fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin", dbUser, dbPass, dbHost, dbPort)
			}
		}

		dbName = env.DatabaseName
		if dbName == "" {
			log.Fatal("DB_NAME is required")
		}

		log.Printf("Connecting to MongoDB at %s", mongoURI)

		opts := options.Client().
			ApplyURI(mongoURI).
			SetConnectTimeout(10 * time.Second)

		client, err := mongo.Connect(context.TODO(), opts)
		if err != nil {
			log.Fatal("MongoDB connection error:", err)
		}

		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = client.Ping(ctx, nil)
		if err != nil {
			log.Fatal("Failed to ping MongoDB: ", err)
		}

		log.Println("DB_NAME =", dbName)
		mongoClient = client
		log.Println("Successfully connected and pinged MongoDB!")
	})

	return mongoClient
}

func GetChatCollection(env *bootstrap.Env) *mongo.Collection {
	chat := InitDB(env)
	return chat.Database(dbName).Collection("chat")
}
func GetConversationCollection(env *bootstrap.Env) *mongo.Collection {
	chat := InitDB(env)
	return chat.Database(dbName).Collection("conversation")
}
func CloseDB() {
	if mongoClient != nil {
		if err := mongoClient.Disconnect(context.TODO()); err != nil {
			log.Fatal("Error disconnecting MongoDB:", err)
		}
		log.Println("MongoDB connection closed.")
	}
}
