package redis

import (
	"context"
	"os"
	"time"

	"github.com/cybercoder/tlscdn-controller/pkg/logger"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

// CreateClient returns a Redis client instance
// It uses a singleton pattern to reuse existing connections
func CreateClient() *redis.Client {
	if redisClient != nil {
		return redisClient
	}

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	logger.Infof("Initializing Redis client connection to %s:%s", redisHost, redisPort)
	
	redisClient = redis.NewClient(&redis.Options{
		Addr:         redisHost + ":" + redisPort,
		Password:     redisPassword,
		DB:           0,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		logger.Warnf("Failed to connect to Redis at %s:%s: %v", redisHost, redisPort, err)
	} else {
		logger.Info("Successfully connected to Redis")
	}
	
	return redisClient
}
