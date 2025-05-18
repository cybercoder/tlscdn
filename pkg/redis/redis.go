package redis

import (
	"os"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func CreateClient() *redis.Client {
	if redisClient != nil {
		return redisClient
	}
	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
}
