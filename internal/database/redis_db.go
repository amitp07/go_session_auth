package database

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func SetupRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "go_session_store",
		DB:       0,
	})

	if rdb.Ping(context.TODO()).Err() != nil {
		panic("Could not connect to Redis..")
	}

	fmt.Printf("Msg: %s\n", "Connection to Redis is successful")

	return rdb
}
