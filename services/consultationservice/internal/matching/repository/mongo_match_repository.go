package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"consultationservice/internal/matching/domain"
	"consultationservice/internal/redis"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoMatchRepository struct {
	collection *mongo.Collection
	redis      redis.RedisStore
}

func NewMongoMatchRepository(collection *mongo.Collection, redis redis.RedisStore) *MongoMatchRepository {
	return &MongoMatchRepository{
		collection: collection,
		redis:      redis,
	}
}

func (r *MongoMatchRepository) keyMatchesByClient(clientId string) string {
	return redis.NewKeyBuilder("psyconnect").Build("consultation", "match", "by_client", clientId)
}

func (r *MongoMatchRepository) keyMatchesByClientPage(clientId string, page int64) string {
	return redis.NewKeyBuilder("psyconnect").Build("consultation", "match", "by_client", clientId, "page", fmt.Sprint(page))
}

func (r *MongoMatchRepository) Create(ctx context.Context, match *domain.Match) error {
	match.MatchedAt = time.Now().UTC()
	match.UpdatedAt = time.Now().UTC()

	_, err := r.collection.InsertOne(ctx, match)
	if err != nil {
		log.Printf("Failed to insert match: %v", err)
		return errors.New("failed to insert match")
	}

	// Invalidate cache
	_ = r.redis.Delete(ctx, r.keyMatchesByClient(match.ClientID))
	_ = r.redis.Delete(ctx, r.keyMatchesByClientPage(match.ClientID, 1)) // Invalidate first page

	return nil
}

func (r *MongoMatchRepository) GetMatchesByClientID(ctx context.Context, clientID string) ([]domain.Match, error) {
	cacheKey := r.keyMatchesByClient(clientID)
	var cached []domain.Match
	if err := r.redis.Get(ctx, cacheKey, &cached); err == nil {
		return cached, nil
	}

	cursor, err := r.collection.Find(ctx, bson.M{"client_id": clientID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var matches []domain.Match
	if err := cursor.All(ctx, &matches); err != nil {
		return nil, err
	}

	_ = r.redis.Set(ctx, cacheKey, matches)
	return matches, nil
}

func (r *MongoMatchRepository) GetMatchesByClientIDPaginated(ctx context.Context, clientID string, page int64, limit int64) ([]domain.Match, error) {
	if page < 1 {
		page = 1
	}
	skip := (page - 1) * limit

	cacheKey := r.keyMatchesByClientPage(clientID, page)
	var cached []domain.Match
	if err := r.redis.Get(ctx, cacheKey, &cached); err == nil {
		return cached, nil
	}

	opts := options.Find().
		SetLimit(limit).
		SetSkip(skip).
		SetSort(bson.M{"matched_at": -1})

	cursor, err := r.collection.Find(ctx, bson.M{"client_id": clientID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var matches []domain.Match
	if err := cursor.All(ctx, &matches); err != nil {
		return nil, err
	}

	_ = r.redis.Set(ctx, cacheKey, matches)
	return matches, nil
}
