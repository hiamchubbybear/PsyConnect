package db

import (
	"context"
	"fmt"
	"github.com/subosito/gotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

var (
	mongoClient *mongo.Client
	dbName      string
)

func InitDB() *mongo.Client {
	err := gotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName = os.Getenv("DB_NAME")
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

	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		log.Fatal("MongoDB ping failed:", err)
	}

	mongoClient = client
	log.Println("Successfully connected to MongoDB!")
	return mongoClient
}

func InitTherapistDB() *mongo.Collection {
	var mongoClient = InitDB()
	if mongoClient == nil {
		log.Fatal("MongoDB client is not initialized! Call InitDB() first.")
	}
	return mongoClient.Database(dbName).Collection("therapists")
}

func CloseDB() {
	if mongoClient != nil {
		if err := mongoClient.Disconnect(context.TODO()); err != nil {
			log.Fatal("Error disconnecting MongoDB:", err)
		}
		log.Println("MongoDB connection closed.")
	}
}
