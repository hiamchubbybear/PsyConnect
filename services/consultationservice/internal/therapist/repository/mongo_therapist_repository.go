package repository

import (
	"consultationservice/internal/redis"
	"consultationservice/internal/therapist/domain"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoTherapistRepository struct {
	collection *mongo.Collection
	redis      redis.RedisStore
}

func NewMongoTherapistRepository(collection *mongo.Collection, redisStore redis.RedisStore) *MongoTherapistRepository {
	return &MongoTherapistRepository{
		collection: collection,
		redis:      redisStore,
	}
}

func (r *MongoTherapistRepository) redisKey(profileID string) string {
	return fmt.Sprintf("psyconnect:therapist:profile:%s", profileID)
}

func (r *MongoTherapistRepository) redisKeyAll() string {
	return "psyconnect:therapist:all"
}


func (r *MongoTherapistRepository) domainToMongo(therapist *domain.Therapist) bson.M {
	doc := bson.M{
		"profile_id":         therapist.ProfileID,
		"address":            therapist.Address,
		"languages":          therapist.Languages,
		"specialization":     therapist.Specialization,
		"consultation_modes": therapist.ConsultationModes,
		"experience":         therapist.Experience,
		"rating":             therapist.Rating,
		"currency":           therapist.Currency,
		"rage_price":         therapist.RagePrice,
		"is_available":       therapist.IsAvailable,
		"availability": bson.M{
			"days":       therapist.Availability.Days,
			"time_slots": therapist.Availability.TimeSlots,
		},
		"current_session": therapist.CurrentSession,
		"matched_clients": therapist.MatchedClients,
		"avatar_override": therapist.AvatarOverride,
		"name":            therapist.Name,
		"created_at":      therapist.CreatedAt,
		"updated_at":      therapist.UpdatedAt,
	}

	
	profInfo := bson.M{
		"title": bson.M{
			"code":    therapist.ProfessionalInfo.Title.Code,
			"display": therapist.ProfessionalInfo.Title.Display,
		},
		"experience_years": therapist.ProfessionalInfo.ExperienceYears,
	}

	var degrees []bson.M
	for _, d := range therapist.ProfessionalInfo.Degrees {
		degrees = append(degrees, bson.M{
			"type":        d.Type,
			"field":       d.Field,
			"institution": d.Institution,
			"year":        d.Year,
		})
	}
	profInfo["degrees"] = degrees

	var certifications []bson.M
	for _, c := range therapist.ProfessionalInfo.Certifications {
		certifications = append(certifications, bson.M{
			"name":   c.Name,
			"issuer": c.Issuer,
			"year":   c.Year,
		})
	}
	profInfo["certifications"] = certifications

	doc["professional_info"] = profInfo

	return doc
}


func (r *MongoTherapistRepository) mongoToDomain(doc bson.M) *domain.Therapist {
	t := &domain.Therapist{
		ProfileID:         getString(doc, "profile_id"),
		Address:           getString(doc, "address"),
		Languages:         getStringArray(doc, "languages"),
		Specialization:    getStringArray(doc, "specialization"),
		ConsultationModes: getStringArray(doc, "consultation_modes"),
		Experience:        getInt(doc, "experience"),
		Rating:            getFloat(doc, "rating"),
		Currency:          getString(doc, "currency"),
		RagePrice:         getInt(doc, "rage_price"),
		IsAvailable:       getBool(doc, "is_available"),
		CurrentSession:    getStringArray(doc, "current_session"),
		MatchedClients:    getStringArray(doc, "matched_clients"),
		AvatarOverride:    getString(doc, "avatar_override"),
		Name:              getString(doc, "name"),
	}

	if avail, ok := doc["availability"].(bson.M); ok {
		t.Availability = domain.Availability{
			Days:      getStringArray(avail, "days"),
			TimeSlots: getStringArray(avail, "time_slots"),
		}
	}

	if prof, ok := doc["professional_info"].(bson.M); ok {
		info := domain.ProfessionalInfo{
			ExperienceYears: getInt(prof, "experience_years"),
		}

		if title, ok := prof["title"].(bson.M); ok {
			info.Title = domain.ProfessionalTitle{
				Code:    getString(title, "code"),
				Display: getString(title, "display"),
			}
		}

		if degrees, ok := prof["degrees"].(bson.A); ok {
			for _, d := range degrees {
				if dm, ok := d.(bson.M); ok {
					info.Degrees = append(info.Degrees, domain.Degree{
						Type:        getString(dm, "type"),
						Field:       getString(dm, "field"),
						Institution: getString(dm, "institution"),
						Year:        getInt(dm, "year"),
					})
				}
			}
		}

		if certs, ok := prof["certifications"].(bson.A); ok {
			for _, c := range certs {
				if cm, ok := c.(bson.M); ok {
					info.Certifications = append(info.Certifications, domain.Certification{
						Name:   getString(cm, "name"),
						Issuer: getString(cm, "issuer"),
						Year:   getInt(cm, "year"),
					})
				}
			}
		}

		t.ProfessionalInfo = info
	}

	return t
}

func (r *MongoTherapistRepository) Create(ctx context.Context, therapist *domain.Therapist) error {
	filter := bson.M{"profile_id": therapist.ProfileID}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return errors.New("failed to check if therapist exists")
	}
	if count > 0 {
		return errors.New("therapist with this profile already exists")
	}

	doc := r.domainToMongo(therapist)
	_, err = r.collection.InsertOne(ctx, doc)
	if err != nil {
		return errors.New("failed to insert therapist")
	}

	_ = r.redis.Set(ctx, r.redisKey(therapist.ProfileID), doc)
	_ = r.redis.Delete(ctx, r.redisKeyAll())

	return nil
}

