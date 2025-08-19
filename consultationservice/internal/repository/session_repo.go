package repository

import (
	"consultationservice/internal/dto"
	"consultationservice/internal/model"
	"consultationservice/internal/redis"
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SessionRepository struct {
	MongoDBCollection *mongo.Collection
	clientRepo        *ClientRepository
	therapistRepo     *TherapistRepository
	matchingRepo      *MatchRepository
	redis             redis.RedisStore
}

func NewSessionRepository(
	collection *mongo.Collection,
	clientRepo *ClientRepository,
	therapistRepo *TherapistRepository,
	matchingRepo *MatchRepository,
	redisStore redis.RedisStore,
) *SessionRepository {
	return &SessionRepository{
		clientRepo:        clientRepo,
		therapistRepo:     therapistRepo,
		MongoDBCollection: collection,
		matchingRepo:      matchingRepo,
		redis:             redisStore,
	}
}

func (r *SessionRepository) CreateNewSession(session dto.SessionRequest) (interface{}, error) {
	ctx := context.Background()
	clientId := session.ClientID
	therapistId := session.TherapistID

	_, terr := r.therapistRepo.FindTherapistMatchingProfile(therapistId)
	if terr != nil {
		log.Printf("Failed to find therapist profile : %v", terr)
		return nil, terr
	}

	client, cerr := r.clientRepo.FindClientMatchingProfile(clientId)
	if cerr != nil || client == nil || client.ProfileId == "" {
		log.Println("Client not found or invalid")
		return nil, errors.New("cannot find client")
	}

	var clientDoc struct {
		CurrentSession []string `bson:"current_session"`
	}
	if err := r.clientRepo.MongoDBCollection.FindOne(ctx, bson.M{"profile_id": clientId}).Decode(&clientDoc); err != nil {
		return nil, errors.New("cannot find client sessions")
	}

	var therapistDoc struct {
		CurrentSession []string `bson:"current_session"`
	}
	if err := r.therapistRepo.MongoDBCollection.FindOne(ctx, bson.M{"profile_id": therapistId}).Decode(&therapistDoc); err != nil {
		return nil, errors.New("cannot find therapist sessions")
	}

	allSessionIDs := append(clientDoc.CurrentSession, therapistDoc.CurrentSession...)
	var sessions []model.Session
	if len(allSessionIDs) > 0 {
		cursor, err := r.MongoDBCollection.Find(ctx, bson.M{"_id": bson.M{"$in": allSessionIDs}})
		if err != nil {
			return nil, errors.New("failed to find sessions")
		}
		if err := cursor.All(ctx, &sessions); err != nil {
			return nil, errors.New("failed to decode sessions")
		}
	}

	overlap, err := r.matchingRepo.CheckTimeOverlap(session.SessionTime.Day.String(), session.SessionTime.StartTime, session.SessionTime.EndTime, sessions)
	if err != nil {
		return nil, err
	}
	if overlap {
		return nil, errors.New("session time overlaps with existing sessions")
	}

	res, err := r.MongoDBCollection.InsertOne(ctx, session)
	if err != nil {
		return nil, err
	}
	sessionID := res.InsertedID.(primitive.ObjectID).Hex()

	update := bson.D{{Key: "$push", Value: bson.D{{Key: "current_session", Value: res.InsertedID}}}}
	_, _ = r.clientRepo.MongoDBCollection.UpdateOne(ctx, bson.M{"profile_id": clientId}, update)
	_, _ = r.therapistRepo.MongoDBCollection.UpdateOne(ctx, bson.M{"profile_id": therapistId}, update)

	_ = r.redis.Set(ctx, "session:"+sessionID, session)

	return res, nil
}

func (r *SessionRepository) DeleteCurrentSession(session dto.DeleteSessionRequest) (bool, error) {
	ctx := context.Background()

	if session.SessionID == "" {
		return false, errors.New("session id couldn't be empty")
	}
	objectID, err := primitive.ObjectIDFromHex(session.SessionID)
	if err != nil {
		return false, errors.New("invalid session id format")
	}

	res, err := r.MongoDBCollection.DeleteOne(ctx, bson.D{{Key: "_id", Value: objectID}})
	if err != nil {
		log.Printf("Failed to delete session: %v", err)
		return false, err
	}

	if res.DeletedCount == 0 {
		log.Println("Session not found")
		return false, errors.New("session not found")
	}

	therapistUpdate := bson.D{{Key: "$pull", Value: bson.D{{Key: "current_session", Value: objectID}}}}
	_, err = r.therapistRepo.MongoDBCollection.UpdateMany(ctx, bson.D{}, therapistUpdate)
	if err != nil {
		log.Printf("Failed to update therapist: %v", err)
		return false, err
	}

	clientUpdate := bson.D{{Key: "$pull", Value: bson.D{{Key: "current_session", Value: objectID}}}}
	_, err = r.clientRepo.MongoDBCollection.UpdateMany(ctx, bson.D{}, clientUpdate)
	if err != nil {
		log.Printf("Failed to update client: %v", err)
		return false, err
	}

	_ = r.redis.Delete(ctx, "session:"+session.SessionID)

	return true, nil
}

func (r *SessionRepository) GetAllSessions() ([]model.Session, error) {
	var sessions []model.Session
	cursor, err := r.MongoDBCollection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(context.Background(), &sessions); err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *SessionRepository) GetSessionByID(id string) (*dto.SessionRequest, error) {
	ctx := context.Background()
	cacheKey := "session:" + id

	var cached dto.SessionRequest
	if err := r.redis.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	var session dto.SessionRequest
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Failed to convert id to hex")
		return nil, err
	}
	err = r.MongoDBCollection.FindOne(ctx, bson.M{"_id": oid}).Decode(&session)
	if err != nil {
		log.Printf("Failed to find session : %v", err)
		return nil, err
	}

	_ = r.redis.Set(ctx, cacheKey, session)

	return &session, nil
}

func (r *SessionRepository) FindAllSessionByProfileId(profileId string, requester string) (*[]dto.SessionResponse, error) {
	var sessions []dto.SessionResponse
	if requester == "role.client" {
		requester = "client_id"
	} else {
		requester = "therapist_id"
	}
	cursor, err := r.MongoDBCollection.Find(context.Background(), bson.M{requester: profileId})
	if err != nil {
		log.Printf("Failed find ")
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var session dto.SessionResponse
		if err := cursor.Decode(&session); err != nil {
			log.Printf("Decode error: %v", err)
			continue
		}
		sessions = append(sessions, session)
	}
	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, err
	}
	return &sessions, nil
}

func (s *SessionRepository) SetMatchingRepository(matchingRepo *MatchRepository) {
	s.matchingRepo = matchingRepo
}
