package repository

import (
	"consultationservice/internal/dto"
	"consultationservice/internal/model"
	"consultationservice/internal/repository/infrastructure/external"
	"context"
	"errors"
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
	// B1: Tìm thông tin client
	client, err := r.clientRepo.FindClientMatchingProfile(clientId)
	if err != nil {
		log.Println("Failed to find client profile", err)
		return nil, errors.New("failed to find client profile")
	}

	// B2: Tìm toàn bộ therapist matching
	therapists, err := r.therapistRepo.FindAllTherapistMatchingProfiles()
	if err != nil {
		log.Println("Failed to find therapist profile", err)
		return nil, errors.New("failed to find therapist profile")
	}

	// B3: Gộp data
	raw, err := AppendDataIntoSwipe(client, therapists)
	if err != nil {
		log.Println("Failed to append swipe data", err)
		return nil, errors.New("failed to append swipe data")
	}

	// B4: Gọi API gợi ý
	response, err := external.RecommendationApi(raw)
	if err != nil {
		log.Println(err)
		return nil, errors.New("failed to fetch recommendation")
	}

	// B5: Tạo Bloom filter từ therapist đã swipe
	existingSwipe := &model.ClientSwipes{}
	_ = r.swipeRepo.FindOne(
		context.Background(),
		bson.D{{Key: "client_id", Value: clientId}},
	).Decode(existingSwipe)

	swipedTherapistIds := make([]string, 0)
	for _, swipe := range existingSwipe.Swiped {
		swipedTherapistIds = append(swipedTherapistIds, swipe.TherapistId)
	}

	// B6: Khởi tạo Bloom Filter và add therapist đã swipe
	bf := bloom.NewWithEstimates(uint(len(swipedTherapistIds)), 0.001)
	for _, id := range swipedTherapistIds {
		bf.Add([]byte(id))
	}

	// B7: Lọc danh sách therapist được recommend nếu đã swipe
	var filteredSwipes []model.TherapistSwipe
	for _, s := range response.Swipes {
		if !bf.Test([]byte(s.TherapistId)) {
			filteredSwipes = append(filteredSwipes, s)
		}
	}

	// Gán lại danh sách đã lọc
	response.ClientId = clientId
	response.Swipes = filteredSwipes
	response.Swiped = swipedTherapistIdsToStruct(swipedTherapistIds) // preserve swiped

	// B8: Lưu vào MongoDB
	ok, err := r.SaveFetchedClient(response)
	if ok {
		return response, nil
	}
	return nil, err
}

func (r *SwipeRepository) SaveFetchedClient(swipesData *model.ClientSwipes) (bool, error) {
	clientId := swipesData.ClientId
	ctx := context.Background()

	// Replace document by client_id
	filter := bson.M{"client_id": clientId}
	_, err := r.swipeRepo.ReplaceOne(
		ctx,
		filter,
		swipesData,
		options.Replace().SetUpsert(true),
	)
	if err != nil {
		log.Println("Failed to replace swipe data:", err)
		return false, errors.New("failed to replace swipe data")
	}
	return true, nil
}

// Gộp client và therapist thành DTO gọi recommendation API
func AppendDataIntoSwipe(client *model.Client, therapist []*model.Therapist) (dto.FilterRawData, error) {
	var therapists []model.Therapist
	for _, t := range therapist {
		therapists = append(therapists, *t)
	}
	return dto.FilterRawData{
		ClientRaw:    *client,
		TherapistRaw: therapists,
	}, nil
}

// Chuyển list therapistId thành []TherapistSwipe dummy để lưu field "swiped"
func swipedTherapistIdsToStruct(ids []string) []model.TherapistSwipe {
	var result []model.TherapistSwipe
	for _, id := range ids {
		result = append(result, model.TherapistSwipe{
			TherapistId: id,
		})
	}
	return result
}
