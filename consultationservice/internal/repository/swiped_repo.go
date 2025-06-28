package repository

import (
	"consultationservice/internal/dto"
	"consultationservice/internal/model"
	"consultationservice/internal/repository/infrastructure/external"
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

type SwipeRepository struct {
	clientRepo      *ClientRepository
	therapistRepo   *TherapistRepository
	sessionRepo     *SessionRepository
	swipeRepository *mongo.Collection
}

func NewSwipeRepository(clientRepo *ClientRepository, therapistRepo *TherapistRepository, sessionRepo *SessionRepository, swipeRepo *mongo.Collection) *SwipeRepository {
	return &SwipeRepository{
		clientRepo:      clientRepo,
		therapistRepo:   therapistRepo,
		sessionRepo:     sessionRepo,
		swipeRepository: swipeRepo,
	}
}

func (r *SwipeRepository) FilterAllTherapist(profileId string, clientId string) (*model.ClientSwipes, error) {
	client, err := r.clientRepo.FindClientMatchingProfile(clientId)
	if err != nil {
		log.Printf("Failed to find client profile", err)
		return nil, errors.New("Failed to find client profile")
	}
	therapists, err := r.therapistRepo.FindAllTherapistMatchingProfiles()
	if err != nil {
		log.Printf("Failed to find therapist profile", err)
		return nil, errors.New("Failed to find therapist profile")
	}
	swipes, err := AppendDataIntoSwipe(client, therapists)
	if err != nil {
		log.Printf("Failed to find therapist profile", err)
		return nil, errors.New("Failed to append swipe data")
	}
	response, err := external.RecommendationApi(swipes)
	if err != nil {
		log.Printf("Failed to fetch recommend")
		return nil, errors.New("Failed to fetch recommend")
	}
	return nil, nil
}
func (r *SwipeRepository) SaveFetchedClient(swipesData model.ClientSwipes) (bool, error) {
	ctx := context.TODO()
	_, err := r.swipeRepository.InsertOne(ctx, swipesData)
	if err != nil {
		log.Print("Failed to insert swipe data after ")
		return false, error.New("Failed to insert swipe data after ")
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
