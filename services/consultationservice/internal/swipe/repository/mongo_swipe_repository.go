package repository

import (
	"context"
	"errors"
	"log"
	"time"

	matchDomain "consultationservice/internal/matching/domain"
	matchRepo "consultationservice/internal/matching/repository"
	"consultationservice/internal/redis"
	"consultationservice/internal/swipe/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoSwipeRepository struct {
	collection *mongo.Collection
	matchRepo  matchRepo.MatchRepository
	redis      redis.RedisStore
}

func NewMongoSwipeRepository(
	collection *mongo.Collection,
	matchRepo matchRepo.MatchRepository,
	redis redis.RedisStore,
) *MongoSwipeRepository {
	return &MongoSwipeRepository{
		collection: collection,
		matchRepo:  matchRepo,
		redis:      redis,
	}
}

func (r *MongoSwipeRepository) InsertSwipe(ctx context.Context, swipe *domain.Swipe) error {
	swipe.CreatedAt = time.Now().UTC()
	swipe.Status = domain.SwipeStatusPending

	_, err := r.collection.InsertOne(ctx, swipe)
	if err != nil {
		log.Printf("Failed to insert swipe: %v", err)
		return errors.New("failed to insert swipe")
	}

	return nil
}

func (r *MongoSwipeRepository) InsertSwipes(ctx context.Context, clientID string, swipes []*domain.Swipe) error {
	if len(swipes) == 0 {
		return nil
	}

	var docs []interface{}
	for _, s := range swipes {
		s.CreatedAt = time.Now().UTC()
		s.Status = domain.SwipeStatusPending
		docs = append(docs, s)
	}

	_, err := r.collection.InsertMany(ctx, docs)
	if err != nil {
		log.Printf("Failed to insert swipes: %v", err)
		return errors.New("failed to insert swipes")
	}

	return nil
}

func (r *MongoSwipeRepository) SwipeAndMatch(ctx context.Context, clientID, therapistID string, points float32, reasons []string) error {
	// Delete existing swipe
	_, err := r.collection.DeleteOne(ctx, bson.M{
		"client_id":    clientID,
		"therapist_id": therapistID,
	})
	if err != nil {
		log.Println("Failed to delete existing swipe:", err)
		return errors.New("failed to delete old swipe")
	}

	// Create match using DDD MatchRepository
	match := matchDomain.NewMatch(clientID, therapistID, "swipe", float64(points), reasons)
	err = r.matchRepo.Create(ctx, match)
	if err != nil {
		log.Println("Failed to insert match:", err)
		return errors.New("failed to insert match")
	}

	log.Printf("Swipe matched: client %s -> therapist %s", clientID, therapistID)
	return nil
}

func (r *MongoSwipeRepository) GetTopSwipes(ctx context.Context, clientID string, limit int) ([]*domain.Swipe, error) {
	opts := options.Find().SetSort(bson.D{{Key: "points", Value: -1}}).SetLimit(int64(limit))
	filter := bson.M{
		"client_id": clientID,
		"status":    domain.SwipeStatusPending,
	}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		log.Printf("Failed to finding top swipes: %v", err)
		return nil, errors.New("failed to find top swipes")
	}
	defer cursor.Close(ctx)

	var swipes []*domain.Swipe
	if err = cursor.All(ctx, &swipes); err != nil {
		return nil, errors.New("failed to decode top swipes")
	}

	return swipes, nil
}

func (r *MongoSwipeRepository) DeleteSwipesByClient(ctx context.Context, clientID string) error {
	filter := bson.M{
		"client_id": clientID,
		"status":    domain.SwipeStatusPending,
	}
	_, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		log.Printf("Failed to delete pending swipes: %v", err)
		return errors.New("failed to delete pending swipes")
	}
	return nil
}