func (r *MongoTherapistRepository) GetByProfileID(ctx context.Context, profileID string) (*domain.Therapist, error) {
	key := r.redisKey(profileID)
	var cachedDoc bson.M
	if err := r.redis.Get(ctx, key, &cachedDoc); err == nil {
		return r.mongoToDomain(cachedDoc), nil
	}

	var doc bson.M
	err := r.collection.FindOne(ctx, bson.M{"profile_id": profileID}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("therapist not found")
		}
		return nil, err
	}

	_ = r.redis.Set(ctx, key, doc)
	return r.mongoToDomain(doc), nil
}

func (r *MongoTherapistRepository) GetAll(ctx context.Context) ([]*domain.Therapist, error) {
	key := r.redisKeyAll()
	var cachedDocs []bson.M
	
	if err := r.redis.Get(ctx, key, &cachedDocs); err == nil && len(cachedDocs) > 0 {
		var results []*domain.Therapist
		for _, doc := range cachedDocs {
			results = append(results, r.mongoToDomain(doc))
		}
		return results, nil
	}

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.New("failed to find all therapists")
	}
	defer cursor.Close(ctx)

	var results []*domain.Therapist
	var docs []bson.M

	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		docs = append(docs, doc)
		results = append(results, r.mongoToDomain(doc))
	}

	if len(docs) > 0 {
		_ = r.redis.Set(ctx, key, docs)
	}

	return results, nil
}

func (r *MongoTherapistRepository) Update(ctx context.Context, profileID string, therapist *domain.Therapist) error {
	therapist.ProfileID = profileID
	therapist.Update()

	update := bson.M{"$set": r.domainToMongo(therapist)}

	result := r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"profile_id": profileID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	if result.Err() != nil {
		return errors.New("failed to update therapist")
	}

	
	_ = r.redis.Delete(ctx, r.redisKey(profileID))
	_ = r.redis.Delete(ctx, r.redisKeyAll())

	return nil
}

func (r *MongoTherapistRepository) UpdateAvailability(ctx context.Context, profileID string, isAvailable bool) error {
	update := bson.M{"$set": bson.M{"is_available": isAvailable}}

	result, err := r.collection.UpdateOne(ctx, bson.M{"profile_id": profileID}, update)
	if err != nil {
		return errors.New("failed to update availability")
	}

	if result.MatchedCount == 0 {
		return errors.New("therapist not found")
	}

	_ = r.redis.Delete(ctx, r.redisKey(profileID))
	_ = r.redis.Delete(ctx, r.redisKeyAll())

	return nil
}

func (r *MongoTherapistRepository) Delete(ctx context.Context, profileID string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"profile_id": profileID})
	if err != nil {
		return errors.New("failed to delete therapist")
	}

	if result.DeletedCount == 0 {
		return errors.New("therapist not found")
	}

	_ = r.redis.Delete(ctx, r.redisKey(profileID))
	_ = r.redis.Delete(ctx, r.redisKeyAll())

	return nil
}

func (r *MongoTherapistRepository) Search(ctx context.Context, query string, limit, skip int64) ([]*domain.Therapist, error) {
	filter := bson.M{
		"is_available": true, 
		"$or": []bson.M{
			{"name": bson.M{"$regex": query, "$options": "i"}},
			{"specialization": bson.M{"$in": []bson.M{{"$regex": query, "$options": "i"}}}},
			{"address": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	opts := options.Find().
		SetLimit(limit).
		SetSkip(skip).
		SetSort(bson.D{{Key: "rating", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to search therapists: %w", err)
	}
	defer cursor.Close(ctx)

	var results []*domain.Therapist
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		results = append(results, r.mongoToDomain(doc))
	}

	return results, nil
}


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
	if val, ok := doc[key].(float64); ok {
		return int(val)
	}
	return 0
}

func getFloat(doc bson.M, key string) float64 {
	if val, ok := doc[key].(float64); ok {
		return val
	}
	if val, ok := doc[key].(int32); ok {
		return float64(val)
	}
	if val, ok := doc[key].(int64); ok {
		return float64(val)
	}
	if val, ok := doc[key].(int); ok {
		return float64(val)
	}
	return 0.0
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
