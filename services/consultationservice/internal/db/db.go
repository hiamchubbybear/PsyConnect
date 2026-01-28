package db

import (
	"context"
	"fmt"
	"log"
	"os"
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

func InitDB() *mongo.Client {
	once.Do(func() {
		mongoURI := os.Getenv("MONGO_URI")
		if mongoURI == "" {
			dbUser := os.Getenv("MONGO_USER")
			if dbUser == "" {
				dbUser = os.Getenv("DB_USER")
			}
			dbPass := os.Getenv("MONGO_PASSWORD")
			if dbPass == "" {
				dbPass = os.Getenv("DB_PASS")
			}
			dbHost := os.Getenv("DB_HOST")
			dbPort := os.Getenv("DB_PORT")

			log.Println("DB_USER =", dbUser)
			log.Println("DB_HOST =", dbHost)
			log.Println("DB_PORT =", dbPort)

			if dbHost == "" {
				dbHost = "mongodb"
			}
			if dbPort == "" {
				dbPort = "27017"
			}

			if dbUser == "" || dbPass == "" {
				log.Println("Warning: DB_USER or DB_PASS not set, connecting without auth")
				mongoURI = fmt.Sprintf("mongodb://%s:%s/", dbHost, dbPort)
			} else {
				mongoURI = fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin", dbUser, dbPass, dbHost, dbPort)
			}
		}

		dbName = os.Getenv("DB_NAME")
		if dbName == "" {
			log.Printf("⚠️  Warning: DB_NAME is required, but missing. Using 'consultationservice' as default.")
			dbName = "consultationservice"
		}

		log.Printf("Connecting to MongoDB at %s", mongoURI)

		opts := options.Client().
			ApplyURI(mongoURI).
			SetConnectTimeout(10 * time.Second)

		client, err := mongo.Connect(context.TODO(), opts)
		if err != nil {
			log.Fatalf("❌ FATAL: MongoDB connection error: %v", err)
		}

		// Verify connection
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = client.Ping(ctx, nil)
		if err != nil {
			log.Fatalf("❌ FATAL: Failed to ping MongoDB: %v", err)
		}

		log.Println("DB_NAME =", dbName)
		mongoClient = client
		log.Println("✅ Successfully connected and pinged MongoDB!")
	})

	return mongoClient
}

func GetTherapistCollection() *mongo.Collection {
	client := InitDB()
	return client.Database(dbName).Collection("therapists")
}
func GetClientCollection() *mongo.Collection {
	client := InitDB()
	return client.Database(dbName).Collection("clients")
}
func GetSessionCollection() *mongo.Collection {
	client := InitDB()
	return client.Database(dbName).Collection("sessions")
}
func GetSwipedCollection() *mongo.Collection {
	client := InitDB()
	return client.Database(dbName).Collection("swipes")
}
func GetMatchCollection() *mongo.Collection {
	client := InitDB()
	return client.Database(dbName).Collection("match")
}

func GetPostCollection() *mongo.Collection {
	chat := InitDB()
	return chat.Database(dbName).Collection("posts")
}

func GetReactionCollection() *mongo.Collection {
	client := InitDB()
	return client.Database(dbName).Collection("reactions")
}

func GetCommentCollection() *mongo.Collection {
	client := InitDB()
	return client.Database(dbName).Collection("comments")
}

func GetFollowCollection() *mongo.Collection {
	client := InitDB()
	return client.Database(dbName).Collection("follows")
}

func GetBookmarkCollection() *mongo.Collection {
	client := InitDB()
	return client.Database(dbName).Collection("bookmarks")
}

func CloseDB() {
	if mongoClient != nil {
		if err := mongoClient.Disconnect(context.TODO()); err != nil {
			log.Printf("⚠️ Warning: Error disconnecting MongoDB: %v", err)
		}
		log.Println("MongoDB connection closed.")
	}
}
