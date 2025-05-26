package repository

import (
	"consultationservice/internal/model"
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ClientRepository struct {
	MongoDBCollection *mongo.Collection
}

func NewClientRepository(collection *mongo.Collection) *ClientRepository {
	return &ClientRepository{MongoDBCollection: collection}
}
func (r *ClientRepository) CreateClientMatchingProfile(client *model.Client) (interface{}, error) {
	filter := bson.D{{Key: "profile_id", Value: client.ProfileId}}
	var existingClient model.Client
	err := r.MongoDBCollection.FindOne(context.Background(), filter).Decode(&existingClient)
	if err == nil {
		log.Println("Client already exists with profile_id:", client.ProfileId)
		return nil, errors.New("Client with this profile already exists")
	}
	if err != mongo.ErrNoDocuments {
		log.Println("Error checking if client exists:", err)
		return nil, errors.New("Failed to check if client exists")
	}
	res, err := r.MongoDBCollection.InsertOne(context.Background(), client)
	if err != nil {
		log.Println("ClientRepository CreateClientMatchingProfile err:", err)
		return nil, errors.New("Failed to insert into client")
	}
	return res, nil
}

func (r *ClientRepository) FindClientMatchingProfile(clientId string) (*model.Client, error) {
	data := new(model.Client)
	err := r.MongoDBCollection.FindOne(context.Background(), bson.D{{Key: "profile_id", Value: clientId}}).Decode(&data)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, errors.New("No client found")
	}
	return data, nil
}
func (r *ClientRepository) FindAllClientMatchingProfiles(client *model.Client) ([]model.Client, error) {
	result, err := r.MongoDBCollection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, errors.New("Failed to find all client")
	}
	var results []model.Client
	err = result.All(context.Background(), &results)
	if err != nil {
		return nil, errors.New("Failed")
	}
	return results, nil
}
func (r *ClientRepository) DeleteMatchingProfile(client *model.Client) (int64, error) {
	_, err := r.MongoDBCollection.DeleteOne(context.Background(),
		bson.D{{Key: "profile_id", Value: client.ProfileId}})
	if err != nil {
		return 0, errors.New("Failed to delete client")
	}
	return 1, nil
}
func (r *ClientRepository) UpdateMatchingProfile(profileId string, client *model.Client) (interface{}, error) {
	var result model.Client
	client.ProfileId = profileId
	update := bson.D{{Key: "$set", Value: client}}

	err := r.MongoDBCollection.FindOneAndUpdate(
		context.Background(),
		bson.D{{Key: "profile_id", Value: profileId}},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&result)
	if err != nil {
		return nil, errors.New("Failed to update client")
	}
	log.Println("ClientRepository UpdateMatchingProfile result:", result)
	return result, nil
}
