package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"time"
)

var rdb *redis.Client

func Init() {
	rdb = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST") + os.Getenv("REDIS_PORT"),
		DB:   0,
	})

	ctx := context.Background()
	err := rdb.Ping(ctx).Err()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	} else {
		fmt.Println("Connected to Redis successfully")
	}
}

func SetCache(ctx context.Context, key string, value string, expiration time.Duration) error {
	err := rdb.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("could not set cache: %v", err)
	}
	return nil
}

func GetCache(ctx context.Context, key string) (string, error) {
	val, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key does not exist")
	} else if err != nil {
		return "", fmt.Errorf("could not get cache: %v", err)
	}
	return val, nil
}
