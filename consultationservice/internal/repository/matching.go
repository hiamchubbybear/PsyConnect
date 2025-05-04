package repository

import (
	"consultationservice/internal/dto"
	"consultationservice/internal/grpc/handler"
	"consultationservice/internal/model"
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MatchingRepository struct {
	MongoDBCollection *mongo.Collection
	clientRepo        *ClientRepository
	grpcProfile       *handler.ProfileGrpc
}

func NewMatchingRepository(collection *mongo.Collection, clientRepo *ClientRepository, grpcProfile *handler.ProfileGrpc) *MatchingRepository {
	return &MatchingRepository{
		MongoDBCollection: collection,
		clientRepo:        clientRepo,
		grpcProfile:       grpcProfile,
	}
}

func (r *MatchingRepository) FilterAllMatching(matchProfileId string, page int64) (interface{}, error) {
	ctx := context.TODO()
	filter := bson.M{}
	matching, err := r.clientRepo.FindClientMatchingProfile(matchProfileId)
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("Client profile not found")
	}
	if len(matching.Languages) > 0 {
		filter["languages"] = bson.M{"$in": matching.Languages}
	}
	if len(matching.TherapistSpecialization) > 0 {
		filter["specialization"] = bson.M{"$in": matching.TherapistSpecialization}
	}
	if matching.RangePrice > 0 {
		filter["rage_price"] = bson.M{
			"$gte": 0,
			"$lte": matching.RangePrice,
		}
	}
	if len(matching.ConsultationModes) > 0 && matching.ConsultationModes[0] != "optional" {
		filter["consultation_modes"] = bson.M{"$in": matching.ConsultationModes}
	}
	if len(matching.PreferredTherapistLanguage) > 0 {
		filter["languages"] = bson.M{"$in": matching.PreferredTherapistLanguage}
	}
	if !matching.IsFlexibleWithSchedule {
		if len(matching.Availability.Days) > 0 {
			filter["availability.days"] = bson.M{"$in": matching.Availability.Days}
		}
		if len(matching.Availability.TimeSlots) > 0 {
			filter["availability.time_slots"] = bson.M{"$in": matching.Availability.TimeSlots}
		}
	}
	if page < 1 {
		page = 1
	}
	limit := int64(10)
	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSkip((page - 1) * limit)

	cursor, err := r.MongoDBCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	var results []model.Therapist
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	log.Printf("[MatchingRepository] Query Results: %+v\n", results)
	return results, nil
}
func (r *MatchingRepository) RequestMatching(request dto.MatchingRequest, profileId string) (bool, error) {
	if request.ThearpistId == profileId {
		return false, errors.New("You can not request to yourself")
	}
	isSuccess, err := r.grpcProfile.FriendRequestAccept(request.ThearpistId, profileId)
	if err != nil {
		return false, err
	}
	return isSuccess, nil
}
