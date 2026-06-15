package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(addr, password string, db int) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,

		PoolSize:     50,
		MinIdleConns: 10,

		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,

		MaxRetries:      3,
		MinRetryBackoff: 8 * time.Millisecond,
		MaxRetryBackoff: 512 * time.Millisecond,
	})

	maxRetries := 10
	retryInterval := 2 * time.Second

	var err error

	for i := 1; i <= maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

		err = rdb.Ping(ctx).Err()
		cancel()

		if err == nil {
			log.Println("Redis connected successfully")
			return &RedisClient{Client: rdb}, nil
		}

		log.Printf("Redis not ready (attempt %d/%d): %v\n", i, maxRetries, err)

		if i < maxRetries {
			time.Sleep(retryInterval)
		}
	}

	return nil, fmt.Errorf("failed to connect to Redis after %d attempts: %w", maxRetries, err)
}
