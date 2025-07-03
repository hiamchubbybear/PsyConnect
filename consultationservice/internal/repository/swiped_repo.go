package repository

import (
	"consultationservice/internal/dto"
	"consultationservice/internal/model"
	"consultationservice/internal/repository/infrastructure/external"
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SwipeRepository struct {
	clientRepo    *ClientRepository
	therapistRepo *TherapistRepository
	sessionRepo   *SessionRepository
	swipeRepo     *mongo.Collection
}

func NewSwipeRepository(clientRepo *ClientRepository, therapistRepo *TherapistRepository, sessionRepo *SessionRepository, swipeRepo *mongo.Collection) *SwipeRepository {
	return &SwipeRepository{
		clientRepo:    clientRepo,
		therapistRepo: therapistRepo,
		sessionRepo:   sessionRepo,
		swipeRepo:     swipeRepo,
	}
}

func (r *SwipeRepository) FilterAllTherapist(clientId string) (*model.ClientSwipes, error) {
	client, err := r.clientRepo.FindClientMatchingProfile(clientId)
	if err != nil {
		log.Println("Failed to find client profile", err)
		return nil, errors.New("Failed to find client profile")
	}
	therapists, err := r.therapistRepo.FindAllTherapistMatchingProfiles()
	if err != nil {
		log.Println("Failed to find therapist profile", err)
		return nil, errors.New("Failed to find therapist profile")
	}
	swipes, err := AppendDataIntoSwipe(client, therapists)
	if err != nil {
		log.Println("Failed to find therapist profile", err)
		return nil, errors.New("Failed to append swipe data")
	}
	response, err := external.RecommendationApi(swipes)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Failed to fetch recommend")
	}
	result, err := r.SaveFetchedClient(response)
	if result {
		return response, nil
	}
	return nil, err
}
func (r *SwipeRepository) SaveFetchedClient(swipesData *model.ClientSwipes) (bool, error) {
	*model.ClientSwipes checkRes
	clientId := swipesData.ClientId

	ctx := context.Background()
	_, err := r.swipeRepo.InsertOne(ctx, swipesData)
	r.swipeRepo.FindOne(
		context.Background(),
		bson.D{{Key: "profile_id", Value: clientId}},
	).Decode(&checkRes)
	if err != nil {
		log.Print("Failed to insert swipe data after ")
		return false, errors.New("Failed to insert swipe data after ")
	}
	return true, nil
}
func AppendDataIntoSwipe(client *model.Client, therapist []*model.Therapist) (dto.FilterRawData, error) {
	var therapists []model.Therapist
	for i := 0; i < len(therapist); i++ {
		therapists = append(therapists, *therapist[i])
	}
	clientSwipes := &dto.FilterRawData{
		ClientRaw:    *client,
		TherapistRaw: therapists,
	}
	return *clientSwipes, nil
}
