package repository

import (
	"consultationservice/internal/dto"
	"consultationservice/internal/model"
	"consultationservice/internal/redis"
	"consultationservice/internal/repository/infrastructure/external"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/willf/bloom"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SwipeRepository struct {
	clientRepo    *ClientRepository
	therapistRepo *TherapistRepository
	sessionRepo   *SessionRepository
	swipeRepo     *mongo.Collection
	redis         redis.RedisStore
}

func NewSwipeRepository(
	clientRepo *ClientRepository,
	therapistRepo *TherapistRepository,
	sessionRepo *SessionRepository,
	swipeRepo *mongo.Collection,
	redis redis.RedisStore,
) *SwipeRepository {
	return &SwipeRepository{
		clientRepo:    clientRepo,
		therapistRepo: therapistRepo,
		sessionRepo:   sessionRepo,
		swipeRepo:     swipeRepo,
		redis:         redis,
	}
}

func (r *SwipeRepository) InsertSwipes(clientId string, swipes []model.ClientSwipe) error {
	ctx := context.Background()
	if len(swipes) == 0 {
		return nil
	}
	var docs []interface{}
	for _, s := range swipes {
		docs = append(docs, s)
	}
	_, err := r.swipeRepo.InsertMany(ctx, docs)
	if err != nil {
		log.Println("Failed to insert swipes:", err)
		return errors.New("failed to insert swipes")
	}
	return nil
}

func (r *SwipeRepository) PopTop5Swipes(clientId string) ([]model.Therapist, error) {
	ctx := context.Background()

	swipes, err := r.getTop5PendingSwipes(clientId)
	if err != nil {
		return nil, err
	}

	if len(swipes) == 0 {
		log.Println("No pending swipes found, resetting and retrying...")

		if err := r.UpdateAllPendingTherapistSwipeProfile(clientId); err != nil {
			log.Println("Failed to reset swipe statuses:", err)
			return nil, err
		}

		swipes, err = r.getTop5PendingSwipes(clientId)
		if err != nil {
			return nil, err
		}

		if len(swipes) == 0 {
			log.Println("No swipes found even after resetting.")
			return nil, nil
		}
	}

	for _, s := range swipes {
		_, err := r.swipeRepo.UpdateOne(ctx,
			bson.M{"client_id": s.ClientId, "therapist_id": s.TherapistId},
			bson.M{"$set": bson.M{"status": "swiped"}},
		)
		if err != nil {
			log.Println("Failed to update swipe to swiped for", s.ClientId, s.TherapistId, err)
			return nil, errors.New("failed to update swipe status")
		}
	}

	var therapistResponse []model.Therapist
	for _, s := range swipes {
		profile, err := r.therapistRepo.FindTherapistMatchingProfile(s.TherapistId)
		if err != nil {
			log.Printf("Failed to find therapist %s: %v", s.TherapistId, err)
			continue
		}
		therapistResponse = append(therapistResponse, *profile)
	}

	return therapistResponse, nil
}

func (r *SwipeRepository) getTop5PendingSwipes(clientId string) ([]model.ClientSwipe, error) {
	ctx := context.Background()

	cursor, err := r.swipeRepo.Find(ctx,
		bson.M{"client_id": clientId, "status": "pending"},
		optsFindTop5(),
	)
	if err != nil {
		log.Println("Failed to find top 5 swipes:", err)
		return nil, errors.New("failed to find swipes")
	}

	var swipes []model.ClientSwipe
	if err := cursor.All(ctx, &swipes); err != nil {
		return nil, err
	}
	return swipes, nil
}

func (r *SwipeRepository) UpdateAllPendingTherapistSwipeProfile(clientId string) error {
	ctx := context.Background()
	if clientId == "" {
		return errors.New("clientId is empty")
	}

	filter := bson.D{{Key: "client_id", Value: clientId}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: "pending"}}}}

	result, err := r.swipeRepo.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println("Failed to update therapist status:", err)
		return err
	}
	if result.ModifiedCount == 0 {
		log.Println("No swipe documents updated")
	}
	return nil
}

func (r *SwipeRepository) GetSwipedTherapistIds(clientId string) ([]string, error) {
	ctx := context.Background()
	cursor, err := r.swipeRepo.Find(ctx,
		bson.M{"client_id": clientId, "status": "swiped"},
		options.Find().SetProjection(bson.M{"therapist_id": 1}),
	)
	if err != nil {
		log.Println("Failed to get swiped therapist ids:", err)
		return nil, errors.New("failed to get swiped therapist ids")
	}

	var ids []string
	for cursor.Next(ctx) {
		var swipe model.ClientSwipe
		if err := cursor.Decode(&swipe); err == nil {
			ids = append(ids, swipe.TherapistId)
		}
	}
	return ids, nil
}

