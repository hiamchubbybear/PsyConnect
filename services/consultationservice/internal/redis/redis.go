package redis

import (
	"consultationservice/bootstrap"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type RedisStore interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string, dest interface{}) error
	Delete(ctx context.Context, key string) error
	SAdd(ctx context.Context, key string, member string) error
	SIsMember(ctx context.Context, key string, member string) (bool, error)
	SMembers(ctx context.Context, key string) ([]string, error)
}

type redisStore struct {
	client *redis.Client
}

func NewRedisStore(env *bootstrap.Env) (RedisStore, error) {

	addr := env.ConsultationRedisAddress
	password := env.ConsultationRedisPassword
	db := env.ConsultationDB
	log.Println("Password" + password)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	if err := redisotel.InstrumentTracing(client); err != nil {
		log.Println(err)
	}
	if err := redisotel.InstrumentMetrics(client); err != nil {
		log.Println(err)
	}

	return &redisStore{client: client}, nil
}

func (r *redisStore) Set(ctx context.Context, key string, value interface{}) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, b, 30*time.Minute).Err()
}

func (r *redisStore) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(val, dest)
}

func (r *redisStore) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *redisStore) SAdd(ctx context.Context, key string, member string) error {
	return r.client.SAdd(ctx, key, member).Err()
}

func (r *redisStore) SIsMember(ctx context.Context, key string, member string) (bool, error) {
	return r.client.SIsMember(ctx, key, member).Result()
}

func (r *redisStore) SMembers(ctx context.Context, key string) ([]string, error) {
	return r.client.SMembers(ctx, key).Result()
}
