package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/subosito/gotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient *mongo.Client
	dbName      string
	once        sync.Once
)

func InitDB() *mongo.Client {
	once.Do(func() {
		err := gotenv.Load("../.env")

		dbUser := os.Getenv("DB_USER")
		dbPass := os.Getenv("DB_PASS")
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbName = os.Getenv("DB_NAME")
		if (dbUser == "" ||
			dbPass == "" ||
			dbHost == "" ||
			dbPort == "" ||
			dbName == "") && err != nil {
			log.Fatal("Error loading .env file")
		}

		uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin", dbUser, dbPass, dbHost, dbPort)

		if dbUser == "" || dbPass == "" {
			uri = fmt.Sprintf("mongodb://%s:%s/", dbHost, dbPort)
		}
		log.Printf("Connecting to MongoDB at %s", uri)

		opts := options.Client().
			ApplyURI(uri).
			SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1)).
			SetConnectTimeout(5 * time.Second)

		client, err := mongo.Connect(context.TODO(), opts)
		if err != nil {
			log.Fatal("MongoDB connection error:", err)
		}
		log.Println("DB_NAME =", dbName)
		mongoClient = client
		log.Println("Successfully connected to MongoDB!")
	})

	return mongoClient
}

func GetChatCollection() *mongo.Collection {
	chat := InitDB()
	return chat.Database(dbName).Collection("chat")
}
func GetConversationCollection() *mongo.Collection {
	chat := InitDB()
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