func (r *SwipeRepository) FilterAllTherapist(clientId string) ([]model.ClientSwipe, error) {
	client, err := r.clientRepo.FindClientMatchingProfile(clientId)
	if err != nil {
		log.Println("Failed to find client profile", err)
		return nil, errors.New("failed to find client profile")
	}

	therapists, err := r.therapistRepo.FindAllTherapistMatchingProfiles()
	if err != nil {
		log.Println("Failed to find therapist profile", err)
		return nil, errors.New("failed to find therapist profile")
	}

	raw, err := AppendDataIntoSwipe(client, therapists)
	if err != nil {
		log.Println("Failed to append swipe data", err)
		return nil, errors.New("failed to append swipe data")
	}

	swipes, err := external.RecommendationApi(raw)
	if err != nil {
		log.Println("Recommendation API error:", err)
		return nil, errors.New("failed to fetch recommendation")
	}

	swipedIds, err := r.GetSwipedTherapistIds(clientId)
	if err != nil {
		log.Println("Failed to get swiped therapist IDs", err)
		return nil, errors.New("failed to get swiped therapist IDs")
	}

	bf := bloom.NewWithEstimates(uint(len(swipedIds)), 0.001)
	for _, id := range swipedIds {
		bf.Add([]byte(id))
	}

	var filtered []model.ClientSwipe
	for _, s := range swipes {
		if !bf.Test([]byte(s.TherapistId)) {
			filtered = append(filtered, s)
		}
	}

	err = r.InsertSwipes(clientId, filtered)
	if err != nil {
		log.Println("Failed to save swipes", err)
		return nil, err
	}

	return filtered, nil
}

func AppendDataIntoSwipe(client *model.Client, therapist []model.Therapist) (dto.FilterRawData, error) {
	return dto.FilterRawData{
		ClientRaw:    *client,
		TherapistRaw: therapist,
	}, nil
}

func optsFindTop5() *options.FindOptions {
	return options.Find().SetSort(bson.M{"points": -1}).SetLimit(5)
}

func (r *SwipeRepository) PopTop5SwipesV1(clientId string) ([]model.TherapistV1, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("psyconnect:swipe:top5:%s", clientId)

	var cached []model.TherapistV1
	if err := r.redis.Get(ctx, cacheKey, &cached); err == nil && len(cached) > 0 {
		log.Println("PopTop5SwipesV1: cache hit")
		return cached, nil
	}

	swipes, err := r.getTop5PendingSwipes(clientId)
	if err != nil {
		return nil, err
	}

	if len(swipes) == 0 {
		if err := r.UpdateAllPendingTherapistSwipeProfile(clientId); err != nil {
			return nil, err
		}
		swipes, err = r.getTop5PendingSwipes(clientId)
		if err != nil {
			return nil, err
		}
		if len(swipes) == 0 {
			return nil, nil
		}
	}

	for _, s := range swipes {
		_, _ = r.swipeRepo.UpdateOne(ctx,
			bson.M{"client_id": s.ClientId, "therapist_id": s.TherapistId},
			bson.M{"$set": bson.M{"status": "swiped"}},
		)
	}

	var therapistResponse []model.TherapistV1
	for _, s := range swipes {
		profile, err := r.therapistRepo.FindTherapistMatchingProfileV1(s.TherapistId)
		if err != nil {
			continue
		}
		therapistResponse = append(therapistResponse, *profile)
	}

	_ = r.redis.Set(ctx, cacheKey, therapistResponse)
	return therapistResponse, nil
}

func (r *SwipeRepository) FilterAllTherapistV1(clientId string) ([]model.ClientSwipeV1, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("psyconnect:swipe:filter:%s", clientId)

	var cached []model.ClientSwipeV1
	if err := r.redis.Get(ctx, cacheKey, &cached); err == nil && len(cached) > 0 {
		log.Println("FilterAllTherapistV1: cache hit")
		return cached, nil
	}

	client, err := r.clientRepo.FindClientMatchingProfile(clientId)
	if err != nil {
		return nil, err
	}

	therapists, err := r.therapistRepo.FindAllTherapistMatchingProfilesV1()
	if err != nil {
		return nil, err
	}

	raw, _ := AppendDataIntoSwipeV1(client, therapists)
	swipes, err := external.RecommendationApiV1(raw)
	if err != nil {
		return nil, err
	}

	swipedIds, _ := r.GetSwipedTherapistIds(clientId)
	bf := bloom.NewWithEstimates(uint(len(swipedIds)), 0.001)
	for _, id := range swipedIds {
		bf.Add([]byte(id))
	}

	var filtered []model.ClientSwipeV1
	for _, s := range swipes {
		if !bf.Test([]byte(s.TherapistId)) {
			filtered = append(filtered, s)
		}
	}

	_ = r.InsertSwipesV1(clientId, filtered)
	_ = r.redis.Set(ctx, cacheKey, filtered)
	return filtered, nil
}

func (r *SwipeRepository) InsertSwipesV1(clientId string, swipes []model.ClientSwipeV1) error {
	ctx := context.Background()
	if len(swipes) == 0 {
		return nil
	}
	var docs []interface{}
	for _, s := range swipes {
		docs = append(docs, s)
	}
	_, err := r.swipeRepo.InsertMany(ctx, docs)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("psyconnect:swipe:filter:%s", clientId)
	_ = r.redis.Delete(ctx, cacheKey)
	return nil
}

func AppendDataIntoSwipeV1(client *model.Client, therapist []model.TherapistV1) (dto.FilterRawDataV1, error) {
	return dto.FilterRawDataV1{
		ClientRaw:    *client,
		TherapistRaw: therapist,
	}, nil
}
