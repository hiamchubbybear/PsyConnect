package repository

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"testing"
	"time"
)

func newMongoClient() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(""))
	if err != nil {
		log.Fatal("Failed to connect MongoDB: ", err)
	}
	log.Println("MongoDB connected!")
	return client
}
func TestMongoOperations(t *testing.T) {
	client := newMongoClient()
	defer client.Disconnect(context.Background())
	db := client.Database("test_db")
	collection := db.Collection("test_collection")
	doc := bson.M{"name": "John Doe", "age": 30}
	insertResult, err := collection.InsertOne(context.Background(), doc)
	assert.Nil(t, err)
	assert.NotNil(t, insertResult.InsertedID)
	filter := bson.M{"name": "John Doe"}
	var result bson.M
	err = collection.FindOne(context.Background(), filter).Decode(&result)
	assert.Nil(t, err)
	assert.Equal(t, "John Doe", result["name"])
	update := bson.M{"$set": bson.M{"age": 35}}
	updateResult, err := collection.UpdateOne(context.Background(), filter, update)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), updateResult.ModifiedCount)
	deleteResult, err := collection.DeleteOne(context.Background(), filter)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), deleteResult.DeletedCount)

	log.Println(" MongoDB Mock Test Passed!")
}
