package repository

import (
	"consultationservice/internal/model"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type TherapistRepository struct {
	MongoDBCollection *mongo.Collection
}

func NewTherapistRepository(collection *mongo.Collection) *TherapistRepository {
	return &TherapistRepository{MongoDBCollection: collection}
}

func (r *TherapistRepository) CreateTherapistMatchingProfile(therapist *model.Therapist) (interface{}, error) {
	filter := bson.D{{Key: "profile_id", Value: therapist.ProfileId}}
	var existingTherapist model.Therapist
	err := r.MongoDBCollection.FindOne(context.Background(), filter).Decode(&existingTherapist)
	if err == nil {
		log.Println("Therapist already exists with profile_id:", therapist.ProfileId)
		return nil, errors.New("Therapist with this profile already exists")
	}
	if err != mongo.ErrNoDocuments {
		log.Println("Error checking if therapist exists:", err)
		return nil, errors.New("Failed to check if therapist exists")
	}
	res, err := r.MongoDBCollection.InsertOne(context.Background(), therapist)
	if err != nil {
		log.Println("TherapistRepository CreateTherapistMatchingProfile err:", err)
		return nil, errors.New("Failed to insert into therapist")
	}
	return res.InsertedID, nil
}

func (r *TherapistRepository) FindTherapistMatchingProfile(therapistId string) (*model.Therapist, error) {
	var data model.Therapist
	err := r.MongoDBCollection.FindOne(context.Background(), bson.D{{Key: "profile_id", Value: therapistId}}).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *TherapistRepository) FindAllTherapistMatchingProfiles() ([]model.Therapist, error) {
	cursor, err := r.MongoDBCollection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, errors.New("Failed to find all therapist")
	}
	var results []model.Therapist
	err = cursor.All(context.Background(), &results)
	if err != nil {
		return nil, errors.New("Failed to decode therapist data")
	}
	return results, nil
}

func (r *TherapistRepository) DeleteMatchingProfile(profileId string) (int64, error) {
	res, err := r.MongoDBCollection.DeleteOne(context.Background(), bson.D{{Key: "profile_id", Value: profileId}})
	if err != nil {
		return 0, errors.New("Failed to delete therapist")
	}
	return res.DeletedCount, nil
}

func (r *TherapistRepository) UpdateMatchingProfile(profileId string, therapist *model.Therapist) (interface{}, error) {
	therapist.ProfileId = profileId
	update := bson.D{{Key: "$set", Value: therapist}}
	var result model.Therapist

	err := r.MongoDBCollection.FindOneAndUpdate(
		context.Background(),
		bson.D{{Key: "profile_id", Value: profileId}},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&result)
	if err != nil {
		return nil, errors.New("Failed to update therapist")
	}
	return result, nil
}

func (r *TherapistRepository) DisableTherapistMatchingProfile(profileId string, isAvailable bool) (bool, error) {
	filter := bson.D{{Key: "profile_id", Value: profileId}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "is_available", Value: isAvailable}}}}

	res, err := r.MongoDBCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return false, errors.New("Failed to update therapist availability")
	}
	return res.MatchedCount > 0, nil
}
