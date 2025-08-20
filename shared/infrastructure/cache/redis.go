package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
	PoolSize int
}

func NewRedisClient(config Config) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
		PoolSize: config.PoolSize,
	})

	ctx := context.Background()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisClient{
		client: client,
		ctx:    ctx,
	}, nil
}

func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return r.client.Set(r.ctx, key, data, expiration).Err()
}

func (r *RedisClient) Get(key string, dest interface{}) error {
	data, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key not found: %s", key)
		}
		return fmt.Errorf("failed to get value: %w", err)
	}

	if err := json.Unmarshal([]byte(data), dest); err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return nil
}

func (r *RedisClient) Delete(keys ...string) error {
	return r.client.Del(r.ctx, keys...).Err()
}

func (r *RedisClient) Exists(key string) (bool, error) {
	result, err := r.client.Exists(r.ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

func (r *RedisClient) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal value: %w", err)
	}

	return r.client.SetNX(r.ctx, key, data, expiration).Result()
}

func (r *RedisClient) Expire(key string, expiration time.Duration) error {
	return r.client.Expire(r.ctx, key, expiration).Err()
}

func (r *RedisClient) TTL(key string) (time.Duration, error) {
	return r.client.TTL(r.ctx, key).Result()
}

func (r *RedisClient) Increment(key string) (int64, error) {
	return r.client.Incr(r.ctx, key).Result()
}

func (r *RedisClient) Decrement(key string) (int64, error) {
	return r.client.Decr(r.ctx, key).Result()
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}

func (r *RedisClient) FlushAll() error {
	return r.client.FlushAll(r.ctx).Err()
}

func (r *RedisClient) GetClient() *redis.Client {
	return r.client
}
