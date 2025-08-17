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

/*
// Deprecated : Unable to use . Reason : Scale -> Change last commit : 4fc7730 -- Remove

	func (r *MatchRepository) FilterAllMatching(matchProfileId string, page int64) (interface{}, error) {
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
*/
