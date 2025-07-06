package repository

import (
	"consultationservice/internal/model"
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MatchRepository struct {
	MatchCollection *mongo.Collection
}

func NewMatchRepository(MatchCollection *mongo.Collection) *MatchRepository {
	repo := &MatchRepository{
		MatchCollection: MatchCollection,
	}
	return repo
}

func (r *MatchRepository) CreateMatch(match model.Match) error {
	if r == nil {
		log.Println("MatchRepository is nil")
		return errors.New("match repository is nil")
	}

	if r.MatchCollection == nil {
		log.Println("MatchCollection is nil in CreateMatch")
		return errors.New("match collection is nil")
	}

	ctx := context.Background()
	match.MatchedAt = time.Now().UTC()
	match.UpdatedAt = time.Now().UTC()

	log.Printf("Attempting to insert match into collection: %s", r.MatchCollection.Name())
	log.Printf("Match data: %+v", match)
	_, err := r.MatchCollection.InsertOne(ctx, match)
	if err != nil {
		log.Printf("Failed to insert match: %v", err)
		return errors.New("failed to insert match")
	}

	log.Println("Match inserted successfully")
	return nil
}
func (r *MatchRepository) GetMatchesByClientId(clientId string) ([]model.Match, error) {
	ctx := context.Background()
	cursor, err := r.MatchCollection.Find(ctx, bson.M{"client_id": clientId})
	if err != nil {
		log.Println("Failed to find matches:", err)
		return nil, errors.New("failed to get matches")
	}
	var matches []model.Match
	if err := cursor.All(ctx, &matches); err != nil {
		log.Println("Failed to decode matches:", err)
		return nil, err
	}
	return matches, nil
}

func (r *SwipeRepository) InsertSwipe(swipe model.ClientSwipe) error {
	ctx := context.Background()
	_, err := r.swipeRepo.InsertOne(ctx, swipe)
	if err != nil {
		log.Println("Failed to insert swipe:", err)
		return errors.New("failed to insert swipe")
	}
	return nil
}

func (r *SwipeRepository) SwipeAndMatch(clientId, therapistId string, points float32, reasons []string) error {
	ctx := context.Background()

	_, err := r.swipeRepo.DeleteOne(ctx, bson.M{
		"client_id":    clientId,
		"therapist_id": therapistId,
	})
	if err != nil {
		log.Println("Failed to delete existing swipe:", err)
		return errors.New("failed to delete old swipe")
	}

	match := model.NewMatch(clientId, therapistId, "swipe", float64(points), reasons)

	err = r.sessionRepo.matchingRepo.CreateMatch(*match)
	if err != nil {
		log.Println("Failed to insert match:", err)
		return errors.New("failed to insert match")
	}

	log.Printf("Swipe matched: client %s -> therapist %s", clientId, therapistId)
	return nil
}

func (r *MatchRepository) CheckTimeOverlap(day string, newStart, newEnd string, sessions []model.Session) (bool, error) {
	parsedNewStart, err := parseTime(newStart)
	if err != nil {
		return false, errors.New("invalid start time format")
	}

	parsedNewEnd, err := parseTime(newEnd)
	if err != nil {
		return false, errors.New("invalid end time format")
	}

	if !parsedNewStart.Before(parsedNewEnd) {
		return false, errors.New("start time must be before end time")
	}

	for _, session := range sessions {
		if session.StartTime.Weekday().String() == day {
			if isOverlap(parsedNewStart, parsedNewEnd, session.StartTime, session.EndTime) {
				return true, nil
			}
		}
	}
	return false, nil
}
func (r *MatchRepository) FilterAllTherapist(clientId string, page int64) ([]model.Match, error) {
	ctx := context.Background()

	if page < 1 {
		page = 1
	}
	limit := int64(10)
	skip := (page - 1) * limit
	opts := options.Find().
		SetLimit(limit).
		SetSkip(skip).
		SetSort(bson.M{"matched_at": -1})

	filter := bson.M{"client_id": clientId}

	cursor, err := r.MatchCollection.Find(ctx, filter, opts)
	if err != nil {
		log.Println("Failed to query match:", err)
		return nil, err
	}

	var matches []model.Match
	if err := cursor.All(ctx, &matches); err != nil {
		log.Println("Failed to decode matches:", err)
		return nil, err
	}

	return matches, nil
}

func parseTime(timeStr string) (time.Time, error) {
	return time.Parse("15:04", timeStr)
}

func isOverlap(start1, end1, start2, end2 time.Time) bool {
	return start1.Before(end2) && start2.Before(end1)
}

// Deprecated : Unable to use . Reason : Scale -> Change last commit : 4fc7730

// func (r *MatchingRepository) FilterAllMatching(matchProfileId string, page int64) (interface{}, error) {
// 	ctx := context.TODO()
// 	filter := bson.M{}
// 	matching, err := r.clientRepo.FindClientMatchingProfile(matchProfileId)
// 	if err != nil {
// 		log.Println(err.Error())
// 		return nil, errors.New("client profile not found")
// 	}
// 	if len(matching.Languages) > 0 {
// 		filter["languages"] = bson.M{"$in": matching.Languages}
// 	}
// 	if len(matching.TherapistSpecialization) > 0 {
// 		filter["specialization"] = bson.M{"$in": matching.TherapistSpecialization}
// 	}
// 	if matching.RangePrice > 0 {
// 		filter["rage_price"] = bson.M{
// 			"$gte": 0,
// 			"$lte": matching.RangePrice,
// 		}
// 	}
// 	filter["matched_clients"] = bson.M{"$nin": []string{matchProfileId}}
// 	if len(matching.ConsultationModes) > 0 && matching.ConsultationModes[0] != "optional" {
// 		filter["consultation_modes"] = bson.M{"$in": matching.ConsultationModes}
// 	}
// 	if len(matching.PreferredTherapistLanguage) > 0 {
// 		filter["languages"] = bson.M{"$in": matching.PreferredTherapistLanguage}
// 	}
// 	if !matching.IsFlexibleWithSchedule {
// 		if len(matching.Availability.Days) > 0 {
// 			filter["availability.days"] = bson.M{"$in": matching.Availability.Days}
// 		}
// 		if len(matching.Availability.TimeSlots) > 0 {
// 			filter["availability.time_slots"] = bson.M{"$in": matching.Availability.TimeSlots}
// 		}
// 	}
// 	if page < 1 {
// 		page = 1
// 	}
// 	limit := int64(10)
// 	findOptions := options.Find()
// 	findOptions.SetLimit(limit)
// 	findOptions.SetSkip((page - 1) * limit)

// 	cursor, err := r.therapistRepo.MongoDBCollection.Find(ctx, filter, findOptions)

// 	if err != nil {
// 		return nil, err
// 	}
// 	var results []model.Therapist
// 	if err := cursor.All(ctx, &results); err != nil {
// 		return nil, err
// 	}
// 	log.Printf("[MatchingRepository] Query Results: %+v\n", results)
// 	return results, nil
// }
// func (r *MatchingRepository) RequestMatching(request dto.MatchingRequest, profileId string) (bool, error) {
// 	if request.ThearpistId == profileId {
// 		return false, errors.New("you can not request to yourself")
// 	}
// 	isSuccess, err := r.grpcProfile.FriendRequestAccept(request.ThearpistId, profileId, request.Message)
// 	if err != nil {
// 		return false, err
// 	}
// 	return isSuccess, nil
// }
// func (r *MatchingRepository) FindClientLogs(profileId string) (interface{}, error) {
// 	ctx := context.TODO()
// 	if profileId == "" {
// 		return nil, errors.New("profile Id is invalid ")
// 	}
// 	opts := bson.M{"profile_id": profileId}
// 	cursor, err := r.therapistRepo.MongoDBCollection.Find(ctx, opts)
// 	if err != nil {
// 		return false, errors.New("failed to find matched client")
// 	}
// 	for cursor.Next(ctx) {
// 		var result bson.M
// 		err := cursor.Decode(&result)
// 		if err != nil {
// 			log.Printf("Failted to find matched object")
// 			return false, errors.New("failed to find matched client")
// 		}
// 		return result, nil
// 	}
// 	return nil, errors.New("failed")
// }
// func (r *MatchingRepository) CheckExistedClientLogs(profileId string, clientId string) (bool, error) {
// 	filter := bson.M{
// 		"profile_id": profileId,
// 		"matched_clients": bson.M{
// 			"$ne": clientId,
// 		},
// 	}
// 	count, err := r.therapistRepo.MongoDBCollection.CountDocuments(context.TODO(), filter)
// 	if err != nil {
// 		return false, err
// 	}
// 	return count > 0, nil
// }

// func (r *MatchingRepository) AddMatchedClient(profileId string, clientId string) (bool, error) {
// 	ctx := context.Background()
// 	existed, err := r.CheckExistedClientLogs(profileId, clientId)
// 	if err != nil {
// 		return false, err
// 	}
// 	if existed {
// 		return true, nil
// 	}
// 	_, error := r.clientRepo.FindClientMatchingProfile(clientId)
// 	if error != nil {
// 		log.Printf("failed find clients: %v", error.Error())
// 		return false, nil
// 	}
// 	_, error = r.therapistRepo.FindTherapistMatchingProfile(profileId)
// 	if error != nil {
// 		log.Printf("Find find therapist : %v", error.Error())
// 		return false, nil
// 	}
// 	filter := bson.M{"profile_id": profileId}
// 	update := bson.M{
// 		"$push": bson.M{
// 			"matched_clients": clientId,
// 		},
// 	}
// 	res, updateError := r.therapistRepo.MongoDBCollection.UpdateOne(ctx, filter, update)
// 	if updateError != nil {
// 		return false, errors.New("failed to add client matched")
// 	}
// 	if res.MatchedCount == 0 {
// 		log.Println("No document matched for profile_id:", profileId)
// 		return false, errors.New("therapist not found")
// 	}
// 	if res.ModifiedCount == 0 {
// 		log.Println("Document matched but not modified, maybe already contains clientId?")
// 	}

// 	return true, nil
// }
// func (r *MatchRepository) CheckValidSessionTimeAndDays(time dto.SessionTime, clientId string, therapistId string) (bool, error) {
// 	ctx := context.Background()
// 	if !IsValidDay(time.Day) {
// 		return false, errors.New("Invalid day of week")
// 	}
// 	if time.StartTime == "" || time.EndTime == "" {
// 		return false, errors.New("Invalid time")
// 	}
// 	var clientDoc struct {
// 		CurrentSession []string `bson:"current_session"`
// 	}
// 	if err := r.clientRepo.MongoDBCollection.FindOne(ctx,
// 		bson.M{"profile_id": clientId},
// 	).Decode(&clientDoc); err != nil {
// 		return false, errors.New("cannot find client")
// 	}
// 	var therapistDoc struct {
// 		CurrentSession []string `bson:"current_session"`
// 	}
// 	if err := r.therapistRepo.MongoDBCollection.FindOne(ctx,
// 		bson.M{"profile_id": therapistId},
// 	).Decode(&therapistDoc); err != nil {
// 		return false, errors.New("cannot find therapist")
// 	}
// 	var sessions []model.Session

// 	allSessionIDs := append(clientDoc.CurrentSession, therapistDoc.CurrentSession...)
// 	if len(allSessionIDs) != 0 {
// 		cursor, err := r.sessionRepo.MongoDBCollection.Find(ctx, bson.M{
// 			"_id": bson.M{"$in": allSessionIDs},
// 		})
// 		if err != nil {
// 			log.Printf(err.Error())
// 			return false, errors.New("failed to query sessions")
// 		}
// 		if err := cursor.All(ctx, &sessions); err != nil {
// 			return false, errors.New("failed to decode sessions")
// 		}

// 	}
// 	newStart, err1 := parseTime(time.StartTime)
// 	newEnd, err2 := parseTime(time.EndTime)
// 	if err1 != nil || err2 != nil {
// 		return false, errors.New("invalid input time format (expected HH:MM)")
// 	}

// 	for _, s := range sessions {
// 		if s.StartTime.Weekday().String() == time.Day.String() {
// 			if isOverlapping(s.StartTime, s.EndTime, newStart, newEnd) {
// 				return false, errors.New("conflict with existing session")
// 			}
// 		}
// 	}
// 	return true, nil
// }
// func IsValidDay(d enum.DayOfWeek) bool {
// 	return d >= enum.Monday && d <= enum.Sunday
// }

// func parseTime(t string) (time.Time, error) {
// 	return time.Parse("15:04", t)
// }

// func isOverlapping(start1, end1, start2, end2 time.Time) bool {
// 	return start1.Before(end2) && start2.Before(end1)
// }
