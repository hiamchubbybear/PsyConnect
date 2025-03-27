package repository

import (
	"consultationservice/model"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type TherapistRepository struct {
	MongoDBCollection *mongo.Collection
}

func NewTherapistRepository(collection *mongo.Collection) *TherapistRepository {
	return &TherapistRepository{MongoDBCollection: collection}
}

func (r *TherapistRepository) CreateTherapistMatchingProfile(therapist *model.Therapist) (interface{}, error) {
	res, err := r.MongoDBCollection.InsertOne(context.Background(), therapist)
	if err != nil {
		log.Println("TherapistRepository CreateTherapistMatchingProfile err:", err)
		return nil, errors.New("Failed to insert into therapist")
	}
	return res, nil
}
func (r *TherapistRepository) FindTherapistMatchingProfile(therapist *model.Therapist) (interface{}, error) {
	data := therapist
	err := r.MongoDBCollection.FindOne(context.Background(), bson.D{{Key: "profile_id", Value: data.ProfileId}}).Decode(&data)
	if err != nil {
		log.Println("TherapistRepository FindTherapistMatchingProfile err:", err)
		return nil, errors.New("No therapist found")
	}
	return data, nil
}
func (r *TherapistRepository) FindAllTherapistMatchingProfiles(therapist *model.Therapist) ([]model.Therapist, error) {
	result, err := r.MongoDBCollection.Find(context.Background(), bson.D{})
	if err != nil {

		return nil, errors.New("Failed to find all therapist")
	}
	var results []model.Therapist
	err = result.All(context.Background(), &results)
	if err != nil {
		return nil, errors.New("Failed")
	}
	return results, nil
}
func (r *TherapistRepository) DeleteMatchingProfile(therapist *model.Therapist) (int64, error) {
	_, err := r.MongoDBCollection.DeleteOne(context.Background(),
		bson.D{{Key: "profile_id", Value: therapist.ProfileId}})
	if err != nil {
		return 0, errors.New("Failed to delete therapist")
	}
	return 1, nil
}
