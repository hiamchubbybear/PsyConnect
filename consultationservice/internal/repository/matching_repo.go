package repository

import (
	"consultationservice/internal/dto"
	"consultationservice/internal/enum"
	grpc "consultationservice/internal/grpc/handler"
	"consultationservice/internal/model"
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MatchingRepository struct {
	clientRepo    *ClientRepository
	therapistRepo *TherapistRepository
	grpcProfile   *grpc.ProfileGrpc
	sessionRepo   *SessionRepository
}

func NewMatchingRepository(clientRepo *ClientRepository, grpcProfile *grpc.ProfileGrpc, therapistRepo *TherapistRepository, sessionRepo *SessionRepository) *MatchingRepository {
	return &MatchingRepository{
		clientRepo:    clientRepo,
		grpcProfile:   grpcProfile,
		therapistRepo: therapistRepo,
		sessionRepo:   sessionRepo,
	}
}

func (r *MatchingRepository) FilterAllMatching(matchProfileId string, page int64) (interface{}, error) {
	ctx := context.TODO()
	filter := bson.M{}
	matching, err := r.clientRepo.FindClientMatchingProfile(matchProfileId)
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("client profile not found")
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
	filter["matched_clients"] = bson.M{"$nin": []string{matchProfileId}}
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

	cursor, err := r.therapistRepo.MongoDBCollection.Find(ctx, filter, findOptions)

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
		return false, errors.New("you can not request to yourself")
	}
	isSuccess, err := r.grpcProfile.FriendRequestAccept(request.ThearpistId, profileId, request.Message)
	if err != nil {
		return false, err
	}
	return isSuccess, nil
}
func (r *MatchingRepository) FindClientLogs(profileId string) (interface{}, error) {
	ctx := context.TODO()
	if profileId == "" {
		return nil, errors.New("profile Id is invalid ")
	}
	opts := bson.M{"profile_id": profileId}
	cursor, err := r.therapistRepo.MongoDBCollection.Find(ctx, opts)
	if err != nil {
		return false, errors.New("failed to find matched client")
	}
	for cursor.Next(ctx) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			log.Printf("Failted to find matched object")
			return false, errors.New("failed to find matched client")
		}
		return result, nil
	}
	return nil, errors.New("failed")
}
func (r *MatchingRepository) CheckExistedClientLogs(profileId string, clientId string) (bool, error) {
	filter := bson.M{
		"profile_id": profileId,
		"matched_clients": bson.M{
			"$ne": clientId,
		},
	}
	count, err := r.therapistRepo.MongoDBCollection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *MatchingRepository) AddMatchedClient(profileId string, clientId string) (bool, error) {
	ctx := context.Background()
	existed, err := r.CheckExistedClientLogs(profileId, clientId)
	if err != nil {
		return false, err
	}
	if existed {
		return true, nil
	}
	_, error := r.clientRepo.FindClientMatchingProfile(clientId)
	if error != nil {
		log.Printf("failed find clients: %v", error.Error())
		return false, nil
	}
	_, error = r.therapistRepo.FindTherapistMatchingProfile(profileId)
	if error != nil {
		log.Printf("Find find therapist : %v", error.Error())
		return false, nil
	}
	filter := bson.M{"profile_id": profileId}
	update := bson.M{
		"$push": bson.M{
			"matched_clients": clientId,
		},
	}
	res, updateError := r.therapistRepo.MongoDBCollection.UpdateOne(ctx, filter, update)
	if updateError != nil {
		return false, errors.New("failed to add client matched")
	}
	if res.MatchedCount == 0 {
		log.Println("No document matched for profile_id:", profileId)
		return false, errors.New("therapist not found")
	}
	if res.ModifiedCount == 0 {
		log.Println("Document matched but not modified, maybe already contains clientId?")
	}

	return true, nil
}
func (r *MatchingRepository) CheckValidSessionTimeAndDays(time dto.SessionTime, clientId string, therapistId string) (bool, error) {
	ctx := context.Background()
	if !IsValidDay(time.Day) {
		return false, errors.New("Invalid day of week")
	}
	if time.StartTime == "" || time.EndTime == "" {
		return false, errors.New("Invalid time")
	}
	var clientDoc struct {
		CurrentSession []string `bson:"current_session"`
	}
	if err := r.clientRepo.MongoDBCollection.FindOne(ctx,
		bson.M{"profile_id": clientId},
	).Decode(&clientDoc); err != nil {
		return false, errors.New("cannot find client")
	}
	var therapistDoc struct {
		CurrentSession []string `bson:"current_session"`
	}
	if err := r.therapistRepo.MongoDBCollection.FindOne(ctx,
		bson.M{"profile_id": therapistId},
	).Decode(&therapistDoc); err != nil {
		return false, errors.New("cannot find therapist")
	}
	allSessionIDs := append(clientDoc.CurrentSession, therapistDoc.CurrentSession...)
	cursor, err := r.sessionRepo.MongoDBCollection.Find(ctx, bson.M{
		"_id": bson.M{"$in": allSessionIDs},
	})
	if err != nil {
		return false, errors.New("failed to query sessions")
	}

	var sessions []model.Session
	if err := cursor.All(ctx, &sessions); err != nil {
		return false, errors.New("failed to decode sessions")
	}

	newStart, err1 := parseTime(time.StartTime)
	newEnd, err2 := parseTime(time.EndTime)
	if err1 != nil || err2 != nil {
		return false, errors.New("invalid input time format (expected HH:MM)")
	}

	for _, s := range sessions {
		if s.StartTime.Weekday().String() == time.Day.String() {
			if isOverlapping(s.StartTime, s.EndTime, newStart, newEnd) {
				return false, errors.New("conflict with existing session")
			}
		}
	}
	return true, nil
}
func IsValidDay(d enum.DayOfWeek) bool {
	return d >= enum.Monday && d <= enum.Sunday
}
func parseTime(t string) (time.Time, error) {
	return time.Parse("15:04", t)
}

func isOverlapping(start1, end1, start2, end2 time.Time) bool {
	return start1.Before(end2) && start2.Before(end1)
}
