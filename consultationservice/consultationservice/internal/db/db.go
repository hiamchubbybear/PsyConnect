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
		if (dbUser == " " ||
			dbPass == " " ||
			dbHost == " " ||
			dbPort == " " ||
			dbName == " ") && err != nil {
			log.Fatal("Error loading .env file")
		}

		uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/", dbUser, dbPass, dbHost, dbPort)

		log.Printf("Connecting to MongoDB at %s", uri)

		opts := options.Client().
			ApplyURI(uri).
			SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1)).
			SetConnectTimeout(5 * time.Second)

		client, err := mongo.Connect(context.TODO(), opts)
		if err != nil {
			log.Fatal("MongoDB connection error:", err)
		}

		mongoClient = client
		log.Println("Successfully connected to MongoDB!")
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
func CloseDB() {
	if mongoClient != nil {
		if err := mongoClient.Disconnect(context.TODO()); err != nil {
			log.Fatal("Error disconnecting MongoDB:", err)
		}
		log.Println("MongoDB connection closed.")
	}
}
