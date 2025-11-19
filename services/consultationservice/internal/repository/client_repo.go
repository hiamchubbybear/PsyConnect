package repository

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"consultationservice/internal/model"
	"consultationservice/internal/redis"
)

type ClientRepository struct {
	MongoDBCollection *mongo.Collection
	Redis             redis.RedisStore
}

func NewClientRepository(collection *mongo.Collection, redisStore redis.RedisStore) *ClientRepository {
	return &ClientRepository{MongoDBCollection: collection, Redis: redisStore}
}

func (r *ClientRepository) redisKey(profileId string) string {
	return redis.NewKeyBuilder("psyconnect").Build("consultation", "client", profileId)
}

func (r *ClientRepository) CreateClientMatchingProfile(client *model.Client) (interface{}, error) {
	filter := bson.D{{Key: "profile_id", Value: client.ProfileId}}
	var existingClient model.Client
	err := r.MongoDBCollection.FindOne(context.Background(), filter).Decode(&existingClient)
	if err == nil {
		log.Println("Client already exists with profile_id:", client.ProfileId)
		return nil, errors.New("client with this profile already exists")
	}
	if err != mongo.ErrNoDocuments {
		log.Println("Error checking if client exists:", err)
		return nil, errors.New("failed to check if client exists")
	}

	res, err := r.MongoDBCollection.InsertOne(context.Background(), client)
	if err != nil {
		log.Println("ClientRepository CreateClientMatchingProfile err:", err)
		return nil, errors.New("failed to insert client")
	}

	_ = r.Redis.Delete(context.Background(), r.redisKey(client.ProfileId))

	return res, nil
}

func (r *ClientRepository) FindClientMatchingProfile(clientId string) (*model.Client, error) {
	key := r.redisKey(clientId)
	ctx := context.Background()

	var cachedClient model.Client
	if err := r.Redis.Get(ctx, key, &cachedClient); err == nil {
		log.Println("Cache hit for client:", clientId)
		return &cachedClient, nil
	}

	var client model.Client
	err := r.MongoDBCollection.FindOne(ctx, bson.M{"profile_id": clientId}).Decode(&client)
	if err != nil {
		return nil, err
	}

	_ = r.Redis.Set(ctx, key, client)

	return &client, nil
}

func (r *ClientRepository) FindAllClientMatchingProfiles() ([]model.Client, error) {
	cursor, err := r.MongoDBCollection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, errors.New("failed to find all clients")
	}
	defer cursor.Close(context.Background())

	var results []model.Client
	if err := cursor.All(context.Background(), &results); err != nil {
		return nil, errors.New("failed to decode results")
	}
	return results, nil
}

func (r *ClientRepository) DeleteMatchingProfile(client *model.Client) (int64, error) {
	res, err := r.MongoDBCollection.DeleteOne(context.Background(),
		bson.D{{Key: "profile_id", Value: client.ProfileId}})
	if err != nil {
		return 0, errors.New("failed to delete client")
	}

	_ = r.Redis.Delete(context.Background(), r.redisKey(client.ProfileId))

	return res.DeletedCount, nil
}

func (r *ClientRepository) UpdateMatchingProfile(profileId string, client *model.Client) (interface{}, error) {
	var updated model.Client
	client.ProfileId = profileId
	update := bson.D{{Key: "$set", Value: client}}

	err := r.MongoDBCollection.FindOneAndUpdate(
		context.Background(),
		bson.D{{Key: "profile_id", Value: profileId}},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updated)
	if err != nil {
		return nil, errors.New("failed to update client")
	}

	_ = r.Redis.Delete(context.Background(), r.redisKey(profileId))

	log.Println("ClientRepository UpdateMatchingProfile result:", updated)
	return updated, nil
}
