package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(ctx context.Context, redisAddr string) (*RedisRepository, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	ctxPing, cancelPing := context.WithTimeout(ctx, 5*time.Second)
	defer cancelPing()

	if err := client.Ping(ctxPing).Err(); err != nil {
		return nil, fmt.Errorf("repository.NewRedisRepository, failed to ping redis: %w", err)
	}

	return &RedisRepository{client: client}, nil
}

func (r *RedisRepository) ListenQueue(ctx context.Context, queueName string) (string, error) {
	result, err := r.client.BLPop(ctx, 0, queueName).Result()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return "", fmt.Errorf("repository.ListenQueue, context is canceled")
		}
		return "", fmt.Errorf("repository.ListenQueue, failed to BLPop: %w", err)
	}

	if len(result) < 2 {
		return "", fmt.Errorf("repository.ListenQueue, empty result from redis")
	}
	return result[1], nil
}

func (r *RedisRepository) PushQueue(ctx context.Context, queueName string, message string) error {
	err := r.client.RPush(ctx, queueName, message).Err()
	if err != nil {
		return fmt.Errorf("repository.PushQueue, failed to RPush: %w", err)
	}
	return nil
}
