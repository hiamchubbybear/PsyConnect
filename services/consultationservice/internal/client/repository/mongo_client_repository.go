package repository

import (
	"consultationservice/internal/client/domain"
	"consultationservice/internal/redis"
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClientRepository struct {
	collection *mongo.Collection
	redis      redis.RedisStore
}

func NewMongoClientRepository(collection *mongo.Collection, redisStore redis.RedisStore) *MongoClientRepository {
	return &MongoClientRepository{
		collection: collection,
		redis:      redisStore,
	}
}

func (r *MongoClientRepository) redisKey(profileID string) string {
	return redis.NewKeyBuilder("psyconnect").Build("consultation", "client", profileID)
}

// domainToMongo converts domain entity to MongoDB document
func (r *MongoClientRepository) domainToMongo(client *domain.Client) bson.M {
	return bson.M{
		"profile_id":         client.ProfileID,
		"address":            client.Address,
		"languages":          client.Languages,
		"issue_detail":       client.IssueDetail,
		"consultation_modes": client.ConsultationModes,
		"rage_price":         client.RangePrice,
		"availability": bson.M{
			"days":       client.Availability.Days,
			"time_slots": client.Availability.TimeSlots,
		},
		"preferred_therapist_gender":   client.PreferredTherapistGender,
		"experience_level":             client.ExperienceLevel,
		"therapist_specialization":     client.TherapistSpecialization,
		"urgency_level":                client.UrgencyLevel,
		"session_duration":             client.SessionDuration,
		"preferred_therapist_language": client.PreferredTherapistLanguage,
		"is_flexible_with_schedule":    client.IsFlexibleWithSchedule,
		"created_at":                   client.CreatedAt,
		"updated_at":                   client.UpdatedAt,
	}
}

// mongoToDomain converts MongoDB document to domain entity
func (r *MongoClientRepository) mongoToDomain(doc bson.M) *domain.Client {
	client := &domain.Client{
		ProfileID:                  getString(doc, "profile_id"),
		Address:                    getString(doc, "address"),
		Languages:                  getStringArray(doc, "languages"),
		IssueDetail:                getStringArray(doc, "issue_detail"),
		ConsultationModes:          getStringArray(doc, "consultation_modes"),
		RangePrice:                 getInt(doc, "rage_price"),
		PreferredTherapistGender:   getString(doc, "preferred_therapist_gender"),
		ExperienceLevel:            getString(doc, "experience_level"),
		TherapistSpecialization:    getStringArray(doc, "therapist_specialization"),
		UrgencyLevel:               getString(doc, "urgency_level"),
		SessionDuration:            getInt(doc, "session_duration"),
		PreferredTherapistLanguage: getStringArray(doc, "preferred_therapist_language"),
		IsFlexibleWithSchedule:     getBool(doc, "is_flexible_with_schedule"),
	}

	if avail, ok := doc["availability"].(bson.M); ok {
		client.Availability = domain.Availability{
			Days:      getStringArray(avail, "days"),
			TimeSlots: getStringArray(avail, "time_slots"),
		}
	}

	return client
}

func (r *MongoClientRepository) Create(ctx context.Context, client *domain.Client) error {
	// Check if client already exists
	filter := bson.M{"profile_id": client.ProfileID}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("Error checking if client exists: %v", err)
		return errors.New("failed to check if client exists")
	}
	if count > 0 {
		log.Printf("Client already exists with profile_id: %s", client.ProfileID)
		return errors.New("client with this profile already exists")
	}

	// Insert new client
	doc := r.domainToMongo(client)
	_, err = r.collection.InsertOne(ctx, doc)
	if err != nil {
		log.Printf("Failed to insert client: %v", err)
		return errors.New("failed to insert client")
	}

	// Invalidate cache
	_ = r.redis.Delete(ctx, r.redisKey(client.ProfileID))

	return nil
}

func (r *MongoClientRepository) GetByProfileID(ctx context.Context, profileID string) (*domain.Client, error) {
	key := r.redisKey(profileID)

	// Try cache first
	var cachedDoc bson.M
	if err := r.redis.Get(ctx, key, &cachedDoc); err == nil {
		log.Printf("Cache hit for client: %s", profileID)
		return r.mongoToDomain(cachedDoc), nil
	}

	// Query database
	var doc bson.M
	err := r.collection.FindOne(ctx, bson.M{"profile_id": profileID}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("client not found")
		}
		return nil, err
	}

	// Cache result
	_ = r.redis.Set(ctx, key, doc)

	return r.mongoToDomain(doc), nil
}

func (r *MongoClientRepository) GetAll(ctx context.Context) ([]*domain.Client, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.New("failed to find all clients")
	}
	defer cursor.Close(ctx)

	var results []*domain.Client
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			log.Printf("Failed to decode client: %v", err)
			continue
		}
		results = append(results, r.mongoToDomain(doc))
	}

	if err := cursor.Err(); err != nil {
		return nil, errors.New("cursor error")
	}

	return results, nil
}

func (r *MongoClientRepository) Update(ctx context.Context, profileID string, client *domain.Client) error {
	client.ProfileID = profileID
	client.Update()

	update := bson.M{"$set": r.domainToMongo(client)}

	result := r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"profile_id": profileID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return errors.New("client not found")
		}
		return errors.New("failed to update client")
	}

	// Invalidate cache
	_ = r.redis.Delete(ctx, r.redisKey(profileID))

	log.Printf("Client updated: %s", profileID)
	return nil
}

func (r *MongoClientRepository) Delete(ctx context.Context, profileID string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"profile_id": profileID})
	if err != nil {
		return errors.New("failed to delete client")
	}

	if result.DeletedCount == 0 {
		return errors.New("client not found")
	}

	// Invalidate cache
	_ = r.redis.Delete(ctx, r.redisKey(profileID))

	return nil
}

// Helper functions
func getString(doc bson.M, key string) string {
	if val, ok := doc[key].(string); ok {
		return val
	}
	return ""
}

func getInt(doc bson.M, key string) int {
	if val, ok := doc[key].(int32); ok {
		return int(val)
	}
	if val, ok := doc[key].(int64); ok {
		return int(val)
	}
	if val, ok := doc[key].(int); ok {
		return val
	}
	return 0
}

func getBool(doc bson.M, key string) bool {
	if val, ok := doc[key].(bool); ok {
		return val
	}
	return false
}

func getStringArray(doc bson.M, key string) []string {
	if val, ok := doc[key].(bson.A); ok {
		result := make([]string, 0, len(val))
		for _, v := range val {
			if str, ok := v.(string); ok {
				result = append(result, str)
			}
		}
		return result
	}
	return []string{}
}
