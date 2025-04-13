package database

import (
	"context"
	"fmt"
	"session-auth/internal/utils"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	db *redis.Client
}

func SetupRedis() *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "go_session_store",
		DB:       0,
	})

	if rdb.Ping(context.TODO()).Err() != nil {
		panic("Could not connect to Redis..")
	}

	fmt.Printf("Msg: %s\n", "Connection to Redis is successful")

	return &RedisClient{db: rdb}
}

func (r *RedisClient) SetSession(val string) (string, error) {
	// generate random uuid for session key
	key, err := utils.GenerateRadomId(16)

	if err != nil {
		return "", err
	}

	key = "session:" + key

	if err := r.db.Set(context.TODO(), key, val, -1).Err(); err != nil {
		fmt.Println("Error in setting session to redis", err.Error())
		return "", nil
	}

	return key, nil
}

func (r *RedisClient) SetOtp(val []byte) (string, error) {
	key := "mfa:" + uuid.NewString()
	err := r.db.Set(context.TODO(), key, val, 5*time.Minute).Err()

	if err != nil {
		return "", err
	}

	return key, nil
}

func (r *RedisClient) Get(key string) string {
	res, err := r.db.Get(context.TODO(), key).Result()

	if err != nil {
		return ""
	}

	return res
}

func (r *RedisClient) Delete(key string) error {
	return r.db.Del(context.TODO(), key).Err()
}
