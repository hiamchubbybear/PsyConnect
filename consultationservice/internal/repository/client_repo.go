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
	log.Printf("ClientRepository collection name: %v", r.MongoDBCollection.Name())
	log.Printf("Type of clientId: %T, value: %v", clientId, clientId)

	data := model.Client{}
	err := r.MongoDBCollection.FindOne(context.Background(), bson.M{"profile_id": clientId}).Decode(&data)
	if err != nil {
		log.Printf("Không tìm thấy client với profile_id=%v. Lỗi: %v", clientId, err.Error())
		return nil, err
	}

	log.Printf("Tìm thấy client: %+v", data)
	log.Printf("DEBUG - ProfileId value: '%s'", data.ProfileId)
	log.Printf("DEBUG - ProfileId length: %d", len(data.ProfileId))
	log.Printf("DEBUG - ProfileId is empty: %t", data.ProfileId == "")
	log.Printf("DEBUG - ProfileId bytes: %v", []byte(data.ProfileId))

	return &data, nil
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
